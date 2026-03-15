package handlers

import (
	"net/http"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PetHandler struct {
	PetFactory *factories.PetFactory
	Logger     logging.LoggerInterface
}

func NewPetHandler(factory *factories.PetFactory, logger logging.LoggerInterface) *PetHandler {
	return &PetHandler{
		PetFactory: factory,
		Logger:     logger,
	}
}

// AddPet godoc
// @Summary Adiciona um novo pet
// @Tags Pets
// @Accept json
// @Produce json
// @Param input body usecases.AddPetInputBody true "Dados do pet"
// @Success 201 {object} usecases.AddPetOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /pets [post]
func (h *PetHandler) AddPet(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.AddPet", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	var inputBody usecases.AddPetInputBody
	if err := c.ShouldBindJSON(&inputBody); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados do pet",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.AddPet", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.AddPetInput{
		UserID:    userID.(string),
		Name:      inputBody.Name,
		SpeciesID: inputBody.SpeciesID,
		Breed:     inputBody.Breed,
		Age:       inputBody.Age,
		Weight:    inputBody.Weight,
		Notes:     inputBody.Notes,
	}

	output, errs := h.PetFactory.AddPet.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}

// ListPets godoc
// @Summary Lista pets do usuário autenticado
// @Tags Pets
// @Accept json
// @Produce json
// @Param page query int false "Número da página"
// @Param page_size query int false "Itens por página"
// @Success 200 {object} usecases.ListPetsOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /pets [get]
func (h *PetHandler) ListPets(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.ListPets", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	page := 1
	pageSize := 10
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

	input := usecases.ListPetsInput{
		UserID:   userID.(string),
		Page:     page,
		PageSize: pageSize,
	}

	output, errs := h.PetFactory.ListPets.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// ListPetsByOwnerID godoc
// @Summary Lista pets por ID de owner
// @Tags Pets
// @Accept json
// @Produce json
// @Param user_id path string true "ID do owner"
// @Param page query int false "Número da página"
// @Param page_size query int false "Itens por página"
// @Success 200 {object} usecases.ListPetsOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /users/{user_id}/pets [get]
func (h *PetHandler) ListPetsByOwnerID(c *gin.Context) {
	ctx := c.Request.Context()

	requesterUserID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.ListPetsByOwnerID", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	ownerID := c.Param("user_id")
	if ownerID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do owner ausente",
			Detail: "O ID do owner é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.ListPetsByOwnerID", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	page := 1
	pageSize := 10
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

	input := usecases.ListPetsByOwnerIDInput{
		RequesterUserID: requesterUserID.(string),
		OwnerID:         ownerID,
		Page:            page,
		PageSize:        pageSize,
	}

	output, errs := h.PetFactory.ListPetsByOwner.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// GetPet godoc
// @Summary Obtém detalhes de um pet
// @Tags Pets
// @Accept json
// @Produce json
// @Param pet_id path string true "ID do pet"
// @Success 200 {object} usecases.GetPetOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /pets/{pet_id} [get]
func (h *PetHandler) GetPet(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.GetPet", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	petID := c.Param("pet_id")
	if petID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do pet ausente",
			Detail: "O ID do pet é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.GetPet", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.GetPetInput{
		UserID: userID.(string),
		PetID:  petID,
	}

	output, errs := h.PetFactory.GetPet.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// UpdatePet godoc
// @Summary Atualiza dados de um pet
// @Tags Pets
// @Accept json
// @Produce json
// @Param pet_id path string true "ID do pet"
// @Param input body usecases.UpdatePetInputBody true "Dados do pet"
// @Success 200 {object} usecases.UpdatePetOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /pets/{pet_id} [put]
func (h *PetHandler) UpdatePet(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.UpdatePet", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	petID := c.Param("pet_id")
	if petID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do pet ausente",
			Detail: "O ID do pet é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.UpdatePet", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	var inputBody usecases.UpdatePetInputBody
	if err := c.ShouldBindJSON(&inputBody); err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Erro ao fazer o parser",
			Detail: "Erro ao fazer o parser dos dados do pet",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.UpdatePet", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.UpdatePetInput{
		UserID:    userID.(string),
		PetID:     petID,
		Name:      inputBody.Name,
		SpeciesID: inputBody.SpeciesID,
		Breed:     inputBody.Breed,
		Age:       inputBody.Age,
		Weight:    inputBody.Weight,
		Notes:     inputBody.Notes,
	}

	output, errs := h.PetFactory.UpdatePet.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// DeletePet godoc
// @Summary Remove um pet
// @Tags Pets
// @Accept json
// @Produce json
// @Param pet_id path string true "ID do pet"
// @Success 200 {object} usecases.DeletePetOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /pets/{pet_id} [delete]
func (h *PetHandler) DeletePet(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.DeletePet", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	petID := c.Param("pet_id")
	if petID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do pet ausente",
			Detail: "O ID do pet é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.DeletePet", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.DeletePetInput{
		UserID: userID.(string),
		PetID:  petID,
	}

	output, errs := h.PetFactory.DeletePet.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}

// AddPetPhoto godoc
// @Summary Adiciona foto ao pet
// @Tags Pets
// @Accept multipart/form-data
// @Produce json
// @Param pet_id path string true "ID do pet"
// @Param file formData file true "Imagem"
// @Success 201 {object} usecases.AddPetPhotoOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /pets/{pet_id}/photos [post]
func (h *PetHandler) AddPetPhoto(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.AddPetPhoto", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	petID := c.Param("pet_id")
	if petID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do pet ausente",
			Detail: "O ID do pet é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.AddPetPhoto", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "Arquivo ausente",
			Detail: "A imagem é obrigatória",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.AddPetPhoto", problem.Detail, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	input := usecases.AddPetPhotoInput{
		UserID:      userID.(string),
		PetID:       petID,
		FileName:    header.Filename,
		ContentType: contentType,
		Size:        header.Size,
		Reader:      file,
	}

	output, errs := h.PetFactory.AddPetPhoto.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusCreated, output)
}

// DeletePetPhoto godoc
// @Summary Remove foto do pet
// @Tags Pets
// @Accept json
// @Produce json
// @Param pet_id path string true "ID do pet"
// @Param photo_id path string true "ID da foto"
// @Success 200 {object} usecases.DeletePetPhotoOutput
// @Failure 400 {object} exceptions.ProblemDetails
// @Failure 401 {object} exceptions.ProblemDetails
// @Failure 403 {object} exceptions.ProblemDetails
// @Failure 404 {object} exceptions.ProblemDetails
// @Failure 500 {object} exceptions.ProblemDetails
// @Security Bearer
// @Router /pets/{pet_id}/photos/{photo_id} [delete]
func (h *PetHandler) DeletePetPhoto(c *gin.Context) {
	ctx := c.Request.Context()

	userID, exists := c.Get("user_id")
	if !exists {
		problem := exceptions.NewProblemDetails(exceptions.Unauthorized, exceptions.ErrorMessage{
			Title:  "Usuário não autenticado",
			Detail: "Não foi possível obter o ID do usuário autenticado",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.DeletePetPhoto", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusUnauthorized, problem)
		return
	}

	petID := c.Param("pet_id")
	if petID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID do pet ausente",
			Detail: "O ID do pet é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.DeletePetPhoto", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	photoID := c.Param("photo_id")
	if photoID == "" {
		problem := exceptions.NewProblemDetails(exceptions.BadRequest, exceptions.ErrorMessage{
			Title:  "ID da foto ausente",
			Detail: "O ID da foto é obrigatório",
		})
		h.Logger.LogBadRequest(ctx, "PetHandler.DeletePetPhoto", problem.Detail, nil)
		c.AbortWithStatusJSON(http.StatusBadRequest, problem)
		return
	}

	input := usecases.DeletePetPhotoInput{
		UserID:  userID.(string),
		PetID:   petID,
		PhotoID: photoID,
	}

	output, errs := h.PetFactory.DeletePetPhoto.Execute(ctx, input)
	if len(errs) > 0 {
		exceptions.HandleErrors(c, errs)
		return
	}

	c.JSON(http.StatusOK, output)
}
