package main

import (
	"fmt"
	"net/http"

	"github.com/ethersphere/ethproxy/pkg/rpc"
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
		err = c.WriteJSON(&rpc.JsonrpcMessage{
			Method: rpc.BlockByHash,
		})
		if err != nil {
			fmt.Printf("server: %v\n", err)
			return
		}
	}
}
