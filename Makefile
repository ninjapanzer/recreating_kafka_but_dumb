.PHONY: build-e2e build-server build test
test:
	go test -v ./...
build-e2e:
	go build -o e2e ./cmd/e2e/e2e.go
build-server:
	go build -o server ./cmd/server.go
build: test build-e2e build-server
