package notify

import (
	"fmt"
	"io"
	"strings"

	"github.com/ernierasta/zorix/shared"
)

func createParser(c shared.Check) func(w io.Writer, tag string) (int, error) {
	return func(w io.Writer, tag string) (int, error) {
		switch tag {
		case "check":
			return w.Write([]byte(c.Check))
		case "params":
			if len(c.Params) > 0 {
				return w.Write([]byte(" " + strings.Trim(fmt.Sprint(c.Params), "[]")))
			}
			return w.Write([]byte(""))
		case "timestamp":
			return w.Write([]byte(c.Timestamp.Format("2.1.2006 15:04:05")))
		case "responsecode":
			return w.Write([]byte(fmt.Sprintf("%d", c.ReturnedCode)))
		case "responsetime":
			return w.Write([]byte(fmt.Sprintf("%d", c.ReturnedTime)))
		case "expectedcode":
			return w.Write([]byte(fmt.Sprintf("%d", c.ExpectedCode)))
		case "expectedtime":
			return w.Write([]byte(fmt.Sprintf("%d", c.ExpectedTime)))
		case "error":
			if c.Error != nil {
				return w.Write([]byte(c.Error.Error()))
			}
			return w.Write([]byte(""))
			//TODO: add all fields from shared.Check
		default:
			return w.Write([]byte(fmt.Sprintf("[unknown tag '%s']", tag)))
		}
	}
}
