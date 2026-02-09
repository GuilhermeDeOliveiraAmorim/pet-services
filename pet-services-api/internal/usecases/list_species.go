package usecases

import (
	"context"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type ListSpeciesOutput struct {
	Species []*entities.Species `json:"species"`
}

type ListSpeciesUseCase struct {
	specieRepository entities.SpecieRepository
	logger           logging.LoggerInterface
}

func NewListSpeciesUseCase(specieRepository entities.SpecieRepository, logger logging.LoggerInterface) *ListSpeciesUseCase {
	return &ListSpeciesUseCase{specieRepository: specieRepository, logger: logger}
}

func (uc *ListSpeciesUseCase) Execute(ctx context.Context) (*ListSpeciesOutput, []exceptions.ProblemDetails) {
	species, err := uc.specieRepository.List()
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, "ListSpeciesUseCase.Execute", "Erro ao listar espécies", err)
	}

	uc.logger.LogInfo(ctx, "ListSpeciesUseCase.Execute", "Listagem de espécies executada")
	return &ListSpeciesOutput{Species: species}, nil
}
