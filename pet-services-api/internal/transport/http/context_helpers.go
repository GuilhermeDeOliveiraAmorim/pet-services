package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"pet-services-api/internal/application/exceptions"
	domainprovider "pet-services-api/internal/domain/provider"
	domainuser "pet-services-api/internal/domain/user"
)

// extractUserIDProblems extrai o user_id do contexto, retornando problemas padronizados.
func extractUserIDProblems(c *gin.Context) (uuid.UUID, []exceptions.ProblemDetails) {
	userID, err := extractUserID(c)
	if err != nil {
		// Log detalhado do erro de extração do user_id
		slog.Warn("extractUserIDProblems: user_id não encontrado no contexto", slog.String("error", err.Error()))
		return uuid.Nil, []exceptions.ProblemDetails{{
			Type:   string(exceptions.Unauthorized),
			Title:  "Não autorizado",
			Status: http.StatusUnauthorized,
			Detail: err.Error(),
		}}
	}
	slog.Info("extractUserIDProblems: user_id extraído com sucesso", slog.String("user_id", userID.String()))
	return userID, nil
}

// providerIDFromContextProblems resolve provider vinculado ao usuário autenticado, retornando problemas padronizados.
func providerIDFromContextProblems(c *gin.Context, repo domainprovider.Repository, requireProvider bool) (uuid.UUID, []exceptions.ProblemDetails) {
	userType := extractUserType(c)
	if requireProvider && userType != "" && userType != domainuser.UserTypeProvider {
		return uuid.Nil, []exceptions.ProblemDetails{{
			Type:   string(exceptions.Forbidden),
			Title:  "Apenas prestadores podem executar esta ação",
			Status: http.StatusForbidden,
		}}
	}
	if !requireProvider && userType != domainuser.UserTypeProvider {
		return uuid.Nil, nil
	}
	userID, err := extractUserID(c)
	if err != nil {
		return uuid.Nil, []exceptions.ProblemDetails{{
			Type:   string(exceptions.Unauthorized),
			Title:  "Não autorizado",
			Status: http.StatusUnauthorized,
			Detail: err.Error(),
		}}
	}
	p, err := repo.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound), errors.Is(err, domainprovider.ErrProviderNotFound):
			return uuid.Nil, []exceptions.ProblemDetails{{
				Type:   string(exceptions.NotFound),
				Title:  "Prestador não encontrado para este usuário",
				Status: http.StatusNotFound,
			}}
		default:
			return uuid.Nil, []exceptions.ProblemDetails{{
				Type:   string(exceptions.InternalServerError),
				Title:  "Erro ao buscar prestador",
				Status: http.StatusInternalServerError,
				Detail: err.Error(),
			}}
		}
	}
	return p.ID, nil
}

// extractUserType retorna o tipo de usuário do contexto, se presente.
func extractUserType(c *gin.Context) domainuser.UserType {
	if val, ok := c.Get(ctxUserTypeKey); ok {
		switch v := val.(type) {
		case domainuser.UserType:
			return v
		case string:
			return domainuser.UserType(v)
		}
	}
	return ""
}

// extractUserID obtém user_id do contexto (middleware) ou fallback para header/query, com logs detalhados.
func extractUserID(c *gin.Context) (uuid.UUID, error) {
	if val, ok := c.Get(ctxUserIDKey); ok {
		switch v := val.(type) {
		case string:
			id, err := uuid.Parse(v)
			if err != nil {
				slog.Error("extractUserID: user_id inválido no contexto", slog.String("user_id", v), slog.String("error", err.Error()))
				return uuid.Nil, err
			}
			slog.Info("extractUserID: user_id encontrado no contexto", slog.String("user_id", id.String()))
			return id, nil
		case uuid.UUID:
			slog.Info("extractUserID: user_id encontrado no contexto (uuid.UUID)", slog.String("user_id", v.String()))
			return v, nil
		default:
			slog.Warn("extractUserID: valor encontrado no contexto não é string nem uuid.UUID", slog.Any("value", val))
			return uuid.Nil, errors.New("user_id inválido no contexto")
		}
	}
	slog.Error("extractUserID: user_id não encontrado no contexto")
	return uuid.Nil, errors.New("user_id é obrigatório (Authorization Bearer)")
}
