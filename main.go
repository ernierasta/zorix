package main

import (
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/ernierasta/zorix/check"
	"github.com/ernierasta/zorix/config"
	"github.com/ernierasta/zorix/log"
	"github.com/ernierasta/zorix/notify"
	"github.com/ernierasta/zorix/processor"
	"github.com/ernierasta/zorix/shared"
)

func main() {

	go func() { log.Debug(http.ListenAndServe("localhost:6060", nil)) }()

	c := config.New("config.toml")
	err := c.Read()
	if err != nil {
		log.Fatal(err)
	}

	log.Set(c.Global.Loglevel, true)
	log.Debug("Config loaded. Starting ...")

	err = c.Validate()
	if err != nil {
		log.Fatal(err)
	}
	c.Normalize()

	//all results goes there
	resultsChan := make(chan shared.Check, len(c.Checks)*10)
	notifChan := make(chan shared.NotifiedCheck, len(c.Checks)*10)

	chm := check.NewManager(c.Checks, c.Global.Workers, resultsChan)
	proc := processor.New(resultsChan, notifChan, len(c.Checks), c.Notifications)
	nm := notify.NewManager(notifChan, c.Notifications)

	// start listening for results
	proc.Listen()
	// start listening for notifications
	nm.Listen()
	// run checks
	chm.Register()
	chm.Run()

	for {
		time.Sleep(1 * time.Second)
	}
}
