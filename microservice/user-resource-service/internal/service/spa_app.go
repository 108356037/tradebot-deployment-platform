package service

import (
	"fmt"
	"net/http"

	helmcli "github.com/108356037/v1/user-resource-svc/internal/helm/cli"
	"github.com/gin-gonic/gin"
)

func (srv ServiceLib) ResourceJuypter(c *gin.Context) {

	uuid := c.Param("username")
	method := c.Request.Method
	var result *helmcli.CmdResult
	var jsonResult map[string]interface{}

	switch method {
	case "POST":
		result = helmcli.CreateResourceJupyter(uuid)
		jsonResult = helmcli.ParseCreate(result)
	case "DELETE":
		result = helmcli.DeleteResourceJupyter(uuid)
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

func (srv ServiceLib) ResourceGrafana(c *gin.Context) {

	uuid := c.Param("username")
	method := c.Request.Method
	var result *helmcli.CmdResult
	var jsonResult map[string]interface{}

	switch method {
	case "POST":
		result = helmcli.CreateResourceGrafana(uuid)
		jsonResult = helmcli.ParseCreate(result)
	case "DELETE":
		result = helmcli.DeleteResourceGrafana(uuid)
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

func (srv ServiceLib) ResourceC9(c *gin.Context) {
	uuid := c.Param("username")
	method := c.Request.Method
	var result *helmcli.CmdResult
	var jsonResult map[string]interface{}

	switch method {
	case "POST":
		result = helmcli.CreateResourceC9(uuid) //
		jsonResult = helmcli.ParseCreate(result)
	case "DELETE":
		result = helmcli.DeleteResourceC9(uuid) //
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

func (srv ServiceLib) PatchResourceJupyter(c *gin.Context) {

	uuid := c.Param("username")
	replicas := c.Param("replicas")
	result := helmcli.PatchResourceJupyter(uuid, replicas)

	if result.Stderr != "" {
		c.JSON(500, map[string]string{
			"status": "failed",
			"info":   result.Stderr,
		})
		return
	}

	c.JSON(200, map[string]string{
		"status": "success",
		"info":   result.Stdout,
	})
}

func (srv ServiceLib) PatchResourceGrafana(c *gin.Context) {

	uuid := c.Param("username")
	replicas := c.Param("replicas")
	result := helmcli.PatchResourceGrafana(uuid, replicas)

	if result.Stderr != "" {
		c.JSON(500, map[string]string{
			"status": "failed",
			"info":   result.Stderr,
		})
		return
	}

	c.JSON(200, map[string]string{
		"status": "success",
		"info":   result.Stdout,
	})
}
