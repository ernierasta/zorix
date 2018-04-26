// Package cmd implements external command worker.
// It will run any command you give it, and fail will
// depend on returned status code.
package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/ernierasta/zorix/shared"
	log "github.com/sirupsen/logrus"
)

const (
	timeout = 10 * time.Minute
)

// Cmd worker
type Cmd struct {
}

// New return new Cmd worker instance
func New() *Cmd {
	return &Cmd{}
}

// Send runs selected command with given params.
// Returns returnCode, requestTime and error.
// For convience success returns code 200 and errors:
//   - starting command: 404
//   - non zero status: 500
func (w *Cmd) Send(c shared.CheckConfig) (int, string, int64, error) {
	t0 := time.Now()
	cl := c.Check + " " + c.Params
	o, err := exec.Command("sh", "-c", cl).Output()
	duration := time.Since(t0)
	outs := strings.TrimSpace(string(o))
	log.Debugf("cmd.Send: output: %s", o)
	if err != nil {
		return 500, outs, duration.Nanoseconds() / 1000 / 1000, fmt.Errorf("cmd.Send: process 'sh -c %q',  returned non zero status, err: %v", cl, err)
	}

	return 200, outs, duration.Nanoseconds() / 1000 / 1000, nil

}
