package mgmt

import (
	"testing"
	"time"

	"errors"
	"github.com/mcilloni/go-openbaton/catalogue/messages"
	"github.com/mcilloni/go-openbaton/util"
	"github.com/mcilloni/go-openbaton/vnfm/channel"
	"github.com/streadway/amqp"
    log "github.com/sirupsen/logrus"
)

func init() {
    log.SetLevel(log.DebugLevel)
}

func dialAMQP() (*amqp.Channel, error) {
	conn, err := amqp.Dial(util.AmqpUriBuilder("admin", "openbaton", "localhost", "", 5672, false))
	if err != nil {
		return nil, err
	}

	return conn.Channel()
}

func chanChan() (channel.Channel, error) {
	cnl, err := dialAMQP()
	if err != nil {
		return nil, err
	}

	return testChan{cnl}, nil
}

type testChan struct {
	*amqp.Channel
}

func (tc testChan) Close() error {
	return nil
}

func (tc testChan) temporaryQueue() (string, error) {
	queue, err := tc.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)

	if err != nil {
		return "", err
	}

	return queue.Name, nil
}

func (tc testChan) Exchange(dest string, msg []byte) ([]byte, error) {
	replyQueue, err := tc.temporaryQueue()
	if err != nil {
		return nil, err
	}

	deliveries, err := tc.Consume(
		replyQueue, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return nil, err
	}

	corrID := util.GenerateID()

	err = tc.Publish(
		MgmtExchange,
		dest,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrID,
			ReplyTo:       replyQueue,
			Body:          msg,
		},
	)
	if err != nil {
		return nil, err
	}

	timeout := time.After(10 * time.Second)

DeliveryLoop:
	for {
		select {
		case <-timeout:
			break DeliveryLoop

		case delivery, ok := <-deliveries:
			if !ok {
				break DeliveryLoop
			}

			if delivery.CorrelationId == corrID {
				return delivery.Body, nil
			}
		}
	}

	return nil, errors.New("no reply received")
}

// NFVOExchange sends a message to the NFVO, and then waits for a reply.
// The outgoing message must have From() == messages.VNFR.
func (tc testChan) NFVOExchange(msg messages.NFVMessage) (messages.NFVMessage, error) {
	return nil, nil
}

// NFVOSend sends a message to the NFVO without waiting for a reply.
// A success while sending the message is no guarantee about the NFVO actually receiving it.
func (tc testChan) NFVOSend(msg messages.NFVMessage) error {
	return nil
}

// NotifyReceiver creates a channel on which received messages will be delivered.
// The returned channel will be removed if nobody is listening on it for a while.
func (tc testChan) NotifyReceived() (<-chan messages.NFVMessage, error) {
	return nil, nil
}

// Send sends a message to an implementation defined destination without waiting for a reply.
// A success while sending the message is no guarantee about the destination actually receiving it.
func (tc testChan) Send(dest string, msg []byte) error {
	return nil
}

// Status returns the current status of the Channel.
func (tc testChan) Status() channel.Status {
	return channel.Running
}

type handler chan string

func (h handler) Start(id string) error {
	h <- id

	return nil
}

var testID = "test"

func TestAll(t *testing.T) {
	r := make(chan string, 1)

	m := NewManager(testID, handler(r), dialAMQP, nil)
	c := NewConnector(testID, chanChan)

    time.Sleep(time.Second)

	if err := c.Start("testid"); err != nil {
		t.Fatal(err)
	}

	id := <-r

	t.Logf("recv id: %s", id)

	m.Stop()
}
