package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"pet-services-api/internal/application/exceptions"
	"pet-services-api/internal/application/provider"
	domainprovider "pet-services-api/internal/domain/provider"
	domainuser "pet-services-api/internal/domain/user"
	"pet-services-api/internal/infra/factory"
)

// ProviderHandler expõe endpoints de prestadores.
type ProviderHandler struct {
	uc           factory.ProviderUseCases
	errorService *ErrorService
}

func NewProviderHandler(uc factory.ProviderUseCases, errorService *ErrorService) *ProviderHandler {
	return &ProviderHandler{uc: uc, errorService: errorService}
}

// RegisterProviderRoutes registra rotas autenticadas de prestadores.
func RegisterProviderRoutes(rg *gin.RouterGroup, h *ProviderHandler) {
	rg.POST("", h.Register)
	rg.PUT(":provider_id", h.UpdateProfile)
	rg.POST(":provider_id/services", h.AddService)
	rg.PUT(":provider_id/services", h.UpdateService)
	rg.DELETE(":provider_id/services", h.RemoveService)
	rg.POST(":provider_id/photos", h.AddPhoto)
	rg.DELETE(":provider_id/photos/:photo_id", h.RemovePhoto)
	// rg.POST(":provider_id/photos/reorder", h.ReorderPhotos) // Handler não implementado
	rg.PUT(":provider_id/working-hours", h.UpdateWorkingHours)
	rg.PUT(":provider_id/working-hours/day", h.UpdateDaySchedule)
	rg.POST(":provider_id/activate", h.Activate)
	rg.POST(":provider_id/deactivate", h.Deactivate)
	rg.POST(":provider_id/approve", h.Approve)
	rg.POST(":provider_id/reject", h.Reject)
}

// RegisterProviderPublicRoutes registra rotas públicas de prestadores.
func RegisterProviderPublicRoutes(rg *gin.RouterGroup, h *ProviderHandler) {
	rg.GET("/search/location", h.ListByLocation)
}

