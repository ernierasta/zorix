package shared

// Worker interface has to be fulfilled by every type of worker
type Worker interface {
	//Check(id string, input, output chan Check)
	Send(c Check) (code int, reqDuration int64, err error)
}

type Notifier interface {
	Send(c Check, n NotifConfig)
}
