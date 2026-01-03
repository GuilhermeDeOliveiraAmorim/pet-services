package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/guilherme/pet-services-api/internal/application/provider"
	domainprovider "github.com/guilherme/pet-services-api/internal/domain/provider"
	domainuser "github.com/guilherme/pet-services-api/internal/domain/user"
	"github.com/guilherme/pet-services-api/internal/infra/factory"
)

// ProviderHandler expõe endpoints de prestadores.
type ProviderHandler struct {
	uc factory.ProviderUseCases
}

func NewProviderHandler(uc factory.ProviderUseCases) *ProviderHandler {
	return &ProviderHandler{uc: uc}
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
	rg.POST(":provider_id/photos/reorder", h.ReorderPhotos)
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
// @Summary Create provider profile
// @Tags providers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object true "Provider payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /providers [post]
func (h *ProviderHandler) Register(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse("unauthorized", err.Error()))
		return
	}

	var req registerProviderRequest
	if err := BindAndValidateJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid_payload",
			"fields": ValidationErrorResponse(err),
		})
		return
	}

	out, err := h.uc.Register.Execute(c.Request.Context(), provider.RegisterProviderInput{
		UserID:       userID,
		BusinessName: req.BusinessName,
		Description:  req.Description,
		Address:      req.Address,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		Services:     req.Services,
		PriceRange:   req.PriceRange,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrNoServicesProvided), errors.Is(err, domainprovider.ErrInvalidLocation), errors.Is(err, domainprovider.ErrInvalidPriceRange):
			c.JSON(http.StatusBadRequest, errorResponse("invalid_data", err.Error()))
		case errors.Is(err, domainprovider.ErrProviderAlreadyExists):
			c.JSON(http.StatusConflict, errorResponse("provider_exists", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("register_provider_failed", err.Error()))
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"provider_id":   out.ProviderID,
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
// @Param request body object true "Profile payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /providers/{provider_id} [put]
func (h *ProviderHandler) UpdateProfile(c *gin.Context) {
	providerID, err := uuid.Parse(c.Param("provider_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_provider_id", err.Error()))
		return
	}

	userID, _ := extractUserID(c) // autorização simples, se presente

	var req updateProviderProfileRequest
	if err := BindAndValidateJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid_payload",
			"fields": ValidationErrorResponse(err),
		})
		return
	}

	if err := h.uc.UpdateProfile.Execute(c.Request.Context(), provider.UpdateProviderProfileInput{
		ProviderID:   providerID,
		UserID:       userID,
		BusinessName: req.BusinessName,
		Description:  req.Description,
		Address:      req.Address,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		PriceRange:   req.PriceRange,
	}); err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrProviderNotFound):
			c.JSON(http.StatusNotFound, errorResponse("provider_not_found", err.Error()))
		case errors.Is(err, domainprovider.ErrInvalidPriceRange), errors.Is(err, domainprovider.ErrInvalidLocation):
			c.JSON(http.StatusBadRequest, errorResponse("invalid_data", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("update_provider_failed", err.Error()))
		}
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

