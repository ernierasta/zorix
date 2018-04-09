package shared

// Worker interface has to be fulfilled by every type of worker
type Worker interface {
	Send(c CheckConfig) (code int, reqDuration int64, err error)
}

// Notifier sends notification away
type Notifier interface {
	Send(c CheckConfig, n NotifConfig)
}
