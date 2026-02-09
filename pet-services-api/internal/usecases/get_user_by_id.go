package usecases

import (
	"context"
	"errors"
	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/storage"
)

type GetUserByIDInput struct {
	UserID        string `json:"user_id"`
	RequesterID   string `json:"-"`
	RequesterType string `json:"-"`
}

type GetUserByIDOutput struct {
	User *entities.User `json:"user"`
}

type GetUserByIDUseCase struct {
	userRepository entities.UserRepository
	storage        storage.ObjectStorage
	logger         logging.LoggerInterface
}

func NewGetUserByIDUseCase(userRepository entities.UserRepository, storageService storage.ObjectStorage, logger logging.LoggerInterface) *GetUserByIDUseCase {
	return &GetUserByIDUseCase{
		userRepository: userRepository,
		storage:        storageService,
		logger:         logger,
	}
}

func (uc *GetUserByIDUseCase) Execute(ctx context.Context, input GetUserByIDInput) (*GetUserByIDOutput, []exceptions.ProblemDetails) {
	const from = "GetUserByIDUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório para buscar"))
	}

	if input.RequesterID == "" || input.RequesterType == "" {
		return nil, uc.logger.LogUnauthorized(ctx, from, "Usuário não autenticado", errors.New("Credenciais do usuário autenticado ausentes"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar um usuário com o ID informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if err := signUserPhotos(ctx, uc.storage, user); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao gerar URLs das fotos", err)
	}

	if input.UserID != input.RequesterID {
		switch input.RequesterType {
		case "admin":
		case "owner":
			if user.UserType != "provider" {
				return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Usuários owner só podem consultar usuários provider"))
			}
		default:
			return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Você não tem permissão para acessar dados de outro usuário"))
		}
	}

	return &GetUserByIDOutput{
		User: user,
	}, nil
}
