package service

import (
	"context"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/category/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
)

func (s Service) GetCategory(ctx context.Context,categoryID,userID uint)(param.GetCategoryResponse,error){
	const op=richerror.Op("categoryService.GetCategory")

	s.logger.Debug("Getting category","category_id",categoryID,"user_id",userID)

	category,err:=s.repo.GetByID(ctx,categoryID)

	if err!=nil{
		s.logger.Error("failed to get category","category_id",categoryID,"error",err.Error())
		return param.GetCategoryResponse{},richerror.New(op).WithMessage("category not found").WithKind(richerror.KindNotFound).WithErr(err)
	}

	if category.UserID!=userID {
			s.logger.Warn("Unauthorized category access attempt",
			"category_id", categoryID,
			"category_owner", category.UserID,
			"requesting_user", userID,
		)
		return param.GetCategoryResponse{}, richerror.New(op).
			WithMessage("unauthorized").
			WithKind(richerror.KindForbidden)
	}

	return  param.GetCategoryResponse{
		Category: param.ToCategoryInfo(category),
	},nil
}