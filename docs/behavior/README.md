# Behavior reference

How each public mygoui component behaves: keys, pointer, focus, dismissal,
selection, states, motion and what assistive technology is told. Visual rules
live in [`DESIGN.md`](../../DESIGN.md); capabilities in
[`components.md`](../components.md). The rules this reference follows are the
`behavior-reference` capability in OpenSpec.

| Family | Packages | Document |
| --- | --- | --- |
| Root | `mygoui` | [root.md](root.md) |
| Form | `form` | [form.md](form.md) |
| Navigation | `navigation` | [navigation.md](navigation.md) |
| Overlay | `overlay` | [overlay.md](overlay.md) |
| Layout | `layout` | [layout.md](layout.md) |
| List and data | `list`, `data` | [list-data.md](list-data.md) |
| Date and color | `date`, `color` | [date-color.md](date-color.md) |
| Content | `content` | [content.md](content.md) |
| Chart | `chart` | [chart.md](chart.md) |
| Effect | `effect` | [effect.md](effect.md) |

## Evidence

Every statement ends with one mark.

| Mark | Meaning |
| --- | --- |
| `[T Name]` | An existing test named `Name` asserts it. |
| `[N Name]` | A test added for this reference asserts it. |
| `[C]` | Read from the source; no test asserts it. |

The marks record provenance, not confidence: `[C]` is what the code does
today, and a change may alter it without failing a test.

Tests run headless through `ui.NewTester`. They prove the accessibility tree,
focus and key handling MyGo builds, not what a screen reader announces, how an
operating system routes keys, or what a GPU draws. Evidence of those is
recorded separately and named as such.

The MyGo toolkit owns the caret, text selection, clipboard shortcuts, input
methods, spell checking, the radio group's roving focus and slider key
stepping. Entries say "toolkit" for these and do not list them.

## Entry format

One entry per component, `###` heading, then one table. Rows appear in this
order; a row that does not apply is left out.

```markdown
### ComponentName

Also covers: HelperA, HelperB.

| | |
| --- | --- |
| Role and name | Role, where the accessible name comes from, states exposed. [mark] |
| Keys | Each key and what it does while focused. [mark] |
| Pointer | Clicks, drags, drops, hover. [mark] |
| Focus | Tab stops; where focus goes on open, close, remove. [mark] |
| Dismissal | What closes, cancels or clears it, and what it leaves. [mark] |
| Selection and value | Who owns the value; when it is reported. [mark] |
| States | Disabled, ReadOnly, Invalid, Loading. [mark] |
| Motion | What moves, what is instant, reduced motion. Link DESIGN.md. [mark] |
```

"Also covers" is used only when helper builders share the entry.

Durations and easing are owned by `DESIGN.md` ("State motion", "Loading");
entries say what moves and link there.

Each family document ends with **Gaps**: one line per difference from a
platform convention or the
[WAI-ARIA Authoring Practices](https://www.w3.org/WAI/ARIA/apg/patterns/),
including defects found while writing it. "Probed" marks a gap a
temporary headless test confirmed; no test is kept that asserts a defect.
Closing a gap is its own change.
