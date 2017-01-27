package mgmt

import (
	"time"

	"github.com/mcilloni/go-openbaton/util"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type AMQPChannelAccessor func() (*amqp.Channel, error)

type VIMManager interface {
	Stop()
}

func NewManager(
	vimname string,
	h Handler,
	accessor AMQPChannelAccessor,
	l *log.Logger,
) VIMManager {

	if l == nil {
		l = log.StandardLogger()
	}

	m := &manager{
		accessor: accessor,
		l:        l,
		handl:    h,
		id:       makeID(vimname),
		quitChan: make(chan struct{}),
	}

	go m.serve()

	return m
}

type manager struct {
	accessor AMQPChannelAccessor
	l        *log.Logger
	handl    Handler
	id       string
	quitChan chan struct{}
}

func (m *manager) Stop() {
	m.quitChan <- struct{}{}

	<-m.quitChan
}

func (m *manager) serve() {
	tag := util.FuncName()

ServeLoop:
	for {
		cnl, deliveries, err := m.setup()
		if err != nil {
			m.l.WithError(err).WithField("tag", tag).Error("error during delivery")

			time.Sleep(5 * time.Second)
			continue ServeLoop
		}

		for {
			select {
			case <-m.quitChan:
				close(m.quitChan)
				return

			case delivery, ok := <-deliveries:
				if !ok {
					continue ServeLoop
				}

				m.l.WithFields(log.Fields{
					"tag":     tag,
					"corr-id": delivery.CorrelationId,
				}).Debug("new delivery")

				go m.handle(cnl, delivery)
			}
		}
	}
}

func (m *manager) setup() (cnl *amqp.Channel, deliveries <-chan amqp.Delivery, err error) {
	tag := util.FuncName()

	m.l.WithField("tag", tag).Debug("setting up")

	cnl, err = m.accessor()
	if err != nil {
		return
	}

	_, err = cnl.QueueDeclare(
		m.id,  // name
		false, // durable
		true,  // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return
	}

	if err = cnl.QueueBind(m.id, m.id, MgmtExchange, false, nil); err != nil {
		return
	}

	if err = cnl.Qos(1, 0, false); err != nil {
		return
	}

	deliveries, err = cnl.Consume(
		m.id,  // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	m.l.WithField("tag", tag).Debug("all set up")

	return
}
