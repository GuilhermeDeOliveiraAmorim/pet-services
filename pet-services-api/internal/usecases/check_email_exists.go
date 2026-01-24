package usecases

import (
	"context"

	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type CheckEmailExistsInput struct {
	Email string `json:"email"`
}

type CheckEmailExistsOutput struct {
	Exists  bool   `json:"exists"`
	Message string `json:"message,omitempty"`
}

type CheckEmailExistsUseCase struct {
	userRepository entities.UserRepository
}

func NewCheckEmailExistsUseCase(userRepo entities.UserRepository) *CheckEmailExistsUseCase {
	return &CheckEmailExistsUseCase{
		userRepository: userRepo,
	}
}

func (uc *CheckEmailExistsUseCase) Execute(ctx context.Context, input CheckEmailExistsInput) (*CheckEmailExistsOutput, []exceptions.ProblemDetails) {
	const from = "CheckEmailExistsUseCase.Execute"

	if input.Email == "" {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "Email ausente",
				Detail: "O email é obrigatório",
			}),
		}
	}

	if !entities.IsValidEmail(input.Email) {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "Email inválido",
				Detail: "O formato do email está inválido",
			}),
		}
	}

	exists, err := uc.userRepository.ExistsByEmail(input.Email)
	if err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao verificar email", err)
	}

	message := "Email disponível"
	if exists {
		message = "Email já cadastrado"
	}

	return &CheckEmailExistsOutput{
		Exists:  exists,
		Message: message,
	}, nil
}
