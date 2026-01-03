package user

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// RequestEmailVerificationUseCase envia email de verificação.
type RequestEmailVerificationUseCase struct {
	userRepo            domainUser.Repository
	emailService        domainUser.EmailService
	verificationBaseURL string
	logger              *slog.Logger
}

func NewRequestEmailVerificationUseCase(userRepo domainUser.Repository, emailService domainUser.EmailService, verificationBaseURL string, logger *slog.Logger) *RequestEmailVerificationUseCase {
	return &RequestEmailVerificationUseCase{
		userRepo:            userRepo,
		emailService:        emailService,
		verificationBaseURL: verificationBaseURL,
		logger:              logging.EnsureLogger(logger),
	}
}

// RequestEmailVerificationInput dados para solicitar verificação.
type RequestEmailVerificationInput struct {
	UserID uuid.UUID
}

// Execute gera token, persiste e envia email de verificação.
func (uc *RequestEmailVerificationUseCase) Execute(ctx context.Context, input RequestEmailVerificationInput) error {
	var err error
	defer logging.UseCase(ctx, uc.logger, "RequestEmailVerificationUseCase", slog.String("user_id", input.UserID.String()))(&err)

	if input.UserID == uuid.Nil {
		err = domainUser.ErrUserNotFound
		return err
	}

	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	if user.EmailVerified {
		err = domainUser.ErrEmailAlreadyVerified
		return err
	}

	// Gera token seguro
	token, err := generateSecureToken(32)
	if err != nil {
		err = fmt.Errorf("falha ao gerar token: %w", err)
		return err
	}

	// Token válido por 24 horas
	expiresAt := time.Now().Add(24 * time.Hour)

	verificationToken := domainUser.NewEmailVerificationToken(user.ID, token, expiresAt)
	if err := uc.userRepo.CreateEmailVerificationToken(ctx, verificationToken); err != nil {
		err = fmt.Errorf("falha ao salvar token de verificação: %w", err)
		return err
	}

	// Envia email com link de verificação
	verificationLink := fmt.Sprintf("%s?token=%s", uc.verificationBaseURL, token)
	if err := uc.emailService.SendEmailVerification(user.Email.String(), verificationLink); err != nil {
		err = fmt.Errorf("falha ao enviar email: %w", err)
		return err
	}

	return nil
}
