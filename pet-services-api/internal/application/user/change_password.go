package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	domainAuth "github.com/guilherme/pet-services-api/internal/domain/auth"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// ChangePasswordUseCase permite ao usuário trocar sua senha.
type ChangePasswordUseCase struct {
	userRepo       domainUser.Repository
	refreshRepo    domainAuth.RefreshTokenRepository
	passwordHasher domainAuth.PasswordHasher
}

func NewChangePasswordUseCase(userRepo domainUser.Repository, refreshRepo domainAuth.RefreshTokenRepository, passwordHasher domainAuth.PasswordHasher) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		userRepo:       userRepo,
		refreshRepo:    refreshRepo,
		passwordHasher: passwordHasher,
	}
}

// ChangePasswordInput dados para troca de senha.
type ChangePasswordInput struct {
	UserID          uuid.UUID
	CurrentPassword string
	NewPassword     string
}

// Execute valida senha atual, atualiza com nova e revoga tokens existentes.
func (uc *ChangePasswordUseCase) Execute(ctx context.Context, input ChangePasswordInput) error {
	if input.UserID == uuid.Nil {
		return domainUser.ErrUserNotFound
	}

	currentPassword := strings.TrimSpace(input.CurrentPassword)
	newPassword := strings.TrimSpace(input.NewPassword)

	if currentPassword == "" || newPassword == "" {
		return domainUser.ErrInvalidPassword
	}

	if len(newPassword) < 6 {
		return fmt.Errorf("nova senha deve ter no mínimo 6 caracteres")
	}

	if currentPassword == newPassword {
		return fmt.Errorf("nova senha deve ser diferente da atual")
	}

	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	// Valida senha atual
	if err := uc.passwordHasher.Compare(user.Password, currentPassword); err != nil {
		return domainAuth.ErrInvalidCredentials
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

	// Revoga todos os refresh tokens para forçar re-login
	_ = uc.refreshRepo.RevokeAllByUserID(ctx, user.ID)

	return nil
}
