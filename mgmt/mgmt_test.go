package mgmt

import (
	"fmt"
	"testing"
	"time"

	"github.com/openbaton/go-openbaton/catalogue"
	"github.com/openbaton/go-openbaton/catalogue/messages"
	"github.com/openbaton/go-openbaton/util"
	"github.com/openbaton/go-openbaton/vnfm/channel"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
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
	return nil, nil
}

func (tc testChan) Impl() (interface{}, error) {
	return tc.Channel, nil
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

func (h handler) AddMetadata(id string, entries map[string]string) error {
	h <- fmt.Sprintf("id: %#v", entries)

	return nil
}

func (h handler) Check(id string) (*catalogue.Server, error) {
	return &catalogue.Server{ID: id}, nil
}

func (h handler) Start(id string) error {
	h <- id

	return nil
}

var testID = "4de36375-7514-4c1f-8f5c-e56de8c08dcf"

func TestAll(t *testing.T) {
	r := make(chan string, 1)

	m := NewManager(testID, handler(r), dialAMQP, nil)
	c := NewConnector(testID, chanChan)

	time.Sleep(time.Second)

	sentID := "33"
	srv, err := c.Check(sentID)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("recv %s", srv.ID)

	if srv.ID != sentID {
		t.Errorf("expecting %s, got %s", sentID, srv.ID)
	}

	if err := c.Start(testID); err != nil {
		t.Fatal(err)
	}

	t.Logf("recv id: %s", <-r)

	if err := c.AddMetadata(testID, map[string]string{
		"key":  "value",
		"key2": "value2",
	}); err != nil {
		t.Fatal(err)
	}

	t.Log(<-r)

	m.Stop()
}
