package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/usecases"

	"gorm.io/gorm"
)

type SpecieFactory struct {
	ListSpecies *usecases.ListSpeciesUseCase
}

func NewSpecieFactory(db *gorm.DB, logger logging.LoggerInterface) *SpecieFactory {
	specieRepo := repository_impl.NewSpecieRepository(db)

	return &SpecieFactory{
		ListSpecies: usecases.NewListSpeciesUseCase(specieRepo, logger),
	}
}
