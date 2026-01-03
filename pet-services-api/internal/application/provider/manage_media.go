package provider

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// AddPhotoUseCase adiciona uma foto ao perfil do prestador.
type AddPhotoUseCase struct {
	providerRepo provider.Repository
	logger       *slog.Logger
}

func NewAddPhotoUseCase(providerRepo provider.Repository, logger *slog.Logger) *AddPhotoUseCase {
	return &AddPhotoUseCase{providerRepo: providerRepo, logger: logging.EnsureLogger(logger)}
}

type AddPhotoInput struct {
	ProviderID uuid.UUID
	URL        string
}

func (uc *AddPhotoUseCase) Execute(ctx context.Context, input AddPhotoInput) error {
	var err error
	defer logging.UseCase(ctx, uc.logger, "AddPhotoUseCase", slog.String("provider_id", input.ProviderID.String()))(&err)

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		err = provider.ErrProviderNotFound
		return err
	}

	if err := p.AddPhoto(input.URL); err != nil {
		return err
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		err = fmt.Errorf("falha ao adicionar foto: %w", err)
		return err
	}

	return nil
}

// RemovePhotoUseCase remove uma foto do perfil do prestador.
type RemovePhotoUseCase struct {
	providerRepo provider.Repository
	logger       *slog.Logger
}

func NewRemovePhotoUseCase(providerRepo provider.Repository, logger *slog.Logger) *RemovePhotoUseCase {
	return &RemovePhotoUseCase{providerRepo: providerRepo, logger: logging.EnsureLogger(logger)}
}

type RemovePhotoInput struct {
	ProviderID uuid.UUID
	PhotoID    uuid.UUID
}

func (uc *RemovePhotoUseCase) Execute(ctx context.Context, input RemovePhotoInput) error {
	var err error
	defer logging.UseCase(ctx, uc.logger, "RemovePhotoUseCase", slog.String("provider_id", input.ProviderID.String()), slog.String("photo_id", input.PhotoID.String()))(&err)

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		err = provider.ErrProviderNotFound
		return err
	}

	if err := p.RemovePhoto(input.PhotoID); err != nil {
		return err
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		err = fmt.Errorf("falha ao remover foto: %w", err)
		return err
	}

	return nil
}

// ReorderPhotosUseCase reordena as fotos do prestador.
type ReorderPhotosUseCase struct {
	providerRepo provider.Repository
	logger       *slog.Logger
}

func NewReorderPhotosUseCase(providerRepo provider.Repository, logger *slog.Logger) *ReorderPhotosUseCase {
	return &ReorderPhotosUseCase{providerRepo: providerRepo, logger: logging.EnsureLogger(logger)}
}

type ReorderPhotosInput struct {
	ProviderID uuid.UUID
	Order      []uuid.UUID
}

func (uc *ReorderPhotosUseCase) Execute(ctx context.Context, input ReorderPhotosInput) error {
	var err error
	defer logging.UseCase(ctx, uc.logger, "ReorderPhotosUseCase", slog.String("provider_id", input.ProviderID.String()))(&err)

	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		err = provider.ErrProviderNotFound
		return err
	}

	if err := p.ReorderPhotos(input.Order); err != nil {
		return err
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		err = fmt.Errorf("falha ao reordenar fotos: %w", err)
		return err
	}

	return nil
}
