package handlers

import (
	"context"
	"pet-services-api/internal/database"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"

	"github.com/oklog/ulid/v2"
)

type HandlerFactory struct {
	UserHandler  *UserHandler
	TokenHandler *TokenHandler
}

func NewHandlerFactory(inputFactory database.StorageInput) *HandlerFactory {
	userFactory := factories.NewUserFactory(inputFactory.DB)
	tokenFactory := factories.NewTokenFactory(inputFactory.DB, nil, 0)

	return &HandlerFactory{
		UserHandler:  NewUserHandler(userFactory),
		TokenHandler: NewTokenHandler(tokenFactory),
	}
}

func IsValidULID(ctx context.Context, id string) bool {
	var isValid = true

	_, err := ulid.Parse(id)

	if err != nil {
		logging.NewLogger(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.INTERFACE_HANDLERS,
			Code:    exceptions.RFC500_CODE,
			From:    "IsValidULID",
			Message: "ID inválido: " + id,
			Error:   err,
		})

		isValid = false
	}

	return isValid
}
