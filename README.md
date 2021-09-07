# Ethereum Backend Proxy

To locally run the example client, proxy, and backend:

- `make binary`
- `./dist/backend`
- `PROXY_BACKEND_ENDPOINT=ws://localhost:7000 ./dist/proxy`
- `./dist/client`

Env Flags Needed for the proxy

- `PROXY_API_PORT` Endpoint for executing proxy methods 
- `PROXY_BACKEND_ENDPOINT` Backend endpoint to proxy requests to
- `PROXY_WS_PORT` Websocket port to connect to
- `PROXY_LOG_LEVEL` "info" is default, "debug" prints all `jsonRPC` requests and responses


## Proxy RPC Methods

Methods can be found in /pkg/api/api.go

To interact with the proxy and execute internal commands to, for example,
record or alter the `block number` of the chain, a json RPC call may be sent to the proxy API as such:

```
POST /execute
{"method": "blockNumberRecord"}
```
```
POST /execute
{"method": "blockNumberFreeze", "params": ["some-ip4-address"]}}
```

To get the ip address of a Bee node, you may query the debug api at `/addresses` and use the public ip4 address.

To cancel the above executions, grab the id in the response and cancel it as such:
```
DELETE /cancel/{id}
```

Example commands:
```
curl --data '{"method": "blockNumberRecord"}' ethproxy.localhost/execute
curl --data '{"method": "blockNumberFreeze", "params": ["10.42.0.37"]}' ethproxy.localhost/execute
curl -X DELETE ethproxy.localhost/cancel/5
```


## Deployment

### ethersphere/beelocal
```bash
OPTS=skip-vet ./beelocal.sh
```

### ethersphere/ethproxy
```bash
cd deploy
helmsman -apply -f ethproxy.yaml
```

### ethersphere/beekeeper

Copy `beekeeper/config/local.yml` to `~/.beekeeper/local.yaml` and change `swap-endpoint` to `ws://ethproxy:8546`

```bash
beekeeper create bee-cluster --cluster-name local
```
### Logs

```bash
kubectl get pods -n local
kubectl logs -n local -f {proxy container name}
```