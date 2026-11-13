# Proposal

## Why

mygoui is recorded as "a native rebuild of the coss ui catalog". An application
author, or the agent writing their application, gets components and a colour
palette, but no public motion or effect tokens, no per-component behavior
contract, and no guidance written for them: every skill in the repository
serves maintainers. The project now aims to be a complete native UI solution,
and its product record has to say so before the work is split up.

## What Changes

- Restate the product in `PRODUCT.md` as four pillars with one owner each:
  theme tokens, the component library, a design and behavior reference, and
  agent skills for both application authors and maintainers.
- Record the direction for the visual authority: the coss design language is
  extracted into the project's own record with Impeccable, coss becomes
  provenance, and a later design language can replace it without touching
  component code. Until that extraction lands, the pinned coss reference
  remains the authority and no contract changes.
- Align the one-line descriptions in `README.md` and `AGENTS.md` with that
  product statement.
- Fix the order and boundaries of the follow-up changes (see `design.md`):
  `publish-agent-skills`, `extract-design-language`, `open-theme-tokens`,
  `component-behavior-reference`, `native-effects`.

This change edits product records only. It adds no API, moves no token and
writes no skill; each follow-up change carries its own requirements.

## Capabilities

### New Capabilities

None. The change sets `skip_specs: true`: it alters no observable behavior.

### Modified Capabilities

None. The requirements that name coss as the authority
(`agent-workflow`: "Design work loads durable native context";
`semantic-theme`: "Elevation follows the coss scale in both appearances";
`native-library`: "Control state motion follows coss" and "Anchored popups
share one surface and coss motion"; `native-charts`: "Charts compose with
semantic theme and utilities") stay as they are. `extract-design-language`
owns their deltas.

## Impact

- Records: `PRODUCT.md`, the opening of `README.md`, the first line of
  `AGENTS.md`.
- No Go code, public API, dependency, test or gallery change.
- `openspec/config.yaml`, `DESIGN.md`, `docs/design.md` and
  `docs/references.md` still describe the current coss authority correctly
  and are left alone until `extract-design-language`.
