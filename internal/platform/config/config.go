package config

// TODO: shouldn't we separate defaults, validation and normalization to separate
// files named by check/notification type? Hate it like it is now.

import (
	"fmt"

	"github.com/ernierasta/zorix/internal/template"

	"github.com/BurntSushi/toml"
)

const (
	Loglevel    = "warn"
	HTTPTimeout = 60
	PingTimeout = 60
	PortTimeout = 10

	CheckMethod         = "GET"
	CheckRepeatInterval = 60
	CheckAllowedSlows   = 2
	CheckAllowedFails   = 0

	WebExpectedCode = 200
	WebExpectedTime = 1000

	PortAllowedFails = 1
	PortAllowedSlows = 3
	PortExpectedTime = 150
	PortExpectedCode = 200

	PingAllowedFails = 1
	PingAllowedSlows = 3
	PingExpectedTime = 150
	PingExpectedCode = 0

	CmdExpectedTime = 5000
	CmdExpectedCode = 0

	NotifyType          = "mail"
	NotifyFailSubject   = "{check}{params} problem"
	NotifySlowSubject   = "{check}{params} slow"
	NotifyFailOKSubject = "{check}{params} ok"
	NotifySlowOKSubject = "{check}{params} ok"
	NotifyFailText      = "FAILURE:\n{check}{params}\nTime: {timestamp}\n\nResponse code: {response_code}\nError: {error}\n"
	NotifySlowText      = "SLOW RESPONSE:\n{check}{params}\nTime: {timestamp}\n\nResponse/Expected time: {response_time}/{expected_time}\n"
	NotifyFailOKText    = "RECOVERED:\n{check}{params}\nTime: {timestamp}\n\nResponse code: {response_code}\n"
	NotifySlowOKText    = "RECOVERED:\n{check}{params}\nTime: {timestamp}\n\nResponse/Expected time: {response_time}/{expected_time}\n"
)

var (
	FailSchedule = []string{"1m", "5m", "10m"}
	SlowSchedule = []string{"5m", "0"}
)

// Config represents whole configuration file parsed to stuct.
type Config struct {
	Global Global
	Notify Notify
	Check  Check
	File   string
}

func (c *Config) isNoChecksDefined() bool {
	return len(c.Check.Web) == 0 && len(c.Check.Port) == 0 && len(c.Check.Ping) == 0 && len(c.Check.Cmd) == 0
}

func (c *Config) isNoNotifyDefined() bool {
	return len(c.Notify.Mail) == 0 && len(c.Notify.Jabber) == 0 && len(c.Notify.Cmd) == 0
}

// New returns Config with config file defined
func New(file string) *Config {
	return &Config{File: file}
}

func (c *Config) Read() ([]string, error) {

	unexp := []string{}

	meta, err := toml.DecodeFile(c.File, c)
	if err != nil {
		return unexp, fmt.Errorf("in file %s: err: %s", c.File, err)
	}
	keys := meta.Undecoded()
	for _, key := range keys {
		unexp = append(unexp, key.String())
	}
	return unexp, nil

}

// Validate will check if all necessary fields are given
func (c *Config) Validate() error {
	if err := c.validateGlobal(); err != nil {
		return err
	}
	if err := c.validateChecks(); err != nil {
		return err
	}
	return c.validateNotifications()
}

func (c *Config) validateGlobal() error {
	if c.Global.Workers == 0 {
		return fmt.Errorf("config.validate: [global] workers not defined (cur val: %d), fix config file", c.Global.Workers)
	}

	return nil
}

