package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"data-platform-authenticator/configs"
	"data-platform-authenticator/internal/crypto"
	"data-platform-authenticator/internal/models"
	response "data-platform-authenticator/pkg/response"

	"github.com/form3tech-oss/jwt-go"
	"github.com/labstack/echo/v4"
)

type UserLoginParam struct {
	EmailAddress string `json:"email_address" form:"email_address"`
	Password     string `json:"password" form:"password"`
}

var jwtExp int64
var privateKeyPem string
var publicKeyPem string

type TokenParam struct {
	JwtToken string `json:"jwt_token" form:"jwt_token"`
}

func init() {
	cfgs, err := configs.New()
	if err != nil {
		panic(err)
	}
	jwtExp = cfgs.Jwt.Exp
	privateKeyPem = cfgs.PrivateKey
	publicKeyPem = cfgs.PublicKey
}

func EnsureUser(c echo.Context) error {
	param := &UserLoginParam{}
	err := c.Bind(param)
	if err != nil {
		return c.JSON(response.BadRequestRes.Code, response.Format{
			Code:    response.BadRequestRes.Code,
			Message: response.BadRequestRes.Message,
		})
	}
	user := models.NewUser()
	result, err := user.GetByEmailAddress(param.EmailAddress)

	if err != nil {
		return c.JSON(response.NotFoundErrRes.Code, response.Format{
			Code:    response.NotFoundErrRes.Code,
			Message: response.NotFoundErrRes.Message,
		})
	}
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

	// generate JWT
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email_address"] = user.User().EmailAddress
	claims["exp"] = strconv.FormatInt(time.Now().Add(time.Hour*time.Duration(jwtExp)).Unix(), 10)

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPem))
	if err != nil {
		c.Logger().Printf("Failed to parse private key: %v", err)
		return c.JSON(response.InternalErrRes.Code, response.Format{
			Code:    response.InternalErrRes.Code,
			Message: response.InternalErrRes.Message,
		})
	}
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		c.Logger().Printf("Failed to generate JWT: %v", err)
		return c.JSON(response.InternalErrRes.Code, response.Format{
			Code:    response.InternalErrRes.Code,
			Message: response.InternalErrRes.Message,
		})
	}

	if err := result.Login(); err != nil {
		c.Logger().Printf("Failed to record last_login_at: %v", err)
		return c.JSON(response.InternalErrRes.Code, response.Format{
			Code:    response.InternalErrRes.Code,
			Message: response.InternalErrRes.Message,
		})
	}
	return c.JSON(http.StatusOK, response.JWTResponseFormat{Jwt: signedToken})
}

func VerifyJWTToken(c echo.Context) error {
	param := &TokenParam{}
	err := c.Bind(param)
	if err != nil {
		return c.JSON(response.BadRequestRes.Code, response.Format{
			Code:    response.BadRequestRes.Code,
			Message: response.BadRequestRes.Message,
		})
	}

	var jwtStr string

	for key, values := range c.Request().Header {
		if key == "Authorization" {
			jwtStr = strings.Replace(values[0], "Bearer ", "", 1)
		}
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyPem))
	if err != nil {
		c.Logger().Printf("Failed to parse public key: %v", err)
		return c.JSON(response.InternalErrRes.Code, response.InternalErrRes)
	}

	parsedToken, err := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodRS256 {
			return nil, errors.New("invalid signing method")
		}
		return publicKey, nil
	})

	claims := parsedToken.Claims.(jwt.MapClaims)

	intUnix, err := strconv.Atoi(claims["exp"].(string))
	if err != nil {
		c.Logger().Printf("Failed to parse claims exp: %v", err)
		return c.JSON(response.InternalErrRes.Code, response.InternalErrRes)
	}

	if time.Now().Unix() > int64(intUnix) {
		return c.JSON(response.UnauthorizedRes.Code, response.UnauthorizedRes)
	}

	return c.JSON(http.StatusOK, response.UserVerifyResponseFormat{
		EmailAddress: fmt.Sprintf("%v", claims["email_address"]),
	})
}
