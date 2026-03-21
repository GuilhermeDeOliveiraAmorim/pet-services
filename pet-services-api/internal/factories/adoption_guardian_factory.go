package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/usecases"

	"gorm.io/gorm"
)

type AdoptionGuardianFactory struct {
	CreateAdoptionGuardianProfile  *usecases.CreateAdoptionGuardianProfileUseCase
	GetMyAdoptionGuardianProfile   *usecases.GetMyAdoptionGuardianProfileUseCase
	UpdateAdoptionGuardianProfile  *usecases.UpdateAdoptionGuardianProfileUseCase
	ApproveAdoptionGuardianProfile *usecases.ApproveAdoptionGuardianProfileUseCase
	RejectAdoptionGuardianProfile  *usecases.RejectAdoptionGuardianProfileUseCase
}

func NewAdoptionGuardianFactory(db *gorm.DB, logger logging.LoggerInterface) *AdoptionGuardianFactory {
	userRepo := repository_impl.NewUserRepository(db)
	guardianProfileRepo := repository_impl.NewAdoptionGuardianProfileRepository(db)

	return &AdoptionGuardianFactory{
		CreateAdoptionGuardianProfile:  usecases.NewCreateAdoptionGuardianProfileUseCase(userRepo, guardianProfileRepo, logger),
		GetMyAdoptionGuardianProfile:   usecases.NewGetMyAdoptionGuardianProfileUseCase(guardianProfileRepo, logger),
		UpdateAdoptionGuardianProfile:  usecases.NewUpdateAdoptionGuardianProfileUseCase(guardianProfileRepo, logger),
		ApproveAdoptionGuardianProfile: usecases.NewApproveAdoptionGuardianProfileUseCase(guardianProfileRepo, logger),
		RejectAdoptionGuardianProfile:  usecases.NewRejectAdoptionGuardianProfileUseCase(guardianProfileRepo, logger),
	}
}
