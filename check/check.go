package check

import (
	"strconv"
	"time"

	"github.com/ernierasta/zorix/check/cmd"
	"github.com/ernierasta/zorix/check/web"
	"github.com/ernierasta/zorix/shared"
	log "github.com/sirupsen/logrus"
)

// Manager registers all availabile checks and lounches them
type Manager struct {
	checks           []shared.CheckConfig
	workers          int
	requestedWorkers map[string]worker
	resultsChan      chan shared.CheckConfig
	webChan, cmdChan chan shared.CheckConfig
	httpTimeout      shared.Duration
}

// worker is helper type, every worker type has its own implementation,
// job channel and ammount of checks processed by this type of worker.
type worker struct {
	Worker shared.Worker
	Chan   chan shared.CheckConfig
	Checks int
}

// NewManager initializes Manager.
// checks: is slice of check params from config file
// workers: is number of all shared.Worker concurent workers for selected worker type
// for example, 1 means: 1 web worker, 1 ping worker, ...
func NewManager(checks []shared.CheckConfig, workers int, resultsChan chan shared.CheckConfig, httpTimeout shared.Duration) *Manager {
	return &Manager{
		checks:           checks,
		workers:          workers,
		requestedWorkers: make(map[string]worker),
		resultsChan:      resultsChan,
		webChan:          make(chan shared.CheckConfig, len(checks)),
		cmdChan:          make(chan shared.CheckConfig, len(checks)),
		httpTimeout:      httpTimeout,
	}
}

// Register registers required workers, based on checks.
// Add check numeric id here.
func (cm *Manager) Register() {
	for i, c := range cm.checks {
		cm.registerWorker(c.Type)
		cm.checks[i].ID = i + 1
		//no checkin on start: log.Println(cm.requestedWorkers[c.Type].Worker.Send(c))
	}
}

// registerWorker adds worker to requestedWorkers map.
// map[string]worker where string is type name "web", "cmd", ...
//   and worker is struct with worker interface, channel of checks as fields
//   and number of checks using this type of worker
//   (used to determine how many workers of particullar type are needed).
//
// This funcion is also guard, which exit aplication if unknown type is given.
func (cm *Manager) registerWorker(t string) {
	c := 1
	if w, ok := cm.requestedWorkers[t]; ok {
		c = w.Checks + 1
	}

	switch t {
	case "web":
		cm.requestedWorkers["web"] = worker{Worker: web.New(cm.httpTimeout), Chan: cm.webChan, Checks: c}
	case "cmd":
		cm.requestedWorkers["cmd"] = worker{Worker: cmd.New(), Chan: cm.cmdChan, Checks: c}
	default:
		log.Fatalf("Manager.registerWorker: unknown worker type: '%s', check config file.", t)
	}
}

// Run monitoring workers.
func (cm *Manager) Run() {

	// for every defined check create ticker, which will periodically create jobs for workers
	// send job to apropriate channel (webChan, cmdChan, ...)
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
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			checksChan <- c
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

// startWorker will realize actual check. This method should run as gorutine.
// Method will call specific worker implementation and send data to resultsChan.
func (cm *Manager) startWorker(id string, worker shared.Worker, input, output chan shared.CheckConfig) {
	log.WithFields(log.Fields{"worker_id": id}).Info("starting some work ...")
	for c := range input {
		code, time, err := worker.Send(c)
		c.ReturnedCode = code
		c.ReturnedTime = time
		if err != nil {
			c.Error = err
		}
		output <- c
	}
}
