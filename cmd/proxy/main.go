package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethersphere/ethproxy/pkg/api"
	"github.com/ethersphere/ethproxy/pkg/callback"
	"github.com/ethersphere/ethproxy/pkg/proxy"
)

func main() {

	apiPort := flag.String("apiPort", "6100", "port to listen on")
	port := flag.String("proxyPort", "6000", "port to listen on")
	backend := flag.String("backendEndpoint", "ws://:7000/", "backend endpoint to proxy requests")

	callback := callback.New()

	go log.Fatal(proxy.NewProxy(callback, *port, *backend).ListenAndServe())
	go log.Fatal(api.NewServer(callback, *apiPort).ListenAndServe())

	<-waitTerminate()
}

func waitTerminate() <-chan os.Signal {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, syscall.SIGINT, syscall.SIGTERM)
	return interruptChannel
}