func (c *Config) validateChecks() error {
	if c.isNoChecksDefined() {
		return fmt.Errorf("config.validate: no checks defined, fix config file")
	}

	for i, check := range c.Check.Web {
		i++ // count from 1
		if check.URL == "" {
			return fmt.Errorf("config.validate: empty 'url' for %q. check. This field is mandatory, fix config file", check.ID)
		}
		if err := c.checkGeneral(check.FailNotify, check.SlowNotify, check.ID, "web", i); err != nil {
			return err
		}
	}

	for i, check := range c.Check.Ping {
		i++
		if check.Host == "" {
			return fmt.Errorf("config.validate: empty 'host' for %q. check. This field is mandatory, fix config file", check.ID)
		}
		if err := c.checkGeneral(check.FailNotify, check.SlowNotify, check.ID, "ping", i); err != nil {
			return err
		}
	}
	for i, check := range c.Check.Port {
		i++
		if check.Host == "" {
			return fmt.Errorf("config.validate: empty 'host' for %q. check. This field is mandatory, fix config file", check.ID)
		}
		if check.Port == 0 {
			return fmt.Errorf("config.validate: empty 'port' for %q. check. This field is mandatory, fix config file", check.ID)
		}
		if err := c.checkGeneral(check.FailNotify, check.SlowNotify, check.ID, "port", i); err != nil {
			return err
		}
	}

	for i, check := range c.Check.Cmd {
		i++
		if check.Cmd == "" {
			return fmt.Errorf("config.validate: emtry 'cmd' for %q. check. This field is mandatory, fix config file", check.ID)
		}
		if err := c.checkGeneral(check.FailNotify, check.SlowNotify, check.ID, "cmd", i); err != nil {
			return err
		}

	}
	return nil
}

func (c *Config) checkGeneral(FailNotify, SlowNotify []string, cID, cType string, i int) error {
	if cID == "" {
		return fmt.Errorf("config.validate: empty 'id' in %d. %s check. This field is mandatory, fix config file", i, cType)
	}
	if FailNotify != nil {
		if err := c.validateNotifyIDList(FailNotify); err != nil {
			return fmt.Errorf("config.validate: wrong notification in 'fail_notify' for %q. %s check, err: %v. fix config file", cID, cType, err)
		}
	}
	if SlowNotify != nil {
		if err := c.validateNotifyIDList(SlowNotify); err != nil {
			return fmt.Errorf("config.validate: wrong notification in 'slow_notify' for %q. %s check, err: %v. fix config file", cID, cType, err)
		}
	}
	return nil
}

func (c *Config) validateNotifications() error {
	for i, notif := range c.Notify.Mail {
		i++ //count from 1
		if notif.ID == "" {
			return fmt.Errorf("config.validate: empty 'ID' for %d. notification. This field is mandatory, fix config file", i)
		}
		if notif.Server == "" {
			return fmt.Errorf("config.validate: empty 'server' for %q notification. This field is mandatory, fix config file", notif.ID)
		}
		if notif.Port == 0 {
			return fmt.Errorf("config.validate: Given 0 as 'port' for %q notification. This field must be non-zero, fix config file", notif.ID)
		}
		if notif.From == "" && notif.User == "" {
			return fmt.Errorf("config.validate: empty 'from' for %q notification. This field is mandatory, fix config file", notif.ID)
		}
		if notif.To == nil {
			return fmt.Errorf("config.validate: empty 'to' for %q notification. This field is mandatory, fix config file", notif.ID)
		}

	}

	for _, notif := range c.Notify.Jabber {

		if notif.Server == "" {
			return fmt.Errorf("config.validate: empty 'server' for %q notification. This field is mandatory, fix config file", notif.ID)
		}
		if notif.Port == 0 {
			return fmt.Errorf("config.validate: Given 0 as 'port' for %q notification. This field must be non-zero, fix config file", notif.ID)
		}
		if notif.User == "" {
			return fmt.Errorf("config.validate: empty 'from' for %q notification. This field is mandatory, fix config file", notif.ID)
		}
		if notif.To == nil {
			return fmt.Errorf("config.validate: empty 'to' for %q notification. This field is mandatory, fix config file", notif.ID)
		}
	}

	for _, notif := range c.Notify.Cmd {

		if notif.Cmd == "" {
			return fmt.Errorf("config.validate: empty 'cmd' for %q notification. This field is mandatory, fix config file", notif.ID)
		}
	}
	return nil
}

// Normalize will fill in default values if missing in config
func (c *Config) Normalize() {
	c.normalizeGlobal()
	c.normalizeChecks()
	c.normalizeNotifications()

	c.parseCheckVars()
	c.parseNotifVars()
}

