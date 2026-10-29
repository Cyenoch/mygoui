# OpenSpec and Impeccable

OpenSpec owns behavior contracts and scoped changes. Impeccable owns design
discovery, review and refinement. MyGo maintenance and the Go skills supply
native implementation guidance. This is the project's working agreement,
replacing the Matt Pocock skill workflow.

## Setup

Use Node.js 20.19 or newer for the development CLIs:

```sh
make agents
make spec-check
make design-doctor
```

`package.json` and `package-lock.json` pin the CLI dependencies.
The generated skills, agent profiles and both harnesses' hook manifests are
committed. Setup verifies the npm platform engine against the skill's
committed version and SHA256 manifest, then copies it into the launcher's
ignored project binary directory. Commands and hooks use this binary ahead
of an unversioned home installation. Run setup before using design commands.
Node and these tools are not required to import or build the Go library.

OpenSpec's twelve workflows are installed as skills in `.agents/skills/`.
The project uses skills-only delivery, so setup does not write prompts to
your home directory. Both npm and Make call `scripts/openspec.mjs`, which
scopes `XDG_CONFIG_HOME` to `.agents/config/`; its checked-in profile selects
the complete workflow set without changing your global preferences.
Use `npm run openspec -- <arguments>` for the bare `openspec` examples in
upstream skills so they use this project configuration.

Claude and Codex share physical files:

| Entry | Canonical source |
| --- | --- |
| Codex skills | `.agents/skills/` |
| `.claude/skills` | Symlink to `../.agents/skills` |
| `.claude/agents` | Symlink to `../.agents/claude-agents` (official Claude profiles) |
| `.codex/agents` | Symlink to `../.agents/skills/impeccable/agents` (official Codex profiles) |
| `CLAUDE.md` | Imports `AGENTS.md`; Codex reads `AGENTS.md` directly |

Edit the canonical source; preserve the links. Setup creates missing links
and refuses conflicting directories without replacing user content.
Windows checkouts need directory-symlink support and Git symlink preservation;
plain files containing a link target cannot supply shared skills.

Use `make agents-update` to intentionally refresh generated integrations.
Review the resulting diff: Impeccable's updater fetches the current skill
distribution, independently of the pinned npm CLI. It refreshes the Claude
profiles, then reinstalls the Codex skill variant as the canonical shared
version; this variant includes both harnesses' invocation guidance.
If the refreshed engine version changes, align the pinned npm dependencies
before setup can pass its checksum verification.

Both `.claude/settings.json` and `.codex/hooks.json` wire `PostToolUse` and
`Stop` to the shared launcher. Claude's project-local settings remain ignored.
Codex requires a separate trust decision in `/hooks` before executing them;
Claude also applies its own project trust. Enabling a hook in project
configuration does not establish that the harness has trusted it.
The hook scans supported web files; it does not scan Go UI code.

## Sources of truth

| Record | Responsibility |
| --- | --- |
| `PRODUCT.md` | Users, purpose, platform, constraints and binding references |
| Root `DESIGN.md` | Current visual system, portable default tokens and application guidance |
| `.impeccable/design.json` | Visual extensions such as native motion and compact layout threshold |
| `.impeccable/surfaces/` | Surface-specific audience, mode and workflow |
| `openspec/specs/` | Current observable behavior contracts |
| `openspec/changes/<name>/` | A scoped proposal, delta requirements, decisions, tasks and evidence |
| `docs/design.md` | Native ownership and public engineering contracts |
| `docs/components.md` | Catalog correspondence, implemented subsets and gaps |
| `docs/theme.md`, `theme/` | Token provenance, customization and executable defaults |
| `VALIDATION.md` | Commands actually run, outcomes and platform evidence limits |

The initial native-library, semantic-theme and gallery-workbench specs
summarize existing documented behavior, including the current working tree.
They do not certify new platform results or replace the detailed capability
matrix. Keep these summaries aligned as the implementation changes.

## Choose the amount of process

New capabilities, public behavior/API changes, architectural changes and
substantial fixes need a change. Read the relevant main specs, inspect the
code, then capture only the behavior and decisions the work changes.
Small corrections that preserve the contract can proceed directly with the
appropriate focused check. Do not create a spec delta for a pure cleanup.
An artifact-only change with no behavior delta can use `skip_specs: true`
in its `.openspec.yaml`.

Explicit requests to explore, shape or propose are planning work. A request
to implement authorizes work within that scope. Material ambiguity still
needs an answer; routine implementation choices do not need new approval
rounds. Keep one set of artifacts instead of parallel PRDs and plans.
Test meaningful regressions and boundaries; a mandatory red/green ritual
is not the project workflow.

## OpenSpec workflows

These are agent skills, not terminal subcommands. Codex uses `$skill-name`;
Claude uses `/skill-name`. This shared skills-only setup does not install
the separate `/opsx:` command aliases.

