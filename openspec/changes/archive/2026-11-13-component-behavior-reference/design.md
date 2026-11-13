# Design

## Context

`complete-native-ui-solution` (D5) made behavior a checkable contract per
component, drawn from the native implementation, platform conventions and the
WAI-ARIA Authoring Practices (APG), and left open whether it lives in OpenSpec
per family or under `docs/`. About 230 public builders exist across nine
packages; `form` alone has 39.

## Decisions

### D1. The reference lives in `docs/behavior/`; OpenSpec owns only its rules

A per-key OpenSpec requirement would put hundreds of near-identical
scenarios into main specs, duplicate the tests that already protect them, and
force a spec delta for every key a component gains. The reference instead
lives as one Markdown document per family, and the new `behavior-reference`
capability states what the reference must be: complete, evidenced, and
updated with behavior changes. A key binding becomes binding by being in the
reference with a test mark, not by being a requirement.

Alternative considered: one capability per family with a requirement per
component. Rejected for the volume above, and because OpenSpec scenarios
cannot point at the test that proves them.

### D2. One document per package family

`root`, `form`, `navigation`, `overlay`, `list-data`, `date-color`,
`content`, `chart`. The split follows the import paths an author already
knows; `list`/`data` and `date`/`color` are paired because each pair is small
and shares widgets. One file per component would scatter the gap lists and
make the shared format drift.

### D3. A fixed, compact entry

Each entry is a heading plus one two-column table whose rows always appear
in the same order: Role and name, Keys, Pointer, Focus, Dismissal, Selection
and value, States, Motion. Non-applicable rows are omitted. Motion rows name
what moves and what is instant and link to `DESIGN.md` for durations and
easing, so the numbers have one owner.

### D4. Three evidence marks

`[T]` an existing test, named; `[N]` a test this change added, named; `[C]`
observed in the code only. A reader can tell at a glance which promises a
regression would break a test, and later work can turn `[C]` into `[T]`.
Headless `ui.NewTester` results are accessibility-tree and focus evidence
only; the README says so once.

### D5. Tests only where consumers depend on the behavior

New tests cover documented key, focus and dismissal behavior that no test
protected and that an application would notice breaking. They live in
`behavior_*_test.go` inside the family package and reuse its helpers. They
never edit component code: a failing probe becomes a corrected entry or a gap.

### D6. Gaps are measured, not fixed

Each family document ends with a gap list, one line per difference from a
platform convention or the APG pattern. Closing a gap is a separate change
per component, as D5 of the umbrella change requires.

## Form findings recorded in the gap list

Probes run while writing `form` found two defects and one interaction worth
recording: Select's Escape closes the popup without returning focus to the
trigger; removing a KeyValueEditor row with its button leaves the focus
nowhere; and an ancestor Escape shortcut takes Escape before SearchField can
clear its text. Each is in `docs/behavior/form.md` under Gaps.

## Risks

- Entries drift from code. Mitigation: the `behavior-reference` requirement
  that behavior changes update their entry, and test marks that name tests.
- Later families read `[C]` as weaker than it is. The README defines the
  marks as provenance, not confidence.
