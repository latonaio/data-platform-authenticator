package response

import (
	"net/http"

	customError "data-platform-authenticator/pkg/error"
	"github.com/labstack/echo/v4"
)

type Format struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func CustomHTTPErrorHandler(err error, c echo.Context) {
	res := generateErrorResponse(err)
	_ = c.JSON(res.Code, &res)
	c.Logger().Error(err)
}

func generateErrorResponse(err error) Format {
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
	BadRequestRes = Format{
		Code:    http.StatusBadRequest,
		Message: customError.ErrBadRequest.Error(),
	}

	InternalErrRes = Format{
		Code:    http.StatusInternalServerError,
		Message: customError.ErrInternal.Error(),
	}

	NotFoundErrRes = Format{
		Code:    http.StatusNotFound,
		Message: customError.ErrNotFound.Error(),
	}

	UnauthorizedRes = Format{
		Code:    http.StatusUnauthorized,
		Message: customError.ErrUnauthorized.Error(),
	}

	Conflict = Format{
		Code:    http.StatusConflict,
		Message: customError.ErrConflict.Error(),
	}
)