func (c *Config) normalizeGlobal() {
	if c.Global.Loglevel == "" {
		c.Global.Loglevel = Loglevel
	}
	if c.Global.HTTPTimeout == 0 {
		c.Global.HTTPTimeout = HTTPTimeout
	}
	if c.Global.PingTimeout == 0 {
		c.Global.PingTimeout = PingTimeout
	}
	if c.Global.PortTimeout == 0 {
		c.Global.PortTimeout = PortTimeout
	}
}

func (c *Config) normalizeChecks() {
	notifids := c.getAllNotificationIDs()
	for i, check := range c.Check.Web {
		if check.Method == "" {
			c.Check.Web[i].Method = CheckMethod
		}
		if check.RepeatInterval == 0 {
			c.Check.Web[i].RepeatInterval = CheckRepeatInterval
		}
		if check.ExpectedCode == 0 {
			c.Check.Web[i].ExpectedCode = WebExpectedCode
		}
		if check.ExpectedTime == 0 {
			c.Check.Web[i].ExpectedTime = WebExpectedTime
		}
		if check.AllowedFails < 0 { // TODO: new meaning: 0 is correct option
			c.Check.Web[i].AllowedFails = CheckAllowedFails
		}
		if check.AllowedSlows < 0 {
			c.Check.Web[i].AllowedSlows = CheckAllowedSlows
		}
		if check.FailNotify == nil {
			c.Check.Web[i].FailNotify = notifids
		}
		if check.SlowNotify == nil {
			c.Check.Web[i].SlowNotify = notifids
		}
	}
	for i, check := range c.Check.Ping {
		if check.RepeatInterval == 0 {
			c.Check.Web[i].RepeatInterval = CheckRepeatInterval
		}
		if check.ExpectedCode == 0 {
			c.Check.Web[i].ExpectedCode = PingExpectedCode
		}
		if check.ExpectedTime == 0 {
			c.Check.Web[i].ExpectedTime = PingExpectedTime
		}
		if check.AllowedFails < 0 {
			c.Check.Web[i].AllowedFails = PingAllowedFails
		}
		if check.AllowedSlows < 0 {
			c.Check.Web[i].AllowedSlows = PingAllowedSlows
		}
		if check.FailNotify == nil {
			c.Check.Web[i].FailNotify = notifids
		}
		if check.SlowNotify == nil {
			c.Check.Web[i].SlowNotify = notifids
		}
	}
	for i, check := range c.Check.Port {
		if check.RepeatInterval == 0 {
			c.Check.Web[i].RepeatInterval = CheckRepeatInterval
		}
		if check.ExpectedCode == 0 {
			c.Check.Web[i].ExpectedCode = PortExpectedCode
		}
		if check.ExpectedTime == 0 {
			c.Check.Web[i].ExpectedTime = PortExpectedTime
		}
		if check.AllowedFails < 0 {
			c.Check.Web[i].AllowedFails = PortAllowedFails
		}
		if check.AllowedSlows < 0 {
			c.Check.Web[i].AllowedSlows = PortAllowedSlows
		}
		if check.FailNotify == nil {
			c.Check.Web[i].FailNotify = notifids
		}
		if check.SlowNotify == nil {
			c.Check.Web[i].SlowNotify = notifids
		}
	}
	for i, check := range c.Check.Cmd {
		if check.RepeatInterval == 0 {
			c.Check.Web[i].RepeatInterval = CheckRepeatInterval
		}
		if check.ExpectedCode == 0 {
			c.Check.Web[i].ExpectedCode = CmdExpectedCode
		}
		if check.ExpectedTime == 0 {
			c.Check.Web[i].ExpectedTime = CmdExpectedTime
		}
		if check.AllowedFails < 0 {
			c.Check.Web[i].AllowedFails = CheckAllowedFails
		}
		if check.AllowedSlows < 0 {
			c.Check.Web[i].AllowedSlows = CheckAllowedSlows
		}
		if check.FailNotify == nil {
			c.Check.Web[i].FailNotify = notifids
		}
		if check.SlowNotify == nil {
			c.Check.Web[i].SlowNotify = notifids
		}
	}
}

