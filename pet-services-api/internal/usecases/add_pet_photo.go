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

type AddPetPhotoInput struct {
	UserID      string
	PetID       string
	FileName    string
	ContentType string
	Size        int64
	Reader      io.Reader
}

type AddPetPhotoOutput struct {
	Message string         `json:"message,omitempty"`
	Detail  string         `json:"detail,omitempty"`
	Photo   *entities.Photo `json:"photo,omitempty"`
}

type AddPetPhotoUseCase struct {
	userRepository  entities.UserRepository
	petRepository   entities.PetRepository
	photoRepository entities.PhotoRepository
	storage         storage.ObjectStorage
	logger          logging.LoggerInterface
}

func NewAddPetPhotoUseCase(
	userRepository entities.UserRepository,
	petRepository entities.PetRepository,
	photoRepository entities.PhotoRepository,
	storage storage.ObjectStorage,
	logger logging.LoggerInterface,
) *AddPetPhotoUseCase {
	return &AddPetPhotoUseCase{
		userRepository:  userRepository,
		petRepository:   petRepository,
		photoRepository: photoRepository,
		storage:         storage,
		logger:          logger,
	}
}

func (uc *AddPetPhotoUseCase) Execute(ctx context.Context, input AddPetPhotoInput) (*AddPetPhotoOutput, []exceptions.ProblemDetails) {
	const from = "AddPetPhotoUseCase.Execute"

	if uc.storage == nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Storage não configurado", errors.New("serviço de storage indisponível"))
	}

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.PetID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do pet ausente", errors.New("O ID do pet é obrigatório"))
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

	if !user.IsOwner() {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo owner podem adicionar fotos ao pet"))
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

	ext := strings.ToLower(filepath.Ext(input.FileName))
	if ext == "" {
		ext = ".jpg"
	}

	objectName := fmt.Sprintf("pets/%s/%s%s", pet.ID, ulid.Make().String(), ext)
	if err := uc.storage.Upload(ctx, objectName, input.Reader, input.Size, input.ContentType); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao enviar imagem", err)
	}

	photo, problems := entities.NewPhoto(objectName)
	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Foto inválida", problems)
		return nil, problems
	}

	if err := uc.photoRepository.CreateAndAttachToPet(pet.ID, photo); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao salvar foto", err)
	}

	signedURL, err := uc.storage.GenerateReadURL(ctx, objectName, photoSignedURLTTL)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao gerar URL da foto", err)
	}

	responsePhoto := *photo
	responsePhoto.URL = signedURL

	return &AddPetPhotoOutput{
		Message: "Foto adicionada com sucesso",
		Detail:  "A foto foi vinculada ao pet",
		Photo:   &responsePhoto,
	}, nil
}
