package mail

import (
	"os"
	"testing"

	"github.com/ernierasta/zorix/shared"
)

// not unit test, this is integration test
// to run it: invoke it like that:
// MSERVER= MUSER= MPASS= MTO= go test

var (
	server = os.ExpandEnv("$MSERVER")
	user   = os.ExpandEnv("$MUSER")
	pass   = os.ExpandEnv("$MPASS")
	to     = []string{os.ExpandEnv("$MTO")}
	n, n2  shared.NotifConfig
	c      shared.CheckConfig
)

func Init() {
	n = shared.NotifConfig{
		Server: server,
		User:   user,
		Pass:   pass,
		Port:   587,
		To:     to,
	}
	n2 = n
	n2.Port = 465
	c = shared.CheckConfig{}
}

func TestMail_Send(t *testing.T) {
	Init()
	type args struct {
		c shared.CheckConfig
		n shared.NotifConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"timeouting send", args{c, shared.NotifConfig{Server: "server.com", User: "a", Pass: "b", To: []string{"to"}}}, true},
		{"not existing srv", args{c, shared.NotifConfig{Server: "rntuylmj320n290n03k093km43209d2.com", User: "a", Pass: "b", To: []string{"to"}}}, true},
		{"simple submission send", args{c, n}, false},
		{"simple tls send", args{c, n2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.n.Server == "" {
				t.Skip("you have to set enviroment vars, invoke it like this:\nMSERVER= MUSER= MPASS= MTO= go test\n")
			}
			m := &Mail{}
			t.Log(tt.args.n)
			if err := m.Send(tt.args.c, tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("Mail.Send() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				t.Log(err)
			}
		})
	}
}
