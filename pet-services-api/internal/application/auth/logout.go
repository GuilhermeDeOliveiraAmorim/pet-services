package auth

import (
	"context"
	"errors"

	"github.com/guilherme/pet-services-api/internal/application/exceptions"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainAuth "github.com/guilherme/pet-services-api/internal/domain/auth"
)

// LogoutUseCase revoga um refresh token específico.
type LogoutUseCase struct {
	refreshRepo  domainAuth.RefreshTokenRepository
	tokenService domainAuth.TokenService
	logger       logging.LoggerService
}

func NewLogoutUseCase(refreshRepo domainAuth.RefreshTokenRepository, tokenService domainAuth.TokenService, logger logging.LoggerService) *LogoutUseCase {
	return &LogoutUseCase{refreshRepo: refreshRepo, tokenService: tokenService, logger: logger}
}

// LogoutInput token de renovação do cliente.
type LogoutInput struct {
	RefreshToken string
}

const LOGOUT_USECASE = "LOGOUT_USECASE"

// Execute revoga o token atual, seguindo padrão de erros e logging.
func (uc *LogoutUseCase) Execute(ctx context.Context, input LogoutInput) []exceptions.ProblemDetails {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    LOGOUT_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	claims, err := uc.tokenService.ParseRefreshToken(input.RefreshToken)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC401_CODE,
			From:    LOGOUT_USECASE,
			Message: "Refresh token inválido",
			Error:   errors.New("refresh token inválido"),
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC401,
			Title:  "Refresh token inválido",
			Status: exceptions.RFC401_CODE,
			Detail: "Token de autenticação inválido ou expirado.",
		}}
	}

	if err := uc.refreshRepo.Revoke(ctx, claims.TokenID); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    LOGOUT_USECASE,
			Message: "Falha ao revogar refresh token",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao revogar refresh token",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    LOGOUT_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	return nil
}
