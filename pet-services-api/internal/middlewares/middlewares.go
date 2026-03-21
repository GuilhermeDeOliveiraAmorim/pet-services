package middlewares

import (
	"errors"
	"net/http"
	"pet-services-api/internal/auth"
	"pet-services-api/internal/consts"
	"pet-services-api/internal/entities"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(logger logging.LoggerInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractBearerToken(c, logger)
		if tokenString == "" {
			return
		}

		jwtService, err := auth.NewJWTServiceFromEnv()
		if err != nil {
			abortUnauthorized(c, "Configuração JWT inválida", "Erro ao carregar configuração de autenticação", logger)
			return
		}

		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			abortUnauthorized(c, "Token inválido", err.Error(), logger)
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_type", claims.UserType)
		c.Next()
	}
}

func extractBearerToken(c *gin.Context, logger logging.LoggerInterface) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		abortUnauthorized(c, "Autorização não fornecida", "Cabeçalho de autorização está vazio", logger)
		return ""
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		abortUnauthorized(c, "Bearer token inválido", "Cabeçalho de autorização inválido", logger)
		return ""
	}

	return parts[1]
}

func abortUnauthorized(c *gin.Context, title, detail string, logger logging.LoggerInterface) {
	problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
		Title:  title,
		Detail: detail,
	})

	logger.LogError(c.Request.Context(), "AuthMiddleware", title, nil)

	c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
}

func OwnerOnlyMiddleware(logger logging.LoggerInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists || userType != "owner" {
			problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
				Title:  "Acesso negado",
				Detail: "Acesso permitido apenas para usuários do tipo owner",
			})
			logger.LogError(c.Request.Context(), "AuthMiddleware", problem.Title+": "+problem.Detail, errors.New(problem.Detail))
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": problem})
			return
		}
		c.Next()
	}
}

func ProviderOnlyMiddleware(logger logging.LoggerInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists || userType != "provider" {
			problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
				Title:  "Acesso negado",
				Detail: "Acesso permitido apenas para usuários do tipo provider",
			})
			logger.LogError(c.Request.Context(), "AuthMiddleware", problem.Title+": "+problem.Detail, errors.New(problem.Detail))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": problem})
			return
		}
		c.Next()
	}
}

func AdminOnlyMiddleware(logger logging.LoggerInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists || userType != "admin" {
			problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
				Title:  "Acesso negado",
				Detail: "Acesso permitido apenas para administradores",
			})
			logger.LogError(c.Request.Context(), "AuthMiddleware", problem.Title+": "+problem.Detail, errors.New(problem.Detail))
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": problem})
			return
		}
		c.Set("is_admin", true)
		c.Next()
	}
}

func ProfileCompleteMiddleware(logger logging.LoggerInterface, userRepository entities.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			abortUnauthorized(c, "Usuário não autenticado", "Não foi possível obter o ID do usuário autenticado", logger)
			return
		}

		user, err := userRepository.FindByID(userID.(string))
		if err != nil {
			if err.Error() == consts.UserNotFoundError {
				problem := exceptions.NewProblemDetails(exceptions.NotFound, exceptions.ErrorMessage{
					Title:  "Usuário não encontrado",
					Detail: "Não foi possível encontrar o usuário autenticado",
				})
				logger.LogError(c.Request.Context(), "ProfileCompleteMiddleware", problem.Title, err)
				c.AbortWithStatusJSON(http.StatusNotFound, problem)
				return
			}
			problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
				Title:  "Erro ao buscar usuário",
				Detail: "Não foi possível validar o perfil do usuário",
			})
			logger.LogError(c.Request.Context(), "ProfileCompleteMiddleware", problem.Title, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, problem)
			return
		}

		if user.ProfileComplete {
			c.Next()
			return
		}

		problem := exceptions.NewProblemDetails(exceptions.Forbidden, exceptions.ErrorMessage{
			Title:  "Cadastro incompleto",
			Detail: "Complete seu cadastro para acessar este recurso",
		})
		logger.LogError(c.Request.Context(), "ProfileCompleteMiddleware", problem.Title, errors.New(problem.Detail))
		c.AbortWithStatusJSON(http.StatusForbidden, problem)
	}
}

func AdoptionGuardianApprovedMiddleware(logger logging.LoggerInterface, guardianRepo entities.AdoptionGuardianProfileRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			abortUnauthorized(c, "Usuário não autenticado", "Não foi possível obter o ID do usuário autenticado", logger)
			return
		}

		profile, err := guardianRepo.FindByUserID(userID.(string))
		if err != nil {
			if err.Error() == consts.AdoptionGuardianProfileNotFoundError {
				problem := exceptions.NewProblemDetails(exceptions.Forbidden, exceptions.ErrorMessage{
					Title:  "Perfil de responsável não encontrado",
					Detail: "Você precisa criar e ter seu perfil de responsável por adoção aprovado para realizar esta ação",
				})
				logger.LogError(c.Request.Context(), "AdoptionGuardianApprovedMiddleware", problem.Detail, errors.New(problem.Detail))
				c.AbortWithStatusJSON(http.StatusForbidden, problem)
				return
			}
			problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
				Title:  "Erro ao verificar perfil",
				Detail: "Não foi possível verificar o perfil de responsável por adoção",
			})
			logger.LogError(c.Request.Context(), "AdoptionGuardianApprovedMiddleware", problem.Detail, err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, problem)
			return
		}

		if profile.ApprovalStatus != entities.AdoptionGuardianApprovalStatuses.Approved {
			problem := exceptions.NewProblemDetails(exceptions.Forbidden, exceptions.ErrorMessage{
				Title:  "Perfil não aprovado",
				Detail: "Seu perfil de responsável por adoção ainda não foi aprovado. Aguarde a análise da equipe",
			})
			logger.LogError(c.Request.Context(), "AdoptionGuardianApprovedMiddleware", problem.Detail, errors.New(problem.Detail))
			c.AbortWithStatusJSON(http.StatusForbidden, problem)
			return
		}

		c.Set("guardian_profile_id", profile.ID)
		c.Next()
	}
}
