package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"pet-services-api/internal/application/exceptions"
	"pet-services-api/internal/application/logging"
	domainUser "pet-services-api/internal/domain/user"
)

// RequestPasswordResetUseCase inicia o fluxo de reset de senha.
type RequestPasswordResetUseCase struct {
	userRepo     domainUser.Repository
	emailService domainUser.EmailService
	resetBaseURL string
	logger       logging.LoggerService
}

func NewRequestPasswordResetUseCase(userRepo domainUser.Repository, emailService domainUser.EmailService, resetBaseURL string, logger logging.LoggerService) *RequestPasswordResetUseCase {
	return &RequestPasswordResetUseCase{
		userRepo:     userRepo,
		emailService: emailService,
		resetBaseURL: resetBaseURL,
		logger:       logger,
	}
}

// RequestPasswordResetInput dados para solicitar reset.
type RequestPasswordResetInput struct {
	Email string
}

// Execute gera token, persiste e envia email.
const REQUEST_PASSWORD_RESET_USECASE = "REQUEST_PASSWORD_RESET_USECASE"

func (uc *RequestPasswordResetUseCase) Execute(ctx context.Context, input RequestPasswordResetInput) []exceptions.ProblemDetails {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    REQUEST_PASSWORD_RESET_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	email := strings.TrimSpace(strings.ToLower(input.Email))
	if email == "" {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    REQUEST_PASSWORD_RESET_USECASE,
			Message: "Email é obrigatório",
			Error:   domainUser.ErrInvalidEmail,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Email é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O email é obrigatório.",
		}}
	}

	user, err := uc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Por segurança, não revela se o email existe ou não
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.INFO,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC200_CODE,
			From:    REQUEST_PASSWORD_RESET_USECASE,
			Message: "Solicitação de reset processada (email pode não existir)",
		})
		return nil
	}

	token, err := generateSecureToken(32)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    REQUEST_PASSWORD_RESET_USECASE,
			Message: "Falha ao gerar token de reset",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao gerar token de reset",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	expiresAt := time.Now().Add(1 * time.Hour)
	resetToken := domainUser.NewPasswordResetToken(user.ID, token, expiresAt)
	if err := uc.userRepo.CreatePasswordResetToken(ctx, resetToken); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    REQUEST_PASSWORD_RESET_USECASE,
			Message: "Falha ao salvar token de reset",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao salvar token de reset",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	resetLink := fmt.Sprintf("%s?token=%s", uc.resetBaseURL, token)
	if err := uc.emailService.SendPasswordResetEmail(user.Email.String(), resetLink); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    REQUEST_PASSWORD_RESET_USECASE,
			Message: "Falha ao enviar email de reset",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao enviar email de reset",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    REQUEST_PASSWORD_RESET_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	return nil
}

// generateSecureToken gera um token aleatório seguro.
func generateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
