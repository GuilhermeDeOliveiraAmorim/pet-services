package auth

import (
	"context"
	"fmt"
	"time"

	domainAuth "github.com/guilherme/pet-services-api/internal/domain/auth"
)

// RefreshTokenUseCase rotaciona o refresh token e emite novo par de tokens.
type RefreshTokenUseCase struct {
	refreshRepo  domainAuth.RefreshTokenRepository
	tokenService domainAuth.TokenService
}

func NewRefreshTokenUseCase(refreshRepo domainAuth.RefreshTokenRepository, tokenService domainAuth.TokenService) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		refreshRepo:  refreshRepo,
		tokenService: tokenService,
	}
}

// RefreshInput token de renovação enviado pelo cliente.
type RefreshInput struct {
	RefreshToken string
}

// RefreshOutput novo par de tokens.
type RefreshOutput struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  int64
	RefreshExpiresAt int64
	UserID           string
	UserType         string
}

// Execute valida o refresh token, revoga o antigo e cria um novo.
func (uc *RefreshTokenUseCase) Execute(ctx context.Context, input RefreshInput) (*RefreshOutput, error) {
	claims, err := uc.tokenService.ParseRefreshToken(input.RefreshToken)
	if err != nil {
		return nil, domainAuth.ErrInvalidCredentials
	}

	stored, err := uc.refreshRepo.FindByID(ctx, claims.TokenID)
	if err != nil {
		return nil, domainAuth.ErrRefreshTokenNotFound
	}

	now := time.Now()
	if stored.Revoked {
		return nil, domainAuth.ErrRefreshTokenRevoked
	}
	if now.After(stored.ExpiresAt) || now.After(claims.ExpiresAt) {
		return nil, domainAuth.ErrRefreshTokenExpired
	}

	// Revoga o token anterior
	_ = uc.refreshRepo.Revoke(ctx, stored.ID)

	tokens, err := uc.tokenService.GenerateTokens(claims.UserID, claims.UserType)
	if err != nil {
		return nil, err
	}

	rt := domainAuth.NewRefreshToken(claims.UserID, tokens.RefreshExpiresAt)
	if tokens.RefreshID != claims.TokenID { // usualmente gera um novo
		rt.ID = tokens.RefreshID
	} else {
		rt.ID = claims.TokenID
	}

	if err := uc.refreshRepo.Create(ctx, rt); err != nil {
		return nil, fmt.Errorf("falha ao salvar refresh token: %w", err)
	}

	return &RefreshOutput{
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		AccessExpiresAt:  tokens.AccessExpiresAt.Unix(),
		RefreshExpiresAt: tokens.RefreshExpiresAt.Unix(),
		UserID:           claims.UserID.String(),
		UserType:         string(claims.UserType),
	}, nil
}
