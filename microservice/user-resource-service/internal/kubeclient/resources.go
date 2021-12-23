package kubeclient

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ResourceHard struct {
	LimitCpu   string `json:"limits.cpu"`
	LimitMem   string `json:"limits.memory"`
	RequestCpu string `json:"requests.cpu"`
	RequestMem string `json:"requests.memory"`
}

type ResourceUsed struct {
	LimitCpu   string `json:"limits.cpu"`
	LimitMem   string `json:"limits.memory"`
	RequestCpu string `json:"requests.cpu"`
	RequestMem string `json:"requests.memory"`
}

type ResourceStatus struct {
	Hard ResourceHard `json:"hard"`
	Used ResourceUsed `json:"used"`
}

func CreateReourceQuota(namespace string) error {

	resourceQuotaClient := TypedClient.CoreV1().ResourceQuotas(namespace)

	resourceQuotaConfig := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name: "default-resource-quota",
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: corev1.ResourceList{
				corev1.ResourceRequestsCPU:    *resource.NewMilliQuantity(2600, resource.DecimalSI),     // 2000+(100*3*2) = 2600
				corev1.ResourceRequestsMemory: *resource.NewQuantity(2816*1024*1024, resource.BinarySI), // 2048+768(128*3*2) = 2816
				corev1.ResourceLimitsCPU:      *resource.NewMilliQuantity(6500, resource.DecimalSI),     // 2000+(750*3*2) = 6500
				corev1.ResourceLimitsMemory:   *resource.NewQuantity(5120*1024*1024, resource.BinarySI), // 2048+(512*3*2) = 5120
			},
		},
	}

	_, err := resourceQuotaClient.Create(context.Background(), resourceQuotaConfig, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	return nil
}

func CreateRangeLimit(namespace string) error {
	limitRangeClient := TypedClient.CoreV1().LimitRanges(namespace)

	limitRangeConfig := &corev1.LimitRange{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pod-limit-range",
		},
		Spec: corev1.LimitRangeSpec{
			Limits: []corev1.LimitRangeItem{
				{
					Type: "Container",
					Max: corev1.ResourceList{
						corev1.ResourceCPU:    *resource.NewMilliQuantity(2000, resource.DecimalSI),
						corev1.ResourceMemory: *resource.NewQuantity(2*1024*1024*1024, resource.BinarySI),
					},
					Min: corev1.ResourceList{
						corev1.ResourceCPU:    *resource.NewMilliQuantity(100, resource.DecimalSI),
						corev1.ResourceMemory: *resource.NewQuantity(128*1024*1024, resource.BinarySI),
					},
					Default: corev1.ResourceList{
						corev1.ResourceCPU:    *resource.NewMilliQuantity(750, resource.DecimalSI),
						corev1.ResourceMemory: *resource.NewQuantity(512*1024*1024, resource.BinarySI),
					},
					DefaultRequest: corev1.ResourceList{
						corev1.ResourceCPU:    *resource.NewMilliQuantity(100, resource.DecimalSI),
						corev1.ResourceMemory: *resource.NewQuantity(128*1024*1024, resource.BinarySI),
					},
				},
			},
		},
	}

	_, err := limitRangeClient.Create(context.Background(), limitRangeConfig, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

// get resourceQuota status in human readable format
func GetNsResourceQuota(namespace string) (*ResourceStatus, error) {

	resourceQuotaClient := TypedClient.CoreV1().ResourceQuotas(namespace)

	res, err := resourceQuotaClient.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	if len(res.Items) == 0 {
		return nil, errors.New("no resource qouta found in namespace")
	}

	hardReqCpu := res.Items[0].Status.Hard["requests.cpu"]
	hardReqMem := res.Items[0].Status.Hard["requests.memory"]
	hardLimCpu := res.Items[0].Status.Hard["limits.cpu"]
	hardLimMem := res.Items[0].Status.Hard["limits.memory"]

	usedReqCpu := res.Items[0].Status.Used["requests.cpu"]
	usedReqMem := res.Items[0].Status.Used["requests.memory"]
	usedLimCpu := res.Items[0].Status.Used["limits.cpu"]
	usedLimMem := res.Items[0].Status.Used["limits.memory"]

	// fmt.Println(res.Items[0].Status)
	resourceQuotaStaus := ResourceStatus{
		Hard: ResourceHard{
			LimitCpu:   resourceQuotaParser(hardLimCpu),
			LimitMem:   resourceQuotaParser(hardLimMem),
			RequestCpu: resourceQuotaParser(hardReqCpu),
			RequestMem: resourceQuotaParser(hardReqMem),
		},
		Used: ResourceUsed{
			LimitCpu:   resourceQuotaParser(usedLimCpu),
			LimitMem:   resourceQuotaParser(usedLimMem),
			RequestCpu: resourceQuotaParser(usedReqCpu),
			RequestMem: resourceQuotaParser(usedReqMem),
		},
	}
	fmt.Println(resourceQuotaStaus)

	return &resourceQuotaStaus, nil
}

// get resourceQuota status in string(Int64),
// json has limit size in int, so frontend will decode str to int
func GetNsResourceQuotaInt64(namespace string) (*ResourceStatus, error) {
	resourceQuotaClient := TypedClient.CoreV1().ResourceQuotas(namespace)

	res, err := resourceQuotaClient.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	if len(res.Items) == 0 {
		return nil, errors.New("no resource qouta found in namespace")
	}

	hardReqCpu := res.Items[0].Status.Hard["requests.cpu"]
	hardReqMem := res.Items[0].Status.Hard["requests.memory"]
	hardLimCpu := res.Items[0].Status.Hard["limits.cpu"]
	hardLimMem := res.Items[0].Status.Hard["limits.memory"]

	usedReqCpu := res.Items[0].Status.Used["requests.cpu"]
	usedReqMem := res.Items[0].Status.Used["requests.memory"]
	usedLimCpu := res.Items[0].Status.Used["limits.cpu"]
	usedLimMem := res.Items[0].Status.Used["limits.memory"]

	resourceQuotaStaus := ResourceStatus{
		Hard: ResourceHard{
			LimitCpu:   strconv.FormatInt(hardLimCpu.MilliValue(), 10),
			LimitMem:   strconv.FormatInt(hardLimMem.Value(), 10),
			RequestCpu: strconv.FormatInt(hardReqCpu.MilliValue(), 10),
			RequestMem: strconv.FormatInt(hardReqMem.Value(), 10),
		},
		Used: ResourceUsed{
			LimitCpu:   strconv.FormatInt(usedLimCpu.MilliValue(), 10),
			LimitMem:   strconv.FormatInt(usedLimMem.Value(), 10),
			RequestCpu: strconv.FormatInt(usedReqCpu.MilliValue(), 10),
			RequestMem: strconv.FormatInt(usedReqMem.Value(), 10),
		},
	}

	fmt.Println(resourceQuotaStaus)

	return &resourceQuotaStaus, nil
}

// given value, suffix as (byte[], byte[]), return the human readable value (ex: "1", "1Gi", "500Mi")
func resourceQuotaParser(q resource.Quantity) string {
	val := []byte{}
	val, suffix := q.CanonicalizeBytes(val)

	return string(val[:]) + string(suffix[:])
}

// returns the remaing request resources avab for given ns, returns (cpu(int64), memory(int64), error)
func remainRequest(namespace string) (*int64, *int64, error) {
	resourceQuotaClient := TypedClient.CoreV1().ResourceQuotas(namespace)

	res, err := resourceQuotaClient.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}
	if len(res.Items) == 0 {
		return nil, nil, errors.New("no resource qouta found in namespace")
	}

	hardCpu := res.Items[0].Status.Hard["requests.cpu"]
	hardMem := res.Items[0].Status.Hard["requests.memory"]
	usedCpu := res.Items[0].Status.Used["requests.cpu"]
	usedMem := res.Items[0].Status.Used["requests.memory"]

	remainCpu := hardCpu.MilliValue() - usedCpu.MilliValue()
	remainMem := hardMem.Value() - usedMem.Value()
	return &remainCpu, &remainMem, nil
}

