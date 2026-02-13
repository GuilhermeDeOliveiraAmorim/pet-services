package usecases

import (
	"context"
	"errors"
	"math"
	"strings"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/storage"
)

type ListPetsInput struct {
	UserID   string `json:"user_id"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

type ListPetsOutput struct {
	Pets         []*entities.Pet `json:"pets"`
	TotalItems   int64           `json:"total_items"`
	TotalPages   int             `json:"total_pages"`
	CurrentPage  int             `json:"current_page"`
	ItemsPerPage int             `json:"items_per_page"`
}

type ListPetsUseCase struct {
	userRepository entities.UserRepository
	petRepository  entities.PetRepository
	storage        storage.ObjectStorage
	logger         logging.LoggerInterface
}

func NewListPetsUseCase(
	userRepository entities.UserRepository,
	petRepository entities.PetRepository,
	storage storage.ObjectStorage,
	logger logging.LoggerInterface,
) *ListPetsUseCase {
	return &ListPetsUseCase{
		userRepository: userRepository,
		petRepository:  petRepository,
		storage:        storage,
		logger:         logger,
	}
}

func (uc *ListPetsUseCase) Execute(ctx context.Context, input ListPetsInput) (*ListPetsOutput, []exceptions.ProblemDetails) {
	const from = "ListPetsUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.IsOwner() {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo owner podem listar pets"))
	}

	if input.Page < 1 {
		input.Page = 1
	}

	if input.PageSize < 1 {
		input.PageSize = 10
	}

	if input.PageSize > 100 {
		input.PageSize = 100
	}

	pets, total, err := uc.petRepository.ListByUser(user.ID, input.Page, input.PageSize)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao listar pets", err)
	}

	for i := range pets {
		if len(pets[i].Photos) == 0 {
			continue
		}
		for j := range pets[i].Photos {
			key := pets[i].Photos[j].URL
			if key != "" && !strings.HasPrefix(key, "http") {
				if !strings.Contains(key, "/") {
					key = "pets/" + pets[i].ID + "/" + key
				}
				url, err := uc.storage.GenerateReadURL(ctx, key, storage.PhotoSignedURLTTL)
				if err == nil {
					pets[i].Photos[j].URL = url
				}
			}
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(input.PageSize)))

	return &ListPetsOutput{
		Pets:         pets,
		TotalItems:   total,
		TotalPages:   totalPages,
		CurrentPage:  input.Page,
		ItemsPerPage: input.PageSize,
	}, nil
}
