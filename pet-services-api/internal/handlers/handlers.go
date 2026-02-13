package handlers

import (
	"context"
	"pet-services-api/internal/database"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/mail"
	"pet-services-api/internal/storage"

	"github.com/oklog/ulid/v2"
)

type HandlerFactory struct {
	UserHandler      *UserHandler
	TokenHandler     *TokenHandler
	HealthHandler    *HealthHandler
	ReferenceHandler *ReferenceHandler
	SpecieHandler    *SpecieHandler
	PetHandler       *PetHandler
	ProviderHandler  *ProviderHandler
	ServiceHandler   *ServiceHandler
	RequestHandler   *RequestHandler
	CategoryHandler  *CategoryHandler
	Logger           logging.LoggerInterface
}

func NewHandlerFactory(inputFactory database.StorageInput, logger logging.LoggerInterface) *HandlerFactory {
	storageService, err := storage.NewMinioServiceFromEnv()
	if err != nil {
		logger.LogError(context.Background(), "HandlerFactory", "Falha ao configurar MinIO", err)
	}
	userFactory := factories.NewUserFactory(inputFactory.DB, storageService, logger)
	mailService := mail.GetEmailServiceFromEnv()
	tokenFactory := factories.NewTokenFactory(inputFactory.DB, mailService, 0, storageService, logger)
	healthFactory := factories.NewHealthFactory(inputFactory.DB, logger)
	referenceFactory := factories.NewReferenceFactory(logger)
	specieFactory := factories.NewSpecieFactory(inputFactory.DB, logger)
	petFactory := factories.NewPetFactory(inputFactory.DB, storageService, logger)
	providerFactory := factories.NewProviderFactory(inputFactory.DB, logger)
	serviceFactory := factories.NewServiceFactory(inputFactory.DB, storageService, logger)
	requestFactory := factories.NewRequestFactory(inputFactory.DB, logger)
	categoryFactory := factories.NewCategoryFactory(inputFactory.DB, logger)

	return &HandlerFactory{
		UserHandler:      NewUserHandler(userFactory, logger),
		TokenHandler:     NewTokenHandler(tokenFactory, logger),
		HealthHandler:    NewHealthHandler(healthFactory, logger),
		ReferenceHandler: NewReferenceHandler(referenceFactory, logger),
		SpecieHandler:    NewSpecieHandler(specieFactory, logger),
		PetHandler:       NewPetHandler(petFactory, logger),
		ProviderHandler:  NewProviderHandler(providerFactory, logger),
		ServiceHandler:   NewServiceHandler(serviceFactory, logger),
		RequestHandler:   NewRequestHandler(requestFactory, logger),
		CategoryHandler:  NewCategoryHandler(categoryFactory, logger),
		Logger:           logger,
	}
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
