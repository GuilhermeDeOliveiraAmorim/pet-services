package usecases

import (
	"context"
	"errors"
	"strings"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/storage"
)

type GetPetInput struct {
	UserID string `json:"user_id"`
	PetID  string `json:"pet_id"`
}

type GetPetOutput struct {
	Pet *entities.Pet `json:"pet"`
}

type GetPetUseCase struct {
	userRepository entities.UserRepository
	petRepository  entities.PetRepository
	storage        storage.ObjectStorage
	logger         logging.LoggerInterface
}

func NewGetPetUseCase(
	userRepository entities.UserRepository,
	petRepository entities.PetRepository,
	storage storage.ObjectStorage,
	logger logging.LoggerInterface,
) *GetPetUseCase {
	return &GetPetUseCase{
		userRepository: userRepository,
		petRepository:  petRepository,
		storage:        storage,
		logger:         logger,
	}
}

func (uc *GetPetUseCase) Execute(ctx context.Context, input GetPetInput) (*GetPetOutput, []exceptions.ProblemDetails) {
	const from = "GetPetUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.PetID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do pet ausente", errors.New("O ID do pet é obrigatório"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.IsOwner() {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo owner podem acessar pets"))
	}

	pet, err := uc.petRepository.FindByID(input.PetID)
	if err != nil {
		if err.Error() == consts.PetNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Pet não encontrado", errors.New("Não foi possível encontrar o pet informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar pet", err)
	}

	if !pet.Active {
		return nil, uc.logger.LogBadRequest(ctx, from, "Pet inativo", errors.New("O pet informado está inativo"))
	}

	if pet.UserID != user.ID {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("O pet não pertence ao usuário autenticado"))
	}

	if len(pet.Photos) > 0 {
		for i := range pet.Photos {
			key := pet.Photos[i].URL
			if key != "" && !strings.HasPrefix(key, "http") {
				if !strings.Contains(key, "/") {
					key = "pets/" + pet.ID + "/" + key
				}
				url, err := uc.storage.GenerateReadURL(ctx, key, storage.PhotoSignedURLTTL)
				if err == nil {
					pet.Photos[i].URL = url
				}
			}
		}
	}

	return &GetPetOutput{
		Pet: pet,
	}, nil
}
