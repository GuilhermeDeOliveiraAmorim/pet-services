package usecases

import (
	"context"
	"errors"
	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type AddPetInputBody struct {
	Name     string  `json:"name"`
	SpecieID string  `json:"specie_id"`
	Age      int     `json:"age"`
	Weight   float64 `json:"weight"`
	Notes    string  `json:"notes"`
}

type AddPetInput struct {
	UserID   string  `json:"user_id"`
	Name     string  `json:"name"`
	SpecieID string  `json:"specie_id"`
	Age      int     `json:"age"`
	Weight   float64 `json:"weight"`
	Notes    string  `json:"notes"`
}

type AddPetOutput struct {
	Message string        `json:"message,omitempty"`
	Detail  string        `json:"detail,omitempty"`
	Pet     *entities.Pet `json:"pet,omitempty"`
}

type AddPetUseCase struct {
	userRepository   entities.UserRepository
	specieRepository entities.SpecieRepository
	petRepository    entities.PetRepository
	logger           logging.LoggerInterface
}

func NewAddPetUseCase(
	userRepository entities.UserRepository,
	specieRepository entities.SpecieRepository,
	petRepository entities.PetRepository,
	logger logging.LoggerInterface,
) *AddPetUseCase {
	return &AddPetUseCase{
		userRepository:   userRepository,
		specieRepository: specieRepository,
		petRepository:    petRepository,
		logger:           logger,
	}
}

func (uc *AddPetUseCase) Execute(ctx context.Context, input AddPetInput) (*AddPetOutput, []exceptions.ProblemDetails) {
	const from = "AddPetUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.SpecieID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID da espécie ausente", errors.New("O ID da espécie é obrigatório"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.IsOwner() {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo owner podem adicionar pets"))
	}

	specie, err := uc.specieRepository.FindByID(input.SpecieID)
	if err != nil {
		if err.Error() == consts.SpecieNotFoundErr {
			return nil, uc.logger.LogNotFound(ctx, from, "Espécie não encontrada", errors.New("Não foi possível encontrar a espécie informada"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar espécie", err)
	}

	pet, problems := entities.NewPet(input.UserID, input.Name, *specie, input.Age, input.Weight, input.Notes)
	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Pet inválido", problems)
		return nil, problems
	}

	if err := uc.petRepository.Create(pet); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao adicionar pet", err)
	}

	return &AddPetOutput{
		Message: "Pet adicionado com sucesso",
		Detail:  "O pet foi cadastrado para o usuário owner",
		Pet:     pet,
	}, nil
}
