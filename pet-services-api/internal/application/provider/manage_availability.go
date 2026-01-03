package provider

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// UpdateWorkingHoursUseCase ajusta a disponibilidade semanal do prestador.
type UpdateWorkingHoursUseCase struct {
	providerRepo provider.Repository
	logger       *slog.Logger
}

// NewUpdateWorkingHoursUseCase cria uma nova instância do caso de uso.
func NewUpdateWorkingHoursUseCase(providerRepo provider.Repository, logger *slog.Logger) *UpdateWorkingHoursUseCase {
	return &UpdateWorkingHoursUseCase{providerRepo: providerRepo, logger: logging.EnsureLogger(logger)}
}

// UpdateWorkingHoursInput representa o horário semanal completo.
type UpdateWorkingHoursInput struct {
	ProviderID   uuid.UUID
	WorkingHours provider.WorkingHours
}

// Execute valida e persiste o horário semanal.
func (uc *UpdateWorkingHoursUseCase) Execute(ctx context.Context, input UpdateWorkingHoursInput) error {
	var err error
	defer logging.UseCase(ctx, uc.logger, "UpdateWorkingHoursUseCase", slog.String("provider_id", input.ProviderID.String()))(&err)

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		err = provider.ErrProviderNotFound
		return err
	}

	if err := p.SetWorkingHours(input.WorkingHours); err != nil {
		return err
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		err = fmt.Errorf("falha ao atualizar disponibilidade: %w", err)
		return err
	}

	return nil
}

// UpdateDayScheduleUseCase ajusta a disponibilidade de um único dia.
type UpdateDayScheduleUseCase struct {
	providerRepo provider.Repository
	logger       *slog.Logger
}

// NewUpdateDayScheduleUseCase cria uma nova instância do caso de uso.
func NewUpdateDayScheduleUseCase(providerRepo provider.Repository, logger *slog.Logger) *UpdateDayScheduleUseCase {
	return &UpdateDayScheduleUseCase{providerRepo: providerRepo, logger: logging.EnsureLogger(logger)}
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
	var err error
	defer logging.UseCase(ctx, uc.logger, "UpdateDayScheduleUseCase", slog.String("provider_id", input.ProviderID.String()), slog.Int("day", int(input.Day)))(&err)

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		err = provider.ErrProviderNotFound
		return err
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
		err = fmt.Errorf("falha ao atualizar disponibilidade: %w", err)
		return err
	}

	return nil
}
