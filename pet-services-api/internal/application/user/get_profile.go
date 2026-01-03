package user

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// GetProfileUseCase retorna o perfil do usuário autenticado.
type GetProfileUseCase struct {
	userRepo domainUser.Repository
	logger   *slog.Logger
}

func NewGetProfileUseCase(userRepo domainUser.Repository, logger *slog.Logger) *GetProfileUseCase {
	return &GetProfileUseCase{userRepo: userRepo, logger: logging.EnsureLogger(logger)}
}

// GetProfileInput entrada com ID do usuário autenticado.
type GetProfileInput struct {
	UserID uuid.UUID
}

// Execute busca e retorna o perfil completo.
func (uc *GetProfileUseCase) Execute(ctx context.Context, input GetProfileInput) (*domainUser.User, error) {
	var (
		result *domainUser.User
		err    error
	)
	defer logging.UseCase(ctx, uc.logger, "GetProfileUseCase", slog.String("user_id", input.UserID.String()))(&err)

	if input.UserID == uuid.Nil {
		err = domainUser.ErrUserNotFound
		return nil, err
	}

	result, err = uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
