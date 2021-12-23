package service

import (
	"fmt"
	"net/http"

	helmcli "github.com/108356037/v1/user-resource-svc/internal/helm/cli"
	kube "github.com/108356037/v1/user-resource-svc/internal/kubeclient"
	"github.com/gin-gonic/gin"
	//log "github.com/sirupsen/logrus"
)

var statusCodeMap = map[string]int{
	"POST":   201,
	"GET":    200,
	"DELETE": 200,
	"UPDATE": 200,
}

type ServiceLib struct{}

func New() ServiceLib {
	return ServiceLib{}
}

func (src ServiceLib) Namespace(c *gin.Context) {
	uuid := c.Param("username")
	method := c.Request.Method

	switch method {
	case "POST":
		if err := kube.CreateNSWithRegcred(uuid); err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"status": "failed",
				"info":   fmt.Sprintf("Create ns failed, error: %s", err.Error()),
			})
			return
		}
		c.JSON(201, map[string]string{
			"status": "success",
			"info":   fmt.Sprintf("Successfully created ns: %s", uuid),
		})
		return
	case "DELETE":
		if err := kube.DeleteUserNamespace(uuid); err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"status": "failed",
				"info":   fmt.Sprintf("Delete ns failed, error: %s", err.Error()),
			})
			return
		}
		c.JSON(204, map[string]string{
			"status": "success",
			"info":   fmt.Sprintf("Successfully deleted ns: %s", uuid),
		})
		return
	default:
		c.JSON(http.StatusBadRequest, map[string]string{
			"status": "failed",
			"info":   fmt.Sprintf("Method %s denied\n", c.Request.Method),
		})
		return
	}
}

func (src ServiceLib) HelmDepls(c *gin.Context) {

	uuid := c.Param("username")
	method := c.Request.Method
	var result *helmcli.CmdResult
	var jsonResult map[string]interface{}

	switch method {
	case "POST":
		result = helmcli.CreateUserDepls(uuid)
		jsonResult = helmcli.ParseCreate(result)
	case "GET":
		result = helmcli.GetUserDepls(uuid)
		jsonResult = helmcli.ParseGet(result)
	case "DELETE":
		result = helmcli.DeleteUserDepls(uuid)
		jsonResult = helmcli.ParseDelete(result)
	default:
		c.JSON(http.StatusBadRequest, map[string]string{
			"status": "failed",
			"info":   fmt.Sprintf("Method %s denied\n", c.Request.Method),
		})
		return
	}

	if result.Stderr != "" {
		jsonResult["status"] = "failed"
		c.JSON(500, jsonResult)
		return
	}

	jsonResult["status"] = "success"
	c.JSON(statusCodeMap[method], jsonResult)
}

func (srv ServiceLib) LivenessProbe(c *gin.Context) {
	c.JSON(200, "server healthy")
}
