package delivery

import (
	"go-backend/internal/common/middlewares"
	"go-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

type OrderDelivery struct {
	orderHandler   *handler.OrderHandler
	authMiddleware *middlewares.AuthMiddleware
}

func NewOrderDelivery(orderHandler *handler.OrderHandler, authMiddleware *middlewares.AuthMiddleware) *OrderDelivery {
	return &OrderDelivery{
		orderHandler:   orderHandler,
		authMiddleware: authMiddleware,
	}
}

func (d *OrderDelivery) RegisterRouter(apiGroup *gin.RouterGroup) {
	OrderGroup := apiGroup.Group("order")
	{
		OrderGroup.Use(d.authMiddleware.Protect)
		OrderGroup.GET("", d.orderHandler.FindAll)
		OrderGroup.POST("", d.orderHandler.Create)
	}
}
