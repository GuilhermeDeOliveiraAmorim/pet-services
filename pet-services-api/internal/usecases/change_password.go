package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"

	"gorm.io/gorm"
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
	logger         logging.LoggerInterface
}

func NewChangePasswordUseCase(userRepo entities.UserRepository, logger logging.LoggerInterface) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		userRepository: userRepo,
		logger:         logger,
	}
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, input ChangePasswordInput) (*ChangePasswordOutput, []exceptions.ProblemDetails) {
	const from = "ChangePasswordUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.CurrentPassword == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "Senha atual ausente", errors.New("A senha atual é obrigatória"))
	}

	if input.NewPassword == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "Nova senha ausente", errors.New("A nova senha é obrigatória"))
	}

	if !entities.IsValidPassword(input.NewPassword) {
		return nil, uc.logger.LogBadRequest(ctx, from, "Senha inválida", errors.New("A senha deve atender aos requisitos mínimos"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar um usuário com o ID informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.Login.DecryptPassword(input.CurrentPassword) {
		return nil, uc.logger.LogUnauthorized(ctx, from, "Senha atual incorreta", errors.New("A senha atual fornecida está incorreta"))
	}

	if err := user.Login.SetPassword(input.NewPassword); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao criptografar senha", err)
	}

	user.Updated()

	if err := uc.userRepository.Update(user); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao salvar usuário", err)
	}

	return &ChangePasswordOutput{
		Message: "Senha atualizada",
		Detail:  "Sua senha foi alterada com sucesso",
	}, nil
}
