package user

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"pet-services-api/internal/application/exceptions"
	"pet-services-api/internal/application/logging"
	domainUser "pet-services-api/internal/domain/user"
)

// RequestEmailVerificationUseCase envia email de verificação.
type RequestEmailVerificationUseCase struct {
	userRepo            domainUser.Repository
	emailService        domainUser.EmailService
	verificationBaseURL string
	logger              logging.LoggerService
}

func NewRequestEmailVerificationUseCase(userRepo domainUser.Repository, emailService domainUser.EmailService, verificationBaseURL string, logger logging.LoggerService) *RequestEmailVerificationUseCase {
	return &RequestEmailVerificationUseCase{
		userRepo:            userRepo,
		emailService:        emailService,
		verificationBaseURL: verificationBaseURL,
		logger:              logger,
	}
}

// RequestEmailVerificationInput dados para solicitar verificação.
type RequestEmailVerificationInput struct {
	UserID uuid.UUID
}

// Execute gera token, persiste e envia email de verificação.
const REQUEST_EMAIL_VERIFICATION_USECASE = "REQUEST_EMAIL_VERIFICATION_USECASE"

func (uc *RequestEmailVerificationUseCase) Execute(ctx context.Context, input RequestEmailVerificationInput) []exceptions.ProblemDetails {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    REQUEST_EMAIL_VERIFICATION_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.UserID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    REQUEST_EMAIL_VERIFICATION_USECASE,
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

	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    REQUEST_EMAIL_VERIFICATION_USECASE,
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

	if user.EmailVerified {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC409_CODE,
			From:    REQUEST_EMAIL_VERIFICATION_USECASE,
			Message: "Email já verificado",
			Error:   domainUser.ErrEmailAlreadyVerified,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC409,
			Title:  "Email já verificado",
			Status: exceptions.RFC409_CODE,
			Detail: "O email já foi verificado anteriormente.",
		}}
	}

	token, err := generateSecureToken(32)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    REQUEST_EMAIL_VERIFICATION_USECASE,
			Message: "Falha ao gerar token de verificação",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao gerar token de verificação",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	verificationToken := domainUser.NewEmailVerificationToken(user.ID, token, expiresAt)
	if err := uc.userRepo.CreateEmailVerificationToken(ctx, verificationToken); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    REQUEST_EMAIL_VERIFICATION_USECASE,
			Message: "Falha ao salvar token de verificação",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao salvar token de verificação",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	verificationLink := fmt.Sprintf("%s?token=%s", uc.verificationBaseURL, token)
	if err := uc.emailService.SendEmailVerification(user.Email.String(), verificationLink); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    REQUEST_EMAIL_VERIFICATION_USECASE,
			Message: "Falha ao enviar email de verificação",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao enviar email de verificação",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    REQUEST_EMAIL_VERIFICATION_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	return nil
}
