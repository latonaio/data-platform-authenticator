package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4/middleware"
	configs "jwt-authentication-golang/configs"
	"jwt-authentication-golang/internal/route"
	customers "jwt-authentication-golang/pkg/response"

	"github.com/labstack/echo/v4"
)

type server struct {
	Server  *echo.Echo
	Context context.Context
	Port    string
}

type Server interface {
	Start(errC chan error)
	Shutdown(ctx context.Context) error
}

func New(ctx context.Context, cfgs *configs.Configs) Server {
	// Echo instance
	e := echo.New()
	// Routes
	route.RegisterRouting(e)
	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	// use echo default logger
	e.Use(middleware.Logger())
	// error handling
	e.HTTPErrorHandler = customers.CustomHTTPErrorHandler

	return &server{
		Server:  e,
		Context: ctx,
		Port:    fmt.Sprintf(":%v", cfgs.Server.Port),
	}
}

func (s *server) Start(errC chan error) {
	err := s.Server.Start(s.Port)
	errC <- err
}

func (s *server) Shutdown(ctx context.Context) error {
	if err := s.Server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}
