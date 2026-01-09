package gormrepo

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	authdom "pet-services-api/internal/domain/auth"
	providerdom "pet-services-api/internal/domain/provider"
	requestdom "pet-services-api/internal/domain/request"
	reviewdom "pet-services-api/internal/domain/review"
	userdom "pet-services-api/internal/domain/user"
	"pet-services-api/internal/models"
)

// toModelUser converte domínio para model.
func toModelUser(u *userdom.User) (*models.User, error) {
	if u == nil {
		return nil, fmt.Errorf("user is nil")
	}

	var lat, lon *float64
	if u.Location != nil {
		lat = &u.Location.Latitude
		lon = &u.Location.Longitude
	}

	addr := userdom.Address{}
	if u.Location != nil {
		addr = u.Location.Address
	}

	return &models.User{
		ID:            u.ID,
		Email:         u.Email.String(),
		EmailVerified: u.EmailVerified,
		Password:      u.Password,
		Name:          u.Name,
		Phone:         u.Phone.String(),
		Type:          models.UserType(u.Type),
		Latitude:      lat,
		Longitude:     lon,
		Street:        addr.Street,
		Number:        addr.Number,
		Complement:    addr.Complement,
		District:      addr.District,
		City:          addr.City,
		State:         addr.State,
		ZipCode:       addr.ZipCode,
		Country:       addr.Country,
		DeletedAt:     u.DeletedAt,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}, nil
}

// toDomainUser converte model para domínio.
func toDomainUser(m *models.User) (*userdom.User, error) {
	if m == nil {
		return nil, fmt.Errorf("model user is nil")
	}

	email, err := userdom.NewEmail(m.Email)
	if err != nil {
		return nil, err
	}

	phone, err := userdom.NewPhone(m.Phone)
	if err != nil {
		return nil, err
	}

	var location *userdom.Location
	if m.Latitude != nil && m.Longitude != nil {
		addr := userdom.Address{
			Street:     m.Street,
			Number:     m.Number,
			Complement: m.Complement,
			District:   m.District,
			City:       m.City,
			State:      m.State,
			ZipCode:    m.ZipCode,
			Country:    m.Country,
		}
		loc, err := userdom.NewLocation(*m.Latitude, *m.Longitude, addr)
		if err != nil {
			return nil, err
		}
		location = &loc
	}

	return &userdom.User{
		ID:            m.ID,
		Email:         email,
		EmailVerified: m.EmailVerified,
		Password:      m.Password,
		Name:          m.Name,
		Phone:         phone,
		Type:          userdom.UserType(m.Type),
		Location:      location,
		DeletedAt:     m.DeletedAt,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}, nil
}

