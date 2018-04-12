// Package processor processes checks results It keeps states of checks and
// dispatch sending messages if needed.
// If desired, it can store results in db.
package processor

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ernierasta/zorix/shared"
	log "github.com/sirupsen/logrus"
)

// Processor analyze and store check results.
// If needed, it sends notifications.
type Processor struct {
	resultChan    chan shared.CheckConfig
	notifChan     chan shared.NotifiedCheck
	recoveryChans map[string]chan bool
	checks        map[string]*shared.CheckConfig
	notifications map[string]*shared.NotifConfig
	mutex         *sync.Mutex
}

// New returns processor instance
func New(resultChan chan shared.CheckConfig, notifChan chan shared.NotifiedCheck, checkAmmount int, notifications []shared.NotifConfig) *Processor {
	notes := make(map[string]*shared.NotifConfig, len(notifications))
	for i, n := range notifications {
		notes[n.ID] = &notifications[i]
	}

	return &Processor{
		resultChan:    resultChan,
		notifChan:     notifChan,
		recoveryChans: make(map[string]chan bool, checkAmmount*len(notifications)),
		checks:        make(map[string]*shared.CheckConfig, checkAmmount),
		notifications: notes,
		mutex:         &sync.Mutex{},
	}
}

// Listen starts listening for notifications.
func (p *Processor) Listen() {
	log.Debug("p.Listen:", "start waiting for results ...")
	go func() {
		for {
			select {
			case c := <-p.resultChan:
				c = p.analyze(c)
				p.updateCheckResult(c)
				log.WithFields(log.Fields{"id": c.ID, "check": c.Check, "code:": c.ReturnedCode, "time": c.ReturnedTime, "fails": p.checks[c.ID].Fails, "allowed_fails": c.AllowedFails, "slows": p.checks[c.ID].Slowdowns, "allowed_slows": c.AllowedSlows}).Debug("p.Listen: new result comes to precessor")
				//log.Debug("response:" + c.Response)
				p.notify(c.ID)
			}

		}
	}()
}

func (p *Processor) analyze(r shared.CheckConfig) shared.CheckConfig {

	if r.Error != nil {
		r.Fails = 1
	}
	if r.ReturnedCode != r.ExpectedCode {
		r.Fails = 1
		if r.Error == nil {
			r.Error = fmt.Errorf("wrong response code")
		}
	}
	if r.LookFor != "" && !strings.Contains(r.Response, r.LookFor) {
		r.Fails = 1
		if r.Error == nil {
			r.Error = fmt.Errorf("response does not contain: %s", r.LookFor)
		}
	}

	if r.ReturnedTime > r.ExpectedTime {
		r.Slowdowns = 1
		if r.Error == nil {
			r.Error = fmt.Errorf("slow response")
		}
	}

	return r

}

// updateCheckResult will store actual result in checks map.
// It will increment fail or slowdown counter if needed.
func (p *Processor) updateCheckResult(r shared.CheckConfig) {
	// if this check failed, add amount of failures to check data
	if r.Fails == 1 {
		if prevResult, ok := p.checks[r.ID]; ok {
			r.Fails = prevResult.Fails + 1
		}
	}
	if r.Slowdowns == 1 {
		if prevResult, ok := p.checks[r.ID]; ok {
			r.Slowdowns = prevResult.Slowdowns + 1
		}
	}

	// detect recovery situation
	if prevResult, ok := p.checks[r.ID]; ok {
		if r.Slowdowns == 0 && r.Fails == 0 {
			if prevResult.Failure {
				r.RecoveryFailure = true
				r.Timestamp = time.Now()
			} else if prevResult.Slow {
				r.RecoverySlow = true
				r.Timestamp = time.Now()
			}
		}
	}

	// mark Check as Fails and add timestamp
	if r.Fails >= r.AllowedFails {
		r.Timestamp = time.Now()
		r.Failure = true
	}

	// mark Check as slow and add timestamp
	if r.Slowdowns >= r.AllowedSlows {
		r.Timestamp = time.Now()
		r.Slow = true
	}

	// add or update result
	// if check do not failed or is not slow it will overwrite previous data completly
	p.checks[r.ID] = &r
}

// notify checks if notification is needed, if so, send them it to notifyGenerator.
// We are sending only CheckConfigs with NotifyFail, NotifySlow, RecoveryFailure, RecoverySlow,
// those messages are sent always only once for given CheckConfig.
func (p *Processor) notify(id string) {

	if len(p.checks[id].NotifyFail) > 0 {
		if p.checks[id].Fails == p.checks[id].AllowedFails && p.checks[id].Fails != 0 {
			log.Debugf("p.notify: f == allowed & not 0, %s sent to generator", id)
			p.notifyGenerator(id, false)
		} else if p.checks[id].RecoveryFailure {
			log.Debugf("p.notify: recovery == true, %s sent to generator", id)
			p.notifyGenerator(id, true)
		}
	}

	if len(p.checks[id].NotifySlow) > 0 {
		if p.checks[id].Slowdowns == p.checks[id].AllowedSlows && p.checks[id].Slowdowns != 0 {
			p.notifyGenerator(id, false)
		} else if p.checks[id].RecoverySlow {
			p.notifyGenerator(id, true)
		}
	}
}

