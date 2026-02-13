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

type DeleteProviderPhotoInput struct {
	UserID     string
	ProviderID string
	PhotoID    string
}

type DeleteProviderPhotoOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type DeleteProviderPhotoUseCase struct {
	userRepository     entities.UserRepository
	providerRepository entities.ProviderRepository
	photoRepository    entities.PhotoRepository
	storage            storage.ObjectStorage
	logger             logging.LoggerInterface
}

func NewDeleteProviderPhotoUseCase(
	userRepository entities.UserRepository,
	providerRepository entities.ProviderRepository,
	photoRepository entities.PhotoRepository,
	storage storage.ObjectStorage,
	logger logging.LoggerInterface,
) *DeleteProviderPhotoUseCase {
	return &DeleteProviderPhotoUseCase{
		userRepository:     userRepository,
		providerRepository: providerRepository,
		photoRepository:    photoRepository,
		storage:            storage,
		logger:             logger,
	}
}

func (uc *DeleteProviderPhotoUseCase) Execute(ctx context.Context, input DeleteProviderPhotoInput) (*DeleteProviderPhotoOutput, []exceptions.ProblemDetails) {
	const from = "DeleteProviderPhotoUseCase.Execute"

	if uc.storage == nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Storage não configurado", errors.New("serviço de storage indisponível"))
	}

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.ProviderID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do provedor ausente", errors.New("O ID do provedor é obrigatório"))
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

	if !user.IsProvider() {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo provider podem remover fotos do provedor"))
	}

	provider, err := uc.providerRepository.FindByUserID(user.ID)
	if err != nil {
		if err.Error() == consts.ProviderNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Provedor não encontrado", errors.New("Não foi possível encontrar o provedor do usuário"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar provedor", err)
	}

	if provider.ID != input.ProviderID {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("O provedor informado não pertence ao usuário"))
	}

	var photoURL string
	for _, photo := range provider.Photos {
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
			objectKey = "providers/" + provider.ID + "/" + objectKey
		}
		_ = uc.storage.Delete(ctx, objectKey)
	}

	if err := uc.photoRepository.DeleteProviderPhoto(provider.ID, input.PhotoID); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao remover foto do provedor", err)
	}

	return &DeleteProviderPhotoOutput{
		Message: "Foto removida com sucesso",
		Detail:  "A foto foi removida do provedor",
	}, nil
}
