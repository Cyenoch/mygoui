# Proposal

## Why

Root `DESIGN.md` says it "does not establish a new visual identity" and
defers to the pinned coss reference for anything it does not record. Its
rules are spread across `DESIGN.md`, `docs/theme.md`, `docs/design.md`,
`theme/`, `internal/motion`, `internal/busy`, `overlay/motion.go` and
individual component files, and about fifty motion literals are written
nowhere else. A design decision therefore still needs the coss site, and a
later design language could not replace the current one without first
asking "what did coss do". Decision D1 of `complete-native-ui-solution`
schedules this extraction before tokens open (`open-theme-tokens`) and
before the behavior reference (`component-behavior-reference`).

## What Changes

- Extract, with Impeccable's `extract` and `document` playbooks, every
  visual and motion rule the implementation follows: colour roles and how
  derived colours are computed, type roles, spacing, radius, elevation,
  every state-motion duration and easing, overlay motion, loading, shader
  effects, control sizing, disabled and focus treatment, density and layout.
- Rewrite root `DESIGN.md` so each rule stands on its own in mygoui's terms
  and `DESIGN.md` is the visual authority. Keep the Impeccable frontmatter
  schema; refresh `.impeccable/design.json`.
- Demote coss to provenance in the records that name it as the authority:
  `docs/references.md` (licensing boundary unchanged), `docs/theme.md`,
  `docs/design.md`, `docs/components.md`, `docs/mygo-vs-css.md`,
  `PRODUCT.md` (Visual Authority, Capabilities, Native Gallery,
  Development Context), `AGENTS.md` (Design Work, Engineering),
  `docs/agent-workflow.md`, `openspec/config.yaml`, `README.md`,
  `CONTRIBUTING.md`, `docs/charts.md` and the gallery surface brief.
- Record in `design.md` the disagreements found between components for the
  same kind of change, and the motion, spacing and effect inventory that
  `open-theme-tokens` will consume, with proposed token names.
- Restate five requirements that cite coss as the authority so they cite
  the recorded design language; observable behavior is unchanged.

No Go source, value, test, file name, environment variable or licensing
statement changes. No new visual identity is introduced.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `agent-workflow`: "Design work loads durable native context" names the
  recorded design language instead of the pinned coss reference.
- `semantic-theme`: "Elevation follows the coss scale in both appearances"
  is renamed and states the shadow layers in its own terms.
- `native-library`: "Control state motion follows coss" and "Anchored
  popups share one surface and coss motion" are renamed; their motion is
  stated without reference to coss.
- `native-charts`: "Charts compose with semantic theme and utilities" uses
  the semantic chart tokens rather than "existing coss" tokens. Its Purpose
  line is edited directly in the main spec.

## Impact

- Records: `DESIGN.md`, `.impeccable/design.json`,
  `.impeccable/surfaces/internal-gallery.md`, `PRODUCT.md`, `AGENTS.md`,
  `README.md`, `docs/agent-workflow.md`, `docs/references.md`,
  `docs/theme.md`, `docs/design.md`, `docs/components.md`,
  `docs/mygo-vs-css.md`, `docs/charts.md`, `CONTRIBUTING.md`,
  `openspec/config.yaml`, four main specs through
  deltas and the `native-charts` Purpose line.
- No Go code, public API, dependency, test or gallery change. Go comments
  that say "coss" are left for `open-theme-tokens`, which edits those files.
- Follow-ups: `open-theme-tokens` consumes the inventory in `design.md`;
  the disagreements are recorded as findings for a later change to settle.
