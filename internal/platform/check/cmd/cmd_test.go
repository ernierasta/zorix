package cmd

import (
	"testing"

	"github.com/ernierasta/zorix/shared"
)

var (
	c1 = shared.CheckConfig{
		Check:  "echo",
		Params: "hello",
	}
	c2 = shared.CheckConfig{
		Check:  "ssh",
		Params: `it 'top -bn2' | awk '/Cpu\(s\):/ { isLoaded=100-$8>50  } END {print isLoaded}'`,
		//Params: `it 'grep 'cpu ' /proc/stat' | awk '{isLoad=($2+$4)*100/($2+$4+$5)} END {print isLoad; exit isLoad}'`,
	}
)

func TestCmd_Send(t *testing.T) {
	type args struct {
		c shared.CheckConfig
	}
	tests := []struct {
		name    string
		args    args
		want    int
		want1   string
		wantErr bool
	}{
		{"simple", args{c1}, 200, "hello", false},
		{"ssh", args{c2}, 200, "a", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Cmd{}
			got, got1, _, err := w.Send(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cmd.Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Cmd.Send() got = '%v', want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Cmd.Send() got1 = '%v', want %v", got1, tt.want1)
			}
		})
	}
}
