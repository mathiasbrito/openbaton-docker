package mgmt

import (
	"encoding/json"
	"fmt"

	vnfmAMQP "github.com/mcilloni/go-openbaton/vnfm/amqp"
)

type (
	checkParams string
	startParams string
)

var (
	fnCheck = "check"
	fnStart = "start"
)

const (
	DefaultTimeout = vnfmAMQP.DefaultTimeout
	MgmtExchange   = vnfmAMQP.ExchangeDefault
)

type request struct {
	Func   string
	Params json.RawMessage
}

type response struct {
	Value json.RawMessage `json:",omitempty"`
	Error string          `json:",omitempty"`
}

func makeID(vimname string) string {
	return fmt.Sprintf("vim-mgmt-%s", vimname)
}
