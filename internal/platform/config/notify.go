package config

// Notify contains all notification types
type Notify struct {
	Mail   []NotifyMail   `toml:"mail"`
	Jabber []NotifyJabber `toml:"jabber"`
	Cmd    []NotifyCmd    `toml:"cmd"`
}

// NotifyShared params common to all notifictation types
type NotifyShared struct {
	ID            string
	FailSubject   string   `toml:"fail_subject"`
	FailOKSubject string   `toml:"fail_ok_subject"`
	FailText      string   `toml:"fail_text"`
	FailOKText    string   `toml:"fail_ok_text"`
	SlowSubject   string   `toml:"slow_subject"`
	SlowOKSubject string   `toml:"slow_ok_subject"`
	SlowText      string   `toml:"slow_text"`
	SlowOKText    string   `toml:"slow_ok_text"`
	FailSchedule  []string `toml:"fail_schedule"`
	SlowSchedule  []string `toml:"slow_schedule"`
	NoRecovery    bool     `toml:"no_recovery"`
}

// NotifyMail notify struct
type NotifyMail struct {
	Server     string
	Port       int
	IgnoreCert bool `toml:"ignore_cert"`
	User, Pass string
	From       string
	To         []string
	NotifyShared
}

// NotifyJabber struct
type NotifyJabber struct {
	Server     string
	Port       int
	IgnoreCert bool `toml:"ignore_cert"`
	User       string
	Pass       string
	To         []string
	NotifyShared
}

//NotifyCmd struct
type NotifyCmd struct {
	Cmd string
	NotifyShared
}
