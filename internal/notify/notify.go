package notify

import (
	"fmt"
	"time"

	"github.com/ernierasta/zorix/notify/cmd"
	"github.com/ernierasta/zorix/notify/jabber"
	"github.com/ernierasta/zorix/notify/mail"
	"github.com/ernierasta/zorix/shared"
	"github.com/ernierasta/zorix/template"

	log "github.com/sirupsen/logrus"
)

// NotificationModules is a map.
// string: name of notification module (f.e: mail, jabber, ...)
// value: Notifier interface, in fact concrete implementation.
// This is only place, where you add new modules.
var NotificationModules = map[string]shared.Notifier{
	"mail":   &mail.Mail{},
	"jabber": &jabber.Jabber{},
	"cmd":    &cmd.Cmd{},
} // TODO: do the same for checks

// Manager sets up notifications and send them if needed
type Manager struct {
	notifChan     chan shared.NotifiedCheck
	notifications map[string]*shared.NotifConfig
}

// NewManager creates new instance of notification manager
func NewManager(notifChan chan shared.NotifiedCheck, notifications []shared.NotifConfig) *Manager {
	notes := make(map[string]*shared.NotifConfig, len(notifications))
	for id, n := range notifications {
		notes[n.ID] = &notifications[id]
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
				m.send(ncheck)
			}

		}
	}()
}

// send determines notification and calls methods that sets subject and title
// and sends Check to dispatch method.
func (m *Manager) send(nc shared.NotifiedCheck) {
	n := m.notifications[nc.NotificationID]
	if isRecovery(nc.CheckConfig) && n.NoRecovery {
		return // do not send recovery if 'no_recovery = true' in notification settings
	}
	n = m.setSubjectAndText(nc.CheckConfig, n)
	n = m.parseNotificationCmd(nc.CheckConfig, n)
	m.dispatch(nc.CheckConfig, n)
}

// setSubjectAndText determines notification type (fail, slow, failOK, slowOK), gets parsed subject
// and text and sets them in returned Notification struct.
func (m *Manager) setSubjectAndText(c shared.CheckConfig, n *shared.NotifConfig) *shared.NotifConfig {
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
		c.Response = "removed output"
		log.Errorf("unknown notification, no known condition found, %+v", c)
		n.Subject = "Unknown notification"
		n.Text = fmt.Sprintf("Programming error, please send bug report containing"+
			" folowing:\nCheck: %+v\n\n Notification(id: %s): %+v\n", c, n.ID, n)
	}
	return n

}

// parseNotificationCmd parses notification Cmd string.
// first it replaces CheckConfig data, then any enviroment data.
func (m *Manager) parseNotificationCmd(c shared.CheckConfig, n *shared.NotifConfig) *shared.NotifConfig {
	n.Cmd = template.Parse(n.CmdTemplate, c, n.ID, "cmd")
	n.Cmd = template.ParseNotif(n.Cmd, n, "cmd")
	n.Cmd = template.ParseEnv(n.Cmd, n.ID, "cmd")
	return n
}

// dispach determines which plugin should be called
func (m *Manager) dispatch(c shared.CheckConfig, n *shared.NotifConfig) {
	if _, ok := NotificationModules[n.Type]; ok {
		err := NotificationModules[n.Type].Send(c, *n)
		if err != nil {
			log.Error(err)
		}

	} else {
		log.Errorf("programming error, check is to late: unknown notification type: '%s'. Check config file.", n.Type)
	}
}

//TestAll sends test message to all configured notifications.
func (m *Manager) TestAll() {
	fc := shared.CheckConfig{}
	fc.Timestamp = time.Now()

	for _, n := range m.notifications {
		n.Subject = "Test notification from ZoriX"
		n.Text = "Hi comrade!\nIf you are reading this, all went good.\nWe are glad you want to give ZoriX a try!\n\nWelcome in ZoriX community.\n\n Yours ZoriX"
		n.Cmd = "echo \"It works!\" > /tmp/zorix.test"
		log.Infof("notify.TestAll: trying to send '%s', check if it arrived!\n", n.ID)
		m.dispatch(fc, n)
	}
}

func isRecovery(c shared.CheckConfig) bool {
	return c.RecoveryFailure || c.RecoverySlow
}
