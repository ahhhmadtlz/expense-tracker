package userhandler

import (
	"net/http"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/user/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) userLogin(c echo.Context) error {
	var req param.LoginRequest

	logger.Info("Login request received")

	if err := c.Bind(&req); err != nil {
		logger.Warn("Failed to bind request", "error", err.Error())
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}

	logger.Debug("Request bound successfully", "phone_number", req.PhoneNumber)

	if fieldErrors, err := h.userValidator.ValidateLoginRequest(c.Request().Context(), req); err != nil {
		logger.Warn("Validation failed",
			"phone_number", req.PhoneNumber,
			"field_errors", fieldErrors,
		)
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}

	logger.Debug("Validation passed, calling service", "phone_number", req.PhoneNumber)

	resp, err := h.userSvc.Login(c.Request().Context(), req)
	if err != nil {
		logger.Error("Login failed",
			"phone_number", req.PhoneNumber,
			"error", err.Error(),
		)
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
		})
	}

	logger.Info("Login successful", "phone_number", req.PhoneNumber)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "user login successfully",
		"data":    resp,
	})
}