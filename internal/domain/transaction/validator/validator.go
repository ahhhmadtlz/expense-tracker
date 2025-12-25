package validator

import (
	"context"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/entity"
)

type Repository interface {
	GetByID(ctx context.Context, transactionID uint) (entity.Transaction,error)
}


type CategoryRepository interface {
	GetByID(ctx context.Context,categoryID uint)(any,error)
}

type Validator struct {
	repo Repository
	catRepo CategoryRepository
}


func New(repo Repository,catRepo CategoryRepository)Validator{
	return Validator{
		repo:repo,
		catRepo: catRepo,
	}
}