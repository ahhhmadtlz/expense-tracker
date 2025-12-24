package categoryhandler

import (
	"net/http"
	"strconv"

	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) deleteCategory(c echo.Context)error {
	logger.Info("Delete category request received")

	paramID:=c.Param("id")
	categoryID,err:=strconv.ParseUint(paramID,10,32)

	if err!=nil{
		logger.Warn("Invalid category ID","id",paramID)

		return c.JSON(http.StatusBadRequest,echo.Map{
			"message":"invalid category id",
		})
	}

	userID:=c.Get("user_id").(uint)

	resp,err:=h.categorySvc.DeleteCategory(c.Request().Context(),uint(categoryID),userID)
	if err!=nil{
		logger.Error("Failed to delete category","category_id",categoryID,"user_id",userID,"error",err.Error())
		msg,code:=httpmsgerrorhandler.Error(err)

		return c.JSON(code,echo.Map{
			"message":msg,
		})
	}
	logger.Info("category deleted successfully","category_id",categoryID,"user_id",userID)

	return c.JSON(http.StatusOK,echo.Map{
		"message":resp.Message,
	})
}