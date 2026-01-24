package usecases

import (
	"context"

	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type GetProfileInput struct {
	UserID string `json:"user_id"`
}

type GetProfileOutput struct {
	User *entities.User `json:"user"`
}

type GetProfileUseCase struct {
	userRepository entities.UserRepository
}

func NewGetProfileUseCase(userRepo entities.UserRepository) *GetProfileUseCase {
	return &GetProfileUseCase{
		userRepository: userRepo,
	}
}

func (uc *GetProfileUseCase) Execute(ctx context.Context, input GetProfileInput) (*GetProfileOutput, []exceptions.ProblemDetails) {
	const from = "GetProfileUseCase.Execute"

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
		return nil, logging.InternalServerError(ctx, from, "Erro ao buscar perfil do usuário", err)
	}

	return &GetProfileOutput{
		User: user,
	}, nil
}
