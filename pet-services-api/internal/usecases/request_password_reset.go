package usecases

import (
	"context"
	"errors"
	"time"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/mail"

	"github.com/oklog/ulid/v2"
)

type RequestPasswordResetInput struct {
	Email     string `json:"email"`
	UserAgent string `json:"user_agent,omitempty"`
	IP        string `json:"ip,omitempty"`
}

type RequestPasswordResetOutput struct {
	Message   string    `json:"message,omitempty"`
	Detail    string    `json:"detail,omitempty"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

type RequestPasswordResetUseCase struct {
	userRepository       entities.UserRepository
	resetTokenRepository entities.RefreshTokenRepository
	emailService         mail.EmailService
	ttl                  time.Duration
	logger               logging.LoggerInterface
}

func NewRequestPasswordResetUseCase(userRepo entities.UserRepository, resetRepo entities.RefreshTokenRepository, emailService mail.EmailService, ttl time.Duration, logger logging.LoggerInterface) *RequestPasswordResetUseCase {
	if ttl <= 0 {
		ttl = time.Hour
	}
	return &RequestPasswordResetUseCase{
		userRepository:       userRepo,
		resetTokenRepository: resetRepo,
		emailService:         emailService,
		ttl:                  ttl,
		logger:               logger,
	}
}

func (uc *RequestPasswordResetUseCase) Execute(ctx context.Context, input RequestPasswordResetInput) (*RequestPasswordResetOutput, []exceptions.ProblemDetails) {
	const from = "RequestPasswordResetUseCase.Execute"

	if input.Email == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "Email ausente", errors.New("O email é obrigatório para redefinição de senha"))
	}

	user, err := uc.userRepository.FindByEmail(input.Email)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar um usuário com o email informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	tokenStr := ulid.Make().String()
	expiresAt := time.Now().Add(uc.ttl)

	if err := uc.resetTokenRepository.RevokeAllPasswordResetByUserID(user.ID); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao revogar tokens anteriores", err)
	}

	resetToken := &entities.PasswordResetToken{
		Token:     tokenStr,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
		UserAgent: input.UserAgent,
		IP:        input.IP,
	}

	if err := uc.resetTokenRepository.CreatePasswordReset(resetToken); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao salvar token de redefinição", err)
	}

	if err := uc.emailService.SendPasswordResetEmail(user.Login.Email, tokenStr); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao enviar email de redefinição", err)
	}

	return &RequestPasswordResetOutput{
		Message:   "Instruções enviadas",
		Detail:    "Se o email existir, um link de redefinição foi gerado",
		ExpiresAt: expiresAt,
	}, nil
}
