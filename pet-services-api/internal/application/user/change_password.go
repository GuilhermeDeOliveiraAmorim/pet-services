package user

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainAuth "github.com/guilherme/pet-services-api/internal/domain/auth"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// ChangePasswordUseCase permite ao usuário trocar sua senha.
type ChangePasswordUseCase struct {
	userRepo       domainUser.Repository
	refreshRepo    domainAuth.RefreshTokenRepository
	passwordHasher domainAuth.PasswordHasher
	logger         *slog.Logger
}

func NewChangePasswordUseCase(userRepo domainUser.Repository, refreshRepo domainAuth.RefreshTokenRepository, passwordHasher domainAuth.PasswordHasher, logger *slog.Logger) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		userRepo:       userRepo,
		refreshRepo:    refreshRepo,
		passwordHasher: passwordHasher,
		logger:         logging.EnsureLogger(logger),
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
	var err error
	defer logging.UseCase(ctx, uc.logger, "ChangePasswordUseCase", slog.String("user_id", input.UserID.String()))(&err)

	if input.UserID == uuid.Nil {
		err = domainUser.ErrUserNotFound
		return err
	}

	currentPassword := strings.TrimSpace(input.CurrentPassword)
	newPassword := strings.TrimSpace(input.NewPassword)

	if currentPassword == "" || newPassword == "" {
		err = domainUser.ErrInvalidPassword
		return err
	}

	if len(newPassword) < 6 {
		err = fmt.Errorf("nova senha deve ter no mínimo 6 caracteres")
		return err
	}

	if currentPassword == newPassword {
		err = fmt.Errorf("nova senha deve ser diferente da atual")
		return err
	}

	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	// Valida senha atual
	if err := uc.passwordHasher.Compare(user.Password, currentPassword); err != nil {
		err = domainAuth.ErrInvalidCredentials
		return err
	}

	// Hash nova senha
	hashed, err := uc.passwordHasher.Hash(newPassword)
	if err != nil {
		return err
	}

	user.SetPassword(hashed)

	if err := uc.userRepo.Update(ctx, user); err != nil {
		err = fmt.Errorf("falha ao atualizar senha: %w", err)
		return err
	}

	// Revoga todos os refresh tokens para forçar re-login
	_ = uc.refreshRepo.RevokeAllByUserID(ctx, user.ID)

	return nil
}
