package service

import (
	"context"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/category/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
)

func (s Service) UpdateCategory(ctx context.Context,req param.UpdateCategoryRequest,categoryID,userID uint)(param.UpdateCategoryResponse,error){
	const op=richerror.Op("categoryservice.UpdateCategory")

	logger.Info("Updating category","category_id",categoryID,"user_id",userID)

	existingCategory,err:=s.repo.GetByID(ctx,categoryID)

	if err!=nil{
		logger.Error("Failed to get category","category_id",categoryID,"error",err.Error())

		return param.UpdateCategoryResponse{},richerror.New(op).WithMessage("category not found").WithKind(richerror.KindNotFound).WithErr(err)
	}

	if existingCategory.UserID !=userID {
		logger.Warn("Unauthorized category update attempt","category_id",categoryID,"category_owner",existingCategory.UserID,"requesting_user",userID)

		return param.UpdateCategoryResponse{},richerror.New(op).WithMessage("unauthorized").WithKind(richerror.KindForbidden)
	}

	if req.Name != nil {
			existingCategory.Name = *req.Name
	}
	if req.Color != nil {
			existingCategory.Color = *req.Color
	}

	updatedCategory,err:=s.repo.Update(ctx,existingCategory)
	
	if err!=nil {
		logger.Error("Failed to update category","category_id",categoryID,"error",err.Error())
		return  param.UpdateCategoryResponse{},richerror.New(op).WithMessage("failed to udpate category").WithKind(richerror.KindUnexpected).WithErr(err)
	}

	logger.Info("Category updated successfully","category_id",categoryID,"user_id",userID)

	return param.UpdateCategoryResponse{
		Category: param.ToCategoryInfo(updatedCategory),
	},nil
}