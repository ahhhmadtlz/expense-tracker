package service

import (
	"context"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/category/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
)

func (s Service) DeleteCategory(ctx context.Context,categoryID,userID uint)(param.DeleteCategoryResponse,error){

	const op=richerror.Op("categoryService.DeleteCategory")

	s.logger.Info("Deleting category","category_id",categoryID,"user_id",userID)

	existingCategory,err:=s.repo.GetByID(ctx,categoryID)
	if err !=nil{
		s.logger.Error("Failed to get category","category_id",categoryID,"error",err.Error())
		return  param.DeleteCategoryResponse{},richerror.New(op).WithMessage("category not found").WithKind(richerror.KindNotFound).WithErr(err)
	}
	if existingCategory.UserID !=userID{
			s.logger.Warn("Unauthorized category deletion attempt",
			"category_id", categoryID,
			"category_owner", existingCategory.UserID,
			"requesting_user", userID,
		)
		return param.DeleteCategoryResponse{}, richerror.New(op).
			WithMessage("unauthorized").
			WithKind(richerror.KindForbidden)
	}

	hasTransaction,err:=s.repo.CategoryHasTransactions(ctx,categoryID)

	if err !=nil{
		s.logger.Error("Failed to check category transactions","category_id",categoryID,"error",err.Error())
		return param.DeleteCategoryResponse{},richerror.New(op).WithMessage("failedto check category usage").WithKind(richerror.KindUnexpected).WithErr(err)
	}

	if hasTransaction {
		s.logger.Warn("Attempted to delete category with transactions","category_id",categoryID,"user_id",userID)
		return  param.DeleteCategoryResponse{},richerror.New(op).WithMessage("cannot delete category with existing transactions ").WithKind(richerror.KindInvalid)
	}

	if err :=s.repo.Delete(ctx,categoryID);err!=nil{
		s.logger.Error("Failed to delete category","category_id",categoryID,"error",err.Error())
		return param.DeleteCategoryResponse{},richerror.New(op).WithMessage("failed to delete category").WithKind(richerror.KindUnexpected).WithErr(err)
	}

	s.logger.Info("category deleted successfully","category_id",categoryID,"user_id",userID)

	return param.DeleteCategoryResponse{
		Message: "category deleted successfully",
	},nil
}