func toModelRefreshToken(t *authdom.RefreshToken) *models.RefreshToken {
	if t == nil {
		return nil
	}
	return &models.RefreshToken{
		ID:        t.ID,
		UserID:    t.UserID,
		ExpiresAt: t.ExpiresAt,
		Revoked:   t.Revoked,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func toDomainRefreshToken(m *models.RefreshToken) *authdom.RefreshToken {
	if m == nil {
		return nil
	}
	return &authdom.RefreshToken{
		ID:        m.ID,
		UserID:    m.UserID,
		ExpiresAt: m.ExpiresAt,
		Revoked:   m.Revoked,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func toModelPasswordResetToken(t *userdom.PasswordResetToken) *models.PasswordResetToken {
	if t == nil {
		return nil
	}
	return &models.PasswordResetToken{
		ID:        t.ID,
		UserID:    t.UserID,
		Token:     t.Token,
		ExpiresAt: t.ExpiresAt,
		Used:      t.Used,
		CreatedAt: t.CreatedAt,
	}
}

func toDomainPasswordResetToken(m *models.PasswordResetToken) *userdom.PasswordResetToken {
	if m == nil {
		return nil
	}
	return &userdom.PasswordResetToken{
		ID:        m.ID,
		UserID:    m.UserID,
		Token:     m.Token,
		ExpiresAt: m.ExpiresAt,
		Used:      m.Used,
		CreatedAt: m.CreatedAt,
	}
}

func toModelEmailVerificationToken(t *userdom.EmailVerificationToken) *models.EmailVerificationToken {
	if t == nil {
		return nil
	}
	return &models.EmailVerificationToken{
		ID:        t.ID,
		UserID:    t.UserID,
		Token:     t.Token,
		ExpiresAt: t.ExpiresAt,
		Used:      t.Used,
		CreatedAt: t.CreatedAt,
	}
}

func toDomainEmailVerificationToken(m *models.EmailVerificationToken) *userdom.EmailVerificationToken {
	if m == nil {
		return nil
	}
	return &userdom.EmailVerificationToken{
		ID:        m.ID,
		UserID:    m.UserID,
		Token:     m.Token,
		ExpiresAt: m.ExpiresAt,
		Used:      m.Used,
		CreatedAt: m.CreatedAt,
	}
}

func toModelProvider(p *providerdom.Provider) (*models.Provider, error) {
	if p == nil {
		return nil, fmt.Errorf("provider is nil")
	}

	var lat, lon *float64
	if p.Location.Latitude != 0 || p.Location.Longitude != 0 {
		lat = &p.Location.Latitude
		lon = &p.Location.Longitude
	}

	m := &models.Provider{
		ID:             p.ID,
		UserID:         p.UserID,
		BusinessName:   p.BusinessName,
		Description:    p.Description,
		Street:         p.Address.Street,
		Number:         p.Address.Number,
		Complement:     p.Address.Complement,
		District:       p.Address.District,
		City:           p.Address.City,
		State:          p.Address.State,
		ZipCode:        p.Address.ZipCode,
		Country:        p.Address.Country,
		Latitude:       lat,
		Longitude:      lon,
		PriceRange:     models.PriceRange(p.PriceRange),
		AvgRating:      p.AvgRating,
		TotalReviews:   p.TotalReviews,
		IsActive:       p.IsActive,
		ApprovalStatus: models.ApprovalStatus(p.ApprovalStatus),
		ModerationNote: p.ModerationNote,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
	}

	// Services
	for _, s := range p.Services {
		m.Services = append(m.Services, models.ProviderService{
			ID:         uuid.New(),
			ProviderID: p.ID,
			Category:   s.Category,
			Name:       s.Name,
			PriceMin:   s.PriceMin,
			PriceMax:   s.PriceMax,
		})
	}

	// Photos
	for _, ph := range p.Photos {
		m.Photos = append(m.Photos, models.ProviderPhoto{
			ID:         ph.ID,
			ProviderID: p.ID,
			URL:        ph.URL,
			SortOrder:  ph.Order,
			CreatedAt:  ph.CreatedAt,
		})
	}

	// Working hours
	schedules := map[int]providerdom.DaySchedule{
		0: p.WorkingHours.Sunday,
		1: p.WorkingHours.Monday,
		2: p.WorkingHours.Tuesday,
		3: p.WorkingHours.Wednesday,
		4: p.WorkingHours.Thursday,
		5: p.WorkingHours.Friday,
		6: p.WorkingHours.Saturday,
	}
	for dow, sch := range schedules {
		m.WorkingHours = append(m.WorkingHours, models.ProviderWorkingHour{
			ID:         uuid.New(),
			ProviderID: p.ID,
			DayOfWeek:  dow,
			IsOpen:     sch.IsOpen,
			OpenTime:   sch.Open,
			CloseTime:  sch.Close,
		})
	}

	return m, nil
}

func toDomainProvider(m *models.Provider) (*providerdom.Provider, error) {
	if m == nil {
		return nil, fmt.Errorf("provider model is nil")
	}

	p := providerdom.NewProvider(m.UserID, m.BusinessName, m.Description)
	p.ID = m.ID
	p.Address = userdom.Address{
		Street:     m.Street,
		Number:     m.Number,
		Complement: m.Complement,
		District:   m.District,
		City:       m.City,
		State:      m.State,
		ZipCode:    m.ZipCode,
		Country:    m.Country,
	}
	if m.Latitude != nil && m.Longitude != nil {
		p.Location = userdom.Location{Latitude: *m.Latitude, Longitude: *m.Longitude, Address: p.Address}
	}
	p.PriceRange = providerdom.PriceRange(m.PriceRange)
	p.AvgRating = m.AvgRating
	p.TotalReviews = m.TotalReviews
	p.IsActive = m.IsActive
	p.ApprovalStatus = providerdom.ApprovalStatus(m.ApprovalStatus)
	p.ModerationNote = m.ModerationNote
	p.CreatedAt = m.CreatedAt
	p.UpdatedAt = m.UpdatedAt

	// Services
	for _, s := range m.Services {
		p.Services = append(p.Services, providerdom.Service{
			Category: s.Category,
			Name:     s.Name,
			PriceMin: s.PriceMin,
			PriceMax: s.PriceMax,
		})
	}

	// Photos
	for _, ph := range m.Photos {
		p.Photos = append(p.Photos, providerdom.Photo{
			ID:        ph.ID,
			URL:       ph.URL,
			Order:     ph.SortOrder,
			CreatedAt: ph.CreatedAt,
		})
	}

	// Working hours
	hours := providerdom.WorkingHours{}
	for _, wh := range m.WorkingHours {
		ds := providerdom.DaySchedule{IsOpen: wh.IsOpen, Open: wh.OpenTime, Close: wh.CloseTime}
		switch wh.DayOfWeek {
		case 1:
			hours.Monday = ds
		case 2:
			hours.Tuesday = ds
		case 3:
			hours.Wednesday = ds
		case 4:
			hours.Thursday = ds
		case 5:
			hours.Friday = ds
		case 6:
			hours.Saturday = ds
		case 0:
			hours.Sunday = ds
		}
	}
	p.WorkingHours = hours

	return p, nil
}

func toModelServiceRequest(r *requestdom.ServiceRequest) *models.ServiceRequest {
	if r == nil {
		return nil
	}
	return &models.ServiceRequest{
		ID:              r.ID,
		OwnerID:         r.OwnerID,
		ProviderID:      r.ProviderID,
		ServiceType:     r.ServiceType,
		PetName:         r.Pet.Name,
		PetType:         models.PetType(r.Pet.Type),
		PetBreed:        r.Pet.Breed,
		PetAge:          r.Pet.Age,
		PetWeight:       r.Pet.Weight,
		PetNotes:        r.Pet.Notes,
		PreferredDate:   r.PreferredDate,
		PreferredTime:   r.PreferredTime,
		AdditionalNotes: r.AdditionalNotes,
		Status:          models.RequestStatus(r.Status),
		RejectionReason: r.RejectionReason,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

func toDomainServiceRequest(m *models.ServiceRequest) (*requestdom.ServiceRequest, error) {
	if m == nil {
		return nil, fmt.Errorf("service request model is nil")
	}

	pet := requestdom.PetInfo{
		Name:   m.PetName,
		Type:   requestdom.PetType(m.PetType),
		Breed:  m.PetBreed,
		Age:    m.PetAge,
		Weight: m.PetWeight,
		Notes:  m.PetNotes,
	}
	if err := pet.IsValid(); err != nil {
		return nil, err
	}

	req := requestdom.NewServiceRequest(
		m.OwnerID,
		m.ProviderID,
		m.ServiceType,
		pet,
		m.PreferredDate,
		m.PreferredTime,
		m.AdditionalNotes,
	)
	req.ID = m.ID
	req.Status = requestdom.Status(m.Status)
	req.RejectionReason = m.RejectionReason
	req.CreatedAt = m.CreatedAt
	req.UpdatedAt = m.UpdatedAt
	return req, nil
}

func toModelReview(r *reviewdom.Review) *models.Review {
	if r == nil {
		return nil
	}
	return &models.Review{
		ID:         r.ID,
		RequestID:  r.RequestID,
		ProviderID: r.ProviderID,
		OwnerID:    r.OwnerID,
		Rating:     r.Rating,
		Comment:    r.Comment,
		CreatedAt:  r.CreatedAt,
	}
}

func toDomainReview(m *models.Review) *reviewdom.Review {
	if m == nil {
		return nil
	}
	return &reviewdom.Review{
		ID:         m.ID,
		RequestID:  m.RequestID,
		ProviderID: m.ProviderID,
		OwnerID:    m.OwnerID,
		Rating:     m.Rating,
		Comment:    m.Comment,
		CreatedAt:  m.CreatedAt,
	}
}

// applyPagination aplica limit/offset padrão.
func applyPagination(tx *gorm.DB, page, limit int) *gorm.DB {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	return tx.Limit(limit).Offset(offset)
}
