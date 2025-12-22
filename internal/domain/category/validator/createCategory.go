package validator

import (
	"context"
	"regexp"
	"strings"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/category/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation"
)

func (v Validator) ValidateCreateCategory(ctx context.Context, req param.CreateCategoryRequest)(map[string]string,error) {
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

	if len(fieldErrors)>0{
		return  fieldErrors,richerror.New(op).WithMessage("invalid input").WithKind(richerror.KindInvalid).WithMeta("fields",fieldErrors)
	}

	return  nil, nil
}