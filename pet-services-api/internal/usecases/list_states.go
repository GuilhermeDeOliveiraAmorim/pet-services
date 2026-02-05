package usecases

import (
	"context"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/reference"
)

type ListStatesOutput struct {
	States []reference.State `json:"states"`
}

type ListStatesUseCase struct {
	service reference.Service
	logger  logging.LoggerInterface
}

func NewListStatesUseCase(service reference.Service, logger logging.LoggerInterface) *ListStatesUseCase {
	return &ListStatesUseCase{service: service, logger: logger}
}

func (uc *ListStatesUseCase) Execute(ctx context.Context) (*ListStatesOutput, []exceptions.ProblemDetails) {
	states, err := uc.service.ListStates(ctx)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, "ListStatesUseCase.Execute", "Erro ao buscar estados", err)
	}

	uc.logger.LogInfo(ctx, "ListStatesUseCase.Execute", "Listagem de estados executada")
	return &ListStatesOutput{States: states}, nil
}