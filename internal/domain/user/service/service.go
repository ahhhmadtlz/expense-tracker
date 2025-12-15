package service

import (
	"context"
	"log/slog"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/auth"
	"github.com/ahhhmadtlz/expense-tracker/internal/domain/user/entity"
)


type AuthService interface {
	CreateAccessToken(user entity.User)(string,error)
	CreateRefreshToken(user entity.User)(string,error)
	ParseRefreshToken(refreshToken string)(*auth.Claims,error)
}


type Repositroy interface {
	RegisterUser(ctx context.Context, user entity.User)(entity.User,error)
	GetUserByID(ctx context.Context,userID uint)(entity.User,error)
	GetUserByPhoneNumber(ctx context.Context,phone string)(entity.User, error)
}


type Service struct {
	repo Repositroy
	auth AuthService
	logger *slog.Logger
}

func New(authService AuthService,repo Repositroy,logger *slog.Logger)Service{
		return Service{
			auth:authService,
			repo:repo,
			logger:logger,
		}
}
