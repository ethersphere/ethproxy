// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc

import (
	"errors"

	"github.com/ethersphere/bee/pkg/logging"
	"github.com/ethersphere/ethproxy/pkg/callback"
)

const (
	BlockNumberRecord = "blockNumberRecord"
	BlockNumberFreeze = "blockNumberFreeze"
)

type State struct {
	BlockNumber uint64
}

type Caller struct {
	call   *callback.Callback
	state  State
	logger logging.Logger
}

func New(call *callback.Callback, logger logging.Logger) *Caller {
	return &Caller{
		call:   call,
		logger: logger,
	}
}

func (c *Caller) Execute(method string, params ...interface{}) (int, error) {
	switch method {
	case BlockNumberRecord:
		return c.blockNumberRecord()
	case BlockNumberFreeze:
		return c.blockNumberFreeze(c.state.BlockNumber, params)
	default:
		return 0, errors.New("bad method")
	}
}

func (c *Caller) GetState() State {
	return c.state
}
