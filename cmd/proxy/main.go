package main

import (
	"fmt"

	"github.com/ethersphere/ethproxy/pkg/callback"
	"github.com/ethersphere/ethproxy/pkg/rpc"
)

func main() {

	callback := callback.New()

	callback.On(rpc.BlockByHash, func(j *rpc.JsonrpcMessage) {
		fmt.Println(j.Method)
		j.SetBlockNumber(12)
	})

	NewProxy(callback).ListenAndServe()
}
