package main

import (
	"fmt"
	"net/http"

	"github.com/ethersphere/ethproxy/pkg/callback"
	"github.com/ethersphere/ethproxy/pkg/rpc"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: false,
}

func NewProxy(call *callback.Callback) *http.Server {
	r := chi.NewRouter()

	r.HandleFunc("/", wsRoute(call))

	return &http.Server{
		Addr:    ":6000",
		Handler: r,
	}
}

func wsRoute(call *callback.Callback) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Printf("proxy: %v\n", err)
			return
		}
		defer conn.Close()

		serverClient, err := serverClient()
		if err != nil {
			fmt.Printf("proxy: %v\n", err)
			return
		}
		defer serverClient.Close()

		go func() {
			for {
				t, msg, err := conn.ReadMessage()
				if err != nil {
					if _, ok := err.(*websocket.CloseError); !ok {
						panic(err)
					} else {
						return
					}
				}

				err = serverClient.WriteMessage(t, msg)
				if err != nil {
					panic(err)
				}
			}
		}()

		for {
			t, msg, err := serverClient.ReadMessage()
			if err != nil {
				panic(err)
			}

			jmsg, err := rpc.Unmarshall(msg)
			if err != nil {
				call.Run(jmsg)
			}

			err = conn.WriteMessage(t, msg)
			if err != nil {
				panic(err)
			}
		}
	}
}

func serverClient() (*websocket.Conn, error) {
	url := fmt.Sprintf("ws://%s/", ":7000")
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
