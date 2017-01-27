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
	// DefaultTimeout is the default timeout for RPC requests.
	DefaultTimeout = vnfmAMQP.DefaultTimeout

	// MgmtExchange is the default exchange to be used.
	MgmtExchange   = vnfmAMQP.ExchangeDefault
)

// Due to a bug with json.RawMessage, Go versions before 1.8 do not show the correct behaviour 
// while serializing/deserializing this structure.
// Please use Go version 1.8 or higher.

type request struct {
	Func   string
	Params json.RawMessage
}

type response struct {
	Value json.RawMessage `json:",omitempty"`
	Error string          `json:",omitempty"`
}

// makeID returns the ID to be used as the queue name for a given VIM instance.
func makeID(vimname string) string {
	return fmt.Sprintf("vim-mgmt-%s", vimname)
}
