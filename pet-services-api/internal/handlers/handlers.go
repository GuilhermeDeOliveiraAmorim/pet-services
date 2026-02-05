package handlers

import (
	"context"
	"pet-services-api/internal/database"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"

	"github.com/oklog/ulid/v2"
)

type HandlerFactory struct {
	UserHandler   *UserHandler
	TokenHandler  *TokenHandler
	HealthHandler *HealthHandler
	ReferenceHandler *ReferenceHandler
	Logger        logging.LoggerInterface
}

func NewHandlerFactory(inputFactory database.StorageInput, logger logging.LoggerInterface) *HandlerFactory {
	userFactory := factories.NewUserFactory(inputFactory.DB, logger)
	tokenFactory := factories.NewTokenFactory(inputFactory.DB, nil, 0, logger)
	healthFactory := factories.NewHealthFactory(inputFactory.DB, logger)
	referenceFactory := factories.NewReferenceFactory(logger)

	return &HandlerFactory{
		UserHandler:   NewUserHandler(userFactory, logger),
		TokenHandler:  NewTokenHandler(tokenFactory, logger),
		HealthHandler: NewHealthHandler(healthFactory, logger),
		ReferenceHandler: NewReferenceHandler(referenceFactory, logger),
		Logger:        logger,
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
