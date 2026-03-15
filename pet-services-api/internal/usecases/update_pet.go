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

type UpdatePetInputBody struct {
	Name      string  `json:"name,omitempty"`
	SpeciesID string  `json:"species_id,omitempty"`
	Breed     string  `json:"breed,omitempty"`
	Age       int     `json:"age,omitempty"`
	Weight    float64 `json:"weight,omitempty"`
	Notes     string  `json:"notes,omitempty"`
}

type UpdatePetInput struct {
	UserID    string  `json:"user_id"`
	PetID     string  `json:"pet_id"`
	Name      string  `json:"name,omitempty"`
	SpeciesID string  `json:"species_id,omitempty"`
	Breed     string  `json:"breed,omitempty"`
	Age       int     `json:"age,omitempty"`
	Weight    float64 `json:"weight,omitempty"`
	Notes     string  `json:"notes,omitempty"`
}

type UpdatePetOutput struct {
	Message string        `json:"message,omitempty"`
	Detail  string        `json:"detail,omitempty"`
	Pet     *entities.Pet `json:"pet,omitempty"`
}

type UpdatePetUseCase struct {
	userRepository   entities.UserRepository
	specieRepository entities.SpecieRepository
	petRepository    entities.PetRepository
	storage          storage.ObjectStorage
	logger           logging.LoggerInterface
}

func NewUpdatePetUseCase(
	userRepository entities.UserRepository,
	specieRepository entities.SpecieRepository,
	petRepository entities.PetRepository,
	storage storage.ObjectStorage,
	logger logging.LoggerInterface,
) *UpdatePetUseCase {
	return &UpdatePetUseCase{
		userRepository:   userRepository,
		specieRepository: specieRepository,
		petRepository:    petRepository,
		storage:          storage,
		logger:           logger,
	}
}

func (uc *UpdatePetUseCase) Execute(ctx context.Context, input UpdatePetInput) (*UpdatePetOutput, []exceptions.ProblemDetails) {
	const from = "UpdatePetUseCase.Execute"

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
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo owner podem atualizar pets"))
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

	if input.Name != "" {
		if len(input.Name) > 50 {
			return nil, uc.logger.LogBadRequest(ctx, from, "Nome muito longo", errors.New("O nome do pet deve ter no máximo 50 caracteres"))
		}
		pet.Name = input.Name
	}

	if input.SpeciesID != "" {
		specie, err := uc.specieRepository.FindByID(input.SpeciesID)
		if err != nil {
			if err.Error() == consts.SpecieNotFoundErr {
				return nil, uc.logger.LogNotFound(ctx, from, "Espécie não encontrada", errors.New("Não foi possível encontrar a espécie informada"))
			}
			return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar espécie", err)
		}
		pet.Species = *specie
	}

	if input.Breed != "" {
		if len(input.Breed) > 100 {
			return nil, uc.logger.LogBadRequest(ctx, from, "Raça muito longa", errors.New("A raça do pet deve ter no máximo 100 caracteres"))
		}
		pet.Breed = input.Breed
	}

	if input.Age != 0 {
		if input.Age < 0 {
			return nil, uc.logger.LogBadRequest(ctx, from, "Idade do pet inválida", errors.New("A idade do pet não pode ser negativa"))
		}
		pet.Age = input.Age
	}

	if input.Weight != 0 {
		if input.Weight < 0 {
			return nil, uc.logger.LogBadRequest(ctx, from, "Peso do pet inválido", errors.New("O peso do pet não pode ser negativo"))
		}
		pet.Weight = input.Weight
	}

	if input.Notes != "" {
		if len(input.Notes) > 500 {
			return nil, uc.logger.LogBadRequest(ctx, from, "Observações muito longas", errors.New("As observações devem ter no máximo 500 caracteres"))
		}
		pet.Notes = input.Notes
	}

	pet.Updated()

	if err := uc.petRepository.Update(pet); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao atualizar pet", err)
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

	return &UpdatePetOutput{
		Message: "Pet atualizado com sucesso",
		Detail:  "Os dados do pet foram atualizados",
		Pet:     pet,
	}, nil
}
