package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type ReactivateUserInput struct {
	UserID string `json:"user_id"`
}

type ReactivateUserOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type ReactivateUserUseCase struct {
	userRepository entities.UserRepository
	logger         logging.LoggerInterface
}

func NewReactivateUserUseCase(userRepo entities.UserRepository, logger logging.LoggerInterface) *ReactivateUserUseCase {
	return &ReactivateUserUseCase{
		userRepository: userRepo,
		logger:         logger,
	}
}

func (uc *ReactivateUserUseCase) Execute(ctx context.Context, input ReactivateUserInput) (*ReactivateUserOutput, []exceptions.ProblemDetails) {
	const from = "ReactivateUserUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar um usuário com o ID informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if user.Active {
		return &ReactivateUserOutput{
			Message: "Conta já ativa",
			Detail:  "Esta conta já está ativa",
		}, nil
	}

	user.Activate()

	if err := uc.userRepository.Update(user); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao reativar usuário", err)
	}

	return &ReactivateUserOutput{
		Message: "Conta reativada com sucesso",
		Detail:  "Sua conta foi reativada e você pode fazer login novamente",
	}, nil
}
