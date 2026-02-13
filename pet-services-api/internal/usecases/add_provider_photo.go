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

type AddProviderPhotoInput struct {
	UserID      string
	FileName    string
	ContentType string
	Size        int64
	Reader      io.Reader
}

type AddProviderPhotoOutput struct {
	Message string          `json:"message,omitempty"`
	Detail  string          `json:"detail,omitempty"`
	Photo   *entities.Photo `json:"photo,omitempty"`
}

type AddProviderPhotoUseCase struct {
	userRepository     entities.UserRepository
	providerRepository entities.ProviderRepository
	photoRepository    entities.PhotoRepository
	storage            storage.ObjectStorage
	logger             logging.LoggerInterface
}

func NewAddProviderPhotoUseCase(
	userRepository entities.UserRepository,
	providerRepository entities.ProviderRepository,
	photoRepository entities.PhotoRepository,
	storage storage.ObjectStorage,
	logger logging.LoggerInterface,
) *AddProviderPhotoUseCase {
	return &AddProviderPhotoUseCase{
		userRepository:     userRepository,
		providerRepository: providerRepository,
		photoRepository:    photoRepository,
		storage:            storage,
		logger:             logger,
	}
}

func (uc *AddProviderPhotoUseCase) Execute(ctx context.Context, input AddProviderPhotoInput) (*AddProviderPhotoOutput, []exceptions.ProblemDetails) {
	const from = "AddProviderPhotoUseCase.Execute"

	if uc.storage == nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Storage não configurado", errors.New("serviço de storage indisponível"))
	}

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
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
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo provider podem adicionar fotos ao provedor"))
	}

	provider, err := uc.providerRepository.FindByUserID(user.ID)
	if err != nil {
		if err.Error() == consts.ProviderNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Provedor não encontrado", errors.New("Não foi possível encontrar o provedor do usuário"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar provedor", err)
	}

	count, err := uc.photoRepository.CountProviderPhotos(provider.ID)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao contar fotos do provedor", err)
	}
	if count >= 10 {
		return nil, uc.logger.LogConflict(ctx, from, "Limite de fotos excedido", errors.New("O provedor já possui o máximo de 10 fotos"))
	}

	ext := strings.ToLower(filepath.Ext(input.FileName))
	if ext == "" {
		ext = ".jpg"
	}

	fileName := fmt.Sprintf("%s%s", ulid.Make().String(), ext)
	objectName := fmt.Sprintf("providers/%s/%s", provider.ID, fileName)
	if err := uc.storage.Upload(ctx, objectName, input.Reader, input.Size, input.ContentType); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao enviar imagem", err)
	}

	photo, problems := entities.NewPhoto(fileName)
	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Foto inválida", problems)
		return nil, problems
	}

	if err := uc.photoRepository.CreateAndAttachToProvider(provider.ID, photo); err != nil {
		_ = uc.storage.Delete(ctx, objectName)
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao salvar foto", err)
	}

	signedURL, err := uc.storage.GenerateReadURL(ctx, objectName, storage.PhotoSignedURLTTL)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao gerar URL da foto", err)
	}

	responsePhoto := *photo
	responsePhoto.URL = signedURL

	return &AddProviderPhotoOutput{
		Message: "Foto adicionada com sucesso",
		Detail:  "A foto foi vinculada ao provedor",
		Photo:   &responsePhoto,
	}, nil
}
