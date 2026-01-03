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

// SignupUseCase registra um novo usuário e retorna tokens.
type SignupUseCase struct {
	userRepo       domainUser.Repository
	refreshRepo    domainAuth.RefreshTokenRepository
	tokenService   domainAuth.TokenService
	passwordHasher domainAuth.PasswordHasher
	logger         *slog.Logger
}

func NewSignupUseCase(userRepo domainUser.Repository, refreshRepo domainAuth.RefreshTokenRepository, tokenService domainAuth.TokenService, passwordHasher domainAuth.PasswordHasher, logger *slog.Logger) *SignupUseCase {
	return &SignupUseCase{
		userRepo:       userRepo,
		refreshRepo:    refreshRepo,
		tokenService:   tokenService,
		passwordHasher: passwordHasher,
		logger:         logging.EnsureLogger(logger),
	}
}

// SignupInput dados para criação de usuário.
type SignupInput struct {
	Email    string
	Name     string
	Phone    string
	Password string
	Type     domainUser.UserType
}

// SignupOutput resultado com usuário e tokens.
type SignupOutput struct {
	UserID           uuid.UUID
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  int64
	RefreshExpiresAt int64
	UserType         domainUser.UserType
}

// Execute cria usuário, hash de senha, tokens e registra refresh.
func (uc *SignupUseCase) Execute(ctx context.Context, input SignupInput) (*SignupOutput, error) {
	var (
		result *SignupOutput
		err    error
	)
	email := strings.TrimSpace(strings.ToLower(input.Email))
	name := strings.TrimSpace(input.Name)
	password := strings.TrimSpace(input.Password)
	defer logging.UseCase(ctx, uc.logger, "SignupUseCase", slog.String("email", email), slog.String("type", string(input.Type)))(&err)

	if password == "" || len(password) < 6 {
		err = domainUser.ErrInvalidPassword
		return nil, err
	}

	exists, err := uc.userRepo.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainUser.ErrUserAlreadyExists
	}

	user, err := domainUser.NewUser(email, name, input.Phone, input.Type)
	if err != nil {
		return nil, err
	}

	hashed, err := uc.passwordHasher.Hash(password)
	if err != nil {
		return nil, err
	}
	user.SetPassword(hashed)

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	tokens, err := uc.tokenService.GenerateTokens(user.ID, user.Type)
	if err != nil {
		return nil, err
	}

	rt := domainAuth.NewRefreshToken(user.ID, tokens.RefreshExpiresAt)
	// usa o ID gerado pelo token service para rastrear o mesmo token
	if tokens.RefreshID != uuid.Nil {
		rt.ID = tokens.RefreshID
	}
	if err := uc.refreshRepo.Create(ctx, rt); err != nil {
		err = fmt.Errorf("falha ao salvar refresh token: %w", err)
		return nil, err
	}

	result = &SignupOutput{
		UserID:           user.ID,
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		AccessExpiresAt:  tokens.AccessExpiresAt.Unix(),
		RefreshExpiresAt: tokens.RefreshExpiresAt.Unix(),
		UserType:         user.Type,
	}

	return result, nil
}
