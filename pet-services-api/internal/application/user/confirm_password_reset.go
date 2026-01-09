package user

import (
	"context"
	"strings"

	"pet-services-api/internal/application/exceptions"
	"pet-services-api/internal/application/logging"
	domainAuth "pet-services-api/internal/domain/auth"
	domainUser "pet-services-api/internal/domain/user"
)

// ConfirmPasswordResetUseCase confirma e aplica o reset de senha.
type ConfirmPasswordResetUseCase struct {
	userRepo       domainUser.Repository
	refreshRepo    domainAuth.RefreshTokenRepository
	passwordHasher domainAuth.PasswordHasher
	logger         logging.LoggerService
}

func NewConfirmPasswordResetUseCase(userRepo domainUser.Repository, refreshRepo domainAuth.RefreshTokenRepository, passwordHasher domainAuth.PasswordHasher, logger logging.LoggerService) *ConfirmPasswordResetUseCase {
	return &ConfirmPasswordResetUseCase{
		userRepo:       userRepo,
		refreshRepo:    refreshRepo,
		passwordHasher: passwordHasher,
		logger:         logger,
	}
}

// ConfirmPasswordResetInput dados para confirmar reset.
type ConfirmPasswordResetInput struct {
	Token       string
	NewPassword string
}

// Execute valida token, atualiza senha e revoga tokens existentes.
const CONFIRM_PASSWORD_RESET_USECASE = "CONFIRM_PASSWORD_RESET_USECASE"

func (uc *ConfirmPasswordResetUseCase) Execute(ctx context.Context, input ConfirmPasswordResetInput) []exceptions.ProblemDetails {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    CONFIRM_PASSWORD_RESET_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	token := strings.TrimSpace(input.Token)
	newPassword := strings.TrimSpace(input.NewPassword)

	if token == "" {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    CONFIRM_PASSWORD_RESET_USECASE,
			Message: "Token de reset é obrigatório",
			Error:   domainUser.ErrPasswordResetTokenInvalid,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Token de reset é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O token de reset é obrigatório.",
		}}
	}

	if newPassword == "" || len(newPassword) < 6 {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    CONFIRM_PASSWORD_RESET_USECASE,
			Message: "Nova senha inválida",
			Error:   domainUser.ErrInvalidPassword,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Nova senha inválida",
			Status: exceptions.RFC400_CODE,
			Detail: "A nova senha deve ter no mínimo 6 caracteres.",
		}}
	}

	resetToken, err := uc.userRepo.FindPasswordResetToken(ctx, token)
	if err != nil || !resetToken.IsValid() {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    CONFIRM_PASSWORD_RESET_USECASE,
			Message: "Token de reset inválido ou expirado",
			Error:   domainUser.ErrPasswordResetTokenInvalid,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Token de reset inválido",
			Status: exceptions.RFC400_CODE,
			Detail: "O token de reset é inválido ou expirou.",
		}}
	}

	user, err := uc.userRepo.FindByID(ctx, resetToken.UserID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    CONFIRM_PASSWORD_RESET_USECASE,
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

	hashed, err := uc.passwordHasher.Hash(newPassword)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    CONFIRM_PASSWORD_RESET_USECASE,
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
			From:    CONFIRM_PASSWORD_RESET_USECASE,
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

	if err := uc.userRepo.MarkPasswordResetTokenAsUsed(ctx, resetToken.ID); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    CONFIRM_PASSWORD_RESET_USECASE,
			Message: "Falha ao marcar token como usado",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao marcar token como usado",
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
		From:    CONFIRM_PASSWORD_RESET_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	return nil
}
