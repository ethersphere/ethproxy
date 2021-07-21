// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"net/http"

	"github.com/ethersphere/ethproxy/pkg/ethrpc"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: false,
}

func NewServer() *http.Server {
	r := chi.NewRouter()
	r.HandleFunc("/", serverRoute)

	return &http.Server{
		Addr:    ":7000",
		Handler: r,
	}
}

func serverRoute(w http.ResponseWriter, r *http.Request) {

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("server: %v\n", err)
		return
	}
	defer c.Close()

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			fmt.Printf("server: %v\n", err)
			return
		}
		err = c.WriteJSON(ethrpc.JsonrpcMessage{
			Method: ethrpc.BlockNumber,
			Result: []byte("10"),
		})
		if err != nil {
			fmt.Printf("server: %v\n", err)
			return
		}
	}
}
