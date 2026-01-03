package provider

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/exceptions"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// RemoveServiceUseCase remove um serviço existente do prestador.
type RemoveServiceUseCase struct {
	providerRepo provider.Repository
	logger       logging.LoggerService
}

// NewRemoveServiceUseCase cria nova instância do caso de uso.
func NewRemoveServiceUseCase(providerRepo provider.Repository, logger *slog.Logger) *RemoveServiceUseCase {
	return &RemoveServiceUseCase{providerRepo: providerRepo, logger: logging.NewSlogLogger()}
}

// RemoveServiceInput representa os dados para remover serviço.
type RemoveServiceInput struct {
	ProviderID uuid.UUID
	Category   string
	Name       string
}

// Execute remove o serviço do prestador seguindo padrão CreateRequestUseCase.
func (uc *RemoveServiceUseCase) Execute(ctx context.Context, input RemoveServiceInput) (*provider.Service, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "RemoveServiceUseCase",
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.ProviderID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "RemoveServiceUseCase",
			Message: "ProviderID é obrigatório",
			Error:   fmt.Errorf("providerID é obrigatório"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "ProviderID é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O ID do prestador é obrigatório.",
		}}
	}
	if input.Name == "" {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "RemoveServiceUseCase",
			Message: "Nome do serviço é obrigatório",
			Error:   provider.NewValidationError("service.name", "nome do serviço é obrigatório"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Nome do serviço é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O nome do serviço é obrigatório.",
		}}
	}

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    "RemoveServiceUseCase",
			Message: "Prestador não encontrado",
			Error:   provider.ErrProviderNotFound,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Prestador não encontrado",
			Status: exceptions.RFC404_CODE,
			Detail: "O prestador informado não foi encontrado.",
		}}
	}

	var removedService *provider.Service
	before := len(p.Services)
	for _, svc := range p.Services {
		if svc.Category == input.Category && svc.Name == input.Name {
			removedService = &svc
			break
		}
	}
	p.RemoveService(input.Category, input.Name)

	if len(p.Services) == before {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    "RemoveServiceUseCase",
			Message: "Serviço não encontrado",
			Error:   provider.ErrServiceNotFound,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Serviço não encontrado",
			Status: exceptions.RFC404_CODE,
			Detail: "O serviço informado não foi encontrado.",
		}}
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    "RemoveServiceUseCase",
			Message: "Falha ao remover serviço",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao remover serviço",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "RemoveServiceUseCase",
		Message: logging.DEFAULTMESSAGES.END,
	})

	return removedService, nil
}
