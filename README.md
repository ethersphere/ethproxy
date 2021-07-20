# Ethereum Backend Proxy

To run the example client, proxy, and backend:

- `make binary`
- `./dist/server`
- `PROXY_BACKEND_ENDPOINT=ws://localhost:7000 ./dist/proxy`
- `./dist/client`

Env Flags Needed for the proxy

- `PROXY_API_PORT` Endpoint for executing proxy methods 
- `PROXY_BACKEND_ENDPOINT` Backend endpoint to proxy requests to
- `PROXY_WS_PORT` Websocket port to connect to


## Proxy Methods

Methods can be found in /pkg/api/api.go

To interact with the proxy and execute internal commands to, for example,
record or alter the `block number` of the chain, a json RPC call may be sent to the proxy API as such:

```
POST /
{"method": "blockNumberRecord"}
```
```
POST /
{"method": "blockNumberFreeze"}
```