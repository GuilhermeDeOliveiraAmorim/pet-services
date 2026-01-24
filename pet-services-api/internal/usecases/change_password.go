package usecases

import (
	"context"

	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type ChangePasswordInput struct {
	UserID          string `json:"user_id"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ChangePasswordOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type ChangePasswordUseCase struct {
	userRepository entities.UserRepository
}

func NewChangePasswordUseCase(userRepo entities.UserRepository) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		userRepository: userRepo,
	}
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, input ChangePasswordInput) (*ChangePasswordOutput, []exceptions.ProblemDetails) {
	const from = "ChangePasswordUseCase.Execute"

	if input.UserID == "" {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "ID do usuário ausente",
				Detail: "O ID do usuário é obrigatório",
			}),
		}
	}

	if input.CurrentPassword == "" {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "Senha atual ausente",
				Detail: "A senha atual é obrigatória",
			}),
		}
	}

	if input.NewPassword == "" {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "Nova senha ausente",
				Detail: "A nova senha é obrigatória",
			}),
		}
	}

	if !entities.IsValidPassword(input.NewPassword) {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "Senha inválida",
				Detail: "A senha deve atender aos requisitos mínimos",
			}),
		}
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.Login.DecryptPassword(input.CurrentPassword) {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
				Title:  "Senha atual incorreta",
				Detail: "A senha atual fornecida está incorreta",
			}),
		}
	}

	if err := user.Login.SetPassword(input.NewPassword); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao criptografar senha", err)
	}

	user.Updated()

	if err := uc.userRepository.Update(user); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao salvar usuário", err)
	}

	return &ChangePasswordOutput{
		Message: "Senha atualizada",
		Detail:  "Sua senha foi alterada com sucesso",
	}, nil
}
