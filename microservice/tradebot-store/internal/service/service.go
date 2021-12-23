package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	kube "github.com/108356037/v1/tradebot-store/kubeclient"
	"github.com/gin-gonic/gin"
)

// var statusCodeMap = map[string]int{
// 	"POST":   201,
// 	"GET":    200,
// 	"DELETE": 200,
// 	"UPDATE": 200,
// }

type ServiceLib struct{}

func New() ServiceLib {
	return ServiceLib{}
}

func (src ServiceLib) TradeBot(c *gin.Context) {
	uuid := c.Param("uuid")
	botname := c.Param("botname")
	method := c.Request.Method

	switch method {
	case "POST":
		jsonData, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{
				"status": "failed",
				"info":   err.Error(),
			})
			return
		}

		var dataHolder map[string]string
		err = json.Unmarshal(jsonData, &dataHolder)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{
				"status": "failed",
				"info":   err.Error(),
			})
			return
		}

		if err := kube.CreateBot(uuid, botname, dataHolder); err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"status": "failed",
				"info":   fmt.Sprintf("Create bot failed, error: %s", err.Error()),
			})
			return
		}
		c.JSON(201, map[string]string{
			"status": "success",
			"info":   fmt.Sprintf("Successfully created bot: %s", botname),
		})
		return

	case "DELETE":
		if err := kube.DeleteBot(uuid, botname); err != nil {
			c.JSON(http.StatusInternalServerError, map[string]string{
				"status": "failed",
				"info":   fmt.Sprintf("Delete bot failed, error: %s", err.Error()),
			})
			return
		}
		c.JSON(204, map[string]string{
			"status": "success",
			"info":   fmt.Sprintf("Successfully deleted bot: %s", botname),
		})
		return

	case "PUT":
		if err := kube.UpdateBot(uuid, botname); err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{
				"status": "failed",
				"info":   fmt.Sprintf("Update bot failed, error: %s", err.Error()),
			})
			return
		}
		c.JSON(200, map[string]string{
			"status": "success",
			"info":   fmt.Sprintf("Successfully updated bot: %s", botname),
		})
		return

	case "GET":
		if err := kube.GetSingleBot(uuid, botname); err != nil {
			c.JSON(400, map[string]string{
				"status": "failed",
				"info":   fmt.Sprintf("Cannot get bot: %s", err.Error()),
			})
			return
		}
		c.JSON(200, map[string]string{
			"status": "success",
			"info":   fmt.Sprintf("Successfully get bot: %s", botname),
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
