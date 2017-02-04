package mgmt

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/mcilloni/go-openbaton/catalogue"
	"github.com/mcilloni/go-openbaton/util"
	"github.com/mcilloni/go-openbaton/vnfm/channel"
	"github.com/streadway/amqp"
)

// VNFMChannelAccessor is a function type that allows to access a channel.Channel instance.
type VNFMChannelAccessor func() (channel.Channel, error)

// VIMConnector is a client for a remote VIM instance.
// Its methods mirror those of Handler; invoking one of them
// actually invokes the method on the remote handler.
type VIMConnector interface {
	AddMetadata(id string, entries map[string]string) error
	Check(id string) (*catalogue.Server, error)
	Start(id string) error
}

// NewConnector creates a new Connector to the Manager for the given VIM instance.
func NewConnector(vimname string, acc VNFMChannelAccessor) VIMConnector {
	return conn{
		acc: acc,
		id:  makeID(vimname),
	}
}

// concrete conn type.
type conn struct {
	acc VNFMChannelAccessor
	id  string
}

func (c conn) AddMetadata(id string, entries map[string]string) error {
	resp, err := c.request(fnAddMetadata, addMetadataParams{
		ID: id,
		Entries: entries,
	})
	
	if err != nil {
		return err
	}

	if resp.Error != "" {
		return errors.New(resp.Error)
	}

	return nil
}

func (c conn) Check(id string) (*catalogue.Server, error) {
	resp, err := c.request(fnCheck, checkParams(id))
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}

	var srv *catalogue.Server
	if err := json.Unmarshal(resp.Value, &srv); err != nil {
		return nil, err
	}

	return srv, nil
}

func (c conn) Start(id string) error {
	resp, err := c.request(fnStart, startParams(id))
	if err != nil {
		return err
	}

	if resp.Error != "" {
		return errors.New(resp.Error)
	}

	return nil
}

// exchange does an RPC call to the Manager.
func (c conn) exchange(req []byte) ([]byte, error) {
	cln, err := c.acc()
	if err != nil {
		return nil, err
	}

	cimpl, err := cln.Impl()
	if err != nil {
		return nil, err
	}

	acnl, ok := cimpl.(*amqp.Channel)
	if !ok {
		return nil, errors.New("invalid channel - only AMQP is supported")
	}

	// check if the wanted queue exists.
	if _, err := acnl.QueueInspect(c.id); err != nil {
		return nil, err
	}

	queue, err := acnl.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)

	if err != nil {
		return nil, err
	}

	deliveries, err := acnl.Consume(
		queue.Name, // queue
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

	err = acnl.Publish(
		MgmtExchange,
		c.id,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrID,
			ReplyTo:       queue.Name,
			Body:          req,
		},
	)
	if err != nil {
		return nil, err
	}

	timeout := time.After(DefaultTimeout)

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

// request marshals the request, does the RPC call and unmarshals the response.
func (c conn) request(fn string, params interface{}) (response, error) {
	sparams, err := json.Marshal(params)
	if err != nil {
		return response{}, err
	}

	sreq, err := json.Marshal(request{Func: fn, Params: sparams})
	if err != nil {
		return response{}, err
	}

	sresp, err := c.exchange(sreq)
	if err != nil {
		return response{}, err
	}

	var resp response
	if err := json.Unmarshal(sresp, &resp); err != nil {
		return response{}, err
	}

	return resp, nil
}
