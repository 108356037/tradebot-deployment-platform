package kubeclient

import (
	"context"
	"fmt"

	"github.com/108356037/v1/strategy-manager/internal/models"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// old func, deprecated for now
// given a user strategy, set the cpu/ram  for it
// user: string, strategy: string, cpu: string(m), ram: string(Mi, Gi)
func SetResourcesForStrategy(user, strategy, cpu, mem string) error {

	deplClient := TypedClient.AppsV1().Deployments(user)
	cpuVal := resource.MustParse(cpu)
	memVal := resource.MustParse(mem)
	patchString := ""

	switch {
	case cpuVal.MilliValue() > 750 && memVal.Value() > 512*1024*1024:
		patchString = fmt.Sprintf("{\"spec\":{\"strategy\":{\"rollingUpdate\":{\"maxUnavailable\":1}},\"template\":{\"spec\":{\"containers\":[{\"name\":\"%s\",\"resources\":{\"limits\":{\"cpu\":\"%s\",\"memory\":\"%s\"},\"requests\":{\"cpu\":\"%s\",\"memory\":\"%s\"}}}]}}}}", strategy, cpu, mem, cpu, mem)
	case cpuVal.MilliValue() > 750 && memVal.Value() < 512*1024*1024:
		patchString = fmt.Sprintf("{\"spec\":{\"strategy\":{\"rollingUpdate\":{\"maxUnavailable\":1}},\"template\":{\"spec\":{\"containers\":[{\"name\":\"%s\",\"resources\":{\"limits\":{\"cpu\":\"%s\",\"memory\":\"512Mi\"},\"requests\":{\"cpu\":\"%s\",\"memory\":\"%s\"}}}]}}}}", strategy, cpu, cpu, mem)
	case cpuVal.MilliValue() < 750 && memVal.Value() > 512*1024*1024:
		patchString = fmt.Sprintf("{\"spec\":{\"strategy\":{\"rollingUpdate\":{\"maxUnavailable\":1}},\"template\":{\"spec\":{\"containers\":[{\"name\":\"%s\",\"resources\":{\"limits\":{\"cpu\":\"750m\",\"memory\":\"%s\"},\"requests\":{\"cpu\":\"%s\",\"memory\":\"%s\"}}}]}}}}", strategy, mem, cpu, mem)
	default:
		patchString = fmt.Sprintf("{\"spec\":{\"strategy\":{\"rollingUpdate\":{\"maxUnavailable\":1}},\"template\":{\"spec\":{\"containers\":[{\"name\":\"%s\",\"resources\":{\"limits\":{\"cpu\":\"750m\",\"memory\":\"512Mi\"},\"requests\":{\"cpu\":\"%s\",\"memory\":\"%s\"}}}]}}}}", strategy, cpu, mem)
	}

	//patchString := fmt.Sprintf("{\"spec\":{\"template\":{\"spec\":{\"containers\":[{\"name\":\"%s\",\"resources\":{\"requests\":{\"cpu\":\"%s\",\"memory\":\"%s\"}}}]}}}}", strategy, cpu, ram)
	//patchString := fmt.Sprintf("{\"spec\":{\"template\":{\"spec\":{\"containers\":[{\"name\":\"%s\",\"resources\":{\"limits\":{\"cpu\":\"1000m\",\"memory\":\"1Gi\"},\"requests\":{\"cpu\":\"%s\",\"memory\":\"%s\"}}}]}}}}", strategy, cpu, ram)
	//patchString := fmt.Sprintf("{\"spec\":{\"strategy\":{\"rollingUpdate\":{\"maxUnavailable\":1}},\"template\":{\"spec\":{\"containers\":[{\"name\":\"%s\",\"resources\":{\"limits\":{\"cpu\":\"1000m\",\"memory\":\"1Gi\"},\"requests\":{\"cpu\":\"%s\",\"memory\":\"%s\"}}}]}}}}", strategy, cpu, ram)

	_, err := deplClient.Patch(context.Background(), strategy, "application/strategic-merge-patch+json", []byte(patchString), v1.PatchOptions{})
	if err != nil {
		return err
	}
	return nil
}

/// given a user strategy, set the cpu/mem request for it
// user: string, strategy: string, cpu: string(m), ram: string(Mi, Gi)
func SetStrategyResourceRequest(user, strategy, cpu, mem string) error {
	cpuLim, memLim := models.GetStrategyLimit(user, strategy)
	deplClient := TypedClient.AppsV1().Deployments(user)
	patchString := fmt.Sprintf("{\"spec\":{\"strategy\":{\"rollingUpdate\":{\"maxUnavailable\":1}},\"template\":{\"spec\":{\"containers\":[{\"name\":\"%s\",\"resources\":{\"requests\":{\"cpu\":\"%s\",\"memory\":\"%s\"},\"limits\":{\"cpu\":\"%s\",\"memory\":\"%s\"}}}]}}}}", strategy, cpu, mem, resourceQuotaParser(*cpuLim), resourceQuotaParser(*memLim))
	_, err := deplClient.Patch(context.Background(), strategy, "application/strategic-merge-patch+json", []byte(patchString), v1.PatchOptions{})
	if err != nil {
		return err
	}
	return nil
}

/// given a user strategy, set the cpu/mem limit for it
// user: string, strategy: string, cpu: string(m), ram: string(Mi, Gi)
func SetStrategyResourceLimit(user, strategy, cpu, mem string) error {
	cpuReq, memReq := models.GetStrategyRequest(user, strategy)
	deplClient := TypedClient.AppsV1().Deployments(user)
	patchString := fmt.Sprintf("{\"spec\":{\"strategy\":{\"rollingUpdate\":{\"maxUnavailable\":1}},\"template\":{\"spec\":{\"containers\":[{\"name\":\"%s\",\"resources\":{\"requests\":{\"cpu\":\"%s\",\"memory\":\"%s\"},\"limits\":{\"cpu\":\"%s\",\"memory\":\"%s\"}}}]}}}}", strategy, resourceQuotaParser(*cpuReq), resourceQuotaParser(*memReq), cpu, mem)
	_, err := deplClient.Patch(context.Background(), strategy, "application/strategic-merge-patch+json", []byte(patchString), v1.PatchOptions{})
	if err != nil {
		return err
	}
	return nil
}

// given a user strategy, set the crontab schedule
func ScheduleStrategy(user, strategy, crontabSchedule string) error {
	deplClient := TypedClient.AppsV1().Deployments(user)
	patchString := fmt.Sprintf("{\"spec\":{\"strategy\":{\"rollingUpdate\":{\"maxUnavailable\":1}},\"template\":{\"metadata\":{\"annotations\":{\"prometheus.io.scrape\":\"false\",\"schedule\":\"%s\",\"topic\":\"cron-function\"}}}}}", crontabSchedule)

	_, err := deplClient.Patch(context.Background(), strategy, "application/strategic-merge-patch+json", []byte(patchString), v1.PatchOptions{})
	if err != nil {
		return err
	}
	return nil
}

func resourceQuotaParser(q resource.Quantity) string {
	val := []byte{}
	val, suffix := q.CanonicalizeBytes(val)

	return string(val[:]) + string(suffix[:])
}
