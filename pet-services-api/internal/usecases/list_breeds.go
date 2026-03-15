package usecases

import (
	"context"
	"errors"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type ListBreedsInput struct {
	SpeciesID string `json:"species_id"`
}

type ListBreedsOutput struct {
	Breeds []*entities.Breed `json:"breeds"`
}

type ListBreedsUseCase struct {
	breedRepository entities.BreedRepository
	logger          logging.LoggerInterface
}

func NewListBreedsUseCase(breedRepository entities.BreedRepository, logger logging.LoggerInterface) *ListBreedsUseCase {
	return &ListBreedsUseCase{breedRepository: breedRepository, logger: logger}
}

func (uc *ListBreedsUseCase) Execute(ctx context.Context, input ListBreedsInput) (*ListBreedsOutput, []exceptions.ProblemDetails) {
	const from = "ListBreedsUseCase.Execute"

	if input.SpeciesID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID da espécie ausente", errors.New("O ID da espécie é obrigatório"))
	}

	breeds, err := uc.breedRepository.ListBySpecies(input.SpeciesID)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao listar raças", err)
	}

	uc.logger.LogInfo(ctx, from, "Listagem de raças executada")
	return &ListBreedsOutput{Breeds: breeds}, nil
}
