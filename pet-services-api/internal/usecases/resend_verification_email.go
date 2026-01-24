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

type ResendVerificationEmailInput struct {
	Email     string `json:"email"`
	UserAgent string `json:"user_agent,omitempty"`
	IP        string `json:"ip,omitempty"`
}

type ResendVerificationEmailOutput struct {
	Message     string    `json:"message,omitempty"`
	Detail      string    `json:"detail,omitempty"`
	VerifyToken string    `json:"verify_token,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
}

type ResendVerificationEmailUseCase struct {
	userRepository        entities.UserRepository
	verifyTokenRepository entities.RefreshTokenRepository
	emailService          mail.EmailService
	ttl                   time.Duration
}

func NewResendVerificationEmailUseCase(userRepo entities.UserRepository, verifyRepo entities.RefreshTokenRepository, emailService mail.EmailService, ttl time.Duration) *ResendVerificationEmailUseCase {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &ResendVerificationEmailUseCase{
		userRepository:        userRepo,
		verifyTokenRepository: verifyRepo,
		emailService:          emailService,
		ttl:                   ttl,
	}
}

func (uc *ResendVerificationEmailUseCase) Execute(ctx context.Context, input ResendVerificationEmailInput) (*ResendVerificationEmailOutput, []exceptions.ProblemDetails) {
	const from = "ResendVerificationEmailUseCase.Execute"

	if input.Email == "" {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "Email ausente",
				Detail: "O email é obrigatório para reenviar verificação",
			}),
		}
	}

	user, err := uc.userRepository.FindByEmail(input.Email)
	if err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if user.EmailVerified {
		return &ResendVerificationEmailOutput{
			Message: "Email já verificado",
			Detail:  "Este email já foi verificado anteriormente",
		}, nil
	}

	tokenStr := ulid.Make().String()
	expiresAt := time.Now().Add(uc.ttl)

	if err := uc.verifyTokenRepository.RevokeAllPasswordResetByUserID(user.ID); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao revogar tokens anteriores", err)
	}

	verifyToken := &entities.PasswordResetToken{
		Token:     tokenStr,
		UserID:    user.ID,
		ExpiresAt: expiresAt,
		UserAgent: input.UserAgent,
		IP:        input.IP,
	}

	if err := uc.verifyTokenRepository.CreatePasswordReset(verifyToken); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao salvar token de verificação", err)
	}

	if err := uc.emailService.SendVerificationEmail(user.Login.Email, tokenStr); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao enviar email de verificação", err)
	}

	return &ResendVerificationEmailOutput{
		Message:     "Email de verificação reenviado",
		Detail:      "Verifique sua caixa de entrada para completar a verificação",
		VerifyToken: tokenStr,
		ExpiresAt:   expiresAt,
	}, nil
}
