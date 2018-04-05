package notify

import (
	"fmt"

	"github.com/ernierasta/zorix/notify/mail"
	"github.com/ernierasta/zorix/shared"
	"github.com/ernierasta/zorix/template"

	log "github.com/sirupsen/logrus"
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
				//log.Debug("we got notification! %+v", ncheck)
				m.send(ncheck)
			}

		}
	}()
}

// send determines notification and calls methods that sets subject and title
// and sends Check to dispatch method.
func (m *Manager) send(nc shared.NotifiedCheck) {
	n := m.notifications[nc.NotificationID]
	n = m.setSubjectAndText(nc.Check, n)
	m.dispatch(nc.Check, n)
}

// setSubjectAndText determines notification type (fail, slow, failOK, slowOK), gets parsed subject
// and text and sets them in returned Notification struct.
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

// dispach determines which plugin should be called
func (m *Manager) dispatch(c shared.Check, n *shared.Notification) {
	switch n.Type {
	case "mail":
		mail.Send(c, *n)
	default:
		log.Errorf("programming error, check is to late: unknown notification type: '%s'. Check config file.", n.Type)
	}
}
