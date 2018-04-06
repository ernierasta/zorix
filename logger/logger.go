package logger

import (
	"io/ioutil"
	"log/syslog"
	"os"

	log "github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

// Set logger configures logrus.
// dest - empty (log to stdout), "syslog" (log to syslog), any other value is treated as filename to log to.
// For syslog we set LOG_INFO level for "debug" level, becouse some syslog setting would ignore LOG_DEBUG entries.
func Set(dest string, level string) {

	var lvl log.Level
	var slvl syslog.Priority
	switch level {
	case "", "debug":
		lvl = log.DebugLevel
		slvl = syslog.LOG_INFO
	case "info":
		lvl = log.InfoLevel
		slvl = syslog.LOG_INFO
	case "warn":
		lvl = log.WarnLevel
		slvl = syslog.LOG_WARNING
	case "error":
		lvl = log.ErrorLevel
		slvl = syslog.LOG_ERR
	case "fatal":
		lvl = log.FatalLevel
		slvl = syslog.LOG_CRIT
	case "panic":
		lvl = log.PanicLevel
		slvl = syslog.LOG_CRIT
	default:
		log.Fatal("unknown loglevel, check config file")
	}

	log.SetLevel(lvl)

	switch dest {
	case "":
		log.SetOutput(os.Stdout)
	case "syslog":
		hook, err := logrus_syslog.NewSyslogHook("", "", slvl, "zorix")
		if err != nil {
			log.SetOutput(os.Stdout) // use stdout as failback
			log.Println("unable to connect to local syslog daemon, logging to stdout")
		} else {
			log.AddHook(hook)
			log.SetFormatter(&log.TextFormatter{DisableColors: true, DisableTimestamp: true})
			log.SetOutput(ioutil.Discard) // do not output to stdout
		}
	default:
		logf, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening log file: %v", err)
		}
		log.SetOutput(logf)

	}
}
