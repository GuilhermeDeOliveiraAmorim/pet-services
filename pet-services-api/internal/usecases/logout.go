package usecases

import (
	"context"

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
}

func NewLogoutUseCase(refreshTokenRepo entities.RefreshTokenRepository) *LogoutUseCase {
	return &LogoutUseCase{
		refreshTokenRepository: refreshTokenRepo,
	}
}

func (uc *LogoutUseCase) Execute(ctx context.Context, input LogoutInput) (*LogoutOutput, []exceptions.ProblemDetails) {
	const from = "LogoutUseCase.Execute"

	if input.UserID == "" {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "ID do usuário ausente",
				Detail: "O ID do usuário é obrigatório",
			}),
		}
	}

	if input.RevokeAll {
		if err := uc.refreshTokenRepository.RevokeAllByUserID(input.UserID); err != nil {
			return nil, logging.InternalServerError(ctx, from, "Erro ao revogar todos os tokens", err)
		}

		return &LogoutOutput{
			Message: "Logout realizado com sucesso",
			Detail:  "Todos os tokens foram revogados em todos os dispositivos",
		}, nil
	}

	if input.TokenID != "" {
		if err := uc.refreshTokenRepository.Revoke(input.TokenID); err != nil {
			return nil, logging.InternalServerError(ctx, from, "Erro ao revogar token", err)
		}

		return &LogoutOutput{
			Message: "Logout realizado com sucesso",
			Detail:  "Token revogado com sucesso",
		}, nil
	}

	return nil, []exceptions.ProblemDetails{
		exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Token ou opção de revogação ausente",
			Detail: "Forneça um TokenID ou marque RevokeAll como true",
		}),
	}
}
