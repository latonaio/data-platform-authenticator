package handler

import (
	"net/http"
	"time"

	"github.com/form3tech-oss/jwt-go"
	"github.com/labstack/echo/v4"
	"jwt-authentication-golang/configs"
	"jwt-authentication-golang/internal/crypto"
	"jwt-authentication-golang/internal/models"
	customers "jwt-authentication-golang/pkg/response"
)

type UserLoginParam struct {
	LoginID  string `json:"login_id" form:"login_id"`
	Password string `json:"password" form:"password"`
}

var jwtExp int64
var privateKeyPem string

func init() {
	cfgs, err := configs.New()
	if err != nil {
		panic(err)
	}
	jwtExp = cfgs.Jwt.Exp
	privateKeyPem = cfgs.PrivateKey
}

func EnsureUser(c echo.Context) error {
	param := &UserLoginParam{}
	err := c.Bind(param)
	if err != nil {
		return c.JSON(customers.BadRequestRes.Code, customers.BadRequestRes)
	}
	user := models.NewUser()
	result, err := user.GetByLoginID(param.LoginID)
	if err != nil {
		return c.JSON(customers.NotFoundErrRes.Code, customers.NotFoundErrRes)
	}
	if !*result.IsEncrypt {
		if result.Password != param.Password {
			c.Logger().Print("Failed to login due to incorrect password")
			return c.JSON(customers.UnauthorizedRes.Code, customers.UnauthorizedRes)
		}
	} else {
		if err := crypto.CompareHashAndPassword(result.Password, param.Password); err != nil {
			c.Logger().Printf("Failed to login: %v", err)
			return c.JSON(customers.UnauthorizedRes.Code, customers.UnauthorizedRes)
		}
	}

	// generate JWT
	token := jwt.New(jwt.SigningMethodRS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.User().ID
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(jwtExp)).Unix()
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPem))
	if err != nil {
		c.Logger().Printf("Failed to parse private key: %v", err)
		return c.JSON(customers.InternalErrRes.Code, customers.InternalErrRes)
	}
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		c.Logger().Printf("Failed to generate JWT: %v", err)
		return c.JSON(customers.InternalErrRes.Code, customers.InternalErrRes)
	}

	if err := result.Login(); err != nil {
		c.Logger().Printf("Failed to record last_login_at: %v", err)
		return c.JSON(customers.InternalErrRes.Code, customers.InternalErrRes)
	}
	return c.JSON(http.StatusOK, customers.JWTResponseFormat{Jwt: signedToken})
}
