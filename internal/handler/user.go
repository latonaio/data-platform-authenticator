package handler

import (
	"errors"
	"net/http"
	"time"

	"data-platform-authenticator/internal/crypto"
	"data-platform-authenticator/internal/models"
	"data-platform-authenticator/pkg/response"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserParam struct {
	EmailAddress          string `json:"email_address" form:"email_address"`
	BusinessPartner       int    `json:"business_partner" form:"business_partner"`
	BusinessPartnerName   string `json:"business_partner_name" form:"business_partner_name"`
	Password              string `json:"password" form:"password"`
	BusinessUserFirstName string `json:"business_user_first_name" form:"business_user_first_name"`
	BusinessUserLastName  string `json:"business_user_last_name" form:"business_user_last_name"`
	BusinessUserFullName  string `json:"business_user_full_name" form:"business_user_full_name"`
	Language              string `json:"language" form:"language"`
	Qos                   string `json:"qos" form:"qos"`
	IsEncrypt             bool   `json:"is_encrypt" form:"is_encrypt"`
}

func RegisterUser(c echo.Context) error {
	param := &UserParam{}
	param.IsEncrypt = true // default value
	err := c.Bind(param)
	if err != nil {
		return c.JSON(response.BadRequestRes.Code, response.Format{
			Code:    response.BadRequestRes.Code,
			Message: response.BadRequestRes.Message,
		})
	}

	// validate input fields.
	unverifiedUser := &models.User{
		BusinessPartner:       param.BusinessPartner,
		BusinessPartnerName:   param.BusinessPartnerName,
		EmailAddress:          param.EmailAddress,
		Password:              param.Password,
		BusinessUserFirstName: param.BusinessUserFirstName,
		BusinessUserLastName:  param.BusinessUserLastName,
		BusinessUserFullName:  param.BusinessUserFullName,
		Language:              param.Language,
		Qos:                   models.ToQos(param.Qos),
	}
	if unverifiedUser.NeedsValidation() {
		if err := unverifiedUser.Validate(); err != nil {
			c.Logger().Printf("Failed to validate input parameter: %v", err)
			return c.JSON(response.BadRequestRes.Code, response.Format{
				Code:    response.BadRequestRes.Code,
				Message: response.BadRequestRes.Message,
			})
		}
	}

	// check registration status of login id.
	user := models.NewUser()
	result, err := user.GetByEmailAddress(param.EmailAddress)
	if result != nil && err == nil {
		c.Logger().Printf("Login id is already used")
		return c.JSON(response.Conflict.Code, response.Format{
			Code:    response.Conflict.Code,
			Message: response.Conflict.Message,
		})
	}

	userImp := &models.User{
		BusinessPartner:       param.BusinessPartner,
		BusinessPartnerName:   param.BusinessPartnerName,
		EmailAddress:          param.EmailAddress,
		Password:              param.Password,
		BusinessUserFirstName: param.BusinessUserFirstName,
		BusinessUserLastName:  param.BusinessUserLastName,
		BusinessUserFullName:  param.BusinessUserFullName,
		Language:              param.Language,
		Qos:                   unverifiedUser.Qos,
		IsEncrypt:             &param.IsEncrypt,
		LastLoginAt:           nil,
	}

	// encrypt password
	if param.IsEncrypt {
		encryptedPassword, err := crypto.Encrypt(param.Password)
		if err != nil {
			c.Logger().Printf("Failed to encrypt password: %v", err)
			return c.JSON(response.InternalErrRes.Code, response.Format{
				Code:    response.InternalErrRes.Code,
				Message: response.InternalErrRes.Message,
			})
		}
		userImp.Password = encryptedPassword
	}
	user.SetUser(userImp)

	err = user.Register()
	if err != nil {
		return err
	}
	return nil
}

type UpdateUserParam struct {
	UserParam
	OldPassword string `json:"old_password" form:"old_password"`
}

func (p *UpdateUserParam) PasswordExists() bool {
	return p.Password != ""
}

func (p *UpdateUserParam) QosExists() bool {
	return p.Qos != ""
}

func UpdateUser(c echo.Context) error {
	param := &UpdateUserParam{}
	param.IsEncrypt = true // default value
	err := c.Bind(param)
	if err != nil {
		return c.JSON(response.BadRequestRes.Code, response.Format{
			Code:    response.BadRequestRes.Code,
			Message: response.BadRequestRes.Message,
		})
	}

	// check existence of user.
	result, err := models.NewUser().GetByEmailAddress(c.Param("email_address"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Logger().Printf("Email Address is not found: %s", param.EmailAddress)
			return c.JSON(response.NotFoundErrRes.Code, response.Format{
				Code:    response.NotFoundErrRes.Code,
				Message: response.NotFoundErrRes.Message,
			})
		} else {
			c.Logger().Printf("Failed to db access: %v", err)
			return c.JSON(response.InternalErrRes.Code, response.Format{
				Code:    response.InternalErrRes.Code,
				Message: response.InternalErrRes.Message,
			})
		}
	}
	if result != nil && result.IsDeleted() {
		c.Logger().Print("User is already deleted.")
		return c.JSON(response.Conflict.Code, response.Format{
			Code:    response.Conflict.Code,
			Message: response.Conflict.Message,
		})
	}

	// authenticate old password.
	if !*result.IsEncrypt {
		if result.Password != param.OldPassword {
			c.Logger().Print("Failed to login due to incorrect password")
			return c.JSON(response.UnauthorizedRes.Code, response.Format{
				Code:    response.UnauthorizedRes.Code,
				Message: response.UnauthorizedRes.Message,
			})
		}
	} else {
		if err := crypto.CompareHashAndPassword(result.Password, param.OldPassword); err != nil {
			c.Logger().Printf("Failed to login: %v", err)
			return c.JSON(response.UnauthorizedRes.Code, response.Format{
				Code:    response.UnauthorizedRes.Code,
				Message: response.UnauthorizedRes.Message,
			})
		}
	}

	userImp := &models.User{
		BusinessPartner: param.BusinessPartner,
		EmailAddress:    param.EmailAddress,
		Password:        param.Password,
		Qos:             models.ToQos(param.Qos),
		IsEncrypt:       &param.IsEncrypt,
	}
	if !param.QosExists() {
		userImp.Qos = result.Qos
	}

	// validate input params.
	if userImp.NeedsValidation() {
		if err := userImp.Validate(); err != nil {
			c.Logger().Printf("Failed to validate input parameter: %v", err)
			return c.JSON(response.BadRequestRes.Code, response.Format{
				Code:    response.BadRequestRes.Code,
				Message: response.BadRequestRes.Message,
			})
		}
	}

	// encrypt password.
	if param.PasswordExists() && param.IsEncrypt {
		encryptedPassword, err := crypto.Encrypt(param.Password)
		if err != nil {
			c.Logger().Printf("Failed to encrypt password: %v", err)
			return c.JSON(response.InternalErrRes.Code, response.Format{
				Code:    response.InternalErrRes.Code,
				Message: response.InternalErrRes.Message,
			})
		}
		userImp.Password = encryptedPassword
	}

	result.SetUser(userImp)
	if err = result.Update(); err != nil {
		return err
	}
	return nil
}

func GetUser(c echo.Context) error {
	EmailAddress := c.Param("email_address")
	result, err := models.NewUser().GetByEmailAddress(EmailAddress)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Logger().Printf("Email Address is not found: %s", EmailAddress)
			return c.JSON(response.NotFoundErrRes.Code, response.Format{
				Code:    response.NotFoundErrRes.Code,
				Message: response.NotFoundErrRes.Message,
			})
		} else {
			c.Logger().Printf("Failed to db access: %v", err)
			return c.JSON(response.InternalErrRes.Code, response.Format{
				Code:    response.InternalErrRes.Code,
				Message: response.InternalErrRes.Message,
			})
		}
	}
	if result != nil && result.IsDeleted() {
		c.Logger().Print("User is already deleted.")
		return c.JSON(response.Conflict.Code, response.Format{
			Code:    response.Conflict.Code,
			Message: response.Conflict.Message,
		})
	}

	return c.JSON(http.StatusOK, response.UserDetailResponseFormat{
		EmailAddress:          result.EmailAddress,
		BusinessPartner:       result.BusinessPartner,
		BusinessPartnerName:   result.BusinessPartnerName,
		BusinessUserFirstName: result.BusinessUserFirstName,
		BusinessUserLastName:  result.BusinessUserLastName,
		BusinessUserFullName:  result.BusinessUserFullName,
		Language:              result.Language,
	})
}

