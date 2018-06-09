package template

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ernierasta/zorix/shared"
	log "github.com/sirupsen/logrus"

	"github.com/valyala/fasttemplate"
)

// Parse parses string using config and result variables.
// Usualy used for notification subject and body generation.
// It looks for {something} tags.
func Parse(ts string, c shared.CheckConfig, nID, field string) string {

	st, err := fasttemplate.NewTemplate(ts, "{", "}")
	if err != nil {
		log.Errorf("error creating template from '%s' for notification ID: %s, err: %v", field, nID, err)
	}
	s := st.ExecuteFuncString(CheckVarsParser(c))
	return s
}

// CheckVarsParser function replaces CheckConfig config and result data.
// Replaces all CheckConfig fields.
// All names are as in config files, exceptions:
// cID = check ID
// ctype = check type
// expected_code = code
// expected_time = time
//
// Those are result data:
// response_code, response_time, response, timestamp
func CheckVarsParser(c shared.CheckConfig) func(w io.Writer, tag string) (int, error) {
	return func(w io.Writer, tag string) (int, error) {
		switch tag {
		case "cID":
			return w.Write([]byte(c.ID))
		case "ctype":
			return w.Write([]byte(c.Type))
		case "check":
			return w.Write([]byte(c.Check))
		case "params":
			return w.Write(spaceIfVal(c.Params))
		case "headers":
			return w.Write(spaceIfVal(c.Headers))
		case "redirs":
			return w.Write([]byte(strconv.Itoa(c.Redirs)))
		case "repeat":
			return w.Write([]byte(c.Repeat.String()))
		case "method":
			return w.Write(spaceIfVal(c.Method))
		case "look_for":
			return w.Write(spaceIfVal(c.LookFor))
		case "response":
			return w.Write([]byte(c.Response))
		case "timestamp":
			return w.Write([]byte(c.Timestamp.Format("2.1.2006 15:04:05")))
		case "allowed_fails":
			return w.Write([]byte(strconv.Itoa(c.AllowedFails)))
		case "allowed_slows":
			return w.Write([]byte(strconv.Itoa(c.AllowedSlows)))
		case "notify_fail":
			return w.Write([]byte(strings.Join(c.NotifyFail, ", ")))
		case "notify_slow":
			return w.Write([]byte(strings.Join(c.NotifySlow, ", ")))
		case "response_code":
			return w.Write([]byte(fmt.Sprintf("%d", c.ReturnedCode)))
		case "response_time":
			return w.Write([]byte(fmt.Sprintf("%d", c.ReturnedTime)))
		case "expected_code":
			return w.Write([]byte(fmt.Sprintf("%d", c.ExpectedCode)))
		case "expected_time":
			return w.Write([]byte(fmt.Sprintf("%d", c.ExpectedTime)))
		case "error":
			if c.Error != nil {
				return w.Write([]byte(c.Error.Error()))
			}
			return w.Write([]byte(""))
			//TODO: add all fields from shared.Check
		default:
			return w.Write([]byte("{" + tag + "}"))
		}
	}
}

func spaceIfVal(s string) []byte {
	if len(s) > 0 {
		return []byte(" " + s)
	}
	return []byte{}

}

func spaceIfValI(i int) []byte {
	if i != 0 {
		return []byte(" " + strconv.Itoa(i))
	}
	return []byte{}

}

// ParseNotif parses notification values into ts.
func ParseNotif(ts string, n *shared.NotifConfig, field string) string {
	st, err := fasttemplate.NewTemplate(ts, "{", "}")
	if err != nil {
		log.Errorf("error creating template from '%s' for notification ID: %s, err: %v", field, n.ID, err)
	}
	s := st.ExecuteFuncString(NotifVarsParser(n))
	return s

}

// NotifVarsParser return parsing function.
// Replaces all CheckConfig fields.
// All names are as in config files, exceptions:
// nID = notification ID
// to = returns recipients separated by space
// to, = returns recipients separated by comma
func NotifVarsParser(n *shared.NotifConfig) func(w io.Writer, tag string) (int, error) {
	return func(w io.Writer, tag string) (int, error) {
		switch tag {
		case "nID":
			return w.Write([]byte(n.ID))
		case "server":
			return w.Write([]byte(n.Server))
		case "port":
			return w.Write([]byte(strconv.Itoa(n.Port)))
		case "user":
			return w.Write([]byte(n.User))
		case "pass":
			return w.Write([]byte(n.Pass))
		case "from":
			return w.Write([]byte(n.From))
		case "to":
			return w.Write([]byte(strings.Join(n.To, " ")))
		case "to,":
			return w.Write([]byte(strings.Join(n.To, ",")))
		case "subject_fail":
			return w.Write([]byte(n.SubjectFail))
		case "subject_slow":
			return w.Write([]byte(n.SubjectSlow))
		case "text_fail":
			return w.Write([]byte(n.TextFail))
		case "text_slow":
			return w.Write([]byte(n.TextSlow))
		case "subject_fail_ok":
			return w.Write([]byte(n.SubjectFailOK))
		case "subject_slow_ok":
			return w.Write([]byte(n.SubjectSlowOK))
		case "text_fail_ok":
			return w.Write([]byte(n.TextFailOK))
		case "text_slow_ok":
			return w.Write([]byte(n.TextSlowOK))
		case "subject":
			return w.Write([]byte(n.Subject))
		case "text":
			return w.Write([]byte(n.Text))
		case "no_recovery":
			return w.Write([]byte(strconv.FormatBool(n.NoRecovery)))
		default:
			return w.Write([]byte("{" + tag + "}"))

		}
	}
}

// ParseEnv parses given template and replace all
// occurence of ${var} to enviroment vars values.
func ParseEnv(ts string, ID string, field string) string {

	st, err := fasttemplate.NewTemplate(ts, "${", "}")
	if err != nil {
		log.Errorf("error creating template for check/notification ID: %s, field '%s', err: %v", ID, field, err)
	}

	s := st.ExecuteFuncString(ConfEnvVarsParser())

	return s
}

// ConfEnvVarsParser returns env parser
func ConfEnvVarsParser() func(w io.Writer, tag string) (int, error) {
	return func(w io.Writer, tag string) (int, error) {
		if val, ok := os.LookupEnv(tag); ok {
			return w.Write([]byte(val))
		}
		return w.Write([]byte(tag))
	}
}
