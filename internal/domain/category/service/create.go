package service

import (
	"context"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/category/entity"
	"github.com/ahhhmadtlz/expense-tracker/internal/domain/category/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
)

func (s Service) CreateCategory(ctx context.Context,req param.CreateCategoryRequest,userID uint)(param.CreateCategoryResponse,error){
	const op=richerror.Op("categoryService.CreateCategory")

	logger.Info("Creating category","user_id",userID,"name",req.Name,"type",req.Type)


	catType:=entity.MapToCategoryType(req.Type)

	category:=entity.Category{
		UserID: userID,
		Name: req.Name,
		Type: catType,
		Color: req.Color,
	}

	createdCategory,err:=s.repo.Create(ctx,category)

	if err!=nil{
		logger.Error("Failed to create category",
		"user_id",userID,
		"name",req.Name,
		"error",err.Error())

		return param.CreateCategoryResponse{},richerror.New(op).WithMessage("failed to create category").WithKind(richerror.KindUnexpected).WithErr(err)
	}
	logger.Info("Category created successfully","category_id",createdCategory.ID,"user_id",userID)

	return param.CreateCategoryResponse{
		Category: param.ToCategoryInfo(createdCategory),
	},nil

}