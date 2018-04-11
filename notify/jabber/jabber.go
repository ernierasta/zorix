package jabber

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/ernierasta/zorix/shared"
	xmpp "github.com/mattn/go-xmpp"
	log "github.com/sirupsen/logrus"
)

// Jabber notifier
type Jabber struct {
}

// Send sends jabber message to server.
// Currently it always log in/send/logout, maybe
// we should keep connection open?
func (j *Jabber) Send(c shared.CheckConfig, n shared.NotifConfig) error {

	options := xmpp.Options{
		Host:          fmt.Sprintf("%s:%d", n.Server, n.Port),
		User:          n.User,
		Password:      n.Pass,
		NoTLS:         n.IgnoreCert,
		Debug:         false,
		Session:       false,
		Status:        "xa",
		StatusMessage: "",
	}

	xmpp.DefaultConfig = tls.Config{
		ServerName:         n.Server,
		InsecureSkipVerify: n.IgnoreCert,
	}

	talk, err := options.NewClient()
	if err != nil {
		return fmt.Errorf("jabber.Send: wrong options, err: %v", err)
	}

	quitCh := make(chan bool)
	var recvErr error // this would be normaly terrible idea, but we have only one goroutine ...
	go func(quit chan bool) {
		for {
			select {
			case <-quit:
				recvErr = nil
				return
			default:
				_, recvErr = talk.Recv()
			}
		}
	}(quitCh)

	if recvErr != nil {
		return fmt.Errorf("jabber.Send: error receiving massages, err: %v", err)
	}
	for _, recipient := range n.To {
		i, err := talk.Send(xmpp.Chat{Remote: recipient, Type: "chat", Text: n.Text})
		log.Debugf("jabber.Send: message sent, %d bytes went to server\n", i)
		if err != nil {
			return fmt.Errorf("jabber.Send: error sending message to %s, err: %v", recipient, err)
		}
	}
	quitCh <- true
	time.Sleep(500 * time.Millisecond)
	talk.Close()
	return nil
}
