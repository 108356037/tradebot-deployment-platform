package sdk

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/108356037/v1/user-resource-svc/global"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
)

type globalVal struct {
	user string
}

type jupyterVal struct {
	enabled bool
	user    string
}

// Below grafana struct needs to be refactor
type grafanaVal struct {
	enabled   bool
	podLabels grafanaPodLabelInner
}

type grafanaPodLabelInner struct {
	user string
}

type HelmService struct {
	ChartPath string `json:"chartpath"`
	Namespace string `json:"namespace"`
	global    globalVal
	jupyter   jupyterVal
	grafana   grafanaVal
}

func New() HelmService {
	return HelmService{
		ChartPath: global.HelmSetting.ChartPath.AllInOne,
		Namespace: global.HelmSetting.Namespace,
	}
}

func newActionConfig(namespace string) (*action.Configuration, error) {
	settings := cli.New()
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(settings.RESTClientGetter(), namespace,
		os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, err
	}
	return actionConfig, nil
}

func (hs HelmService) CreateResourse(uuid string, enableJupyter, enableGrafana bool) error {

	hs.global.user = uuid
	hs.jupyter.user = uuid
	hs.grafana.podLabels.user = uuid

	hs.jupyter.enabled = enableJupyter
	hs.grafana.enabled = enableGrafana

	customVals := map[string]interface{}{
		"global":  hs.global,
		"jupyter": hs.jupyter,
		"grafana": hs.grafana,
	}

	actionConfig, err := newActionConfig(hs.Namespace)
	if err != nil {
		return err
	}

	chart, err := loader.Load(hs.ChartPath)
	if err != nil {
		return err
	}

	client := action.NewInstall(actionConfig)
	client.Namespace = hs.Namespace
	client.ReleaseName = uuid
	_, err = client.Run(chart, customVals)
	if err != nil {
		return err
	}
	return nil
}
