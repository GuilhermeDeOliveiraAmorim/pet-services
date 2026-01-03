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
	uc factory.AuthUseCases
}

// NewAuthHandler cria um novo handler de autenticação.
func NewAuthHandler(uc factory.AuthUseCases) *AuthHandler {
	return &AuthHandler{uc: uc}
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
	Email    string `json:"email"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Type     string `json:"type"`
}

// Signup cadastra um usuário e retorna tokens.
func (h *AuthHandler) Signup(c *gin.Context) {
	var req signupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
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
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login autentica e retorna tokens.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
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
	RefreshToken string `json:"refresh_token"`
}

// Refresh troca o refresh token por um novo par.
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
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
	RefreshToken string `json:"refresh_token"`
}

// Logout revoga o refresh token atual.
func (h *AuthHandler) Logout(c *gin.Context) {
	var req logoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
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
