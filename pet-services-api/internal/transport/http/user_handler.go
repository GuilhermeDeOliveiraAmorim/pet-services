package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/guilherme/pet-services-api/internal/application/user"
	domainuser "github.com/guilherme/pet-services-api/internal/domain/user"
	"github.com/guilherme/pet-services-api/internal/infra/factory"
)

// UserHandler expõe endpoints relacionados a usuário.
type UserHandler struct {
	uc           factory.UserUseCases
	errorService *ErrorService
}

// NewUserHandler cria handler de usuário.
func NewUserHandler(uc factory.UserUseCases, errorService *ErrorService) *UserHandler {
	return &UserHandler{uc: uc, errorService: errorService}
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
// GetProfile retorna o perfil do usuário autenticado.
// @Summary Get current user profile
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /users/me [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		h.errorService.RespondWithError(c, err, "invalid_user_id", http.StatusBadRequest)
		return
	}

	profile, problems := h.uc.GetProfile.Execute(c.Request.Context(), user.GetProfileInput{UserID: userID})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "get_profile_failed", http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, userToResponse(profile))
}

// updateProfileRequest payload para atualização parcial.
// UpdateProfile atualiza dados do usuário autenticado.
// @Summary Update current user profile
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body updateProfileRequest true "Profile payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /users/me [put]
type updateProfileRequest struct {
	Name      *string             `json:"name" validate:"omitempty,min=3,max=100"`
	Phone     *string             `json:"phone" validate:"omitempty,min=10,max=20"`
	Address   *domainuser.Address `json:"address"`
	Latitude  *float64            `json:"latitude" validate:"omitempty,min=-90,max=90"`
	Longitude *float64            `json:"longitude" validate:"omitempty,min=-180,max=180"`
}

// UpdateProfile atualiza dados do usuário autenticado.
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_user_id", err.Error()))
		return
	}

	var req updateProfileRequest
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}

	problems := h.uc.UpdateProfile.Execute(c.Request.Context(), user.UpdateProfileInput{
		UserID:    userID,
		Name:      req.Name,
		Phone:     req.Phone,
		Address:   req.Address,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "update_profile_failed", http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// changePasswordRequest payload.
type changePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required,min=1"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=128"`
}

// ChangePassword troca a senha do usuário autenticado.
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_user_id", err.Error()))
		return
	}

	var req changePasswordRequest
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}

	problems := h.uc.ChangePassword.Execute(c.Request.Context(), user.ChangePasswordInput{
		UserID:          userID,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "change_password_failed", http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// requestPasswordResetRequest payload.
type requestPasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// RequestPasswordReset inicia fluxo de reset por email.
func (h *UserHandler) RequestPasswordReset(c *gin.Context) {
	var req requestPasswordResetRequest
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}

	problems := h.uc.RequestPasswordReset.Execute(c.Request.Context(), user.RequestPasswordResetInput{Email: req.Email})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "request_password_reset_failed", http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// confirmPasswordResetRequest payload.
type confirmPasswordResetRequest struct {
	Token       string `json:"token" validate:"required,min=10"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=128"`
}

// RequestEmailVerification envia email com link de verificação.
// @Summary Request email verification
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /users/email/verification/request [post]

// ConfirmPasswordReset aplica o reset de senha.
func (h *UserHandler) ConfirmPasswordReset(c *gin.Context) {
	var req confirmPasswordResetRequest
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}

	problems := h.uc.ConfirmPasswordReset.Execute(c.Request.Context(), user.ConfirmPasswordResetInput{
		Token:       req.Token,
		NewPassword: req.NewPassword,
	})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "confirm_password_reset_failed", http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	// ConfirmEmailVerification confirma email com código/token.
	// @Summary Confirm email verification
	// @Tags users
	// @Accept json
	// @Produce json
	// @Param request body confirmEmailVerificationRequest true "Verification confirm"
	// @Success 200 {object} map[string]interface{}
	// @Failure 400 {object} map[string]interface{}
	// @Router /users/email/verification/confirm [post]
}

// RequestEmailVerification dispara email de verificação.
func (h *UserHandler) RequestEmailVerification(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_user_id", err.Error()))
		return
	}

	problems := h.uc.RequestEmailVerification.Execute(c.Request.Context(), user.RequestEmailVerificationInput{UserID: userID})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "request_email_verification_failed", http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
	// DeleteAccount deleta (soft ou hard) a conta do usuário autenticado.
	// @Summary Delete account
	// @Tags users
	// @Security BearerAuth
	// @Produce json
	// @Param hard query bool false "Hard delete"
	// @Success 200 {object} map[string]interface{}
	// @Failure 400 {object} map[string]interface{}
	// @Router /users/me [delete]
}

// confirmEmailVerificationRequest payload.
type confirmEmailVerificationRequest struct {
	Token string `json:"token" validate:"required,min=10"`
}

// ConfirmEmailVerification confirma email.
func (h *UserHandler) ConfirmEmailVerification(c *gin.Context) {
	var req confirmEmailVerificationRequest
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}

	problems := h.uc.ConfirmEmailVerification.Execute(c.Request.Context(), user.ConfirmEmailVerificationInput{Token: req.Token})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "confirm_email_verification_failed", http.StatusBadRequest)
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

	problems := h.uc.DeleteAccount.Execute(c.Request.Context(), user.DeleteAccountInput{UserID: userID, HardDelete: hard})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "delete_account_failed", http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// extractUserID obtém user_id do contexto (middleware) ou fallback para header/query.
func extractUserID(c *gin.Context) (uuid.UUID, error) {
	if val, ok := c.Get(ctxUserIDKey); ok {
		if id, ok := val.(uuid.UUID); ok {
			return id, nil
		}
	}

	uid := c.GetHeader("X-User-ID")
	if uid == "" {
		uid = c.Query("user_id")
	}
	if uid == "" {
		return uuid.Nil, errors.New("user_id é obrigatório (Authorization Bearer, X-User-ID ou query)")
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
