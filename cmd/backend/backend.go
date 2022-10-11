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

func NewBackend() *http.Server {
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
		_, msg, err := c.ReadMessage()
		if err != nil {
			fmt.Printf("server: %v\n", err)
			return
		}

		jmsg, err := ethrpc.Unmarshall(msg)
		if err != nil {
			return
		}

		id, err := jmsg.GetID()
		if err != nil {
			return
		}

		resp := []byte(fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":"0x3ec"}`, id))

		err = c.WriteMessage(websocket.BinaryMessage, resp)

		// err = c.WriteJSON(ethrpc.JsonrpcMessage{
		// 	Method: ethrpc.BlockNumber,
		// 	Result: json.RawMessage(fmt.Sprintf(`"%s"`, hexutil.Uint64(blockN).String())),
		// })
		if err != nil {
			fmt.Printf("server: %v\n", err)
			return
		}
	}
}
