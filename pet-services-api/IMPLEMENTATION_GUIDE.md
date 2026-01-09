# 🎯 IMPLEMENTAÇÕES IMEDIATAS - MVP PRODUCTION

## 1️⃣ RATE LIMITING (1 hora)

```go
// internal/transport/http/middleware_ratelimit.go
package http

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimitMiddleware implementa rate limiting por IP/User
func RateLimitMiddleware(requestsPerMinute int) gin.HandlerFunc {
	limiters := make(map[string]*rate.Limiter)
	mu := sync.RWMutex{}

	return func(c *gin.Context) {
		// Usar user ID se autenticado, caso contrário usar IP
		key := c.ClientIP()
		if userID, ok := c.Get("user_id"); ok {
			key = fmt.Sprintf("user:%v", userID)
		}

		mu.RLock()
		limiter, exists := limiters[key]
		mu.RUnlock()

		if !exists {
			limiter = rate.NewLimiter(
				rate.Limit(float64(requestsPerMinute)/60),
				requestsPerMinute,
			)
			mu.Lock()
			limiters[key] = limiter
			mu.Unlock()
		}

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
				"retry_after": 60,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
```

## 2️⃣ INPUT VALIDATION (1 hora)

```go
// internal/transport/http/validators.go
package http

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateSignup valida signup payload
func ValidateSignup(req signupRequest) error {
	return validate.Struct(req)
	// Com tags: email `validate:"required,email"`
}

// Auth handler com validação
func (h *AuthHandler) Signup(c *gin.Context) {
	var req signupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_payload", err.Error()))
		return
	}

	// ADICIONAR: Validação estrutural
	if err := ValidateSignup(req); err != nil {
		validationErrors := extractValidationErrors(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "validation_failed",
			"details": validationErrors,
		})
		return
	}

	// ... rest of handler
}
```

## 3️⃣ TRANSAÇÕES DE BANCO (2 horas)

```go
// internal/infra/repository/gorm/request_repository.go

func (r *RequestRepository) CreateWithTransaction(ctx context.Context, req *requestdom.ServiceRequest) error {
	// Usar transação para operações múltiplas
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Verificar disponibilidade do provider
		var provider models.Provider
		if err := tx.First(&provider, "id = ? AND is_active = true", req.ProviderID).Error; err != nil {
			return requestdom.ErrProviderNotFound
		}

		// 2. Criar request
		m := toModelServiceRequest(req)
		if err := tx.Create(m).Error; err != nil {
			return err
		}

		// 3. Atualizar request ID
		req.ID = m.ID

		return nil
	})
}
```

## 4️⃣ SECURITY HEADERS (30 min)

```go
// Add ao cmd/api/main.go na função corsMiddleware()
func securityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		if os.Getenv("ENABLE_HSTS") == "true" {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		c.Next()
	}
}

// No main(), adicionar após outras middlewares:
router.Use(securityHeadersMiddleware())
```

## 5️⃣ SOFT DELETE FILTERING (1 hora)

```go
// internal/infra/repository/gorm/user_repository.go

// Adicionar scope helper
func notDeleted(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at IS NULL")
}

// Usar em todas as queries:
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*userdom.User, error) {
	var m models.User
	if err := r.db.WithContext(ctx).
		Scopes(notDeleted).  // ← ADICIONAR
		Where("email = ?", email).
		First(&m).Error; err != nil {
		return nil, err
	}
	return toDomainUser(&m)
}

// Aplicar para: FindByID, FindAll, List queries
```

## 6️⃣ FILE UPLOAD PARA PROVIDER PHOTOS (2-3 horas)

```go
// internal/transport/http/provider_handler.go

// AddPhoto upload foto do provider
func (h *ProviderHandler) AddPhoto(c *gin.Context) {
	providerID, err := uuid.Parse(c.Param("provider_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("invalid_provider_id", err.Error()))
		return
	}

	// Verificar autorização (seu próprio provider)
	userID, ok := userIDFromContext(c)
	if !ok {
		return
	}

	// Validar provider ownership
	provider, err := h.providerRepo.FindByID(c.Request.Context(), providerID)
	if err != nil || provider.UserID != userID {
		c.JSON(http.StatusForbidden, errorResponse("access_denied", ""))
		return
	}

	// Limpar arquivo
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse("missing_photo", err.Error()))
		return
	}

	// Validar
	if file.Size > 5*1024*1024 { // 5MB
		c.JSON(http.StatusBadRequest, errorResponse("file_too_large", "max 5MB"))
		return
	}

	// Upload para S3/Storage
	photoURL, err := uploadToStorage(file, providerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("upload_failed", err.Error()))
		return
	}

	// Salvar no banco
	photo := domainprovider.Photo{
		ID:  uuid.New(),
		URL: photoURL,
		Order: 0,
		CreatedAt: time.Now(),
	}

	provider.Photos = append(provider.Photos, photo)
	if err := h.providerRepo.Update(c.Request.Context(), provider); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("save_photo_failed", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"photo_url": photoURL})
}

// Helper: upload para storage backend
func uploadToStorage(file *multipart.FileHeader, providerID uuid.UUID) (string, error) {
	// TODO: Implementar com S3 ou Azure Blob Storage
	// Retornar URL public
	return "", nil
}
```

## 7️⃣ ESTRUTURA DE ERROS MELHORADA (1 hora)

```go
// internal/transport/http/error_response.go

type ErrorResponse struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message,omitempty"`
	Details []ValidationErrorDetail `json:"details,omitempty"`
	TraceID string                 `json:"trace_id"`
}

type ValidationErrorDetail struct {
	Field string `json:"field"`
	Rule  string `json:"rule"`
	Value string `json:"value"`
}

