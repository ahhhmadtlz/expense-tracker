package service

import (
	"context"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/entity"
)

type Repository interface {
	Create(ctx context.Context,tx entity.Transaction)(entity.Transaction,error)
	GetByID(ctx context.Context,transactionID uint)(entity.Transaction,error)
	GetByUserID(ctx context.Context,userID uint,filters map[string]any)([]entity.Transaction,error)
	Update(ctx context.Context,tx entity.Transaction)(entity.Transaction,error)
	Delete(ctx context.Context,transactionID uint)error
}

type Service struct {
	repo Repository
}


func New (repo Repository)Service{
	return Service{
		repo: repo,
	}
}