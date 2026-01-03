package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/guilherme/pet-services-api/internal/application/user"
	domainauth "github.com/guilherme/pet-services-api/internal/domain/auth"
	domainuser "github.com/guilherme/pet-services-api/internal/domain/user"
	"github.com/guilherme/pet-services-api/internal/infra/factory"
)

// UserHandler expõe endpoints relacionados a usuário.
type UserHandler struct {
	uc factory.UserUseCases
}

// NewUserHandler cria handler de usuário.
func NewUserHandler(uc factory.UserUseCases) *UserHandler {
	return &UserHandler{uc: uc}
}

// RegisterUserRoutes registra rotas de usuário.
func RegisterUserRoutes(rg *gin.RouterGroup, h *UserHandler) {
	rg.GET("/me", h.GetProfile)
	rg.PUT("/me", h.UpdateProfile)
	rg.POST("/change-password", h.ChangePassword)
	rg.POST("/password-reset/request", h.RequestPasswordReset)
	rg.POST("/password-reset/confirm", h.ConfirmPasswordReset)
	rg.POST("/email/verification/request", h.RequestEmailVerification)
	rg.POST("/email/verification/confirm", h.ConfirmEmailVerification)
	rg.DELETE("/me", h.DeleteAccount)
}

// GetProfile retorna o perfil do usuário autenticado.
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_user_id", err.Error()))
		return
	}

	profile, err := h.uc.GetProfile.Execute(c.Request.Context(), user.GetProfileInput{UserID: userID})
	if err != nil {
		switch {
		case errors.Is(err, domainuser.ErrUserNotFound):
			c.JSON(http.StatusNotFound, errorResponse("user_not_found", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("get_profile_failed", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, userToResponse(profile))
}

// updateProfileRequest payload para atualização parcial.
type updateProfileRequest struct {
	Name      *string             `json:"name"`
	Phone     *string             `json:"phone"`
	Address   *domainuser.Address `json:"address"`
	Latitude  *float64            `json:"latitude"`
	Longitude *float64            `json:"longitude"`
}

// UpdateProfile atualiza dados do usuário autenticado.
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_user_id", err.Error()))
		return
	}

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
		return
	}

	if err := h.uc.UpdateProfile.Execute(c.Request.Context(), user.UpdateProfileInput{
		UserID:    userID,
		Name:      req.Name,
		Phone:     req.Phone,
		Address:   req.Address,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}); err != nil {
		switch {
		case errors.Is(err, domainuser.ErrUserNotFound):
			c.JSON(http.StatusNotFound, errorResponse("user_not_found", err.Error()))
		case errors.Is(err, domainuser.ErrInvalidPhone), strings.Contains(err.Error(), "coordenadas inválidas"):
			c.JSON(http.StatusBadRequest, errorResponse("invalid_data", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("update_profile_failed", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// changePasswordRequest payload.
type changePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// ChangePassword troca a senha do usuário autenticado.
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_user_id", err.Error()))
		return
	}

	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
		return
	}

	if err := h.uc.ChangePassword.Execute(c.Request.Context(), user.ChangePasswordInput{
		UserID:          userID,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}); err != nil {
		switch {
		case errors.Is(err, domainuser.ErrUserNotFound):
			c.JSON(http.StatusNotFound, errorResponse("user_not_found", err.Error()))
		case errors.Is(err, domainauth.ErrInvalidCredentials), errors.Is(err, domainuser.ErrInvalidPassword):
			c.JSON(http.StatusUnauthorized, errorResponse("invalid_credentials", err.Error()))
		default:
			c.JSON(http.StatusBadRequest, errorResponse("change_password_failed", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// requestPasswordResetRequest payload.
type requestPasswordResetRequest struct {
	Email string `json:"email"`
}

// RequestPasswordReset inicia fluxo de reset por email.
func (h *UserHandler) RequestPasswordReset(c *gin.Context) {
	var req requestPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
		return
	}

	if err := h.uc.RequestPasswordReset.Execute(c.Request.Context(), user.RequestPasswordResetInput{Email: req.Email}); err != nil {
		switch {
		case errors.Is(err, domainuser.ErrInvalidEmail):
			c.JSON(http.StatusBadRequest, errorResponse("invalid_email", err.Error()))
		default:
			// Por segurança, respostas genéricas
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// confirmPasswordResetRequest payload.
type confirmPasswordResetRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

// ConfirmPasswordReset aplica o reset de senha.
func (h *UserHandler) ConfirmPasswordReset(c *gin.Context) {
	var req confirmPasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
		return
	}

	if err := h.uc.ConfirmPasswordReset.Execute(c.Request.Context(), user.ConfirmPasswordResetInput{
		Token:       req.Token,
		NewPassword: req.NewPassword,
	}); err != nil {
		switch {
		case errors.Is(err, domainuser.ErrPasswordResetTokenInvalid), errors.Is(err, domainuser.ErrInvalidPassword):
			c.JSON(http.StatusBadRequest, errorResponse("invalid_token_or_password", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("confirm_reset_failed", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// RequestEmailVerification dispara email de verificação.
func (h *UserHandler) RequestEmailVerification(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_user_id", err.Error()))
		return
	}

	if err := h.uc.RequestEmailVerification.Execute(c.Request.Context(), user.RequestEmailVerificationInput{UserID: userID}); err != nil {
		switch {
		case errors.Is(err, domainuser.ErrUserNotFound):
			c.JSON(http.StatusNotFound, errorResponse("user_not_found", err.Error()))
		case errors.Is(err, domainuser.ErrEmailAlreadyVerified):
			c.JSON(http.StatusConflict, errorResponse("email_already_verified", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("request_email_verification_failed", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// confirmEmailVerificationRequest payload.
type confirmEmailVerificationRequest struct {
	Token string `json:"token"`
}

// ConfirmEmailVerification confirma email.
func (h *UserHandler) ConfirmEmailVerification(c *gin.Context) {
	var req confirmEmailVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
		return
	}

	if err := h.uc.ConfirmEmailVerification.Execute(c.Request.Context(), user.ConfirmEmailVerificationInput{Token: req.Token}); err != nil {
		switch {
		case errors.Is(err, domainuser.ErrEmailVerificationTokenInvalid):
			c.JSON(http.StatusBadRequest, errorResponse("invalid_token", err.Error()))
		case errors.Is(err, domainuser.ErrEmailAlreadyVerified):
			c.JSON(http.StatusConflict, errorResponse("email_already_verified", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("confirm_email_verification_failed", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// DeleteAccount deleta (soft ou hard) a conta do usuário autenticado.
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_user_id", err.Error()))
		return
	}

	hard := c.Query("hard") == "true"

	if err := h.uc.DeleteAccount.Execute(c.Request.Context(), user.DeleteAccountInput{UserID: userID, HardDelete: hard}); err != nil {
		switch {
		case errors.Is(err, domainuser.ErrUserNotFound):
			c.JSON(http.StatusNotFound, errorResponse("user_not_found", err.Error()))
		default:
			c.JSON(http.StatusBadRequest, errorResponse("delete_account_failed", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// extractUserID obtém user_id do header X-User-ID ou query user_id.
func extractUserID(c *gin.Context) (uuid.UUID, error) {
	uid := c.GetHeader("X-User-ID")
	if uid == "" {
		uid = c.Query("user_id")
	}
	if uid == "" {
		return uuid.Nil, errors.New("user_id é obrigatório (X-User-ID ou query)")
	}
	return uuid.Parse(uid)
}

// userToResponse converte domínio para resposta JSON amigável.
func userToResponse(u *domainuser.User) gin.H {
	resp := gin.H{
		"id":              u.ID,
		"email":           u.Email.String(),
		"email_verified":  u.EmailVerified,
		"name":            u.Name,
		"phone":           u.Phone.String(),
		"phone_formatted": u.Phone.Formatted(),
		"type":            u.Type,
		"created_at":      u.CreatedAt.Unix(),
		"updated_at":      u.UpdatedAt.Unix(),
	}

	if u.DeletedAt != nil {
		resp["deleted_at"] = u.DeletedAt.Unix()
	}

	if u.Location != nil {
		resp["location"] = gin.H{
			"latitude":  u.Location.Latitude,
			"longitude": u.Location.Longitude,
			"address":   u.Location.Address,
		}
	}

	return resp
}
