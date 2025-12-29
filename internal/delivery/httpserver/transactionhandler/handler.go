package transactionhandler

import (
	"github.com/ahhhmadtlz/expense-tracker/internal/domain/auth"
	transactionservice "github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/service"
	transactionvalidator "github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/validator"
)

type Handler struct {
	authConfig auth.Config
	authSvc auth.Service
	transactionSvc transactionservice.Service
	transactionValidator transactionvalidator.Validator
}

func New(
	authConfig auth.Config,
	authSvc auth.Service,
	transactionSvc transactionservice.Service,
	transactionValidator transactionvalidator.Validator,
)Handler {
	return Handler{
		authConfig: authConfig,
		authSvc: authSvc,
		transactionSvc: transactionSvc,
		transactionValidator:transactionValidator ,
	}
}