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

type AcceptRequestInput struct {
	UserID    string
	RequestID string
}

type AcceptRequestOutput struct {
	Message string            `json:"message"`
	Detail  string            `json:"detail"`
	Request *entities.Request `json:"request"`
}

type AcceptRequestUseCase struct {
	userRepository     entities.UserRepository
	providerRepository entities.ProviderRepository
	requestRepository  entities.RequestRepository
	emailService       mail.EmailService
	logger             logging.LoggerInterface
}

func NewAcceptRequestUseCase(
	userRepository entities.UserRepository,
	providerRepository entities.ProviderRepository,
	requestRepository entities.RequestRepository,
	emailService mail.EmailService,
	logger logging.LoggerInterface,
) *AcceptRequestUseCase {
	return &AcceptRequestUseCase{
		userRepository:     userRepository,
		providerRepository: providerRepository,
		requestRepository:  requestRepository,
		emailService:       emailService,
		logger:             logger,
	}
}

func (uc *AcceptRequestUseCase) Execute(ctx context.Context, input AcceptRequestInput) (*AcceptRequestOutput, []exceptions.ProblemDetails) {
	const from = "AcceptRequestUseCase.Execute"

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
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo provider podem aceitar solicitações"))
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

	request.Accept()

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

	if err := uc.emailService.SendRequestAcceptedEmail(
		owner.Login.Email,
		owner.Name,
		user.Name,
		request.Pet.Name,
		request.ID,
	); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao enviar email de solicitação aceita", err)
	}

	uc.logger.LogInfo(ctx, from, "Solicitação aceita com sucesso")

	return &AcceptRequestOutput{
		Message: "Solicitação aceita com sucesso",
		Detail:  "A solicitação foi aceita e o status foi atualizado para 'accepted'",
		Request: request,
	}, nil
}
