package kubeclient

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/108356037/v1/user-resource-svc/mq"
	"github.com/lithammer/shortuuid/v3"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func PatchSvcAcc(uuid string) error {
	svcAccClient := TypedClient.CoreV1().ServiceAccounts(uuid)
	patch := []byte(`{"imagePullSecrets": [{"name": "dockerhub-regcred"}]}`)
	result, err := svcAccClient.Patch(context.Background(), "default", "application/strategic-merge-patch+json", patch, metav1.PatchOptions{})
	if err != nil {
		return err
	}

	log.Infof("Successfully patched ns %s", result.Namespace)
	return nil
}

func CreateUserNamespace(uuid string) error {
	namespaceClient := TypedClient.CoreV1().Namespaces()

	res, _ := namespaceClient.Get(context.Background(), uuid, metav1.GetOptions{})

	if res.Name == uuid {
		log.Infof("namespace %s already exists.\n", uuid)
		return errors.New("namespace " + uuid + " already exists.")
	}

	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: uuid,
			Labels: map[string]string{
				"istio-injection": "enabled",
			},
			Annotations: map[string]string{
				"openfaas": "1",
			},
		},
	}

	result, err := namespaceClient.Create(context.Background(), namespace, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	log.Infof("Successfully created namespace: %s\n", result.Name)
	return nil
}

func CreateNSWithRegcred(uuid string) error {
	err := CreateUserNamespace(uuid)
	if err != nil {
		return err
	}

	err = CreateRegSecret(uuid)
	if err != nil {
		return err
	}

	time.Sleep(2 * time.Second)

	err = PatchSvcAcc(uuid)
	if err != nil {
		return err
	}

	err = CreateRangeLimit(uuid)
	if err != nil {
		return err
	}

	err = CreateReourceQuota(uuid)
	if err != nil {
		return err
	}

	event := mq.CreateEvent{
		BasicEvent: mq.BasicEvent{
			EventId:    shortuuid.New(),
			EventType:  mq.ResourceCreate,
			OccurredAt: time.Now().Format(time.RFC3339),
		},
		UserId:         uuid,
		TargetResource: "namespace",
	}
	mqPayload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = mq.PublishMsgNoKey(os.Getenv("PUBLISH_TOPIC"), mqPayload)
	if err != nil {
		return err
	}
	log.Infof("Successfully published event %s", event.EventId)
	return nil
}

func DeleteUserNamespace(uuid string) error {
	namespaceClient := TypedClient.CoreV1().Namespaces()

	res, err := namespaceClient.Get(context.Background(), uuid, metav1.GetOptions{})

	if err != nil {
		return err
	}

	if res.Name == uuid {
		if delErr := namespaceClient.Delete(context.Background(), uuid, metav1.DeleteOptions{}); delErr != nil {
			return delErr
		}
		log.Infof("Successfully deleted namespace: %s\n", uuid)

		event := mq.DeleteEvent{
			BasicEvent: mq.BasicEvent{
				EventId:    shortuuid.New(),
				EventType:  mq.ResourceDelete,
				OccurredAt: time.Now().Format(time.RFC3339),
			},
			UserId:         uuid,
			TargetResource: "namespace",
		}
		mqPayload, err := json.Marshal(event)
		if err != nil {
			return err
		}

		err = mq.PublishMsgNoKey(os.Getenv("PUBLISH_TOPIC"), mqPayload)
		if err != nil {
			return err
		}
		log.Infof("Successfully published event %s", event.EventId)
		return nil

	} else {
		log.Infof("Namespace %s not found\n", uuid)
		return errors.New("Namespace " + uuid + "not found")
	}
}
