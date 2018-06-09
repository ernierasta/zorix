// Package ping implements external ping worker.
// We are using system ping, becouse on every platform,
// there is some form of setuid on this binary, so ordinary user can run it.
package ping

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/ernierasta/zorix/shared"
)

const (
	timeout = 10 * time.Minute
)

// Ping worker
type Ping struct {
	timeout shared.Duration
}

// New return new Ping worker instance
func New(timeout shared.Duration) *Ping {
	return &Ping{timeout}
}

// Send runs ping command with reasonable.
// Returns returnCode, ping output, requestTime and error.
// For convince success returns code 200 and errors:
//   - non zero ping status: 500
func (w *Ping) Send(c shared.CheckConfig) (int, string, int64, error) {
	t0 := time.Now()
	ping := exec.Command("ping", c.Check, "-c1", fmt.Sprintf("-t%.0f", w.timeout.Seconds()))
	out, err := ping.Output()
	duration := time.Since(t0)
	if err != nil {
		return 500, "", 0, fmt.Errorf("ping.Send: ping returned non zero status, err: %v", err)
	}
	return 200, string(out), duration.Nanoseconds() / 1000 / 1000, nil

}
