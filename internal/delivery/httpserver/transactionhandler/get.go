package transactionhandler

import (
	"net/http"
	"strconv"

	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) getTransaction(c echo.Context)error{
	logger.Info("Get transaction request received")

	paramID:=c.Param("id")
	transactionID,err:=strconv.ParseUint(paramID,10,32)

	if err!=nil{
		logger.Warn("invalid transaction id ","id",paramID)
		return c.JSON(http.StatusBadRequest,echo.Map{
			"message":"invalid transactin id",
		})
	}

	userID:=c.Get("user_id").(uint)

	resp,err:=h.transactionSvc.GetTransaction(c.Request().Context(),uint(transactionID),userID)

	if err!=nil{
		logger.Error("Failed to get transaction","transaction_id",transactionID,"user_id",userID,"error",err.Error())
		msg,code:=httpmsgerrorhandler.Error(err)
		return c.JSON(code,echo.Map{
			"message":msg,
		})
	}

	return c.JSON(http.StatusOK,echo.Map{
		"message":"transaction retrieved successfully",
		"data":resp,
	})
}