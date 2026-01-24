package usecases

import (
	"context"

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
}

func NewReactivateUserUseCase(userRepo entities.UserRepository) *ReactivateUserUseCase {
	return &ReactivateUserUseCase{
		userRepository: userRepo,
	}
}

func (uc *ReactivateUserUseCase) Execute(ctx context.Context, input ReactivateUserInput) (*ReactivateUserOutput, []exceptions.ProblemDetails) {
	const from = "ReactivateUserUseCase.Execute"

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

	if user.Active {
		return &ReactivateUserOutput{
			Message: "Conta já ativa",
			Detail:  "Esta conta já está ativa",
		}, nil
	}

	user.Activate()

	if err := uc.userRepository.Update(user); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao reativar usuário", err)
	}

	return &ReactivateUserOutput{
		Message: "Conta reativada com sucesso",
		Detail:  "Sua conta foi reativada e você pode fazer login novamente",
	}, nil
}
