package handlers

import (
	"net/http"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"

	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	TokenFactory *factories.TokenFactory
	Logger       logging.LoggerInterface
}

func NewTokenHandler(factory *factories.TokenFactory, logger logging.LoggerInterface) *TokenHandler {
	return &TokenHandler{
		TokenFactory: factory,
		Logger:       logger,
	}
}

// LoginUser godoc
// @Summary Realiza login do usuário
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body usecases.LoginUserInput true "Dados de login"
// @Success 200 {object} usecases.LoginUserOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Router /auth/login [post]
func (h *TokenHandler) LoginUser(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.LoginUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do login",
		})
		h.Logger.LogBadRequest(ctx, "TokenHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}
	output, errs := h.TokenFactory.LoginUser.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

// Logout godoc
// @Summary Realiza logout do usuário
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body usecases.LogoutInput true "Dados de logout"
// @Success 200 {object} usecases.LogoutOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /auth/logout [post]
func (h *TokenHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.LogoutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do logout",
		})
		h.Logger.LogBadRequest(ctx, "TokenHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}
	output, errs := h.TokenFactory.Logout.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

// RequestPasswordReset godoc
// @Summary Solicita reset de senha
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body usecases.RequestPasswordResetInput true "Dados para reset de senha"
// @Success 200 {object} usecases.RequestPasswordResetOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Router /auth/request-password-reset [post]
func (h *TokenHandler) RequestPasswordReset(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.RequestPasswordResetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do request de reset de senha",
		})
		h.Logger.LogBadRequest(ctx, "TokenHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}
	output, errs := h.TokenFactory.RequestPasswordReset.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

// ResetPassword godoc
// @Summary Realiza reset de senha
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body usecases.ResetPasswordInput true "Dados para reset de senha"
// @Success 200 {object} usecases.ResetPasswordOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Router /auth/reset-password [post]
func (h *TokenHandler) ResetPassword(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do reset de senha",
		})
		h.Logger.LogBadRequest(ctx, "TokenHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}
	output, errs := h.TokenFactory.ResetPassword.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

// ResendVerificationEmail godoc
// @Summary Reenvia email de verificação
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body usecases.ResendVerificationEmailInput true "Dados para reenvio de verificação"
// @Success 200 {object} usecases.ResendVerificationEmailOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Router /auth/resend-verification-email [post]
func (h *TokenHandler) ResendVerificationEmail(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.ResendVerificationEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do reenvio de verificação",
		})
		h.Logger.LogBadRequest(ctx, "TokenHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}
	output, errs := h.TokenFactory.ResendVerificationEmail.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

// VerifyEmail godoc
// @Summary Verifica email do usuário
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body usecases.VerifyEmailInput true "Dados para verificação de email"
// @Success 200 {object} usecases.VerifyEmailOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Router /auth/verify-email [post]
func (h *TokenHandler) VerifyEmail(c *gin.Context) {
	ctx := c.Request.Context()

	var input usecases.VerifyEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser da verificação de email",
		})
		h.Logger.LogBadRequest(ctx, "TokenHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	output, errs := h.TokenFactory.VerifyEmail.Execute(ctx, input)

	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}
