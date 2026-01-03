package provider

import (
	"errors"
	"fmt"
)

// Erros de domínio do Provider
var (
	// ErrProviderNotFound indica que o prestador não foi encontrado
	ErrProviderNotFound = errors.New("prestador não encontrado")

	// ErrProviderNotActive indica que o prestador não está ativo
	ErrProviderNotActive = errors.New("prestador não está ativo")

	// ErrProviderAlreadyExists indica que já existe um prestador para o usuário
	ErrProviderAlreadyExists = errors.New("prestador já existe para este usuário")

	// ErrMaxPhotosReached indica que o limite de fotos foi atingido
	ErrMaxPhotosReached = errors.New("número máximo de fotos atingido (10)")

	// ErrPhotoNotFound indica que a foto não foi encontrada
	ErrPhotoNotFound = errors.New("foto não encontrada")

	// ErrInvalidLocation indica que a localização é inválida
	ErrInvalidLocation = errors.New("localização inválida")

	// ErrInvalidPriceRange indica que a faixa de preço é inválida
	ErrInvalidPriceRange = errors.New("faixa de preço inválida")

	// ErrInvalidWorkingHours indica que o horário de funcionamento é inválido
	ErrInvalidWorkingHours = errors.New("horário de funcionamento inválido")

	// ErrNoServicesProvided indica que nenhum serviço foi cadastrado
	ErrNoServicesProvided = errors.New("nenhum serviço cadastrado")

	// ErrServiceNotFound indica que o serviço não foi encontrado
	ErrServiceNotFound = errors.New("serviço não encontrado")

	// ErrInvalidBusinessName indica que o nome do negócio é inválido
	ErrInvalidBusinessName = errors.New("nome do negócio inválido")

	// ErrInvalidRating indica que a avaliação é inválida
	ErrInvalidRating = errors.New("avaliação deve estar entre 0 e 5")
)

// ValidationError representa um erro de validação com contexto
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("validação falhou no campo '%s': %s", e.Field, e.Message)
}

// NewValidationError cria um novo erro de validação
func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}

// LocationError representa um erro relacionado à localização
type LocationError struct {
	Latitude  float64
	Longitude float64
	Message   string
}

func (e LocationError) Error() string {
	return fmt.Sprintf("erro de localização (lat: %.6f, lon: %.6f): %s",
		e.Latitude, e.Longitude, e.Message)
}

// NewLocationError cria um novo erro de localização
func NewLocationError(lat, lon float64, message string) LocationError {
	return LocationError{
		Latitude:  lat,
		Longitude: lon,
		Message:   message,
	}
}