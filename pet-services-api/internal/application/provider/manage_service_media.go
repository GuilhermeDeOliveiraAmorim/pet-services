package provider

import (
	"context"
	"fmt"
	"pet-services-api/internal/application/exceptions"
	"pet-services-api/internal/application/logging"
	"pet-services-api/internal/domain/provider"

	"github.com/google/uuid"
)

type AddServicePhotoInput struct {
	ServiceID uuid.UUID
	URL       string
}

type AddServicePhotoUseCase struct {
	logger       logging.LoggerService
	providerRepo provider.Repository
}

func (uc *AddServicePhotoUseCase) Execute(ctx context.Context, input AddServicePhotoInput) (*provider.Provider, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "AddServicePhotoUseCase",
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.ServiceID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "AddServicePhotoUseCase",
			Message: "ServiceID é obrigatório",
			Error:   fmt.Errorf("serviceID é obrigatório"),
		})
		return nil, []exceptions.ProblemDetails{
			{
				Type:   exceptions.RFC400,
				Title:  "ServiceID é obrigatório",
				Status: exceptions.RFC400_CODE,
				Detail: "O ID do serviço é obrigatório.",
			},
		}
	}

	service, err := uc.providerRepo.FindByID(ctx, input.ServiceID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    "AddServicePhotoUseCase",
			Message: "Serviço não encontrado",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{
			{
				Type:   exceptions.RFC404,
				Title:  "Serviço não encontrado",
				Status: exceptions.RFC404_CODE,
				Detail: "O serviço informado não foi encontrado.",
			},
		}
	}

	if err := service.AddPhoto(input.URL); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "AddServicePhotoUseCase",
			Message: "Erro ao adicionar foto ao serviço",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{
			{
				Type:   exceptions.RFC400,
				Title:  "Erro ao adicionar foto ao serviço",
				Status: exceptions.RFC400_CODE,
				Detail: err.Error(),
			},
		}
	}

	if err := uc.providerRepo.Update(ctx, service); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    "AddServicePhotoUseCase",
			Message: "Erro ao persistir serviço com foto",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{
			{
				Type:   exceptions.RFC500,
				Title:  "Erro ao persistir serviço com foto",
				Status: exceptions.RFC500_CODE,
				Detail: err.Error(),
			},
		}
	}

	return service, nil
}
