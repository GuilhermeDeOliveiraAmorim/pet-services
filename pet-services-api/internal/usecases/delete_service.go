package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type DeleteServiceInput struct {
	UserID    string `json:"user_id"`
	ServiceID string `json:"service_id"`
}

type DeleteServiceOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type DeleteServiceUseCase struct {
	userRepository     entities.UserRepository
	providerRepository entities.ProviderRepository
	serviceRepository  entities.ServiceRepository
	logger             logging.LoggerInterface
}

func NewDeleteServiceUseCase(
	userRepository entities.UserRepository,
	providerRepository entities.ProviderRepository,
	serviceRepository entities.ServiceRepository,
	logger logging.LoggerInterface,
) *DeleteServiceUseCase {
	return &DeleteServiceUseCase{
		userRepository:     userRepository,
		providerRepository: providerRepository,
		serviceRepository:  serviceRepository,
		logger:             logger,
	}
}

func (uc *DeleteServiceUseCase) Execute(ctx context.Context, input DeleteServiceInput) (*DeleteServiceOutput, []exceptions.ProblemDetails) {
	const from = "DeleteServiceUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.ServiceID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do serviço ausente", errors.New("O ID do serviço é obrigatório"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.IsProvider() {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo provider podem remover serviços"))
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
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("O serviço não pertence ao provedor autenticado"))
	}

	if !service.Active {
		return nil, uc.logger.LogBadRequest(ctx, from, "Serviço já inativo", errors.New("O serviço já está inativo no sistema"))
	}

	service.Deactivate()

	if err := uc.serviceRepository.Update(service); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao remover serviço", err)
	}

	return &DeleteServiceOutput{
		Message: "Serviço removido com sucesso",
		Detail:  "O serviço foi removido do sistema com sucesso",
	}, nil
}
