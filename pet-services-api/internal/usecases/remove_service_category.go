package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type RemoveServiceCategoryInput struct {
	UserID     string `json:"user_id"`
	ServiceID  string `json:"service_id"`
	CategoryID string `json:"category_id"`
}

type RemoveServiceCategoryOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type RemoveServiceCategoryUseCase struct {
	userRepository     entities.UserRepository
	serviceRepository  entities.ServiceRepository
	providerRepository entities.ProviderRepository
	categoryRepository entities.CategoryRepository
	logger             logging.LoggerInterface
}

func NewRemoveServiceCategoryUseCase(
	userRepository entities.UserRepository,
	serviceRepository entities.ServiceRepository,
	providerRepository entities.ProviderRepository,
	categoryRepository entities.CategoryRepository,
	logger logging.LoggerInterface,
) *RemoveServiceCategoryUseCase {
	return &RemoveServiceCategoryUseCase{
		userRepository:     userRepository,
		serviceRepository:  serviceRepository,
		providerRepository: providerRepository,
		categoryRepository: categoryRepository,
		logger:             logger,
	}
}

func (uc *RemoveServiceCategoryUseCase) Execute(ctx context.Context, input RemoveServiceCategoryInput) (*RemoveServiceCategoryOutput, []exceptions.ProblemDetails) {
	const from = "RemoveServiceCategoryUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.ServiceID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do serviço ausente", errors.New("O ID do serviço é obrigatório"))
	}

	if input.CategoryID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID da categoria ausente", errors.New("O ID da categoria é obrigatório"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.IsProvider() {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo provider podem remover categorias do serviço"))
	}

	provider, err := uc.providerRepository.FindByUserID(user.ID)
	if err != nil {
		if err.Error() == consts.ProviderNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Provedor não encontrado", errors.New("Não foi possível encontrar o provedor do usuário"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar provedor", err)
	}

	service, err := uc.serviceRepository.FindByID(input.ServiceID)
	if err != nil {
		if err.Error() == consts.ServiceNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Serviço não encontrado", errors.New("Não foi possível encontrar o serviço informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar serviço", err)
	}

	if service.ProviderID != provider.ID {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("O serviço informado não pertence ao provedor autenticado"))
	}

	if !service.Active {
		return nil, uc.logger.LogBadRequest(ctx, from, "Serviço inativo", errors.New("O serviço informado está inativo"))
	}

	category, err := uc.categoryRepository.FindByID(input.CategoryID)
	if err != nil {
		if err.Error() == consts.CategoryNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Categoria não encontrada", errors.New("Não foi possível encontrar a categoria informada"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar categoria", err)
	}

	if !category.Active {
		return nil, uc.logger.LogBadRequest(ctx, from, "Categoria inativa", errors.New("A categoria informada está inativa"))
	}

	exists, err := uc.serviceRepository.HasCategory(service.ID, category.ID)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao verificar categorias do serviço", err)
	}
	if !exists {
		return nil, uc.logger.LogConflict(ctx, from, "Categoria não vinculada", errors.New("A categoria informada não está vinculada ao serviço"))
	}

	if err := uc.serviceRepository.RemoveCategory(service.ID, category.ID); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao remover categoria do serviço", err)
	}

	return &RemoveServiceCategoryOutput{
		Message: "Categoria removida com sucesso",
		Detail:  "A categoria foi desassociada do serviço",
	}, nil
}
