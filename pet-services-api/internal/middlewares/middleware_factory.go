package middlewares

import "github.com/gin-gonic/gin"

// MiddlewareFactory centraliza o acesso aos middlewares customizados
// Útil para injeção e testes, além de manter padrão de uso.
type MiddlewareFactory struct{}

func NewMiddlewareFactory() *MiddlewareFactory {
	return &MiddlewareFactory{}
}

func (f *MiddlewareFactory) AuthMiddleware() gin.HandlerFunc {
	return AuthMiddleware()
}

func (f *MiddlewareFactory) OwnerOnlyMiddleware() gin.HandlerFunc {
	return OwnerOnlyMiddleware()
}

func (f *MiddlewareFactory) ProviderOnlyMiddleware() gin.HandlerFunc {
	return ProviderOnlyMiddleware()
}

func (f *MiddlewareFactory) AdminOnlyMiddleware() gin.HandlerFunc {
	return AdminOnlyMiddleware()
}

func (f *MiddlewareFactory) DefaultRateLimitMiddleware() gin.HandlerFunc {
	return DefaultRateLimitMiddleware()
}

func (f *MiddlewareFactory) StrictRateLimitMiddleware() gin.HandlerFunc {
	return StrictRateLimitMiddleware()
}

func (f *MiddlewareFactory) RateLimitMiddleware(config RateLimiterConfig) gin.HandlerFunc {
	return RateLimitMiddleware(config)
}
