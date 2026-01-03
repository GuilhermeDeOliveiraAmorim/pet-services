package request

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Status representa o status de uma solicitação
type Status string

const (
	StatusPending   Status = "pending"
	StatusAccepted  Status = "accepted"
	StatusRejected  Status = "rejected"
	StatusCompleted Status = "completed"
	StatusCancelled Status = "cancelled"
)

// ServiceRequest representa uma solicitação de serviço
type ServiceRequest struct {
	ID              uuid.UUID
	OwnerID         uuid.UUID
	ProviderID      uuid.UUID
	ServiceType     string
	Pet             PetInfo
	PreferredDate   time.Time
	PreferredTime   string
	AdditionalNotes string
	Status          Status
	RejectionReason string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// NewServiceRequest cria uma nova solicitação
func NewServiceRequest(
	ownerID, providerID uuid.UUID,
	serviceType string,
	pet PetInfo,
	preferredDate time.Time,
	preferredTime string,
	notes string,
) *ServiceRequest {
	now := time.Now()
	return &ServiceRequest{
		ID:              uuid.New(),
		OwnerID:         ownerID,
		ProviderID:      providerID,
		ServiceType:     serviceType,
		Pet:             pet,
		PreferredDate:   preferredDate,
		PreferredTime:   preferredTime,
		AdditionalNotes: notes,
		Status:          StatusPending,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// Accept aceita a solicitação
func (r *ServiceRequest) Accept() error {
	if r.Status != StatusPending {
		return fmt.Errorf("%w: não pode aceitar solicitação com status %s", ErrInvalidStatusTransition, r.Status)
	}
	r.Status = StatusAccepted
	r.UpdatedAt = time.Now()
	return nil
}

// Reject rejeita a solicitação
func (r *ServiceRequest) Reject(reason string) error {
	if r.Status != StatusPending {
		return fmt.Errorf("%w: não pode rejeitar solicitação com status %s", ErrInvalidStatusTransition, r.Status)
	}
	r.Status = StatusRejected
	r.RejectionReason = reason
	r.UpdatedAt = time.Now()
	return nil
}

// Complete marca a solicitação como concluída
func (r *ServiceRequest) Complete() error {
	if r.Status != StatusAccepted {
		return fmt.Errorf("%w: não pode completar solicitação com status %s", ErrInvalidStatusTransition, r.Status)
	}
	r.Status = StatusCompleted
	r.UpdatedAt = time.Now()
	return nil
}

// Cancel cancela a solicitação
func (r *ServiceRequest) Cancel() error {
	if r.Status == StatusCompleted {
		return fmt.Errorf("%w: não pode cancelar solicitação concluída", ErrInvalidStatusTransition)
	}
	r.Status = StatusCancelled
	r.UpdatedAt = time.Now()
	return nil
}

// CanBeAccepted verifica se pode ser aceita
func (r *ServiceRequest) CanBeAccepted() bool {
	return r.Status == StatusPending
}

// CanBeRejected verifica se pode ser rejeitada
func (r *ServiceRequest) CanBeRejected() bool {
	return r.Status == StatusPending
}

// CanBeCompleted verifica se pode ser completada
func (r *ServiceRequest) CanBeCompleted() bool {
	return r.Status == StatusAccepted
}

// CanBeCancelled verifica se pode ser cancelada
func (r *ServiceRequest) CanBeCancelled() bool {
	return r.Status != StatusCompleted
}
