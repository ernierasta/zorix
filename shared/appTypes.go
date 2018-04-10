package shared

// CMConfig is used for configuring check.Manager.
type CMConfig struct {
	Checks      []CheckConfig
	Workers     int
	ResultsChan chan CheckConfig
	HTTPTimeout Duration
	PingTimeout Duration
	PortTimeout Duration
}
