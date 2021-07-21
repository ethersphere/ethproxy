// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpc_test

import (
	"testing"

	"github.com/ethersphere/ethproxy/pkg/callback"
	"github.com/ethersphere/ethproxy/pkg/ethrpc"
	"github.com/ethersphere/ethproxy/pkg/rpc"
)

func TestBlockNumberRecord(t *testing.T) {

	var (
		call          = callback.New()
		r             = rpc.New(call)
		blockN uint64 = 10
	)

	_, err := r.Execute(rpc.BlockNumberRecord)
	if err != nil {
		t.Fatal(err)
	}

	resp := &callback.Response{
		Body: &ethrpc.JsonrpcMessage{
			Method: ethrpc.BlockNumber,
		},
	}
	resp.Body.SetBlockNumber(blockN)

	call.Run(resp)

	if r.GetState().BlockNumber != blockN {
		t.Fatalf("got %v, expected %v", r.GetState().BlockNumber, blockN)
	}
}

func TestBlockNumberFreeze(t *testing.T) {

	var (
		call          = callback.New()
		r             = rpc.New(call)
		ip            = "1.0.0.0"
		blockN uint64 = 10
	)

	_, err := r.Execute(rpc.BlockNumberRecord)
	if err != nil {
		t.Fatal(err)
	}

	resp := &callback.Response{
		Body: &ethrpc.JsonrpcMessage{
			Method: ethrpc.BlockNumber,
		},
		IP: ip,
	}

	resp.Body.SetBlockNumber(blockN)
	call.Run(resp)

	_, err = r.Execute(rpc.BlockNumberFreeze, ip)
	if err != nil {
		t.Fatal(err)
	}

	resp.Body.SetBlockNumber(blockN + 1)
	call.Run(resp)

	if r.GetState().BlockNumber != blockN {
		t.Fatalf("got %v, expected %v", r.GetState().BlockNumber, blockN)
	}
}

func TestCancel(t *testing.T) {

	var (
		call          = callback.New()
		r             = rpc.New(call)
		blockN uint64 = 10
	)

	id, err := r.Execute(rpc.BlockNumberRecord)
	if err != nil {
		t.Fatal(err)
	}

	resp := &callback.Response{
		Body: &ethrpc.JsonrpcMessage{
			Method: ethrpc.BlockNumber,
		},
	}

	resp.Body.SetBlockNumber(blockN)
	call.Run(resp)

	call.Cancel(id)

	resp.Body.SetBlockNumber(blockN + 1)
	call.Run(resp)

	if r.GetState().BlockNumber != blockN {
		t.Fatalf("got %v, expected %v", r.GetState().BlockNumber, blockN)
	}
}
