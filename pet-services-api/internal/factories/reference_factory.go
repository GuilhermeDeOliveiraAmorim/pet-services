package factories

import (
	"net/http"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/reference"
	"pet-services-api/internal/usecases"
	"time"
)

type ReferenceFactory struct {
	ListCountries *usecases.ListCountriesUseCase
	ListStates    *usecases.ListStatesUseCase
	ListCities    *usecases.ListCitiesUseCase
}

func NewReferenceFactory(logger logging.LoggerInterface) *ReferenceFactory {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	ibgeClient := reference.NewIBGEClient(httpClient, 24*time.Hour)
	service := reference.NewService(ibgeClient)

	return &ReferenceFactory{
		ListCountries: usecases.NewListCountriesUseCase(service, logger),
		ListStates:    usecases.NewListStatesUseCase(service, logger),
		ListCities:    usecases.NewListCitiesUseCase(service, logger),
	}
}