func errorResponse(code string, message string) gin.H {
	return gin.H{
		"error":   code,
		"message": message,
	}
}

func validationErrorResponse(c *gin.Context, errors []ValidationErrorDetail) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error:   "validation_failed",
		Details: errors,
		TraceID: c.GetString("request_id"),
	})
}
```

## 8️⃣ ADMIN ENDPOINTS BÁSICOS (2 horas)

```go
// internal/transport/http/admin_handler.go

// AdminHandler operações administrativas
type AdminHandler struct {
	userRepo     domain.UserRepository
	providerRepo domain.ProviderRepository
}

// ListUsers lista usuários (admin only)
func (h *AdminHandler) ListUsers(c *gin.Context) {
	// TODO: Verificar se user é admin (adicionar role)

	page := parseIntDefault(c.Query("page"), 1)
	limit := parseIntDefault(c.Query("limit"), 20)

	users, total, err := h.userRepo.List(c.Request.Context(), page, limit)
	// ...
}

// DeleteUser hard delete (admin only)
func (h *AdminHandler) DeleteUser(c *gin.Context) {
	userID, _ := uuid.Parse(c.Param("user_id"))
	if err := h.userRepo.HardDelete(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse("delete_failed", err.Error()))
		return
	}
	c.Status(http.StatusNoContent)
}
```

## 9️⃣ UNIT TESTS BÁSICOS (3 horas)

```go
// internal/application/auth/login_test.go

package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"pet-services-api/internal/domain/auth"
	"pet-services-api/internal/domain/user"
)

// MockUserRepository para testes
type MockUserRepository struct {
	users map[string]*user.User
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	u, ok := m.users[email]
	if !ok {
		return nil, user.ErrUserNotFound
	}
	return u, nil
}

func TestLogin_Success(t *testing.T) {
	mockRepo := &MockUserRepository{
		users: map[string]*user.User{
			"test@example.com": {
				ID:       uuid.New(),
				Email:    "test@example.com",
				Password: "hashed_password",
				Type:     user.UserTypeOwner,
			},
		},
	}

	// Arrange
	uc := NewLoginUseCase(mockRepo, nil, nil)
	input := LoginInput{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Act
	_, err := uc.Execute(context.Background(), input)

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mockRepo := &MockUserRepository{users: map[string]*user.User{}}

	uc := NewLoginUseCase(mockRepo, nil, nil)
	input := LoginInput{
		Email:    "nonexistent@example.com",
		Password: "password",
	}

	_, err := uc.Execute(context.Background(), input)

	if !errors.Is(err, auth.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}
```

## 🔟 ENVIRONMENT VARIABLES COMPLETO (30 min)

```bash
# .env.production
# ─────────────────────────────────────────────────────

# Database
DATABASE_URL=postgres://user:pass@prod-db.example.com:5432/pet_services?sslmode=require

# Server
PORT=8080
GIN_MODE=release

# JWT
JWT_ACCESS_SECRET=<random-256-bit-secret>
JWT_REFRESH_SECRET=<random-256-bit-secret>
JWT_ACCESS_DURATION=15m
JWT_REFRESH_DURATION=7d

# Email
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASS=<sendgrid-key>
SMTP_FROM=noreply@petservices.com

# URLs
PASSWORD_RESET_BASE_URL=https://app.petservices.com/reset-password
EMAIL_VERIFY_BASE_URL=https://app.petservices.com/verify-email

# CORS
CORS_ORIGINS=https://app.petservices.com,https://www.petservices.com

# Security
ENABLE_HSTS=true
```

---

## 📅 CRONOGRAMA SUGERIDO

| Tarefa                | Tempo | Prioridade |
| --------------------- | ----- | ---------- |
| Rate limiting         | 1h    | ⭐⭐⭐     |
| Input validation      | 1h    | ⭐⭐⭐     |
| Security headers      | 30m   | ⭐⭐⭐     |
| Transações            | 2h    | ⭐⭐⭐     |
| Soft delete filtering | 1h    | ⭐⭐       |
| File upload           | 3h    | ⭐⭐       |
| Error responses       | 1h    | ⭐⭐       |
| Admin endpoints       | 2h    | ⭐⭐       |
| Unit tests            | 3h    | ⭐⭐       |

**Total: ~15 horas = 2 dias com 1 dev**

---

## ✅ VALIDAÇÃO PRÉ-DEPLOY

Antes de fazer deploy:

```bash
# 1. Tests passando
make test

# 2. Lint clean
make lint

# 3. Build sem warnings
make build

# 4. Migrations rodando
make migrate

# 5. Endpoints testados manualmente
curl http://localhost:8080/health
curl http://localhost:8080/ready

# 6. Documentação atualizada
# Verificar MVP_ANALYSIS.md

# 7. Environment variables configuradas
cat .env.production

# 8. Backup database configurado
# Verificar backup scripts

# 9. Monitoring setup
# Prometheus, logs centralizados, alertas

# 10. Load testing
# Hey, wrk, ou locust
```

---

## 🚀 GO LIVE CHECKLIST

- [ ] Rate limiting ativado
- [ ] Validação de entrada 100%
- [ ] HTTPS/TLS configurado
- [ ] Security headers ativados
- [ ] SMTP real funcionando
- [ ] Transações multi-tabela implementadas
- [ ] Soft deletes funcionando
- [ ] File uploads testado
- [ ] Admin API ativo
- [ ] Tests passando (>70% coverage)
- [ ] Backup automático
- [ ] Monitoring ativo
- [ ] Logs centralizados
- [ ] Alertas configurados
- [ ] Runbooks criados
- [ ] On-call rotation setup
