package usecases

import (
	"context"
	"errors"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type AddRequestInputBody struct {
	ServiceID string `json:"service_id"`
	PetID     string `json:"pet_id"`
	Notes     string `json:"notes"`
}

type AddRequestInput struct {
	UserID    string `json:"user_id"`
	ServiceID string `json:"service_id"`
	PetID     string `json:"pet_id"`
	Notes     string `json:"notes"`
}

type AddRequestOutput struct {
	Message string             `json:"message,omitempty"`
	Detail  string             `json:"detail,omitempty"`
	Request *entities.Request  `json:"request,omitempty"`
}

type AddRequestUseCase struct {
	userRepository     entities.UserRepository
	petRepository      entities.PetRepository
	serviceRepository  entities.ServiceRepository
	providerRepository entities.ProviderRepository
	requestRepository  entities.RequestRepository
	logger             logging.LoggerInterface
}

func NewAddRequestUseCase(
	userRepository entities.UserRepository,
	petRepository entities.PetRepository,
	serviceRepository entities.ServiceRepository,
	providerRepository entities.ProviderRepository,
	requestRepository entities.RequestRepository,
	logger logging.LoggerInterface,
) *AddRequestUseCase {
	return &AddRequestUseCase{
		userRepository:     userRepository,
		petRepository:      petRepository,
		serviceRepository:  serviceRepository,
		providerRepository: providerRepository,
		requestRepository:  requestRepository,
		logger:             logger,
	}
}

func (uc *AddRequestUseCase) Execute(ctx context.Context, input AddRequestInput) (*AddRequestOutput, []exceptions.ProblemDetails) {
	const from = "AddRequestUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.ServiceID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do serviço ausente", errors.New("O ID do serviço é obrigatório"))
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
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo owner podem criar solicitações"))
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

	service, err := uc.serviceRepository.FindByID(input.ServiceID)
	if err != nil {
		if err.Error() == consts.ServiceNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Serviço não encontrado", errors.New("Não foi possível encontrar o serviço informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar serviço", err)
	}

	if !service.Active {
		return nil, uc.logger.LogBadRequest(ctx, from, "Serviço inativo", errors.New("O serviço informado está inativo"))
	}

	provider, err := uc.providerRepository.FindByID(service.ProviderID)
	if err != nil {
		if err.Error() == consts.ProviderNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Provedor não encontrado", errors.New("Não foi possível encontrar o provedor do serviço"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar provedor", err)
	}

	if !provider.Active {
		return nil, uc.logger.LogBadRequest(ctx, from, "Provedor inativo", errors.New("O provedor do serviço está inativo"))
	}

	exists, err := uc.requestRepository.ExistsPending(user.ID, service.ID, pet.ID)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao verificar solicitações", err)
	}
	if exists {
		return nil, uc.logger.LogConflict(ctx, from, "Solicitação já existente", errors.New("Já existe uma solicitação pendente para este pet e serviço"))
	}

	requestEntity, problems := entities.NewRequest(user.ID, provider.ID, service.ID, *pet, input.Notes)
	if len(problems) > 0 {
		uc.logger.LogMultipleBadRequests(ctx, from, "Solicitação inválida", problems)
		return nil, problems
	}

	if err := uc.requestRepository.Create(requestEntity); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao criar solicitação", err)
	}

	return &AddRequestOutput{
		Message: "Solicitação criada com sucesso",
		Detail:  "A solicitação foi enviada para o provedor",
		Request: requestEntity,
	}, nil
}