// parseCheckVars will expand all $var or ${var} to actual
// enviroment variable.
func (c *Config) parseCheckVars() {
	for i, check := range c.Check.Web {
		c.Check.Web[i].FormParams = template.ParseEnv(check.FormParams, check.ID, "params")
		c.Check.Web[i].Headers = template.ParseEnv(check.Headers, check.ID, "headers")
	}
	// TODO: add parsing for ping, port, cmd

}

func (c *Config) parseNotifVars() {
	for i, notif := range c.Notify.Mail {
		c.Notify.Mail[i].User = template.ParseEnv(notif.User, notif.ID, "user")
		c.Notify.Mail[i].Pass = template.ParseEnv(notif.Pass, notif.ID, "pass")
		c.Notify.Mail[i].Server = template.ParseEnv(notif.Server, notif.ID, "server")
	}
	// TODO: add all notifications
}

// getAllNotificationIDs returns slice of all notification IDs.
func (c *Config) getAllNotificationIDs() []string {
	ids := []string{}
	for _, mail := range c.Notify.Mail {
		ids = append(ids, mail.ID)
	}
	for _, jabber := range c.Notify.Jabber {
		ids = append(ids, jabber.ID)
	}
	for _, cmd := range c.Notify.Cmd {
		ids = append(ids, cmd.ID)
	}
	return ids
}

// validateNotifyIDList returns error if any notification on list
// is not found in defined notifications.
func (c *Config) validateNotifyIDList(ss []string) error {
	for _, nID := range ss {
		if !found(nID, c.getAllNotificationIDs()) {
			return fmt.Errorf("notification %q is not defined", nID)
		}
	}
	return nil
}

func (c *Config) normalizeNotifications() {
	for i, notif := range c.Notify.Mail {
		if notif.From == "" {
			c.Notify.Mail[i].From = notif.User
		}
		c.Notify.Mail[i].FailSubject = sets(notif.FailSubject, c.Global.NotifyFailOKSubject, NotifyFailSubject)
		c.Notify.Mail[i].FailOKSubject = sets(notif.FailOKSubject, c.Global.NotifyFailOKSubject, NotifyFailOKSubject)
		c.Notify.Mail[i].FailText = sets(notif.FailText, c.Global.NotifyFailText, NotifyFailText)
		c.Notify.Mail[i].FailOKText = sets(notif.FailOKText, c.Global.NotifyFailOKText, NotifyFailOKText)
		c.Notify.Mail[i].SlowSubject = sets(notif.SlowSubject, c.Global.NotifySlowSubject, NotifySlowSubject)
		c.Notify.Mail[i].SlowOKSubject = sets(notif.SlowOKSubject, c.Global.NotifySlowOKSubject, NotifySlowOKSubject)
		c.Notify.Mail[i].SlowText = sets(notif.SlowText, c.Global.NotifySlowText, NotifySlowText)
		c.Notify.Mail[i].SlowOKText = sets(notif.SlowOKText, c.Global.NotifySlowOKText, NotifySlowOKText)

		c.Notify.Mail[i].FailSchedule = setss(notif.FailSchedule, c.Global.FailSchedule, FailSchedule)
		c.Notify.Mail[i].SlowSchedule = setss(notif.SlowSchedule, c.Global.SlowSchedule, SlowSchedule)
	}
}

func valStr(s, n, t, id string) error {
	if s == "" {
		return fmt.Errorf("config.validate: missing or empty '%s' for %q %s check, it is mandatory, fix config", n, id, t)
	}
	return nil
}

func valInt(i int, n, t, id string) error {
	if i == 0 {
		return fmt.Errorf("config.validate: missing or zero '%s' for %q %s check, it is mandatory, fix config", n, id, t)
	}
	return nil
}

func sets(t, glob, def string) string {
	if t == "" {
		if glob != "" {
			return glob
		}
		return def
	}
	return t
}

func setss(t, glob, def []string) []string {
	if len(t) == 0 {
		if len(glob) != 0 {
			return glob
		}
		return def
	}
	return t
}

func found(s string, ss []string) bool {
	found := false
	for _, t := range ss {
		if s == t {
			found = true
		}
	}
	return found
}
