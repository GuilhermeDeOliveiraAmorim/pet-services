package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// UpdateWorkingHoursUseCase ajusta a disponibilidade semanal do prestador.
type UpdateWorkingHoursUseCase struct {
	providerRepo provider.Repository
}

// NewUpdateWorkingHoursUseCase cria uma nova instância do caso de uso.
func NewUpdateWorkingHoursUseCase(providerRepo provider.Repository) *UpdateWorkingHoursUseCase {
	return &UpdateWorkingHoursUseCase{providerRepo: providerRepo}
}

// UpdateWorkingHoursInput representa o horário semanal completo.
type UpdateWorkingHoursInput struct {
	ProviderID    uuid.UUID
	WorkingHours provider.WorkingHours
}

// Execute valida e persiste o horário semanal.
func (uc *UpdateWorkingHoursUseCase) Execute(ctx context.Context, input UpdateWorkingHoursInput) error {
	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		return provider.ErrProviderNotFound
	}

	if err := p.SetWorkingHours(input.WorkingHours); err != nil {
		return err
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("falha ao atualizar disponibilidade: %w", err)
	}

	return nil
}

// UpdateDayScheduleUseCase ajusta a disponibilidade de um único dia.
type UpdateDayScheduleUseCase struct {
	providerRepo provider.Repository
}

// NewUpdateDayScheduleUseCase cria uma nova instância do caso de uso.
func NewUpdateDayScheduleUseCase(providerRepo provider.Repository) *UpdateDayScheduleUseCase {
	return &UpdateDayScheduleUseCase{providerRepo: providerRepo}
}

// UpdateDayScheduleInput representa o horário de um dia específico.
type UpdateDayScheduleInput struct {
	ProviderID uuid.UUID
	Day        time.Weekday
	IsOpen     bool
	Open       string
	Close      string
}

// Execute valida e persiste a disponibilidade do dia.
func (uc *UpdateDayScheduleUseCase) Execute(ctx context.Context, input UpdateDayScheduleInput) error {
	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		return provider.ErrProviderNotFound
	}

	schedule := provider.DaySchedule{
		IsOpen: input.IsOpen,
		Open:   input.Open,
		Close:  input.Close,
	}

	if err := p.SetDaySchedule(input.Day, schedule); err != nil {
		return err
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("falha ao atualizar disponibilidade: %w", err)
	}

	return nil
}
