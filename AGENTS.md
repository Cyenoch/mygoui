# mygoui

A generic native MyGo component library with a shadcn-inspired style. Read `docs/design.md` for ownership, `docs/components.md` for capabilities and `docs/mygo-0.7.md` before dependency work. Public APIs live at the module root and in coherent family packages. `internal/gallery` and `examples/` are consumers.

## Engineering

- Implement with typed Go options, explicit errors, caller-owned state and native MyGo controls. Visual references establish appearance; each layer uses its own idioms.
- Keep `ui.Element`, custom parts and `*ui.Context` build-scoped. Retain `ui.Handle` and widget state per window. Use native services and `Window.Update` for background results.
- Reuse the standard library, MyGo and mygoicons before adding infrastructure. Preserve native focus, selection, editing and accessibility.
- Update the corresponding design/capability documentation when public behavior changes. Tests assert consumer-visible behavior, boundaries and real regressions.
- A library import must remain independent of application initialization, application language registration, I/O, persistent storage and background jobs.
- Keep release dependencies resolvable: module metadata contains no local replacements. Use a private development workspace for local dependencies.

## Skills

Read `.agents/skills/mygo-maintenance/SKILL.md` for UI changes and upgrades. Use the copied Go skills for code, errors, naming, tests, dependencies, lint and CI as relevant. Project rules and the user's authorized scope take precedence over generic skill examples.

## Validation

`make check` checks the headless suite, build, Go vet, MyGo lifetimes, formatting and module integrity. Also run `make race` for shared-state changes, `golangci-lint run ./...`, and the WebAssembly build before release. Exercise changed UI behavior with `ui.NewTester`. Record OS/GPU/screen-reader evidence separately from headless rendering. Read `CONTRIBUTING.md` before publishing.
