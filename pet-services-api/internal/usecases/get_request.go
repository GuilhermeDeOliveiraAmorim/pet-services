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

type GetRequestInput struct {
	UserID    string `json:"user_id"`
	RequestID string `json:"request_id"`
}

type GetRequestOutput struct {
	Request *entities.Request `json:"request"`
}

type GetRequestUseCase struct {
	userRepository     entities.UserRepository
	providerRepository entities.ProviderRepository
	requestRepository  entities.RequestRepository
	storage            storage.ObjectStorage
	logger             logging.LoggerInterface
}

func NewGetRequestUseCase(
	userRepository entities.UserRepository,
	providerRepository entities.ProviderRepository,
	requestRepository entities.RequestRepository,
	storage storage.ObjectStorage,
	logger logging.LoggerInterface,
) *GetRequestUseCase {
	return &GetRequestUseCase{
		userRepository:     userRepository,
		providerRepository: providerRepository,
		requestRepository:  requestRepository,
		storage:            storage,
		logger:             logger,
	}
}

func (uc *GetRequestUseCase) Execute(ctx context.Context, input GetRequestInput) (*GetRequestOutput, []exceptions.ProblemDetails) {
	const from = "GetRequestUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.RequestID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID da solicitação ausente", errors.New("O ID da solicitação é obrigatório"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	request, err := uc.requestRepository.FindByID(input.RequestID)
	if err != nil {
		if err.Error() == consts.RequestNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Solicitação não encontrada", errors.New("Não foi possível encontrar a solicitação informada"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar solicitação", err)
	}

	if user.IsOwner() {
		if request.UserID != user.ID {
			return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("A solicitação não pertence ao usuário"))
		}
	} else if user.IsProvider() {
		provider, err := uc.providerRepository.FindByUserID(user.ID)
		if err != nil {
			if err.Error() == consts.ProviderNotFoundError {
				return nil, uc.logger.LogNotFound(ctx, from, "Provedor não encontrado", errors.New("Não foi possível encontrar o provedor do usuário"))
			}
			return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar provedor", err)
		}
		if request.ProviderID != provider.ID {
			return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("A solicitação não pertence ao provedor"))
		}
	} else {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente owners ou providers podem acessar solicitações"))
	}

	if len(request.Pet.Photos) > 0 {
		for i := range request.Pet.Photos {
			key := request.Pet.Photos[i].URL
			if key != "" && !strings.HasPrefix(key, "http") {
				if !strings.Contains(key, "/") {
					key = "pets/" + request.Pet.ID + "/" + key
				}
				url, err := uc.storage.GenerateReadURL(ctx, key, storage.PhotoSignedURLTTL)
				if err == nil {
					request.Pet.Photos[i].URL = url
				}
			}
		}
	}

	return &GetRequestOutput{
		Request: request,
	}, nil
}
