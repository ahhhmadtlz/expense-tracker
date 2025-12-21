package service

import (
	"context"
	"github.com/ahhhmadtlz/expense-tracker/internal/domain/category/entity"
)

type Repository interface {
	Create(ctx context.Context, category entity.Category)(entity.Category,error)
	GetByID(ctx context.Context,categoryID uint)(entity.Category,error)
	GetByUserID(ctx context.Context,userID uint)([]entity.Category,error)
	GetByUserIDAndType(ctx context.Context,userID uint ,catType entity.CategoryType)([]entity.Category,error)
	update(ctx context.Context,category entity.Category)(entity.Category,error)
	Delete(ctx context.Context,categoryID uint)error
	CategoryHasTransactions(ctx context.Context,categoryID uint)(bool,error)
}


type Service struct {
	repo Repository
	
}


func New(repo Repository)Service{
	return Service{
		repo:repo,
	}
}