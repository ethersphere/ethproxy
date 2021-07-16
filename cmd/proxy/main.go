package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethersphere/ethproxy/pkg/callback"
	"github.com/ethersphere/ethproxy/pkg/rpc"
)

func main() {

	// apiPort := flag.String("apiPort", "6000", "port to listen on")
	port := flag.String("proxyPort", "6000", "port to listen on")
	backend := flag.String("backendEndpoint", "ws://:7000/", "backend endpoint to proxy requests")

	callback := callback.New()

	callback.On(rpc.BlockNumber, func(j *rpc.JsonrpcMessage) {
		j.SetBlockNumber(12)
	})

	callback.On(rpc.BlockNumber, func(j *rpc.JsonrpcMessage) {
		fmt.Println(j.BlockNumber())
	})

	go NewProxy(callback, *port, *backend).ListenAndServe()
	// go NewProxy(callback, *port, *backend).ListenAndServe()

	<-waitTerminate()
}

func waitTerminate() <-chan os.Signal {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, syscall.SIGINT, syscall.SIGTERM)
	return interruptChannel
}
