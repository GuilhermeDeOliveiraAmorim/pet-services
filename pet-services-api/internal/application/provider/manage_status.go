package provider

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// ActivateProviderUseCase ativa o perfil do prestador.
type ActivateProviderUseCase struct {
	providerRepo provider.Repository
	logger       *slog.Logger
}

func NewActivateProviderUseCase(providerRepo provider.Repository, logger *slog.Logger) *ActivateProviderUseCase {
	return &ActivateProviderUseCase{providerRepo: providerRepo, logger: logging.EnsureLogger(logger)}
}

type ActivateProviderInput struct {
	ProviderID uuid.UUID
}

func (uc *ActivateProviderUseCase) Execute(ctx context.Context, input ActivateProviderInput) error {
	var err error
	defer logging.UseCase(ctx, uc.logger, "ActivateProviderUseCase", slog.String("provider_id", input.ProviderID.String()))(&err)

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		err = provider.ErrProviderNotFound
		return err
	}

	if p.IsActive {
		return nil
	}

	p.Activate()

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		err = fmt.Errorf("falha ao ativar prestador: %w", err)
		return err
	}

	return nil
}

// DeactivateProviderUseCase desativa o perfil do prestador.
type DeactivateProviderUseCase struct {
	providerRepo provider.Repository
	logger       *slog.Logger
}

func NewDeactivateProviderUseCase(providerRepo provider.Repository, logger *slog.Logger) *DeactivateProviderUseCase {
	return &DeactivateProviderUseCase{providerRepo: providerRepo, logger: logging.EnsureLogger(logger)}
}

type DeactivateProviderInput struct {
	ProviderID uuid.UUID
}

func (uc *DeactivateProviderUseCase) Execute(ctx context.Context, input DeactivateProviderInput) error {
	var err error
	defer logging.UseCase(ctx, uc.logger, "DeactivateProviderUseCase", slog.String("provider_id", input.ProviderID.String()))(&err)

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		err = provider.ErrProviderNotFound
		return err
	}

	if !p.IsActive {
		return nil
	}

	p.Deactivate()

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		err = fmt.Errorf("falha ao desativar prestador: %w", err)
		return err
	}

	return nil
}
