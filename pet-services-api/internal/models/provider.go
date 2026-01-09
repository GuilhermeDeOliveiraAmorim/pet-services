package models

import (
	"time"

	"github.com/google/uuid"
)

// PriceRange define faixas de preço exibidas ao cliente.
type PriceRange string

const (
	PriceRangeLow    PriceRange = "$"
	PriceRangeMedium PriceRange = "$$"
	PriceRangeHigh   PriceRange = "$$$"
)

// ApprovalStatus controla o estado de moderação do prestador.
type ApprovalStatus string

const (
	ApprovalPending  ApprovalStatus = "pending"
	ApprovalApproved ApprovalStatus = "approved"
	ApprovalRejected ApprovalStatus = "rejected"
)

// Provider representa o prestador persistido.
type Provider struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex:idx_providers_user"`
	BusinessName   string         `gorm:"size:140;not null"`
	Description    string         `gorm:"size:1000"`
	Street         string         `gorm:"size:150"`
	Number         string         `gorm:"size:30"`
	Complement     string         `gorm:"size:150"`
	District       string         `gorm:"size:120"`
	City           string         `gorm:"size:120"`
	State          string         `gorm:"size:60"`
	ZipCode        string         `gorm:"size:20;index:idx_providers_zip"`
	Country        string         `gorm:"size:80;default:Brasil"`
	Latitude       *float64       `gorm:"type:decimal(10,7);index:idx_providers_location"`
	Longitude      *float64       `gorm:"type:decimal(10,7);index:idx_providers_location"`
	PriceRange     PriceRange     `gorm:"type:varchar(8);default:'$';index:idx_providers_price"`
	AvgRating      float64        `gorm:"type:decimal(3,2);default:0;index:idx_providers_rating"`
	TotalReviews   int            `gorm:"default:0"`
	IsActive       bool           `gorm:"default:false;index:idx_providers_active"`
	ApprovalStatus ApprovalStatus `gorm:"type:varchar(16);default:'pending';index:idx_providers_approval;index:idx_providers_active_approval,composite:idx_providers_active_approval"`
	ModerationNote string         `gorm:"size:500"`
	CreatedAt      time.Time      `gorm:"autoCreateTime;index:idx_providers_created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime"`

	User         *User                 `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:RESTRICT"`
	Services     []ProviderService     `gorm:"foreignKey:ProviderID;constraint:OnDelete:CASCADE"`
	Photos       []ProviderPhoto       `gorm:"foreignKey:ProviderID;constraint:OnDelete:CASCADE"`
	WorkingHours []ProviderWorkingHour `gorm:"foreignKey:ProviderID;constraint:OnDelete:CASCADE"`
}

// TableName define o nome da tabela no banco.
func (Provider) TableName() string {
	return "providers"
}

// ProviderService representa um serviço ofertado pelo prestador.
type ProviderService struct {
	ID         uuid.UUID              `gorm:"type:uuid;primaryKey"`
	ProviderID uuid.UUID              `gorm:"type:uuid;not null;index:idx_provider_service_provider,priority:1;uniqueIndex:idx_provider_service_unique,priority:1"`
	Category   string                 `gorm:"size:100;not null;uniqueIndex:idx_provider_service_unique,priority:2"`
	Name       string                 `gorm:"size:120;not null;uniqueIndex:idx_provider_service_unique,priority:3"`
	PriceMin   float64                `gorm:"type:numeric(10,2)"`
	PriceMax   float64                `gorm:"type:numeric(10,2)"`
	CreatedAt  time.Time              `gorm:"autoCreateTime"`
	UpdatedAt  time.Time              `gorm:"autoUpdateTime"`
	Photos     []ProviderServicePhoto `gorm:"foreignKey:ProviderServiceID;constraint:OnDelete:CASCADE"`
}

// ProviderServicePhoto representa fotos associadas a um serviço do prestador.
type ProviderServicePhoto struct {
	ID                uuid.UUID        `gorm:"type:uuid;primaryKey"`
	ProviderServiceID uuid.UUID        `gorm:"type:uuid;not null;index:idx_service_photo_service,priority:1"`
	URL               string           `gorm:"size:500;not null"`
	SortOrder         int              `gorm:"not null;default:0;index:idx_service_photo_service,priority:2"`
	CreatedAt         time.Time        `gorm:"autoCreateTime"`
	ProviderService   *ProviderService `gorm:"foreignKey:ProviderServiceID;references:ID;constraint:OnDelete:CASCADE"`
}

func (ProviderServicePhoto) TableName() string {
	return "provider_service_photos"
}

// TableName define o nome da tabela no banco.
func (ProviderService) TableName() string {
	return "provider_services"
}

// ProviderPhoto representa fotos do prestador.
type ProviderPhoto struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProviderID uuid.UUID `gorm:"type:uuid;not null;index:idx_provider_photo_provider,priority:1"`
	URL        string    `gorm:"size:500;not null"`
	SortOrder  int       `gorm:"not null;default:0;index:idx_provider_photo_provider,priority:2"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}

// TableName define o nome da tabela no banco.
func (ProviderPhoto) TableName() string {
	return "provider_photos"
}

// ProviderWorkingHour representa o horário de funcionamento por dia.
type ProviderWorkingHour struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProviderID uuid.UUID `gorm:"type:uuid;not null;index:idx_provider_working_hour_day,priority:1"`
	DayOfWeek  int       `gorm:"not null;index:idx_provider_working_hour_day,priority:2"` // 0=domingo
	IsOpen     bool      `gorm:"default:false"`
	OpenTime   string    `gorm:"size:5"`
	CloseTime  string    `gorm:"size:5"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

// TableName define o nome da tabela no banco.
func (ProviderWorkingHour) TableName() string {
	return "provider_working_hours"
}
