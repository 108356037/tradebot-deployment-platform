package cli

import (
	"fmt"
	"os/exec"

	"bytes"

	"github.com/108356037/v1/user-resource-svc/global"
	//log "github.com/sirupsen/logrus"
)

func runCmd(cmd *exec.Cmd) *CmdResult {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()
	return &CmdResult{
		Stdout: outStr,
		Stderr: errStr,
	}
}

func CreateUserDepls(uuid string) *CmdResult {

	cmd := exec.Command("helm", "install", uuid,
		global.HelmSetting.ChartPath.AllInOne,
		"--set", "jupyter.user="+uuid,
		"--set", "grafana.podLabels.user="+uuid,
		"--set", "global.user="+uuid,
		"--namespace", uuid)

	return runCmd(cmd)
}

func GetUserDepls(uuid string) *CmdResult {

	cmd := exec.Command("helm", "list",
		"--filter", uuid,
		"--namespace", uuid, "-o", "json", "-q")

	return runCmd(cmd)
}

func DeleteUserDepls(uuid string) *CmdResult {

	cmd := exec.Command("helm", "uninstall", uuid,
		"--namespace", uuid)

	return runCmd(cmd)
}

func CreateResourceC9(uuid string) *CmdResult {
	cmd := exec.Command("helm", "install",
		fmt.Sprintf("%s-c9", uuid), // release name
		global.HelmSetting.ChartPath.C9,
		"--set", fmt.Sprintf("user=%s", uuid),
		"--namespace", uuid, "--wait")
	result := runCmd(cmd)

	return result
}

func DeleteResourceC9(uuid string) *CmdResult {
	cmd := exec.Command("helm", "uninstall",
		fmt.Sprintf("%s-c9", uuid), // release name
		"--namespace", uuid)
	return runCmd(cmd)
}

func CreateResourceJupyter(uuid string) *CmdResult {
	cmd := exec.Command("helm", "install",
		fmt.Sprintf("%s-jupyter", uuid), // release name
		global.HelmSetting.ChartPath.Jupyter,
		"--set", fmt.Sprintf("user=%s", uuid),
		"--set", fmt.Sprintf("jupyter.user=%s-jupyter", uuid),
		"--namespace", uuid, "--wait")
	result := runCmd(cmd)

	return result
}

func DeleteResourceJupyter(uuid string) *CmdResult {
	cmd := exec.Command("helm", "uninstall",
		fmt.Sprintf("%s-jupyter", uuid), // release name
		"--namespace", uuid)
	return runCmd(cmd)
}

func PatchResourceJupyter(uuid, replica string) *CmdResult {
	cmd := exec.Command("kubectl", "scale", "deployment",
		"--replicas", replica,
		fmt.Sprintf("%s-jupyter", uuid), // deployment name
		"--namespace", uuid)
	return runCmd(cmd)
}

func CreateResourceGrafana(uuid string) *CmdResult {
	cmd := exec.Command("helm", "install",
		fmt.Sprintf("%s-grafana", uuid), // release name
		global.HelmSetting.ChartPath.Grafana,
		"--set", fmt.Sprintf("user=%s", uuid),
		"--set", fmt.Sprintf("grafana.podLabels.user=%s-grafana", uuid),
		"--namespace", uuid, "--wait")
	return runCmd(cmd)
}

func DeleteResourceGrafana(uuid string) *CmdResult {
	cmd := exec.Command("helm", "uninstall",
		fmt.Sprintf("%s-grafana", uuid), // release name
		"--namespace", uuid)
	return runCmd(cmd)
}

func PatchResourceGrafana(uuid, replica string) *CmdResult {
	cmd := exec.Command("kubectl", "scale", "deployment",
		"--replicas", replica,
		fmt.Sprintf("%s-grafana", uuid), // deployment name
		"--namespace", uuid)
	return runCmd(cmd)
}
