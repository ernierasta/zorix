package notify

import (
	"fmt"
	"log"

	"github.com/ernierasta/spock/shared"
	"github.com/ernierasta/spock/shared/mail"
	"github.com/valyala/fasttemplate"
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
	log.Println("start waiting for notifications")
	go func() {
		for {
			select {
			case ncheck := <-m.notifChan:
				ncheck.Check.Debug = append(ncheck.Check.Debug, "notify.Listen")
				//log.Printf("we got notification! %+v", ncheck)

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
		n.Subject = m.parseTemplate(n.SubjectFail, c, n.ID, "subject_fail")
		n.Text = m.parseTemplate(n.TextFail, c, n.ID, "text_fail")
	case c.Slow:
		n.Subject = m.parseTemplate(n.SubjectSlow, c, n.ID, "subject_slow")
		n.Text = m.parseTemplate(n.TextSlow, c, n.ID, "text_slow")
	case c.RecoveryFailure:
		n.Subject = m.parseTemplate(n.SubjectFailOK, c, n.ID, "subject_ok")
		n.Text = m.parseTemplate(n.TextFailOK, c, n.ID, "text_ok")
	case c.RecoverySlow:
		n.Subject = m.parseTemplate(n.SubjectSlowOK, c, n.ID, "subject_ok")
		n.Text = m.parseTemplate(n.TextSlowOK, c, n.ID, "text_ok")
	default:
		log.Printf("unknown notification, no known condition found, %+v", c)
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
		log.Printf("unknown notification type: '%s'. Check config file.", n.Type)
	}
}

func (m *Manager) parseTemplate(ts string, c shared.Check, nID, field string) string {

	st, err := fasttemplate.NewTemplate(ts, "{", "}")
	if err != nil {
		log.Printf("error creating template from '%s' for notification ID: %s, err: %v", field, nID, err)
	}
	s := st.ExecuteFuncString(createParser(c))
	return s
}
