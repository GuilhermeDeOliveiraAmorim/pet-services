package auth

import (
	"context"
	"errors"
	"time"

	"github.com/guilherme/pet-services-api/internal/application/exceptions"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainAuth "github.com/guilherme/pet-services-api/internal/domain/auth"
)

// RefreshTokenUseCase rotaciona o refresh token e emite novo par de tokens.
type RefreshTokenUseCase struct {
	refreshRepo  domainAuth.RefreshTokenRepository
	tokenService domainAuth.TokenService
	logger       logging.LoggerService
}

func NewRefreshTokenUseCase(refreshRepo domainAuth.RefreshTokenRepository, tokenService domainAuth.TokenService, logger logging.LoggerService) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		refreshRepo:  refreshRepo,
		tokenService: tokenService,
		logger:       logger,
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

const REFRESH_USECASE = "REFRESH_USECASE"

// Execute valida o refresh token, revoga o antigo e cria um novo, seguindo padrão de erros e logging.
func (uc *RefreshTokenUseCase) Execute(ctx context.Context, input RefreshInput) (*RefreshOutput, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    REFRESH_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	claims, err := uc.tokenService.ParseRefreshToken(input.RefreshToken)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC401_CODE,
			From:    REFRESH_USECASE,
			Message: "Refresh token inválido",
			Error:   errors.New("refresh token inválido"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC401,
			Title:  "Refresh token inválido",
			Status: exceptions.RFC401_CODE,
			Detail: "Token de autenticação inválido ou expirado.",
		}}
	}

	stored, err := uc.refreshRepo.FindByID(ctx, claims.TokenID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC401_CODE,
			From:    REFRESH_USECASE,
			Message: "Refresh token não encontrado",
			Error:   errors.New("refresh token não encontrado"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC401,
			Title:  "Refresh token não encontrado",
			Status: exceptions.RFC401_CODE,
			Detail: "Token de autenticação não encontrado.",
		}}
	}

	now := time.Now()
	if stored.Revoked {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC401_CODE,
			From:    REFRESH_USECASE,
			Message: "Refresh token revogado",
			Error:   errors.New("refresh token revogado"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC401,
			Title:  "Refresh token revogado",
			Status: exceptions.RFC401_CODE,
			Detail: "Token de autenticação revogado.",
		}}
	}
	if now.After(stored.ExpiresAt) || now.After(claims.ExpiresAt) {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC401_CODE,
			From:    REFRESH_USECASE,
			Message: "Refresh token expirado",
			Error:   errors.New("refresh token expirado"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC401,
			Title:  "Refresh token expirado",
			Status: exceptions.RFC401_CODE,
			Detail: "Token de autenticação expirado.",
		}}
	}

	// Revoga o token anterior
	_ = uc.refreshRepo.Revoke(ctx, stored.ID)

	tokens, err := uc.tokenService.GenerateTokens(claims.UserID, claims.UserType)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    REFRESH_USECASE,
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

	rt := domainAuth.NewRefreshToken(claims.UserID, tokens.RefreshExpiresAt)
	if tokens.RefreshID != claims.TokenID { // usualmente gera um novo
		rt.ID = tokens.RefreshID
	} else {
		rt.ID = claims.TokenID
	}

	if err := uc.refreshRepo.Create(ctx, rt); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    REFRESH_USECASE,
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
		From:    REFRESH_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	result := &RefreshOutput{
		AccessToken:      tokens.AccessToken,
		RefreshToken:     tokens.RefreshToken,
		AccessExpiresAt:  tokens.AccessExpiresAt.Unix(),
		RefreshExpiresAt: tokens.RefreshExpiresAt.Unix(),
		UserID:           claims.UserID.String(),
		UserType:         string(claims.UserType),
	}

	return result, nil
}
