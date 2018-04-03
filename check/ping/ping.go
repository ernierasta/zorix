// Package web implements ping worker
package ping

import (
	"time"

	"github.com/ernierasta/zorix/shared"
)

const (
	timeout = 60 * time.Second
)

// Ping worker
type Ping struct {
}

// New return new ping worker instance
func New() *Ping {
	return &Ping{}
}

// Send ping
func (w *Ping) Send(c shared.Check) (int, int64, error) {
	t0 := time.Now()

	t1 := time.Now()
	return resp.StatusCode, t1.Sub(t0).Nanoseconds() / 1000 / 1000, nil

}
