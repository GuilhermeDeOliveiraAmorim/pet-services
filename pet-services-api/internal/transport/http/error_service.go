package http

import (
	"github.com/gin-gonic/gin"
)

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
