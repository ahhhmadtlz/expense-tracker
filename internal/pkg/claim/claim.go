package claim

import (
	"github.com/ahhhmadtlz/expense-tracker/internal/config"
	"github.com/ahhhmadtlz/expense-tracker/internal/domain/auth"
	"github.com/labstack/echo/v4"
)

func GetClaimsFromEchoContext(c echo.Context) *auth.Claims {
	return c.Get(config.AuthMiddlewareContextKey).(*auth.Claims)
}