// returns the remaing limit resources avab for given ns, returns (cpu(int64), memory(int64), error)
func remainLimit(namespace string) (*int64, *int64, error) {
	resourceQuotaClient := TypedClient.CoreV1().ResourceQuotas(namespace)

	res, err := resourceQuotaClient.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, nil, err
	}
	if len(res.Items) == 0 {
		return nil, nil, errors.New("no resource qouta found in namespace")
	}

	hardCpu := res.Items[0].Status.Hard["limits.cpu"]
	hardMem := res.Items[0].Status.Hard["limits.memory"]
	usedCpu := res.Items[0].Status.Used["limits.cpu"]
	usedMem := res.Items[0].Status.Used["limits.memory"]

	remainCpu := hardCpu.MilliValue() - usedCpu.MilliValue()
	remainMem := hardMem.Value() - usedMem.Value()
	return &remainCpu, &remainMem, nil
}

// check if remain resource.requests is enough
func ResourceRequestValidator(cpu, memory, namespace string) bool {
	remainCpu, remainMem, err := remainRequest(namespace)
	if err != nil {
		return false
	}

	// reqCpu := resource.MustParse(cpu)
	// reqMem := resource.MustParse(memory)
	// if *remainCpu >= reqCpu.MilliValue() && *remainMem >= reqMem.Value() {
	// 	return true
	// }

	reqCpu, err := strconv.ParseInt(cpu, 10, 64)
	if err != nil {
		return false
	}

	reqMem, err := strconv.ParseInt(memory, 10, 64)
	if err != nil {
		return false
	}

	if *remainCpu >= reqCpu && *remainMem >= reqMem {
		return true
	}

	return false
}

// check if remain resource.limits is enough
func ResourceLimitValidator(cpu, memory, namespace string) bool {
	remainCpu, remainMem, err := remainLimit(namespace)
	if err != nil {
		return false
	}

	reqCpu, err := strconv.ParseInt(cpu, 10, 64)
	if err != nil {
		return false
	}

	reqMem, err := strconv.ParseInt(memory, 10, 64)
	if err != nil {
		return false
	}

	if *remainCpu >= reqCpu && *remainMem >= reqMem {
		return true
	}

	return false
}
