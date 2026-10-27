# Contributing

The `coss` branch reconstructs coss ui with native Go/MyGo controls. Read
[design](docs/design.md), the [component matrix](docs/components.md),
[references](docs/references.md) and [MyGo 0.7 guidance](docs/mygo-0.7.md)
before changing public behavior. The current coss contract supersedes the
older visual-reference wording in the MyGo migration notes.

Use the module's Go 1.27.3 target and pinned MyGo v0.7.1/mygoicons v0.2.1.
Published module metadata must have no local replacements. Use a private
development workspace when testing local dependencies.

## References and licensing

Use coss commit `3d8ca28e7d90123e9e0f44665944dc240b026ccb`.
Run `make refs` to populate the pinned `references/coss` submodule.
Reference sources are development inputs, not runtime dependencies.
[docs/references.md](docs/references.md) identifies the source paths and
catalog discrepancies.

Read relevant `apps/ui/registry/default/ui/*.tsx`, MDX documentation and
examples. That subtree is declared MIT by upstream. The monorepo default
and `packages/*` are AGPL; do not copy their implementation or font assets
under the MIT attribution. Use the MIT styling documentation for token
provenance and record any observed runtime differences separately.

Fonts are an application responsibility. Independently licensed Cal Sans v2
and Paper Mono sources are listed in [docs/fonts.md](docs/fonts.md).
The library must not register fonts on import. Preserve the original OFL
notices if an application elects to distribute those fonts.

## Implementation

- Reuse standard Go, MyGo and mygoicons before adding infrastructure.
  Prefer typed options and explicit window-owned state.
- Keep `ui.Element`, custom parts and `*ui.Context` build-scoped. Use stable
  keys and handles; capture services for persistent native callbacks.
- Use native bases for editing, focus, selection and accessibility. Paint
  callbacks must not retain or query expired elements.
- Use `theme.Colors(c)` and canonical coss roles. `Primary` and `Accent`,
  and `MutedBackground` and `MutedForeground`, are separate roles.
- Set layout alignment explicitly. `ui.Start` is a real zero value;
  `ui.Row` cross-axis centering differs from CSS.
- Keep UI labels concise. Document toolkit limits instead of implying that
  every upstream part, CSS state or gesture exists.
- Update capability, customization and provenance documentation alongside
  public API changes. Legacy extra families need no invented coss equivalent.

## Validate

Run the checks appropriate to the change. Before release:

```sh
make check
make race
golangci-lint run ./...
make wasm
make workflows
```

`make check` runs headless tests, ordinary build, Go vet, MyGo lifetime vet,
formatting and module verification. Race checks matter for shared-state
changes. `make race` runs test binaries serially because each headless
rasterizer already uses multiple cores; package-local concurrency stays
enabled. Inspect module integrity when changing dependencies; run
`go mod tidy` only as part of intentional module work.

Exercise changed controls with `ui.NewTester`: keyboard/pointer activation,
focus, read-only/disabled/invalid states, relevant light/dark paint and
retained exits. Tests should protect a current consumer-visible contract.
Do not add assertions of deleted APIs or arbitrary visual/count gates.

Record commands, result and revision in [VALIDATION.md](VALIDATION.md).
Keep pending checks pending. Gallery page presence or headless images do not
establish coss parity, GPU correctness, OS menu appearance, IME/gesture
behavior or actual screen-reader operation.

## Release

1. Review the public API, capability matrix, notices and `CHANGELOG.md`.
   While v0, incompatible changes require a minor bump.
2. Complete the relevant validation and inspect Linux/macOS/Windows CI.
3. Check resolvable module metadata, formatting and `git diff --check`.
   Verify copied notices match the exact distributed reference/assets.
4. Record the version/date before the maintainer tags and publishes.
   Local checks alone do not establish public availability.
