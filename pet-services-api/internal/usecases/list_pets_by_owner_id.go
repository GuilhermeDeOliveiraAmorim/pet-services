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

type ListPetsByOwnerIDInput struct {
	RequesterUserID string `json:"requester_user_id"`
	OwnerID         string `json:"owner_id"`
	Page            int    `json:"page"`
	PageSize        int    `json:"page_size"`
}

type ListPetsByOwnerIDUseCase struct {
	userRepository entities.UserRepository
	petRepository  entities.PetRepository
	storage        storage.ObjectStorage
	logger         logging.LoggerInterface
}

func NewListPetsByOwnerIDUseCase(
	userRepository entities.UserRepository,
	petRepository entities.PetRepository,
	storage storage.ObjectStorage,
	logger logging.LoggerInterface,
) *ListPetsByOwnerIDUseCase {
	return &ListPetsByOwnerIDUseCase{
		userRepository: userRepository,
		petRepository:  petRepository,
		storage:        storage,
		logger:         logger,
	}
}

func (uc *ListPetsByOwnerIDUseCase) Execute(ctx context.Context, input ListPetsByOwnerIDInput) (*ListPetsOutput, []exceptions.ProblemDetails) {
	const from = "ListPetsByOwnerIDUseCase.Execute"

	if input.RequesterUserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário autenticado ausente", errors.New("O ID do usuário autenticado é obrigatório"))
	}

	if input.OwnerID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do owner ausente", errors.New("O ID do owner é obrigatório"))
	}

	requester, err := uc.userRepository.FindByID(input.RequesterUserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário autenticado não encontrado", errors.New("Não foi possível encontrar o usuário autenticado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário autenticado", err)
	}

	owner, err := uc.userRepository.FindByID(input.OwnerID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Owner não encontrado", errors.New("Não foi possível encontrar o owner informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar owner", err)
	}

	if !owner.IsOwner() {
		return nil, uc.logger.LogBadRequest(ctx, from, "Usuário informado não é owner", errors.New("O usuário informado não é do tipo owner"))
	}

	canAccess := requester.ID == owner.ID || requester.IsProvider() || requester.UserType == entities.UserTypes.Admin
	if !canAccess {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Sem permissão para listar pets deste owner"))
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

	pets, total, err := uc.petRepository.ListByUser(owner.ID, input.Page, input.PageSize)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao listar pets do owner", err)
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
