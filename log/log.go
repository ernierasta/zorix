package log

import (
	"log"
	"os"
)

var (
	lever  = "debug"
	levels = []string{"debug", "warning", "error", "silent"}
	debug  = false
	warn   = false
	errors = false
	dlog   = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	wlog   = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	elog   = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

// Set will set log level and color on/off.
// You can call this function any time, to reset logging options.
func Set(level string, color bool) {
	setColor(color)
	found := false
	for _, l := range levels {
		if l == level {
			found = true
		}
	}
	if !found {
		if level == "" {
			setLevel("debug")
		} else {
			Warnf("Wrong error level! Given: %s, availabile: %v. Setting 'debug'.", level, levels)
			setLevel("debug")
		}
	} else {
		setLevel(level)
	}
}

func setColor(c bool) {
	if c {
		dlog = log.New(os.Stdout, "\x1B[36mDEBUG: \x1B[0m", log.Ldate|log.Ltime|log.Lshortfile)
		wlog = log.New(os.Stdout, "\x1B[35mWARN: \x1B[0m", log.Ldate|log.Ltime|log.Lshortfile)
		elog = log.New(os.Stderr, "\x1B[31mERROR: \x1B[0m", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		dlog = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
		wlog = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
		elog = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

func setLevel(level string) {
	switch level {
	case "debug":
		debug, warn, errors = true, true, true
	case "warn":
		debug, warn, errors = false, true, true
	case "error":
		debug, warn, errors = false, false, true
	case "silent":
		debug, warn, errors = false, false, false

	}
}

// Debug logging
func Debug(a ...interface{}) {
	if debug {
		dlog.Println(a...)
	}
}

// Debugf formated logging
func Debugf(format string, a ...interface{}) {
	if debug {
		dlog.Printf(format, a...)
	}
}

// Warn logging
func Warn(a ...interface{}) {
	if warn {
		wlog.Println(a...)
	}
}

// Warnf formated logging
func Warnf(format string, a ...interface{}) {
	if warn {
		wlog.Printf(format, a...)
	}
}

// Error logging
func Error(a ...interface{}) {
	if errors {
		elog.Println(a...)
	}
}

// Errorf formated logging
func Errorf(format string, a ...interface{}) {
	if errors {
		elog.Printf(format, a...)
	}
}

// Fatal log and exit
func Fatal(a ...interface{}) {
	elog.Fatal(a...)
}

// Fatalf formated log and exit
func Fatalf(format string, a ...interface{}) {
	elog.Fatalf(format, a...)
}
