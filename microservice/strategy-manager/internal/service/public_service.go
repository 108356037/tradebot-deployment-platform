package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/108356037/v1/strategy-manager/external"
	"github.com/108356037/v1/strategy-manager/internal/grpcclient"
	"github.com/108356037/v1/strategy-manager/internal/kubeclient"
	"github.com/108356037/v1/strategy-manager/internal/models"
	"github.com/108356037/v1/strategy-manager/mq"
	pb "github.com/108356037/v1/strategy-manager/proto"
	"github.com/adhocore/gronx"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/resource"
)

type ServiceLib struct{}

func New() ServiceLib {
	return ServiceLib{}
}

// List all strategies under user namespace
func (src ServiceLib) Strategies(c *gin.Context) {

	uuid := c.Param("username")
	strategySlice := make([]models.StrategyDoc, 0)
	//method := c.Request.Method

	strategyList, err := models.ListUserStrategies(uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"status": "failed",
			"info":   err.Error(),
		})
		return
	}

	strategySlice = append(strategySlice, *strategyList...)

	c.JSON(200, map[string]interface{}{
		"status": "success",
		"info":   strategySlice,
	})
}

// Get a specific strategy under user namespace with strategy id
func (src ServiceLib) StrategyById(c *gin.Context) {
	uuid := c.Param("username")
	strategyId := c.Param("strategyId")

	result, err := models.GetSingleStrategyById(uuid, strategyId)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"status": "failed",
			"info":   err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"status": "success",
		"info":   result,
	})
}

// List a specific strategy under user namespace
func (src ServiceLib) Strategy(c *gin.Context) {
	uuid := c.Param("username")
	strategyName := c.Param("strategyName")

	result, err := models.GetSingleStrategy(uuid, strategyName)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"status": "failed",
			"info":   err.Error(),
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"status": "success",
		"info":   result,
	})
}

// Delete a specfic strategy under user namespace
func (src ServiceLib) DeleteStrategy(c *gin.Context) {
	uuid := c.Param("username")
	strategyName := c.Param("strategyName")

	client := pb.NewRemoveClient(grpcclient.GrpcConn)
	grpcCallResult, err := client.RemoveFunc(
		context.Background(),
		&pb.RemoveReq{
			FuncName: strategyName,
			UserNS:   uuid,
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"status": "failed",
			"info":   err.Error(),
		})
		return
	}

	if grpcCallResult.Code != pb.StatusCode_Ok {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"status": "failed",
			"info":   grpcCallResult.Message,
		})
		return
	}

	c.JSON(202, map[string]interface{}{
		"status": "success",
		"info":   grpcCallResult.Message,
	})
}

// Schedule a specific user strategy
func (src ServiceLib) ScheduleStrategy(c *gin.Context) {

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

	uuid := c.Param("username")
	strategyName := c.Param("strategyName")
	crontabSechdule := dataHolder["schedule"]
	gron := gronx.New()
	if !gron.IsValid(crontabSechdule) {
		c.JSON(http.StatusBadRequest, map[string]string{
			"status": "failed",
			"info":   "Invalid Crontab syntax",
		})
		return
	}

	err = kubeclient.ScheduleStrategy(uuid, strategyName, crontabSechdule)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"status": "failed",
			"info":   err.Error(),
		})
		return
	}

	resourceEventInfo := map[string]string{
		"strategy": strategyName,
	}
	updateInfo := map[string]interface{}{
		"schedule": crontabSechdule,
	}
	mq.PublishUpdateEvent(uuid, mq.Strategy, resourceEventInfo, updateInfo)

	c.JSON(202, map[string]string{
		"status": "success",
		"info":   fmt.Sprintf("Scheduled strategy %s at %s", strategyName, crontabSechdule),
	})
}

