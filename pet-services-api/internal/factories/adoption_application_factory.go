package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/usecases"

	"gorm.io/gorm"
)

type AdoptionApplicationFactory struct {
	CreateAdoptionApplication         *usecases.CreateAdoptionApplicationUseCase
	ListMyAdoptionApplications        *usecases.ListMyAdoptionApplicationsUseCase
	ListAdoptionApplicationsByListing *usecases.ListAdoptionApplicationsByListingUseCase
	ReviewAdoptionApplication         *usecases.ReviewAdoptionApplicationUseCase
	WithdrawAdoptionApplication       *usecases.WithdrawAdoptionApplicationUseCase
}

func NewAdoptionApplicationFactory(db *gorm.DB, logger logging.LoggerInterface) *AdoptionApplicationFactory {
	listingRepo := repository_impl.NewAdoptionListingRepository(db)
	applicationRepo := repository_impl.NewAdoptionApplicationRepository(db)

	return &AdoptionApplicationFactory{
		CreateAdoptionApplication: usecases.NewCreateAdoptionApplicationUseCase(
			listingRepo,
			applicationRepo,
			logger,
		),
		ListMyAdoptionApplications: usecases.NewListMyAdoptionApplicationsUseCase(
			applicationRepo,
			logger,
		),
		ListAdoptionApplicationsByListing: usecases.NewListAdoptionApplicationsByListingUseCase(
			listingRepo,
			applicationRepo,
			logger,
		),
		ReviewAdoptionApplication: usecases.NewReviewAdoptionApplicationUseCase(
			applicationRepo,
			logger,
		),
		WithdrawAdoptionApplication: usecases.NewWithdrawAdoptionApplicationUseCase(
			applicationRepo,
			logger,
		),
	}
}
