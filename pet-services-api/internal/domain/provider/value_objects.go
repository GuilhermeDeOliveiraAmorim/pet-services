package provider

import (
	"time"

	"github.com/google/uuid"
)

// Service representa um serviço oferecido pelo prestador
type Service struct {
	Category string
	Name     string
	PriceMin float64
	PriceMax float64
}

// Photo representa uma foto do prestador
type Photo struct {
	ID        uuid.UUID
	URL       string
	Order     int
	CreatedAt time.Time
}

// WorkingHours representa o horário de funcionamento
type WorkingHours struct {
	Monday    DaySchedule
	Tuesday   DaySchedule
	Wednesday DaySchedule
	Thursday  DaySchedule
	Friday    DaySchedule
	Saturday  DaySchedule
	Sunday    DaySchedule
}

// DaySchedule representa o horário de um dia específico
type DaySchedule struct {
	IsOpen bool
	Open   string // Formato: "08:00"
	Close  string // Formato: "18:00"
}

// NewDefaultWorkingHours cria horário padrão (Seg-Sex 8h-18h)
func NewDefaultWorkingHours() WorkingHours {
	weekdaySchedule := DaySchedule{
		IsOpen: true,
		Open:   "08:00",
		Close:  "18:00",
	}

	weekendSchedule := DaySchedule{
		IsOpen: false,
	}

	return WorkingHours{
		Monday:    weekdaySchedule,
		Tuesday:   weekdaySchedule,
		Wednesday: weekdaySchedule,
		Thursday:  weekdaySchedule,
		Friday:    weekdaySchedule,
		Saturday:  weekendSchedule,
		Sunday:    weekendSchedule,
	}
}
