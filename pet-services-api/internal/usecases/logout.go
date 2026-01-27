package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type LogoutInput struct {
	UserID    string `json:"user_id"`
	TokenID   string `json:"token_id,omitempty"`
	RevokeAll bool   `json:"revoke_all"`
}

type LogoutOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type LogoutUseCase struct {
	refreshTokenRepository entities.RefreshTokenRepository
	logger                 logging.LoggerInterface
}

func NewLogoutUseCase(refreshTokenRepo entities.RefreshTokenRepository, logger logging.LoggerInterface) *LogoutUseCase {
	return &LogoutUseCase{
		refreshTokenRepository: refreshTokenRepo,
		logger:                 logger,
	}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, input LogoutInput) (*LogoutOutput, []exceptions.ProblemDetails) {
	const from = "LogoutUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.RevokeAll {
		if err := uc.refreshTokenRepository.RevokeAllByUserID(input.UserID); err != nil {
			return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao revogar todos os tokens", err)
		}

		return &LogoutOutput{
			Message: "Logout realizado com sucesso",
			Detail:  "Todos os tokens foram revogados em todos os dispositivos",
		}, nil
	}

	if input.TokenID != "" {
		if err := uc.refreshTokenRepository.Revoke(input.TokenID); err != nil {
			return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao revogar token", err)
		}

		return &LogoutOutput{
			Message: "Logout realizado com sucesso",
			Detail:  "Token revogado com sucesso",
		}, nil
	}

	return nil, uc.logger.LogBadRequest(ctx, from, "Token ou opção de revogação ausente", errors.New("Forneça um TokenID ou marque RevokeAll como true"))
}
