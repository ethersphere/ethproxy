// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proxy

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/ethersphere/ethproxy/pkg/callback"
	"github.com/ethersphere/ethproxy/pkg/ethrpc"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: false,
}

type proxy struct {
	call            *callback.Callback
	backendEndpoint string
}

func NewProxy(call *callback.Callback, backendEndpoing string) *proxy {

	return &proxy{
		call:            call,
		backendEndpoint: backendEndpoing,
	}
}

func (p *proxy) Serve(port string) error {
	m := http.NewServeMux()

	m.HandleFunc("/", p.Handle)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: m,
	}

	return server.ListenAndServe()
}

func (p *proxy) Handle(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("proxy: %v\n", err)
		return
	}
	defer conn.Close()

	backend, err := p.backendClient()
	if err != nil {
		fmt.Printf("proxy: %v\n", err)
		return
	}
	defer backend.Close()

	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			t, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			p.rpcRequest(msg)

			fmt.Printf("CLIENT %v\n", string(msg))

			err = backend.WriteMessage(t, msg)
			if err != nil {
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			t, msg, err := backend.ReadMessage()
			if err != nil {
				return
			}

			fmt.Printf("BACKEND %v\n", string(msg))

			fmt.Println("IP:", ip)
			msg, err = p.rpcResponse(ip, msg)
			if err != nil {
				log.Print(err)
			}

			err = conn.WriteMessage(t, msg)
			if err != nil {
				return
			}
		}
	}()

	wg.Wait()
}

func (p *proxy) rpcRequest(msg []byte) error {

	jmsg, err := ethrpc.Unmarshall(msg)
	if err != nil {
		return err
	}

	id, err := jmsg.GetID()
	if err != nil {
		return err
	}

	p.call.Register(id, jmsg.Method)

	return nil
}

func (p *proxy) rpcResponse(ip string, msg []byte) ([]byte, error) {

	jmsg, err := ethrpc.Unmarshall(msg)
	if err != nil {
		return msg, err
	}

	p.call.Run(&callback.Response{
		Body: jmsg,
		IP:   ip,
	})

	bjmsg, err := jmsg.Marshall()
	if err != nil {
		return msg, err
	}

	return bjmsg, nil
}

func (p *proxy) backendClient() (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(p.backendEndpoint, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
