package user

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainAuth "github.com/guilherme/pet-services-api/internal/domain/auth"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// DeleteAccountUseCase remove ou desativa conta do usuário.
type DeleteAccountUseCase struct {
	userRepo    domainUser.Repository
	refreshRepo domainAuth.RefreshTokenRepository
	logger      *slog.Logger
}

func NewDeleteAccountUseCase(userRepo domainUser.Repository, refreshRepo domainAuth.RefreshTokenRepository, logger *slog.Logger) *DeleteAccountUseCase {
	return &DeleteAccountUseCase{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		logger:      logging.EnsureLogger(logger),
	}
}

// DeleteAccountInput dados para deletar conta.
type DeleteAccountInput struct {
	UserID     uuid.UUID
	HardDelete bool // se true, deleta permanentemente; se false, soft delete
}

// Execute marca usuário como deletado (soft) ou remove permanentemente (hard).
func (uc *DeleteAccountUseCase) Execute(ctx context.Context, input DeleteAccountInput) error {
	var err error
	defer logging.UseCase(ctx, uc.logger, "DeleteAccountUseCase", slog.String("user_id", input.UserID.String()))(&err)

	if input.UserID == uuid.Nil {
		err = domainUser.ErrUserNotFound
		return err
	}

	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	if user.IsDeleted() {
		err = fmt.Errorf("conta já foi deletada")
		return err
	}

	// Revoga todos os refresh tokens
	_ = uc.refreshRepo.RevokeAllByUserID(ctx, user.ID)

	if input.HardDelete {
		// Hard delete: remove permanentemente do banco
		if err := uc.userRepo.Delete(ctx, user.ID); err != nil {
			err = fmt.Errorf("falha ao deletar conta: %w", err)
			return err
		}
	} else {
		// Soft delete: marca DeletedAt
		user.SoftDelete()
		if err := uc.userRepo.Update(ctx, user); err != nil {
			err = fmt.Errorf("falha ao desativar conta: %w", err)
			return err
		}
	}

	return nil
}
