package userhandler

import (

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/auth"
	userservice "github.com/ahhhmadtlz/expense-tracker/internal/domain/user/service"
	uservalidator "github.com/ahhhmadtlz/expense-tracker/internal/domain/user/validator"
)

type Handler struct {
	authSvc       auth.Service
	userSvc       userservice.Service
	userValidator uservalidator.Validator
	authConfig    auth.Config
}

func New(
	authSvc auth.Service,
	userSvc userservice.Service,
	userValidator uservalidator.Validator,
	authConfig auth.Config,
 
) Handler {
	return Handler{
		authSvc:       authSvc,
		userSvc:       userSvc,
		userValidator: userValidator,
		authConfig:    authConfig,
		 
	}
}