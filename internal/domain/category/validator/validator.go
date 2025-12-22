package validator

import (
	"context"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/category/entity"
)

type Repository interface {
	GetByID(ctx context.Context, categoryID uint) (entity.Category,error)
}


type Validator struct {
	repo Repository
}

func New(repo Repository)Validator {
	return Validator{
		repo: repo,
	}
}