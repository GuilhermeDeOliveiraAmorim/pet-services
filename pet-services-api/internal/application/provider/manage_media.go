package provider

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/domain/provider"
)

// AddPhotoUseCase adiciona uma foto ao perfil do prestador.
type AddPhotoUseCase struct {
	providerRepo provider.Repository
}

func NewAddPhotoUseCase(providerRepo provider.Repository) *AddPhotoUseCase {
	return &AddPhotoUseCase{providerRepo: providerRepo}
}

type AddPhotoInput struct {
	ProviderID uuid.UUID
	URL        string
}

func (uc *AddPhotoUseCase) Execute(ctx context.Context, input AddPhotoInput) error {
	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		return provider.ErrProviderNotFound
	}

	if err := p.AddPhoto(input.URL); err != nil {
		return err
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("falha ao adicionar foto: %w", err)
	}

	return nil
}

// RemovePhotoUseCase remove uma foto do perfil do prestador.
type RemovePhotoUseCase struct {
	providerRepo provider.Repository
}

func NewRemovePhotoUseCase(providerRepo provider.Repository) *RemovePhotoUseCase {
	return &RemovePhotoUseCase{providerRepo: providerRepo}
}

type RemovePhotoInput struct {
	ProviderID uuid.UUID
	PhotoID    uuid.UUID
}

func (uc *RemovePhotoUseCase) Execute(ctx context.Context, input RemovePhotoInput) error {
	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		return provider.ErrProviderNotFound
	}

	if err := p.RemovePhoto(input.PhotoID); err != nil {
		return err
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("falha ao remover foto: %w", err)
	}

	return nil
}

// ReorderPhotosUseCase reordena as fotos do prestador.
type ReorderPhotosUseCase struct {
	providerRepo provider.Repository
}

func NewReorderPhotosUseCase(providerRepo provider.Repository) *ReorderPhotosUseCase {
	return &ReorderPhotosUseCase{providerRepo: providerRepo}
}

type ReorderPhotosInput struct {
	ProviderID uuid.UUID
	Order      []uuid.UUID
}

func (uc *ReorderPhotosUseCase) Execute(ctx context.Context, input ReorderPhotosInput) error {
	p, err := uc.providerRepo.FindByID(ctx, input.ProviderID)
	if err != nil {
		return provider.ErrProviderNotFound
	}

	if err := p.ReorderPhotos(input.Order); err != nil {
		return err
	}

	if err := uc.providerRepo.Update(ctx, p); err != nil {
		return fmt.Errorf("falha ao reordenar fotos: %w", err)
	}

	return nil
}
