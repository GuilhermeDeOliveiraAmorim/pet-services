package handlers

import (
	"net/http"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/usecases"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/logging"
	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	TokenFactory *factories.TokenFactory
}

func NewTokenHandler(factory *factories.TokenFactory) *TokenHandler {
	return &TokenHandler{
		TokenFactory: factory,
	}
}

func (h *TokenHandler) LoginUser(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.LoginUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do login",
		})
		logging.BadRequest(ctx, "TokenHandler", problem.Detail, err)
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

func (h *TokenHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.LogoutInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do logout",
		})
		logging.BadRequest(ctx, "TokenHandler", problem.Detail, err)
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

func (h *TokenHandler) RequestPasswordReset(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.RequestPasswordResetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do request de reset de senha",
		})
		logging.BadRequest(ctx, "TokenHandler", problem.Detail, err)
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

func (h *TokenHandler) ResetPassword(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.ResetPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do reset de senha",
		})
		logging.BadRequest(ctx, "TokenHandler", problem.Detail, err)
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

func (h *TokenHandler) ResendVerificationEmail(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.ResendVerificationEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do reenvio de verificação",
		})
		logging.BadRequest(ctx, "TokenHandler", problem.Detail, err)
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

func (h *TokenHandler) VerifyEmail(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.VerifyEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser da verificação de email",
		})
		logging.BadRequest(ctx, "TokenHandler", problem.Detail, err)
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
