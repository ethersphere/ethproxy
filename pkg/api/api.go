package api

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/ethersphere/ethproxy/pkg/callback"
	"github.com/ethersphere/ethproxy/pkg/rpc"
)

type Api struct {
	call *callback.Callback

	mtx         sync.Mutex
	blockNumber uint64
}

const (
	BlockNumberFreeze = "blockNumberFreeze"
	BlockNumberRecord = "blockNumberRecord"
)

func NewServer(call *callback.Callback, port string) *http.Server {

	m := http.NewServeMux()

	api := &Api{
		call: call,
	}

	m.HandleFunc("/", api.handler)

	return &http.Server{
		Addr:    ":" + port,
		Handler: m,
	}
}

type RpcMessage struct {
	Method string          `json:"method,omitempty"`
	Params json.RawMessage `json:"params,omitempty"`
}

func (api *Api) handler(w http.ResponseWriter, r *http.Request) {

	var msg RpcMessage
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		respond(w, http.StatusBadRequest)
		return
	}

	switch msg.Method {

	case BlockNumberRecord:
		api.call.On(rpc.BlockNumber, func(j *rpc.JsonrpcMessage) {
			bN, err := j.BlockNumber()
			if err != nil {
				return
			}
			api.mtx.Lock()
			api.blockNumber = bN
			api.mtx.Unlock()
		})
	case BlockNumberFreeze:
		api.call.On(rpc.BlockNumber, func(j *rpc.JsonrpcMessage) {
			j.SetBlockNumber(api.blockNumber)
		})
	}

}

func respond(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}
