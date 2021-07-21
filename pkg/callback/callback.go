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
	id       int
	handlers map[int]handler
}

func New() *Callback {
	return &Callback{
		handlers: make(map[int]handler),
	}
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

	for _, h := range c.handlers {
		if h.method == resp.Body.Method {
			h.f(resp)
		}
	}
}