| Intent | Codex | Claude |
| --- | --- | --- |
| Investigate an uncertain problem | `$openspec-explore` | `/openspec-explore` |
| Produce a complete proposal | `$openspec-propose` | `/openspec-propose` |
| Start with a scaffold | `$openspec-new-change` | `/openspec-new-change` |
| Build artifacts incrementally | `$openspec-continue-change` | `/openspec-continue-change` |
| Generate the planning artifacts together | `$openspec-ff-change` | `/openspec-ff-change` |
| Implement tracked tasks | `$openspec-apply-change` | `/openspec-apply-change` |
| Revise an existing change as understanding improves | `$openspec-update-change` | `/openspec-update-change` |
| Check completeness, correctness and consistency | `$openspec-verify-change` | `/openspec-verify-change` |
| Merge deltas into current specs before archive | `$openspec-sync-specs` | `/openspec-sync-specs` |
| Complete and archive one change | `$openspec-archive-change` | `/openspec-archive-change` |
| Archive several completed changes with conflict review | `$openspec-bulk-archive-change` | `/openspec-bulk-archive-change` |
| Learn the workflow through a real scoped task | `$openspec-onboard` | `/openspec-onboard` |

A typical substantive change uses propose, apply, verify and archive.
Explore first when the problem is uncertain; use new/continue instead of
propose when individual artifacts need discussion. Update artifacts when
implementation reveals a better answer. Sync is useful before archiving;
archive also handles outstanding delta synchronization.

Agents execute the pinned CLI with:

```sh
npm run openspec -- list --specs
npm run openspec -- list
npm run --silent openspec -- status --change <name> --json
npm run --silent openspec -- instructions apply --change <name> --json
make spec-check
```

Use npm's `--silent` flag when parsing JSON output so npm does not prepend
its script banner.

`openspec/config.yaml` injects native constraints, per-artifact rules and
apply/archive guidance. Keep project customization there and in `AGENTS.md`;
regenerate upstream skills rather than forking their contents.
The standard `spec-driven` schema already provides proposal, specs, design
and tasks; no extra schema or second task tracker is needed.

Before archive, verify the implementation against the scenarios, resolve
review findings, and record the evidence. Structural validation checks the
artifact format; it does not run the application or prove requirements.
After archive, main specs describe current behavior and the archive retains
the completed change and its evidence.

## Impeccable on native desktop

Use the project skill and `make design-context TARGET=internal/gallery` once
at the start of design work. The target selects its surface brief.
Impeccable's `adaptive` selector enables native guidance. Shipping surfaces
are macOS, Windows and Linux; phone touch targets, mobile navigation and
browser conventions do not override the documented desktop controls.
The gallery is an Operate surface; documentation is Read.

For new flows, use `$impeccable shape <target>`. For existing interfaces,
start with `$impeccable critique <target>` for UX or
`$impeccable audit <target>` for technical quality. Put actionable findings,
decisions and checks in the active change. Preserve the coss system during
refinement; a redesign requires a redesign request. In Claude, use
`/impeccable` with the same arguments.

Choose the playbook that answers the finding:

| Work | Impeccable commands |
| --- | --- |
| Structure, typography and density | `layout`, `typeset`, `distill` |
| Visual emphasis and color | `bolder`, `quieter`, `colorize` |
| Interaction and expression | `animate`, `delight`, `overdrive` |
| Labels, recovery and first use | `clarify`, `harden`, `onboard` |
| Window sizes, text scale and runtime cost | `adapt`, `optimize` |
| Final inspection and consistency | `polish` |
| Product context and visual system maintenance | `init`, `document`, `extract`, `doctor` |

Load the command's reference and craft floor before editing UI. Reuse
existing components and `theme` APIs. `extract` needs real repetition;
`bolder`, `colorize` and `overdrive` must still honor the visual reference,
native behavior and reduced motion. Keep unrequested aesthetics out of a
focused fix.

`live` and `generate` require a browser-backed implementation and are
available in the installed skill for a future web surface. They cannot
select or mutate native MyGo controls. Likewise, `impeccable detect` parses
web markup/styles/scripts, not Go drawing or native accessibility.
Do not map `.go` to its text detector or turn unsupported scans into a
passing CI check.

Use the existing native verification instead:

```sh
make design-artifacts
make check
```

`design-artifacts` runs the existing gallery collector and writes
compact/wide, light/dark specimens under `artifacts/design-review/`.
Inspect those images and the affected controls' interaction states.
For the full component/state collector, use
`MYGOUI_COSS_VISUAL_AUDIT_DIR` with `TestCossVisualAudit` as described in
the existing gallery tests and `docs/visual-review.md`.
Actual window behavior can use the `e2e` gallery automation documented in
`docs/mygo-0.7.md`.

Apply `ui.NewTester` to changed keyboard/pointer actions, focus, disabled,
read-only, invalid and retained-exit behavior as relevant. Check actual OS
windows for window chrome, system menus, IME, GPU and assistive-technology
claims. Record what was and was not observed.

`make design-doctor` checks context/configuration health; read its findings
array even on exit zero. Use `document` when visual rules change and `init`
when product facts change. Review upstream updates, notices and generated
files together; local configuration and engine binaries stay ignored.

The doctor can report an unset `buildPath` after a surface brief is
recorded. That preference affects future Persuade/Experience surfaces;
Operate and Read remain code-led. It is intentionally unset here until
such a surface needs a user decision, rather than inventing a standing
mockup-first or code-first preference during native workflow setup.
