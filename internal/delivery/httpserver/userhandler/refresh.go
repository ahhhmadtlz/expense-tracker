package userhandler

import (
	"net/http"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/user/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) refreshAccessToken(c echo.Context) error {
	var req param.RefreshAccessTokenRequest
	
	logger.Info("Refresh token request received")
	
	if err := c.Bind(&req); err != nil {
		logger.Warn("Failed to bind request", "error", err.Error())
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}
	
	resp, err := h.userSvc.RefreshAccessToken(c.Request().Context(), req)
	if err != nil {
		logger.Error("Token refresh failed", "error", err.Error())
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
		})
	}
	
	logger.Info("Token refreshed successfully")
	
	return c.JSON(http.StatusOK, echo.Map{
		"message": "token refreshed successfully",
		"data":    resp,
	})
}