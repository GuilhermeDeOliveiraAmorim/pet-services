package usecases

import (
	"context"
	"errors"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type CreateCategoryInput struct {
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
}

type CreateCategoryOutput struct {
	Message  string             `json:"message,omitempty"`
	Detail   string             `json:"detail,omitempty"`
	Category *entities.Category `json:"category,omitempty"`
}

type CreateCategoryUseCase struct {
	categoryRepo entities.CategoryRepository
	logger       logging.LoggerInterface
}

func NewCreateCategoryUseCase(categoryRepo entities.CategoryRepository, logger logging.LoggerInterface) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (uc *CreateCategoryUseCase) Execute(ctx context.Context, input CreateCategoryInput) (*CreateCategoryOutput, []exceptions.ProblemDetails) {
	const from = "CreateCategoryUseCase.Execute"

	if !input.IsAdmin {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Apenas administradores podem criar categorias."))
	}

	existing, err := uc.categoryRepo.FindByName(input.Name)
	if err == nil && existing != nil {
		return nil, uc.logger.LogConflict(ctx, from, "Categoria já existe", errors.New("Já existe uma categoria com esse nome."))
	}

	category, problems := entities.NewCategory(input.Name)
	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Categoria inválida", problems)
		return nil, problems
	}

	if err := uc.categoryRepo.Create(category); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao criar categoria", err)
	}

	return &CreateCategoryOutput{
		Message:  "Categoria criada com sucesso",
		Detail:   "A categoria foi cadastrada",
		Category: category,
	}, nil
}
