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
	resultsChan              chan shared.CheckConfig
	webChan, iwebChan        chan shared.CheckConfig
	pingChan, cmdChan        chan shared.CheckConfig
	portChan                 chan shared.CheckConfig
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
// This function is also guard, which exit apllication if unknown type is given.
func (cm *Manager) registerWorker(t string) {
	c := 1
	if w, ok := cm.requestedWorkers[t]; ok {
		c = w.Checks + 1
	}

	switch t {
	case "web":
		cm.requestedWorkers["web"] = worker{Worker: web.New(cm.httpTimeout, false), Chan: cm.webChan, Checks: c}
	case "insecureweb":
		cm.requestedWorkers["insecureweb"] = worker{Worker: web.New(cm.httpTimeout, true), Chan: cm.iwebChan, Checks: c}
	case "cmd":
		cm.requestedWorkers["cmd"] = worker{Worker: cmd.New(), Chan: cm.cmdChan, Checks: c}
	case "ping":
		cm.requestedWorkers["ping"] = worker{Worker: ping.New(cm.pingTimeout), Chan: cm.pingChan, Checks: c}
	case "port":
		cm.requestedWorkers["port"] = worker{Worker: port.New(cm.portTimeout), Chan: cm.portChan, Checks: c}
	default:
		log.Fatalf("check.registerWorker: unknown worker type: '%s', check config file.", t)
	}
}

// worker is helper type, every worker type has its own implementation,
// job channel and amount of checks processed by this type of worker.
type worker struct {
	Worker shared.Worker
	Chan   chan shared.CheckConfig
	Checks int
}

// NewManager initializes Manager.
// checks: is slice of check params from config file
// workers: is number of all shared.Worker concurrent workers for selected worker type
// for example, 1 means: 1 web worker, 1 ping worker, ...
func NewManager(cc shared.CMConfig) *Manager {
	return &Manager{
		checks:           cc.Checks,
		workers:          cc.Workers,
		requestedWorkers: make(map[string]worker),
		resultsChan:      cc.ResultsChan,
		webChan:          make(chan shared.CheckConfig, len(cc.Checks)),
		iwebChan:         make(chan shared.CheckConfig, len(cc.Checks)),
		cmdChan:          make(chan shared.CheckConfig, len(cc.Checks)),
		pingChan:         make(chan shared.CheckConfig, len(cc.Checks)),
		portChan:         make(chan shared.CheckConfig, len(cc.Checks)),
		httpTimeout:      cc.HTTPTimeout,
		pingTimeout:      cc.PingTimeout,
		portTimeout:      cc.PortTimeout,
	}
}

// Register registers required workers, based on checks.
// Add check numeric id here.
func (cm *Manager) Register() {
	for _, c := range cm.checks {
		cm.registerWorker(c.Type)
		//no checking on start: log.Println(cm.requestedWorkers[c.Type].Worker.Send(c))
	}
}

// Run monitoring workers.
func (cm *Manager) Run() {

	// for every defined check create ticker, which will periodically create jobs for workers
	// send job to appropriate channel (webChan, cmdChan, ...)
	for _, c := range cm.checks {
		go cm.createTicker(c, cm.requestedWorkers[c.Type].Chan)
	}

	for name, wrk := range cm.requestedWorkers {
		for w := 1; w <= cm.workers && w <= wrk.Checks; w++ {
			workerName := name + strconv.Itoa(w)
			go cm.startWorker(workerName, wrk.Worker, wrk.Chan, cm.resultsChan)
		}
	}
}

//A time ticker writes data to request channel for every request.CheckEvery seconds
func (cm *Manager) createTicker(c shared.CheckConfig, checksChan chan shared.CheckConfig) {
	ticker := time.NewTicker(c.Repeat.Duration)
	for {
		checksChan <- c
		<-ticker.C
	}
}

// startWorker will realize actual check. This method should run as goroutine.
// Method will call specific worker implementation and send data to resultsChan.
func (cm *Manager) startWorker(id string, worker shared.Worker, checksChan, resultsChan chan shared.CheckConfig) {
	log.WithFields(log.Fields{"worker_id": id}).Info("starting some work ...")
	for c := range checksChan {
		code, body, time, err := worker.Send(c)
		c.ReturnedCode = code
		c.ReturnedTime = time
		c.Response = body
		if err != nil {
			c.Error = err
		}
		resultsChan <- c
	}
}
