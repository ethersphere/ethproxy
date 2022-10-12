// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethersphere/beekeeper/pkg/logging"
	"github.com/ethersphere/ethproxy"
	"github.com/ethersphere/ethproxy/pkg/callback"
	"github.com/ethersphere/ethproxy/pkg/rpc"
	"github.com/go-chi/chi"
)

type Api struct {
	rpc    *rpc.Caller
	call   *callback.Callback
	logger logging.Logger
}

const (
	JSONContent = "application/json; charset=utf-8"
)

const (
	BlockNumberFreeze = "blockNumberFreeze"
	BlockNumberRecord = "blockNumberRecord"
)

func NewApi(call *callback.Callback, rpc *rpc.Caller, logger logging.Logger) *Api {
	return &Api{
		rpc:    rpc,
		call:   call,
		logger: logger,
	}
}

func (api *Api) Server(port string) *http.Server {
	r := chi.NewRouter()
	r.Get("/health", api.status)
	r.Get("/readiness", api.status)
	r.Get("/state", api.state)
	r.Post("/execute", api.Execute)
	r.Delete("/cancel/{ID}", api.cancel)

	api.logger.Infof("API listing on %s", port)

	return &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
}

type RpcMessage struct {
	Method string        `json:"method,omitempty"`
	Params []interface{} `json:"params,omitempty"`
}

func (api *Api) Execute(w http.ResponseWriter, r *http.Request) {
	var msg RpcMessage
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	api.logger.Infof("api: execute %s %v", msg.Method, msg.Params)

	id, err := api.rpc.Execute(msg.Method, msg.Params...)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	respond(w, map[string]int{"id": id})
}

func (api *Api) cancel(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "ID")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
	}

	api.logger.Infof("api: cancel %d", id)

	err = api.call.Cancel(int(id))
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
	}
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
