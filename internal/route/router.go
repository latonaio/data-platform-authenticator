package route

import (
	"github.com/labstack/echo/v4"
	"jwt-authentication-golang/internal/handler"
)

func RegisterRouting(e *echo.Echo) {
	// ユーザー作成
	e.POST("/users", handler.RegisterUser)
	// ユーザー認証
	e.POST("/login", handler.EnsureUser)
	// ユーザー更新
	e.PUT("/users/login_id/:login_id", handler.UpdateUser)
	// ユーザー取得
	e.GET("/users/login_id/:login_id", handler.GetUser)
	// ユーザー削除
	e.POST("/users/login_id/:login_id", handler.DeleteUser)
}
