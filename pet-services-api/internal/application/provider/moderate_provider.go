package provider

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// ApproveProviderUseCase aprova um prestador (ativa o perfil).
type ApproveProviderUseCase struct {
	providerRepo provider.Repository
	logger       *slog.Logger
}

// NewApproveProviderUseCase cria instância.
func NewApproveProviderUseCase(providerRepo provider.Repository, logger *slog.Logger) *ApproveProviderUseCase {
	return &ApproveProviderUseCase{providerRepo: providerRepo, logger: logging.EnsureLogger(logger)}
}

// ApproveProviderInput dados para aprovação.
type ApproveProviderInput struct {
	ProviderID uuid.UUID
	Note       string // opcional
}

// Execute aprova e persiste.
func (uc *ApproveProviderUseCase) Execute(ctx context.Context, input ApproveProviderInput) error {
	var err error
	defer logging.UseCase(ctx, uc.logger, "ApproveProviderUseCase", slog.String("provider_id", input.ProviderID.String()))(&err)

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		err = provider.ErrProviderNotFound
		return err
	}

	if err := p.Approve(input.Note); err != nil {
		return err
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		err = fmt.Errorf("falha ao aprovar prestador: %w", err)
		return err
	}

	return nil
}

// RejectProviderUseCase reprova/bloqueia um prestador.
type RejectProviderUseCase struct {
	providerRepo provider.Repository
	logger       *slog.Logger
}

// NewRejectProviderUseCase cria instância.
func NewRejectProviderUseCase(providerRepo provider.Repository, logger *slog.Logger) *RejectProviderUseCase {
	return &RejectProviderUseCase{providerRepo: providerRepo, logger: logging.EnsureLogger(logger)}
}

// RejectProviderInput dados para reprovação.
type RejectProviderInput struct {
	ProviderID uuid.UUID
	Reason     string // obrigatório
}

// Execute rejeita e persiste.
func (uc *RejectProviderUseCase) Execute(ctx context.Context, input RejectProviderInput) error {
	var err error
	defer logging.UseCase(ctx, uc.logger, "RejectProviderUseCase", slog.String("provider_id", input.ProviderID.String()))(&err)

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		err = provider.ErrProviderNotFound
		return err
	}

	if err := p.Reject(input.Reason); err != nil {
		return err
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		err = fmt.Errorf("falha ao rejeitar prestador: %w", err)
		return err
	}

	return nil
}
