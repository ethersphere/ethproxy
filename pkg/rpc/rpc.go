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
	BlockNumberRecord = "blockNumberRecord"
	BlockNumberFreeze = "blockNumberFreeze"
)

type State struct {
	BlockNumber uint64
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

func (c *Caller) Execute(method string, params ...interface{}) (int, error) {
	switch method {

	case BlockNumberRecord:

		return c.call.On(ethrpc.BlockNumber, func(resp *callback.Response) {
			bN, err := resp.Body.BlockNumber()
			if err != nil {
				return
			}

			c.state.BlockNumber = bN
		}), nil

	case BlockNumberFreeze:

		if len(params) == 0 {
			return c.call.On(ethrpc.BlockNumber, func(resp *callback.Response) {
				resp.Body.SetBlockNumber(c.state.BlockNumber)
			}), nil
		} else {

			ips, err := stringArray(params)
			if err != nil {
				return 0, err
			}

			return func(ips []string) int {
				return c.call.On(ethrpc.BlockNumber, func(resp *callback.Response) {
					for _, ip := range ips {
						if resp.IP == ip {
							resp.Body.SetBlockNumber(c.state.BlockNumber)
						}
					}
				})
			}(ips), nil
		}

	default:
		return 0, errors.New("bad method")
	}
}

func (c *Caller) GetState() State {
	return c.state
}

func stringArray(args []interface{}) ([]string, error) {

	ret := make([]string, len(args))

	for _, arg := range args {
		str, ok := arg.(string)
		if !ok {
			return nil, errors.New("bad param")
		}
		ret = append(ret, str)
	}

	return ret, nil
}
