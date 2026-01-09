package user

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"pet-services-api/internal/application/exceptions"
	"pet-services-api/internal/application/logging"
	domainAuth "pet-services-api/internal/domain/auth"
	domainUser "pet-services-api/internal/domain/user"
)

// ChangePasswordUseCase permite ao usuário trocar sua senha.
type ChangePasswordUseCase struct {
	userRepo       domainUser.Repository
	refreshRepo    domainAuth.RefreshTokenRepository
	passwordHasher domainAuth.PasswordHasher
	logger         logging.LoggerService
}

func NewChangePasswordUseCase(userRepo domainUser.Repository, refreshRepo domainAuth.RefreshTokenRepository, passwordHasher domainAuth.PasswordHasher, logger logging.LoggerService) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		userRepo:       userRepo,
		refreshRepo:    refreshRepo,
		passwordHasher: passwordHasher,
		logger:         logger,
	}
}

// ChangePasswordInput dados para troca de senha.
type ChangePasswordInput struct {
	UserID          uuid.UUID
	CurrentPassword string
	NewPassword     string
}

// Execute valida senha atual, atualiza com nova e revoga tokens existentes.
const CHANGE_PASSWORD_USECASE = "CHANGE_PASSWORD_USECASE"

func (uc *ChangePasswordUseCase) Execute(ctx context.Context, input ChangePasswordInput) []exceptions.ProblemDetails {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    CHANGE_PASSWORD_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.UserID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    CHANGE_PASSWORD_USECASE,
			Message: "Usuário não encontrado",
			Error:   domainUser.ErrUserNotFound,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Usuário não encontrado",
			Status: exceptions.RFC404_CODE,
			Detail: "O usuário informado não foi encontrado.",
		}}
	}

	currentPassword := strings.TrimSpace(input.CurrentPassword)
	newPassword := strings.TrimSpace(input.NewPassword)

	if currentPassword == "" || newPassword == "" {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    CHANGE_PASSWORD_USECASE,
			Message: "Senha atual e nova senha são obrigatórias",
			Error:   domainUser.ErrInvalidPassword,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Senha obrigatória",
			Status: exceptions.RFC400_CODE,
			Detail: "A senha atual e a nova senha são obrigatórias.",
		}}
	}

	if len(newPassword) < 6 {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    CHANGE_PASSWORD_USECASE,
			Message: "Nova senha deve ter no mínimo 6 caracteres",
			Error:   errors.New("nova senha deve ter no mínimo 6 caracteres"),
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Senha muito curta",
			Status: exceptions.RFC400_CODE,
			Detail: "A nova senha deve ter no mínimo 6 caracteres.",
		}}
	}

	if currentPassword == newPassword {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    CHANGE_PASSWORD_USECASE,
			Message: "Nova senha deve ser diferente da atual",
			Error:   errors.New("nova senha deve ser diferente da atual"),
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Senha igual à anterior",
			Status: exceptions.RFC400_CODE,
			Detail: "A nova senha deve ser diferente da atual.",
		}}
	}

	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    CHANGE_PASSWORD_USECASE,
			Message: "Usuário não encontrado",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Usuário não encontrado",
			Status: exceptions.RFC404_CODE,
			Detail: "O usuário informado não foi encontrado.",
		}}
	}

	if err := uc.passwordHasher.Compare(user.Password, currentPassword); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC401_CODE,
			From:    CHANGE_PASSWORD_USECASE,
			Message: "Senha atual incorreta",
			Error:   domainAuth.ErrInvalidCredentials,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC401,
			Title:  "Senha atual incorreta",
			Status: exceptions.RFC401_CODE,
			Detail: "A senha atual informada está incorreta.",
		}}
	}

	hashed, err := uc.passwordHasher.Hash(newPassword)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    CHANGE_PASSWORD_USECASE,
			Message: "Falha ao gerar hash da nova senha",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao atualizar senha",
			Status: exceptions.RFC500_CODE,
			Detail: "Erro ao gerar hash da nova senha.",
		}}
	}

	user.SetPassword(hashed)

	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    CHANGE_PASSWORD_USECASE,
			Message: "Falha ao atualizar senha",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao atualizar senha",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	_ = uc.refreshRepo.RevokeAllByUserID(ctx, user.ID)

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    CHANGE_PASSWORD_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	return nil
}
