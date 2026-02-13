package handlers

import (
	"net/http"
	"pet-services-api/internal/exceptions"
	"pet-services-api/internal/factories"
	"pet-services-api/internal/logging"
	"pet-services-api/internal/usecases"

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
