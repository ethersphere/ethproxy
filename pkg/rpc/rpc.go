// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"errors"

	"github.com/ethersphere/ethproxy/pkg/callback"
	"github.com/ethersphere/ethproxy/pkg/ethrpc"
)

const (
	BlockNumberFreeze = "blockNumberFreeze"
	BlockNumberRecord = "blockNumberRecord"
)

type State struct {
	BlockNumber       uint64
	FrozenBlockNumber bool
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

func (c *Caller) Execute(method string, params ...interface{}) error {
	switch method {

	case BlockNumberRecord:

		c.call.On(ethrpc.BlockNumber, func(resp *callback.Response) {
			bN, err := resp.Body.BlockNumber()
			if err != nil {
				return
			}

			if !c.state.FrozenBlockNumber {
				c.state.BlockNumber = bN
			}
		})

	case BlockNumberFreeze:

		if len(params) == 0 {
			c.state.FrozenBlockNumber = true
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
			c.state.FrozenBlockNumber = true
		}

	default:
		return errors.New("bad method")
	}

	return nil
}

func (c *Caller) GetState() State {
	return c.state
}
