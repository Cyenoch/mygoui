# The tool directive resolves the same MyGo version as the library.
MYGO = go tool mygo

.PHONY: check test race vet build lint wasm workflows gallery basic dev app refs

check: vet test build
	@test -z "$$(gofmt -l .)"
	go mod verify

vet:
	CGO_ENABLED=0 go vet ./...
	CGO_ENABLED=0 $(MYGO) vet .

test:
	CGO_ENABLED=0 go test ./... -count=1 -timeout=180s

race:
	CGO_ENABLED=1 go test -race -p 1 ./... -count=1 -timeout=300s

build:
	CGO_ENABLED=0 go build ./...

lint:
	golangci-lint run ./...

wasm:
	GOOS=js GOARCH=wasm CGO_ENABLED=0 go build ./...

workflows:
	actionlint

gallery:
	go run ./examples/gallery

basic:
	go run ./examples/basic

dev:
	$(MYGO) dev

app:
	$(MYGO) build

refs:
	git submodule update --init --depth 1
