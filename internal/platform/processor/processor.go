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
				//log.WithFields(log.Fields{"id": c.ID, "check": c.Check, "code:": c.ReturnedCode, "time": c.ReturnedTime, "fails": p.checks[c.ID].Fails, "allowed_fails": c.AllowedFails, "slows": p.checks[c.ID].Slowdowns, "allowed_slows": c.AllowedSlows}).Debug("p.Listen: new result comes to precessor")
				log.WithFields(log.Fields{"check": c.Check}).Infof("t: %d, c: %d, fails: %d, slows: %d", c.ReturnedTime, c.ReturnedCode, p.checks[c.ID].Fails, p.checks[c.ID].Slowdowns)
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
		if p.checks[id].Fails == p.checks[id].AllowedFails && p.checks[id].Fails != 0 { //TODO: ref: should check only for r.Failure, new logic: if fails == AllowedFails +1
			log.Debugf("p.notify: f == allowed & not 0, %s sent to generator", id)
			p.notifyGenerator(*p.checks[id], false)
		} else if p.checks[id].RecoveryFailure {
			log.Debugf("p.notify: recovery == true, %s sent to generator", id)
			p.notifyGenerator(*p.checks[id], true)
		}
	}

	if len(p.checks[id].NotifySlow) > 0 {
		if p.checks[id].Slowdowns == p.checks[id].AllowedSlows && p.checks[id].Slowdowns != 0 {
			p.notifyGenerator(*p.checks[id], false)
		} else if p.checks[id].RecoverySlow {
			p.notifyGenerator(*p.checks[id], true)
		}
	}
}

// notifyGenerator will create notification for every required notification type
// (f.e: mail, jabber, ...)
// note, that notifyGenerator takes CheckConfig by value, so it copied. It is
// intensional, to be sure, that struct will not change before (if) notified.
func (p *Processor) notifyGenerator(c shared.CheckConfig, isRecovery bool) {
	notifs := p.notifsCleanup(&c)
	for _, nID := range notifs {

		schedule := p.schedule(&c, nID)
		cnID := fmt.Sprintf("%s_%s", c.ID, nID) // make unique ID string for this notification
		if isRecovery {
			// if recovery is BEFORE creating goroutine, it will stuck, send only if channel created = goroutine exists
			if _, ok := p.recoveryChans[cnID]; ok {
				log.Debugf("p.notifyGenerator: recovery message for %s (notification: %s) sending to the channel: recoveryChans[%s]", c.ID, nID, cnID)
				p.recoveryChans[cnID] <- true // every CheckConfig's notification has uniq quit channel
			}
		} else {
			log.Debugf("p.notifyGenerator: start go notificationTimer with recoveryChans[%s] for %s (notification: %s)", cnID, c.ID, nID)
			p.notifChan <- shared.NotifiedCheck{CheckConfig: c, NotificationID: nID} // send first notification directly
			p.recoveryChans[cnID] = make(chan bool, 1)
			go p.notificationTimer(c.ID, schedule, p.notifications[nID].ID, p.notifChan, p.recoveryChans[cnID])
		}
	}
}

// notificationTimer is running as goroutine for every CheckConfigs notification. There is always
// as many instances as CheckConfig has notifications set.
// It takes CheckConfig and schedule in form [ 1m, 5m, 10m ], where last interval is repeated until the end.
// If last interval is 0s, then it will stop notifications and terminate goroutine.
// Recovery message will be sent and will also terminate goroutine.
//
// TODO: getting c from p.checks is problematic, it can be already in different state (f.e.: RecoveryFailure: false for recovery message), we also cannot pass c here, because it will be reused for every notifications.
// Solutions:
// - make notifications concurent (channel will never be full)
// - use additional chanell here and send all check results after 1. notification to it, then check for Recovery check result, and keep it until we are sure notifications has been sent ... sounds overcomplicated
// Spliting CheckConfig and result data will solve this?
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
				if timer != nil {
					timer.Stop()
				}
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

// schedule returns notification schedule based on failure type (Slow or Fail)
func (p *Processor) schedule(c *shared.CheckConfig, nID string) []shared.Duration {
	if c.Failure || c.RecoveryFailure {
		return p.notifications[nID].RepeatFail
	}
	return p.notifications[nID].RepeatSlow
}

// notifsCleanup validates all defined notifications and returns slice of
// valid notifications based on failure type: Slow, Fail
func (p *Processor) notifsCleanup(c *shared.CheckConfig) []string {
	// check notification list, clean it if needed
	if c.Failure || c.RecoveryFailure {
		return p.validateNotifyIDList(c.NotifyFail)
	} else if c.Slow || c.RecoverySlow {
		return p.validateNotifyIDList(c.NotifySlow)
	}

	log.Errorf("processor.notifyGenerator: unknown failure type (not Fail or Slow)")
	return []string{}
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
