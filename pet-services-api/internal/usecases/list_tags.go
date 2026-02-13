package usecases

import (
	"context"
	"pet-services-api/internal/entities"
)

// ListTagsInput define os parâmetros de entrada para paginação
// Page: número da página (1-indexed)
// PageSize: quantidade de itens por página
// Name: filtro opcional por nome
// ProviderID: ID do usuário provider autenticado
// (ProviderID pode ser usado para validação, não para filtrar tags)
type ListTagsInput struct {
	Page       int
	PageSize   int
	Name       string
	ProviderID string
}

// ListTagsOutput define o resultado da listagem
// Tags: lista de tags
// Total: total de tags disponíveis (para paginação)
type ListTagsOutput struct {
	Tags  []entities.Tag
	Total int
}

// TagRepository define interface para buscar tags
// Precisa implementar ListTagsPaginated e CountTags
// (CountTags pode ser embutido no ListTagsPaginated se preferir)
type TagRepository interface {
	ListTagsPaginated(ctx context.Context, name string, offset, limit int) ([]entities.Tag, error)
	CountTags(ctx context.Context, name string) (int, error)
}

// ListTagsUseCase executa a listagem paginada de tags
// Só providers autenticados podem acessar (validação fora do usecase)
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
