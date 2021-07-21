// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ethrpc

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type jsonError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (err *jsonError) Error() string {
	if err.Message == "" {
		return fmt.Sprintf("json-rpc error %d", err.Code)
	}
	return err.Message
}

type JsonrpcMessage struct {
	Version string          `json:"jsonrpc,omitempty"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Error   *jsonError      `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}

const (
	BlockByHash        = "eth_getBlockByHash"
	BlockNumber        = "eth_blockNumber"
	TransactionReceipt = "eth_getTransactionReceipt"
)

func (j *JsonrpcMessage) BlockNumber() (uint64, error) {
	var result hexutil.Uint64
	err := json.Unmarshal(j.Result, &result)
	return uint64(result), err
}

func (j *JsonrpcMessage) SetBlockNumber(n uint64) error {
	var hexN hexutil.Uint64 = hexutil.Uint64(n)
	j.Result = json.RawMessage(fmt.Sprintf(`"%s"`, hexN.String()))
	return nil
}

func Unmarshall(data json.RawMessage) (*JsonrpcMessage, error) {
	var msg JsonrpcMessage
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

func (j *JsonrpcMessage) Marshall() (json.RawMessage, error) {
	return json.Marshal(j)
}
