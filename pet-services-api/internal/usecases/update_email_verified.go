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
}

func NewUpdateEmailVerifiedUseCase(userRepo entities.UserRepository) *UpdateEmailVerifiedUseCase {
	return &UpdateEmailVerifiedUseCase{userRepository: userRepo}
}

func (uc *UpdateEmailVerifiedUseCase) Execute(ctx context.Context, input UpdateEmailVerifiedInput) (*UpdateEmailVerifiedOutput, []exceptions.ProblemDetails) {
	const from = "UpdateEmailVerifiedUseCase.Execute"

	if input.UserID == "" {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "ID do usuário ausente",
				Detail: "O ID do usuário é obrigatório",
			}),
		}
	}

	if _, err := uc.userRepository.FindByID(input.UserID); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if err := uc.userRepository.UpdateEmailVerified(input.UserID, input.Verified); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao atualizar verificação de email", err)
	}

	return &UpdateEmailVerifiedOutput{
		Message: "Status de email atualizado",
		Detail:  "A verificação de email foi atualizada com sucesso",
	}, nil
}
