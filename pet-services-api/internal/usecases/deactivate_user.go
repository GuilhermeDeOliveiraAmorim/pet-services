package usecases

import (
	"context"

	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type DeactivateUserInput struct {
	UserID string `json:"user_id"`
}

type DeactivateUserOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type DeactivateUserUseCase struct {
	userRepository         entities.UserRepository
	refreshTokenRepository entities.RefreshTokenRepository
}

func NewDeactivateUserUseCase(userRepo entities.UserRepository, refreshTokenRepo entities.RefreshTokenRepository) *DeactivateUserUseCase {
	return &DeactivateUserUseCase{
		userRepository:         userRepo,
		refreshTokenRepository: refreshTokenRepo,
	}
}

func (uc *DeactivateUserUseCase) Execute(ctx context.Context, input DeactivateUserInput) (*DeactivateUserOutput, []exceptions.ProblemDetails) {
	const from = "DeactivateUserUseCase.Execute"

	if input.UserID == "" {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "ID do usuário ausente",
				Detail: "O ID do usuário é obrigatório",
			}),
		}
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.Active {
		return &DeactivateUserOutput{
			Message: "Conta já desativada",
			Detail:  "Esta conta já foi desativada anteriormente",
		}, nil
	}

	user.Deactivate()

	if err := uc.userRepository.Update(user); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao desativar usuário", err)
	}

	if err := uc.refreshTokenRepository.RevokeAllByUserID(input.UserID); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao revogar tokens", err)
	}

	return &DeactivateUserOutput{
		Message: "Conta desativada com sucesso",
		Detail:  "Sua conta foi desativada e todos os tokens foram revogados",
	}, nil
}
