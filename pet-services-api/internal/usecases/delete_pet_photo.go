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

type DeletePetPhotoInput struct {
	UserID  string
	PetID   string
	PhotoID string
}

type DeletePetPhotoOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type DeletePetPhotoUseCase struct {
	userRepository  entities.UserRepository
	petRepository   entities.PetRepository
	photoRepository entities.PhotoRepository
	storage         storage.ObjectStorage
	logger          logging.LoggerInterface
}

func NewDeletePetPhotoUseCase(
	userRepository entities.UserRepository,
	petRepository entities.PetRepository,
	photoRepository entities.PhotoRepository,
	storage storage.ObjectStorage,
	logger logging.LoggerInterface,
) *DeletePetPhotoUseCase {
	return &DeletePetPhotoUseCase{
		userRepository:  userRepository,
		petRepository:   petRepository,
		photoRepository: photoRepository,
		storage:         storage,
		logger:          logger,
	}
}

func (uc *DeletePetPhotoUseCase) Execute(ctx context.Context, input DeletePetPhotoInput) (*DeletePetPhotoOutput, []exceptions.ProblemDetails) {
	const from = "DeletePetPhotoUseCase.Execute"

	if uc.storage == nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Storage não configurado", errors.New("serviço de storage indisponível"))
	}

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.PetID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do pet ausente", errors.New("O ID do pet é obrigatório"))
	}

	if input.PhotoID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID da foto ausente", errors.New("O ID da foto é obrigatório"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.IsOwner() {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo owner podem remover fotos do pet"))
	}

	pet, err := uc.petRepository.FindByID(input.PetID)
	if err != nil {
		if err.Error() == consts.PetNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Pet não encontrado", errors.New("Não foi possível encontrar o pet informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar pet", err)
	}

	if pet.UserID != user.ID {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("O pet informado não pertence ao usuário"))
	}

	var photoURL string
	for _, photo := range pet.Photos {
		if photo.ID == input.PhotoID {
			photoURL = photo.URL
			break
		}
	}

	if photoURL == "" {
		return nil, uc.logger.LogNotFound(ctx, from, "Foto não encontrada", errors.New(consts.PhotoNotFoundError))
	}

	objectKey := strings.TrimSpace(photoURL)
	if objectKey != "" && !strings.HasPrefix(objectKey, "http") {
		if !strings.Contains(objectKey, "/") {
			objectKey = "pets/" + pet.ID + "/" + objectKey
		}
		_ = uc.storage.Delete(ctx, objectKey)
	}

	if err := uc.photoRepository.DeletePetPhoto(pet.ID, input.PhotoID); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao remover foto do pet", err)
	}

	return &DeletePetPhotoOutput{
		Message: "Foto removida com sucesso",
		Detail:  "A foto foi removida do pet",
	}, nil
}
