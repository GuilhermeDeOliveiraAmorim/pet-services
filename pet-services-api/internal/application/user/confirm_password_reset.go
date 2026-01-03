package user

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainAuth "github.com/guilherme/pet-services-api/internal/domain/auth"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// ConfirmPasswordResetUseCase confirma e aplica o reset de senha.
type ConfirmPasswordResetUseCase struct {
	userRepo       domainUser.Repository
	refreshRepo    domainAuth.RefreshTokenRepository
	passwordHasher domainAuth.PasswordHasher
	logger         *slog.Logger
}

func NewConfirmPasswordResetUseCase(userRepo domainUser.Repository, refreshRepo domainAuth.RefreshTokenRepository, passwordHasher domainAuth.PasswordHasher, logger *slog.Logger) *ConfirmPasswordResetUseCase {
	return &ConfirmPasswordResetUseCase{
		userRepo:       userRepo,
		refreshRepo:    refreshRepo,
		passwordHasher: passwordHasher,
		logger:         logging.EnsureLogger(logger),
	}
}

// ConfirmPasswordResetInput dados para confirmar reset.
type ConfirmPasswordResetInput struct {
	Token       string
	NewPassword string
}

// Execute valida token, atualiza senha e revoga tokens existentes.
func (uc *ConfirmPasswordResetUseCase) Execute(ctx context.Context, input ConfirmPasswordResetInput) error {
	var err error
	defer logging.UseCase(ctx, uc.logger, "ConfirmPasswordResetUseCase")(&err)

	token := strings.TrimSpace(input.Token)
	newPassword := strings.TrimSpace(input.NewPassword)

	if token == "" {
		err = domainUser.ErrPasswordResetTokenInvalid
		return err
	}

	if newPassword == "" || len(newPassword) < 6 {
		err = domainUser.ErrInvalidPassword
		return err
	}

	// Busca token
	resetToken, err := uc.userRepo.FindPasswordResetToken(ctx, token)
	if err != nil {
		err = domainUser.ErrPasswordResetTokenInvalid
		return err
	}

	if !resetToken.IsValid() {
		err = domainUser.ErrPasswordResetTokenInvalid
		return err
	}

	// Busca usuário
	user, err := uc.userRepo.FindByID(ctx, resetToken.UserID)
	if err != nil {
		return err
	}

	// Hash nova senha
	hashed, err := uc.passwordHasher.Hash(newPassword)
	if err != nil {
		return err
	}

	user.SetPassword(hashed)

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("falha ao atualizar senha: %w", err)
	}

	// Marca token como usado
	if err := uc.userRepo.MarkPasswordResetTokenAsUsed(ctx, resetToken.ID); err != nil {
		err = fmt.Errorf("falha ao marcar token como usado: %w", err)
		return err
	}

	// Revoga todos os refresh tokens
	_ = uc.refreshRepo.RevokeAllByUserID(ctx, user.ID)

	return nil
}
