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
	Logger      logging.LoggerInterface
}

func NewUserHandler(factory *factories.UserFactory, logger logging.LoggerInterface) *UserHandler {
	return &UserHandler{
		UserFactory: factory,
		Logger:      logger,
	}
}

// RegisterUser godoc
// @Summary Registra um novo usuário
// @Tags Users
// @Accept json
// @Produce json
// @Param input body usecases.RegisterUserInput true "Dados do usuário"
// @Success 201 {object} usecases.RegisterUserOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Router /users/register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.RegisterUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do usuário",
		})
		h.Logger.LogBadRequest(ctx, "UserHandler", problem.Detail, err)
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

// GetProfile godoc
// @Summary Retorna o perfil do usuário autenticado
// @Tags Users
// @Produce json
// @Success 200 {object} usecases.GetProfileOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	ctx := c.Request.Context()
	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "UserHandler", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}
	input := usecases.GetProfileInput{UserID: userID.(string)}
	output, errs := h.UserFactory.GetProfile.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

// ListUsers godoc
// @Summary Lista usuários
// @Tags Users
// @Produce json
// @Param page query int false "Página"
// @Param limit query int false "Limite"
// @Success 200 {array} usecases.ListUsersOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.ListUsersInput
	if err := c.ShouldBindQuery(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos parâmetros de listagem",
		})
		h.Logger.LogBadRequest(ctx, "UserHandler", problem.Detail, err)
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

// UpdateUser godoc
// @Summary Atualiza dados do usuário
// @Tags Users
// @Accept json
// @Produce json
// @Param input body usecases.UpdateUserInput true "Dados do usuário"
// @Success 200 {object} usecases.UpdateUserOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /users [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados do usuário",
		})
		h.Logger.LogBadRequest(ctx, "UserHandler", problem.Detail, err)
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

// DeleteUser godoc
// @Summary Deleta usuário
// @Tags Users
// @Accept json
// @Produce json
// @Param input body usecases.DeleteUserInput true "ID do usuário"
// @Success 200 {object} usecases.DeleteUserOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /users [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.DeleteUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do ID do usuário",
		})
		h.Logger.LogBadRequest(ctx, "UserHandler", problem.Detail, err)
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

// ReactivateUser godoc
// @Summary Reativa usuário
// @Tags Users
// @Accept json
// @Produce json
// @Param input body usecases.ReactivateUserInput true "ID do usuário"
// @Success 200 {object} usecases.ReactivateUserOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /users/reactivate [post]
func (h *UserHandler) ReactivateUser(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.ReactivateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do ID do usuário",
		})
		h.Logger.LogBadRequest(ctx, "UserHandler", problem.Detail, err)
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

// DeactivateUser godoc
// @Summary Desativa usuário
// @Tags Users
// @Accept json
// @Produce json
// @Param input body usecases.DeactivateUserInput true "ID do usuário"
// @Success 200 {object} usecases.DeactivateUserOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /users/deactivate [post]
func (h *UserHandler) DeactivateUser(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.DeactivateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do ID do usuário",
		})
		h.Logger.LogBadRequest(ctx, "UserHandler", problem.Detail, err)
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

// GetUserByID godoc
// @Summary Busca usuário por ID
// @Tags Users
// @Produce json
// @Param user_id path string true "ID do usuário"
// @Success 200 {object} usecases.GetUserByIDOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /users/{user_id} [get]
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

// CheckEmailExists godoc
// @Summary Verifica se email existe
// @Tags Users
// @Accept json
// @Produce json
// @Param input body usecases.CheckEmailExistsInput true "Email"
// @Success 200 {object} usecases.CheckEmailExistsOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Router /users/check-email [post]
func (h *UserHandler) CheckEmailExists(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.CheckEmailExistsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do email",
		})
		h.Logger.LogBadRequest(ctx, "UserHandler", problem.Detail, err)
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

// CheckPhoneExists godoc
// @Summary Verifica se telefone existe
// @Tags Users
// @Accept json
// @Produce json
// @Param input body usecases.CheckPhoneExistsInput true "Telefone"
// @Success 200 {object} usecases.CheckPhoneExistsOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Router /users/check-phone [post]
func (h *UserHandler) CheckPhoneExists(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.CheckPhoneExistsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser do telefone",
		})
		h.Logger.LogBadRequest(ctx, "UserHandler", problem.Detail, err)
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

// UpdateEmailVerified godoc
// @Summary Atualiza verificação de email
// @Tags Users
// @Accept json
// @Produce json
// @Param input body usecases.UpdateEmailVerifiedInput true "Dados de verificação"
// @Success 200 {object} usecases.UpdateEmailVerifiedOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Router /users/update-email-verified [post]
func (h *UserHandler) UpdateEmailVerified(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.UpdateEmailVerifiedInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados de verificação de email",
		})
		h.Logger.LogBadRequest(ctx, "UserHandler", problem.Detail, err)
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

// ChangePassword godoc
// @Summary Altera senha do usuário
// @Tags Users
// @Accept json
// @Produce json
// @Param input body usecases.ChangePasswordInput true "Dados de senha"
// @Success 200 {object} usecases.ChangePasswordOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /users/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados de senha",
		})
		h.Logger.LogBadRequest(ctx, "UserHandler", problem.Detail, err)
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

// CreateAdmin godoc
// @Summary Cria um novo administrador
// @Tags Admin
// @Accept json
// @Produce json
// @Param input body usecases.CreateAdminInput true "Dados do admin"
// @Success 201 {object} usecases.CreateAdminOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /admin [post]
func (h *UserHandler) CreateAdmin(c *gin.Context) {
	ctx := c.Request.Context()
	var input usecases.CreateAdminInput
	if err := c.ShouldBindJSON(&input); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados do admin",
		})
		h.Logger.LogBadRequest(ctx, "UserHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	isAdmin := false
	if claims, ok := c.Get("is_admin"); ok {
		if v, ok := claims.(bool); ok && v {
			isAdmin = true
		}
	}

	output, errs := h.UserFactory.RegisterAdmin.Execute(ctx, input, isAdmin)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusCreated, output)
}
