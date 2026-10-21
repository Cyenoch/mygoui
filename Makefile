.PHONY: check test race vet build lint gallery basic

check: vet test build
	@test -z "$$(gofmt -l .)"
	go mod verify

vet:
	CGO_ENABLED=0 go vet ./...
	CGO_ENABLED=0 go tool mygo vet .

test:
	CGO_ENABLED=0 go test ./... -count=1 -timeout=180s

race:
	CGO_ENABLED=1 go test -race ./... -count=1 -timeout=300s

build:
	CGO_ENABLED=0 go build ./...

lint:
	golangci-lint run ./...

gallery:
	go run ./examples/gallery

basic:
	go run ./examples/basic
