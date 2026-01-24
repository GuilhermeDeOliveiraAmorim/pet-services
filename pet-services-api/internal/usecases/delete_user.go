package usecases

import (
	"context"
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
}

func NewDeleteUserUseCase(userRepository entities.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		userRepository: userRepository,
	}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, input DeleteUserInput) (*DeleteUserOutput, []exceptions.ProblemDetails) {
	const from = "DeleteUserUseCase.Execute"

	if input.UserID == "" {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "ID do usuário ausente",
				Detail: "O ID do usuário é obrigatório para deletar",
			}),
		}
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, []exceptions.ProblemDetails{
				exceptions.NewProblemDetails(exceptions.NotFound, exceptions.ErrorMessage{
					Title:  "Usuário não encontrado",
					Detail: "Não foi possível encontrar um usuário com o ID informado",
				}),
			}
		}
		return nil, logging.InternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.Active {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "Usuário já inativo",
				Detail: "O usuário já está inativo no sistema",
			}),
		}
	}

	user.Deactivate()

	if err := uc.userRepository.Update(user); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao deletar usuário", err)
	}

	return &DeleteUserOutput{
		Message: "Usuário deletado com sucesso",
		Detail:  "O usuário foi removido do sistema com sucesso",
	}, nil
}
