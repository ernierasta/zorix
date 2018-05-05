package check

import (
	"strconv"
	"time"

	"github.com/ernierasta/zorix/check/cmd"
	"github.com/ernierasta/zorix/check/ping"
	"github.com/ernierasta/zorix/check/port"
	"github.com/ernierasta/zorix/check/web"
	"github.com/ernierasta/zorix/shared"
	log "github.com/sirupsen/logrus"
)

// Manager registers all available checks and launches them
type Manager struct {
	checks                   []shared.CheckConfig
	workers                  int
	requestedWorkers         map[string]worker
	quitTickerChannels       map[string]chan bool
	resultsChan              chan shared.CheckConfig
	httpTimeout, pingTimeout shared.Duration
	portTimeout              shared.Duration
}

// registerWorker adds worker to requestedWorkers map.
// map[string]worker where key is type name "web", "cmd", ...
// and worker is struct with worker interface, channel of checks as fields.
//
// Number of checks using this type of worker is used to determine how many
// workers of particular type are needed).
//
// It you adding new worker this is the only place you need to add code.
//
// This function is also guard, which exit apllication if unknown type is given.
func (cm *Manager) registerWorker(t string) {
	c := 1
	if w, ok := cm.requestedWorkers[t]; ok {
		c = w.checks + 1
	}

	switch t {
	case "web":
		cm.requestedWorkers["web"] = worker{worker: web.New(cm.httpTimeout, false), typeChan: make(chan shared.CheckConfig, len(cm.checks)), checks: c}
	case "insecureweb":
		cm.requestedWorkers["insecureweb"] = worker{worker: web.New(cm.httpTimeout, true), typeChan: make(chan shared.CheckConfig, len(cm.checks)), checks: c}
	case "cmd":
		cm.requestedWorkers["cmd"] = worker{worker: cmd.New(), typeChan: make(chan shared.CheckConfig, len(cm.checks)), checks: c}
	case "ping":
		cm.requestedWorkers["ping"] = worker{worker: ping.New(cm.pingTimeout), typeChan: make(chan shared.CheckConfig, len(cm.checks)), checks: c}
	case "port":
		cm.requestedWorkers["port"] = worker{worker: port.New(cm.portTimeout), typeChan: make(chan shared.CheckConfig, len(cm.checks)), checks: c}
	default:
		log.Fatalf("check.registerWorker: unknown worker type: '%s', check config file.", t)
	}
}

// NewManager initializes Manager.
// checks: is slice of check params from config file
// workers: is number of all shared.Worker concurrent workers for selected worker type
// for example, 1 means: 1 web worker, 1 ping worker, ...
func NewManager(cc shared.CMConfig) *Manager {
	return &Manager{
		checks:             cc.Checks,
		workers:            cc.Workers,
		requestedWorkers:   make(map[string]worker),
		quitTickerChannels: make(map[string]chan bool),
		resultsChan:        cc.ResultsChan,
		httpTimeout:        cc.HTTPTimeout,
		pingTimeout:        cc.PingTimeout,
		portTimeout:        cc.PortTimeout,
	}
}

// Register registers required workers, based on checks.
// Add check numeric id here.
func (cm *Manager) Register() {
	for _, c := range cm.checks {
		cm.registerWorker(c.Type)
	}
}

// Run monitoring workers.
func (cm *Manager) Run() {

	// for every defined check create ticker, which will periodically create jobs for workers
	// send job to appropriate channel
	for _, c := range cm.checks {
		cm.quitTickerChannels[c.ID] = make(chan bool) // quit channel is unique for every ticker
		go cm.runTicker(c, cm.requestedWorkers[c.Type].typeChan, cm.quitTickerChannels[c.ID])
	}

	for name, wrk := range cm.requestedWorkers {
		for w := 1; w <= cm.workers && w <= wrk.checks; w++ {
			workerName := name + strconv.Itoa(w)
			go wrk.start(workerName, cm.resultsChan)
		}
	}
}

// runTicker writes data to request channel for every request.
// It periodically sends CheckConfig to appropriate channel.
func (cm *Manager) runTicker(c shared.CheckConfig, typeChan chan shared.CheckConfig, quitChan chan bool) {
	if c.Repeat.Duration <= 0 {
		log.WithFields(log.Fields{"check_id": c.ID, "repeat": c.Repeat}).Error("check.runTicker: non positive repeat interval, setting 60s.")
		c.Repeat.ParseDuration("1m")
	}
	ticker := time.NewTicker(c.Repeat.Duration)
	for {
		select {
		case <-quitChan:
			ticker.Stop()
			return
		default:
			typeChan <- c
			<-ticker.C

		}
	}
}

// worker is helper type, every worker type has its own implementation,
// job channel and amount of checks processed by this type of worker.
type worker struct {
	worker   shared.Worker
	typeChan chan shared.CheckConfig
	checks   int
}

// start will realize actual check. This method should run as goroutine.
// Method will call specific worker implementation and send data to resultsChan.
func (w *worker) start(id string, resultsChan chan shared.CheckConfig) {
	log.WithFields(log.Fields{"worker_id": id}).Info("starting some work ...")
	for c := range w.typeChan {
		code, body, time, err := w.worker.Send(c)
		c.ReturnedCode = code
		c.ReturnedTime = time
		c.Response = body
		if err != nil {
			c.Error = err
		}
		resultsChan <- c
	}
}
