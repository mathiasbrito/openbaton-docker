package server

import (
	log "github.com/sirupsen/logrus"
)

// newLogger inits a logger
func newLogger(level log.Level) *log.Logger {
	color := level >= log.DebugLevel // enable forced color for debug

	l := log.New()

	l.Formatter = &log.TextFormatter{
		DisableColors: !color,
		ForceColors:   color,
	}

	l.Level = level

	return l
}
