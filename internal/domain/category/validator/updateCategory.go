package validator

import (
	"context"
	"regexp"
	"strings"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/category/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
	validation "github.com/go-ozzo/ozzo-validation"
)

func (v Validator) ValidateUpdateCategory(ctx context.Context,req param.UpdateCategoryRequest)(map[string]string ,error){
	const op="categoryvalidator.ValidateUpdateCategory"
	fieldErrors:=make(map[string]string)
	
	if req.Name !=nil{
		*req.Name=strings.TrimSpace(*req.Name)
		err:=validation.Validate(req.Name,
			validation.Required.Error("name cannot be empty"),
			validation.Length(3,100).Error("name must be betwewn 3 and 100 characters"),
		)
		if err!=nil{
			fieldErrors["name"]=err.Error()
		}
	}

	if req.Color !=nil{
		*req.Color=strings.TrimSpace(*req.Color)
		err:=validation.Validate(req.Color,
		validation.Match(regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)).Error("color must be valid hex (e.g , #ff3300)"),
		)

		if err!=nil{
			fieldErrors["color"]=err.Error()
		}
	}

	if len(fieldErrors)>0 {
		return fieldErrors ,richerror.New(op).WithMessage("invalid input").WithKind(richerror.KindInvalid).WithMeta("fields",fieldErrors)
	}

	return  nil,nil
}