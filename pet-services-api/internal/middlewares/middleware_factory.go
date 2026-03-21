package middlewares

import (
	"pet-services-api/internal/entities"
	"pet-services-api/internal/logging"

	"github.com/gin-gonic/gin"
)

type MiddlewareFactory struct {
	Logger logging.LoggerInterface
}

func NewMiddlewareFactory(logger logging.LoggerInterface) *MiddlewareFactory {
	return &MiddlewareFactory{Logger: logger}
}

func (f *MiddlewareFactory) AuthMiddleware() gin.HandlerFunc {
	return AuthMiddleware(f.Logger)
}

func (f *MiddlewareFactory) OwnerOnlyMiddleware() gin.HandlerFunc {
	return OwnerOnlyMiddleware(f.Logger)
}

func (f *MiddlewareFactory) ProviderOnlyMiddleware() gin.HandlerFunc {
	return ProviderOnlyMiddleware(f.Logger)
}

func (f *MiddlewareFactory) AdminOnlyMiddleware() gin.HandlerFunc {
	return AdminOnlyMiddleware(f.Logger)
}

func (f *MiddlewareFactory) DefaultRateLimitMiddleware() gin.HandlerFunc {
	return DefaultRateLimitMiddleware(f.Logger)
}

func (f *MiddlewareFactory) StrictRateLimitMiddleware() gin.HandlerFunc {
	return StrictRateLimitMiddleware(f.Logger)
}

func (f *MiddlewareFactory) RateLimitMiddleware(config RateLimiterConfig) gin.HandlerFunc {
	return RateLimitMiddleware(config, f.Logger)
}

func (f *MiddlewareFactory) AdoptionGuardianApprovedMiddleware(guardianRepo entities.AdoptionGuardianProfileRepository) gin.HandlerFunc {
	return AdoptionGuardianApprovedMiddleware(f.Logger, guardianRepo)
}
