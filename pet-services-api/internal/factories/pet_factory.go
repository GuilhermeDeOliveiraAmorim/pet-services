package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/storage"
	"pet-services-api/internal/usecases"

	"gorm.io/gorm"
)

type PetFactory struct {
	AddPet         *usecases.AddPetUseCase
	AddPetPhoto    *usecases.AddPetPhotoUseCase
	DeletePetPhoto *usecases.DeletePetPhotoUseCase
	ListPets       *usecases.ListPetsUseCase
	GetPet         *usecases.GetPetUseCase
}

func NewPetFactory(db *gorm.DB, storageService storage.ObjectStorage, logger logging.LoggerInterface) *PetFactory {
	userRepo := repository_impl.NewUserRepository(db)
	specieRepo := repository_impl.NewSpecieRepository(db)
	petRepo := repository_impl.NewPetRepository(db)
	photoRepo := repository_impl.NewPhotoRepository(db)

	return &PetFactory{
		AddPet:         usecases.NewAddPetUseCase(userRepo, specieRepo, petRepo, logger),
		AddPetPhoto:    usecases.NewAddPetPhotoUseCase(userRepo, petRepo, photoRepo, storageService, logger),
		DeletePetPhoto: usecases.NewDeletePetPhotoUseCase(userRepo, petRepo, photoRepo, storageService, logger),
		ListPets:       usecases.NewListPetsUseCase(userRepo, petRepo, storageService, logger),
		GetPet:         usecases.NewGetPetUseCase(userRepo, petRepo, storageService, logger),
	}
}
