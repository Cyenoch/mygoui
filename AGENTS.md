# mygoui

A generic native MyGo component library with a shadcn-inspired style. Read `docs/design.md` for ownership, `docs/components.md` for capabilities and `docs/mygo-0.7.md` before dependency work. Public APIs live at the module root and in coherent family packages. `internal/gallery` and `examples/` are consumers.

## Working Agreement

- OpenSpec owns requirements and change tracking; Impeccable owns design work.
  Read [the workflow](docs/agent-workflow.md) and the relevant
  `openspec/specs/*/spec.md` before substantive changes.
- New capabilities, public behavior/API changes, architectural changes and
  substantial fixes use `openspec/changes/<name>/`. Keep proposal, delta specs,
  design decisions and tasks aligned with implementation; verify and archive
  completed work. Small corrections that preserve the contract need no new change.
- Explicit planning requests stop at planning. Implementation requests authorize
  implementation within their stated scope; do not add approval rounds that the
  user did not request. Ask only about material unresolved decisions.
- This replaces Matt Pocock's skill workflow for this project. Do not invoke
  `grill-me`, `write-a-prd`, `prd-to-plan`, `prd-to-issues`,
  `design-an-interface`, `improve-codebase-architecture`, or a mandatory TDD
  ritual as a parallel process. Preserve focused regression testing.
- Preserve existing uncommitted work. Do not create a branch unless requested.

## Design Work

- Use the project `.agents/skills/impeccable/SKILL.md`. Load its context once
  per session with `make design-context` (add `TARGET=path` when scoped).
- `PRODUCT.md` records product truth; root `DESIGN.md` records visual decisions
  extracted from the implementation. `docs/design.md` owns native engineering
  contracts. Surface briefs live under `.impeccable/surfaces/`.
- This is an existing native desktop system. Preserve coss authority and use
  Impeccable's Operate mode for the gallery. Desktop MyGo contracts override
  phone-only examples in its adaptive platform references.
- Route design work through the relevant Impeccable playbook: shape/critique,
  audit, focused refinements, then polish. Read its craft floor before UI edits.
  Record actionable findings and their verification in the active OpenSpec
  change; use existing tokens and components.
- The web detector and browser live/generate do not analyze Go/MyGo UI.
  Use `ui.NewTester`, native captures and real-window checks for those surfaces.
  A detector run over zero supported files is not design validation.

## Engineering

- Implement with typed Go options, explicit errors, caller-owned state and native MyGo controls. Visual references establish appearance; each layer uses its own idioms.
- Keep `ui.Element`, custom parts and `*ui.Context` build-scoped. Retain `ui.Handle` and widget state per window. Use native services and `Window.Update` for background results.
- Reuse the standard library, MyGo and mygoicons before adding infrastructure. Preserve native focus, selection, editing and accessibility.
- Update the corresponding design/capability documentation when public behavior changes. Tests assert consumer-visible behavior, boundaries and real regressions.
- A library import must remain independent of application initialization, application language registration, I/O, persistent storage and background jobs.
- Keep release dependencies resolvable: module metadata contains no local replacements. Use a private development workspace for local dependencies.

## Skills

Read `.agents/skills/mygo-maintenance/SKILL.md` for UI changes and upgrades. Use the copied Go skills for code, errors, naming, tests, dependencies, lint and CI as relevant. Project rules and the user's authorized scope take precedence over generic skill examples.

All twelve OpenSpec workflows are installed as project skills, including
verify, sync and bulk archive. Run the pinned CLI with
`npm run openspec -- <arguments>`; run `make agents` after cloning and
`make agents-update` when intentionally refreshing generated integrations.
Use this project entry point instead of a bare `openspec` command in generated
examples: it applies the shared profile on both Claude and Codex.
`.claude/skills` links to `.agents/skills`; `.claude/agents` links to
`.agents/claude-agents`. Codex agent profiles link from `.codex/agents` to
the upstream Impeccable skill's `agents` directory.
Edit the canonical targets, preserving the links.
Claude loads these rules through `CLAUDE.md`; Codex reads this file directly.
Do not edit generated upstream skills for project-specific rules; keep those
in this file and `openspec/config.yaml`.

## Validation

`make check` checks the headless suite, build, Go vet, MyGo lifetimes, formatting and module integrity. Also run `make race` for shared-state changes, `golangci-lint run ./...`, and the WebAssembly build before release. Exercise changed UI behavior with `ui.NewTester`. Record OS/GPU/screen-reader evidence separately from headless rendering. Read `CONTRIBUTING.md` before publishing.

`make spec-check` validates main specs, active deltas and archived task completion.
`make design-doctor` reports Impeccable context/configuration drift; inspect its
findings, since a zero exit code alone does not mean the report is clean.
