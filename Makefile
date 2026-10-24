# The MyGo CLI is run from the module version below, so no tool install is
# needed. Pin it for reproducibility: make dev MYGO_VERSION=v0.6.3
MYGO_VERSION ?= latest
MYGO = go run github.com/egoist/mygo/cmd/mygo@$(MYGO_VERSION)

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
	CGO_ENABLED=1 go test -race ./... -count=1 -timeout=300s

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
