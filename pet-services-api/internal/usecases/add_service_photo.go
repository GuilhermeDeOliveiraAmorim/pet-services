package usecases

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/storage"

	"github.com/oklog/ulid/v2"
)

type AddServicePhotoInput struct {
	UserID      string
	ServiceID   string
	FileName    string
	ContentType string
	Size        int64
	Reader      io.Reader
}

type AddServicePhotoOutput struct {
	Message string          `json:"message,omitempty"`
	Detail  string          `json:"detail,omitempty"`
	Photo   *entities.Photo `json:"photo,omitempty"`
}

type AddServicePhotoUseCase struct {
	userRepository     entities.UserRepository
	serviceRepository  entities.ServiceRepository
	providerRepository entities.ProviderRepository
	photoRepository    entities.PhotoRepository
	storage            storage.ObjectStorage
	logger             logging.LoggerInterface
}

func NewAddServicePhotoUseCase(
	userRepository entities.UserRepository,
	serviceRepository entities.ServiceRepository,
	providerRepository entities.ProviderRepository,
	photoRepository entities.PhotoRepository,
	storage storage.ObjectStorage,
	logger logging.LoggerInterface,
) *AddServicePhotoUseCase {
	return &AddServicePhotoUseCase{
		userRepository:     userRepository,
		serviceRepository:  serviceRepository,
		providerRepository: providerRepository,
		photoRepository:    photoRepository,
		storage:            storage,
		logger:             logger,
	}
}

func (uc *AddServicePhotoUseCase) Execute(ctx context.Context, input AddServicePhotoInput) (*AddServicePhotoOutput, []exceptions.ProblemDetails) {
	const from = "AddServicePhotoUseCase.Execute"

	if uc.storage == nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Storage não configurado", errors.New("serviço de storage indisponível"))
	}

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.ServiceID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do serviço ausente", errors.New("O ID do serviço é obrigatório"))
	}

	if input.Reader == nil {
		return nil, uc.logger.LogBadRequest(ctx, from, "Arquivo ausente", errors.New("A imagem é obrigatória"))
	}

	if input.ContentType == "" || !strings.HasPrefix(input.ContentType, "image/") {
		return nil, uc.logger.LogBadRequest(ctx, from, "Tipo de arquivo inválido", errors.New("Apenas imagens são permitidas"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.IsProvider() {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo provider podem adicionar fotos ao serviço"))
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

	count, err := uc.photoRepository.CountServicePhotos(service.ID)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao contar fotos do serviço", err)
	}
	if count >= 10 {
		return nil, uc.logger.LogConflict(ctx, from, "Limite de fotos excedido", errors.New("O serviço já possui o máximo de 10 fotos"))
	}

	ext := strings.ToLower(filepath.Ext(input.FileName))
	if ext == "" {
		ext = ".jpg"
	}

	fileName := fmt.Sprintf("%s%s", ulid.Make().String(), ext)
	objectName := fmt.Sprintf("services/%s/%s", service.ID, fileName)
	if err := uc.storage.Upload(ctx, objectName, input.Reader, input.Size, input.ContentType); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao enviar imagem", err)
	}

	photo, problems := entities.NewPhoto(fileName)
	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Foto inválida", problems)
		return nil, problems
	}

	if err := uc.photoRepository.CreateAndAttachToService(service.ID, photo); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao salvar foto", err)
	}

	signedURL, err := uc.storage.GenerateReadURL(ctx, objectName, storage.PhotoSignedURLTTL)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao gerar URL da foto", err)
	}

	responsePhoto := *photo
	responsePhoto.URL = signedURL

	return &AddServicePhotoOutput{
		Message: "Foto adicionada com sucesso",
		Detail:  "A foto foi vinculada ao serviço",
		Photo:   &responsePhoto,
	}, nil
}
