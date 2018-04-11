package jabber

import (
	"os"
	"testing"

	"github.com/ernierasta/zorix/shared"
)

var (
	jserver = os.ExpandEnv("$JSERVER")
	juser   = os.ExpandEnv("$JUSER")
	jpass   = os.ExpandEnv("$JPASS")
	jto     = os.ExpandEnv("$JTO")
)

func Init() {
	if juser == "" {
		t := &testing.T{}
		t.Errorf("No enviroment variables set. Run:\n JSERVER=%q JUSER=u JPASS=p JTO=r go test", "talk.google.com")
	}
}

func TestSend(t *testing.T) {
	type args struct {
		c shared.CheckConfig
		n shared.NotifConfig
	}
	tests := []struct {
		name string
		args args
	}{
		{"send message", args{c: shared.CheckConfig{}, n: shared.NotifConfig{Server: "talk.google.com", Port: 5223, User: os.ExpandEnv("$JUSER"), Pass: os.ExpandEnv("$JPASS"), To: []string{os.ExpandEnv("$JTO")}, Subject: "test message", Text: "hi, this is jabber test!"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := Jabber{}
			j.Send(tt.args.c, tt.args.n)
		})
	}
}
