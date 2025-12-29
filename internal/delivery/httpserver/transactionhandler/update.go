package transactionhandler

import (
	"net/http"
	"strconv"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) updateTransaction(c echo.Context) error {
	var req param.UpdateTransactionRequest

	logger.Info("Update transaction request received")

	if err := c.Bind(&req); err != nil {
		logger.Warn("Failed to bind request", "error", err.Error())
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}

	idParam := c.Param("id")
	transactionID, err := strconv.ParseUint(idParam, 10, 32)

	if err != nil {
		logger.Warn("Invalid transaction ID", "id", idParam)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid transaction id",
		})
	}

	userID := c.Get("user_id").(uint)

	if fieldErrors, err := h.transactionValidator.ValidateUpdateTransaction(c.Request().Context(), req); err != nil {
		logger.Warn("Validation failed", "transaction_id", transactionID, "user_id", userID, "field_errors", fieldErrors)
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fieldErrors,
		})
	}

	resp, err := h.transactionSvc.UpdateTransaction(c.Request().Context(), req, uint(transactionID), userID)

	if err != nil {
		logger.Error("Failed to update transaction", "transaction_id", transactionID, "user_id", userID, "error", err.Error())
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
		})
	}

	logger.Info("Transaction updated successfully", "transaction_id", transactionID, "user_id", userID)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "transaction updated successfully",
		"data":    resp,
	})
}