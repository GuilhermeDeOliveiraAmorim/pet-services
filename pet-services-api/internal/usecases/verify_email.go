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

type VerifyEmailInput struct {
	Token string `json:"token"`
}

type VerifyEmailOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type VerifyEmailUseCase struct {
	userRepository       entities.UserRepository
	verifyTokenRepository entities.RefreshTokenRepository
}

func NewVerifyEmailUseCase(userRepo entities.UserRepository, verifyRepo entities.RefreshTokenRepository) *VerifyEmailUseCase {
	return &VerifyEmailUseCase{
		userRepository:       userRepo,
		verifyTokenRepository: verifyRepo,
	}
}

func (uc *VerifyEmailUseCase) Execute(ctx context.Context, input VerifyEmailInput) (*VerifyEmailOutput, []exceptions.ProblemDetails) {
	const from = "VerifyEmailUseCase.Execute"

	if input.Token == "" {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "Token ausente",
				Detail: "O token de verificação é obrigatório",
			}),
		}
	}

	verifyToken, err := uc.verifyTokenRepository.FindValidPasswordResetByToken(input.Token)
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

	if verifyToken.RevokedAt != nil || time.Now().After(verifyToken.ExpiresAt) {
		return nil, []exceptions.ProblemDetails{
			exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
				Title:  "Token expirado ou revogado",
				Detail: "Solicite um novo email de verificação",
			}),
		}
	}

	user, err := uc.userRepository.FindByID(verifyToken.UserID)
	if err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if user.EmailVerified {
		return &VerifyEmailOutput{
			Message: "Email já verificado",
			Detail:  "Este email já foi verificado anteriormente",
		}, nil
	}

	if err := uc.userRepository.UpdateEmailVerified(user.ID, true); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao atualizar email", err)
	}

	if err := uc.verifyTokenRepository.RevokePasswordResetByToken(input.Token); err != nil {
		return nil, logging.InternalServerError(ctx, from, "Erro ao revogar token", err)
	}

	return &VerifyEmailOutput{
		Message: "Email verificado com sucesso",
		Detail:  "Sua conta foi ativada completamente",
	}, nil
}
