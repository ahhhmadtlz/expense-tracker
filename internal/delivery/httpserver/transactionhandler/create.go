package transactionhandler

import (
	"net/http"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/httpmsgerrorhandler"
	"github.com/labstack/echo/v4"
)

func (h Handler) createTransaction(c echo.Context) error {
	var req param.CreateTransactionRequest

	logger.Info("Create transaction requestt received")

	if err:=c.Bind(&req);err!=nil{
		logger.Warn("Failed to bind request","error",err.Error())

		return c.JSON(http.StatusBadRequest,echo.Map{
			"message":"invalid request body",
		})
	}

	userID:=c.Get("user_id").(uint)

	if fieldErrors,err:=h.transactionValidator.ValidateCreateTransaction(c.Request().Context(),req,userID);err!=nil{
		logger.Warn("Validation failed","user_id",userID,"field_erros",fieldErrors)
		msg,code:=httpmsgerrorhandler.Error(err)

		return c.JSON(code,echo.Map{
			"message":msg,
			"errors":fieldErrors,
		})
	}

	resp,err:=h.transactionSvc.CreateTransaction(c.Request().Context(),req,userID)

	if err!=nil{
		logger.Error("Failed to create transaction","user_id",userID,"error",err.Error())
		msg,code:=httpmsgerrorhandler.Error(err)

		return c.JSON(code,echo.Map{
			"message":msg,
		})
	}

	logger.Info("Transaction created successfully","user_id",userID,"transaction_id",resp.Transaction.ID)

	return c.JSON(http.StatusCreated,echo.Map{
		"message":"transaction created successfuly",
		"data":resp,
	})

}