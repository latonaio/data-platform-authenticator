package route

import (
	"data-platform-authenticator/internal/handler"
	"github.com/labstack/echo/v4"
)

func RegisterRouting(e *echo.Echo) {
	// ユーザー作成
	e.POST("/users", handler.RegisterUser)
	// ユーザー認証
	e.POST("/login", handler.EnsureUser)
	// ユーザー更新
	e.PUT("/users/login_id/:email_address", handler.UpdateUser)
	// ユーザー取得
	e.GET("/users/login_id/:email_address", handler.GetUser)
	// ユーザー削除
	e.POST("/users/login_id/:email_address", handler.DeleteUser)
	// トークン確認
	e.POST("/token/verify", handler.VerifyJWTToken)
}
