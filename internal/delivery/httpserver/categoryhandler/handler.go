package categoryhandler

import (
	"github.com/ahhhmadtlz/expense-tracker/internal/domain/auth"
	categoryservice "github.com/ahhhmadtlz/expense-tracker/internal/domain/category/service"
	categoryvalidator "github.com/ahhhmadtlz/expense-tracker/internal/domain/category/validator"
)
type Handler struct {
	authConfig auth.Config
	authSvc auth.Service
	categorySvc categoryservice.Service
	categoryValidator categoryvalidator.Validator
}


func New(
	authConfig auth.Config,
	authSvc auth.Service,
	categorySvc categoryservice.Service,
	categoryValidator categoryvalidator.Validator,
)Handler {
	return Handler{
		authConfig: authConfig,
		authSvc: authSvc,
		categorySvc: categorySvc,
		categoryValidator: categoryValidator,
	}
}