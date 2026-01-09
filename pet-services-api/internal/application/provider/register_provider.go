package provider

import (
	"context"
	"fmt"

	"pet-services-api/internal/application/exceptions"
	"pet-services-api/internal/application/logging"
	"pet-services-api/internal/domain/provider"
	"pet-services-api/internal/domain/user"

	"github.com/google/uuid"
)

// RegisterProviderUseCase orquestra a criaç
type RegisterProviderUseCase struct {
	providerRepo provider.Repository
	userRepo     user.Repository
	logger       logging.LoggerService
}

// NewRegisterProviderUseCase cria um novo caso de uso.
func NewRegisterProviderUseCase(providerRepo provider.Repository, userRepo user.Repository, logger logging.LoggerService) *RegisterProviderUseCase {
	return &RegisterProviderUseCase{
		providerRepo: providerRepo,
		userRepo:     userRepo,
		logger:       logging.NewSlogLogger(),
	}
}

// RegisterProviderInput representa os dados necessários para criar um prestador.
type RegisterProviderInput struct {
	UserID       uuid.UUID
	BusinessName string
	Description  string
	Address      user.Address
	Latitude     float64
	Longitude    float64
	Services     []ServiceInput
	PriceRange   provider.PriceRange
	Photos       []string // URLs das fotos do prestador
}

// ServiceInput representa um serviço a ser cadastrado.
type ServiceInput struct {
	Category string
	Name     string
	PriceMin float64
	PriceMax float64
}

// RegisterProviderOutput representa o resultado da criação.
type RegisterProviderOutput struct {
	ProviderID   uuid.UUID
	BusinessName string
	IsActive     bool
}

// Execute cria um perfil de prestador validando regras básicas de domínio.
func (uc *RegisterProviderUseCase) Execute(ctx context.Context, input RegisterProviderInput) (*provider.Provider, []exceptions.ProblemDetails) {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "RegisterProviderUseCase",
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.UserID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "RegisterProviderUseCase",
			Message: "UserID é obrigatório",
			Error:   fmt.Errorf("userID é obrigatório"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "UserID é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O ID do usuário é obrigatório.",
		}}
	}
	if input.BusinessName == "" {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "RegisterProviderUseCase",
			Message: "Nome do negócio é obrigatório",
			Error:   fmt.Errorf("nome do negócio é obrigatório"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Nome do negócio é obrigatório",
			Status: exceptions.RFC400_CODE,
			Detail: "O nome do negócio é obrigatório.",
		}}
	}
	if len(input.Services) == 0 {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "RegisterProviderUseCase",
			Message: "Nenhum serviço informado",
			Error:   provider.ErrNoServicesProvided,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Nenhum serviço informado",
			Status: exceptions.RFC400_CODE,
			Detail: "É necessário informar ao menos um serviço.",
		}}
	}
	for _, svc := range input.Services {
		if svc.Name == "" {
			uc.logger.Log(logging.Logger{
				Context: ctx,
				TypeLog: logging.LoggerTypes.ERROR,
				Layer:   logging.LoggerLayers.USECASES,
				Code:    exceptions.RFC400_CODE,
				From:    "RegisterProviderUseCase",
				Message: "Nome do serviço é obrigatório",
				Error:   fmt.Errorf("nome do serviço é obrigatório"),
			})
			return nil, []exceptions.ProblemDetails{{
				Type:   exceptions.RFC400,
				Title:  "Nome do serviço é obrigatório",
				Status: exceptions.RFC400_CODE,
				Detail: "O nome do serviço é obrigatório.",
			}}
		}
		if svc.PriceMin < 0 || svc.PriceMax < 0 {
			uc.logger.Log(logging.Logger{
				Context: ctx,
				TypeLog: logging.LoggerTypes.ERROR,
				Layer:   logging.LoggerLayers.USECASES,
				Code:    exceptions.RFC400_CODE,
				From:    "RegisterProviderUseCase",
				Message: "Preço não pode ser negativo",
				Error:   fmt.Errorf("preço não pode ser negativo"),
			})
			return nil, []exceptions.ProblemDetails{{
				Type:   exceptions.RFC400,
				Title:  "Preço não pode ser negativo",
				Status: exceptions.RFC400_CODE,
				Detail: "O preço do serviço não pode ser negativo.",
			}}
		}
		if svc.PriceMax > 0 && svc.PriceMin > svc.PriceMax {
			uc.logger.Log(logging.Logger{
				Context: ctx,
				TypeLog: logging.LoggerTypes.ERROR,
				Layer:   logging.LoggerLayers.USECASES,
				Code:    exceptions.RFC400_CODE,
				From:    "RegisterProviderUseCase",
				Message: "Preço máximo deve ser maior ou igual ao mínimo",
				Error:   fmt.Errorf("preço máximo deve ser maior ou igual ao mínimo"),
			})
			return nil, []exceptions.ProblemDetails{{
				Type:   exceptions.RFC400,
				Title:  "Preço máximo deve ser maior ou igual ao mínimo",
				Status: exceptions.RFC400_CODE,
				Detail: "O preço máximo deve ser maior ou igual ao mínimo.",
			}}
		}
	}

	u, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    "RegisterProviderUseCase",
			Message: "Usuário não encontrado",
			Error:   user.ErrUserNotFound,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Usuário não encontrado",
			Status: exceptions.RFC404_CODE,
			Detail: "O usuário informado não foi encontrado.",
		}}
	}
	if u.Type != user.UserTypeProvider {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC400_CODE,
			From:    "RegisterProviderUseCase",
			Message: "Usuário não é do tipo provider",
			Error:   fmt.Errorf("usuário não é do tipo provider"),
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC400,
			Title:  "Usuário não é do tipo provider",
			Status: exceptions.RFC400_CODE,
			Detail: "O usuário informado não é do tipo provider.",
		}}
	}

	p := provider.NewProvider(input.UserID, input.BusinessName, input.Description)
	p.SetLocation(input.Latitude, input.Longitude, input.Address)
	for _, svc := range input.Services {
		p.AddService(svc.Category, svc.Name, svc.PriceMin, svc.PriceMax)
	}

	if input.PriceRange != "" {
		p.PriceRange = input.PriceRange
	}
	// Adiciona fotos ao domínio
	for _, url := range input.Photos {
		_ = p.AddPhoto(url)
	}

	if err := uc.providerRepo.Create(ctx, p); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    "RegisterProviderUseCase",
			Message: "Falha ao registrar prestador",
			Error:   err,
		})
		return nil, []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao registrar prestador",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    "RegisterProviderUseCase",
		Message: logging.DEFAULTMESSAGES.END,
	})

	return p, nil
}
