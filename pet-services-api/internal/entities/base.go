package entities

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Base struct {
	ID            string     `json:"id"`
	Active        bool       `json:"active"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	DeactivatedAt *time.Time `json:"deactivated_at"`
}

func NewBase() *Base {
	timeNow := time.Now()

	return &Base{
		ID:            ulid.Make().String(),
		Active:        true,
		CreatedAt:     timeNow,
		UpdatedAt:     nil,
		DeactivatedAt: nil,
	}
}

func (se *Base) Activate() {
	timeNow := time.Now()
	se.DeactivatedAt = nil
	se.UpdatedAt = &timeNow
	se.Active = true
}

func (se *Base) Deactivate() {
	timeNow := time.Now()
	se.DeactivatedAt = &timeNow
	se.UpdatedAt = &timeNow
	se.Active = false
}
