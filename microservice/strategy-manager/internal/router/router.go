package router

import (
	"github.com/108356037/v1/strategy-manager/internal/service"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	svc := service.New()

	apiv1 := r.Group("/v1")
	{
		apiv1.GET("/:username/strategies", svc.Strategies)
		apiv1.GET("/:username/strategies/:strategyName", svc.Strategy)
		apiv1.GET("/:username/strategies/id/:strategyId", svc.StrategyById)
		apiv1.DELETE("/:username/strategies/:strategyName", svc.DeleteStrategy)
		apiv1.PATCH("/:username/strategies/:strategyName/schedule", svc.ScheduleStrategy)
		//apiv1.PATCH("/:username/strategies/:strategyName/resources", svc.SetStrategyResources)
		apiv1.PATCH("/:username/strategies/:strategyName/resources/requests", svc.SetStrategyRequest)
		apiv1.PATCH("/:username/strategies/:strategyName/resources/limits", svc.SetStrategyLimit)
	}

	return r

}
