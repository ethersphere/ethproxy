// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ethersphere/ethproxy/pkg/api"
	"github.com/ethersphere/ethproxy/pkg/rpc"
)

var ErrStatusNotOK = errors.New("not STATUSOK")

const (
	BlockNumberRecord = rpc.BlockNumberRecord
	BlockNumberFreeze = rpc.BlockNumberFreeze
)

type Client struct {
	endpoint string
	client   *http.Client
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		client:   &http.Client{Timeout: time.Second * 10},
	}
}

func (c *Client) Execute(method string, params ...interface{}) (int, error) {
	b, err := json.Marshal(api.RpcMessage{
		Method: method,
		Params: params,
	})
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/execute", c.endpoint), bytes.NewReader(b))
	if err != nil {
		return 0, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return 0, ErrStatusNotOK
	}

	idmap := make(map[string]int)
	err = json.NewDecoder(res.Body).Decode(&idmap)
	if err != nil {
		return 0, err
	}

	return idmap["id"], nil
}

func (c *Client) Cancel(id int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/cancel/%d", c.endpoint, id), nil)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ErrStatusNotOK
	}

	return nil
}

func (c *Client) State() (*rpc.State, error) {
	req, err := http.NewRequest("GET", c.endpoint+"/state", nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var state rpc.State
	json.NewDecoder(res.Body).Decode(&state)
	if err != nil {
		return nil, err
	}

	return &state, nil
}
