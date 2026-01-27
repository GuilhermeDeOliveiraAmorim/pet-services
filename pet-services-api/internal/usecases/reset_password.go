package usecases

import (
	"context"
	"errors"
	"time"

	"pet-services-api/internal/consts"
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
	logger               logging.LoggerInterface
}

func NewResetPasswordUseCase(userRepo entities.UserRepository, resetRepo entities.RefreshTokenRepository, logger logging.LoggerInterface) *ResetPasswordUseCase {
	return &ResetPasswordUseCase{
		userRepository:       userRepo,
		resetTokenRepository: resetRepo,
		logger:               logger,
	}
}

func (uc *ResetPasswordUseCase) Execute(ctx context.Context, input ResetPasswordInput) (*ResetPasswordOutput, []exceptions.ProblemDetails) {
	const from = "ResetPasswordUseCase.Execute"

	if input.Token == "" || input.NewPassword == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "Dados ausentes", errors.New("Token e nova senha são obrigatórios"))
	}

	if !entities.IsValidPassword(input.NewPassword) {
		return nil, uc.logger.LogBadRequest(ctx, from, "Senha inválida", errors.New("A senha deve atender aos requisitos mínimos"))
	}

	resetToken, err := uc.resetTokenRepository.FindValidPasswordResetByToken(input.Token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, uc.logger.LogNotFound(ctx, from, "Token inválido", errors.New("Token não encontrado ou inválido"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao validar token", err)
	}

	if resetToken.RevokedAt != nil || time.Now().After(resetToken.ExpiresAt) {
		return nil, uc.logger.LogBadRequest(ctx, from, "Token expirado ou revogado", errors.New("Solicite um novo reset de senha"))
	}

	user, err := uc.userRepository.FindByID(resetToken.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar um usuário com o ID informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if err := user.Login.SetPassword(input.NewPassword); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao criptografar senha", err)
	}

	if err := uc.resetTokenRepository.RevokePasswordResetByToken(input.Token); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao revogar token", err)
	}

	user.Updated()

	if err := uc.userRepository.Update(user); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao salvar usuário", err)
	}

	return &ResetPasswordOutput{
		Message: "Senha atualizada",
		Detail:  "A senha foi redefinida com sucesso",
	}, nil
}
