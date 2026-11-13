# Proposal

## Why

An agent building a MyGo application with mygoui has no skill to load: every
skill in this repository serves maintainers, and the mygoui-specific rules for
components, theme, design and behavior, motion and shader effects reach an
agent only by reading the library's records cold. `complete-native-ui-solution`
makes agent skills a product pillar and schedules this change first (D2, D3).

## What Changes

- Add five model-invoked skills, `mygoui-components`, `mygoui-theme`,
  `mygoui-design`, `mygoui-motion` and `mygoui-effects`, each serving both
  audiences from one source: `SKILL.md` states the public contract for an
  application author's agent; a sibling `MAINTAINING.md`, reached by one
  pointer, adds the maintainer-only contracts and workflow.
- Publish them from a root `skills/` directory that the `skills` CLI
  discovers, so `npx skills add Cyenoch/mygoui` installs them once pushed;
  link each into `.agents/skills/` so the project's harnesses read them
  beside the generated upstream skills.
- Point the skills at the owning records shipped with the module
  (`DESIGN.md`, `docs/*.md`) and at `go doc` instead of copying values.
- Record installation in `README.md`, the layout in `docs/agent-workflow.md`
  and `AGENTS.md`, and the delivered pillar in `PRODUCT.md`.

## Capabilities

### New Capabilities

- `agent-skills`: published mygoui skills: installable from the repository,
  linked into the project's skill directory, public contract for application authors
  with disclosed maintainer material, and names that exist in the API.

### Modified Capabilities

None. `agent-workflow` keeps its requirements: the harnesses still resolve
skills through `.agents/skills`, and setup and refresh behave as specified.
The published-skill requirements belong to the new capability.

## Impact

- New: `skills/mygoui-*/`, five links under `.agents/skills/`.
- Records: `README.md`, `docs/agent-workflow.md`, `AGENTS.md` (Skills),
  `PRODUCT.md` (Agent skills row).
- No Go code, public API, dependency or gallery change. The skills describe
  the system as it is; `extract-design-language`, `open-theme-tokens`,
  `component-behavior-reference` and `native-effects` update the affected
  skill in their own tasks.
