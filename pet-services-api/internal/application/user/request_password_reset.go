package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// RequestPasswordResetUseCase inicia o fluxo de reset de senha.
type RequestPasswordResetUseCase struct {
	userRepo     domainUser.Repository
	emailService domainUser.EmailService
	resetBaseURL string
}

func NewRequestPasswordResetUseCase(userRepo domainUser.Repository, emailService domainUser.EmailService, resetBaseURL string) *RequestPasswordResetUseCase {
	return &RequestPasswordResetUseCase{
		userRepo:     userRepo,
		emailService: emailService,
		resetBaseURL: resetBaseURL,
	}
}

// RequestPasswordResetInput dados para solicitar reset.
type RequestPasswordResetInput struct {
	Email string
}

// Execute gera token, persiste e envia email.
func (uc *RequestPasswordResetUseCase) Execute(ctx context.Context, input RequestPasswordResetInput) error {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" {
		return domainUser.ErrInvalidEmail
	}

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Por segurança, não revela se o email existe ou não
		return nil
	}

	// Gera token seguro
	token, err := generateSecureToken(32)
	if err != nil {
		return fmt.Errorf("falha ao gerar token: %w", err)
	}

	// Token válido por 1 hora
	expiresAt := time.Now().Add(1 * time.Hour)

	resetToken := domainUser.NewPasswordResetToken(user.ID, token, expiresAt)
	if err := uc.userRepo.CreatePasswordResetToken(ctx, resetToken); err != nil {
		return fmt.Errorf("falha ao salvar token de reset: %w", err)
	}

	// Envia email com link de reset
	resetLink := fmt.Sprintf("%s?token=%s", uc.resetBaseURL, token)
	if err := uc.emailService.SendPasswordResetEmail(user.Email.String(), resetLink); err != nil {
		return fmt.Errorf("falha ao enviar email: %w", err)
	}

	return nil
}

// generateSecureToken gera um token aleatório seguro.
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
