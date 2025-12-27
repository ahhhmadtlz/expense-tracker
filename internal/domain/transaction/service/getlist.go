package service

import (
	"context"
	"time"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/entity"
	"github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
)

func (s Service) ListTransactions(ctx context.Context , userID uint ,req param.ListTransactionsRequest)(param.ListTransactionsResponse,error){
	const op=richerror.Op("transactionservice.ListTransactions")

	logger.Debug("Listing transactions","user_id",userID,"type_filter",req.Type)

	filters:=make(map[string]any)

	if req.Type !=""{
		filters["type"]=req.Type
	}

	if req.CategoryID!=nil &&*req.CategoryID>0{
		filters["category_id"]=*req.CategoryID
	}

	if req.StartDate != ""{
		startDate,err:=time.Parse("2006-01-02",req.StartDate)

		if err!=nil{
			logger.Error("Invalid start date format","start_date",req.StartDate,"error",err.Error())
			return  param.ListTransactionsResponse{},richerror.New(op).WithMessage("invalid start date format").WithKind(richerror.KindInvalid).WithErr(err)
		}
		filters["start_date"]=startDate
	}

	if req.EndDate !=""{
		endDate,err:=time.Parse("2006-01-02",req.EndDate)
		if err!=nil{
			logger.Error("Invalid end date format","end_date",req.EndDate,"error",err.Error())
			return param.ListTransactionsResponse{},richerror.New(op).WithMessage("invalid end_date format").WithKind(richerror.KindInvalid).WithErr(err)
		}
		filters["end_date"]=endDate
	}

	transactions,err:=s.repo.GetByUserID(ctx,userID,filters)

	if err!=nil{
		logger.Error("Failed to list transactions ","user_id",userID,"error",err.Error())

		return  param.ListTransactionsResponse{},richerror.New(op).WithMessage("failed to list transaction").WithKind(richerror.KindUnexpected).WithErr(err) 
	}
	transactionInfos:=make([]param.TransactionInfo,0,len(transactions))

	var totalIncome,totalExpense float64

	for _,tx:=range transactions {
		transactionInfos=append(transactionInfos, param.ToTransactionInfo(tx))
		if tx.Type==entity.TypeIncome {
			totalIncome +=tx.Amount
		}else{
			totalExpense +=tx.Amount
		}
	}
	balance:=totalIncome - totalExpense

	logger.Debug("Transactions listed successfully","user_id",userID,"count",len(transactionInfos),"total_income",totalIncome,"total_expense",totalExpense,"balance",balance)

	return  param.ListTransactionsResponse{
		Transactions: transactionInfos,
		TotalIncome: totalIncome,
		TotalExpense: totalExpense,
		Balance: balance,
	},nil
}