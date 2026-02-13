package usecases

import (
	"context"
	"errors"
	"math"
	"slices"
	"strings"

	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/storage"
)

type ListRequestsInput struct {
	UserID     string
	ProviderID string
	Status     string
	Page       int
	PageSize   int
}

type RequestListItem struct {
	ID           string        `json:"id"`
	UserID       string        `json:"user_id"`
	UserName     string        `json:"user_name"`
	ProviderID   string        `json:"provider_id"`
	BusinessName string        `json:"business_name"`
	ServiceID    string        `json:"service_id"`
	ServiceName  string        `json:"service_name"`
	Pet          *entities.Pet `json:"pet"`
	Notes        string        `json:"notes"`
	Status       string        `json:"status"`
	RejectReason string        `json:"reject_reason,omitempty"`
	CreatedAt    string        `json:"created_at"`
	UpdatedAt    string        `json:"updated_at,omitempty"`
}

type ListRequestsOutput struct {
	Requests     []*RequestListItem `json:"requests"`
	TotalItems   int64              `json:"total_items"`
	TotalPages   int                `json:"total_pages"`
	CurrentPage  int                `json:"current_page"`
	ItemsPerPage int                `json:"items_per_page"`
}

type ListRequestsUseCase struct {
	userRepository     entities.UserRepository
	requestRepository  entities.RequestRepository
	providerRepository entities.ProviderRepository
	serviceRepository  entities.ServiceRepository
	storage            storage.ObjectStorage
	logger             logging.LoggerInterface
}

func NewListRequestsUseCase(
	userRepository entities.UserRepository,
	requestRepository entities.RequestRepository,
	providerRepository entities.ProviderRepository,
	serviceRepository entities.ServiceRepository,
	storage storage.ObjectStorage,
	logger logging.LoggerInterface,
) *ListRequestsUseCase {
	return &ListRequestsUseCase{
		userRepository:     userRepository,
		requestRepository:  requestRepository,
		providerRepository: providerRepository,
		serviceRepository:  serviceRepository,
		storage:            storage,
		logger:             logger,
	}
}

func (uc *ListRequestsUseCase) Execute(ctx context.Context, input ListRequestsInput) (*ListRequestsOutput, []exceptions.ProblemDetails) {
	const from = "ListRequestsUseCase.Execute"

	if input.UserID == "" {
		return nil, uc.logger.LogBadRequest(ctx, from, "ID do usuário ausente", errors.New("O ID do usuário é obrigatório"))
	}

	user, err := uc.userRepository.FindByID(input.UserID)
	if err != nil {
		if err.Error() == consts.UserNotFoundError {
			return nil, uc.logger.LogNotFound(ctx, from, "Usuário não encontrado", errors.New("Não foi possível encontrar o usuário informado"))
		}
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar usuário", err)
	}

	if input.Page < 1 {
		input.Page = 1
	}

	if input.PageSize < 1 {
		input.PageSize = 10
	}

	if input.PageSize > 100 {
		input.PageSize = 100
	}

	if input.Status != "" {
		validStatuses := []string{
			entities.RequestStatuses.Pending,
			entities.RequestStatuses.Accepted,
			entities.RequestStatuses.Rejected,
			entities.RequestStatuses.Completed,
		}
		isValid := slices.Contains(validStatuses, input.Status)
		if !isValid {
			return nil, uc.logger.LogBadRequest(ctx, from, "Status inválido", errors.New("O status informado não é válido"))
		}
	}

	var userIDFilter, providerIDFilter string
	if user.IsOwner() {
		userIDFilter = input.UserID
		providerIDFilter = ""
	} else if user.IsProvider() {
		provider, err := uc.providerRepository.FindByUserID(input.UserID)
		if err != nil {
			if err.Error() == consts.ProviderNotFoundError {
				return nil, uc.logger.LogNotFound(ctx, from, "Provedor não encontrado", errors.New("Não foi possível encontrar o provedor do usuário"))
			}
			return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao buscar provedor", err)
		}
		userIDFilter = ""
		providerIDFilter = provider.ID
	} else {
		return nil, uc.logger.LogForbidden(ctx, from, "Acesso negado", errors.New("Somente usuários do tipo owner ou provider podem listar solicitações"))
	}

	requests, total, err := uc.requestRepository.List(userIDFilter, providerIDFilter, input.Status, input.Page, input.PageSize)
	if err != nil {
		return nil, uc.logger.LogInternalServerError(ctx, from, "Erro ao listar solicitações", err)
	}

	totalPages := int(math.Ceil(float64(total) / float64(input.PageSize)))

	requestItems := make([]*RequestListItem, len(requests))
	for i, req := range requests {
		user, _ := uc.userRepository.FindByID(req.UserID)
		provider, _ := uc.providerRepository.FindByID(req.ProviderID)
		service, _ := uc.serviceRepository.FindByID(req.ServiceID)

		userName := ""
		if user != nil {
			userName = user.Name
		}

		businessName := ""
		if provider != nil {
			businessName = provider.BusinessName
		}

		serviceName := ""
		if service != nil {
			serviceName = service.Name
		}

		if len(req.Pet.Photos) > 0 {
			for j := range req.Pet.Photos {
				key := req.Pet.Photos[j].URL
				if key != "" && !strings.HasPrefix(key, "http") {
					if !strings.Contains(key, "/") {
						key = "pets/" + req.Pet.ID + "/" + key
					}
					url, err := uc.storage.GenerateReadURL(ctx, key, storage.PhotoSignedURLTTL)
					if err == nil {
						req.Pet.Photos[j].URL = url
					}
				}
			}
		}

		updatedAt := ""
		if req.UpdatedAt != nil {
			updatedAt = req.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		requestItems[i] = &RequestListItem{
			ID:           req.ID,
			UserID:       req.UserID,
			UserName:     userName,
			ProviderID:   req.ProviderID,
			BusinessName: businessName,
			ServiceID:    req.ServiceID,
			ServiceName:  serviceName,
			Pet:          &req.Pet,
			Notes:        req.Notes,
			Status:       req.Status,
			RejectReason: req.RejectReason,
			CreatedAt:    req.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    updatedAt,
		}
	}

	return &ListRequestsOutput{
		Requests:     requestItems,
		TotalItems:   total,
		TotalPages:   totalPages,
		CurrentPage:  input.Page,
		ItemsPerPage: input.PageSize,
	}, nil
}
