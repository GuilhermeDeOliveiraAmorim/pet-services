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

type AddUserPhotoInput struct {
	UserID      string
	FileName    string
	ContentType string
	Size        int64
	Reader      io.Reader
}

type AddUserPhotoOutput struct {
	Message string        `json:"message,omitempty"`
	Detail  string        `json:"detail,omitempty"`
	Photo   *entities.Photo `json:"photo,omitempty"`
}

type AddUserPhotoUseCase struct {
	userRepository  entities.UserRepository
	photoRepository entities.PhotoRepository
	storage         storage.ObjectStorage
	logger          logging.LoggerInterface
}

func NewAddUserPhotoUseCase(
	userRepository entities.UserRepository,
	photoRepository entities.PhotoRepository,
	storage storage.ObjectStorage,
	logger logging.LoggerInterface,
) *AddUserPhotoUseCase {
	return &AddUserPhotoUseCase{
		userRepository:  userRepository,
		photoRepository: photoRepository,
		storage:         storage,
		logger:          logger,
	}
}

func (uc *AddUserPhotoUseCase) Execute(ctx context.Context, input AddUserPhotoInput) (*AddUserPhotoOutput, []exceptions.ProblemDetails) {
	const from = "AddUserPhotoUseCase.Execute"

	if uc.storage == nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Storage não configurado", errors.New("serviço de storage indisponível"))
	}

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.Reader == nil {
		return nil, uc.logger.LogBadRequest(ctx, from, "Arquivo ausente", errors.New("O arquivo é obrigatório"))
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

	ext := strings.ToLower(filepath.Ext(input.FileName))
	if ext == "" {
		ext = ".jpg"
	}

	objectName := fmt.Sprintf("users/%s/%s%s", user.ID, ulid.Make().String(), ext)
	url, err := uc.storage.UploadImage(ctx, objectName, input.Reader, input.Size, input.ContentType)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao enviar imagem", err)
	}

	photo, problems := entities.NewPhoto(url)
	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Foto inválida", problems)
		return nil, problems
	}

	if err := uc.photoRepository.CreateAndAttachToUser(user.ID, photo); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao salvar foto", err)
	}

	return &AddUserPhotoOutput{
		Message: "Foto adicionada com sucesso",
		Detail:  "A foto foi vinculada ao usuário",
		Photo:   photo,
	}, nil
}
