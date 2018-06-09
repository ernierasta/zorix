//this program will generate config boilerplate functions based on config.go
package main

import (
	"fmt"
	"reflect"

	"github.com/ernierasta/zorix/internal/platform/config"
)

// plan:
// - get all check types
// - get tags for every type (Web, Ping, ...)
// - generate validation and defaults.

func main() {
	c := config.Check{}

	t := reflect.TypeOf(c)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Println(f.Name, f.Tag)
	}
}
