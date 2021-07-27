package rpc

import (
	"errors"

	"github.com/ethersphere/ethproxy/pkg/callback"
	"github.com/ethersphere/ethproxy/pkg/ethrpc"
)

func (c *Caller) blockNumberRecord() (int, error) {
	return c.call.On(ethrpc.BlockNumber, func(resp *callback.Response) {
		bN, err := resp.Body.BlockNumber()
		if err != nil {
			return
		}
		c.logger.Infof("block number update: %d", bN)
		c.state.BlockNumber = bN
	}), nil
}

func (c *Caller) blockNumberFreeze(blockN uint64, params []interface{}) (int, error) {
	if len(params) == 0 {
		return c.call.On(ethrpc.BlockNumber, func(resp *callback.Response) {
			resp.Body.SetBlockNumber(blockN)
		}), nil
	} else {

		ips, err := stringArray(params)
		if err != nil {
			return 0, err
		}

		return func(ips []string) int {
			return c.call.On(ethrpc.BlockNumber, func(resp *callback.Response) {
				for _, ip := range ips {
					if resp.IP == ip {
						resp.Body.SetBlockNumber(blockN)
						c.logger.Infof("ip %s: frozen block number %d", ip, blockN)
					}
				}
			})
		}(ips), nil
	}
}

func stringArray(args []interface{}) ([]string, error) {

	ret := make([]string, len(args))

	for _, arg := range args {
		str, ok := arg.(string)
		if !ok {
			return nil, errors.New("bad param")
		}
		ret = append(ret, str)
	}

	return ret, nil
}
