package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/guilherme/pet-services-api/internal/application/exceptions"
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
	logger         logging.LoggerService
}

func NewLoginUseCase(userRepo domainUser.Repository, refreshRepo domainAuth.RefreshTokenRepository, tokenService domainAuth.TokenService, passwordHasher domainAuth.PasswordHasher, logger logging.LoggerService) *LoginUseCase {
	return &LoginUseCase{
		userRepo:       userRepo,
		refreshRepo:    refreshRepo,
		tokenService:   tokenService,
		passwordHasher: passwordHasher,
		logger:         logger,
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

const LOGIN_USECASE = "LOGIN_USECASE"

// Execute valida credenciais, gera tokens e registra refresh.
func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInput) (*LoginOutput, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    LOGIN_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	email := strings.TrimSpace(strings.ToLower(input.Email))
	password := strings.TrimSpace(input.Password)

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC401_CODE,
			From:    LOGIN_USECASE,
			Message: "Credenciais inválidas",
			Error:   errors.New("usuário não encontrado"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC401,
			Title:  "Credenciais inválidas",
			Status: exceptions.RFC401_CODE,
			Detail: "Usuário não encontrado",
		}}
	}

	if err := uc.passwordHasher.Compare(user.Password, password); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC401_CODE,
			From:    LOGIN_USECASE,
			Message: "Senha inválida",
			Error:   errors.New("senha inválida"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC401,
			Title:  "Credenciais inválidas",
			Status: exceptions.RFC401_CODE,
			Detail: "Senha inválida",
		}}
	}

	tokens, err := uc.tokenService.GenerateTokens(user.ID, user.Type)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    LOGIN_USECASE,
			Message: "Erro ao gerar tokens",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Erro ao gerar tokens",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	rt := domainAuth.NewRefreshToken(user.ID, tokens.RefreshExpiresAt)
	if tokens.RefreshID != uuid.Nil {
		rt.ID = tokens.RefreshID
	}
	if err := uc.refreshRepo.Create(ctx, rt); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    LOGIN_USECASE,
			Message: "Falha ao salvar refresh token",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao salvar refresh token",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    LOGIN_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	result := &LoginOutput{
		UserID:           user.ID,
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		AccessExpiresAt:  tokens.AccessExpiresAt.Unix(),
		RefreshExpiresAt: tokens.RefreshExpiresAt.Unix(),
		UserType:         user.Type,
	}

	return result, nil
}
