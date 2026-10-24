# Contributing

Use Go 1.27.1 or later and MyGo v0.6.3. The library is a native Go translation
of HeroUI v3, not a React port: typed options, caller-owned state, native MyGo
controls and no CSS/DOM runtime.

## Set up

```sh
git submodule update --init --depth 1
```

`references/heroui` and `references/heroui-native` are pinned visual/API
references. Read the relevant CSS, variant definitions and TSX before changing
a component. `docs/references.md` records the update procedure.

## Rules

- Read `docs/design.md`, `docs/components.md` and `docs/mygo-0.6.md` first.
- Keep components native: no providers, ambient scopes, build values stored in
  structs, parallel render caches or per-frame allocation in paint callbacks.
- Coordinated collections take typed slices. Freely composed content uses body
  callbacks. Caller state stays explicit and window-scoped.
- HeroUI `Start` is the zero `ui.Align` value; do not use zero as an "unset"
  sentinel for alignment. MyGo `ui.Row` centers across its axis by default, so
  set explicit alignment where CSS would not center.
- Use `theme.Colors(c)` live, `internal/style` for text/paint, stable keys in
  collections, and `FocusRing(true)` for keyboard focus.
- Update the component catalog and limits when public behavior changes.

## Validate

```sh
go mod tidy
make check
make race
make dev      # go run .../cmd/mygo@latest dev
golangci-lint run ./...
make wasm
make workflows
make vet
go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...
```

The headless gallery suite renders every registered page at compact and wide
widths in light and dark modes. Page-local tests exercise real controls. Tests
must protect observable behavior, boundaries, focus, validation and paint
states; do not assert implementation details or duplicate coverage.

Real OS gestures, native menu appearance, GPU presentation and screen-reader
operation require separate platform checks. Record limits honestly in
`docs/components.md` and `VALIDATION.md`.

## Upstream updates

```sh
git -C references/heroui fetch --depth 1 origin v3
git -C references/heroui checkout <commit>
git -C references/heroui-native fetch --depth 1 origin main
git -C references/heroui-native checkout <commit>
git add references/heroui references/heroui-native
```

Review upstream logs and changelogs, then re-read every affected CSS/TSX file.
Do not infer current behavior from an older reference copy.

## Release

1. Review `CHANGELOG.md` and the public API. While v0, incompatible changes
   require a minor bump; patches preserve behavior.
2. Run the validation commands and inspect CI on Linux, macOS and Windows.
3. Verify `go.mod` has no replacements or workspace requirements, and that
   `git diff --check` and `go mod verify` are clean.
4. Replace Unreleased with the version/date, then tag and publish. Release
   publication is a maintainer action; local checks do not establish public
   availability.
