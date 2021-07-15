.PHONY: build
binary:
	go build -o ./dist/client ./cmd/client
	go build -o ./dist/server ./cmd/server
	go build -o ./dist/proxy ./cmd/proxy

	