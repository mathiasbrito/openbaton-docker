package server

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// newLogger inits a logger
func newLogger(cfg Config) *log.Logger {
	color := cfg.LogLevel >= log.DebugLevel // enable forced color for debug

	l := log.New()

	l.Formatter = &log.TextFormatter{
		DisableColors:    !color,
		DisableTimestamp: !cfg.LogTimestamps,
		ForceColors:      color,
		FullTimestamp:    cfg.LogTimestamps,
		TimestampFormat:  time.RFC3339Nano,
	}

	l.Level = cfg.LogLevel

	return l
}
