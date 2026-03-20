package factories

import (
	"pet-services-api/internal/logging"
	"pet-services-api/internal/mail"
	"pet-services-api/internal/repository_impl"
	"pet-services-api/internal/usecases"

	"gorm.io/gorm"
)

type ReviewFactory struct {
	CreateReview *usecases.CreateReviewUseCase
	ListReviews  *usecases.ListReviewsUseCase
}

func NewReviewFactory(db *gorm.DB, mailService mail.EmailService, logger logging.LoggerInterface) *ReviewFactory {
	userRepo := repository_impl.NewUserRepository(db)
	providerRepo := repository_impl.NewProviderRepository(db)
	requestRepo := repository_impl.NewRequestRepository(db)
	reviewRepo := repository_impl.NewReviewRepository(db)

	return &ReviewFactory{
		CreateReview: usecases.NewCreateReviewUseCase(userRepo, providerRepo, requestRepo, reviewRepo, mailService, logger),
		ListReviews:  usecases.NewListReviewsUseCase(reviewRepo, logger),
	}
}
