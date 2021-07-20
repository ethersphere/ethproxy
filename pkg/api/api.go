package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ethersphere/ethproxy/pkg/callback"
	"github.com/ethersphere/ethproxy/pkg/rpc"
)

type Api struct {
	rpc *rpc.Caller
}

const (
	BlockNumberFreeze = "blockNumberFreeze"
	BlockNumberRecord = "blockNumberRecord"
)

func NewServer(call *callback.Callback, port string) *http.Server {

	m := http.NewServeMux()

	api := &Api{
		rpc: rpc.New(call),
	}

	m.HandleFunc("/", api.handler)

	return &http.Server{
		Addr:    ":" + port,
		Handler: m,
	}
}

type RpcMessage struct {
	Method string        `json:"method,omitempty"`
	Params []interface{} `json:"params,omitempty"`
}

func (api *Api) handler(w http.ResponseWriter, r *http.Request) {

	var msg RpcMessage
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Println(err)
		respondError(w, http.StatusBadRequest, err)
		return
	}

	err = api.rpc.Register(msg.Method, msg.Params...)
	if err != nil {
		log.Println(err)
		respondError(w, http.StatusBadRequest, err)
		return
	}
}

func respondError(w http.ResponseWriter, status int, err error) {

	w.WriteHeader(status)

	b, _ := json.Marshal(map[string]string{
		"error": err.Error(),
	})

	w.Write(b)
}
