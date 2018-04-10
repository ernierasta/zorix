// Package web implements web worker
package web

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/textproto"
	"strings"
	"time"

	"github.com/ernierasta/zorix/shared"
)

const (
	formContentType = "application/x-www-form-urlencoded"
	jsonContentType = "application/json"
	contentType     = "Content-Type"
)

// Web worker
type Web struct {
	timeout    time.Duration
	ignoreCert bool
}

// New return new web worker instance
func New(t shared.Duration, ignoreCert bool) *Web {
	return &Web{t.Duration, ignoreCert}
}

// Send web request to destination url
func (w *Web) Send(c shared.CheckConfig) (int, string, int64, error) {

	request, err := w.newRequest(&c)
	if err != nil {
		return 0, "", 0, err
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: w.ignoreCert,
		},
	}

	client := http.Client{
		Transport:     transport,
		Timeout:       w.timeout,
		CheckRedirect: redirectGuard(c.Redirs),
	}

	t0 := time.Now()
	//resp, err := client.Get(c.Check)
	resp, err := client.Do(request)
	reqDur := time.Since(t0)
	if err != nil {
		return 0, "", 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, "", 0, fmt.Errorf("web.Test: can not read response")
	}
	if len(body) == 0 {
		return 0, "", 0, fmt.Errorf("web.Test: returned body is empty")
	}
	defer resp.Body.Close()

	return resp.StatusCode, string(body), reqDur.Nanoseconds() / 1000 / 1000, nil

}

// newRequest prepares http.Request with needed headers and body
func (w *Web) newRequest(c *shared.CheckConfig) (*http.Request, error) {
	request := &http.Request{}
	headers := parseHeaders(c.Headers)

	// we have Form params
	if len(c.Params) != 0 {
		// determine ContentType by ':' (json) or '=' (url encoded)
		switch {
		case strings.Contains(c.Params, ":"):
			// we have json
			setContentType(headers, jsonContentType) // maybe add header
		case strings.Contains(c.Params, "="):
			// we have urlencoded params
			setContentType(headers, formContentType) // maybe add header
		default:
			return request, fmt.Errorf("unknown parameters for %s, params: %s, check config file", c.Check, c.Params)
		}

	}

	rb := bytes.NewBuffer([]byte(c.Params))
	request, err := http.NewRequest(c.Method, c.Check, rb)
	if err != nil {
		return request, err
	}

	addHeaders(request, headers)

	return request, nil
}

func redirectGuard(r int) func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if len(via) > r {
			return fmt.Errorf("to many redirects, only %d allowed", r)
		}
		return nil
	}
}

// setContentType adds "Content-Type" value to headers
func setContentType(headers map[string][]string, contTypeVal string) {
	if _, ok := headers[contentType]; !ok {
		headers[contentType] = []string{contTypeVal}
	}
}

// parseHeaders will take string and parse it to map, which
// can be used in request struct
func parseHeaders(s string) map[string][]string {
	s = strings.Replace(s, "\n", "\r\n", -1)
	s = strings.Replace(s, `"`, "", -1)
	r := bufio.NewReader(strings.NewReader(s + "\r\n\r\n"))
	tp := textproto.NewReader(r)
	if h, err := tp.ReadMIMEHeader(); err == nil {
		return h
	}

	return map[string][]string{}
}

// addHeaders adds headers to the request
func addHeaders(r *http.Request, h map[string][]string) {
	//r.Header = map[string][]string{} // NewRequest initializes this map, so not needed
	if len(h) > 0 {
		for k, vv := range h {
			for _, v := range vv { // go through []string
				r.Header.Add(k, v)
			}
		}
	}
}
