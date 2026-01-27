package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"

	"gorm.io/gorm"
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
	logger                 logging.LoggerInterface
}

func NewDeactivateUserUseCase(userRepo entities.UserRepository, refreshTokenRepo entities.RefreshTokenRepository, logger logging.LoggerInterface) *DeactivateUserUseCase {
	return &DeactivateUserUseCase{
		userRepository:         userRepo,
		refreshTokenRepository: refreshTokenRepo,
		logger:                 logger,
	}
}

func (uc *DeactivateUserUseCase) Execute(ctx context.Context, input DeactivateUserInput) (*DeactivateUserOutput, []exceptions.ProblemDetails) {
	const from = "DeactivateUserUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar um usuário com o ID informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.Active {
		return &DeactivateUserOutput{
			Message: "Conta já desativada",
			Detail:  "Esta conta já foi desativada anteriormente",
		}, nil
	}

	user.Deactivate()

	if err := uc.userRepository.Update(user); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao desativar usuário", err)
	}

	if err := uc.refreshTokenRepository.RevokeAllByUserID(input.UserID); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao revogar tokens", err)
	}

	return &DeactivateUserOutput{
		Message: "Conta desativada com sucesso",
		Detail:  "Sua conta foi desativada e todos os tokens foram revogados",
	}, nil
}
