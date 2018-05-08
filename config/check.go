package config

// Check contains all check types
type Check struct {
	Web  []CheckWeb  `toml:"web"`
	Ping []CheckPing `toml:"ping"`
	Port []CheckPort `toml:"port"`
	Cmd  []CheckCmd  `cmd:"cmd"`
}

// CheckShared keeps common check params
type CheckShared struct {
	ID             string
	ExpectedCode   int      `toml:"expected_code"`
	ExpectedTime   int      `toml:"timeout"`
	RepeatInterval int      `toml:"repeat_interval"`
	AllowedFails   int      `toml:"allowed_fails"`
	AllowedSlows   int      `toml:"allowed_slows"`
	FailNotify     []string `toml:"fail_notify"`
	SlowNotify     []string `toml:"slow_notify"`
}

//CheckWeb check struct
type CheckWeb struct {
	URL        string
	IgnoreCert bool   `toml:"ignore_cert"`
	FormParams string `toml:"form_params"`
	Headers    string
	Method     string
	LookFor    string `toml:"look_for"`
	Redirects  int
	CheckShared
}

//CheckPort check struct
type CheckPort struct {
	Host string
	Port int
	CheckShared
}

//CheckPing check struct
type CheckPing struct {
	Host string
	CheckShared
}

//CheckCmd check struct
type CheckCmd struct {
	Cmd string
	CheckShared
}
