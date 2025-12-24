package categoryhandler

import (
	"net/http"
	"strconv"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/category/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) updateCategory(c echo.Context) error {
	var req param.UpdateCategoryRequest

	logger.Info("Update category request received")

	if err := c.Bind(&req); err != nil {
		logger.Warn("Failed to bind request", "error", err.Error())
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid request body",
		})
	}

	idParam:=c.Param("id")
	categoryID, err := strconv.ParseUint(idParam, 10, 32)

	if err !=nil{
		logger.Warn("Invalid category ID","id",idParam)
		return  c.JSON(http.StatusBadRequest,echo.Map{
			"message":"invalid category id",
		})
	}

	userID:=c.Get("user_id").(uint)

	if fieldErrors,err:=h.categoryValidator.ValidateUpdateCategory(c.Request().Context(),req);err!=nil{
		logger.Warn("validation failed","category_id",categoryID,"user_id",userID,"field_errors",fieldErrors)

		msg,code:=httpmsgerrorhandler.Error(err)

		return c.JSON(code,echo.Map{
			"message":msg,
			"errors":fieldErrors,
		})
	}

	resp,err:=h.categorySvc.UpdateCategory(c.Request().Context(),req,uint(categoryID),userID)
	
	if err!=nil{
		logger.Error("Failed to udpate category","category_id",categoryID,"user_id",userID,"error",err.Error())
		msg,code:=httpmsgerrorhandler.Error(err)
		return c.JSON(code,echo.Map{
			"message":msg,
		})
	}

	logger.Info("Category udpate successfully","category_id",categoryID,"user_id",userID)
	
	return c.JSON(http.StatusOK,echo.Map{
		"message":"category updated successfully",
		"data":resp,
	})

}