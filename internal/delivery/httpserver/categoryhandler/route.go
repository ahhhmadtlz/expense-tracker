package categoryhandler

import (
	"github.com/ahhhmadtlz/expense-tracker/internal/delivery/httpserver/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	categoryGroup := e.Group("/categories")

	categoryGroup.Use(middleware.Auth(h.authSvc,h.authConfig),middleware.UserContext())

	categoryGroup.POST("",h.createCategory)
	categoryGroup.GET("",h.listCategories)
	categoryGroup.GET("/:id",h.getCategory)
	categoryGroup.PUT("/:id",h.updateCategory)
	categoryGroup.DELETE("/:id",h.deleteCategory)
}