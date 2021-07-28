// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"io/ioutil"
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

	var (
		port     = getEnv("PROXY_WS_PORT", "6000")
		apiPort  = getEnv("PROXY_API_PORT", "6100")
		backend  = getEnv("PROXY_BACKEND_ENDPOINT", "ws://geth-swap8546")
		logLevel = getEnv("PROXY_LOG_LEVEL", "info")

		logger, _ = newLogger(logLevel)
		callback  = callback.New(logger)
		rpc       = rpc.New(callback, logger)

		apiServer   = api.NewApi(callback, rpc, logger).Server(apiPort)
		proxyServer = proxy.NewProxy(callback, backend, logger).Server(port)

		ctx = context.Background()
	)

	go func() { apiServer.ListenAndServe() }()
	go func() { proxyServer.ListenAndServe() }()

	<-terminateChan()

	apiServer.Shutdown(ctx)
	proxyServer.Shutdown(ctx)
}

func terminateChan() <-chan os.Signal {
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
