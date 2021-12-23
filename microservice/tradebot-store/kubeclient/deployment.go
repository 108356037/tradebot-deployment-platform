package kubeclient

import (
	"context"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// var (
// 	ImageLookupMap map[string][]string
// )

func int32Ptr(i int32) *int32 { return &i }
func int64Ptr(i int64) *int64 { return &i }

func CreateBot(uuid, botname string, envs map[string]string) error {

	publisher := envs["PUBLISHER"]
	envsHolder := []corev1.EnvVar{}
	if len(envs) > 0 {
		for k, v := range envs {
			if k != "PUBLISHER" {
				envsHolder = append(envsHolder, corev1.EnvVar{
					Name:  k,
					Value: v,
				})
			}
		}
	}

	deplClient := TypedClient.AppsV1().Deployments(uuid)
	botDepl := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: botname,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": botname,
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				RollingUpdate: &appsv1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{IntVal: 1},
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":           botname,
						"faas_function": "exist-to-avoid-sidecar-injection",
						"identity":      "tradebot",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  botname,
							Image: os.Getenv("STORE_IMAGE_REPO") + "/" + publisher + "/" + botname + "-bot",
							SecurityContext: &corev1.SecurityContext{
								RunAsUser:  int64Ptr(101),
								RunAsGroup: int64Ptr(101),
							},
							Env: envsHolder,
						},
					},
				},
			},
		},
	}
	_, err := deplClient.Create(context.Background(), botDepl, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return nil
}

func DeleteBot(uuid, botname string) error {
	deplClient := TypedClient.AppsV1().Deployments(uuid)
	err := deplClient.Delete(context.Background(), botname, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func GetAllBots(uuid string) error {
	deplClient := TypedClient.AppsV1().Deployments(uuid)
	res, err := deplClient.List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	log.Info(res.Items)
	return nil
}

func GetSingleBot(uuid, botname string) error {
	deplClient := TypedClient.AppsV1().Deployments(uuid)
	res, err := deplClient.Get(context.Background(), botname, metav1.GetOptions{})
	if err != nil {
		return err
	}
	log.Info(res.Status)
	return nil
}

func UpdateBot(uuid, botname string) error {
	data := fmt.Sprintf(`{"spec":{"template":{"metadata":{"annotations":{"kubectl.kubernetes.io/restartedAt":"%s"}}}},"strategy":{"type":"RollingUpdate","rollingUpdate":{"maxUnavailable":1,"maxSurge": "%s"}}}`, time.Now().String(), "25%")
	deplClient := TypedClient.AppsV1().Deployments(uuid)
	newDepl, err := deplClient.Patch(context.Background(), botname, types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-rollout"})
	if err != nil {
		return err
	}
	log.Info(newDepl)
	return nil
	//newDeployment, err := clientImpl.ClientSet.AppsV1().Deployments(item.Pod.Namespace).Patch(context.Background(), deployment.Name, types.StrategicMergePatchType, []byte(data), metav1.PatchOptions{FieldManager: "kubectl-rollout"})
}
