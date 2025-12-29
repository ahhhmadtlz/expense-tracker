package validator

import (
	"context"

	categoryentity "github.com/ahhhmadtlz/expense-tracker/internal/domain/category/entity"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/entity"
)

type Repository interface {
	GetByID(ctx context.Context, transactionID uint) (entity.Transaction,error)
}


type CategoryRepository interface {
	GetByID(ctx context.Context,categoryID uint)(categoryentity.Category,error)
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