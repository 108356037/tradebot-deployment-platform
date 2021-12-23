package router

import (
	"github.com/108356037/v1/user-resource-svc/internal/service"
	"github.com/108356037/v1/user-resource-svc/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(middleware.LoggerMiddleware())
	svc := service.New()

	apiv1 := r.Group("/resources")
	{

		apiv1.POST("/:username/ns", svc.Namespace)
		apiv1.DELETE("/:username/ns", svc.Namespace)

		apiv1.POST("/:username/helm", svc.HelmDepls)
		apiv1.GET("/:username/helm", svc.HelmDepls)
		//apiv1.PATCH("/:username/all/:replicas")
		apiv1.DELETE("/:username/helm", svc.HelmDepls)

		apiv1.POST("/:username/jupyter", svc.ResourceJuypter)
		apiv1.DELETE("/:username/jupyter", svc.ResourceJuypter)
		//apiv1.GET("/jupyter/:username")
		apiv1.PATCH("/:username/jupyter/:replicas", svc.PatchResourceJupyter)

		apiv1.POST("/:username/grafana", svc.ResourceGrafana)
		apiv1.DELETE("/:username/grafana", svc.ResourceGrafana)
		//apiv1.GET("/grafana/:username")
		apiv1.PATCH("/:username/grafana/:replicas", svc.PatchResourceGrafana)

		apiv1.POST("/:username/c9", svc.ResourceC9)
		apiv1.DELETE("/:username/c9", svc.ResourceC9)

		apiv1.GET("/:username/resourcequota", svc.GetResourceQuota)
		apiv1.GET("/:username/validate/requests", svc.ResourceRequestIsSufficient)
		apiv1.GET("/:username/validate/limits", svc.ResourceLimitIsSufficient)
	}

	// for k8s livenessProbe
	r.GET("/healthz", svc.LivenessProbe)

	return r
}
