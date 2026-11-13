# Proposal

## Why

An application author, or the agent writing their application, can see what a
mygoui component looks like (`DESIGN.md`) and what it can do
(`docs/components.md`), but not how it behaves: which keys it answers, where
focus goes when it opens, closes or loses a row, what dismisses it, what
Loading, Disabled and ReadOnly change, and what assistive technology is told.
That knowledge lives only in the source and its tests, so an author cannot
rely on it and a maintainer cannot tell a regression from an undocumented
choice. The umbrella change `complete-native-ui-solution` (D5) named this the
third pillar and left one question open: requirements per family in OpenSpec,
or a reference under `docs/`.

## What Changes

- Add `docs/behavior/`, one document per package family (root, form,
  navigation, overlay, list-data, date-color, content, chart) and a README
  that fixes the entry format and the evidence marks.
- Each public component gets one entry with the same headings in the same
  order. Every statement says how it is known: protected by a named test,
  protected by a test this change adds, or observed in the code only.
  Toolkit-owned behavior (caret, selection, IME, the radio group's roving
  focus) is named once and left to MyGo.
- Each family document ends with a gap list measured against platform
  conventions and the WAI-ARIA Authoring Practices.
- Deliver the `form` family complete as the exemplar, with headless tests
  for the documented key, focus and dismissal behavior no existing test
  protected. Component code does not change; defects found go to the gap
  list.
- Link the reference from `docs/components.md` and from the "Design and
  behavior reference" pillar in `PRODUCT.md`.
- The other seven families are tasks for later sessions.

## Capabilities

### New Capabilities

- `behavior-reference`: the reference itself — every public component has an
  entry, statements carry evidence, and a behavior change updates its entry.
  Individual key bindings are not requirements: they live in the reference,
  under the evidence that protects them.

### Modified Capabilities

None.

## Impact

- New: `docs/behavior/README.md`, `docs/behavior/form.md`,
  `form/behavior_keyboard_test.go`, `form/behavior_dismissal_test.go`.
- One sentence each in `docs/components.md` and the `PRODUCT.md` pillar
  table.
- No public API, no component code and no main spec changes.
