package auth

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

// LoginUseCase autentica e emite tokens.
type LoginUseCase struct {
	userRepo       domainUser.Repository
	refreshRepo    domainAuth.RefreshTokenRepository
	tokenService   domainAuth.TokenService
	passwordHasher domainAuth.PasswordHasher
	logger         *slog.Logger
}

func NewLoginUseCase(userRepo domainUser.Repository, refreshRepo domainAuth.RefreshTokenRepository, tokenService domainAuth.TokenService, passwordHasher domainAuth.PasswordHasher, logger *slog.Logger) *LoginUseCase {
	return &LoginUseCase{
		userRepo:       userRepo,
		refreshRepo:    refreshRepo,
		tokenService:   tokenService,
		passwordHasher: passwordHasher,
		logger:         logging.EnsureLogger(logger),
	}
}

// LoginInput dados de credenciais.
type LoginInput struct {
	Email    string
	Password string
}

// LoginOutput tokens emitidos.
type LoginOutput struct {
	UserID           uuid.UUID
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  int64
	RefreshExpiresAt int64
	UserType         domainUser.UserType
}

// Execute valida credenciais, gera tokens e registra refresh.
func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	var (
		result *LoginOutput
		err    error
	)
	email := strings.TrimSpace(strings.ToLower(input.Email))
	password := strings.TrimSpace(input.Password)
	defer logging.UseCase(ctx, uc.logger, "LoginUseCase", slog.String("email", email))(&err)

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		err = domainAuth.ErrInvalidCredentials
		return nil, err
	}

	if err := uc.passwordHasher.Compare(user.Password, password); err != nil {
		err = domainAuth.ErrInvalidCredentials
		return nil, err
	}

	tokens, err := uc.tokenService.GenerateTokens(user.ID, user.Type)
	if err != nil {
		return nil, err
	}

	rt := domainAuth.NewRefreshToken(user.ID, tokens.RefreshExpiresAt)
	if tokens.RefreshID != uuid.Nil {
		rt.ID = tokens.RefreshID
	}
	if err := uc.refreshRepo.Create(ctx, rt); err != nil {
		err = fmt.Errorf("falha ao salvar refresh token: %w", err)
		return nil, err
	}

	result = &LoginOutput{
		UserID:           user.ID,
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		AccessExpiresAt:  tokens.AccessExpiresAt.Unix(),
		RefreshExpiresAt: tokens.RefreshExpiresAt.Unix(),
		UserType:         user.Type,
	}

	return result, nil
}