type DeleteUserParam struct {
	Password string `json:"password" form:"password"`
}

func DeleteUser(c echo.Context) error {
	// ユーザの理論削除を行います。
	// 具体的な処理としては、ユーザモデルの deleted_at フラグに現在時刻を登録します。

	param := &DeleteUserParam{}
	err := c.Bind(param)
	if err != nil {
		return c.JSON(response.BadRequestRes.Code, response.Format{
			Code:    response.BadRequestRes.Code,
			Message: response.BadRequestRes.Message,
		})
	}
	EmailAddress := c.Param("email_address")

	result, err := models.NewUser().GetByEmailAddress(EmailAddress)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Logger().Printf("Email Address is not found: %s", EmailAddress)
			return c.JSON(response.NotFoundErrRes.Code, response.Format{
				Code:    response.NotFoundErrRes.Code,
				Message: response.NotFoundErrRes.Message,
			})
		} else {
			c.Logger().Printf("Failed to db access: %v", err)
			return c.JSON(response.InternalErrRes.Code, response.Format{
				Code:    response.InternalErrRes.Code,
				Message: response.InternalErrRes.Message,
			})
		}
	}
	if result != nil && result.IsDeleted() {
		c.Logger().Print("User is already deleted.")
		return c.JSON(response.Conflict.Code, response.Format{
			Code:    response.Conflict.Code,
			Message: response.Conflict.Message,
		})
	}

	// authenticate password.
	if !*result.IsEncrypt {
		if result.Password != param.Password {
			c.Logger().Print("Failed to login due to incorrect password")
			return c.JSON(response.UnauthorizedRes.Code, response.Format{
				Code:    response.UnauthorizedRes.Code,
				Message: response.UnauthorizedRes.Message,
			})
		}
	} else {
		if err := crypto.CompareHashAndPassword(result.Password, param.Password); err != nil {
			c.Logger().Printf("Failed to login: %v", err)
			return c.JSON(response.UnauthorizedRes.Code, response.Format{
				Code:    response.UnauthorizedRes.Code,
				Message: response.UnauthorizedRes.Message,
			})
		}
	}
	now := time.Now()
	userImp := &models.User{EmailAddress: result.EmailAddress, DeletedAt: &now}
	if err := userImp.Update(); err != nil {
		return err
	}
	return nil
}
