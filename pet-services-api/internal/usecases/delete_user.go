package usecases

import (
	"context"
	"errors"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"

	"gorm.io/gorm"
)

type DeleteUserInput struct {
	UserID string `json:"user_id"`
}

type DeleteUserOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type DeleteUserUseCase struct {
	userRepository entities.UserRepository
	logger         logging.LoggerInterface
}

func NewDeleteUserUseCase(userRepository entities.UserRepository, logger logging.LoggerInterface) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, input DeleteUserInput) (*DeleteUserOutput, []exceptions.ProblemDetails) {
	const from = "DeleteUserUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório para deletar"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar um usuário com o ID informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.Active {
		return nil, uc.logger.LogBadRequest(ctx, from, "Usuário já inativo", errors.New("O usuário já está inativo no sistema"))
	}

	user.Deactivate()

	if err := uc.userRepository.Update(user); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao deletar usuário", err)
	}

	return &DeleteUserOutput{
		Message: "Usuário deletado com sucesso",
		Detail:  "O usuário foi removido do sistema com sucesso",
	}, nil
}
