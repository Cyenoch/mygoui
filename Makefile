# The tool directive resolves the same MyGo version as the library.
MYGO = go tool mygo
OPENSPEC = node scripts/openspec.mjs
IMPECCABLE = .agents/skills/impeccable/scripts/impeccable

.PHONY: check test race vet build lint wasm workflows gallery basic dev app refs agents agents-update spec-check design-context design-doctor design-artifacts

check: vet test build
	@test -z "$$(gofmt -l .)"
	go mod verify

vet:
	CGO_ENABLED=0 go vet ./...
	CGO_ENABLED=0 $(MYGO) vet .

test:
	CGO_ENABLED=0 go test ./... -count=1 -timeout=180s

race:
	CGO_ENABLED=1 go test -race -p 1 ./... -count=1 -timeout=600s

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

# Agent tooling is optional for Go library consumers.
agents:
	npm ci --ignore-scripts --no-audit --no-fund
	npm run agents:setup
	$(OPENSPEC) update --force
	$(IMPECCABLE) hooks on

agents-update:
	npm exec --no -- impeccable update --project --yes --no-hooks
	npm exec --no -- impeccable install --providers=codex --scope=project --yes --force --no-hooks
	npm run agents:setup
	$(OPENSPEC) update --force
	$(IMPECCABLE) hooks on

spec-check:
	$(OPENSPEC) validate --all --strict --no-interactive
	$(OPENSPEC) validate --archived --no-interactive

design-context:
	$(IMPECCABLE) context $(if $(TARGET),--target "$(TARGET)")

design-doctor:
	$(IMPECCABLE) doctor --json

design-artifacts:
	MYGOUI_COSS_ARTIFACT_DIR="$(CURDIR)/artifacts/design-review" CGO_ENABLED=0 go test ./internal/gallery -run '^TestCossVisualArtifacts$$' -count=1 -timeout=180s
