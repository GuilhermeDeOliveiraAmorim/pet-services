package usecases

import (
	"context"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"

	"gorm.io/gorm"
)

type GetUserByIDInput struct {
	UserID string `json:"user_id"`
}

type GetUserByIDOutput struct {
	User *entities.User `json:"user"`
}

type GetUserByIDUseCase struct {
	userRepository entities.UserRepository
}

func NewGetUserByIDUseCase(userRepository entities.UserRepository) *GetUserByIDUseCase {
	return &GetUserByIDUseCase{
		userRepository: userRepository,
	}
}

func (uc *GetUserByIDUseCase) Execute(ctx context.Context, input GetUserByIDInput) (*GetUserByIDOutput, []exceptions.ProblemDetails) {
	const from = "GetUserByIDUseCase.Execute"

	if input.UserID == "" {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "ID do usuário ausente",
				Detail: "O ID do usuário é obrigatório para buscar",
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

	return &GetUserByIDOutput{
		User: user,
	}, nil
}
