.PHONY: build

.PHONY: binary
binary: export CGO_ENABLED=0
binary: dist FORCE
	go build -o ./dist/client ./cmd/client
	go build -o ./dist/backend ./cmd/backend
	go build -o ./dist/proxy ./cmd/proxy

.PHONY: proxy
proxy: export CGO_ENABLED=0
proxy: dist FORCE
	go build -o ./dist/proxy ./cmd/proxy

dist:
	mkdir $@

.PHONY: test
test:
	go test -v -failfast -count 1 ./...

FORCE: