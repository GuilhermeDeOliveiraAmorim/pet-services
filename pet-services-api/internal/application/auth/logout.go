package auth

import (
	"context"

	domainAuth "github.com/guilherme/pet-services-api/internal/domain/auth"
)

// LogoutUseCase revoga um refresh token específico.
type LogoutUseCase struct {
	refreshRepo  domainAuth.RefreshTokenRepository
	tokenService domainAuth.TokenService
}

func NewLogoutUseCase(refreshRepo domainAuth.RefreshTokenRepository, tokenService domainAuth.TokenService) *LogoutUseCase {
	return &LogoutUseCase{refreshRepo: refreshRepo, tokenService: tokenService}
}

// LogoutInput token de renovação do cliente.
type LogoutInput struct {
	RefreshToken string
}

// Execute revoga o token atual.
func (uc *LogoutUseCase) Execute(ctx context.Context, input LogoutInput) error {
	claims, err := uc.tokenService.ParseRefreshToken(input.RefreshToken)
	if err != nil {
		return domainAuth.ErrInvalidCredentials
	}

	if err := uc.refreshRepo.Revoke(ctx, claims.TokenID); err != nil {
		return err
	}
	return nil
}
