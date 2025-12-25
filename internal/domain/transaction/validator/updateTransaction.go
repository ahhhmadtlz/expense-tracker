package validator

import (
	"context"
	"strings"
	"time"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation"
)

func (v Validator) ValidateUpdateTransaction(ctx context.Context, req param.UpdateTransactionRequest)(map[string]string,error){
const op = "transactionvalidator.ValidateUpdateTransaction"

	fieldErrors := make(map[string]string)


	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		req.Description = &trimmed
	}
	if req.Type != nil {
		trimmed := strings.TrimSpace(*req.Type)
		req.Type = &trimmed
	}
	if req.Date != nil {
		trimmed := strings.TrimSpace(*req.Date)
		req.Date = &trimmed
	}

	// Validate CategoryID if provided
	if req.CategoryID != nil {
		err := validation.Validate(req.CategoryID,
			validation.Min(uint(1)).Error("category_id must be greater than 0"),
		)
		if err != nil {
			fieldErrors["category_id"] = err.Error()
		} else {
			// Validate category exists
			_, err := v.catRepo.GetByID(ctx, *req.CategoryID)
			if err != nil {
				if rErr, ok := err.(*richerror.RichError); ok {
					if rErr.GetKind() == richerror.KindNotFound {
						fieldErrors["category_id"] = "category not found"
					} else {
						return nil, richerror.New(op).
							WithErr(err).
							WithMessage("failed to verify category").
							WithKind(richerror.KindUnexpected)
					}
				}
			}
		}
	}

	
	if req.Type != nil {
		err := validation.Validate(req.Type,
			validation.In("income", "expense").Error("type must be either 'income' or 'expense'"),
		)
		if err != nil {
			fieldErrors["type"] = err.Error()
		}
	}

	
	if req.Amount != nil {
		err := validation.Validate(req.Amount,
			validation.Min(0.01).Error("amount must be greater than 0"),
		)
		if err != nil {
			fieldErrors["amount"] = err.Error()
		}
	}

	
	if req.Description != nil {
		err := validation.Validate(req.Description,
			validation.Length(0, 500).Error("description must not exceed 500 characters"),
		)
		if err != nil {
			fieldErrors["description"] = err.Error()
		}
	}

	
	if req.Date != nil && *req.Date != "" {
		_, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			fieldErrors["date"] = "date must be in format YYYY-MM-DD"
		}
	}

	if len(fieldErrors) > 0 {
		return fieldErrors, richerror.New(op).WithMessage("invalid input").WithKind(richerror.KindInvalid).WithMeta("fields", fieldErrors)
	}

	return nil, nil
}
