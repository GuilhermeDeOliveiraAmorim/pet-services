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

type DeleteServicePhotoInput struct {
	UserID    string
	ServiceID string
	PhotoID   string
}

type DeleteServicePhotoOutput struct {
	Message string `json:"message,omitempty"`
	Detail  string `json:"detail,omitempty"`
}

type DeleteServicePhotoUseCase struct {
	userRepository     entities.UserRepository
	serviceRepository  entities.ServiceRepository
	providerRepository entities.ProviderRepository
	photoRepository    entities.PhotoRepository
	storage            storage.ObjectStorage
	logger             logging.LoggerInterface
}

func NewDeleteServicePhotoUseCase(
	userRepository entities.UserRepository,
	serviceRepository entities.ServiceRepository,
	providerRepository entities.ProviderRepository,
	photoRepository entities.PhotoRepository,
	storage storage.ObjectStorage,
	logger logging.LoggerInterface,
) *DeleteServicePhotoUseCase {
	return &DeleteServicePhotoUseCase{
		userRepository:     userRepository,
		serviceRepository:  serviceRepository,
		providerRepository: providerRepository,
		photoRepository:    photoRepository,
		storage:            storage,
		logger:             logger,
	}
}

func (uc *DeleteServicePhotoUseCase) Execute(ctx context.Context, input DeleteServicePhotoInput) (*DeleteServicePhotoOutput, []exceptions.ProblemDetails) {
	const from = "DeleteServicePhotoUseCase.Execute"

	if uc.storage == nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Storage não configurado", errors.New("serviço de storage indisponível"))
	}

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.ServiceID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do serviço ausente", errors.New("O ID do serviço é obrigatório"))
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
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo provider podem remover fotos do serviço"))
	}

	provider, err := uc.providerRepository.FindByUserID(user.ID)
	if err != nil {
		if err.Error() == consts.ProviderNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Provedor não encontrado", errors.New("Não foi possível encontrar o provedor do usuário"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar provedor", err)
	}

	service, err := uc.serviceRepository.FindByID(input.ServiceID)
	if err != nil {
		if err.Error() == consts.ServiceNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Serviço não encontrado", errors.New("Não foi possível encontrar o serviço informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar serviço", err)
	}

	if service.ProviderID != provider.ID {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("O serviço informado não pertence ao provedor autenticado"))
	}

	if !service.Active {
		return nil, uc.logger.LogBadRequest(ctx, from, "Serviço inativo", errors.New("O serviço informado está inativo"))
	}

	var photoURL string
	for _, photo := range service.Photos {
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
			objectKey = "services/" + service.ID + "/" + objectKey
		}
		_ = uc.storage.Delete(ctx, objectKey)
	}

	if err := uc.photoRepository.DeleteServicePhoto(service.ID, input.PhotoID); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao remover foto do serviço", err)
	}

	return &DeleteServicePhotoOutput{
		Message: "Foto removida com sucesso",
		Detail:  "A foto foi removida do serviço",
	}, nil
}
