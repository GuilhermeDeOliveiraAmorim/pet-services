package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/mail"
)

type CompleteRequestInput struct {
	UserID    string
	RequestID string
}

type CompleteRequestOutput struct {
	Message string            `json:"message"`
	Detail  string            `json:"detail"`
	Request *entities.Request `json:"request"`
}

type CompleteRequestUseCase struct {
	userRepository     entities.UserRepository
	providerRepository entities.ProviderRepository
	requestRepository  entities.RequestRepository
	emailService       mail.EmailService
	logger             logging.LoggerInterface
}

func NewCompleteRequestUseCase(
	userRepository entities.UserRepository,
	providerRepository entities.ProviderRepository,
	requestRepository entities.RequestRepository,
	emailService mail.EmailService,
	logger logging.LoggerInterface,
) *CompleteRequestUseCase {
	return &CompleteRequestUseCase{
		userRepository:     userRepository,
		providerRepository: providerRepository,
		requestRepository:  requestRepository,
		emailService:       emailService,
		logger:             logger,
	}
}

func (uc *CompleteRequestUseCase) Execute(ctx context.Context, input CompleteRequestInput) (*CompleteRequestOutput, []exceptions.ProblemDetails) {
	const from = "CompleteRequestUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.RequestID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID da solicitação ausente", errors.New("O ID da solicitação é obrigatório"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.IsProvider() {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo provider podem concluir solicitações"))
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

	if request.Status != entities.RequestStatuses.Accepted {
		return nil, uc.logger.LogConflict(ctx, from, "Solicitação não aceita", errors.New("A solicitação deve estar com status 'accepted' para ser concluída"))
	}

	request.Complete()

	if err := uc.requestRepository.Update(request); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao atualizar solicitação", err)
	}

	owner, err := uc.userRepository.FindByID(request.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário dono da solicitação não encontrado", errors.New("Não foi possível encontrar o usuário dono da solicitação"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar dono da solicitação", err)
	}

	if err := uc.emailService.SendRequestCompletedEmail(
		owner.Login.Email,
		owner.Name,
		user.Name,
		request.Pet.Name,
		request.ID,
	); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao enviar email de solicitação concluída", err)
	}

	uc.logger.LogInfo(ctx, from, "Solicitação concluída com sucesso")

	return &CompleteRequestOutput{
		Message: "Solicitação concluída com sucesso",
		Detail:  "A solicitação foi concluída e o status foi atualizado para 'completed'",
		Request: request,
	}, nil
}
