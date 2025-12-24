package categoryhandler

import (
	"net/http"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/category/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) createCategory(c echo.Context) error {
	var req param.CreateCategoryRequest

	logger.Info("Create category request received")

	if err :=c.Bind(&req);err !=nil{
		logger.Warn("Failed to bind request","error",err.Error())

		return  c.JSON(http.StatusBadRequest,echo.Map{
			"message":"invalid request body",
		})
	}

	userID:=c.Get("user_id").(uint)

	if fieldErrors,err:=h.categoryValidator.ValidateCreateCategory(c.Request().Context(),req);err!=nil{
		logger.Warn("Validation failed",
	 "user_id",userID,
	 "field_erros",fieldErrors,
	)
	msg,code:=httpmsgerrorhandler.Error(err)
	return  c.JSON(code,echo.Map{
		"message":msg,
		"errors":fieldErrors,
	})
	}

	resp,err:=h.categorySvc.CreateCategory(c.Request().Context(),req,userID)
	
	if err!=nil{
		logger.Error("Failed to create category",
	"user_id",userID,"error",err.Error())
	msg,code:=httpmsgerrorhandler.Error(err)
	return c.JSON(code,echo.Map{
		"message":msg,
	})
	}

	logger.Info("Category created successfully",
		"user_id",userID,
		"category_id",resp.Category.ID,
	)

	return c.JSON(http.StatusCreated,echo.Map{
		"message":"category created successfully",
		"data":resp,
	})

}