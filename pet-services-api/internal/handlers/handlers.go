package handlers

import (
	"context"
	"pet-services-api/internal/database"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/mail"
	"pet-services-api/internal/storage"
	"time"

	"github.com/oklog/ulid/v2"
)

type HandlerFactory struct {
	UserHandler             *UserHandler
	TokenHandler            *TokenHandler
	HealthHandler           *HealthHandler
	ReferenceHandler        *ReferenceHandler
	SpecieHandler           *SpecieHandler
	BreedHandler            *BreedHandler
	PetHandler              *PetHandler
	ProviderHandler         *ProviderHandler
	ServiceHandler          *ServiceHandler
	RequestHandler          *RequestHandler
	CategoryHandler         *CategoryHandler
	ReviewHandler           *ReviewHandler
	AdoptionGuardianHandler *AdoptionGuardianHandler
	AdoptionListingHandler  *AdoptionListingHandler
	Logger                  logging.LoggerInterface
}

func NewHandlerFactory(inputFactory database.StorageInput, logger logging.LoggerInterface) *HandlerFactory {
	storageService := initializeStorageService(logger)
	mailService := mail.GetEmailServiceFromEnv()

	userFactory := factories.NewUserFactory(inputFactory.DB, storageService, mailService, 24*time.Hour, logger)
	tokenFactory := factories.NewTokenFactory(inputFactory.DB, mailService, 0, storageService, logger)
	healthFactory := factories.NewHealthFactory(inputFactory.DB, logger)
	referenceFactory := factories.NewReferenceFactory(logger)
	specieFactory := factories.NewSpecieFactory(inputFactory.DB, logger)
	breedFactory := factories.NewBreedFactory(inputFactory.DB, logger)
	petFactory := factories.NewPetFactory(inputFactory.DB, storageService, logger)
	providerFactory := factories.NewProviderFactory(inputFactory.DB, storageService, logger)
	serviceFactory := factories.NewServiceFactory(inputFactory.DB, storageService, logger)
	requestFactory := factories.NewRequestFactory(inputFactory.DB, storageService, mailService, logger)
	categoryFactory := factories.NewCategoryFactory(inputFactory.DB, logger)
	reviewFactory := factories.NewReviewFactory(inputFactory.DB, mailService, logger)
	adoptionGuardianFactory := factories.NewAdoptionGuardianFactory(inputFactory.DB, logger)
	adoptionListingFactory := factories.NewAdoptionListingFactory(inputFactory.DB, logger)

	return &HandlerFactory{
		UserHandler:             NewUserHandler(userFactory, logger),
		TokenHandler:            NewTokenHandler(tokenFactory, logger),
		HealthHandler:           NewHealthHandler(healthFactory, logger),
		ReferenceHandler:        NewReferenceHandler(referenceFactory, logger),
		SpecieHandler:           NewSpecieHandler(specieFactory, logger),
		BreedHandler:            NewBreedHandler(breedFactory, logger),
		PetHandler:              NewPetHandler(petFactory, logger),
		ProviderHandler:         NewProviderHandler(providerFactory, logger),
		ServiceHandler:          NewServiceHandler(serviceFactory, logger),
		RequestHandler:          NewRequestHandler(requestFactory, logger),
		CategoryHandler:         NewCategoryHandler(categoryFactory, logger),
		ReviewHandler:           NewReviewHandler(reviewFactory, logger),
		AdoptionGuardianHandler: NewAdoptionGuardianHandler(adoptionGuardianFactory, logger),
		AdoptionListingHandler:  NewAdoptionListingHandler(adoptionListingFactory, logger),
		Logger:                  logger,
	}
}
func initializeStorageService(logger logging.LoggerInterface) storage.ObjectStorage {
	if gcsService, err := storage.NewGCSServiceFromEnv(); err == nil {
		logger.LogInfo(context.Background(), "StorageService", "Google Cloud Storage inicializado")
		return gcsService
	}

	minioService, err := storage.NewMinioServiceFromEnv()
	if err != nil {
		logger.LogError(context.Background(), "StorageService", "Falha ao configurar armazenamento", err)
		return nil
	}

	logger.LogInfo(context.Background(), "StorageService", "MinIO inicializado")
	return minioService
}
func (hf *HandlerFactory) IsValidULID(ctx context.Context, id string) bool {
	var isValid = true

	_, err := ulid.Parse(id)

	if err != nil {
		hf.Logger.LogBadRequest(ctx, "IsValidULID", "ID inválido: "+id, err)
		isValid = false
	}

	return isValid
}
