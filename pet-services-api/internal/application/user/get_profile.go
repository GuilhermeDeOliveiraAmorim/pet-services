package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/exceptions"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// GetProfileUseCase retorna o perfil do usuário autenticado.
type GetProfileUseCase struct {
	userRepo domainUser.Repository
	logger   logging.LoggerService
}

func NewGetProfileUseCase(userRepo domainUser.Repository, logger logging.LoggerService) *GetProfileUseCase {
	return &GetProfileUseCase{userRepo: userRepo, logger: logger}
}

// GetProfileInput entrada com ID do usuário autenticado.
type GetProfileInput struct {
	UserID uuid.UUID
}

// Execute busca e retorna o perfil completo.
const GET_PROFILE_USECASE = "GET_PROFILE_USECASE"

func (uc *GetProfileUseCase) Execute(ctx context.Context, input GetProfileInput) (*domainUser.User, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    GET_PROFILE_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.UserID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    GET_PROFILE_USECASE,
			Message: "Usuário não encontrado",
			Error:   domainUser.ErrUserNotFound,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Usuário não encontrado",
			Status: exceptions.RFC404_CODE,
			Detail: "O usuário informado não foi encontrado.",
		}}
	}

	result, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    GET_PROFILE_USECASE,
			Message: "Usuário não encontrado",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Usuário não encontrado",
			Status: exceptions.RFC404_CODE,
			Detail: "O usuário informado não foi encontrado.",
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    GET_PROFILE_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	return result, nil
}
