package cmd

import (
	"fmt"
	"os/exec"

	"github.com/ernierasta/zorix/shared"
	log "github.com/sirupsen/logrus"
)

// Cmd notifier
type Cmd struct {
}

// Send sends message running some command.
func (cd *Cmd) Send(c shared.CheckConfig, n shared.NotifConfig) error {
	log.Debugf("cmd.Send: command: sh -c %s", n.Cmd)
	cmd := exec.Command("sh", "-c", n.Cmd)
	out, err := cmd.Output()
	log.Debugf("cmd.Send: output: %s", out)
	if err != nil {
		return fmt.Errorf("cmd.Send: process 'sh -c %q',  returned non zero status, err: %v", n.Cmd, err)
	}
	return nil
}
