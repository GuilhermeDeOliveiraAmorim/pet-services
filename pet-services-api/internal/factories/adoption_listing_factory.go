package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/mail"
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/usecases"

	"gorm.io/gorm"
)

type AdoptionListingFactory struct {
	CreateAdoptionListing        *usecases.CreateAdoptionListingUseCase
	UpdateAdoptionListing        *usecases.UpdateAdoptionListingUseCase
	ChangeAdoptionListingStatus  *usecases.ChangeAdoptionListingStatusUseCase
	ListPublicAdoptionListings   *usecases.ListPublicAdoptionListingsUseCase
	GetPublicAdoptionListing     *usecases.GetPublicAdoptionListingUseCase
	ListMyAdoptionListings       *usecases.ListMyAdoptionListingsUseCase
	MarkAdoptionListingAsAdopted *usecases.MarkAdoptionListingAsAdoptedUseCase
}

func NewAdoptionListingFactory(db *gorm.DB, mailService mail.EmailService, logger logging.LoggerInterface) *AdoptionListingFactory {
	petRepo := repository_impl.NewPetRepository(db)
	listingRepo := repository_impl.NewAdoptionListingRepository(db)
	applicationRepo := repository_impl.NewAdoptionApplicationRepository(db)
	guardianRepo := repository_impl.NewAdoptionGuardianProfileRepository(db)
	userRepo := repository_impl.NewUserRepository(db)

	return &AdoptionListingFactory{
		CreateAdoptionListing:        usecases.NewCreateAdoptionListingUseCase(petRepo, listingRepo, logger),
		UpdateAdoptionListing:        usecases.NewUpdateAdoptionListingUseCase(listingRepo, logger),
		ChangeAdoptionListingStatus:  usecases.NewChangeAdoptionListingStatusUseCase(listingRepo, logger),
		ListPublicAdoptionListings:   usecases.NewListPublicAdoptionListingsUseCase(listingRepo, logger),
		GetPublicAdoptionListing:     usecases.NewGetPublicAdoptionListingUseCase(listingRepo, logger),
		ListMyAdoptionListings:       usecases.NewListMyAdoptionListingsUseCase(listingRepo, logger),
		MarkAdoptionListingAsAdopted: usecases.NewMarkAdoptionListingAsAdoptedUseCase(listingRepo, applicationRepo, guardianRepo, userRepo, mailService, logger),
	}
}
