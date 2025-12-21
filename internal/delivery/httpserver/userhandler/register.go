package userhandler

import (
	"net/http"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/user/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) userRegister(c echo.Context) error {
	var req param.RegisterRequest
	
	logger.Info("Register request received")
	
	if err := c.Bind(&req); err != nil {
		logger.Warn("Failed to bind request", "error", err.Error())
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}
	
	logger.Debug("Request bound successfully", "phone_number", req.PhoneNumber)
	
	if fieldErrors, err := h.userValidator.ValidateRegisterRequest(c.Request().Context(), req); err != nil {
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
	
	resp, err := h.userSvc.Register(c.Request().Context(), req)
	if err != nil {
		logger.Error("Registration failed",
			"phone_number", req.PhoneNumber,
			"error", err.Error(),
		)
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
		})
	}
	
	logger.Info("User registered successfully", "phone_number", req.PhoneNumber)
	
	return c.JSON(http.StatusCreated, echo.Map{
		"message": "user registered successfully",
		"data":    resp,
	})
}