package httpserver

import (
	"github.com/ahhhmadtlz/expense-tracker/internal/config"
	"github.com/ahhhmadtlz/expense-tracker/internal/delivery/httpserver/userhandler"
	"github.com/labstack/echo/v4"
)

type Server struct {
	config          config.Config
	userHandler     userhandler.Handler

	Router          *echo.Echo
}
