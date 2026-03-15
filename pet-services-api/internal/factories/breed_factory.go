package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/usecases"

	"gorm.io/gorm"
)

type BreedFactory struct {
	ListBreeds *usecases.ListBreedsUseCase
}

func NewBreedFactory(db *gorm.DB, logger logging.LoggerInterface) *BreedFactory {
	breedRepo := repository_impl.NewBreedRepository(db)

	return &BreedFactory{
		ListBreeds: usecases.NewListBreedsUseCase(breedRepo, logger),
	}
}
