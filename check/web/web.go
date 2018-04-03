// Package web implements web worker
package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ernierasta/spock/shared"
)

const (
	timeout = 60 * time.Second
)

// Web worker
type Web struct {
}

// New return new web worker instance
func New() *Web {
	return &Web{}
}

// Send web request to destination url
func (w *Web) Send(c shared.Check) (int, int64, error) {
	t0 := time.Now()
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(c.Check) //needs port too? what aboult url?
	t1 := time.Now()
	if err != nil {
		return 0, 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("web.Test: can not read response")
	}
	if len(body) == 0 {
		return 0, 0, fmt.Errorf("web.Test: returned body is empty")
	}

	defer resp.Body.Close()
	return resp.StatusCode, t1.Sub(t0).Nanoseconds() / 1000 / 1000, nil

}
