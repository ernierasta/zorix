package logger

import (
	"log/syslog"
	"os"

	log "github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

// Set logger configures logrus.
// dest - empty (log to stdout), "syslog" (log to syslog), any other value is treated as
func Set(dest string, level string) {

	var lvl log.Level
	var slvl syslog.Priority
	switch level {
	case "", "debug":
		lvl = log.DebugLevel
		slvl = syslog.LOG_DEBUG
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
		hook, err := logrus_syslog.NewSyslogHook("udp", "localhost:514", slvl, "")
		if err != nil {
			log.Error("unable to connect to local syslog daemon")
			log.SetOutput(os.Stdout)
		} else {
			log.AddHook(hook)
		}
	default:
		logf, err := os.OpenFile(dest, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening log file: %v", err)
		}
		log.SetOutput(logf)

	}
}
