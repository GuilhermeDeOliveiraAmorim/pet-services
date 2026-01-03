package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/exceptions"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainAuth "github.com/guilherme/pet-services-api/internal/domain/auth"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// DeleteAccountUseCase remove ou desativa conta do usuário.
type DeleteAccountUseCase struct {
	userRepo    domainUser.Repository
	refreshRepo domainAuth.RefreshTokenRepository
	logger      logging.LoggerService
}

func NewDeleteAccountUseCase(userRepo domainUser.Repository, refreshRepo domainAuth.RefreshTokenRepository, logger logging.LoggerService) *DeleteAccountUseCase {
	return &DeleteAccountUseCase{
		userRepo:    userRepo,
		refreshRepo: refreshRepo,
		logger:      logger,
	}
}

// DeleteAccountInput dados para deletar conta.
type DeleteAccountInput struct {
	UserID     uuid.UUID
	HardDelete bool // se true, deleta permanentemente; se false, soft delete
}

// Execute marca usuário como deletado (soft) ou remove permanentemente (hard).
const DELETE_ACCOUNT_USECASE = "DELETE_ACCOUNT_USECASE"

func (uc *DeleteAccountUseCase) Execute(ctx context.Context, input DeleteAccountInput) []exceptions.ProblemDetails {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    DELETE_ACCOUNT_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.UserID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    DELETE_ACCOUNT_USECASE,
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
			From:    DELETE_ACCOUNT_USECASE,
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

	if user.IsDeleted() {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC409_CODE,
			From:    DELETE_ACCOUNT_USECASE,
			Message: "Conta já foi deletada",
			Error:   errors.New("conta já foi deletada"),
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC409,
			Title:  "Conta já deletada",
			Status: exceptions.RFC409_CODE,
			Detail: "A conta já foi deletada anteriormente.",
		}}
	}

	_ = uc.refreshRepo.RevokeAllByUserID(ctx, user.ID)

	if input.HardDelete {
		if err := uc.userRepo.Delete(ctx, user.ID); err != nil {
			uc.logger.Log(logging.Logger{
				Context: ctx,
				TypeLog: logging.LoggerTypes.ERROR,
				Layer:   logging.LoggerLayers.USECASES,
				Code:    exceptions.RFC500_CODE,
				From:    DELETE_ACCOUNT_USECASE,
				Message: "Falha ao deletar conta",
				Error:   err,
			})
			return []exceptions.ProblemDetails{{
				Type:   exceptions.RFC500,
				Title:  "Falha ao deletar conta",
				Status: exceptions.RFC500_CODE,
				Detail: err.Error(),
			}}
		}
	} else {
		user.SoftDelete()
		if err := uc.userRepo.Update(ctx, user); err != nil {
			uc.logger.Log(logging.Logger{
				Context: ctx,
				TypeLog: logging.LoggerTypes.ERROR,
				Layer:   logging.LoggerLayers.USECASES,
				Code:    exceptions.RFC500_CODE,
				From:    DELETE_ACCOUNT_USECASE,
				Message: "Falha ao desativar conta",
				Error:   err,
			})
			return []exceptions.ProblemDetails{{
				Type:   exceptions.RFC500,
				Title:  "Falha ao desativar conta",
				Status: exceptions.RFC500_CODE,
				Detail: err.Error(),
			}}
		}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    DELETE_ACCOUNT_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	return nil
}
