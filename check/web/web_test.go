package web

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/ernierasta/zorix/shared"
)

var (
	commonHeaders  = "Authorization:Bearer abcbc123123abc\nContent-Type:application/json"
	commonQHeaders = fmt.Sprintf("%q:%q\n%q:%q", "Authorization", "Bearer abcbc123123abc", "Content-Type", "application/json")

	commonHParsed = map[string][]string{
		"Authorization": []string{"Bearer abcbc123123abc"},
		"Content-Type":  []string{"application/json"},
	}
	commonLowerHParsed = map[string][]string{
		"authorization": []string{"Bearer abcbc123123abc"},
		"content-type":  []string{"application/json"},
	}
	commonHPWithoutType = map[string][]string{
		"Authorization": []string{"Bearer abcbc123123abc"},
	}
)

func TestNew(t *testing.T) {
	type args struct {
		t shared.Duration
	}
	tests := []struct {
		name string
		args args
		want *Web
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeb_Send(t *testing.T) {
	type fields struct {
		timeout time.Duration
	}
	type args struct {
		c shared.CheckConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Web{
				timeout: tt.fields.timeout,
			}
			got, got1, err := w.Send(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Web.Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Web.Send() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Web.Send() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_redirectGuard(t *testing.T) {
	type args struct {
		r int
	}
	tests := []struct {
		name string
		args args
		want func(req *http.Request, via []*http.Request) error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := redirectGuard(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				//				t.Errorf("redirectGuard() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseHeaders(t *testing.T) {

	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{name: "empty headers", args: args{s: ""}, want: map[string][]string{}},
		{name: "common headers", args: args{s: commonHeaders}, want: commonHParsed},
		{name: "common quoted", args: args{s: commonQHeaders}, want: commonHParsed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseHeaders(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseHeaders() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addHeaders(t *testing.T) {

	r, err := http.NewRequest("GET", "http://something.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	r2, err := http.NewRequest("GET", "http://something.com", nil)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		r *http.Request
		h map[string][]string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{name: "add common headers", args: args{r: r, h: commonHParsed}, want: commonHParsed},
		{name: "add common lower case headers", args: args{r: r2, h: commonLowerHParsed}, want: commonHParsed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addHeaders(tt.args.r, tt.args.h)
			if !reflect.DeepEqual(map[string][]string(tt.args.r.Header), tt.want) {
				t.Errorf("%+v, want: %+v", tt.args.r.Header, tt.want)
			}
		})
	}
}

func Test_setContentType(t *testing.T) {
	type args struct {
		headers     map[string][]string
		contTypeVal string
	}
	tests := []struct {
		name string
		args args
		want map[string][]string
	}{
		{"add content type", args{headers: commonHPWithoutType, contTypeVal: jsonContentType}, commonHParsed},
		{"content type already set", args{headers: commonHParsed, contTypeVal: jsonContentType}, commonHParsed},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setContentType(tt.args.headers, tt.args.contTypeVal)
			if !reflect.DeepEqual(tt.args.headers, tt.want) {
				t.Errorf("%+v, want: %+v", tt.args.headers, tt.want)
			}
		})
	}
}

func TestWeb_newRequest(t *testing.T) {

	jsonCheck := &shared.CheckConfig{
		Params: `{"data": [{"password": "xxx", "email": "xxx"}]}`,
		Method: "POST",
	}

	req1 := &http.Request{
		Method: "POST",
		Header: map[string][]string{contentType: []string{jsonContentType}},
		Body:   io.ReadCloser(bytes.NewBuffer([]byte(""))),
	}

	type fields struct {
		timeout time.Duration
	}
	type args struct {
		c *shared.CheckConfig
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Request
		wantErr bool
	}{

		{name: "web request", fields: fields{timeout: 60 * time.Second}, args: args{c: jsonCheck}, want: req1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Web{
				timeout: tt.fields.timeout,
			}
			got, err := w.newRequest(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Web.newRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Proto, tt.want.Proto) {
				t.Errorf("Web.newRequest() = %v\nwant:%+v\n", got.Proto, tt.want.Proto)
			}
			if !reflect.DeepEqual(got.Header, tt.want.Header) {
				t.Errorf("Web.newRequest() = %v\nwant:%+v\n", got.Header, tt.want.Header)
			}
			if !reflect.DeepEqual(got.Body, tt.want.Body) {
				t.Errorf("Web.newRequest() = %v\nwant:%+v\n", got.Body, tt.want.Body)
			}
			if !reflect.DeepEqual(got.Proto, tt.want.Proto) {
				t.Errorf("Web.newRequest() = %v\nwant:%+v\n", got.Proto, tt.want.Proto)
			}
		})
	}
}
