package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// ApproveProviderUseCase aprova um prestador (ativa o perfil).
type ApproveProviderUseCase struct {
	providerRepo provider.Repository
}

// NewApproveProviderUseCase cria instância.
func NewApproveProviderUseCase(providerRepo provider.Repository) *ApproveProviderUseCase {
	return &ApproveProviderUseCase{providerRepo: providerRepo}
}

// ApproveProviderInput dados para aprovação.
type ApproveProviderInput struct {
	ProviderID uuid.UUID
	Note       string // opcional
}

// Execute aprova e persiste.
func (uc *ApproveProviderUseCase) Execute(ctx context.Context, input ApproveProviderInput) error {
	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		return provider.ErrProviderNotFound
	}

	if err := p.Approve(input.Note); err != nil {
		return err
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("falha ao aprovar prestador: %w", err)
	}

	return nil
}

// RejectProviderUseCase reprova/bloqueia um prestador.
type RejectProviderUseCase struct {
	providerRepo provider.Repository
}

// NewRejectProviderUseCase cria instância.
func NewRejectProviderUseCase(providerRepo provider.Repository) *RejectProviderUseCase {
	return &RejectProviderUseCase{providerRepo: providerRepo}
}

// RejectProviderInput dados para reprovação.
type RejectProviderInput struct {
	ProviderID uuid.UUID
	Reason     string // obrigatório
}

// Execute rejeita e persiste.
func (uc *RejectProviderUseCase) Execute(ctx context.Context, input RejectProviderInput) error {
	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		return provider.ErrProviderNotFound
	}

	if err := p.Reject(input.Reason); err != nil {
		return err
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("falha ao rejeitar prestador: %w", err)
	}

	return nil
}
