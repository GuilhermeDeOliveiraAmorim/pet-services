package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type DeletePetInput struct {
	UserID string `json:"user_id"`
	PetID  string `json:"pet_id"`
}

type DeletePetOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type DeletePetUseCase struct {
	userRepository entities.UserRepository
	petRepository  entities.PetRepository
	logger         logging.LoggerInterface
}

func NewDeletePetUseCase(
	userRepository entities.UserRepository,
	petRepository entities.PetRepository,
	logger logging.LoggerInterface,
) *DeletePetUseCase {
	return &DeletePetUseCase{
		userRepository: userRepository,
		petRepository:  petRepository,
		logger:         logger,
	}
}

func (uc *DeletePetUseCase) Execute(ctx context.Context, input DeletePetInput) (*DeletePetOutput, []exceptions.ProblemDetails) {
	const from = "DeletePetUseCase.Execute"

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
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo owner podem remover pets"))
	}

	pet, err := uc.petRepository.FindByID(input.PetID)
	if err != nil {
		if err.Error() == consts.PetNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Pet não encontrado", errors.New("Não foi possível encontrar o pet informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar pet", err)
	}

	if pet.UserID != user.ID {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("O pet não pertence ao usuário autenticado"))
	}

	if !pet.Active {
		return nil, uc.logger.LogBadRequest(ctx, from, "Pet já inativo", errors.New("O pet já está inativo no sistema"))
	}

	pet.Deactivate()

	if err := uc.petRepository.Update(pet); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao remover pet", err)
	}

	return &DeletePetOutput{
		Message: "Pet removido com sucesso",
		Detail:  "O pet foi removido do sistema com sucesso",
	}, nil
}
