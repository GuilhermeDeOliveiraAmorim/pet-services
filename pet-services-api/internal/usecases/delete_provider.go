package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type DeleteProviderInput struct {
	UserID     string `json:"user_id"`
	ProviderID string `json:"provider_id"`
}

type DeleteProviderOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type DeleteProviderUseCase struct {
	userRepository     entities.UserRepository
	providerRepository entities.ProviderRepository
	logger             logging.LoggerInterface
}

func NewDeleteProviderUseCase(
	userRepository entities.UserRepository,
	providerRepository entities.ProviderRepository,
	logger logging.LoggerInterface,
) *DeleteProviderUseCase {
	return &DeleteProviderUseCase{
		userRepository:     userRepository,
		providerRepository: providerRepository,
		logger:             logger,
	}
}

func (uc *DeleteProviderUseCase) Execute(ctx context.Context, input DeleteProviderInput) (*DeleteProviderOutput, []exceptions.ProblemDetails) {
	const from = "DeleteProviderUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.ProviderID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do provedor ausente", errors.New("O ID do provedor é obrigatório"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.IsProvider() {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo provider podem remover provedores"))
	}

	provider, err := uc.providerRepository.FindByID(input.ProviderID)
	if err != nil {
		if err.Error() == consts.ProviderNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Provedor não encontrado", errors.New("Não foi possível encontrar o provedor informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar provedor", err)
	}

	if provider.UserID != user.ID {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("O provedor não pertence ao usuário autenticado"))
	}

	if !provider.Active {
		return nil, uc.logger.LogBadRequest(ctx, from, "Provedor já inativo", errors.New("O provedor já está inativo no sistema"))
	}

	provider.Deactivate()

	if err := uc.providerRepository.Update(provider); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao remover provedor", err)
	}

	return &DeleteProviderOutput{
		Message: "Provedor removido com sucesso",
		Detail:  "O provedor foi removido do sistema com sucesso",
	}, nil
}
