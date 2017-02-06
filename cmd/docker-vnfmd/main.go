package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/openbaton/go-openbaton/util"
	"github.com/openbaton/go-openbaton/vnfm"
	_ "github.com/openbaton/go-openbaton/vnfm/amqp" // import needed to load the driver
	"github.com/openbaton/go-openbaton/vnfm/config"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var confPath = flag.String("cfg", "config.toml", "a TOML file to be loaded as config")

func main() {
	tag := util.FuncName()

	flag.Parse()

	cfg, err := config.LoadFile(*confPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: while loading config file %s: %v\n", *confPath, err)
		os.Exit(1)
	}

	h := &handl{}

	svc, err := vnfm.New("amqp", h, cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: while initialising VNFM %s: %v\n", *confPath, err)
		os.Exit(1)
	}

	l := svc.Logger()

	h.Logger = l
	h.acc = func() (*amqp.Channel, error) {
		vcnl, err := svc.ChannelAccessor()()
		if err != nil {
			return nil, err
		}

		icnl, err := vcnl.Impl()
		if err != nil {
			return nil, err
		}

		return icnl.(*amqp.Channel), nil
	}

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
		}).Fatal("VNFM failed during execution")
	}

	<-join

	l.WithFields(log.Fields{
		"tag": tag,
	}).Info("exiting cleanly")
}
