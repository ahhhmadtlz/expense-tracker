package param

import "github.com/ahhhmadtlz/expense-tracker/internal/domain/category/entity"

type CreateCategoryRequest struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Color string `json:"color"`
}

type CreateCategoryResponse struct {
	Category CategoryInfo `json:"category"`
}

type UpdateCategoryRequest struct {
	Name  *string `json:"name"`
	Color *string `json:"color"`
}

type UpdateCategoryResponse struct {
	Category CategoryInfo `json:"category"`
}


type GetCategoryResponse struct {
	Category CategoryInfo `json:"category"`
}

type ListCategoriesRequest struct {
	Type string `json:"type"`
}

type ListCategoriesResponse struct {
	Categories []CategoryInfo `json:"categories"`
}


type DeleteCategoryResponse struct {
	Message string `json:"message"`
}

type CategoryInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Color string `json:"color"`
}

func ToCategoryInfo(cat entity.Category) CategoryInfo {
	return CategoryInfo{
		ID:    cat.ID,
		Name:  cat.Name,
		Type:  cat.Type.String(),
		Color: cat.Color,
	}
}