func (h *ProviderHandler) AddService(c *gin.Context) {
	providerID, err := uuid.Parse(c.Param("provider_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_provider_id", err.Error()))
		return
	}

	var req serviceRequest
	if err := BindAndValidateJSON(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "invalid_payload",
			"fields": ValidationErrorResponse(err),
		})
		return
	}

	if err := h.uc.AddService.Execute(c.Request.Context(), provider.AddServiceInput{
		ProviderID: providerID,
		Category:   req.Category,
		Name:       req.Name,
		PriceMin:   req.PriceMin,
		PriceMax:   req.PriceMax,
	}); err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrProviderNotFound), errors.Is(err, domainprovider.ErrServiceNotFound):
			c.JSON(http.StatusNotFound, errorResponse("provider_or_service_not_found", err.Error()))
		default:
			c.JSON(http.StatusBadRequest, errorResponse("add_service_failed", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ProviderHandler) UpdateService(c *gin.Context) {
	providerID, err := uuid.Parse(c.Param("provider_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_provider_id", err.Error()))
		return
	}
	var req serviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
		return
	}

	if err := h.uc.UpdateService.Execute(c.Request.Context(), provider.UpdateServiceInput{
		ProviderID: providerID,
		Category:   req.Category,
		Name:       req.Name,
		PriceMin:   req.PriceMin,
		PriceMax:   req.PriceMax,
	}); err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrProviderNotFound), errors.Is(err, domainprovider.ErrServiceNotFound):
			c.JSON(http.StatusNotFound, errorResponse("provider_or_service_not_found", err.Error()))
		default:
			c.JSON(http.StatusBadRequest, errorResponse("update_service_failed", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ProviderHandler) RemoveService(c *gin.Context) {
	providerID, err := uuid.Parse(c.Param("provider_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_provider_id", err.Error()))
		return
	}
	var req struct {
		Category string `json:"category"`
		Name     string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
		return
	}

	if err := h.uc.RemoveService.Execute(c.Request.Context(), provider.RemoveServiceInput{
		ProviderID: providerID,
		Category:   req.Category,
		Name:       req.Name,
	}); err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrProviderNotFound), errors.Is(err, domainprovider.ErrServiceNotFound):
			c.JSON(http.StatusNotFound, errorResponse("provider_or_service_not_found", err.Error()))
		default:
			c.JSON(http.StatusBadRequest, errorResponse("remove_service_failed", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ProviderHandler) AddPhoto(c *gin.Context) {
	providerID, err := uuid.Parse(c.Param("provider_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_provider_id", err.Error()))
		return
	}
	var req struct {
		URL string `json:"url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
		return
	}

	if err := h.uc.AddPhoto.Execute(c.Request.Context(), provider.AddPhotoInput{ProviderID: providerID, URL: req.URL}); err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrProviderNotFound):
			c.JSON(http.StatusNotFound, errorResponse("provider_not_found", err.Error()))
		case errors.Is(err, domainprovider.ErrMaxPhotosReached):
			c.JSON(http.StatusBadRequest, errorResponse("max_photos_reached", err.Error()))
		default:
			c.JSON(http.StatusBadRequest, errorResponse("add_photo_failed", err.Error()))
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ProviderHandler) RemovePhoto(c *gin.Context) {
	providerID, err := uuid.Parse(c.Param("provider_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_provider_id", err.Error()))
		return
	}
	photoID, err := uuid.Parse(c.Param("photo_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_photo_id", err.Error()))
		return
	}

	if err := h.uc.RemovePhoto.Execute(c.Request.Context(), provider.RemovePhotoInput{ProviderID: providerID, PhotoID: photoID}); err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrProviderNotFound), errors.Is(err, domainprovider.ErrPhotoNotFound):
			c.JSON(http.StatusNotFound, errorResponse("provider_or_photo_not_found", err.Error()))
		default:
			c.JSON(http.StatusBadRequest, errorResponse("remove_photo_failed", err.Error()))
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ProviderHandler) ReorderPhotos(c *gin.Context) {
	providerID, err := uuid.Parse(c.Param("provider_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_provider_id", err.Error()))
		return
	}
	var req struct {
		Order []uuid.UUID `json:"order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
		return
	}

	if err := h.uc.ReorderPhotos.Execute(c.Request.Context(), provider.ReorderPhotosInput{ProviderID: providerID, Order: req.Order}); err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrProviderNotFound), errors.Is(err, domainprovider.ErrPhotoNotFound):
			c.JSON(http.StatusNotFound, errorResponse("provider_or_photo_not_found", err.Error()))
		default:
			c.JSON(http.StatusBadRequest, errorResponse("reorder_photos_failed", err.Error()))
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ProviderHandler) UpdateWorkingHours(c *gin.Context) {
	providerID, err := uuid.Parse(c.Param("provider_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_provider_id", err.Error()))
		return
	}
	var req struct {
		WorkingHours domainprovider.WorkingHours `json:"working_hours"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
		return
	}

	if err := h.uc.UpdateWorkingHours.Execute(c.Request.Context(), provider.UpdateWorkingHoursInput{
		ProviderID:   providerID,
		WorkingHours: req.WorkingHours,
	}); err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrProviderNotFound), errors.Is(err, domainprovider.ErrInvalidWorkingHours):
			c.JSON(http.StatusBadRequest, errorResponse("invalid_working_hours", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("update_working_hours_failed", err.Error()))
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ProviderHandler) UpdateDaySchedule(c *gin.Context) {
	providerID, err := uuid.Parse(c.Param("provider_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_provider_id", err.Error()))
		return
	}
	var req struct {
		Day    int    `json:"day"`
		IsOpen bool   `json:"is_open"`
		Open   string `json:"open"`
		Close  string `json:"close"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
		return
	}

	weekday := time.Weekday(req.Day)
	schedule := domainprovider.DaySchedule{IsOpen: req.IsOpen, Open: req.Open, Close: req.Close}

	if err := h.uc.UpdateDaySchedule.Execute(c.Request.Context(), provider.UpdateDayScheduleInput{
		ProviderID: providerID,
		Day:        weekday,
		IsOpen:     req.IsOpen,
		Open:       req.Open,
		Close:      req.Close,
	}); err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrProviderNotFound), errors.Is(err, domainprovider.ErrInvalidWorkingHours):
			c.JSON(http.StatusBadRequest, errorResponse("invalid_day_schedule", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("update_day_schedule_failed", err.Error()))
		}
		return
	}
	_ = schedule // keep reference for future logging; avoid unused warning if we change signature
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ProviderHandler) Activate(c *gin.Context) {
	providerID, err := uuid.Parse(c.Param("provider_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_provider_id", err.Error()))
		return
	}
	if err := h.uc.Activate.Execute(c.Request.Context(), provider.ActivateProviderInput{ProviderID: providerID}); err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrProviderNotFound):
			c.JSON(http.StatusNotFound, errorResponse("provider_not_found", err.Error()))
		default:
			c.JSON(http.StatusBadRequest, errorResponse("activate_provider_failed", err.Error()))
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ProviderHandler) Deactivate(c *gin.Context) {
	providerID, err := uuid.Parse(c.Param("provider_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_provider_id", err.Error()))
		return
	}
	if err := h.uc.Deactivate.Execute(c.Request.Context(), provider.DeactivateProviderInput{ProviderID: providerID}); err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrProviderNotFound):
			c.JSON(http.StatusNotFound, errorResponse("provider_not_found", err.Error()))
		default:
			c.JSON(http.StatusBadRequest, errorResponse("deactivate_provider_failed", err.Error()))
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ProviderHandler) Approve(c *gin.Context) {
	providerID, err := uuid.Parse(c.Param("provider_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_provider_id", err.Error()))
		return
	}
	var req struct {
		Note string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
		return
	}
	if err := h.uc.Approve.Execute(c.Request.Context(), provider.ApproveProviderInput{ProviderID: providerID, Note: req.Note}); err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrProviderNotFound):
			c.JSON(http.StatusNotFound, errorResponse("provider_not_found", err.Error()))
		default:
			c.JSON(http.StatusBadRequest, errorResponse("approve_provider_failed", err.Error()))
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *ProviderHandler) Reject(c *gin.Context) {
	providerID, err := uuid.Parse(c.Param("provider_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_provider_id", err.Error()))
		return
	}
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
		return
	}
	if err := h.uc.Reject.Execute(c.Request.Context(), provider.RejectProviderInput{ProviderID: providerID, Reason: req.Reason}); err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrProviderNotFound):
			c.JSON(http.StatusNotFound, errorResponse("provider_not_found", err.Error()))
		default:
			c.JSON(http.StatusBadRequest, errorResponse("reject_provider_failed", err.Error()))
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
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
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /providers/search/location [get]
func (h *ProviderHandler) ListByLocation(c *gin.Context) {
	lat, lon, radius := c.Query("lat"), c.Query("lon"), c.Query("radius_km")
	page := parseIntDefault(c.Query("page"), 1)
	limit := parseIntDefault(c.Query("limit"), 20)

	latF, lonF, radF, err := parseFloats(lat, lon, radius)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_coordinates", err.Error()))
		return
	}

	providers, total, err := h.uc.ListByLocation.Execute(c.Request.Context(), provider.ListProvidersByLocationInput{
		Latitude:  latF,
		Longitude: lonF,
		RadiusKM:  radF,
		Page:      page,
		Limit:     limit,
	})
	if err != nil {
		switch {
		case errors.Is(err, domainprovider.ErrInvalidLocation):
			c.JSON(http.StatusBadRequest, errorResponse("invalid_location", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, errorResponse("list_providers_failed", err.Error()))
		}
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
