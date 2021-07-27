// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc_test

import (
	"io/ioutil"
	"testing"

	"github.com/ethersphere/bee/pkg/logging"
	"github.com/ethersphere/ethproxy/pkg/callback"
	"github.com/ethersphere/ethproxy/pkg/ethrpc"
	"github.com/ethersphere/ethproxy/pkg/rpc"
)

func TestBlockNumberRecord(t *testing.T) {

	var (
		logger        = logging.New(ioutil.Discard, 0)
		call          = callback.New(logger)
		r             = rpc.New(call, logger)
		blockN uint64 = 10
		method        = ethrpc.BlockNumber
	)

	_, err := r.Execute(rpc.BlockNumberRecord)
	if err != nil {
		t.Fatal(err)
	}

	resp := &callback.Response{
		Body: &ethrpc.JsonrpcMessage{
			Method: method,
			ID:     []byte("0"),
		},
	}
	resp.Body.SetBlockNumber(blockN)
	call.Register(0, method)
	call.Run(resp)

	if r.GetState().BlockNumber != blockN {
		t.Fatalf("got %v, expected %v", r.GetState().BlockNumber, blockN)
	}
}

func TestBlockNumberRecordCancel(t *testing.T) {

	var (
		logger        = logging.New(ioutil.Discard, 0)
		call          = callback.New(logger)
		r             = rpc.New(call, logger)
		ip            = "1.0.0.0"
		blockN uint64 = 10
		method        = ethrpc.BlockNumber
	)

	recordID, err := r.Execute(rpc.BlockNumberRecord)
	if err != nil {
		t.Fatal(err)
	}

	resp := &callback.Response{
		Body: &ethrpc.JsonrpcMessage{
			Method: method,
			ID:     []byte("0"),
		},
		IP: ip,
	}

	resp.Body.SetBlockNumber(blockN)
	call.Register(0, method)
	call.Run(resp)

	call.Cancel(recordID)

	_, err = r.Execute(rpc.BlockNumberFreeze, ip)
	if err != nil {
		t.Fatal(err)
	}

	resp.Body.SetBlockNumber(blockN + 1)
	call.Register(0, method)
	call.Run(resp)

	if r.GetState().BlockNumber != blockN {
		t.Fatalf("got %v, expected %v", r.GetState().BlockNumber, blockN)
	}
}
