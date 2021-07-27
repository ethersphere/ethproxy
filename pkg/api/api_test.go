// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ethersphere/bee/pkg/logging"
	"github.com/ethersphere/ethproxy/pkg/api"
	"github.com/ethersphere/ethproxy/pkg/callback"
	"github.com/ethersphere/ethproxy/pkg/ethrpc"
	"github.com/ethersphere/ethproxy/pkg/rpc"
)

func TestBlockNumberFreezeForIP(t *testing.T) {

	var (
		logger        = logging.New(ioutil.Discard, 0)
		call          = callback.New(logger)
		r             = rpc.New(call, logger)
		a             = api.NewApi(call, r, logger)
		ip            = "10.0.0.0"
		id     uint64 = 0
		block  uint64 = 0
	)

	// STEP X: set block number
	_, err := r.Execute(rpc.BlockNumberRecord)
	if err != nil {
		t.Fatal(err)
	}
	runBlockNumber(call, id, block, "")
	block++
	id++

	// STEP X: send api request to freeze block # for ip
	req := api.RpcMessage{
		Method: api.BlockNumberFreeze,
		Params: []interface{}{ip},
	}
	doRequest(t, http.HandlerFunc(a.Execute), req, nil)

	call.On(ethrpc.BlockNumber, func(resp *callback.Response) {
		bN, err := resp.Body.BlockNumber()
		if err != nil {
			t.Fatal(err)
		}

		if bN != block-1 {
			t.Fatalf("got %v, expected %v", bN, block-1)
		}
	})

	// STEP X: set new block number
	runBlockNumber(call, id, block, ip)
}

func runBlockNumber(call *callback.Callback, ethID uint64, blockN uint64, ip string) {
	resp := &callback.Response{
		Body: &ethrpc.JsonrpcMessage{
			Method: ethrpc.BlockNumber,
			ID:     json.RawMessage(fmt.Sprintf("%d", ethID)),
		},
		IP: ip,
	}
	resp.Body.SetBlockNumber(blockN)
	call.Register(ethID, ethrpc.BlockNumber)
	call.Run(resp)
}

func doRequest(t *testing.T, h http.Handler, req, res interface{}) {
	w := httptest.NewRecorder()

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(req); err != nil {
		t.Fatalf("encode request body: %v", err)
	}

	r, err := http.NewRequest("", "", &buf)
	if err != nil {
		t.Fatalf("new request: %v", err)
	}

	h.ServeHTTP(w, r)

	if res != nil {
		err = json.NewDecoder(w.Body).Decode(res)
		if err != nil {
			t.Fatalf("decode response: %v", err)
		}
	}

}

// func ethBlockNumberReq(id int) []byte {
// 	return []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"method":"eth_blockNumber"}`, id))
// }

// func ethBlockNumberRes(id int, hex string) []byte {
// 	return []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":%s}`, id, hex))
// }
