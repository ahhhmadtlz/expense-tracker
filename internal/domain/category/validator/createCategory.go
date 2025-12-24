package validator

import (
	"context"
	"regexp"
	"strings"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/category/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation"
)

func (v Validator) ValidateCreateCategory(ctx context.Context, req param.CreateCategoryRequest,userID uint)(map[string]string,error) {
	const op="categoryvalidator.ValidateCreateCategory"

	req.Name = strings.TrimSpace(req.Name)
	req.Type = strings.TrimSpace(req.Type)
	req.Color = strings.TrimSpace(req.Color)

	err:=validation.ValidateStruct(&req,validation.Field(&req.Name,
			validation.Required.Error("name is required"),
			validation.Length(3,100).Error("name must be between 3 and 100 characters"),
			),
			validation.Field(&req.Type,validation.Required.Error("type is required"),
			validation.In("income","expense").Error("type must either 'income' or 'expense ' ")),
			validation.Field(&req.Color,validation.Match(regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)).Error("color must be a valid hex code (e.g., #FF5733)"),
		),
	)
	fieldErrors:=make(map[string]string)

	if err!=nil{
		if errV,ok:=err.(validation.Errors);ok{
			for key,value:=range errV {
				if value!=nil{
					fieldErrors[key]=value.Error()
				}
			}
		}
	}

_, err = v.repo.GetByUserIDAndName(ctx, userID, req.Name)
	if err == nil {
		// Category found - this is a duplicate!
		fieldErrors["name"] = "category name already exists"
	} else {
		// Check if it's a "not found" error (which is good)
		if rErr, ok := err.(*richerror.RichError); ok {
			if rErr.GetKind() != richerror.KindNotFound {
				// Some other error occurred
				return nil, richerror.New(op).
					WithErr(err).
					WithMessage("failed to check category name").
					WithKind(richerror.KindUnexpected)
			}
		}
	}

	if len(fieldErrors)>0{
		return  fieldErrors,richerror.New(op).WithMessage("invalid input").WithKind(richerror.KindInvalid).WithMeta("fields",fieldErrors)
	}

	return  nil, nil
}