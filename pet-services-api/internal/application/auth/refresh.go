package auth

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainAuth "github.com/guilherme/pet-services-api/internal/domain/auth"
)

// RefreshTokenUseCase rotaciona o refresh token e emite novo par de tokens.
type RefreshTokenUseCase struct {
	refreshRepo  domainAuth.RefreshTokenRepository
	tokenService domainAuth.TokenService
	logger       *slog.Logger
}

func NewRefreshTokenUseCase(refreshRepo domainAuth.RefreshTokenRepository, tokenService domainAuth.TokenService, logger *slog.Logger) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		refreshRepo:  refreshRepo,
		tokenService: tokenService,
		logger:       logging.EnsureLogger(logger),
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
	var (
		result *RefreshOutput
		err    error
	)
	defer logging.UseCase(ctx, uc.logger, "RefreshTokenUseCase")(&err)

	claims, err := uc.tokenService.ParseRefreshToken(input.RefreshToken)
	if err != nil {
		err = domainAuth.ErrInvalidCredentials
		return nil, err
	}

	stored, err := uc.refreshRepo.FindByID(ctx, claims.TokenID)
	if err != nil {
		err = domainAuth.ErrRefreshTokenNotFound
		return nil, err
	}

	now := time.Now()
	if stored.Revoked {
		err = domainAuth.ErrRefreshTokenRevoked
		return nil, err
	}
	if now.After(stored.ExpiresAt) || now.After(claims.ExpiresAt) {
		err = domainAuth.ErrRefreshTokenExpired
		return nil, err
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
		err = fmt.Errorf("falha ao salvar refresh token: %w", err)
		return nil, err
	}

	result = &RefreshOutput{
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		AccessExpiresAt:  tokens.AccessExpiresAt.Unix(),
		RefreshExpiresAt: tokens.RefreshExpiresAt.Unix(),
		UserID:           claims.UserID.String(),
		UserType:         string(claims.UserType),
	}

	return result, nil
}
