package entities

import (
	"pet-services-api/internal/exceptions"
	"time"
)

type Provider struct {
	Base
	UserID        string    `json:"user_id"`
	BusinessName  string    `json:"business_name"`
	Address       Address   `json:"address"`
	Description   string    `json:"description"`
	PriceRange    string    `json:"price_range"`
	AverageRating float64   `json:"average_rating"`
	Photos        []Photo   `json:"photos"`
	Reviews       []Review  `json:"reviews"`
	Requests      []Request `json:"requests"`
}

func NewProvider(userID string, businessName string, address Address, description string, priceRange string) (*Provider, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if userID == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do usuário ausente",
			Detail: "O ID do usuário é obrigatório",
		}))
	}

	if businessName == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome comercial ausente",
			Detail: "O nome comercial é obrigatório",
		}))
	} else if len(businessName) > 100 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Nome comercial muito longo",
			Detail: "O nome comercial deve ter no máximo 100 caracteres",
		}))
	}

	if len(description) > 1000 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Descrição muito longa",
			Detail: "A descrição deve ter no máximo 1000 caracteres",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &Provider{
		Base:         *NewBase(),
		UserID:       userID,
		BusinessName: businessName,
		Address:      address,
		Description:  description,
		PriceRange:   priceRange,
	}, nil
}

func (p *Provider) AddPhoto(photo Photo) {
	timeNow := time.Now()
	p.Photos = append(p.Photos, photo)
	p.UpdatedAt = &timeNow
}

func (p *Provider) RemovePhoto(photo Photo) {
	timeNow := time.Now()
	var updatedPhotos []Photo
	for _, pht := range p.Photos {
		if pht.ID != photo.ID {
			updatedPhotos = append(updatedPhotos, pht)
		}
	}
	p.Photos = updatedPhotos
	p.UpdatedAt = &timeNow
}

func (p *Provider) AddReview(review Review) {
	timeNow := time.Now()
	p.Reviews = append(p.Reviews, review)
	p.UpdateAverageRating()
	p.UpdatedAt = &timeNow
}

func (p *Provider) UpdateAverageRating() {
	if len(p.Reviews) == 0 {
		p.AverageRating = 0
		return
	}

	var total float64
	for _, review := range p.Reviews {
		total += review.Rating
	}

	p.AverageRating = total / float64(len(p.Reviews))
}

func (p *Provider) AddRequest(request Request) {
	timeNow := time.Now()
	p.Requests = append(p.Requests, request)
	p.UpdatedAt = &timeNow
}
