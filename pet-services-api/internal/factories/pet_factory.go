package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/usecases"

	"gorm.io/gorm"
)

type PetFactory struct {
	AddPet *usecases.AddPetUseCase
}

func NewPetFactory(db *gorm.DB, logger logging.LoggerInterface) *PetFactory {
	userRepo := repository_impl.NewUserRepository(db)
	specieRepo := repository_impl.NewSpecieRepository(db)
	petRepo := repository_impl.NewPetRepository(db)

	return &PetFactory{
		AddPet: usecases.NewAddPetUseCase(userRepo, specieRepo, petRepo, logger),
	}
}
