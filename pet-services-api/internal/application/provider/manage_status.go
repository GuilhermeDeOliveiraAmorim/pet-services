package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// ActivateProviderUseCase ativa o perfil do prestador.
type ActivateProviderUseCase struct {
	providerRepo provider.Repository
}

func NewActivateProviderUseCase(providerRepo provider.Repository) *ActivateProviderUseCase {
	return &ActivateProviderUseCase{providerRepo: providerRepo}
}

type ActivateProviderInput struct {
	ProviderID uuid.UUID
}

func (uc *ActivateProviderUseCase) Execute(ctx context.Context, input ActivateProviderInput) error {
	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		return provider.ErrProviderNotFound
	}

	if p.IsActive {
		return nil
	}

	p.Activate()

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("falha ao ativar prestador: %w", err)
	}

	return nil
}

// DeactivateProviderUseCase desativa o perfil do prestador.
type DeactivateProviderUseCase struct {
	providerRepo provider.Repository
}

func NewDeactivateProviderUseCase(providerRepo provider.Repository) *DeactivateProviderUseCase {
	return &DeactivateProviderUseCase{providerRepo: providerRepo}
}

type DeactivateProviderInput struct {
	ProviderID uuid.UUID
}

func (uc *DeactivateProviderUseCase) Execute(ctx context.Context, input DeactivateProviderInput) error {
	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		return provider.ErrProviderNotFound
	}

	if !p.IsActive {
		return nil
	}

	p.Deactivate()

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("falha ao desativar prestador: %w", err)
	}

	return nil
}
