package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	domainAuth "github.com/guilherme/pet-services-api/internal/domain/auth"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// DeleteAccountUseCase remove ou desativa conta do usuário.
type DeleteAccountUseCase struct {
	userRepo    domainUser.Repository
	refreshRepo domainAuth.RefreshTokenRepository
}

func NewDeleteAccountUseCase(userRepo domainUser.Repository, refreshRepo domainAuth.RefreshTokenRepository) *DeleteAccountUseCase {
	return &DeleteAccountUseCase{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
	}
}

// DeleteAccountInput dados para deletar conta.
type DeleteAccountInput struct {
	UserID     uuid.UUID
	HardDelete bool // se true, deleta permanentemente; se false, soft delete
}

// Execute marca usuário como deletado (soft) ou remove permanentemente (hard).
func (uc *DeleteAccountUseCase) Execute(ctx context.Context, input DeleteAccountInput) error {
	if input.UserID == uuid.Nil {
		return domainUser.ErrUserNotFound
	}

	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	if user.IsDeleted() {
		return fmt.Errorf("conta já foi deletada")
	}

	// Revoga todos os refresh tokens
	_ = uc.refreshRepo.RevokeAllByUserID(ctx, user.ID)

	if input.HardDelete {
		// Hard delete: remove permanentemente do banco
		if err := uc.userRepo.Delete(ctx, user.ID); err != nil {
			return fmt.Errorf("falha ao deletar conta: %w", err)
		}
	} else {
		// Soft delete: marca DeletedAt
		user.SoftDelete()
		if err := uc.userRepo.Update(ctx, user); err != nil {
			return fmt.Errorf("falha ao desativar conta: %w", err)
		}
	}

	return nil
}
