package mgmt

import (
	"encoding/json"
	"fmt"

	vnfmAMQP "github.com/mcilloni/go-openbaton/vnfm/amqp"
)

type (
	startParams string
)

var (
	fnStart = "start"
)

const (
	MgmtExchange = vnfmAMQP.ExchangeDefault
)

type request struct {
	Func   string
	Params json.RawMessage
}

type response struct {
	Value interface{} `json:",omitempty"`
	Error string      `json:",omitempty"`
}

func makeID(vimname string) string {
	return fmt.Sprintf("vim-mgmt-%s", vimname)
}
