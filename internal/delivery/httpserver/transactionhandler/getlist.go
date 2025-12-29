package transactionhandler

import (
	"net/http"
	"strconv"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) listTransactions(c echo.Context) error {
	logger.Info("List transactions request received")

	userID := c.Get("user_id").(uint)

	// Build request with query parameters
	req :=param.ListTransactionsRequest{
		Type:      c.QueryParam("type"),       // Filter by income/expense
		StartDate: c.QueryParam("start_date"), // Format: YYYY-MM-DD
		EndDate:   c.QueryParam("end_date"),   // Format: YYYY-MM-DD
	}

	// Parse category_id if provided
	if categoryIDStr := c.QueryParam("category_id"); categoryIDStr != "" {
		categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err != nil {
			logger.Warn("Invalid category_id", "category_id", categoryIDStr)
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "invalid category_id",
			})
		}
		catID := uint(categoryID)
		req.CategoryID = &catID
	}

	resp, err := h.transactionSvc.ListTransactions(c.Request().Context(), userID, req)

	if err != nil {
		logger.Error("Failed to list transactions", "user_id", userID, "error", err.Error())
		msg, code := httpmsgerrorhandler.Error(err)
		return c.JSON(code, echo.Map{
			"message": msg,
		})
	}

	logger.Info("Transactions listed successfully",
		"user_id", userID,
		"count", len(resp.Transactions),
		"total_income", resp.TotalIncome,
		"total_expense", resp.TotalExpense,
		"balance", resp.Balance,
	)

	return c.JSON(http.StatusOK, echo.Map{
		"message": "transactions retrieved successfully",
		"data":    resp,
	})
}