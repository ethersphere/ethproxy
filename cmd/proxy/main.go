// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethersphere/bee/pkg/logging"
	"github.com/ethersphere/ethproxy/pkg/api"
	"github.com/ethersphere/ethproxy/pkg/callback"
	"github.com/ethersphere/ethproxy/pkg/proxy"
	"github.com/ethersphere/ethproxy/pkg/rpc"
	"github.com/sirupsen/logrus"
)

func main() {

	port := getEnv("PROXY_WS_PORT", "6000")
	apiPort := getEnv("PROXY_API_PORT", "6100")
	backend := getEnv("PROXY_BACKEND_ENDPOINT", "ws://geth-swap:8546")
	logLevel := getEnv("PROXY_LOG_LEVEL", "info")

	logger, _ := newLogger(logLevel)
	callback := callback.New(logger)
	rpc := rpc.New(callback, logger)

	go func() { log.Fatal(proxy.NewProxy(callback, backend).Serve(port)) }()
	go func() { log.Fatal(api.NewApi(callback, rpc, logger).Serve(apiPort)) }()

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

func newLogger(verbosity string) (logging.Logger, error) {
	var logger logging.Logger
	switch verbosity {
	case "0", "silent":
		logger = logging.New(ioutil.Discard, 0)
	case "1", "error":
		logger = logging.New(os.Stdout, logrus.ErrorLevel)
	case "2", "warn":
		logger = logging.New(os.Stdout, logrus.WarnLevel)
	case "3", "info":
		logger = logging.New(os.Stdout, logrus.InfoLevel)
	case "4", "debug":
		logger = logging.New(os.Stdout, logrus.DebugLevel)
	case "5", "trace":
		logger = logging.New(os.Stdout, logrus.TraceLevel)
	default:
		return nil, fmt.Errorf("unknown verbosity level %q", verbosity)
	}
	return logger, nil
}
