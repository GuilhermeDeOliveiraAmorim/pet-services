package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/exceptions"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// AddServiceUseCase adiciona um serviço ao prestador.
type AddServiceUseCase struct {
	providerRepo provider.Repository
	logger       logging.LoggerService
}

// NewAddServiceUseCase cria uma nova instância do caso de uso.
func NewAddServiceUseCase(providerRepo provider.Repository, logger logging.LoggerService) *AddServiceUseCase {
	return &AddServiceUseCase{providerRepo: providerRepo, logger: logging.NewSlogLogger()}
}

// AddServiceInput representa os dados para adicionar serviço.
type AddServiceInput struct {
	ProviderID uuid.UUID
	Category   string
	Name       string
	PriceMin   float64
	PriceMax   float64
}

// Execute valida e adiciona o serviço ao prestador seguindo padrão CreateRequestUseCase.
func (uc *AddServiceUseCase) Execute(ctx context.Context, input AddServiceInput) (*provider.Service, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "AddServiceUseCase",
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.ProviderID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "AddServiceUseCase",
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
			From:    "AddServiceUseCase",
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
	if input.PriceMin < 0 || input.PriceMax < 0 {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "AddServiceUseCase",
			Message: "Preço não pode ser negativo",
			Error:   provider.NewValidationError("service.price", "preço não pode ser negativo"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Preço não pode ser negativo",
			Status: exceptions.RFC400_CODE,
			Detail: "O preço não pode ser negativo.",
		}}
	}
	if input.PriceMax > 0 && input.PriceMin > input.PriceMax {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "AddServiceUseCase",
			Message: "Preço máximo deve ser maior ou igual ao mínimo",
			Error:   provider.NewValidationError("service.price", "preço máximo deve ser maior ou igual ao mínimo"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Preço máximo deve ser maior ou igual ao mínimo",
			Status: exceptions.RFC400_CODE,
			Detail: "O preço máximo deve ser maior ou igual ao mínimo.",
		}}
	}

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    "AddServiceUseCase",
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

	// Checagem simples para evitar duplicidade (categoria + nome)
	for _, svc := range p.Services {
		if svc.Category == input.Category && svc.Name == input.Name {
			uc.logger.Log(logging.Logger{
				Context: ctx,
				TypeLog: logging.LoggerTypes.ERROR,
				Layer:   logging.LoggerLayers.USECASES,
				Code:    exceptions.RFC409_CODE,
				From:    "AddServiceUseCase",
				Message: "Serviço já cadastrado",
				Error:   provider.NewValidationError("service", "serviço já cadastrado"),
			})
			return nil, []exceptions.ProblemDetails{{
				Type:   exceptions.RFC409,
				Title:  "Serviço já cadastrado",
				Status: exceptions.RFC409_CODE,
				Detail: "O serviço já está cadastrado para este prestador.",
			}}
		}
	}

	p.AddService(input.Category, input.Name, input.PriceMin, input.PriceMax)
	var addedService *provider.Service
	for i, svc := range p.Services {
		if svc.Category == input.Category && svc.Name == input.Name {
			addedService = &p.Services[i]
			break
		}
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    "AddServiceUseCase",
			Message: "Falha ao salvar serviço",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao salvar serviço",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "AddServiceUseCase",
		Message: logging.DEFAULTMESSAGES.END,
	})

	return addedService, nil
}
