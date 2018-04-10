// Package port implements portscanner worker.
package port

import (
	"fmt"
	"net"
	"time"

	"github.com/ernierasta/zorix/shared"
)

// Port worker
type Port struct {
	timeout shared.Duration
}

// New return new Port worker instance
func New(timeout shared.Duration) *Port {
	return &Port{timeout}
}

// Send runs portscan.
// Returns returnCode, empty body, requestTime and error.
// For convince success returns code 200 and errors:
//   - closed: 500
func (p *Port) Send(c shared.CheckConfig) (int, string, int64, error) {
	t0 := time.Now()
	fmt.Println("tcp", c.Check, p.timeout.String())
	conn, err := net.DialTimeout("tcp", c.Check, p.timeout.Duration)
	duration := time.Since(t0)
	if err != nil {
		return 500, "", 0, fmt.Errorf("port.Send: port returned non zero status, err: %v", err)
	}
	conn.Close()
	return 200, "", duration.Nanoseconds() / 1000 / 1000, nil

}
