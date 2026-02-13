package usecases

import (
	"context"
	"errors"
	"strings"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type RejectRequestInput struct {
	UserID    string `json:"user_id"`
	RequestID string `json:"request_id"`
	Reason    string `json:"reason"`
}

type RejectRequestInputBody struct {
	Reason string `json:"reason"`
}

type RejectRequestOutput struct {
	Message string            `json:"message"`
	Detail  string            `json:"detail"`
	Request *entities.Request `json:"request"`
}

type RejectRequestUseCase struct {
	userRepository     entities.UserRepository
	providerRepository entities.ProviderRepository
	requestRepository  entities.RequestRepository
	logger             logging.LoggerInterface
}

func NewRejectRequestUseCase(
	userRepository entities.UserRepository,
	providerRepository entities.ProviderRepository,
	requestRepository entities.RequestRepository,
	logger logging.LoggerInterface,
) *RejectRequestUseCase {
	return &RejectRequestUseCase{
		userRepository:     userRepository,
		providerRepository: providerRepository,
		requestRepository:  requestRepository,
		logger:             logger,
	}
}

func (uc *RejectRequestUseCase) Execute(ctx context.Context, input RejectRequestInput) (*RejectRequestOutput, []exceptions.ProblemDetails) {
	const from = "RejectRequestUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.RequestID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID da solicitação ausente", errors.New("O ID da solicitação é obrigatório"))
	}

	reason := strings.TrimSpace(input.Reason)
	if reason == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "Motivo ausente", errors.New("O motivo da rejeição é obrigatório"))
	}

	if len(reason) > 500 {
		return nil, uc.logger.LogBadRequest(ctx, from, "Motivo muito longo", errors.New("O motivo da rejeição deve ter no máximo 500 caracteres"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.IsProvider() {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo provider podem rejeitar solicitações"))
	}

	provider, err := uc.providerRepository.FindByUserID(input.UserID)
	if err != nil {
		if err.Error() == consts.ProviderNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Provedor não encontrado", errors.New("Não foi possível encontrar o provedor do usuário"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar provedor", err)
	}

	request, err := uc.requestRepository.FindByID(input.RequestID)
	if err != nil {
		if err.Error() == consts.RequestNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Solicitação não encontrada", errors.New("Não foi possível encontrar a solicitação informada"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar solicitação", err)
	}

	if request.ProviderID != provider.ID {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("A solicitação não pertence ao provedor"))
	}

	if request.Status != entities.RequestStatuses.Pending {
		return nil, uc.logger.LogConflict(ctx, from, "Solicitação já processada", errors.New("A solicitação não está mais pendente"))
	}

	request.Reject(reason)

	if err := uc.requestRepository.Update(request); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao atualizar solicitação", err)
	}

	uc.logger.LogInfo(ctx, from, "Solicitação rejeitada com sucesso")

	return &RejectRequestOutput{
		Message: "Solicitação rejeitada com sucesso",
		Detail:  "A solicitação foi rejeitada e o status foi atualizado para 'rejected'",
		Request: request,
	}, nil
}