// Set the cpu/memory request for strategy pod
func (src ServiceLib) SetStrategyRequest(c *gin.Context) {
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

	uuid := c.Param("username")
	strategyName := c.Param("strategyName")

	cpuReq := dataHolder["cpu"]
	memReq := dataHolder["mem"]
	if cpuReq == "" && memReq == "" {
		c.JSON(http.StatusBadRequest, map[string]string{
			"status": "failed",
			"info":   "please ensure cpu/mem body is included",
		})
		return
	}

	cpuReqVal := resource.MustParse(cpuReq)
	memReqVal := resource.MustParse(memReq)

	cpuCurrUsage, memCurrUsage := models.GetStrategyRequest(uuid, strategyName)
	if cpuCurrUsage == nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"status": "failed",
			"info":   "cannot set resources for non-existing strategy",
		})
		return
	}

	cpuOffset := cpuReqVal.MilliValue() - (*cpuCurrUsage).MilliValue()
	memOffset := memReqVal.Value() - (*memCurrUsage).Value()

	if cpuOffset >= 0 || memOffset >= 0 {
		if !external.RequestValidation(uuid, strconv.FormatInt(cpuOffset, 10), strconv.FormatInt(memOffset, 10)) {
			c.JSON(http.StatusBadRequest, map[string]string{
				"status": "failed",
				"info":   "insufficient resources",
			})
			return
		}
	}

	err = kubeclient.SetStrategyResourceRequest(uuid, strategyName, cpuReq, memReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"status": "failed",
			"info":   err.Error(),
		})
		return
	}

	resourceEventInfo := map[string]string{
		"strategy": strategyName,
	}
	updateInfo := map[string]interface{}{
		"cpu_request": cpuReq,
		"mem_request": memReq,
	}
	mq.PublishUpdateEvent(uuid, mq.Strategy, resourceEventInfo, updateInfo)

	c.JSON(202, map[string]string{
		"status": "success",
		"info":   fmt.Sprintf("Submitted resource request for strategy %s", strategyName),
	})
}

// Set the cpu/memory limit for strategy pod
func (src ServiceLib) SetStrategyLimit(c *gin.Context) {
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

	uuid := c.Param("username")
	strategyName := c.Param("strategyName")

	cpuReq := dataHolder["cpu"]
	memReq := dataHolder["mem"]
	if cpuReq == "" && memReq == "" {
		c.JSON(http.StatusBadRequest, map[string]string{
			"status": "failed",
			"info":   "please ensure cpu/mem body is included",
		})
		return
	}

	cpuReqVal := resource.MustParse(cpuReq)
	memReqVal := resource.MustParse(memReq)

	cpuCurrUsage, memCurrUsage := models.GetStrategyLimit(uuid, strategyName)
	if cpuCurrUsage == nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"status": "failed",
			"info":   "cannot set limit for non-existing strategy",
		})
		return
	}

	cpuOffset := cpuReqVal.MilliValue() - (*cpuCurrUsage).MilliValue()
	memOffset := memReqVal.Value() - (*memCurrUsage).Value()

	if cpuOffset >= 0 || memOffset >= 0 {
		if !external.LimitValidation(uuid, strconv.FormatInt(cpuOffset, 10), strconv.FormatInt(memOffset, 10)) {
			c.JSON(http.StatusBadRequest, map[string]string{
				"status": "failed",
				"info":   "insufficient limit resources",
			})
			return
		}
	}

	err = kubeclient.SetStrategyResourceLimit(uuid, strategyName, cpuReq, memReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{
			"status": "failed",
			"info":   err.Error(),
		})
		return
	}

	resourceEventInfo := map[string]string{
		"strategy": strategyName,
	}
	updateInfo := map[string]interface{}{
		"cpu_limit": cpuReq,
		"mem_limit": memReq,
	}
	mq.PublishUpdateEvent(uuid, mq.Strategy, resourceEventInfo, updateInfo)

	c.JSON(202, map[string]string{
		"status": "success",
		"info":   fmt.Sprintf("Submitted resource limit for strategy %s", strategyName),
	})
}
