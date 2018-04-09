package template

import (
	"fmt"
	"io"

	"github.com/ernierasta/zorix/shared"
	log "github.com/sirupsen/logrus"

	"github.com/valyala/fasttemplate"
)

func Parse(ts string, c shared.CheckConfig, nID, field string) string {

	st, err := fasttemplate.NewTemplate(ts, "{", "}")
	if err != nil {
		log.Errorf("error creating template from '%s' for notification ID: %s, err: %v", field, nID, err)
	}
	s := st.ExecuteFuncString(createParser(c))
	return s
}

func createParser(c shared.CheckConfig) func(w io.Writer, tag string) (int, error) {
	return func(w io.Writer, tag string) (int, error) {
		switch tag {
		case "check":
			return w.Write([]byte(c.Check))
		case "params":
			return w.Write(spaceIfVal(c.Params))
		case "headers":
			return w.Write(spaceIfVal(c.Headers))
		case "look_for":
			return w.Write(spaceIfVal(c.LookFor))
		case "response":
			return w.Write([]byte(c.Response))
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

func spaceIfVal(s string) []byte {
	if len(s) > 0 {
		return []byte(" " + s)
	}
	return []byte{}

}
