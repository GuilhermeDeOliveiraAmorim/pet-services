package handlers

import (
	"net/http"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserFactory *factories.UserFactory
}

func NewUserHandler(factory *factories.UserFactory) *UserHandler {
	return &UserHandler{
		UserFactory: factory,
	}
}

func (h *UserHandler) RegisterUser(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.RegisterUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do usuário",
		})
		logging.BadRequest(ctx, "UserHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}
	output, errs := h.UserFactory.RegisterUser.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusCreated, output)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("user_id")
	input := usecases.GetProfileInput{UserID: userID}
	output, errs := h.UserFactory.GetProfile.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.ListUsersInput
	if err := c.ShouldBindQuery(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos parâmetros de listagem",
		})
		logging.BadRequest(ctx, "UserHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}
	output, errs := h.UserFactory.ListUsers.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados do usuário",
		})
		logging.BadRequest(ctx, "UserHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}
	output, errs := h.UserFactory.UpdateUser.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.DeleteUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do ID do usuário",
		})
		logging.BadRequest(ctx, "UserHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}
	output, errs := h.UserFactory.DeleteUser.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) ReactivateUser(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.ReactivateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do ID do usuário",
		})
		logging.BadRequest(ctx, "UserHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}
	output, errs := h.UserFactory.ReactivateUser.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) DeactivateUser(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.DeactivateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do ID do usuário",
		})
		logging.BadRequest(ctx, "UserHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}
	output, errs := h.UserFactory.DeactivateUser.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("user_id")
	input := usecases.GetUserByIDInput{UserID: userID}
	output, errs := h.UserFactory.GetUserByID.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) CheckEmailExists(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.CheckEmailExistsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do email",
		})
		logging.BadRequest(ctx, "UserHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}
	output, errs := h.UserFactory.CheckEmailExists.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) CheckPhoneExists(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.CheckPhoneExistsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do telefone",
		})
		logging.BadRequest(ctx, "UserHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}
	output, errs := h.UserFactory.CheckPhoneExists.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) UpdateEmailVerified(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.UpdateEmailVerifiedInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados de verificação de email",
		})
		logging.BadRequest(ctx, "UserHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}
	output, errs := h.UserFactory.UpdateEmailVerified.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados de senha",
		})
		logging.BadRequest(ctx, "UserHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}
	output, errs := h.UserFactory.ChangePassword.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}
