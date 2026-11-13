## Why

The behavior reference (`docs/behavior/`) and the design-language extraction
(findings F1–F17 in the archived `extract-design-language` design) recorded
defects and inconsistencies without fixing them: focus lost after dismissal,
keys that pass bounds, values clamped per key, one motion expressed with
several timings. A reference that lists known defects is not a contract an
application can rely on.

## What Changes

- Fix every defect listed in the gap lists of `docs/behavior/*.md`, each with
  a test, and move the row from the gap list to the entry with its evidence.
- Settle design findings F1–F17: one treatment per motion or effect, read
  from theme tokens, recorded in `DESIGN.md`.
- A gap that the toolkit owns and the library cannot fix stays, marked as a
  toolkit limit with the reason.
- Remove waste found on the way (duplicate helpers, per-frame work at rest).

## Capabilities

### New Capabilities

None.

### Modified Capabilities

None: the requirements in `behavior-reference`, `native-library` and
`semantic-theme` already demand this behavior; this change makes the
implementation and the reference meet them.

## Impact

Component packages (root, form, list, data, navigation, overlay, date, color,
chart, content, layout), `docs/behavior/`, `DESIGN.md`, `CHANGELOG.md`, and
the `mygoui-*` skills where they quote changed behavior.
