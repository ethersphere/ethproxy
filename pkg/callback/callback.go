// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package callback

import (
	"sync"

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
}

func New() *Callback {
	return &Callback{
		handlers: make(map[int]handler),
		methods:  make(map[uint64]string),
	}
}

func (c *Callback) Register(id uint64, method string) {
	c.Lock()
	defer c.Unlock()

	c.methods[id] = method
}

func (c *Callback) On(method string, f handlerFunc) int {
	c.Lock()
	defer c.Unlock()
	defer func() { c.id++ }()

	c.handlers[c.id] = handler{f: f, method: method}
	return c.id
}

func (c *Callback) Cancel(id int) {
	c.Lock()
	defer c.Unlock()
	delete(c.handlers, id)
}

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
