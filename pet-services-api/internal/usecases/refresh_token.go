package usecases

import (
	"context"
	"errors"
	"strings"
	"time"

	"pet-services-api/internal/auth"
	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"

	"gorm.io/gorm"
)

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token"`
	UserAgent    string `json:"user_agent,omitempty"`
	IP           string `json:"ip,omitempty"`
}

type RefreshTokenOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type RefreshTokenUseCase struct {
	userRepository         entities.UserRepository
	refreshTokenRepository entities.RefreshTokenRepository
	logger                 logging.LoggerInterface
}

func NewRefreshTokenUseCase(userRepo entities.UserRepository, refreshRepo entities.RefreshTokenRepository, logger logging.LoggerInterface) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		userRepository:         userRepo,
		refreshTokenRepository: refreshRepo,
		logger:                 logger,
	}
}

func (uc *RefreshTokenUseCase) Execute(ctx context.Context, input RefreshTokenInput) (*RefreshTokenOutput, []exceptions.ProblemDetails) {
	const from = "RefreshTokenUseCase.Execute"

	refreshToken := strings.TrimSpace(input.RefreshToken)
	if refreshToken == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "Refresh token ausente", errors.New("O refresh token é obrigatório"))
	}
	input.RefreshToken = refreshToken

	jwtSvc, err := auth.NewJWTServiceFromEnv()
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao carregar configuração JWT", err)
	}

	claims, err := jwtSvc.ValidateRefreshToken(input.RefreshToken)
	if err != nil {
		return nil, uc.logger.LogUnauthorized(ctx, from, "Refresh token inválido", errors.New("Refresh token inválido ou expirado"))
	}

	storedToken, err := uc.refreshTokenRepository.FindByToken(input.RefreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, uc.logger.LogUnauthorized(ctx, from, "Refresh token inválido", errors.New("Refresh token não encontrado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar refresh token", err)
	}

	if !storedToken.IsValid() {
		return nil, uc.logger.LogUnauthorized(ctx, from, "Refresh token inválido", errors.New("Refresh token revogado ou expirado"))
	}

	user, err := uc.userRepository.FindByID(claims.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar um usuário com o ID informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.Active {
		return nil, uc.logger.LogForbidden(ctx, from, "Conta desativada", errors.New("Sua conta foi desativada. Entre em contato com o suporte para reativar"))
	}

	if !user.EmailVerified {
		return nil, uc.logger.LogForbidden(ctx, from, "Email não verificado", errors.New("Verifique seu email antes de continuar"))
	}

	accessToken, err := jwtSvc.GenerateAccessToken(user.ID, user.Login.Email, user.UserType)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao gerar access token", err)
	}

	newRefreshToken, err := jwtSvc.GenerateRefreshToken(user.ID, user.Login.Email, user.UserType)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao gerar refresh token", err)
	}

	accessClaims, err := jwtSvc.ValidateAccessToken(accessToken)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao validar access token", err)
	}

	expiresIn := max(int64(time.Until(accessClaims.ExpiresAt.Time).Seconds()), 0)

	refreshClaims, err := jwtSvc.ValidateRefreshToken(newRefreshToken)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao validar refresh token", err)
	}

	refreshEntity, problems := entities.NewRefreshToken(user.ID, newRefreshToken, refreshClaims.ExpiresAt.Time, input.UserAgent, input.IP)
	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Refresh token inválido", problems)
		return nil, problems
	}

	if err := uc.refreshTokenRepository.Revoke(storedToken.ID); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao revogar refresh token antigo", err)
	}

	if err := uc.refreshTokenRepository.Create(refreshEntity); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao salvar refresh token", err)
	}

	return &RefreshTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}
