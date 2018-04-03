package config

import (
	"fmt"

	"github.com/ernierasta/zorix/shared"

	"github.com/BurntSushi/toml"
)

const (
	CheckType         = "web"
	CheckRepeat       = "60s"
	CheckExpectedCode = 200
	CheckExpectedTime = 1000
	CheckAllowedSlows = 3

	NotifType        = "mail"
	NotifSubjectFail = "{check}{params} problem"
	NotifSubjectSlow = "{check}{params} slow"
	NotifTextFail    = "FAILURE:\n{check}{params}\nTime: {timestamp}\n\nResponse code: {responsecode}\nError: {error}\n"
	NotifTextSlow    = "SLOW RESPONSE:\n{check}{params}\nTime: {timestamp}\n\nResponse time: {responsetime}\nExpected time: {expectedtime}"
)

type Config struct {
	Workers       int
	Notifications []shared.Notification `toml:"notify"`
	Checks        []shared.Check        `toml:"check"`
	file          string
}

func New(file string) *Config {
	return &Config{file: file}
}

func (c *Config) Read() error {

	_, err := toml.DecodeFile(c.file, c)
	if err != nil {
		return fmt.Errorf("in file %s: err: %s", c.file, err)
	}
	return nil

}

// Validate will check if all necessary fields are given
func (c *Config) Validate() error {
	if err := c.validateChecks(); err != nil {
		return err
	}

	return c.validateNotifications()
}

func (c *Config) validateChecks() error {
	for i, check := range c.Checks {
		i++ //count from 1
		if check.Check == "" {
			return fmt.Errorf("config.validate: empty 'check' for %d. check. This field is mandatory, fix config file", i)
		}
	}
	return nil
}

func (c *Config) validateNotifications() error {
	for i, notif := range c.Notifications {
		i++ //count from 1
		if notif.ID == "" {
			return fmt.Errorf("config.validate: empty 'ID' for %d. notification. This field is mandatory, fix config file", i)
		}
		if notif.Server == "" {
			return fmt.Errorf("config.validate: empty 'server' for %d. notification. This field is mandatory, fix config file", i)
		}
		if notif.Port == 0 {
			return fmt.Errorf("config.validate: Given 0 as 'port' for %d. notification. This field must be non-zero, fix config file", i)
		}
		if notif.From == "" && notif.User == "" {
			return fmt.Errorf("config.validate: empty 'from' for %d. notification. This field is mandatory, fix config file", i)
		}
		if notif.To == nil {
			return fmt.Errorf("config.validate: empty 'to' for %d. notification. This field is mandatory, fix config file", i)
		}

	}
	return nil

}

// Normalize will fill in default values if missing in config
func (c *Config) Normalize() {
	c.normalizeChecks()
	c.normalizeNotifications()
}

func (c *Config) normalizeChecks() {
	notifids := c.getAllNotificationIDs()
	for i, check := range c.Checks {
		if check.Type == "" {
			c.Checks[i].Type = CheckType
		}
		if check.Repeat.Duration == 0 {
			c.Checks[i].Repeat.ParseDuration(CheckRepeat)
		}
		if check.ExpectedCode == 0 {
			c.Checks[i].ExpectedCode = CheckExpectedCode
		}
		if check.ExpectedTime == 0 {
			c.Checks[i].ExpectedTime = CheckExpectedTime
		}
		if check.AllowedSlows == 0 {
			c.Checks[i].AllowedSlows = CheckAllowedSlows
		}
		if check.NotifyFail == nil {
			c.Checks[i].NotifyFail = notifids
		}
		if check.NotifySlow == nil {
			c.Checks[i].NotifySlow = notifids
		}
	}
}

func (c *Config) getAllNotificationIDs() []string {
	ids := []string{}
	for _, notif := range c.Notifications {
		ids = append(ids, notif.ID)
	}
	return ids
}

func (c *Config) normalizeNotifications() {
	for i, notif := range c.Notifications {
		if notif.Type == "" {
			c.Notifications[i].Type = NotifType
		}
		if notif.From == "" {
			c.Notifications[i].From = notif.User
		}
		if notif.SubjectFail == "" {
			c.Notifications[i].SubjectFail = NotifSubjectFail
		}
		if notif.SubjectSlow == "" {
			c.Notifications[i].SubjectSlow = NotifSubjectSlow
		}
		if notif.TextFail == "" {
			c.Notifications[i].TextFail = NotifTextFail
		}
		if notif.TextSlow == "" {
			c.Notifications[i].TextSlow = NotifTextSlow
		}
		if notif.RepeatFail == nil {
			c.Notifications[i].RepeatFail = []shared.Duration{
				shared.Duration{Duration: 60},
				shared.Duration{Duration: 300},
				shared.Duration{Duration: 600},
			}
		}
		if notif.RepeatSlow == nil {
			c.Notifications[i].RepeatSlow = []shared.Duration{
				shared.Duration{Duration: 300},
				shared.Duration{Duration: 0},
			}
		}
	}
}
