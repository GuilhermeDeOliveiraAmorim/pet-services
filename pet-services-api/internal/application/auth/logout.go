package auth

import (
	"context"
	"log/slog"

	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainAuth "github.com/guilherme/pet-services-api/internal/domain/auth"
)

// LogoutUseCase revoga um refresh token específico.
type LogoutUseCase struct {
	refreshRepo  domainAuth.RefreshTokenRepository
	tokenService domainAuth.TokenService
	logger       *slog.Logger
}

func NewLogoutUseCase(refreshRepo domainAuth.RefreshTokenRepository, tokenService domainAuth.TokenService, logger *slog.Logger) *LogoutUseCase {
	return &LogoutUseCase{refreshRepo: refreshRepo, tokenService: tokenService, logger: logging.EnsureLogger(logger)}
}

// LogoutInput token de renovação do cliente.
type LogoutInput struct {
	RefreshToken string
}

// Execute revoga o token atual.
func (uc *LogoutUseCase) Execute(ctx context.Context, input LogoutInput) error {
	var (
		err    error
		userID string
	)
	defer logging.UseCase(ctx, uc.logger, "LogoutUseCase", slog.String("user_id", userID))(&err)

	claims, err := uc.tokenService.ParseRefreshToken(input.RefreshToken)
	if err != nil {
		err = domainAuth.ErrInvalidCredentials
		return err
	}
	userID = claims.UserID.String()

	if err := uc.refreshRepo.Revoke(ctx, claims.TokenID); err != nil {
		return err
	}
	return nil
}
