package usecases

import (
	"context"
	"errors"
	"strings"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
)

type AddServiceTagInput struct {
	UserID    string `json:"user_id"`
	ServiceID string `json:"service_id"`
	TagID     string `json:"tag_id"`
	TagName   string `json:"tag_name"`
}

type AddServiceTagOutput struct {
	Message string         `json:"message,omitempty"`
	Detail  string         `json:"detail,omitempty"`
	Tag     *entities.Tag  `json:"tag,omitempty"`
}

type AddServiceTagUseCase struct {
	userRepository     entities.UserRepository
	serviceRepository  entities.ServiceRepository
	providerRepository entities.ProviderRepository
	tagRepository      entities.TagRepository
	logger             logging.LoggerInterface
}

func NewAddServiceTagUseCase(
	userRepository entities.UserRepository,
	serviceRepository entities.ServiceRepository,
	providerRepository entities.ProviderRepository,
	tagRepository entities.TagRepository,
	logger logging.LoggerInterface,
) *AddServiceTagUseCase {
	return &AddServiceTagUseCase{
		userRepository:     userRepository,
		serviceRepository:  serviceRepository,
		providerRepository: providerRepository,
		tagRepository:      tagRepository,
		logger:             logger,
	}
}

func (uc *AddServiceTagUseCase) Execute(ctx context.Context, input AddServiceTagInput) (*AddServiceTagOutput, []exceptions.ProblemDetails) {
	const from = "AddServiceTagUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	if input.ServiceID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do serviço ausente", errors.New("O ID do serviço é obrigatório"))
	}

	if input.TagID == "" && strings.TrimSpace(input.TagName) == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "Tag ausente", errors.New("Informe o ID da tag ou o nome da tag"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if !user.IsProvider() {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo provider podem adicionar tags ao serviço"))
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

	var tag *entities.Tag
	if input.TagID != "" {
		found, err := uc.tagRepository.FindByID(input.TagID)
		if err != nil {
			if err.Error() != consts.TagNotFoundError {
				return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar tag", err)
			}
		} else {
			tag = found
		}
	}

	if tag == nil {
		name := strings.TrimSpace(input.TagName)
		if name == "" {
			return nil, uc.logger.LogNotFound(ctx, from, "Tag não encontrada", errors.New("Não foi possível encontrar a tag informada"))
		}
		found, err := uc.tagRepository.FindByName(name)
		if err != nil {
			if err.Error() != consts.TagNotFoundError {
				return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar tag", err)
			}
		} else {
			tag = found
		}
	}

	if tag == nil {
		newTag, problems := entities.NewTag(strings.TrimSpace(input.TagName))
		if len(problems) > 0 {
			uc.logger.LogMultipleBadRequests(ctx, from, "Tag inválida", problems)
			return nil, problems
		}
		if err := uc.tagRepository.Create(newTag); err != nil {
			return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao criar tag", err)
		}
		tag = newTag
	}

	if !tag.Active {
		return nil, uc.logger.LogBadRequest(ctx, from, "Tag inativa", errors.New("A tag informada está inativa"))
	}

	exists, err := uc.serviceRepository.HasTag(service.ID, tag.ID)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao verificar tags do serviço", err)
	}
	if exists {
		return nil, uc.logger.LogConflict(ctx, from, "Tag já vinculada", errors.New("A tag já está vinculada ao serviço"))
	}

	if err := uc.serviceRepository.AddTag(service.ID, tag.ID); err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao adicionar tag ao serviço", err)
	}

	return &AddServiceTagOutput{
		Message: "Tag adicionada com sucesso",
		Detail:  "A tag foi vinculada ao serviço",
		Tag:     tag,
	}, nil
}
