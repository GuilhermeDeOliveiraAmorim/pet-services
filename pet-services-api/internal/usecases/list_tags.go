package usecases

import (
	"context"
	"pet-services-api/internal/entities"
)

type ListTagsInput struct {
	Page       int
	PageSize   int
	Name       string
	ProviderID string
}

type ListTagsOutput struct {
	Tags  []entities.Tag
	Total int
}

type TagRepository interface {
	ListTagsPaginated(ctx context.Context, name string, offset, limit int) ([]entities.Tag, error)
	CountTags(ctx context.Context, name string) (int, error)
}

type ListTagsUseCase struct {
	tagRepo TagRepository
}

func NewListTagsUseCase(tagRepo TagRepository) *ListTagsUseCase {
	return &ListTagsUseCase{tagRepo: tagRepo}
}

func (uc *ListTagsUseCase) Execute(ctx context.Context, input ListTagsInput) (ListTagsOutput, error) {
	if input.Page < 1 {
		input.Page = 1
	}
	if input.PageSize < 1 {
		input.PageSize = 10
	}
	offset := (input.Page - 1) * input.PageSize

	tags, err := uc.tagRepo.ListTagsPaginated(ctx, input.Name, offset, input.PageSize)
	if err != nil {
		return ListTagsOutput{}, err
	}
	total, err := uc.tagRepo.CountTags(ctx, input.Name)
	if err != nil {
		return ListTagsOutput{}, err
	}
	return ListTagsOutput{Tags: tags, Total: total}, nil
}
