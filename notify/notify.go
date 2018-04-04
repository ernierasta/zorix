package notify

import (
	"fmt"
	"io"
	"strings"

	"github.com/ernierasta/zorix/log"
	"github.com/ernierasta/zorix/notify/mail"
	"github.com/ernierasta/zorix/shared"
	"github.com/ernierasta/zorix/template"
)

// Manager sets up notifications and send them if needed
type Manager struct {
	notifChan     chan shared.NotifiedCheck
	notifications map[string]*shared.Notification
}

// NewManager creates new instance of notification manager
func NewManager(notifChan chan shared.NotifiedCheck, notifications []shared.Notification) *Manager {
	notes := make(map[string]*shared.Notification, len(notifications))
	for _, n := range notifications {
		notes[n.ID] = &n
	}
	return &Manager{
		notifChan:     notifChan,
		notifications: notes,
	}
}

// Listen starts listening for notifications.
func (m *Manager) Listen() {
	log.Debug("start waiting for notifications")
	go func() {
		for {
			select {
			case ncheck := <-m.notifChan:
				ncheck.Check.Debug = append(ncheck.Check.Debug, "notify.Listen")
				//log.Debug("we got notification! %+v", ncheck)

				m.send(ncheck)
			}

		}
	}()
}

// sendAll takes slice of defined notificiations for this Check
// and runs them all.
// TODO: remove this, from now, every notificion will come from processor, for every type
func (m *Manager) send(c shared.NotifiedCheck) {
	n := m.notifications[c.NotificationID]
	n = m.setSubjectAndText(c.Check, n)
	m.dispatch(c.Check, n)
}

func (m *Manager) setSubjectAndText(c shared.Check, n *shared.Notification) *shared.Notification {
	switch {
	case c.Failure:
		n.Subject = template.Parse(n.SubjectFail, c, n.ID, "subject_fail")
		n.Text = template.Parse(n.TextFail, c, n.ID, "text_fail")
	case c.Slow:
		n.Subject = template.Parse(n.SubjectSlow, c, n.ID, "subject_slow")
		n.Text = template.Parse(n.TextSlow, c, n.ID, "text_slow")
	case c.RecoveryFailure:
		n.Subject = template.Parse(n.SubjectFailOK, c, n.ID, "subject_ok")
		n.Text = template.Parse(n.TextFailOK, c, n.ID, "text_ok")
	case c.RecoverySlow:
		n.Subject = template.Parse(n.SubjectSlowOK, c, n.ID, "subject_ok")
		n.Text = template.Parse(n.TextSlowOK, c, n.ID, "text_ok")
	default:
		log.Errorf("unknown notification, no known condition found, %+v", c)
		n.Subject = "Unknown notification"
		n.Text = fmt.Sprintf("Programming error, please send bug report containing"+
			" folowing:\nCheck: %+v\n\n Notification(id: %s): %+v\n", c, n.ID, n)
	}
	return n

}

func (m *Manager) dispatch(c shared.Check, n *shared.Notification) {
	switch n.Type {
	case "mail":
		mail.Send(c, *n)
	default:
		log.Errorf("programming error, check is to late: unknown notification type: '%s'. Check config file.", n.Type)
	}
}

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
