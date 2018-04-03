package shared

import "time"

//TODO: add nice comments to every struct field

// Check type represents all check attributes
type Check struct {
	ID           int
	Type         string
	Check        string
	Params       []string
	Repeat       Duration
	ExpectedCode int      `toml:"code"`
	ExpectedTime int64    `toml:"time"`
	AllowedFails int      `toml:"fails"`
	AllowedSlows int      `toml:"slows"`
	NotifyFail   []string `toml:"notify_fail"`
	NotifySlow   []string `toml:"notify_slow"`
	ResultData
}

// ResultData contains Check additional data.
// Probably not used separetly.
type ResultData struct {
	WorkerType      Worker
	ReturnedCode    int
	Error           error
	ReturnedTime    int64
	Slowdowns       int
	Failed          int
	Timestamp       time.Time
	Failure         bool
	Slow            bool
	RecoveryFailure bool
	RecoverySlow    bool
	Debug           []string
}

// Notification type represent all notification attributes
type Notification struct {
	ID, Type      string
	Server        string
	Port          int
	User, Pass    string
	From          string
	To            []string
	SubjectFail   string `toml:"subject_fail"`
	SubjectSlow   string `toml:"subject_slow"`
	SubjectFailOK string `toml:"subject_fail_ok"`
	SubjectSlowOK string `toml:"subject_slow_ok"`
	TextFail      string `toml:"text_fail"`
	TextSlow      string `toml:"text_slow"`
	TextFailOK    string `toml:"text_fail_ok"`
	TextSlowOK    string `toml:"text_slow_ok"`
	Subject       string
	Text          string
	RepeatFail    []Duration `toml:"repeat_fail"`
	RepeatSlow    []Duration `toml:"repeat_slow"`
}

// NotifiedCheck is Check with notification ID string.
// It is used for sending notification.
type NotifiedCheck struct {
	Check
	NotificationID string
}

// Duration is custom time.Duration like type
type Duration struct {
	time.Duration
}

// UnmarshalText fulfils toml interface
func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err

}

// ParseDuration is used in code. Parses string to duration
func (d *Duration) ParseDuration(text string) error {
	var err error
	d.Duration, err = time.ParseDuration(text)
	return err

}
