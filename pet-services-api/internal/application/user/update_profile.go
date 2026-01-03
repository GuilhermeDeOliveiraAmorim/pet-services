package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	domainUser "github.com/guilherme/pet-services-api/internal/domain/user"
)

// UpdateProfileUseCase atualiza informações do perfil do usuário.
type UpdateProfileUseCase struct {
	userRepo domainUser.Repository
}

func NewUpdateProfileUseCase(userRepo domainUser.Repository) *UpdateProfileUseCase {
	return &UpdateProfileUseCase{userRepo: userRepo}
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
func (uc *UpdateProfileUseCase) Execute(ctx context.Context, input UpdateProfileInput) error {
	if input.UserID == uuid.Nil {
		return domainUser.ErrUserNotFound
	}

	user, err := uc.userRepo.FindByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return fmt.Errorf("nome não pode ser vazio")
		}
		user.Name = name
	}

	if input.Phone != nil {
		phoneStr := strings.TrimSpace(*input.Phone)
		if phoneStr != "" {
			phone, err := domainUser.NewPhone(phoneStr)
			if err != nil {
				return err
			}
			user.Phone = phone
		}
	}

	locationUpdate := input.Latitude != nil || input.Longitude != nil || input.Address != nil
	if locationUpdate {
		if input.Latitude == nil || input.Longitude == nil || input.Address == nil {
			return fmt.Errorf("para atualizar localização, latitude, longitude e endereço são obrigatórios")
		}

		lat := *input.Latitude
		lon := *input.Longitude
		addr := *input.Address

		if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
			return fmt.Errorf("coordenadas inválidas")
		}

		user.SetLocation(lat, lon, addr)
	}

	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("falha ao atualizar perfil: %w", err)
	}

	return nil
}
