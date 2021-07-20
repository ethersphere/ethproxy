package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ethersphere/ethproxy/pkg/api"
)

var ErrStatusNotOK = errors.New("not STATUSOK")

type Client struct {
	endpoint string
	client   *http.Client
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		client:   &http.Client{Timeout: time.Second * 10},
	}
}

func (c *Client) Execute(method api.Method, args ...interface{}) error {

	b, err := json.Marshal(api.RpcMessage{
		Method: method,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ErrStatusNotOK
	}

	return nil
}
