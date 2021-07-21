// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ethersphere/ethproxy"
	"github.com/ethersphere/ethproxy/pkg/callback"
	"github.com/ethersphere/ethproxy/pkg/rpc"
)

type Api struct {
	rpc *rpc.Caller
}

const (
	JSONContent = "application/json; charset=utf-8"
)

const (
	BlockNumberFreeze = "blockNumberFreeze"
	BlockNumberRecord = "blockNumberRecord"
)

func NewServer(call *callback.Callback, port string) *http.Server {

	m := http.NewServeMux()

	api := &Api{
		rpc: rpc.New(call),
	}

	m.HandleFunc("/health", api.status)
	m.HandleFunc("/readiness", api.status)
	m.HandleFunc("/state", api.state)
	m.HandleFunc("/", api.handler)

	fmt.Printf("API listing on %v\n", port)

	return &http.Server{
		Addr:    ":" + port,
		Handler: m,
	}
}

type RpcMessage struct {
	Method string        `json:"method,omitempty"`
	Params []interface{} `json:"params,omitempty"`
}

func (api *Api) handler(w http.ResponseWriter, r *http.Request) {

	var msg RpcMessage
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Println(err)
		respondError(w, http.StatusBadRequest, err)
		return
	}

	err = api.rpc.Execute(msg.Method, msg.Params...)
	if err != nil {
		log.Println(err)
		respondError(w, http.StatusBadRequest, err)
		return
	}
}

func (api *Api) state(w http.ResponseWriter, r *http.Request) {

	b, err := json.Marshal(api.rpc.GetState())
	if err != nil {
		log.Println(err)
		respondError(w, http.StatusBadRequest, err)
		return
	}

	w.Header().Set("Content-Type", JSONContent)

	w.Write(b)
}

type statusResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func (api *Api) status(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", JSONContent)

	b, _ := json.Marshal(statusResponse{
		Status:  "ok",
		Version: ethproxy.Version,
	})

	w.Write(b)
}

func respondError(w http.ResponseWriter, status int, err error) {

	w.WriteHeader(status)
	w.Header().Set("Content-Type", JSONContent)

	b, _ := json.Marshal(map[string]string{
		"error": err.Error(),
	})

	w.Write(b)
}
