package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"

	"github.com/mcilloni/go-openbaton/plugin"
	log "github.com/sirupsen/logrus"
	"github.com/mcilloni/go-openbaton/util"
	"github.com/mcilloni/openbaton-docker/pop/client"
)

var logPath = flag.String("log", "", "path to the optional logfile")

var defaultParams = &plugin.Params{
	BrokerAddress: "localhost",
	LogFile:  "-",
	LogLevel: log.DebugLevel,
	Name:     "openbaton",
	Port:     5672,
	Workers:  10,
	Type:     "docker-pop",
	Username: "admin",
	Password: "openbaton",
}

func main() {
	tag := util.FuncName()
	
	flag.Parse()

	args := flag.Args()

	params := defaultParams

	if len(args) == 6 {
		port, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: port must be int")
			os.Exit(1)
		}

		workers, err := strconv.Atoi(args[3])
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: workers number must be int")
			os.Exit(1)
		}

		params = &plugin.Params{
			Name:     args[0],
			BrokerAddress: args[1],
			Port:     port,
			Workers:  workers,
			Username: args[4],
			Password: args[5],
			LogLevel: log.InfoLevel,
			LogFile: *logPath,
			Type: "test",
		}
	}

	var d driver

	if params == defaultParams {
		fmt.Fprintln(os.Stderr, "warn: using default parameters")
	}

	svc, err := plugin.New(&d, params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: while initialising Plugin %s: %v\n", params.Name, err)
		os.Exit(1)
	}

	l := svc.Logger()

	d.Logger = l

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	join := make(chan struct{})

	go func() {
		tag := util.FuncName()

		<-sigChan

		l.WithFields(log.Fields{
			"tag": tag,
		}).Warn("interrupt signal received, quitting")

		if err := svc.Stop(); err != nil {
			l.WithError(err).WithFields(log.Fields{
				"tag": tag,
			}).Fatal("stopping service failed")
		}

		close(join)
	}()

	if err = svc.Serve(); err != nil {
		l.WithError(err).WithFields(log.Fields{
			"tag": tag,
		}).Fatal("Plugin failed during execution")
	}

	<-join

	client.FlushSessions()
}
