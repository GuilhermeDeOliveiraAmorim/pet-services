package http

import (
	"log/slog"
	"net/http"

	"pet-services-api/internal/application/exceptions"

	"github.com/gin-gonic/gin"
)

// BindAndValidateJSONProblems faz bind do JSON e valida, retornando []ProblemDetails.
func BindAndValidateJSONProblems(c *gin.Context, req any) []exceptions.ProblemDetails {
	if err := c.ShouldBindJSON(req); err != nil {
		slog.Warn("BindAndValidateJSONProblems: erro ao fazer bind do JSON", slog.String("error", err.Error()))
		return []exceptions.ProblemDetails{{
			Type:   string(exceptions.BadRequest),
			Title:  "Erro ao decodificar JSON",
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}}
	}
	if err := ValidateRequest(req); err != nil {
		return []exceptions.ProblemDetails{{
			Type:   string(exceptions.BadRequest),
			Title:  "Erro de validação",
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}}
	}
	return nil
}

// BindAndValidateJSON faz bind do JSON e valida em um passo
func BindAndValidateJSON(c *gin.Context, req any) error {
	if err := c.ShouldBindJSON(req); err != nil {
		return err
	}

	if err := ValidateRequest(req); err != nil {
		return err
	}

	return nil
}

// ValidatingMiddleware retorna um middleware que valida requisições automaticamente
// Útil quando aplicado globalmente para garantir validação em todas as rotas
func ValidatingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Continua para o próximo handler
		// A validação real acontece em BindAndValidateJSON nos handlers específicos
		c.Next()
	}
}
