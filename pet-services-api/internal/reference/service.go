package reference

import "context"

type Service interface {
	ListCountries() []Country
	ListStates(ctx context.Context) ([]State, error)
	ListCities(ctx context.Context, stateID *int) ([]City, error)
}

type service struct {
	client *IBGEClient
}

func NewService(client *IBGEClient) Service {
	return &service{client: client}
}

func (s *service) ListCountries() []Country {
	return Countries
}

func (s *service) ListStates(ctx context.Context) ([]State, error) {
	return s.client.ListStates(ctx)
}

func (s *service) ListCities(ctx context.Context, stateID *int) ([]City, error) {
	return s.client.ListCities(ctx, stateID)
}