package auth

import (
    "context"
    "fmt"
    "strings"

    "github.com/google/uuid"

    domainAuth "github.com/guilherme/pet-services-api/internal/domain/auth"
    domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// LoginUseCase autentica e emite tokens.
type LoginUseCase struct {
    userRepo       domainUser.Repository
    refreshRepo    domainAuth.RefreshTokenRepository
    tokenService   domainAuth.TokenService
    passwordHasher domainAuth.PasswordHasher
}

func NewLoginUseCase(userRepo domainUser.Repository, refreshRepo domainAuth.RefreshTokenRepository, tokenService domainAuth.TokenService, passwordHasher domainAuth.PasswordHasher) *LoginUseCase {
    return &LoginUseCase{
        userRepo:       userRepo,
        refreshRepo:    refreshRepo,
        tokenService:   tokenService,
        passwordHasher: passwordHasher,
    }
}

// LoginInput dados de credenciais.
type LoginInput struct {
    Email    string
    Password string
}

// LoginOutput tokens emitidos.
type LoginOutput struct {
    UserID          uuid.UUID
    AccessToken     string
    RefreshToken    string
    AccessExpiresAt  int64
    RefreshExpiresAt int64
    UserType        domainUser.UserType
}

// Execute valida credenciais, gera tokens e registra refresh.
func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, error) {
    email := strings.TrimSpace(strings.ToLower(input.Email))
    password := strings.TrimSpace(input.Password)

    user, err := uc.userRepo.FindByEmail(ctx, email)
    if err != nil {
        return nil, domainAuth.ErrInvalidCredentials
    }

    if err := uc.passwordHasher.Compare(user.Password, password); err != nil {
        return nil, domainAuth.ErrInvalidCredentials
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
        return nil, fmt.Errorf("falha ao salvar refresh token: %w", err)
    }

    return &LoginOutput{
        UserID:          user.ID,
        AccessToken:     tokens.AccessToken,
        RefreshToken:    tokens.RefreshToken,
        AccessExpiresAt:  tokens.AccessExpiresAt.Unix(),
        RefreshExpiresAt: tokens.RefreshExpiresAt.Unix(),
        UserType:        user.Type,
    }, nil
}
