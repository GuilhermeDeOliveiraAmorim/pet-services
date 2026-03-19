package usecases

import (
	"context"
	"errors"

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
	logger         logging.LoggerInterface
}

func NewCheckPhoneExistsUseCase(userRepo entities.UserRepository, logger logging.LoggerInterface) *CheckPhoneExistsUseCase {
	return &CheckPhoneExistsUseCase{
		userRepository: userRepo,
		logger:         logger,
	}
}

func (uc *CheckPhoneExistsUseCase) Execute(ctx context.Context, input CheckPhoneExistsInput) (*CheckPhoneExistsOutput, []exceptions.ProblemDetails) {
	const from = "CheckPhoneExistsUseCase.Execute"

	if input.CountryCode == "" || input.AreaCode == "" || input.Number == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "Dados do telefone incompletos", errors.New("Código do país, código de área e número são obrigatórios"))
	}

	_, validationErrors := entities.NewPhone(input.CountryCode, input.AreaCode, input.Number)
	if len(validationErrors) > 0 {
		var loggedErrors []exceptions.ProblemDetails
		for _, err := range validationErrors {
			loggedErrors = append(loggedErrors, uc.logger.LogBadRequest(ctx, from, "Erro de validação de telefone", errors.New(err.Detail))...)
		}
		return nil, loggedErrors
	}

	exists, err := uc.userRepository.ExistsByPhone(input.CountryCode, input.AreaCode, input.Number)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao verificar telefone", err)
	}

	message := "Solicitação processada"

	return &CheckPhoneExistsOutput{
		Exists:  exists,
		Message: message,
	}, nil
}
