package user

import (
	"context"
	"fmt"
	"strings"

	domainAuth "github.com/guilherme/pet-services-api/internal/domain/auth"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// ConfirmPasswordResetUseCase confirma e aplica o reset de senha.
type ConfirmPasswordResetUseCase struct {
	userRepo       domainUser.Repository
	refreshRepo    domainAuth.RefreshTokenRepository
	passwordHasher domainAuth.PasswordHasher
}

func NewConfirmPasswordResetUseCase(userRepo domainUser.Repository, refreshRepo domainAuth.RefreshTokenRepository, passwordHasher domainAuth.PasswordHasher) *ConfirmPasswordResetUseCase {
	return &ConfirmPasswordResetUseCase{
		userRepo:       userRepo,
		refreshRepo:    refreshRepo,
		passwordHasher: passwordHasher,
	}
}

// ConfirmPasswordResetInput dados para confirmar reset.
type ConfirmPasswordResetInput struct {
	Token       string
	NewPassword string
}

// Execute valida token, atualiza senha e revoga tokens existentes.
func (uc *ConfirmPasswordResetUseCase) Execute(ctx context.Context, input ConfirmPasswordResetInput) error {
	token := strings.TrimSpace(input.Token)
	newPassword := strings.TrimSpace(input.NewPassword)

	if token == "" {
		return domainUser.ErrPasswordResetTokenInvalid
	}

	if newPassword == "" || len(newPassword) < 6 {
		return domainUser.ErrInvalidPassword
	}

	// Busca token
	resetToken, err := uc.userRepo.FindPasswordResetToken(ctx, token)
	if err != nil {
		return domainUser.ErrPasswordResetTokenInvalid
	}

	if !resetToken.IsValid() {
		return domainUser.ErrPasswordResetTokenInvalid
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
		return fmt.Errorf("falha ao marcar token como usado: %w", err)
	}

	// Revoga todos os refresh tokens
	_ = uc.refreshRepo.RevokeAllByUserID(ctx, user.ID)

	return nil
}