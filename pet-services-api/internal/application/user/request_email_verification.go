package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// RequestEmailVerificationUseCase envia email de verificação.
type RequestEmailVerificationUseCase struct {
	userRepo            domainUser.Repository
	emailService        domainUser.EmailService
	verificationBaseURL string
}

func NewRequestEmailVerificationUseCase(userRepo domainUser.Repository, emailService domainUser.EmailService, verificationBaseURL string) *RequestEmailVerificationUseCase {
	return &RequestEmailVerificationUseCase{
		userRepo:            userRepo,
		emailService:        emailService,
		verificationBaseURL: verificationBaseURL,
	}
}

// RequestEmailVerificationInput dados para solicitar verificação.
type RequestEmailVerificationInput struct {
	UserID uuid.UUID
}

// Execute gera token, persiste e envia email de verificação.
func (uc *RequestEmailVerificationUseCase) Execute(ctx context.Context, input RequestEmailVerificationInput) error {
	if input.UserID == uuid.Nil {
		return domainUser.ErrUserNotFound
	}

	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	if user.EmailVerified {
		return domainUser.ErrEmailAlreadyVerified
	}

	// Gera token seguro
	token, err := generateSecureToken(32)
	if err != nil {
		return fmt.Errorf("falha ao gerar token: %w", err)
	}

	// Token válido por 24 horas
	expiresAt := time.Now().Add(24 * time.Hour)

	verificationToken := domainUser.NewEmailVerificationToken(user.ID, token, expiresAt)
	if err := uc.userRepo.CreateEmailVerificationToken(ctx, verificationToken); err != nil {
		return fmt.Errorf("falha ao salvar token de verificação: %w", err)
	}

	// Envia email com link de verificação
	verificationLink := fmt.Sprintf("%s?token=%s", uc.verificationBaseURL, token)
	if err := uc.emailService.SendEmailVerification(user.Email.String(), verificationLink); err != nil {
		return fmt.Errorf("falha ao enviar email: %w", err)
	}

	return nil
}
