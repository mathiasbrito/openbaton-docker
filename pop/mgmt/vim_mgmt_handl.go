package mgmt

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/streadway/amqp"
	"github.com/mcilloni/go-openbaton/util"
)

var (
	ErrTooFewParams    = errors.New("not enough parameters for function")
	ErrMalformedParams = errors.New("malformed parameters")
)

type Handler interface {
	Start(id string) error
}

func (m *manager) handle(cnl *amqp.Channel, delivery amqp.Delivery) {
	tag := util.FuncName()

	var req request
	if err := json.Unmarshal(delivery.Body, &req); err != nil {
		m.l.WithError(err).WithField("tag", tag).Error("error while handling delivery")
		return
	}

	var resp response

	switch strings.ToLower(req.Func) {
	case fnStart:
		if err := m.handleStart(req.Params); err != nil {
			resp.Error = err.Error()
		}		
	}

	respBytes, err := json.Marshal(resp)
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

func (m *manager) handleStart(params json.RawMessage) error {
	var id startParams
	if err := json.Unmarshal(params, &id); err != nil {
		return ErrMalformedParams
	}

	return m.handl.Start(string(id))
}
