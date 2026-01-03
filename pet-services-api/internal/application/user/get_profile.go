package user

import (
	"context"

	"github.com/google/uuid"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// GetProfileUseCase retorna o perfil do usuário autenticado.
type GetProfileUseCase struct {
	userRepo domainUser.Repository
}

func NewGetProfileUseCase(userRepo domainUser.Repository) *GetProfileUseCase {
	return &GetProfileUseCase{userRepo: userRepo}
}

// GetProfileInput entrada com ID do usuário autenticado.
type GetProfileInput struct {
	UserID uuid.UUID
}

// Execute busca e retorna o perfil completo.
func (uc *GetProfileUseCase) Execute(ctx context.Context, input GetProfileInput) (*domainUser.User, error) {
	if input.UserID == uuid.Nil {
		return nil, domainUser.ErrUserNotFound
	}

	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
