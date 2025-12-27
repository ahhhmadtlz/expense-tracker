package service

import (
	"context"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
)

func (s Service) DeleteTransaction(ctx context.Context,transactionID,userID uint)(param.DeleteTransactionResponse,error){
	const op=richerror.Op("transactionservice.DeleteTransaction")

	logger.Info("Deleting transaction","transaction_id",transactionID,"user_id",userID)

	existingTx,err:=s.repo.GetByID(ctx,transactionID)

	if err!=nil{
		logger.Error("Failed to get transaction","transaction_id",transactionID,"error",err.Error())
		return param.DeleteTransactionResponse{},richerror.New(op).WithMessage("transaction not found").WithKind(richerror.KindNotFound).WithErr(err)
	}

	if existingTx.UserID !=userID{
		logger.Warn("Unauthorized transaction deletion attempt","transaction_id",transactionID,"transaction_owner",existingTx.UserID,"requesting_user",userID)
		return param.DeleteTransactionResponse{},richerror.New(op).WithMessage("unauthorized").WithKind(richerror.KindForbidden)
	}

	if err :=s.repo.Delete(ctx,transactionID);err !=nil{
		logger.Error("Failed to delete transaction","transaction_id",transactionID,"error",err.Error())

		return  param.DeleteTransactionResponse{},richerror.New(op).WithMessage("failed to delete transaction").WithKind(richerror.KindUnexpected).WithErr(err)
	}

	logger.Info("Transaction deleted successfully","transaction_id",transactionID,"user_id",userID)



	return param.DeleteTransactionResponse{
		Message: "transaction deleted successfully",
	},nil
}