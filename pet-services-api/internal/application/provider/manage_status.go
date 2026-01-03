package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/exceptions"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// ActivateProviderUseCase ativa o perfil do prestador.
type ActivateProviderUseCase struct {
	providerRepo provider.Repository
	logger       logging.LoggerService
}

func NewActivateProviderUseCase(providerRepo provider.Repository, logger logging.LoggerService) *ActivateProviderUseCase {
	return &ActivateProviderUseCase{providerRepo: providerRepo, logger: logging.NewSlogLogger()}
}

type ActivateProviderInput struct {
	ProviderID uuid.UUID
}

func (uc *ActivateProviderUseCase) Execute(ctx context.Context, input ActivateProviderInput) (*provider.Provider, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "ActivateProviderUseCase",
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.ProviderID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "ActivateProviderUseCase",
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

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    "ActivateProviderUseCase",
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

	if p.IsActive {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.INFO,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC200_CODE,
			From:    "ActivateProviderUseCase",
			Message: "Prestador já está ativo",
		})
		return p, nil
	}

	p.Activate()

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    "ActivateProviderUseCase",
			Message: "Falha ao ativar prestador",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao ativar prestador",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "ActivateProviderUseCase",
		Message: logging.DEFAULTMESSAGES.END,
	})

	return p, nil
}

// DeactivateProviderUseCase desativa o perfil do prestador.
type DeactivateProviderUseCase struct {
	providerRepo provider.Repository
	logger       logging.LoggerService
}

func NewDeactivateProviderUseCase(providerRepo provider.Repository, logger logging.LoggerService) *DeactivateProviderUseCase {
	return &DeactivateProviderUseCase{providerRepo: providerRepo, logger: logging.NewSlogLogger()}
}

type DeactivateProviderInput struct {
	ProviderID uuid.UUID
}

func (uc *DeactivateProviderUseCase) Execute(ctx context.Context, input DeactivateProviderInput) (*provider.Provider, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "DeactivateProviderUseCase",
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.ProviderID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "DeactivateProviderUseCase",
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

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    "DeactivateProviderUseCase",
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

	if !p.IsActive {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.INFO,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC200_CODE,
			From:    "DeactivateProviderUseCase",
			Message: "Prestador já está inativo",
		})
		return p, nil
	}

	p.Deactivate()

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    "DeactivateProviderUseCase",
			Message: "Falha ao desativar prestador",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao desativar prestador",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "DeactivateProviderUseCase",
		Message: logging.DEFAULTMESSAGES.END,
	})

	return p, nil
}
