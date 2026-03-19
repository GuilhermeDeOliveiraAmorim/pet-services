package usecases

import (
	"context"
	"errors"

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
	logger         logging.LoggerInterface
}

func NewCheckEmailExistsUseCase(userRepo entities.UserRepository, logger logging.LoggerInterface) *CheckEmailExistsUseCase {
	return &CheckEmailExistsUseCase{
		userRepository: userRepo,
		logger:         logger,
	}
}

func (uc *CheckEmailExistsUseCase) Execute(ctx context.Context, input CheckEmailExistsInput) (*CheckEmailExistsOutput, []exceptions.ProblemDetails) {
	const from = "CheckEmailExistsUseCase.Execute"

	if input.Email == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "Email ausente", errors.New("O email é obrigatório"))
	}

	if !entities.IsValidEmail(input.Email) {
		return nil, uc.logger.LogBadRequest(ctx, from, "Email inválido", errors.New("O formato do email está inválido"))
	}

	exists, err := uc.userRepository.ExistsByEmail(input.Email)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao verificar email", err)
	}

	message := "Solicitação processada"

	return &CheckEmailExistsOutput{
		Exists:  exists,
		Message: message,
	}, nil
}
