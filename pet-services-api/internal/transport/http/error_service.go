package http

import (
	"github.com/gin-gonic/gin"
	"pet-services-api/internal/application/exceptions"
)

// ProblemDetailsDTO representa um erro detalhado para resposta HTTP.
type ProblemDetailsDTO struct {
	Type   string `json:"type,omitempty"`
	Title  string `json:"title"`
	Status int    `json:"status"`
	Detail string `json:"detail,omitempty"`
}

// RespondWithProblems responde com um ou mais ProblemDetails.
func (s *ErrorService) RespondWithProblems(c *gin.Context, problems []exceptions.ProblemDetails, defaultCode string, defaultStatus int) {
	if len(problems) == 0 {
		c.JSON(defaultStatus, ErrorDTO{Message: "Erro desconhecido", Code: defaultCode})
		return
	}
	// Se só há um problema, responde direto
	if len(problems) == 1 {
		p := problems[0]
		c.JSON(p.Status, ProblemDetailsDTO{
			Type:   p.Type,
			Title:  p.Title,
			Status: p.Status,
			Detail: p.Detail,
		})
		return
	}
	// Caso múltiplos problemas, retorna array
	dtos := make([]ProblemDetailsDTO, 0, len(problems))
	for _, p := range problems {
		dtos = append(dtos, ProblemDetailsDTO{
			Type:   p.Type,
			Title:  p.Title,
			Status: p.Status,
			Detail: p.Detail,
		})
	}
	c.JSON(defaultStatus, gin.H{"errors": dtos})
}

type ErrorDTO struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

type ErrorService struct{}

func NewErrorService() *ErrorService {
	return &ErrorService{}
}

func (s *ErrorService) ToDTO(err error, code string) ErrorDTO {
	return ErrorDTO{
		Message: err.Error(),
		Code:    code,
	}
}

func (s *ErrorService) RespondWithError(c *gin.Context, err error, code string, status int) {
	c.JSON(status, s.ToDTO(err, code))
}
