# Contributing

Use Go 1.27.1 or later. UI changes should use the pinned MyGo public API, typed options and caller-owned state. Read `docs/design.md` before changing a public contract. Keep examples as real consumers and update the component catalog when behavior or limitations change.

## Validate a change

```sh
go mod tidy
make check
make race
golangci-lint run ./...
GOOS=js GOARCH=wasm CGO_ENABLED=0 go build ./...
go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...
```

The ordinary suite runs without windows, external services or user configuration. Linux software rendering requires Pango; race detection requires a C compiler. Check the modified control with real native input on the relevant OS before claiming platform accessibility, GPU redraw or gesture acceptance. Use a disposable application, with no user data.

Tests should verify actions, editing, focus, selection, source fidelity, state transitions and real regressions. Preserve stable identities while filtering/reordering or hiding/reopening widgets. Prefer the toolkit's virtual clock to sleeps. Optional renderer artifacts use `MYGOUI_*` environment variables and paths outside tracked source.

## Dependencies

Pin stable MyGo and mygoicons releases in `go.mod`. Read upstream release notes and migration guidance, run `go tool mygo migrate-ui .` in preview, then verify public source before changing assumptions. Record relevant integration decisions in `docs/mygo-0.6.md` or a successor upgrade note. Run tidy, integrity, vulnerability, build and behavior checks. Keep releases free of local replacements and uncommitted workspaces.

For sibling-module development, use an untracked `go.work`. The published library must build without that workspace. Before release, test an external consumer from the module archive with `GOWORK=off`; after publishing, also test fetching the exact public tag from a fresh module with no replace directives.

## Release

1. Review public API changes and `CHANGELOG.md`; use semantic versioning. While v0, incompatible changes require a minor bump, and patches preserve public behavior.
2. Run the validation commands and inspect CI on Linux, macOS and Windows. Record remaining platform limits honestly in `VALIDATION.md`.
3. Verify module metadata with `go mod tidy`, `go mod verify` and a clean diff. Check license/notices and example import paths. Confirm the external archive consumer uses only resolvable dependencies.
4. Replace the Unreleased heading with the chosen version/date and update installation examples to an actual version.
5. Publish the reviewed repository to `github.com/Cyenoch/mygoui`. Create and push the reviewed `v0.x.y` tag, then create its GitHub release. There is no automatic publishing workflow.
6. From a fresh external module, fetch `github.com/Cyenoch/mygoui@v0.x.y`, build a sample application and run a headless component action. Verify the exact version and module checksum without local replacements before announcing availability.

Repository creation, push, tagging and release publication are separate maintainer actions. Local validation does not establish public availability.
