package usecases

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

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
	reviewRepository   entities.ReviewRepository
	emailService       mail.EmailService
	logger             logging.LoggerInterface
}

var reviewReminderDelay = resolveReviewReminderDelay()

func resolveReviewReminderDelay() time.Duration {
	const defaultDelay = 24 * time.Hour

	raw := strings.TrimSpace(os.Getenv("REVIEW_REMINDER_DELAY_MINUTES"))
	if raw == "" {
		return defaultDelay
	}

	minutes, err := strconv.Atoi(raw)
	if err != nil || minutes <= 0 {
		return defaultDelay
	}

	return time.Duration(minutes) * time.Minute
}

func NewCompleteRequestUseCase(
	userRepository entities.UserRepository,
	providerRepository entities.ProviderRepository,
	requestRepository entities.RequestRepository,
	reviewRepository entities.ReviewRepository,
	emailService mail.EmailService,
	logger logging.LoggerInterface,
) *CompleteRequestUseCase {
	return &CompleteRequestUseCase{
		userRepository:     userRepository,
		providerRepository: providerRepository,
		requestRepository:  requestRepository,
		reviewRepository:   reviewRepository,
		emailService:       emailService,
		logger:             logger,
	}
}

func (uc *CompleteRequestUseCase) scheduleReviewReminder(owner *entities.User, providerName, petName, requestID, providerID string) {
	if owner == nil || owner.Login.Email == "" {
		return
	}

	ownerID := owner.ID
	ownerName := owner.Name
	ownerEmail := owner.Login.Email

	go func() {
		const from = "CompleteRequestUseCase.scheduleReviewReminder"

		timer := time.NewTimer(reviewReminderDelay)
		defer timer.Stop()

		<-timer.C

		reviews, total, err := uc.reviewRepository.List(providerID, ownerID, 1, 1)
		if err != nil {
			uc.logger.LogWarning(context.Background(), from, "Erro ao verificar reviews para lembrete", err)
			return
		}

		if total > 0 || len(reviews) > 0 {
			return
		}

		if err := uc.emailService.SendReviewReminderEmail(ownerEmail, ownerName, providerName, petName, requestID); err != nil {
			uc.logger.LogWarning(context.Background(), from, "Erro ao enviar lembrete de avaliação", err)
		}
	}()
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

	uc.scheduleReviewReminder(owner, user.Name, request.Pet.Name, request.ID, provider.ID)

	uc.logger.LogInfo(ctx, from, "Solicitação concluída com sucesso")

	return &CompleteRequestOutput{
		Message: "Solicitação concluída com sucesso",
		Detail:  "A solicitação foi concluída e o status foi atualizado para 'completed'",
		Request: request,
	}, nil
}
