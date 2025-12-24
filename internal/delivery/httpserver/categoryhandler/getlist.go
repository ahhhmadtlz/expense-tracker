package categoryhandler

import (
	"net/http"

	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) listCategories(c echo.Context) error {
	logger.Info("List categories request received")

	userID:=c.Get("user_id").(uint)

	catType:=c.QueryParam("type")

	resp,err:=h.categorySvc.ListCategories(c.Request().Context(),userID,catType)

	if err!=nil{
		logger.Error("Failed to list categories","user_id",userID,"error",err.Error())
		msg,code:=httpmsgerrorhandler.Error(err)
		return c.JSON(code,echo.Map{
			"message":msg,
		})
	}
	logger.Info("Categories listed successfully","user_id",userID,"count",len(resp.Categories),
	)

	return c.JSON(http.StatusOK,echo.Map{
		"message":"categories retrieved successfully",
		"data":resp,
	})

}