// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/ethersphere/ethproxy"
	"github.com/ethersphere/ethproxy/pkg/callback"
	"github.com/ethersphere/ethproxy/pkg/rpc"
	"github.com/go-chi/chi"
)

type Api struct {
	rpc  *rpc.Caller
	call *callback.Callback
}

const (
	JSONContent = "application/json; charset=utf-8"
)

const (
	BlockNumberFreeze = "blockNumberFreeze"
	BlockNumberRecord = "blockNumberRecord"
)

func NewServer(call *callback.Callback, port string) *http.Server {

	api := &Api{
		rpc:  rpc.New(call),
		call: call,
	}
	r := chi.NewRouter()

	r.Get("/health", api.status)
	r.Get("/readiness", api.status)
	r.Get("/state", api.state)
	r.Post("/execute", api.execute)
	r.Delete("/cancel/{ID}", api.cancel)

	fmt.Printf("API listing on %v\n", port)

	return &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
}

type RpcMessage struct {
	Method string        `json:"method,omitempty"`
	Params []interface{} `json:"params,omitempty"`
}

func (api *Api) execute(w http.ResponseWriter, r *http.Request) {

	var msg RpcMessage
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	id, err := api.rpc.Execute(msg.Method, msg.Params...)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	respond(w, map[string]int{"handler": id})
}

func (api *Api) cancel(w http.ResponseWriter, r *http.Request) {

	idStr := chi.URLParam(r, "ID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
	}

	api.call.Cancel(int(id))
}

func (api *Api) state(w http.ResponseWriter, r *http.Request) {
	respond(w, api.rpc.GetState())
}

type statusResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func (api *Api) status(w http.ResponseWriter, r *http.Request) {
	respond(w, statusResponse{
		Status:  "ok",
		Version: ethproxy.Version,
	})
}

func respond(w http.ResponseWriter, body interface{}) error {

	w.Header().Set("Content-Type", JSONContent)

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	return err
}

func respondError(w http.ResponseWriter, status int, err error) error {

	w.WriteHeader(status)
	return respond(w, map[string]string{"error": err.Error()})
}
