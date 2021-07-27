// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package callback

import (
	"errors"
	"sync"

	"github.com/ethersphere/bee/pkg/logging"
	"github.com/ethersphere/ethproxy/pkg/ethrpc"
)

type Response struct {
	Body *ethrpc.JsonrpcMessage
	IP   string
}

type handlerFunc func(resp *Response)

type handler struct {
	f      handlerFunc
	method string
}

type Callback struct {
	sync.Mutex
	methods  map[uint64]string
	id       int
	handlers map[int]handler
	logger   logging.Logger
}

func New(logger logging.Logger) *Callback {
	return &Callback{
		handlers: make(map[int]handler),
		methods:  make(map[uint64]string),
		logger:   logger,
	}
}

// Register keeps track of the ethrpc requests with ethrpc ID to method name mapping.
func (c *Callback) Register(id uint64, method string) {
	c.Lock()
	defer c.Unlock()

	c.methods[id] = method
}

// On adds a new callback based on an ethrpc method and returns the associated callback ID
func (c *Callback) On(method string, f handlerFunc) int {
	c.Lock()
	defer c.Unlock()
	defer func() { c.id++ }()

	c.handlers[c.id] = handler{f: f, method: method}
	return c.id
}

// Cancel removes a callback based on an callback ID
func (c *Callback) Cancel(id int) error {
	c.Lock()
	defer c.Unlock()

	_, ok := c.handlers[id]
	if !ok {
		return errors.New("callback not found")
	}
	delete(c.handlers, id)
	return nil
}

// Run, with the ID from the ethrpc response, finds the registered method name, and executes
// all callbacks assoicated to the method.
func (c *Callback) Run(resp *Response) {
	c.Lock()
	defer c.Unlock()

	id, err := resp.Body.GetID()
	if err != nil {
		return
	}

	method, ok := c.methods[id]
	if !ok {
		return
	}
	delete(c.methods, id)

	for _, h := range c.handlers {
		if h.method == method {
			h.f(resp)
		}
	}
}
