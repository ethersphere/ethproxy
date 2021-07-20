package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethersphere/ethproxy/pkg/api"
	"github.com/ethersphere/ethproxy/pkg/callback"
	"github.com/ethersphere/ethproxy/pkg/proxy"
)

func main() {

	port := getEnv("PROXY_WS_PORT", "6000")
	apiPort := getEnv("PROXY_API_PORT", "6100")
	backend := getEnv("PROXY_BACKEND_ENDPOINT", "ws://geth-swap:8546")

	callback := callback.New()

	go log.Fatal(proxy.NewProxy(callback, port, backend).ListenAndServe())
	go log.Fatal(api.NewServer(callback, apiPort).ListenAndServe())

	<-waitTerminate()
}

func waitTerminate() <-chan os.Signal {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, syscall.SIGINT, syscall.SIGTERM)
	return interruptChannel
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
