package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/mail"
)

type ChangePasswordInput struct {
	UserID      string `json:"user_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ChangePasswordInputBody struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type ChangePasswordOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type ChangePasswordUseCase struct {
	userRepository entities.UserRepository
	emailService   mail.EmailService
	logger         logging.LoggerInterface
}

func NewChangePasswordUseCase(userRepo entities.UserRepository, emailService mail.EmailService, logger logging.LoggerInterface) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		userRepository: userRepo,
		emailService:   emailService,
		logger:         logger,
	}
}

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, input ChangePasswordInput) (*ChangePasswordOutput, []exceptions.ProblemDetails) {
	const from = "ChangePasswordUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.OldPassword == "" {
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
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar um usuário com o ID informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.Login.DecryptPassword(input.OldPassword) {
		return nil, uc.logger.LogUnauthorized(ctx, from, "Senha atual incorreta", errors.New("A senha atual fornecida está incorreta"))
	}

	if err := user.Login.SetPassword(input.NewPassword); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao criptografar senha", err)
	}

	user.Updated()

	if err := uc.userRepository.Update(user); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao salvar usuário", err)
	}

	if err := uc.emailService.SendPasswordChangedAlertEmail(user.Login.Email, user.Name); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao enviar alerta de segurança por email", err)
	}

	return &ChangePasswordOutput{
		Message: "Senha atualizada",
		Detail:  "Sua senha foi alterada com sucesso",
	}, nil
}
