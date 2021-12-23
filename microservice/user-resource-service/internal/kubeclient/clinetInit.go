package kubeclient

import (
	"os"
	"path/filepath"

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

// func RunInitWithExec() error {
// 	cmd := exec.Command("kubectl", "apply", "-f", global.HelmSetting.InitYamlPath)
// 	stdoutReader, _ := cmd.StdoutPipe()
// 	stdoutScanner := bufio.NewScanner(stdoutReader)
// 	go func() {
// 		for stdoutScanner.Scan() {
// 			log.Info(stdoutScanner.Text())
// 		}
// 	}()
// 	stderrReader, _ := cmd.StderrPipe()
// 	stderrScanner := bufio.NewScanner(stderrReader)
// 	go func() {
// 		for stderrScanner.Scan() {
// 			log.Warning(stderrScanner.Text())
// 		}
// 	}()
// 	err := cmd.Start()
// 	if err != nil {
// 		log.Error("Error : %v \n", err)
// 		return err
// 	}
// 	err = cmd.Wait()
// 	if err != nil {
// 		log.Error("Error: %v \n", err)
// 		return err
// 	}

// 	return nil
// }
