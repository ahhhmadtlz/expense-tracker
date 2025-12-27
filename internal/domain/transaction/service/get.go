package service

import (
	"context"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
)

func (s Service) GetTransaction(ctx context.Context, transactionID,userID uint)(param.GetTransactionResponse,error){
	const op=richerror.Op("transactionservice.GetTransaction")

	logger.Debug("Getting transaction","transaction_id",transactionID,"user_id",userID)

	tx,err:=s.repo.GetByID(ctx,transactionID)

	if err!=nil{
		logger.Error("Failed to get transaction","transaction_id",transactionID,"error",err.Error())
		return param.GetTransactionResponse{},richerror.New(op).WithMessage("transaction not found").WithKind(richerror.KindNotFound).WithErr(err)
	}

	if tx.UserID !=userID {
		logger.Warn("Unauthorized transaction access attempt","transaction_id",transactionID,"transaction_owner",tx.UserID,"requesting_user",userID)
	
		return param.GetTransactionResponse{},richerror.New(op).WithMessage("unauthorized").WithKind(richerror.KindForbidden)
	}

	return param.GetTransactionResponse{
		Transaction: param.ToTransactionInfo(tx),
	},nil
}