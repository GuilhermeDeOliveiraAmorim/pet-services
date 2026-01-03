package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/guilherme/pet-services-api/internal/application/auth"
	"github.com/guilherme/pet-services-api/internal/application/exceptions"
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
// @Success 201 {object} auth.SignupOutput
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Router /auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	var req signupRequest
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}

	userType, err := parseUserType(req.Type)
	if err != nil {
		h.errorService.RespondWithProblems(c, []exceptions.ProblemDetails{{
			Type:   string(exceptions.BadRequest),
			Title:  "Tipo de usuário inválido",
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}}, "invalid_user_type", http.StatusBadRequest)
		return
	}

	out, problems := h.uc.Signup.Execute(c.Request.Context(), auth.SignupInput{
		Email:    req.Email,
		Name:     req.Name,
		Phone:    req.Phone,
		Password: req.Password,
		Type:     userType,
	})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "signup_failed", http.StatusBadRequest)
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
// @Success 200 {object} auth.LoginOutput
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}

	out, problems := h.uc.Login.Execute(c.Request.Context(), auth.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "login_failed", http.StatusUnauthorized)
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
// @Success 200 {object} auth.RefreshOutput
// @Failure 401 {object} exceptions.ProblemDetailsOutputDTO
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req refreshRequest
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}

	out, problems := h.uc.Refresh.Execute(c.Request.Context(), auth.RefreshInput{RefreshToken: req.RefreshToken})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_refresh_token", http.StatusUnauthorized)
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
// @Success 200 {object} map[string]interface{}  // logout não retorna body customizado
// @Failure 401 {object} exceptions.ProblemDetailsOutputDTO
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req logoutRequest
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}

	problems := h.uc.Logout.Execute(c.Request.Context(), auth.LogoutInput{RefreshToken: req.RefreshToken})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_refresh_token", http.StatusUnauthorized)
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