type registerProviderRequest struct {
	BusinessName string                    `json:"business_name" validate:"required,min=3,max=150"`
	Description  string                    `json:"description" validate:"required,min=10,max=1000"`
	Address      domainuser.Address        `json:"address" validate:"required"`
	Latitude     float64                   `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude    float64                   `json:"longitude" validate:"required,min=-180,max=180"`
	Services     []provider.ServiceInput   `json:"services" validate:"required,min=1"`
	PriceRange   domainprovider.PriceRange `json:"price_range" validate:"required"`
}

// Register cria perfil de prestador para o usuário autenticado.
// @Summary Register provider profile
// @Tags providers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body registerProviderRequest true "Provider payload"
// @Success 201 {object} provider.RegisterProviderOutput
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Router /providers [post]
func (h *ProviderHandler) Register(c *gin.Context) {
	userID, problems := extractUserIDProblems(c)
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req registerProviderRequest
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}

	out, problems := h.uc.Register.Execute(c.Request.Context(), provider.RegisterProviderInput{
		UserID:       userID,
		BusinessName: req.BusinessName,
		Description:  req.Description,
		Address:      req.Address,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Services:     req.Services,
		PriceRange:   req.PriceRange,
	})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "register_provider_failed", http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"provider_id":   out.ID,
		"business_name": out.BusinessName,
		"is_active":     out.IsActive,
	})
}

type updateProviderProfileRequest struct {
	BusinessName *string                    `json:"business_name" validate:"omitempty,min=3,max=150"`
	Description  *string                    `json:"description" validate:"omitempty,min=10,max=1000"`
	Address      *domainuser.Address        `json:"address"`
	Latitude     *float64                   `json:"latitude" validate:"omitempty,min=-90,max=90"`
	Longitude    *float64                   `json:"longitude" validate:"omitempty,min=-180,max=180"`
	PriceRange   *domainprovider.PriceRange `json:"price_range"`
}

// UpdateProfile atualiza dados do prestador.
// @Summary Update provider profile
// @Tags providers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param provider_id path string true "Provider ID"
// @Param request body updateProviderProfileRequest true "Profile payload"
// @Success 200 {object} domainprovider.Provider
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Router /providers/{provider_id} [put]
func (h *ProviderHandler) UpdateProfile(c *gin.Context) {
	providerID, problems := parseUUIDParamProblems(c, "provider_id", "invalid_provider_id")
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_provider_id", http.StatusBadRequest)
		return
	}

	userID, _ := extractUserID(c)

	var req updateProviderProfileRequest
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}

	_, problems = h.uc.UpdateProfile.Execute(c.Request.Context(), provider.UpdateProviderProfileInput{
		ProviderID:   providerID,
		UserID:       userID,
		BusinessName: req.BusinessName,
		Description:  req.Description,
		Address:      req.Address,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		PriceRange:   req.PriceRange,
	})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "update_provider_failed", http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

type serviceRequest struct {
	Category string  `json:"category" validate:"required,min=2,max=50"`
	Name     string  `json:"name" validate:"required,min=3,max=100"`
	PriceMin float64 `json:"price_min" validate:"required,min=0"`
	PriceMax float64 `json:"price_max" validate:"required,gtfield=PriceMin"`
}

// AddService adiciona um serviço ao prestador.
// @Summary Add service to provider
// @Tags providers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param provider_id path string true "Provider ID"
// @Param request body serviceRequest true "Service payload"
// @Success 200 {object} domainprovider.Provider
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Router /providers/{provider_id}/services [post]
func (h *ProviderHandler) AddService(c *gin.Context) {
	providerID, problems := parseUUIDParamProblems(c, "provider_id", "invalid_provider_id")
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_provider_id", http.StatusBadRequest)
		return
	}

	var req serviceRequest
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}

	_, problems = h.uc.AddService.Execute(c.Request.Context(), provider.AddServiceInput{
		ProviderID: providerID,
		Category:   req.Category,
		Name:       req.Name,
		PriceMin:   req.PriceMin,
		PriceMax:   req.PriceMax,
	})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "add_service_failed", http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// UpdateService atualiza um serviço do prestador.
// @Summary Update provider service
// @Tags providers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param provider_id path string true "Provider ID"
// @Param request body serviceRequest true "Service payload"
// @Success 200 {object} domainprovider.Provider
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Router /providers/{provider_id}/services [put]
func (h *ProviderHandler) UpdateService(c *gin.Context) {
	providerID, problems := parseUUIDParamProblems(c, "provider_id", "invalid_provider_id")
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_provider_id", http.StatusBadRequest)
		return
	}
	var req serviceRequest
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}

	_, problems = h.uc.UpdateService.Execute(c.Request.Context(), provider.UpdateServiceInput{
		ProviderID: providerID,
		Category:   req.Category,
		Name:       req.Name,
		PriceMin:   req.PriceMin,
		PriceMax:   req.PriceMax,
	})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "update_service_failed", http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// RemoveService remove um serviço do prestador.
// @Summary Remove provider service
// @Tags providers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param provider_id path string true "Provider ID"
// @Param request body object true "Service identifier"
// @Success 200 {object} domainprovider.Provider
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Router /providers/{provider_id}/services [delete]
func (h *ProviderHandler) RemoveService(c *gin.Context) {
	providerID, problems := parseUUIDParamProblems(c, "provider_id", "invalid_provider_id")
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_provider_id", http.StatusBadRequest)
		return
	}
	var req struct {
		Category string `json:"category"`
		Name     string `json:"name"`
	}
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}

	_, problems = h.uc.RemoveService.Execute(c.Request.Context(), provider.RemoveServiceInput{
		ProviderID: providerID,
		Category:   req.Category,
		Name:       req.Name,
	})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "remove_service_failed", http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// AddPhoto adiciona uma foto ao prestador.
// @Summary Add photo to provider
// @Tags providers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param provider_id path string true "Provider ID"
// @Param request body object true "Photo payload"
// @Success 200 {object} domainprovider.Provider
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Router /providers/{provider_id}/photos [post]
func (h *ProviderHandler) AddPhoto(c *gin.Context) {
	providerID, problems := parseUUIDParamProblems(c, "provider_id", "invalid_provider_id")
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_provider_id", http.StatusBadRequest)
		return
	}
	var req struct {
		URL string `json:"url"`
	}
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}

	_, problems = h.uc.AddPhoto.Execute(c.Request.Context(), provider.AddPhotoInput{ProviderID: providerID, URL: req.URL})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "add_photo_failed", http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// RemovePhoto remove uma foto do prestador.
// @Summary Remove photo from provider
// @Tags providers
// @Security BearerAuth
// @Produce json
// @Param provider_id path string true "Provider ID"
// @Param photo_id path string true "Photo ID"
// @Success 200 {object} domainprovider.Provider
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Router /providers/{provider_id}/photos/{photo_id} [delete]
func (h *ProviderHandler) RemovePhoto(c *gin.Context) {
	providerID, problems := parseUUIDParamProblems(c, "provider_id", "invalid_provider_id")
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_provider_id", http.StatusBadRequest)
		return
	}
	photoID, problems := parseUUIDParamProblems(c, "photo_id", "invalid_photo_id")
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_photo_id", http.StatusBadRequest)
		return
	}

	_, problems = h.uc.RemovePhoto.Execute(c.Request.Context(), provider.RemovePhotoInput{
		ProviderID: providerID,
		PhotoID:    photoID,
	})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "remove_photo_failed", http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// UpdateWorkingHours atualiza o horário de funcionamento do prestador.
// @Summary Update provider working hours
// @Tags providers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param provider_id path string true "Provider ID"
// @Param request body object true "Working hours payload"
// @Success 200 {object} domainprovider.Provider
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Router /providers/{provider_id}/working-hours [put]
func (h *ProviderHandler) UpdateWorkingHours(c *gin.Context) {
	providerID, problems := parseUUIDParamProblems(c, "provider_id", "invalid_provider_id")
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_provider_id", http.StatusBadRequest)
		return
	}
	var req struct {
		WorkingHours domainprovider.WorkingHours `json:"working_hours"`
	}
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}

	_, problems = h.uc.UpdateWorkingHours.Execute(c.Request.Context(), provider.UpdateWorkingHoursInput{
		ProviderID:   providerID,
		WorkingHours: req.WorkingHours,
	})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "update_working_hours_failed", http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// UpdateDaySchedule atualiza o horário de um dia específico do prestador.
// @Summary Update provider day schedule
// @Tags providers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param provider_id path string true "Provider ID"
// @Param request body object true "Day schedule payload"
// @Success 200 {object} domainprovider.Provider
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Router /providers/{provider_id}/working-hours/day [put]
func (h *ProviderHandler) UpdateDaySchedule(c *gin.Context) {
	providerID, problems := parseUUIDParamProblems(c, "provider_id", "invalid_provider_id")
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_provider_id", http.StatusBadRequest)
		return
	}
	var req struct {
		Day    int    `json:"day"`
		IsOpen bool   `json:"is_open"`
		Open   string `json:"open"`
		Close  string `json:"close"`
	}
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}

	weekday := time.Weekday(req.Day)

	_, problems = h.uc.UpdateDaySchedule.Execute(c.Request.Context(), provider.UpdateDayScheduleInput{
		ProviderID: providerID,
		Day:        weekday,
		IsOpen:     req.IsOpen,
		Open:       req.Open,
		Close:      req.Close,
	})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "update_day_schedule_failed", http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Activate ativa o perfil do prestador.
// @Summary Activate provider profile
// @Tags providers
// @Security BearerAuth
// @Produce json
// @Param provider_id path string true "Provider ID"
// @Success 200 {object} domainprovider.Provider
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Router /providers/{provider_id}/activate [post]
func (h *ProviderHandler) Activate(c *gin.Context) {
	providerID, problems := parseUUIDParamProblems(c, "provider_id", "invalid_provider_id")
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_provider_id", http.StatusBadRequest)
		return
	}
	_, problems = h.uc.Activate.Execute(c.Request.Context(), provider.ActivateProviderInput{ProviderID: providerID})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "activate_provider_failed", http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Deactivate desativa o perfil do prestador.
// @Summary Deactivate provider profile
// @Tags providers
// @Security BearerAuth
// @Produce json
// @Param provider_id path string true "Provider ID"
// @Success 200 {object} domainprovider.Provider
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Router /providers/{provider_id}/deactivate [post]
func (h *ProviderHandler) Deactivate(c *gin.Context) {
	providerID, problems := parseUUIDParamProblems(c, "provider_id", "invalid_provider_id")
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_provider_id", http.StatusBadRequest)
		return
	}
	_, problems = h.uc.Deactivate.Execute(c.Request.Context(), provider.DeactivateProviderInput{ProviderID: providerID})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "deactivate_provider_failed", http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Approve aprova o perfil do prestador.
// @Summary Approve provider profile
// @Tags providers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param provider_id path string true "Provider ID"
// @Param request body object true "Approval note"
// @Success 200 {object} domainprovider.Provider
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Router /providers/{provider_id}/approve [post]
func (h *ProviderHandler) Approve(c *gin.Context) {
	providerID, problems := parseUUIDParamProblems(c, "provider_id", "invalid_provider_id")
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_provider_id", http.StatusBadRequest)
		return
	}
	var req struct {
		Note string `json:"note"`
	}
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}
	_, problems = h.uc.Approve.Execute(c.Request.Context(), provider.ApproveProviderInput{ProviderID: providerID, Note: req.Note})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "approve_provider_failed", http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Reject rejeita o perfil do prestador.
// @Summary Reject provider profile
// @Tags providers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param provider_id path string true "Provider ID"
// @Param request body object true "Rejection reason"
// @Success 200 {object} domainprovider.Provider
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Router /providers/{provider_id}/reject [post]
func (h *ProviderHandler) Reject(c *gin.Context) {
	providerID, problems := parseUUIDParamProblems(c, "provider_id", "invalid_provider_id")
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_provider_id", http.StatusBadRequest)
		return
	}
	var req struct {
		Reason string `json:"reason"`
	}
	if problems := BindAndValidateJSONProblems(c, &req); len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "invalid_payload", http.StatusBadRequest)
		return
	}
	_, problems = h.uc.Reject.Execute(c.Request.Context(), provider.RejectProviderInput{ProviderID: providerID, Reason: req.Reason})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "reject_provider_failed", http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// ListProvidersByLocationResponseDTO representa a resposta da listagem de prestadores por localização.
type ListProvidersByLocationResponseDTO struct {
	Items []domainprovider.Provider `json:"items"`
	Total int64                     `json:"total"`
}

// ListByLocation retorna prestadores ativos próximos à localização.
// @Summary List providers by location
// @Tags providers
// @Produce json
// @Param lat query number true "Latitude"
// @Param lon query number true "Longitude"
// @Param radius_km query number true "Raio em km"
// @Param page query int false "Página"
// @Param limit query int false "Limite"
// @Success 200 {object} http.ListProvidersByLocationResponseDTO
// @Failure 400 {object} exceptions.ProblemDetailsOutputDTO
// @Router /providers/search/location [get]
func (h *ProviderHandler) ListByLocation(c *gin.Context) {
	lat, lon, radius := c.Query("lat"), c.Query("lon"), c.Query("radius_km")
	page := parseIntDefault(c.Query("page"), 1)
	limit := parseIntDefault(c.Query("limit"), 20)

	latF, lonF, radF, err := parseFloats(lat, lon, radius)
	if err != nil {
		problems := []exceptions.ProblemDetails{{
			Type:   string(exceptions.BadRequest),
			Title:  "Coordenadas inválidas",
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}}
		h.errorService.RespondWithProblems(c, problems, "invalid_coordinates", http.StatusBadRequest)
		return
	}

	providers, total, problems := h.uc.ListByLocation.Execute(c.Request.Context(), provider.ListProvidersByLocationInput{
		Latitude:  latF,
		Longitude: lonF,
		RadiusKM:  radF,
		Page:      page,
		Limit:     limit,
	})
	if len(problems) > 0 {
		h.errorService.RespondWithProblems(c, problems, "list_providers_failed", http.StatusBadRequest)
		return
	}

	resp := make([]gin.H, 0, len(providers))
	for _, p := range providers {
		resp = append(resp, gin.H{
			"id":              p.ID,
			"business_name":   p.BusinessName,
			"description":     p.Description,
			"price_range":     p.PriceRange,
			"avg_rating":      p.AvgRating,
			"total_reviews":   p.TotalReviews,
			"is_active":       p.IsActive,
			"approval_status": p.ApprovalStatus,
			"location": gin.H{
				"latitude":  p.Location.Latitude,
				"longitude": p.Location.Longitude,
				"address":   p.Location.Address,
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"items": resp,
		"total": total,
	})
}

// parseIntDefault converte string para int com default.
func parseIntDefault(v string, def int) int {
	if v == "" {
		return def
	}
	if n, err := strconv.Atoi(v); err == nil && n > 0 {
		return n
	}
	return def
}

// parseFloats converte strings para floats.
func parseFloats(lat, lon, radius string) (float64, float64, float64, error) {
	latF, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	lonF, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	radF, err := strconv.ParseFloat(radius, 64)
	if err != nil {
		return 0, 0, 0, err
	}
	return latF, lonF, radF, nil
}
