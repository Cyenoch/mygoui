# Maintaining mygoui design and behavior

For work inside the mygoui repository. Paths are relative to its root.

## Where rules live

| Rule | Owner |
| --- | --- |
| Product facts, pillars, binding constraints | `PRODUCT.md` |
| Visual rules, state motion, loading, default values | `DESIGN.md` (with `.impeccable/design.json` for native extensions) |
| Native ownership and engineering contracts | `docs/design.md` |
| Component subsets, adaptations and gaps | `docs/components.md` |
| CSS-to-MyGo translation | `docs/mygo-vs-css.md` |
| Observable requirements | `openspec/specs/*/spec.md` |
| Evidence actually obtained | `VALIDATION.md` |

`SKILL.md` restates a rule only as an instruction an application author can act on and points at the owner for the rest. A new rule goes into its owning record first, through an OpenSpec change; then the skill points at it.

## Workflow

- Design work uses the `impeccable` skill: `make design-context TARGET=<path>` once per session, then the playbook the finding needs (`docs/agent-workflow.md` has the table). The gallery is an Operate surface; its brief is under `.impeccable/surfaces/`.
- Web detectors and browser live mode do not analyze Go UI. Use `ui.NewTester`, `make design-artifacts` and real windows; record OS, GPU and screen-reader observations separately in `VALIDATION.md`.
- Behavior checks follow `mygo-maintenance` and the existing tests' patterns: drive input, read `Access`, advance the clock, assert `FramePending`.

## Keep this skill current

When `DESIGN.md`, `docs/design.md` or a behavior requirement changes a rule this skill states (loading kinds, accessible-name sources, dismissal, copy, the Tester boundary), update `SKILL.md` in the same change and list it in that change's tasks. The planned `extract-design-language` and `component-behavior-reference` changes carry that task (decision D2 of `complete-native-ui-solution`). Rerun the example test against the checkout when the APIs it uses change.
