package categoryhandler

import (
	"net/http"
	"strconv"

	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) getCategory(c echo.Context) error {
	logger.Info("Get category request received")

	paramID:=c.Param("id")
	categoryID,err:=strconv.ParseUint(paramID,10,32)
	if err!=nil{
		logger.Warn("Invalid category ID","id",paramID)
		return  c.JSON(http.StatusBadRequest,echo.Map{
			"message":"invalid category id",
		})
	}

	userID:=c.Get("user_id").(uint)

	resp,err:=h.categorySvc.GetCategory(c.Request().Context(),uint(categoryID),userID)

	if err!=nil{
		logger.Error("Failed to get category","category_id",categoryID,"user_id",userID,"error",err.Error())
		msg,code:=httpmsgerrorhandler.Error(err)
		return  c.JSON(code,echo.Map{
			"message":msg,
		})
	}

	return c.JSON(http.StatusOK,echo.Map{
		"message":"category retrieved successfully",
		"data":resp,
	})
}