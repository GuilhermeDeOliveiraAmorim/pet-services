package provider

import (
	"strings"
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

// Validate garante que o horário do dia é válido.
func (ds DaySchedule) Validate(dayLabel string) error {
	if !ds.IsOpen {
		return nil
	}

	field := "workingHours." + strings.ToLower(dayLabel)

	if ds.Open == "" || ds.Close == "" {
		return NewValidationError(field, "horários de abertura e fechamento são obrigatórios")
	}

	openTime, err := time.Parse("15:04", ds.Open)
	if err != nil {
		return NewValidationError(field, "formato de horário inválido, use HH:MM")
	}

	closeTime, err := time.Parse("15:04", ds.Close)
	if err != nil {
		return NewValidationError(field, "formato de horário inválido, use HH:MM")
	}

	if !openTime.Before(closeTime) {
		return NewValidationError(field, "horário de abertura deve ser antes do fechamento")
	}

	return nil
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

// Validate garante que todos os dias possuem horários válidos.
func (w WorkingHours) Validate() error {
	days := map[string]DaySchedule{
		"monday":    w.Monday,
		"tuesday":   w.Tuesday,
		"wednesday": w.Wednesday,
		"thursday":  w.Thursday,
		"friday":    w.Friday,
		"saturday":  w.Saturday,
		"sunday":    w.Sunday,
	}

	for day, schedule := range days {
		if err := schedule.Validate(day); err != nil {
			return err
		}
	}

	return nil
}
