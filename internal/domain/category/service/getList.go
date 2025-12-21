package service

import (
	"context"

	"github.com/ahhhmadtlz/expense-tracker/internal/domain/category/entity"
	"github.com/ahhhmadtlz/expense-tracker/internal/domain/category/param"
	"github.com/ahhhmadtlz/expense-tracker/internal/observability/logger"
	"github.com/ahhhmadtlz/expense-tracker/internal/pkg/richerror"
)

func (s Service) ListCategories(ctx context.Context, userID uint, catType string) (param.ListCategoriesResponse,error) {
	const op=richerror.Op("categoryService.ListCategories")

	logger.Debug("Listing categories","user_id",userID,"type_filter",catType)

	var categories []entity.Category
	var err error

	if catType !=""{
		categoryType:=entity.MapToCategoryType(catType)
		categories,err=s.repo.GetByUserIDAndType(ctx,userID,categoryType)
	}else{
		categories,err=s.repo.GetByUserID(ctx,userID)
	}

	if err!=nil{
		logger.Error("Failed to lsit categories","user_id",userID,"error",err.Error())
		return param.ListCategoriesResponse{},richerror.New(op).WithMessage("failed to list categories").WithKind(richerror.KindUnexpected).WithErr(err)
	}

	categoryInfos:=make([]param.CategoryInfo,0,len(categories))
	for _,cat:=range categories {
		categoryInfos=append(categoryInfos, param.ToCategoryInfo(cat))
	}

	logger.Debug("Categories listed successfully","user_id",userID,"count",len(categoryInfos))

	return param.ListCategoriesResponse{
		Categories: categoryInfos,
	},nil

}