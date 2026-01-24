package usecases

import (
	"context"
	"errors"
	"time"

	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"

	"gorm.io/gorm"
)

type ResetPasswordInput struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

type ResetPasswordOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type ResetPasswordUseCase struct {
	userRepository       entities.UserRepository
	resetTokenRepository entities.RefreshTokenRepository
}

func NewResetPasswordUseCase(userRepo entities.UserRepository, resetRepo entities.RefreshTokenRepository) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{
		userRepository:       userRepo,
		resetTokenRepository: resetRepo,
	}
}

func (uc *ResetPasswordUseCase) Execute(ctx context.Context, input ResetPasswordInput) (*ResetPasswordOutput, []exceptions.ProblemDetails) {
	const from = "ResetPasswordUseCase.Execute"

	if input.Token == "" || input.NewPassword == "" {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "Dados ausentes",
				Detail: "Token e nova senha são obrigatórios",
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

	resetToken, err := uc.resetTokenRepository.FindValidPasswordResetByToken(input.Token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, []exceptions.ProblemDetails{
				exceptions.NewProblemDetails(exceptions.NotFound, exceptions.ErrorMessage{
					Title:  "Token inválido",
					Detail: "Token não encontrado ou inválido",
				}),
			}
		}
		return nil, logging.InternalServerError(ctx, from, "Erro ao validar token", err)
	}

	if resetToken.RevokedAt != nil || time.Now().After(resetToken.ExpiresAt) {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "Token expirado ou revogado",
				Detail: "Solicite um novo reset de senha",
			}),
		}
	}

	user, err := uc.userRepository.FindByID(resetToken.UserID)
	if err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if err := user.Login.SetPassword(input.NewPassword); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao criptografar senha", err)
	}

	user.Updated()

	if err := uc.userRepository.Update(user); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao salvar usuário", err)
	}

	if err := uc.resetTokenRepository.RevokePasswordResetByToken(input.Token); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao revogar token", err)
	}

	return &ResetPasswordOutput{
		Message: "Senha atualizada",
		Detail:  "A senha foi redefinida com sucesso",
	}, nil
}
