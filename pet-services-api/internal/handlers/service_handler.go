package handlers

import (
	"net/http"
	"strconv"

	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	ServiceFactory *factories.ServiceFactory
	Logger         logging.LoggerInterface
}

func NewServiceHandler(factory *factories.ServiceFactory, logger logging.LoggerInterface) *ServiceHandler {
	return &ServiceHandler{
		ServiceFactory: factory,
		Logger:         logger,
	}
}

// AddService godoc
// @Summary Adiciona um novo serviço
// @Tags Serviços
// @Accept json
// @Produce json
// @Param input body usecases.AddServiceInputBody true "Dados do serviço"
// @Success 201 {object} usecases.AddServiceOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /services [post]
// ListTags godoc
// @Summary Lista tags com paginação
// @Tags Tags
// @Accept json
// @Produce json
// @Param page query int false "Número da página"
// @Param page_size query int false "Itens por página"
// @Param name query string false "Filtro por nome"
// @Success 200 {object} usecases.ListTagsOutput
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /tags [get]
func (h *ServiceHandler) ListTags(c *gin.Context) {
	ctx := c.Request.Context()

	userType, exists := c.Get("user_type")
	if !exists || userType != "provider" {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Acesso negado",
			Detail: "Acesso permitido apenas para usuários do tipo provider",
		})
		h.Logger.LogError(ctx, "ServiceHandler.ListTags", problem.Title+": "+problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusForbidden, problem)
		return
	}

	page := 1
	pageSize := 10
	name := c.Query("name")
	if p := c.Query("page"); p != "" {
		if val, err := strconv.Atoi(p); err == nil && val > 0 {
			page = val
		}
	}
	if ps := c.Query("page_size"); ps != "" {
		if val, err := strconv.Atoi(ps); err == nil && val > 0 {
			pageSize = val
		}
	}

	providerID, _ := c.Get("user_id")
	input := usecases.ListTagsInput{
		Page:       page,
		PageSize:   pageSize,
		Name:       name,
		ProviderID: providerID.(string),
	}

	output, err := h.ServiceFactory.ListTags.Execute(ctx, input)
	if err != nil {
		problem := exceptions.NewProblemDetails(exceptions.InternalServerError, exceptions.ErrorMessage{
			Title:  "Erro ao listar tags",
			Detail: err.Error(),
		})
		h.Logger.LogError(ctx, "ServiceHandler.ListTags", problem.Title+": "+problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, problem)
		return
	}

	c.JSON(http.StatusOK, output)
}
func (h *ServiceHandler) AddService(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "ServiceHandler.AddService", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	var inputBody usecases.AddServiceInputBody
	if err := c.ShouldBindJSON(&inputBody); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados do serviço",
		})
		h.Logger.LogBadRequest(ctx, "ServiceHandler.AddService", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.AddServiceInput{
		UserID:       userID.(string),
		Name:         inputBody.Name,
		Description:  inputBody.Description,
		Price:        inputBody.Price,
		PriceMinimum: inputBody.PriceMinimum,
		PriceMaximum: inputBody.PriceMaximum,
		Duration:     inputBody.Duration,
	}

	output, errs := h.ServiceFactory.AddService.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}

// AddServicePhoto godoc
// @Summary Adiciona foto ao serviço
// @Tags Serviços
// @Accept multipart/form-data
// @Produce json
// @Param service_id path string true "ID do serviço"
// @Param file formData file true "Imagem"
// @Success 201 {object} usecases.AddServicePhotoOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /services/{service_id}/photos [post]
func (h *ServiceHandler) AddServicePhoto(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "ServiceHandler.AddServicePhoto", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	serviceID := c.Param("service_id")
	if serviceID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do serviço ausente",
			Detail: "O ID do serviço é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "ServiceHandler.AddServicePhoto", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Arquivo ausente",
			Detail: "A imagem é obrigatória",
		})
		h.Logger.LogBadRequest(ctx, "ServiceHandler.AddServicePhoto", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	input := usecases.AddServicePhotoInput{
		UserID:      userID.(string),
		ServiceID:   serviceID,
		FileName:    header.Filename,
		ContentType: contentType,
		Size:        header.Size,
		Reader:      file,
	}

	output, errs := h.ServiceFactory.AddServicePhoto.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}

// AddServiceTag godoc
// @Summary Adiciona tag ao serviço
// @Tags Serviços
// @Accept json
// @Produce json
// @Param service_id path string true "ID do serviço"
// @Param input body usecases.AddServiceTagInput true "Dados da tag"
// @Success 201 {object} usecases.AddServiceTagOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 409 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /services/{service_id}/tags [post]
func (h *ServiceHandler) AddServiceTag(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "ServiceHandler.AddServiceTag", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	serviceID := c.Param("service_id")
	if serviceID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do serviço ausente",
			Detail: "O ID do serviço é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "ServiceHandler.AddServiceTag", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	var inputBody usecases.AddServiceTagInput
	if err := c.ShouldBindJSON(&inputBody); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados da tag",
		})
		h.Logger.LogBadRequest(ctx, "ServiceHandler.AddServiceTag", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.AddServiceTagInput{
		UserID:    userID.(string),
		ServiceID: serviceID,
		TagID:     inputBody.TagID,
		TagName:   inputBody.TagName,
	}

	output, errs := h.ServiceFactory.AddServiceTag.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}

// AddServiceCategory godoc
// @Summary Adiciona categoria ao serviço
// @Tags Serviços
// @Accept json
// @Produce json
// @Param service_id path string true "ID do serviço"
// @Param input body usecases.AddServiceCategoryInput true "Dados da categoria"
// @Success 201 {object} usecases.AddServiceCategoryOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 409 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /services/{service_id}/categories [post]
func (h *ServiceHandler) AddServiceCategory(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "ServiceHandler.AddServiceCategory", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	serviceID := c.Param("service_id")
	if serviceID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do serviço ausente",
			Detail: "O ID do serviço é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "ServiceHandler.AddServiceCategory", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	var inputBody usecases.AddServiceCategoryInput
	if err := c.ShouldBindJSON(&inputBody); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados da categoria",
		})
		h.Logger.LogBadRequest(ctx, "ServiceHandler.AddServiceCategory", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.AddServiceCategoryInput{
		UserID:     userID.(string),
		ServiceID:  serviceID,
		CategoryID: inputBody.CategoryID,
	}

	output, errs := h.ServiceFactory.AddServiceCategory.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}
