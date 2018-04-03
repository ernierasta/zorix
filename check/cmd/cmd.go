// Package cmd implements external command worker.
// It will run any command you give it, and fail will
// depend on returned status code.
package cmd

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/ernierasta/spock/shared"
)

const (
	timeout = 60 * time.Second
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
func (w *Cmd) Send(c shared.Check) (int, int64, error) {
	t0 := time.Now()
	cmd := exec.Command(c.Check, c.Params...)
	err := cmd.Start()
	if err != nil {
		return 404, 0, fmt.Errorf("cmd.Send: problem starting command '%s', params: '%s', err: %v", c.Check, c.Params, err)
	}
	err = cmd.Wait()
	t1 := time.Now()
	if err != nil {
		return 500, 0, fmt.Errorf("cmd.Send: process '%s', params: '%s',  returned non zero status, err: %v", c.Check, c.Params, err)
	}

	return 200, t1.Sub(t0).Nanoseconds() / 1000 / 1000, nil

}