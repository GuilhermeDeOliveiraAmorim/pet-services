package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type RemoveServiceTagInput struct {
	UserID    string `json:"user_id"`
	ServiceID string `json:"service_id"`
	TagID     string `json:"tag_id"`
}

type RemoveServiceTagOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type RemoveServiceTagUseCase struct {
	userRepository     entities.UserRepository
	serviceRepository  entities.ServiceRepository
	providerRepository entities.ProviderRepository
	tagRepository      entities.TagRepository
	logger             logging.LoggerInterface
}

func NewRemoveServiceTagUseCase(
	userRepository entities.UserRepository,
	serviceRepository entities.ServiceRepository,
	providerRepository entities.ProviderRepository,
	tagRepository entities.TagRepository,
	logger logging.LoggerInterface,
) *RemoveServiceTagUseCase {
	return &RemoveServiceTagUseCase{
		userRepository:     userRepository,
		serviceRepository:  serviceRepository,
		providerRepository: providerRepository,
		tagRepository:      tagRepository,
		logger:             logger,
	}
}

func (uc *RemoveServiceTagUseCase) Execute(ctx context.Context, input RemoveServiceTagInput) (*RemoveServiceTagOutput, []exceptions.ProblemDetails) {
	const from = "RemoveServiceTagUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.ServiceID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do serviço ausente", errors.New("O ID do serviço é obrigatório"))
	}

	if input.TagID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID da tag ausente", errors.New("O ID da tag é obrigatório"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.IsProvider() {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo provider podem remover tags do serviço"))
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

	tag, err := uc.tagRepository.FindByID(input.TagID)
	if err != nil {
		if err.Error() == consts.TagNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Tag não encontrada", errors.New("Não foi possível encontrar a tag informada"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar tag", err)
	}

	if !tag.Active {
		return nil, uc.logger.LogBadRequest(ctx, from, "Tag inativa", errors.New("A tag informada está inativa"))
	}

	exists, err := uc.serviceRepository.HasTag(service.ID, tag.ID)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao verificar tags do serviço", err)
	}
	if !exists {
		return nil, uc.logger.LogConflict(ctx, from, "Tag não vinculada", errors.New("A tag informada não está vinculada ao serviço"))
	}

	if err := uc.serviceRepository.RemoveTag(service.ID, tag.ID); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao remover tag do serviço", err)
	}

	return &RemoveServiceTagOutput{
		Message: "Tag removida com sucesso",
		Detail:  "A tag foi desassociada do serviço",
	}, nil
}
