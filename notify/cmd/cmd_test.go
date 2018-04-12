package cmd

import (
	"testing"

	"github.com/ernierasta/zorix/shared"
)

func TestCmd_Send(t *testing.T) {
	type args struct {
		c shared.CheckConfig
		n shared.NotifConfig
	}
	tests := []struct {
		name    string
		j       Cmd
		args    args
		wantErr bool
	}{
		{"run echo", Cmd{}, args{shared.CheckConfig{}, shared.NotifConfig{Cmd: "echo Test"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := Cmd{}
			if err := j.Send(tt.args.c, tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("Cmd.Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
