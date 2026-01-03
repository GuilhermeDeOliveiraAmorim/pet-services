package provider

import (
	"time"

	"github.com/google/uuid"
	"github.com/guilherme/pet-services-api/internal/domain/user"
)

// PriceRange representa a faixa de preço dos serviços
type PriceRange string

const (
	PriceRangeLow    PriceRange = "$"
	PriceRangeMedium PriceRange = "$$"
	PriceRangeHigh   PriceRange = "$$$"
)

// Provider representa um prestador de serviços
type Provider struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	User         *user.User
	BusinessName string
	Description  string
	Address      user.Address
	Location     user.Location
	Services     []Service
	Photos       []Photo
	WorkingHours WorkingHours
	PriceRange   PriceRange
	AvgRating    float64
	TotalReviews int
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewProvider cria um novo prestador
func NewProvider(userID uuid.UUID, businessName, description string) *Provider {
	now := time.Now()
	return &Provider{
		ID:           uuid.New(),
		UserID:       userID,
		BusinessName: businessName,
		Description:  description,
		Services:     []Service{},
		Photos:       []Photo{},
		WorkingHours: NewDefaultWorkingHours(),
		IsActive:     false, // Requer aprovação
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// AddService adiciona um serviço ao prestador
func (p *Provider) AddService(category, name string, priceMin, priceMax float64) {
	service := Service{
		Category: category,
		Name:     name,
		PriceMin: priceMin,
		PriceMax: priceMax,
	}
	p.Services = append(p.Services, service)
	p.UpdatedAt = time.Now()
}

// RemoveService remove um serviço por categoria e nome
func (p *Provider) RemoveService(category, name string) {
	filtered := []Service{}
	for _, s := range p.Services {
		if s.Category != category || s.Name != name {
			filtered = append(filtered, s)
		}
	}
	p.Services = filtered
	p.UpdatedAt = time.Now()
}

// AddPhoto adiciona uma foto ao prestador
func (p *Provider) AddPhoto(url string) error {
	if url == "" {
		return NewValidationError("photo.url", "url da foto é obrigatória")
	}

	if len(p.Photos) >= 10 {
		return ErrMaxPhotosReached
	}

	photo := Photo{
		ID:        uuid.New(),
		URL:       url,
		Order:     len(p.Photos),
		CreatedAt: time.Now(),
	}
	p.Photos = append(p.Photos, photo)
	p.UpdatedAt = time.Now()
	return nil
}

// RemovePhoto remove uma foto por ID

func (p *Provider) RemovePhoto(photoID uuid.UUID) error {
	filtered := []Photo{}
	found := false
	for _, photo := range p.Photos {
		if photo.ID != photoID {
			filtered = append(filtered, photo)
			continue
		}
		found = true
	}
	p.Photos = filtered
	p.reorderPhotos()
	p.UpdatedAt = time.Now()

	if !found {
		return ErrPhotoNotFound
	}

	return nil
}

func (p *Provider) reorderPhotos() {
	for i := range p.Photos {
		p.Photos[i].Order = i
	}
}

// ReorderPhotos ajusta a ordem das fotos pelo array de IDs
func (p *Provider) ReorderPhotos(order []uuid.UUID) error {
	if len(order) != len(p.Photos) {
		return NewValidationError("photos.order", "lista de ordenação não corresponde ao número de fotos")
	}

	index := make(map[uuid.UUID]Photo, len(p.Photos))
	for _, photo := range p.Photos {
		index[photo.ID] = photo
	}

	reordered := make([]Photo, 0, len(p.Photos))
	for i, id := range order {
		photo, ok := index[id]
		if !ok {
			return ErrPhotoNotFound
		}
		photo.Order = i
		reordered = append(reordered, photo)
	}

	p.Photos = reordered
	p.UpdatedAt = time.Now()
	return nil
}

// SetLocation define a localização do prestador
func (p *Provider) SetLocation(latitude, longitude float64, address user.Address) {
	p.Location = user.Location{
		Latitude:  latitude,
		Longitude: longitude,
		Address:   address,
	}
	p.Address = address
	p.UpdatedAt = time.Now()
}

// SetWorkingHours define o horário de funcionamento com validação
func (p *Provider) SetWorkingHours(hours WorkingHours) error {
	if err := hours.Validate(); err != nil {
		return err
	}

	p.WorkingHours = hours
	p.UpdatedAt = time.Now()
	return nil
}

// SetDaySchedule atualiza o horário de um dia específico
func (p *Provider) SetDaySchedule(day time.Weekday, schedule DaySchedule) error {
	if err := schedule.Validate(day.String()); err != nil {
		return err
	}

	switch day {
	case time.Monday:
		p.WorkingHours.Monday = schedule
	case time.Tuesday:
		p.WorkingHours.Tuesday = schedule
	case time.Wednesday:
		p.WorkingHours.Wednesday = schedule
	case time.Thursday:
		p.WorkingHours.Thursday = schedule
	case time.Friday:
		p.WorkingHours.Friday = schedule
	case time.Saturday:
		p.WorkingHours.Saturday = schedule
	case time.Sunday:
		p.WorkingHours.Sunday = schedule
	default:
		return ErrInvalidWorkingHours
	}

	p.UpdatedAt = time.Now()
	return nil
}

// Activate ativa o prestador
func (p *Provider) Activate() {
	p.IsActive = true
	p.UpdatedAt = time.Now()
}

// Deactivate desativa o prestador
func (p *Provider) Deactivate() {
	p.IsActive = false
	p.UpdatedAt = time.Now()
}

// UpdateRating atualiza a avaliação média
func (p *Provider) UpdateRating(newRating float64) {
	totalRating := p.AvgRating * float64(p.TotalReviews)
	p.TotalReviews++
	p.AvgRating = (totalRating + newRating) / float64(p.TotalReviews)
	p.UpdatedAt = time.Now()
}
