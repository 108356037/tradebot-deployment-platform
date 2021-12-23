package service

import (
	"fmt"
	"net/http"

	kube "github.com/108356037/v1/user-resource-svc/internal/kubeclient"
	"github.com/gin-gonic/gin"
)

func (src ServiceLib) GetResourceQuota(c *gin.Context) {
	uuid := c.Param("username")
	params := c.Request.URL.Query()
	var returnVal *kube.ResourceStatus
	var parseErr error

	if params.Get("int64") != "" {
		res, err := kube.GetNsResourceQuotaInt64(uuid)
		returnVal = res
		parseErr = err
	} else {
		res, err := kube.GetNsResourceQuota(uuid)
		returnVal = res
		parseErr = err
	}

	if parseErr != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"status": "failed",
			"info":   fmt.Sprintf("Failed to get resource quota, error: %s", parseErr.Error()),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"status": "success",
		"info":   *returnVal,
	})
}

func (src ServiceLib) ResourceRequestIsSufficient(c *gin.Context) {
	uuid := c.Param("username")
	params := c.Request.URL.Query()

	isValid := kube.ResourceRequestValidator(params.Get("cpu_req"), params.Get("mem_req"), uuid)
	if !isValid {
		c.JSON(http.StatusBadRequest, map[string]string{
			"status": "failed",
			"info":   "insufficient request resource",
		})
		return
	}

	c.JSON(200, map[string]string{
		"status": "success",
		"info":   "sufficient request resource",
	})
}

func (src ServiceLib) ResourceLimitIsSufficient(c *gin.Context) {
	uuid := c.Param("username")
	params := c.Request.URL.Query()

	isValid := kube.ResourceLimitValidator(params.Get("cpu_req"), params.Get("mem_req"), uuid)
	if !isValid {
		c.JSON(http.StatusBadRequest, map[string]string{
			"status": "failed",
			"info":   "insufficient limit resource",
		})
		return
	}

	c.JSON(200, map[string]string{
		"status": "success",
		"info":   "sufficient limit resource",
	})
}
