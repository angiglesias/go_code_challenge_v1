// Package logging provides simple logging global package implementation
package logging

import "log"

type Level int

const (
	ERROR Level = iota
	WARN
	INFO
	DEBUG
	TRACE
)

var (
	Verbosity Level = INFO
)

func Setup(level Level) {
	Verbosity = level
}

func Errorf(format string, args ...interface{}) {
	Logf(ERROR, format, args...)
}

func Warnf(format string, args ...interface{}) {
	Logf(WARN, format, args...)
}

func Infof(format string, args ...interface{}) {
	Logf(INFO, format, args...)
}

func Debugf(format string, args ...interface{}) {
	Logf(DEBUG, format, args...)
}

func Tracef(format string, args ...interface{}) {
	Logf(TRACE, format, args...)
}

func Logf(level Level, format string, args ...interface{}) {
	if level <= Verbosity {
		log.Printf(getLevelPrefix(level)+format+"\n", args...)
	}
}

func getLevelPrefix(level Level) string {
	switch level {
	case ERROR:
		return "[ERROR] "
	case WARN:
		return "[WARN] "
	case INFO:
		return "[INFO] "
	case DEBUG:
		return "[DEBUG] "
	case TRACE:
		return "[TRACE] "
	default:
		return "[INFO] "
	}
}
