package config

// Global settings
type Global struct {
	Workers             int
	Loglevel            string
	HTTPTimeout         int      `toml:"http_timeout"`
	PingTimeout         int      `toml:"ping_timeout"`
	PortTimeout         int      `toml:"port_timeout"`
	NotifyFailSubject   string   `toml:"notify_fail_subject"`
	NotifySlowSubject   string   `toml:"notify_slow_subject"`
	NotifyFailOKSubject string   `toml:"notify_fail_ok_subject"`
	NotifySlowOKSubject string   `toml:"notify_slow_ok_subject"`
	NotifyFailText      string   `toml:"notify_fail_text"`
	NotifySlowText      string   `toml:"notify_slow_text"`
	NotifyFailOKText    string   `toml:"notify_fail_ok_text"`
	NotifySlowOKText    string   `toml:"notify_slow_ok_text"`
	FailSchedule        []string `toml:"fail_schedule"`
	SlowSchedule        []string `toml:"slow_schedule"`
}
