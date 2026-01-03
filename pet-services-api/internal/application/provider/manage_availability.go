package provider

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/exceptions"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// UpdateWorkingHoursUseCase ajusta a disponibilidade semanal do prestador.
type UpdateWorkingHoursUseCase struct {
	providerRepo provider.Repository
	logger       logging.LoggerService
}

// NewUpdateWorkingHoursUseCase cria uma nova instância do caso de uso.
func NewUpdateWorkingHoursUseCase(providerRepo provider.Repository, logger *slog.Logger) *UpdateWorkingHoursUseCase {
	return &UpdateWorkingHoursUseCase{providerRepo: providerRepo, logger: logging.NewSlogLogger()}
}

// UpdateWorkingHoursInput representa o horário semanal completo.
type UpdateWorkingHoursInput struct {
	ProviderID   uuid.UUID
	WorkingHours provider.WorkingHours
}

// Execute valida e persiste o horário semanal seguindo padrão CreateRequestUseCase.
func (uc *UpdateWorkingHoursUseCase) Execute(ctx context.Context, input UpdateWorkingHoursInput) (*provider.Provider, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "UpdateWorkingHoursUseCase",
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.ProviderID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "UpdateWorkingHoursUseCase",
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
			From:    "UpdateWorkingHoursUseCase",
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

	if err := p.SetWorkingHours(input.WorkingHours); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "UpdateWorkingHoursUseCase",
			Message: "Horário de trabalho inválido",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Horário de trabalho inválido",
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
			From:    "UpdateWorkingHoursUseCase",
			Message: "Falha ao atualizar disponibilidade",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao atualizar disponibilidade",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "UpdateWorkingHoursUseCase",
		Message: logging.DEFAULTMESSAGES.END,
	})

	return p, nil
}

// UpdateDayScheduleUseCase ajusta a disponibilidade de um único dia.
type UpdateDayScheduleUseCase struct {
	providerRepo provider.Repository
	logger       logging.LoggerService
}

// NewUpdateDayScheduleUseCase cria uma nova instância do caso de uso.
func NewUpdateDayScheduleUseCase(providerRepo provider.Repository, logger *slog.Logger) *UpdateDayScheduleUseCase {
	return &UpdateDayScheduleUseCase{providerRepo: providerRepo, logger: logging.NewSlogLogger()}
}

// UpdateDayScheduleInput representa o horário de um dia específico.
type UpdateDayScheduleInput struct {
	ProviderID uuid.UUID
	Day        time.Weekday
	IsOpen     bool
	Open       string
	Close      string
}

// Execute valida e persiste a disponibilidade do dia seguindo padrão CreateRequestUseCase.
func (uc *UpdateDayScheduleUseCase) Execute(ctx context.Context, input UpdateDayScheduleInput) (*provider.Provider, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "UpdateDayScheduleUseCase",
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.ProviderID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "UpdateDayScheduleUseCase",
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
			From:    "UpdateDayScheduleUseCase",
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

	schedule := provider.DaySchedule{
		IsOpen: input.IsOpen,
		Open:   input.Open,
		Close:  input.Close,
	}

	if err := p.SetDaySchedule(input.Day, schedule); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "UpdateDayScheduleUseCase",
			Message: "Horário do dia inválido",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Horário do dia inválido",
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
			From:    "UpdateDayScheduleUseCase",
			Message: "Falha ao atualizar disponibilidade",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao atualizar disponibilidade",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "UpdateDayScheduleUseCase",
		Message: logging.DEFAULTMESSAGES.END,
	})

	return p, nil
}
