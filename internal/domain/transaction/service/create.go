package service

import (
	"context"
	"time"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/entity"
	"github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
)



func (s Service) CreateTransaction(ctx context.Context , req param.CreateTransactionRequest,userID uint)(param.CreateTransactionResponse,error){
	const op=richerror.Op("tranactionService.CreateTransaction")
	logger.Info("Creating transaction","user_id",userID,"type",req.Type,"amount",req.Amount)

	date,err:=time.Parse("2006-01-02",req.Date)

	if err!=nil{
		logger.Error("Failed to parse date","date",req.Date,"error",err.Error())

		return param.CreateTransactionResponse{},richerror.New(op).WithMessage("invalid date format").WithKind(richerror.KindInvalid).WithErr(err)
	}

	txType :=entity.MapToTransactionType(req.Type)

	tx:=entity.Transaction{
		UserID: userID,
		CategoryID: req.CategoryID,
		Type: txType,
		Amount: req.Amount,
		Description: req.Description,
		Date: date,
	}

	createdTx,err:=s.repo.Create(ctx,tx)

	if err!=nil{
		logger.Error("Failed to create transaction","user_id",userID,"type",req.Type,"error",err.Error())
		return param.CreateTransactionResponse{},richerror.New(op).WithMessage("failed to create transaction").WithKind(richerror.KindUnexpected).WithErr(err)
	}

	logger.Info("Transaction created successfully","transaction_id",createdTx.ID,"user_id",userID)

	return param.CreateTransactionResponse{
		Transaction: param.ToTransactionInfo(createdTx),
	},nil

	
}