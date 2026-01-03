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

// ApproveProviderUseCase aprova um prestador (ativa o perfil).
type ApproveProviderUseCase struct {
	providerRepo provider.Repository
	logger       logging.LoggerService
}

// NewApproveProviderUseCase cria instância.
func NewApproveProviderUseCase(providerRepo provider.Repository, logger *slog.Logger) *ApproveProviderUseCase {
	return &ApproveProviderUseCase{providerRepo: providerRepo, logger: logging.NewSlogLogger()}
}

// ApproveProviderInput dados
type ApproveProviderInput struct {
	ProviderID uuid.UUID
	Note       string // opcional
}

// Execute aprova e persiste.
func (uc *ApproveProviderUseCase) Execute(ctx context.Context, input ApproveProviderInput) (*provider.Provider, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "ApproveProviderUseCase",
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.ProviderID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "ApproveProviderUseCase",
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
			From:    "ApproveProviderUseCase",
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

	if err := p.Approve(input.Note); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "ApproveProviderUseCase",
			Message: "Erro ao aprovar prestador",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Erro ao aprovar prestador",
			Status: exceptions.RFC400_CODE,
			Detail: err.Error(),
		}}
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    "ApproveProviderUseCase",
			Message: "Falha ao aprovar prestador",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao aprovar prestador",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "ApproveProviderUseCase",
		Message: logging.DEFAULTMESSAGES.END,
	})

	return p, nil
}

// RejectProviderUseCase reprova/bloqueia um prestador.
type RejectProviderUseCase struct {
	providerRepo provider.Repository
	logger       logging.LoggerService
}

// NewRejectProviderUseCase cria instância.
func NewRejectProviderUseCase(providerRepo provider.Repository, logger *slog.Logger) *RejectProviderUseCase {
	return &RejectProviderUseCase{providerRepo: providerRepo, logger: logging.NewSlogLogger()}
}

// RejectProviderInput dados para reprovação.
type RejectProviderInput struct {
	ProviderID uuid.UUID
	Reason     string // obrigatório
}

// Execute rejeita e persiste.
func (uc *RejectProviderUseCase) Execute(ctx context.Context, input RejectProviderInput) (*provider.Provider, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "RejectProviderUseCase",
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.ProviderID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "RejectProviderUseCase",
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

	if input.Reason == "" {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "RejectProviderUseCase",
			Message: "Motivo é obrigatório",
			Error:   fmt.Errorf("motivo é obrigatório"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Motivo é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O motivo da reprovação é obrigatório.",
		}}
	}

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    "RejectProviderUseCase",
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

	if err := p.Reject(input.Reason); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "RejectProviderUseCase",
			Message: "Erro ao rejeitar prestador",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Erro ao rejeitar prestador",
			Status: exceptions.RFC400_CODE,
			Detail: err.Error(),
		}}
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    "RejectProviderUseCase",
			Message: "Falha ao rejeitar prestador",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao rejeitar prestador",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "RejectProviderUseCase",
		Message: logging.DEFAULTMESSAGES.END,
	})

	return p, nil
}
