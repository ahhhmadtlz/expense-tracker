package transactionhandler

import (
	"github.com/ahhhmadtlz/expense-tracker/internal/delivery/httpserver/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	transactionGroup := e.Group("/transactions")

	transactionGroup.Use(middleware.Auth(h.authSvc, h.authConfig), middleware.UserContext())

	transactionGroup.POST("", h.createTransaction)
	transactionGroup.GET("", h.listTransactions)
	transactionGroup.GET("/:id", h.getTransaction)
	transactionGroup.PUT("/:id", h.updateTransaction)
	transactionGroup.DELETE("/:id", h.deleteTransaction)
}