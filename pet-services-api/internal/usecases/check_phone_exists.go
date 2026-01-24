package usecases

import (
	"context"

	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type CheckPhoneExistsInput struct {
	CountryCode string `json:"country_code"`
	AreaCode    string `json:"area_code"`
	Number      string `json:"number"`
}

type CheckPhoneExistsOutput struct {
	Exists  bool   `json:"exists"`
	Message string `json:"message,omitempty"`
}

type CheckPhoneExistsUseCase struct {
	userRepository entities.UserRepository
}

func NewCheckPhoneExistsUseCase(userRepo entities.UserRepository) *CheckPhoneExistsUseCase {
	return &CheckPhoneExistsUseCase{
		userRepository: userRepo,
	}
}

func (uc *CheckPhoneExistsUseCase) Execute(ctx context.Context, input CheckPhoneExistsInput) (*CheckPhoneExistsOutput, []exceptions.ProblemDetails) {
	const from = "CheckPhoneExistsUseCase.Execute"

	if input.CountryCode == "" || input.AreaCode == "" || input.Number == "" {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "Dados do telefone incompletos",
				Detail: "Código do país, código de área e número são obrigatórios",
			}),
		}
	}

	_, validationErrors := entities.NewPhone(input.CountryCode, input.AreaCode, input.Number)
	if len(validationErrors) > 0 {
		return nil, validationErrors
	}

	exists, err := uc.userRepository.ExistsByPhone(input.CountryCode, input.AreaCode, input.Number)
	if err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao verificar telefone", err)
	}

	message := "Telefone disponível"
	if exists {
		message = "Telefone já cadastrado"
	}

	return &CheckPhoneExistsOutput{
		Exists:  exists,
		Message: message,
	}, nil
}
