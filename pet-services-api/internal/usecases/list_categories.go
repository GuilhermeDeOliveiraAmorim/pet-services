package usecases

import (
	"context"
	"pet-services-api/internal/entities"
)

type ListCategoriesInput struct {
	Page       int
	PageSize   int
	Name       string
	ProviderID string
}

type ListCategoriesOutput struct {
	Categories []entities.Category `json:"categories"`
	Total      int                 `json:"total"`
}

type ListCategoriesUseCase struct {
	categoryRepo entities.CategoryRepository
}

func NewListCategoriesUseCase(categoryRepo entities.CategoryRepository) *ListCategoriesUseCase {
	return &ListCategoriesUseCase{categoryRepo: categoryRepo}
}

func (uc *ListCategoriesUseCase) Execute(ctx context.Context, input ListCategoriesInput) (ListCategoriesOutput, error) {
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 {
		input.PageSize = 10
	}
	offset := (input.Page - 1) * input.PageSize

	categories, err := uc.categoryRepo.ListCategoriesPaginated(ctx, input.Name, offset, input.PageSize)
	if err != nil {
		return ListCategoriesOutput{}, err
	}
	total, err := uc.categoryRepo.CountCategories(ctx, input.Name)
	if err != nil {
		return ListCategoriesOutput{}, err
	}
	return ListCategoriesOutput{Categories: categories, Total: total}, nil
}
