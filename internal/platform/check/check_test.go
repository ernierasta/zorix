package check

import (
	"reflect"
	"testing"
	"time"

	"github.com/ernierasta/zorix/shared"
)

func TestManager_registerWorker(t *testing.T) {
	type fields struct {
		checks           []shared.CheckConfig
		workers          int
		requestedWorkers map[string]worker
		resultsChan      chan shared.CheckConfig
		httpTimeout      shared.Duration
		pingTimeout      shared.Duration
		portTimeout      shared.Duration
	}
	type args struct {
		t string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &Manager{
				checks:           tt.fields.checks,
				workers:          tt.fields.workers,
				requestedWorkers: tt.fields.requestedWorkers,
				resultsChan:      tt.fields.resultsChan,
				httpTimeout:      tt.fields.httpTimeout,
				pingTimeout:      tt.fields.pingTimeout,
				portTimeout:      tt.fields.portTimeout,
			}
			cm.registerWorker(tt.args.t)
		})
	}
}

func TestNewManager(t *testing.T) {
	type args struct {
		cc shared.CMConfig
	}
	tests := []struct {
		name string
		args args
		want *Manager
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewManager(tt.args.cc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewManager() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestManager_Register(t *testing.T) {
	type fields struct {
		checks           []shared.CheckConfig
		workers          int
		requestedWorkers map[string]worker
		resultsChan      chan shared.CheckConfig
		httpTimeout      shared.Duration
		pingTimeout      shared.Duration
		portTimeout      shared.Duration
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &Manager{
				checks:           tt.fields.checks,
				workers:          tt.fields.workers,
				requestedWorkers: tt.fields.requestedWorkers,
				resultsChan:      tt.fields.resultsChan,
				httpTimeout:      tt.fields.httpTimeout,
				pingTimeout:      tt.fields.pingTimeout,
				portTimeout:      tt.fields.portTimeout,
			}
			cm.Register()
		})
	}
}

func TestManager_Run(t *testing.T) {
	type fields struct {
		checks           []shared.CheckConfig
		workers          int
		requestedWorkers map[string]worker
		resultsChan      chan shared.CheckConfig
		httpTimeout      shared.Duration
		pingTimeout      shared.Duration
		portTimeout      shared.Duration
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &Manager{
				checks:           tt.fields.checks,
				workers:          tt.fields.workers,
				requestedWorkers: tt.fields.requestedWorkers,
				resultsChan:      tt.fields.resultsChan,
				httpTimeout:      tt.fields.httpTimeout,
				pingTimeout:      tt.fields.pingTimeout,
				portTimeout:      tt.fields.portTimeout,
			}
			cm.Run()
		})
	}
}

func TestManager_runTicker(t *testing.T) {

	t1s := shared.Duration{}
	_ = t1s.ParseDuration("1s")
	t1m := shared.Duration{}
	_ = t1m.ParseDuration("1m")

	type fields struct {
		checks           []shared.CheckConfig
		workers          int
		requestedWorkers map[string]worker
		resultsChan      chan shared.CheckConfig
		httpTimeout      shared.Duration
		pingTimeout      shared.Duration
		portTimeout      shared.Duration
	}
	type args struct {
		c        shared.CheckConfig
		typeChan chan shared.CheckConfig
		quitChan chan bool
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		closeAfter time.Duration
		wantRepeat shared.Duration
		wantAmount int
	}{
		{
			name: "no Duration given",
			fields: fields{
				checks:           []shared.CheckConfig{},
				workers:          0,
				requestedWorkers: map[string]worker{},
				resultsChan:      make(chan shared.CheckConfig, 2),
			},
			args: args{
				c:        shared.CheckConfig{ID: "t1"},
				typeChan: make(chan shared.CheckConfig, 2),
				quitChan: make(chan bool, 1),
			},
			closeAfter: 1 * time.Second,
			wantRepeat: t1m,
			wantAmount: 1,
		},
		{
			name: "1s Duration",
			fields: fields{
				checks:           []shared.CheckConfig{},
				workers:          0,
				requestedWorkers: map[string]worker{},
				resultsChan:      make(chan shared.CheckConfig, 2),
			},
			args: args{
				c:        shared.CheckConfig{ID: "t2", Repeat: t1s},
				typeChan: make(chan shared.CheckConfig, 2),
				quitChan: make(chan bool, 1),
			},
			closeAfter: 2 * time.Second,
			wantRepeat: t1s,
			wantAmount: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cm := &Manager{
				checks:           tt.fields.checks,
				workers:          tt.fields.workers,
				requestedWorkers: tt.fields.requestedWorkers,
				resultsChan:      tt.fields.resultsChan,
				httpTimeout:      tt.fields.httpTimeout,
				pingTimeout:      tt.fields.pingTimeout,
				portTimeout:      tt.fields.portTimeout,
			}

			go cm.runTicker(tt.args.c, tt.args.typeChan, tt.args.quitChan)
			start := time.Now()
			i := 0
		l:
			for {
				select {
				case c := <-tt.args.typeChan:
					i++
					if c.Repeat != tt.wantRepeat {
						t.Errorf("wrong interwal, want: %v, got: %v", tt.wantRepeat, c.Repeat)
					}
				default:
					if time.Since(start) > tt.closeAfter {
						tt.args.quitChan <- true
						break l
					}
					time.Sleep(500 * time.Millisecond)
				}
			}
			if i != tt.wantAmount {
				t.Errorf("wrong amount of ticks returned, want: %d, got: %d", tt.wantAmount, i)
			}
		})
	}
}

func Test_worker_start(t *testing.T) {
	type fields struct {
		worker   shared.Worker
		typeChan chan shared.CheckConfig
		checks   int
	}
	type args struct {
		id          string
		resultsChan chan shared.CheckConfig
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &worker{
				worker:   tt.fields.worker,
				typeChan: tt.fields.typeChan,
				checks:   tt.fields.checks,
			}
			w.start(tt.args.id, tt.args.resultsChan)
		})
	}
}
