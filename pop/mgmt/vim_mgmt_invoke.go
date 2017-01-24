package mgmt

import (
	"encoding/json"
	"errors"

	"github.com/mcilloni/go-openbaton/vnfm/channel"
)

type VNFMChannelAccessor func() (channel.Channel, error)

type VIMConnector interface {
	Start(id string) error
}

func NewConnector(vimname string, acc VNFMChannelAccessor) VIMConnector {
	return conn{
		acc: acc,
		id:  makeID(vimname),
	}
}

type conn struct {
	acc VNFMChannelAccessor
	id  string
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

func (c conn) exchange(req []byte) ([]byte, error) {
	cln, err := c.acc()
	if err != nil {
		return nil, err
	}

	return cln.Exchange(c.id, req)
}

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
