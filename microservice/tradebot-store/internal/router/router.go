package router

import (
	"github.com/108356037/v1/tradebot-store/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	svc := service.New()

	apiv1 := r.Group("/tradebot")
	{
		apiv1.GET("/:uuid/all")
		apiv1.POST("/:uuid/:botname", svc.TradeBot)
		apiv1.PUT("/:uuid/:botname", svc.TradeBot)
		apiv1.DELETE("/:uuid/:botname", svc.TradeBot)
		apiv1.GET("/:uuid/:botname/logs")
	}

	// for k8s livenessProbe
	r.GET("/healthz")

	return r
}
