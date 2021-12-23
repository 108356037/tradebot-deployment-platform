package kubeclient

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var TypedClient *kubernetes.Clientset
var DynamicClient dynamic.Interface

func Init() error {

	config, err := rest.InClusterConfig() // service account for in cluster pod
	if err != nil {
		log.Info("server inside local env")
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		if envvar := os.Getenv("KUBECONFIG"); len(envvar) > 0 {
			kubeconfig = envvar
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			//log.Warn("Kubeconfig cannot be loaded, err: %s", err.Error())
			return err
		}
	}

	// set up typed client
	typedClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	TypedClient = typedClientset

	// set up dynamic client (for crds)
	dynamicClientset, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}
	DynamicClient = dynamicClientset

	return nil
}
