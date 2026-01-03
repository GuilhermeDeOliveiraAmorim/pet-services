package user

import (
	"context"
	"time"

	"github.com/guilherme/pet-services-api/internal/application/exceptions"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// ConfirmEmailVerificationUseCase confirma a verificação de email.
type ConfirmEmailVerificationUseCase struct {
	userRepo domainUser.Repository
	logger   logging.LoggerService
}

func NewConfirmEmailVerificationUseCase(userRepo domainUser.Repository, logger logging.LoggerService) *ConfirmEmailVerificationUseCase {
	return &ConfirmEmailVerificationUseCase{
		userRepo: userRepo,
		logger:   logger,
	}
}

// ConfirmEmailVerificationInput dados para confirmar verificação.
type ConfirmEmailVerificationInput struct {
	Token string
}

// Execute valida o token e marca o email como verificado.
const CONFIRM_EMAIL_VERIFICATION_USECASE = "CONFIRM_EMAIL_VERIFICATION_USECASE"

func (uc *ConfirmEmailVerificationUseCase) Execute(ctx context.Context, input ConfirmEmailVerificationInput) []exceptions.ProblemDetails {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    CONFIRM_EMAIL_VERIFICATION_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.Token == "" {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    CONFIRM_EMAIL_VERIFICATION_USECASE,
			Message: "Token de verificação é obrigatório",
			Error:   domainUser.ErrEmailVerificationTokenInvalid,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Token de verificação é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O token de verificação é obrigatório.",
		}}
	}

	verificationToken, err := uc.userRepo.FindEmailVerificationToken(ctx, input.Token)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    CONFIRM_EMAIL_VERIFICATION_USECASE,
			Message: "Token de verificação inválido",
			Error:   domainUser.ErrEmailVerificationTokenInvalid,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Token de verificação inválido",
			Status: exceptions.RFC400_CODE,
			Detail: "O token de verificação é inválido ou expirou.",
		}}
	}

	if time.Now().After(verificationToken.ExpiresAt) {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    CONFIRM_EMAIL_VERIFICATION_USECASE,
			Message: "Token de verificação expirado",
			Error:   domainUser.ErrEmailVerificationTokenInvalid,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Token de verificação expirado",
			Status: exceptions.RFC400_CODE,
			Detail: "O token de verificação expirou.",
		}}
	}

	user, err := uc.userRepo.FindByID(ctx, verificationToken.UserID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    CONFIRM_EMAIL_VERIFICATION_USECASE,
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
			From:    CONFIRM_EMAIL_VERIFICATION_USECASE,
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

	user.VerifyEmail()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    CONFIRM_EMAIL_VERIFICATION_USECASE,
			Message: "Falha ao atualizar usuário",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao atualizar usuário",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	_ = uc.userRepo.MarkEmailVerificationTokenAsUsed(ctx, verificationToken.ID)

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    CONFIRM_EMAIL_VERIFICATION_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	return nil
}
