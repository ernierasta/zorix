// Package web implements web worker
package web

import (
	"bufio"
	"bytes"
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
	timeout time.Duration
}

// New return new web worker instance
func New(t shared.Duration) *Web {
	return &Web{t.Duration}
}

// Send web request to destination url
func (w *Web) Send(c shared.CheckConfig) (int, int64, error) {

	request, err := w.newRequest(&c)
	if err != nil {
		return 0, 0, err
	}

	client := http.Client{
		Timeout:       w.timeout,
		CheckRedirect: redirectGuard(c.Redirs),
	}

	t0 := time.Now()
	//resp, err := client.Get(c.Check)
	resp, err := client.Do(request)
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

/*
func PerformRequest(c shared.CheckConfig) error {

	var request *http.Request
	var reqErr error

	if len(c.FormParams) == 0 {
		//formParams create a request
		request, reqErr = http.NewRequest(c.RequestType,
			c.Url,
			nil)

	} else {
		if c.Headers[ContentType] == JsonContentType {
			//create a request using using formParams

			jsonBody, jsonErr := GetJsonParamsBody(c.FormParams)
			if jsonErr != nil {
				//Not able to create Request object.Add Error to Database

				return jsonErr
			}
			request, reqErr = http.NewRequest(c.RequestType,
				c.Url,
				jsonBody)

		} else {
			//create a request using formParams
			formParams := GetUrlValues(c.FormParams)

			request, reqErr = http.NewRequest(c.RequestType,
				c.Url,
				bytes.NewBufferString(formParams.Encode()))

			request.Header.Add(ContentLength, strconv.Itoa(len(formParams.Encode())))

			if c.Headers[ContentType] != "" {
				//Add content type to header if user doesnt mention it config file
				//Default content type application/x-www-form-urlencoded
				request.Header.Add(ContentType, FormContentType)
			}
		}
	}

	if reqErr != nil {
		//Not able to create Request object.Add Error to Database

		return reqErr
	}

	//add url parameters to query if present
	if len(c.UrlParams) != 0 {
		urlParams := GetUrlValues(c.UrlParams)
		request.URL.RawQuery = urlParams.Encode()
	}

	//Add headers to the request
	AddHeaders(request, c.Headers)

	client := &http.Client{}
	start := time.Now()

	getResponse, respErr := client.Do(request)

	if respErr != nil {
		//Request failed . Add error info to database
		var statusCode int
		if getResponse == nil {
			statusCode = 0
		} else {
			statusCode = getResponse.StatusCode
		}
		return respErr
	}

	defer getResponse.Body.Close()

	if getResponse.StatusCode != c.ResponseCode {
		//Response code is not the expected one .Add Error to database
		return errResposeCode(getResponse.StatusCode, c.ResponseCode)
	}

	elapsed := time.Since(start)

	//Request succesfull . Add infomartion to Database

	return nil
}
*/
