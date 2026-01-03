package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/guilherme/pet-services-api/internal/application/auth"
	domainauth "github.com/guilherme/pet-services-api/internal/domain/auth"
	"github.com/guilherme/pet-services-api/internal/domain/user"
	"github.com/guilherme/pet-services-api/internal/infra/factory"
)

// AuthHandler expõe endpoints de autenticação.
type AuthHandler struct {
	uc           factory.AuthUseCases
	errorService *ErrorService
}

// NewAuthHandler cria um novo handler de autenticação.
func NewAuthHandler(uc factory.AuthUseCases, errorService *ErrorService) *AuthHandler {
	return &AuthHandler{uc: uc, errorService: errorService}
}

// RegisterAuthRoutes registra rotas de autenticação em um grupo.
func RegisterAuthRoutes(rg *gin.RouterGroup, handler *AuthHandler) {
	rg.POST("/signup", handler.Signup)
	rg.POST("/login", handler.Login)
	rg.POST("/refresh", handler.Refresh)
	rg.POST("/logout", handler.Logout)
}

// signupRequest representa a entrada do endpoint de cadastro.
type signupRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Phone    string `json:"phone" validate:"required,min=10,max=20"`
	Password string `json:"password" validate:"required,min=8,max=128"`
	Type     string `json:"type" validate:"required,oneof=owner provider"`
}

// Signup cadastra um usuário e retorna tokens.
// @Summary Signup
// @Tags auth
// @Accept json
// @Produce json
// @Param request body signupRequest true "Signup payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	var req signupRequest
	if err := BindAndValidateJSON(c, &req); err != nil {
		h.errorService.RespondWithError(c, err, "invalid_payload", http.StatusBadRequest)
		return
	}

	userType, err := parseUserType(req.Type)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_user_type", err.Error()))
		return
	}

	out, err := h.uc.Signup.Execute(c.Request.Context(), auth.SignupInput{
		Email:    req.Email,
		Name:     req.Name,
		Phone:    req.Phone,
		Password: req.Password,
		Type:     userType,
	})
	if err != nil {
		switch {
		case errors.Is(err, user.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, errorResponse("user_exists", err.Error()))
		case errors.Is(err, user.ErrInvalidPassword), errors.Is(err, user.ErrInvalidEmail):
			c.JSON(http.StatusBadRequest, errorResponse("invalid_credentials", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("signup_failed", err.Error()))
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user_id":            out.UserID,
		"user_type":          out.UserType,
		"access_token":       out.AccessToken,
		"refresh_token":      out.RefreshToken,
		"access_expires_at":  out.AccessExpiresAt,
		"refresh_expires_at": out.RefreshExpiresAt,
	})
}

// loginRequest representa a entrada de login.
type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=1"`
}

// Login autentica e retorna tokens.
// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Login payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := BindAndValidateJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid_payload",
			"fields": ValidationErrorResponse(err),
		})
		return
	}

	out, err := h.uc.Login.Execute(c.Request.Context(), auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainauth.ErrInvalidCredentials), errors.Is(err, user.ErrUserNotFound):
			c.JSON(http.StatusUnauthorized, errorResponse("invalid_credentials", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("login_failed", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":            out.UserID,
		"user_type":          out.UserType,
		"access_token":       out.AccessToken,
		"refresh_token":      out.RefreshToken,
		"access_expires_at":  out.AccessExpiresAt,
		"refresh_expires_at": out.RefreshExpiresAt,
	})
}

// refreshRequest representa a entrada de refresh token.
type refreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Refresh troca o refresh token por um novo par.
// @Summary Refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body refreshRequest true "Refresh payload"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := BindAndValidateJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid_payload",
			"fields": ValidationErrorResponse(err),
		})
		return
	}

	out, err := h.uc.Refresh.Execute(c.Request.Context(), auth.RefreshInput{RefreshToken: req.RefreshToken})
	if err != nil {
		switch {
		case errors.Is(err, domainauth.ErrInvalidCredentials),
			errors.Is(err, domainauth.ErrRefreshTokenNotFound),
			errors.Is(err, domainauth.ErrRefreshTokenRevoked),
			errors.Is(err, domainauth.ErrRefreshTokenExpired):
			c.JSON(http.StatusUnauthorized, errorResponse("invalid_refresh_token", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("refresh_failed", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":            out.UserID,
		"user_type":          out.UserType,
		"access_token":       out.AccessToken,
		"refresh_token":      out.RefreshToken,
		"access_expires_at":  out.AccessExpiresAt,
		"refresh_expires_at": out.RefreshExpiresAt,
	})
}

// logoutRequest representa a entrada de logout.
type logoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// Logout revoga o refresh token atual.
// @Summary Logout
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <refresh_token>"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req logoutRequest
	if err := BindAndValidateJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid_payload",
			"fields": ValidationErrorResponse(err),
		})
		return
	}

	if err := h.uc.Logout.Execute(c.Request.Context(), auth.LogoutInput{RefreshToken: req.RefreshToken}); err != nil {
		switch {
		case errors.Is(err, domainauth.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, errorResponse("invalid_refresh_token", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("logout_failed", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// parseUserType converte string para user.UserType.
func parseUserType(t string) (user.UserType, error) {
	switch t {
	case string(user.UserTypeOwner):
		return user.UserTypeOwner, nil
	case string(user.UserTypeProvider):
		return user.UserTypeProvider, nil
	default:
		return "", user.ErrInvalidUserType
	}
}

// errorResponse padroniza resposta de erro.
func errorResponse(code, message string) gin.H {
	return gin.H{"error": code, "message": message}
}
