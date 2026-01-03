package review

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Review representa uma avaliação de um serviço
type Review struct {
	ID         uuid.UUID
	RequestID  uuid.UUID
	ProviderID uuid.UUID
	OwnerID    uuid.UUID
	Rating     int
	Comment    string
	CreatedAt  time.Time
}

// NewReview cria uma nova avaliação
func NewReview(requestID, providerID, ownerID uuid.UUID, rating int, comment string) (*Review, error) {
	if rating < 1 || rating > 5 {
		return nil, ErrInvalidRating
	}

	return &Review{
		ID:         uuid.New(),
		RequestID:  requestID,
		ProviderID: providerID,
		OwnerID:    ownerID,
		Rating:     rating,
		Comment:    comment,
		CreatedAt:  time.Now(),
	}, nil
}

// Validate valida a avaliação
func (r *Review) Validate() error {
	if r.Rating < 1 || r.Rating > 5 {
		return fmt.Errorf("%w: nota deve estar entre 1 e 5", ErrInvalidRating)
	}
	return nil
}
