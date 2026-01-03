package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// RequestPasswordResetUseCase inicia o fluxo de reset de senha.
type RequestPasswordResetUseCase struct {
	userRepo     domainUser.Repository
	emailService domainUser.EmailService
	resetBaseURL string
	logger       *slog.Logger
}

func NewRequestPasswordResetUseCase(userRepo domainUser.Repository, emailService domainUser.EmailService, resetBaseURL string, logger *slog.Logger) *RequestPasswordResetUseCase {
	return &RequestPasswordResetUseCase{
		userRepo:     userRepo,
		emailService: emailService,
		resetBaseURL: resetBaseURL,
		logger:       logging.EnsureLogger(logger),
	}
}

// RequestPasswordResetInput dados para solicitar reset.
type RequestPasswordResetInput struct {
	Email string
}

// Execute gera token, persiste e envia email.
func (uc *RequestPasswordResetUseCase) Execute(ctx context.Context, input RequestPasswordResetInput) error {
	var err error
	email := strings.TrimSpace(strings.ToLower(input.Email))
	defer logging.UseCase(ctx, uc.logger, "RequestPasswordResetUseCase", slog.String("email", email))(&err)
	if email == "" {
		err = domainUser.ErrInvalidEmail
		return err
	}

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Por segurança, não revela se o email existe ou não
		return nil
	}

	// Gera token seguro
	token, err := generateSecureToken(32)
	if err != nil {
		err = fmt.Errorf("falha ao gerar token: %w", err)
		return err
	}

	// Token válido por 1 hora
	expiresAt := time.Now().Add(1 * time.Hour)

	resetToken := domainUser.NewPasswordResetToken(user.ID, token, expiresAt)
	if err := uc.userRepo.CreatePasswordResetToken(ctx, resetToken); err != nil {
		err = fmt.Errorf("falha ao salvar token de reset: %w", err)
		return err
	}

	// Envia email com link de reset
	resetLink := fmt.Sprintf("%s?token=%s", uc.resetBaseURL, token)
	if err := uc.emailService.SendPasswordResetEmail(user.Email.String(), resetLink); err != nil {
		err = fmt.Errorf("falha ao enviar email: %w", err)
		return err
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
