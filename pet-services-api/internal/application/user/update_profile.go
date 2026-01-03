package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/application/exceptions"
	"github.com/guilherme/pet-services-api/internal/application/logging"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// UpdateProfileUseCase atualiza informações do perfil do usuário.
type UpdateProfileUseCase struct {
	userRepo domainUser.Repository
	logger   logging.LoggerService
}

func NewUpdateProfileUseCase(userRepo domainUser.Repository, logger logging.LoggerService) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{userRepo: userRepo, logger: logger}
}

// UpdateProfileInput campos opcionais para atualização.
type UpdateProfileInput struct {
	UserID    uuid.UUID
	Name      *string
	Phone     *string
	Address   *domainUser.Address
	Latitude  *float64
	Longitude *float64
}

// Execute valida e atualiza o perfil do usuário.
const UPDATE_PROFILE_USECASE = "UPDATE_PROFILE_USECASE"

func (uc *UpdateProfileUseCase) Execute(ctx context.Context, input UpdateProfileInput) []exceptions.ProblemDetails {
	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    UPDATE_PROFILE_USECASE,
		Message: logging.DEFAULTMESSAGES.START,
	})

	if input.UserID == uuid.Nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    UPDATE_PROFILE_USECASE,
			Message: "Usuário não encontrado",
			Error:   domainUser.ErrUserNotFound,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Usuário não encontrado",
			Status: exceptions.RFC404_CODE,
			Detail: "O usuário informado não foi encontrado.",
		}}
	}

	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC404_CODE,
			From:    UPDATE_PROFILE_USECASE,
			Message: "Usuário não encontrado",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC404,
			Title:  "Usuário não encontrado",
			Status: exceptions.RFC404_CODE,
			Detail: "O usuário informado não foi encontrado.",
		}}
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			uc.logger.Log(logging.Logger{
				Context: ctx,
				TypeLog: logging.LoggerTypes.ERROR,
				Layer:   logging.LoggerLayers.USECASES,
				Code:    exceptions.RFC400_CODE,
				From:    UPDATE_PROFILE_USECASE,
				Message: "Nome não pode ser vazio",
				Error:   errors.New("nome não pode ser vazio"),
			})
			return []exceptions.ProblemDetails{{
				Type:   exceptions.RFC400,
				Title:  "Nome não pode ser vazio",
				Status: exceptions.RFC400_CODE,
				Detail: "O nome não pode ser vazio.",
			}}
		}
		user.Name = name
	}

	if input.Phone != nil {
		phoneStr := strings.TrimSpace(*input.Phone)
		if phoneStr != "" {
			phone, err := domainUser.NewPhone(phoneStr)
			if err != nil {
				uc.logger.Log(logging.Logger{
					Context: ctx,
					TypeLog: logging.LoggerTypes.ERROR,
					Layer:   logging.LoggerLayers.USECASES,
					Code:    exceptions.RFC400_CODE,
					From:    UPDATE_PROFILE_USECASE,
					Message: "Telefone inválido",
					Error:   err,
				})
				return []exceptions.ProblemDetails{{
					Type:   exceptions.RFC400,
					Title:  "Telefone inválido",
					Status: exceptions.RFC400_CODE,
					Detail: err.Error(),
				}}
			}
			user.Phone = phone
		}
	}

	locationUpdate := input.Latitude != nil || input.Longitude != nil || input.Address != nil
	if locationUpdate {
		if input.Latitude == nil || input.Longitude == nil || input.Address == nil {
			uc.logger.Log(logging.Logger{
				Context: ctx,
				TypeLog: logging.LoggerTypes.ERROR,
				Layer:   logging.LoggerLayers.USECASES,
				Code:    exceptions.RFC400_CODE,
				From:    UPDATE_PROFILE_USECASE,
				Message: "Para atualizar localização, latitude, longitude e endereço são obrigatórios",
				Error:   errors.New("para atualizar localização, latitude, longitude e endereço são obrigatórios"),
			})
			return []exceptions.ProblemDetails{{
				Type:   exceptions.RFC400,
				Title:  "Dados de localização incompletos",
				Status: exceptions.RFC400_CODE,
				Detail: "Para atualizar localização, latitude, longitude e endereço são obrigatórios.",
			}}
		}

		lat := *input.Latitude
		lon := *input.Longitude
		addr := *input.Address

		if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
			uc.logger.Log(logging.Logger{
				Context: ctx,
				TypeLog: logging.LoggerTypes.ERROR,
				Layer:   logging.LoggerLayers.USECASES,
				Code:    exceptions.RFC400_CODE,
				From:    UPDATE_PROFILE_USECASE,
				Message: "Coordenadas inválidas",
				Error:   errors.New("coordenadas inválidas"),
			})
			return []exceptions.ProblemDetails{{
				Type:   exceptions.RFC400,
				Title:  "Coordenadas inválidas",
				Status: exceptions.RFC400_CODE,
				Detail: "Latitude ou longitude inválidas.",
			}}
		}

		user.SetLocation(lat, lon, addr)
	}

	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		uc.logger.Log(logging.Logger{
			Context: ctx,
			TypeLog: logging.LoggerTypes.ERROR,
			Layer:   logging.LoggerLayers.USECASES,
			Code:    exceptions.RFC500_CODE,
			From:    UPDATE_PROFILE_USECASE,
			Message: "Falha ao atualizar perfil",
			Error:   err,
		})
		return []exceptions.ProblemDetails{{
			Type:   exceptions.RFC500,
			Title:  "Falha ao atualizar perfil",
			Status: exceptions.RFC500_CODE,
			Detail: err.Error(),
		}}
	}

	uc.logger.Log(logging.Logger{
		Context: ctx,
		TypeLog: logging.LoggerTypes.INFO,
		Layer:   logging.LoggerLayers.USECASES,
		Code:    exceptions.RFC200_CODE,
		From:    UPDATE_PROFILE_USECASE,
		Message: logging.DEFAULTMESSAGES.END,
	})

	return nil
}
