package mgmt

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/openbaton/go-openbaton/catalogue"
	"github.com/openbaton/go-openbaton/util"
	"github.com/streadway/amqp"
)

var (
	// ErrInternal represents an internal failure.
	ErrInternal = errors.New("internal error")

	// ErrTooFewParams is returned when a request has too few parameters for the
	// requested function.
	ErrTooFewParams = errors.New("not enough parameters for function")

	// ErrMalformedParams generically notifies the caller of malformed parameters.
	ErrMalformedParams = errors.New("malformed parameters")
)

// Handler is an interface that represents the concrete functions to be remotely
// invoked. The Manager primary task is to deliver requests to the Handler and send its
// responses to the caller.
type Handler interface {
	// AddMetadata sets the given entries as metadata for the server specified by the given
	// id, merging them with existing ones. An empty value will erase the key from the metadata.
	AddMetadata(id string, entries map[string]string) error

	// Check checks if a Server with the given ID is up in the VIM.
	Check(id string) (*catalogue.Server, error)

	// Start starts an already created Server with the given ID.
	Start(id string) error
}

// doRequest dispatches a request to the Handler.
func (m *manager) doRequest(req request) response {
	var val interface{}
	var err error

	switch fn := strings.ToLower(req.Func); fn {
	case fnAddMetadata:
		err = m.handleAddMetadata(req.Params)

	case fnCheck:
		val, err = m.handleCheck(req.Params)

	case fnStart:
		err = m.handleStart(req.Params)

	default:
		err = fmt.Errorf("unknown function: %s", fn)
	}

	var valB json.RawMessage
	if val != nil {
		valB, err = json.Marshal(val)
	}

	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	return response{Value: valB, Error: errStr}
}

// handle handles a delivery, by unmarshling the request and  dispatching it to the handler.
func (m *manager) handle(cnl *amqp.Channel, delivery amqp.Delivery) {
	tag := util.FuncName()

	var req request
	if err := json.Unmarshal(delivery.Body, &req); err != nil {
		m.l.WithError(err).WithField("tag", tag).Error("error while handling delivery")
		return
	}

	respBytes, err := json.Marshal(m.doRequest(req))
	if err != nil {
		m.l.WithError(err).WithField("tag", tag).Error("error while handling delivery")
		return
	}

	err = cnl.Publish(
		"",
		delivery.ReplyTo,
		false,
		false,
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: delivery.CorrelationId,
			Body:          respBytes,
		},
	)

	if err != nil {
		m.l.WithError(err).WithField("tag", tag).Error("error while publishing response")
		return
	}

	err = delivery.Ack(false)
	if err != nil {
		m.l.WithError(err).WithField("tag", tag).Error("error while acknowledging delivery")
		return
	}
}

func (m *manager) handleAddMetadata(params json.RawMessage) error {
	var ap addMetadataParams
	if err := json.Unmarshal(params, &ap); err != nil {
		return ErrMalformedParams
	}

	return m.handl.AddMetadata(ap.ID, ap.Entries)
}

func (m *manager) handleCheck(params json.RawMessage) (*catalogue.Server, error) {
	var id checkParams
	if err := json.Unmarshal(params, &id); err != nil {
		return nil, ErrMalformedParams
	}

	return m.handl.Check(string(id))
}

func (m *manager) handleStart(params json.RawMessage) error {
	var id startParams
	if err := json.Unmarshal(params, &id); err != nil {
		return ErrMalformedParams
	}

	return m.handl.Start(string(id))
}
