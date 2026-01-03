package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	domainprovider "github.com/guilherme/pet-services-api/internal/domain/provider"
	domainuser "github.com/guilherme/pet-services-api/internal/domain/user"
)

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

// providerIDFromContext resolve o provider vinculado ao usuário autenticado.
// Se requireProvider for true, também valida se o user_type é provider.
func providerIDFromContext(c *gin.Context, repo domainprovider.Repository, requireProvider bool) (uuid.UUID, bool) {
	userType := extractUserType(c)
	if requireProvider && userType != "" && userType != domainuser.UserTypeProvider {
		c.JSON(http.StatusForbidden, errorResponse("forbidden", "apenas prestadores podem executar esta ação"))
		return uuid.Nil, false
	}

	// Para ações opcionais, apenas resolve se de fato é provider.
	if !requireProvider && userType != domainuser.UserTypeProvider {
		return uuid.Nil, true
	}

	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse("unauthorized", err.Error()))
		return uuid.Nil, false
	}

	p, err := repo.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound), errors.Is(err, domainprovider.ErrProviderNotFound):
			c.JSON(http.StatusNotFound, errorResponse("provider_not_found", "prestador não encontrado para este usuário"))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("provider_lookup_failed", err.Error()))
		}
		return uuid.Nil, false
	}

	return p.ID, true
}
