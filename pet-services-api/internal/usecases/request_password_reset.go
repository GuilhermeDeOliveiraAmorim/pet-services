package usecases

import (
	"context"
	"time"

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
	Message    string    `json:"message,omitempty"`
	Detail     string    `json:"detail,omitempty"`
	ResetToken string    `json:"reset_token,omitempty"`
	ExpiresAt  time.Time `json:"expires_at,omitempty"`
}

type RequestPasswordResetUseCase struct {
	userRepository       entities.UserRepository
	resetTokenRepository entities.RefreshTokenRepository
	emailService         mail.EmailService
	ttl                  time.Duration
}

func NewRequestPasswordResetUseCase(userRepo entities.UserRepository, resetRepo entities.RefreshTokenRepository, emailService mail.EmailService, ttl time.Duration) *RequestPasswordResetUseCase {
	if ttl <= 0 {
		ttl = time.Hour
	}
	return &RequestPasswordResetUseCase{
		userRepository:       userRepo,
		resetTokenRepository: resetRepo,
		emailService:         emailService,
		ttl:                  ttl,
	}
}

func (uc *RequestPasswordResetUseCase) Execute(ctx context.Context, input RequestPasswordResetInput) (*RequestPasswordResetOutput, []exceptions.ProblemDetails) {
	const from = "RequestPasswordResetUseCase.Execute"

	if input.Email == "" {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "Email ausente",
				Detail: "O email é obrigatório para reset de senha",
			}),
		}
	}

	user, err := uc.userRepository.FindByEmail(input.Email)
	if err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	tokenStr := ulid.Make().String()
	expiresAt := time.Now().Add(uc.ttl)

	if err := uc.resetTokenRepository.RevokeAllPasswordResetByUserID(user.ID); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao revogar tokens anteriores", err)
	}

	resetToken := &entities.PasswordResetToken{
		Token:     tokenStr,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
		UserAgent: input.UserAgent,
		IP:        input.IP,
	}

	if err := uc.resetTokenRepository.CreatePasswordReset(resetToken); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao salvar token de reset", err)
	}

	if err := uc.emailService.SendPasswordResetEmail(user.Login.Email, tokenStr); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao enviar email de reset", err)
	}

	return &RequestPasswordResetOutput{
		Message:    "Instruções enviadas",
		Detail:     "Se o email existir, um link de reset foi gerado",
		ResetToken: tokenStr,
		ExpiresAt:  expiresAt,
	}, nil
}
