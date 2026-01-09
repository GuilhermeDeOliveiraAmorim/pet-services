package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"pet-services-api/internal/application/exceptions"
	"pet-services-api/internal/application/logging"
	"pet-services-api/internal/domain/provider"
)

// UpdateServiceUseCase atualiza um serviço existente do prestador.
type UpdateServiceUseCase struct {
	providerRepo provider.Repository
	logger       logging.LoggerService
}

// NewUpdateServiceUseCase cria nova instância do caso de uso.
func NewUpdateServiceUseCase(providerRepo provider.Repository, logger logging.LoggerService) *UpdateServiceUseCase {
	return &UpdateServiceUseCase{providerRepo: providerRepo, logger: logging.NewSlogLogger()}
}

// UpdateServiceInput representa os dados para atualizar serviço.
type UpdateServiceInput struct {
	ProviderID uuid.UUID
	Category   string
	Name       string
	PriceMin   float64
	PriceMax   float64
}

// Execute valida e atualiza o serviço seguindo padrão CreateRequestUseCase.
func (uc *UpdateServiceUseCase) Execute(ctx context.Context, input UpdateServiceInput) (*provider.Service, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "UpdateServiceUseCase",
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.ProviderID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "UpdateServiceUseCase",
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
			From:    "UpdateServiceUseCase",
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
			From:    "UpdateServiceUseCase",
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
			From:    "UpdateServiceUseCase",
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
			From:    "UpdateServiceUseCase",
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

	var updatedService *provider.Service
	updated := false
	for i, svc := range p.Services {
		if svc.Category == input.Category && svc.Name == input.Name {
			p.Services[i].PriceMin = input.PriceMin
			p.Services[i].PriceMax = input.PriceMax
			updatedService = &p.Services[i]
			updated = true
			break
		}
	}

	if !updated {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    "UpdateServiceUseCase",
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
			From:    "UpdateServiceUseCase",
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
		From:    "UpdateServiceUseCase",
		Message: logging.DEFAULTMESSAGES.END,
	})

	return updatedService, nil
}
