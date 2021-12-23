package kubeclient

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var TypedClient *kubernetes.Clientset
var DynamicClient dynamic.Interface

func Init() {

	config, err := rest.InClusterConfig() // service account for in cluster pod
	if err != nil {
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		if envvar := os.Getenv("KUBECONFIG"); len(envvar) > 0 {
			kubeconfig = envvar
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatal(err)
		}
	}

	// set up typed client
	typedClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	TypedClient = typedClientset

	// set up dynamic client (for crds)
	dynamicClientset, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	DynamicClient = dynamicClientset
}
