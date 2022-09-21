package custmres

import (
	"net/http"

	"github.com/labstack/echo/v4"
	customError "jwt-authentication-golang/pkg/error"
)

type ResponseFormat struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	res := generateErrorResponse(err)
	_ = c.JSON(res.Code, &res)
	c.Logger().Error(err)
}

func generateErrorResponse(err error) ResponseFormat {
	customerErr, ok := err.(customError.CustomErrMessage)
	if !ok {
		return InternalErrRes
	}
	switch customerErr {
	case customError.ErrBadRequest:
		return BadRequestRes
	case customError.ErrNotFound:
		return NotFoundErrRes
	}
	return InternalErrRes
}

var (
	BadRequestRes = ResponseFormat{
		Code:    http.StatusBadRequest,
		Message: customError.ErrBadRequest.Error(),
	}

	InternalErrRes = ResponseFormat{
		Code:    http.StatusInternalServerError,
		Message: customError.ErrInternal.Error(),
	}

	NotFoundErrRes = ResponseFormat{
		Code:    http.StatusNotFound,
		Message: customError.ErrNotFound.Error(),
	}

	UnauthorizedRes = ResponseFormat{
		Code:    http.StatusUnauthorized,
		Message: customError.ErrUnauthorized.Error(),
	}

	Conflict = ResponseFormat{
		Code:    http.StatusConflict,
		Message: customError.ErrConflict.Error(),
	}
)
