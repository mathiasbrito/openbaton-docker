package mgmt

import (
	"fmt"
	"testing"
	"time"

	"github.com/openbaton/go-openbaton/catalogue"
	"github.com/openbaton/go-openbaton/util"
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
	c := NewConnector(testID, dialAMQP)

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

func TestCachedAccessor(t *testing.T) {
	acc := cachingAccessor(dialAMQP)

	cnl1, err := acc()
	if err != nil {
		t.Fatal(err)
	}

	cnl2, err := acc()
	if err != nil {
		t.Fatal(err)
	}

	if cnl1 != cnl2 {
		t.Fatal("allocated 2 channels, 1 expected")
	}

	t.Log("both channels are the same one")
}
