package entities

import "pet-services-api/internal/exceptions"

type Review struct {
	Base
	UserID     string  `json:"user_id"`
	ProviderID string  `json:"provider_id"`
	Rating     float64 `json:"rating"`
	Comment    string  `json:"comment"`
}

func NewReview(userID string, providerID string, rating float64, comment string) (*Review, []exceptions.ProblemDetails) {
	var problems []exceptions.ProblemDetails

	if userID == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do usuário ausente",
			Detail: "O ID do usuário é obrigatório",
		}))
	}

	if providerID == "" {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do provedor ausente",
			Detail: "O ID do provedor é obrigatório",
		}))
	}

	if rating < 1 || rating > 5 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Avaliação inválida",
			Detail: "A avaliação deve estar entre 1 e 5",
		}))
	}

	if len(comment) > 500 {
		problems = append(problems, exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Comentário muito longo",
			Detail: "O comentário deve ter no máximo 500 caracteres",
		}))
	}

	if len(problems) > 0 {
		return nil, problems
	}

	return &Review{
		Base:       *NewBase(),
		UserID:     userID,
		ProviderID: providerID,
		Rating:     rating,
		Comment:    comment,
	}, nil
}
