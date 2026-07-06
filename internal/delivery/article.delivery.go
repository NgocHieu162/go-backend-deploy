package delivery

import (
	"go-backend/internal/common/cache"
	"go-backend/internal/common/middlewares"
	"go-backend/internal/handler"
	"time"

	"github.com/gin-gonic/gin"
)

type articleDelivery struct {
	articleHandler *handler.ArticleHandler
	cache          *cache.Cache
	authMiddleware *middlewares.AuthMiddleware
}

func NewArticleDelivery(articleHandler *handler.ArticleHandler, cache *cache.Cache, authMiddleware *middlewares.AuthMiddleware) *articleDelivery {
	return &articleDelivery{
		articleHandler: articleHandler,
		cache:          cache,
		authMiddleware: authMiddleware,
	}
}

func (a *articleDelivery) RegisterRouter(apiGroup *gin.RouterGroup) {
	articleGroup := apiGroup.Group("article")
	{
		articleGroup.Use(a.authMiddleware.Protect)
		articleGroup.POST("", a.articleHandler.Create)
		articleGroup.GET("", a.cache.RedisCache(5*time.Second), a.articleHandler.FindAll)
		articleGroup.DELETE(":id", a.articleHandler.Delete)
	}
}
