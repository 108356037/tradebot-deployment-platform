package kubeclient

import (
	"context"
	"encoding/json"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var SecretInfo = outer{
	Auths: middle{
		RepoUrl: dockerInfo{
			Username: "108356037",
			Password: "dockerdcc123!",
			Email:    "108356037@nccu.edu.tw",
			Auth:     "MTA4MzU2MDM3OmRvY2tlcmRjYzEyMyE=",
		},
	},
}

type dockerInfo struct {
	Auth     string `json:"auth"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type middle struct {
	RepoUrl dockerInfo `json:"https://index.docker.io/v1/"`
}

type outer struct {
	Auths middle `json:"auths"`
}

func CreateRegSecret(uuid string) error {
	secretData, _ := json.Marshal(SecretInfo)

	secretClient := TypedClient.CoreV1().Secrets(uuid)
	secret := &corev1.Secret{
		Type: "kubernetes.io/dockerconfigjson",
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dockerhub-regcred",
			Namespace: uuid,
		},
		Data: map[string][]byte{".dockerconfigjson": secretData},
	}
	result, err := secretClient.Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	log.Infof("Successfully created secret %s in ns %s", result.Name, result.Namespace)
	return nil
}

func CreateCertificate() error {

	certificateRes := schema.GroupVersionResource{Group: "cert-manager.io", Version: "v1", Resource: "certificates"}
	certifcate := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "cert-manager.io/v1",
			"kind":       "Certificate",
			"metadata": map[string]interface{}{
				"name":      "ingress-cert",
				"namespace": "istio-system",
			},
			"spec": map[string]interface{}{
				"secretName":  "ingress-cert",
				"duration":    "2160h",
				"renewBefore": "360h",
				"subject": map[string]interface{}{
					"organizations": []string{
						"myhome",
					},
				},
				"isCA": false,
				"privateKey": map[string]interface{}{
					"algorithm": "RSA",
					"encoding":  "PKCS1",
					"size":      2048,
				},
				"usages": []string{
					"server auth",
					"client auth",
				},
				"dnsNames": []string{
					"*.algotrade.dev",
				},
				"issuerRef": map[string]interface{}{
					"name": "cluster-ca-issuer",
					"kind": "ClusterIssuer",
				},
			},
		},
	}

	result, err := DynamicClient.Resource(certificateRes).Namespace("istio-system").Create(context.TODO(), certifcate, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	log.Infof("Successfully created Certificate: %s\n", result.GetName())
	return nil
}
