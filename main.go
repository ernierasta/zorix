package main

import (
	"log"
	"time"

	"github.com/bcicen/grmon/agent"
	"github.com/ernierasta/zorix/check"
	"github.com/ernierasta/zorix/config"
	"github.com/ernierasta/zorix/notify"
	"github.com/ernierasta/zorix/processor"
	"github.com/ernierasta/zorix/shared"
)

func main() {
	log.Println("Starting ...")

	grmon.Start()
	log.Println("Started grmon")

	c := config.New("config.toml")
	err := c.Read()
	if err != nil {
		log.Fatal(err)
	}
	err = c.Validate()
	if err != nil {
		log.Fatal(err)
	}
	c.Normalize()

	//all results goes there
	resultsChan := make(chan shared.Check, len(c.Checks))
	notifChan := make(chan shared.NotifiedCheck, len(c.Checks))

	chm := check.NewManager(c.Checks, c.Workers, resultsChan)
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
