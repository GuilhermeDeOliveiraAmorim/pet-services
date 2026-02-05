package usecases

import (
	"context"
	"errors"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type ListUsersInput struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type PaginationMetadata struct {
	CurrentPage  int   `json:"current_page"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
	PerPage      int   `json:"per_page"`
}

type ListUsersOutput struct {
	Users      []*entities.User   `json:"users"`
	Pagination PaginationMetadata `json:"pagination"`
}

type ListUsersUseCase struct {
	userRepository entities.UserRepository
	logger         logging.LoggerInterface
}

func NewListUsersUseCase(userRepository entities.UserRepository, logger logging.LoggerInterface) *ListUsersUseCase {
	return &ListUsersUseCase{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (uc *ListUsersUseCase) Execute(ctx context.Context, input ListUsersInput) (*ListUsersOutput, []exceptions.ProblemDetails) {
	const from = "ListUsersUseCase.Execute"
	const maxLimit = 100

	if input.Page < 1 {
		input.Page = 1
	}

	if input.Limit < 1 {
		input.Limit = 10
	}

	if input.Limit > maxLimit {
		return nil, uc.logger.LogBadRequest(ctx, from, "Limite de paginação excedido", errors.New("O limite máximo permitido é 100"))
	}

	users, total, err := uc.userRepository.List(input.Page, input.Limit)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao listar usuários", err)
	}

	totalPages := int(total) / input.Limit
	if int(total)%input.Limit != 0 {
		totalPages++
	}

	return &ListUsersOutput{
		Users: users,
		Pagination: PaginationMetadata{
			CurrentPage:  input.Page,
			TotalPages:   totalPages,
			TotalRecords: total,
			PerPage:      input.Limit,
		},
	}, nil
}
