package server

import (
	"errors"
	"os/exec"
)

func runCmd(cmd *exec.Cmd) error {
	output, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(output))
	}
	return nil
}
