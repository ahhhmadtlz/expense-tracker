package service

import (
	"context"
	"time"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/entity"
	"github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
)

func (s Service) UpdateTransaction(ctx context.Context, req param.UpdateTransactionRequest,transactionID,userID uint)(param.UpdateTransactionResponse,error){
	const op =richerror.Op("transactionservice.UpdateTransaction")
	
	logger.Info("Updating transaction","transaction_id",transactionID,"user_id",userID)

	existingTx,err:=s.repo.GetByID(ctx,transactionID)
	if err!=nil{
		logger.Error("Failed to get transaction","transaction_id",transactionID,"error",err.Error())
		return  param.UpdateTransactionResponse{},richerror.New(op).WithMessage("transaction not found").WithKind(richerror.KindNotFound).WithErr(err)
	}

	if existingTx.UserID !=userID {
		logger.Warn("Unauthorized transaction update attempt","transaction_id",transactionID,"transaction_owneer",existingTx.UserID,"requesting_user",userID)

		return  param.UpdateTransactionResponse{},richerror.New(op).WithMessage("unauthorized").WithKind(richerror.KindForbidden)
	}

	if req.CategoryID !=nil{
		existingTx.CategoryID =*req.CategoryID
	}

	if req.Type !=nil{
		existingTx.Type=entity.MapToTransactionType(*req.Type)
	}

	if req.Amount !=nil{
		existingTx.Amount =*req.Amount
	} 

	if req.Description !=nil{
		existingTx.Description=*req.Description
	}

	if req.Date !=nil && * req.Date !=""{
		date ,err:=time.Parse("2006-01-02",*req.Date)
		if err!=nil{
			logger.Error("Failed to parse date","date",*req.Date,"error",err.Error())
			return param.UpdateTransactionResponse{},richerror.New(op).WithMessage("invalid date format").WithKind(richerror.KindInvalid).WithErr(err)
		}
		existingTx.Date=date
	}

	updatedTx,err:=s.repo.Update(ctx,existingTx)

	if err!=nil{
		logger.Error("Failed to update transaction","transaction_id",transactionID,"error",err.Error())
		return param.UpdateTransactionResponse{},richerror.New(op).WithMessage("failed to udpate transaction").WithKind(richerror.KindUnexpected).WithErr(err)
	}
	logger.Info("Transaction updated successfully","transaction_id",transactionID,"user_id",userID)

	return param.UpdateTransactionResponse{
		Transaction: param.ToTransactionInfo(updatedTx),
	},nil
}