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

type handler func(resp *Response)

type Callback struct {
	mtx      sync.Mutex
	handlers map[string][]handler
}

func New() *Callback {
	return &Callback{
		handlers: make(map[string][]handler),
	}
}

func (c *Callback) On(method string, f handler) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.handlers[method] = append(c.handlers[method], f)
}

func (c *Callback) Remove(method string, f handler) {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	delete(c.handlers, method)
}

func (c *Callback) Run(resp *Response) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	for _, h := range c.handlers[resp.Body.Method] {
		h(resp)
	}
}
