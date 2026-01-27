package usecases

import (
	"context"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type UpdateEmailVerifiedInput struct {
	UserID   string `json:"user_id"`
	Verified bool   `json:"verified"`
}

type UpdateEmailVerifiedOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type UpdateEmailVerifiedUseCase struct {
	userRepository entities.UserRepository
	logger         logging.LoggerInterface
}

func NewUpdateEmailVerifiedUseCase(userRepo entities.UserRepository, logger logging.LoggerInterface) *UpdateEmailVerifiedUseCase {
	return &UpdateEmailVerifiedUseCase{
		userRepository: userRepo,
		logger:         logger,
	}
}

func (uc *UpdateEmailVerifiedUseCase) Execute(ctx context.Context, input UpdateEmailVerifiedInput) (*UpdateEmailVerifiedOutput, []exceptions.ProblemDetails) {
	const from = "UpdateEmailVerifiedUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", nil)
	}

	if _, err := uc.userRepository.FindByID(input.UserID); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if err := uc.userRepository.UpdateEmailVerified(input.UserID, input.Verified); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao atualizar verificação de email", err)
	}

	uc.logger.LogInfo(ctx, from, "Status de email atualizado: user_id="+input.UserID+", verified=%v")

	return &UpdateEmailVerifiedOutput{
		Message: "Status de email atualizado",
		Detail:  "A verificação de email foi atualizada com sucesso",
	}, nil
}
