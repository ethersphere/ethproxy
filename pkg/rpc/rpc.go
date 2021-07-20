package rpc

import (
	"errors"
	"sync"

	"github.com/ethersphere/ethproxy/pkg/callback"
	"github.com/ethersphere/ethproxy/pkg/ethrpc"
)

const (
	BlockNumberFreeze = "blockNumberFreeze"
	BlockNumberRecord = "blockNumberRecord"
)

type State struct {
	mtx               sync.Mutex
	BlockNumber       uint64
	FreezeBlockNumber bool
}

type Caller struct {
	call  *callback.Callback
	state State
}

func New(call *callback.Callback) *Caller {
	return &Caller{
		call: call,
	}
}

func (c *Caller) GetState() State {
	return c.state
}

func (c *Caller) Register(method string, params ...interface{}) error {
	switch method {

	case BlockNumberRecord:
		c.call.On(ethrpc.BlockNumber, func(resp *callback.Response) {
			bN, err := resp.Body.BlockNumber()
			if err != nil {
				return
			}
			if !c.frozenBlockNumber() {
				c.state.BlockNumber = bN
			}
		})
	case BlockNumberFreeze:
		if len(params) == 0 {
			c.freezeBlockNumber()
			c.call.On(ethrpc.BlockNumber, func(resp *callback.Response) {
				resp.Body.SetBlockNumber(c.state.BlockNumber)
			})
		} else {
			for _, param := range params {
				ip, ok := param.(string)
				if !ok {
					return errors.New("bad param")
				}
				func(ip string) {
					c.call.On(ethrpc.BlockNumber, func(resp *callback.Response) {
						if resp.IP == ip {
							resp.Body.SetBlockNumber(c.state.BlockNumber)
						}
					})
				}(ip)
			}
			c.freezeBlockNumber()
		}
	default:
		return errors.New("bad method")
	}

	return nil
}

func (c *Caller) freezeBlockNumber() {
	c.state.mtx.Lock()
	defer c.state.mtx.Unlock()
	c.state.FreezeBlockNumber = true
}

func (c *Caller) frozenBlockNumber() bool {
	c.state.mtx.Lock()
	defer c.state.mtx.Unlock()
	return c.state.FreezeBlockNumber
}
