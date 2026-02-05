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
// @Tags Usuários
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
// @Tags Usuários
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
// @Tags Usuários
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
// @Summary Atualiza dados do usuário autenticado
// @Tags Usuários
// @Accept json
// @Produce json
// @Param input body usecases.UpdateUserInputBody true "Dados do usuário"
// @Success 200 {object} usecases.UpdateUserOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /users [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
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

	var inputBody usecases.UpdateUserInputBody
	if err := c.ShouldBindJSON(&inputBody); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados do usuário",
		})
		h.Logger.LogBadRequest(ctx, "UserHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.UpdateUserInput{
		UserID:  userID.(string),
		Name:    inputBody.Name,
		Phone:   inputBody.Phone,
		Address: inputBody.Address,
	}

	output, errs := h.UserFactory.UpdateUser.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

// DeleteUser godoc
// @Summary Deleta a conta do usuário autenticado
// @Tags Usuários
// @Produce json
// @Success 200 {object} usecases.DeleteUserOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /users [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
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

	input := usecases.DeleteUserInput{UserID: userID.(string)}

	output, errs := h.UserFactory.DeleteUser.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

// ReactivateUser godoc
// @Summary Reativa a conta do usuário autenticado
// @Tags Usuários
// @Produce json
// @Success 200 {object} usecases.ReactivateUserOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /users/reactivate [post]
func (h *UserHandler) ReactivateUser(c *gin.Context) {
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

	input := usecases.ReactivateUserInput{UserID: userID.(string)}

	output, errs := h.UserFactory.ReactivateUser.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

// DeactivateUser godoc
// @Summary Desativa a conta do usuário autenticado
// @Tags Usuários
// @Produce json
// @Success 200 {object} usecases.DeactivateUserOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /users/deactivate [post]
func (h *UserHandler) DeactivateUser(c *gin.Context) {
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

	input := usecases.DeactivateUserInput{UserID: userID.(string)}

	output, errs := h.UserFactory.DeactivateUser.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

// GetUserByID godoc
// @Summary Busca usuário por ID
// @Tags Usuários
// @Produce json
// @Param user_id path string true "ID do usuário"
// @Success 200 {object} usecases.GetUserByIDOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /users/{user_id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.Param("user_id")

	authUserID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "UserHandler", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	authUserType, exists := c.Get("user_type")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o tipo do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "UserHandler", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	input := usecases.GetUserByIDInput{
		UserID:        userID,
		RequesterID:   authUserID.(string),
		RequesterType: authUserType.(string),
	}

	output, errs := h.UserFactory.GetUserByID.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}
	c.JSON(http.StatusOK, output)
}

// CheckEmailExists godoc
// @Summary Verifica se email existe
// @Tags Usuários
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
// @Tags Usuários
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
// @Tags Usuários
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
// @Summary Altera a senha do usuário autenticado
// @Tags Usuários
// @Accept json
// @Produce json
// @Param input body usecases.ChangePasswordInputBody true "Senhas atual e nova"
// @Success 200 {object} usecases.ChangePasswordOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /users/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
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

	var body usecases.ChangePasswordInputBody
	if err := c.ShouldBindJSON(&body); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser das senhas",
		})
		h.Logger.LogBadRequest(ctx, "UserHandler", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.ChangePasswordInput{
		UserID:      userID.(string),
		OldPassword: body.OldPassword,
		NewPassword: body.NewPassword,
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
// @Tags Administração
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