// notifyGenerator will create notification for every required notification type
// (f.e: mail, jabber, ...)
func (p *Processor) notifyGenerator(cID string, isRecovery bool) {

	source := []string{}
	if p.checks[cID].Failure || p.checks[cID].RecoveryFailure {
		source = p.validateNotifyIDList(p.checks[cID].NotifyFail)
	} else if p.checks[cID].Slow || p.checks[cID].RecoverySlow {
		source = p.validateNotifyIDList(p.checks[cID].NotifySlow)
	} else {
		log.Errorf("processor.notifyGenerator: unknown failure type (not Fail or Slow)")
	}

	for _, nID := range source {

		schedule := []shared.Duration{}
		if p.checks[cID].Failure || p.checks[cID].RecoveryFailure {
			schedule = p.notifications[nID].RepeatFail
		} else {
			schedule = p.notifications[nID].RepeatSlow //TODO: segfault
		}

		cnID := fmt.Sprintf("%s_%s", cID, nID) // make unique ID string for this notification
		if isRecovery {
			// if recovery is BEFORE creating goroutine, it will stuck, send only if channel created = goroutine exists
			if _, ok := p.recoveryChans[cnID]; ok {
				p.recoveryChans[cnID] <- true // every CheckConfig's notification has uniq quit channel
				log.Debugf("p.notifyGenerator: recovery message for %s (notification: %s) sent to channel: recoveryChans[%s]", cID, nID, cnID)
			}
		} else {
			p.notifChan <- shared.NotifiedCheck{CheckConfig: *p.checks[cID], NotificationID: nID} // send first notification directly
			p.recoveryChans[cnID] = make(chan bool, 1)
			log.Debugf("p.notifyGenerator: start go notificationTimer with recoveryChans[%s] for %s (notification: %s)", cnID, cID, nID)
			go p.notificationTimer(cID, schedule, p.notifications[nID].ID, p.notifChan, p.recoveryChans[cnID])
		}
	}
}

// notificationTimer is running as goroutine for every CheckConfig. There is always
// only one instance for CheckConfig.
// It takes CheckConfig and schedule in form [ 1m, 5m, 10m ], where last interval is repeated until the end.
// If last interval is 0s, then it will stop notifications and terminate goroutine.
// Recovery message will be sent and will also terminate goroutine.
func (p *Processor) notificationTimer(cID string, schedule []shared.Duration, nID string, notifChan chan<- shared.NotifiedCheck, recovery chan bool) {
	var timer *time.Timer
	//for _, interval := range schedule {
	cnt := 0
	for {
		runTimer := true
		sent := false
		interval := schedule[cnt]
		isLast := cnt == len(schedule)-1
		isStopNotifications := schedule[cnt].Duration == 0
		if isLast {
			if isStopNotifications {
				runTimer = false
			} else {
				cnt-- // decrement counter always keep last schedule value
			}
		}
		cnt++
		log.Debug("p.notificationTimer: starting inner for loop for interval: ", interval)
		for {
			select {
			case <-recovery:
				log.Debug("p.notificationTimer: sending recovery to channel")
				p.mutex.Lock()
				notifChan <- shared.NotifiedCheck{CheckConfig: *p.checks[cID], NotificationID: nID}
				timer.Stop()
				p.mutex.Unlock()
				log.Debugf("p.notificationTimer: recovery message received for %s (notification: %s)", cID, nID)
				return
			default:
				if runTimer {
					timer = time.NewTimer(interval.Duration)

					go func() { // we need to be able to cancel timer if recovery came
						runTimer = false
						<-timer.C
						log.Debugf("p.notificationTimer: sending message %s(nID: %s) to notifChan", cID, nID)
						p.mutex.Lock()
						notifChan <- shared.NotifiedCheck{CheckConfig: *p.checks[cID], NotificationID: nID}
						p.mutex.Unlock()
						sent = true
					}()

				}
			}
			if sent {
				log.Debugf("notificationTimer: sent problem notification for %s (notification: %s)", cID, nID)
				log.Debug("notificationTimer: breaking from inner for loop")
				break
			}
			time.Sleep(1 * time.Second)
		}
	}
}

// validateNotifyIDList checks if all notification in slice
// are valid notification ID's. If not, remove them from slice.
func (p *Processor) validateNotifyIDList(ss []string) []string {
	res := []string{}
	for _, nID := range ss {
		if _, ok := p.notifications[nID]; ok {
			res = append(res, nID)
		} else {
			log.WithFields(log.Fields{"given": nID}).Warn("p.validateNotifyIDList: found strange NotificationID")
		}
	}
	return res
}
