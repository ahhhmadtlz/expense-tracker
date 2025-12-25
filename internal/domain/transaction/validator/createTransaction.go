package validator

import (
	"context"
	"strings"
	"time"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/transaction/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation"
)


func (v Validator)ValidateCreateTransaction(ctx context.Context,req param.CreateTransactionRequest,userID uint)(map[string]string ,error){
	const op="transactionvalidator.ValidateCreateTransaction"

	req.Description=strings.TrimSpace(req.Description)
	req.Type=strings.TrimSpace(req.Type)
	req.Date=strings.TrimSpace(req.Date)

	fieldErrors:=make(map[string]string)

	err:=validation.ValidateStruct(&req,validation.Field(&req.CategoryID,
		validation.Required.Error("category_id is required"),
		validation.Min(uint(1)).Error("category_id must be greater than 0"),
		),
		validation.Field(&req.Type,
		validation.Required.Error("type is required"),
		validation.In("income","expense").Error("type must be either 'income' or 'expense' "),
		),
		validation.Field(&req.Amount,validation.Required.Error("amount is required"),
		validation.Min(0.01).Error("amount must be greater than 0"),
		),
		validation.Field(&req.Description,
				validation.Length(0, 500).Error("description must not exceed 500 characters"),
		),
		validation.Field(&req.Date,
			validation.Required.Error("date is required"),
		),
	)

	if err != nil {
		if errV, ok := err.(validation.Errors); ok {
			for key, value := range errV {
				if value != nil {
					fieldErrors[key] = value.Error()
				}
			}
		}
	}

	if req.Date != "" {
		_, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			fieldErrors["date"] = "date must be in format YYYY-MM-DD"
		}
	}
	
	if req.CategoryID > 0 {
			_, err := v.catRepo.GetByID(ctx, req.CategoryID)
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

	if len(fieldErrors) > 0 {
		return fieldErrors, richerror.New(op).WithMessage("invalid input").WithKind(richerror.KindInvalid).WithMeta("fields", fieldErrors)
	}

	return nil, nil
}