## Context

`docs/references.md` already fixes the rule: Lobe MyGo supplies capability
and behavior, coss the appearance. None of the sixteen has a coss
counterpart, so each uses the nearest coss anatomy (button, input, card,
popover surface, empty state).

## Decisions

- **Placement.** Code and content components form a new `content` package so
  the QR dependency and syntax colouring stay out of the root import. Inputs
  join `form`, Steps `navigation`, SortableList `list`, floating surfaces and
  the image viewer `overlay`; Image, Result, GuideCard and ResizablePanel are
  root composition.
- **State.** Values, order, sizes, open flags and file lists are caller-owned
  pointers or callbacks. Components keep only transient interaction state
  (hover, drag, edit draft) in MyGo element state.
- **I/O.** `Upload` reports chosen or dropped paths; the application reads
  and uploads. `Image` takes a decoded `ui.ImageSource` and a caller status;
  the library starts no loads or goroutines.
- **Highlighter.** A small built-in tokenizer colours comments, strings,
  numbers and keywords for common languages with theme roles; it is not a
  grammar engine. Unknown languages render plain.
- **Result.** Built on `EmptyState` with `StatusIcon`, not a second layout.
- **ResizablePanel.** One component covers Lobe's DraggableSideNav and
  DraggablePanel: an edge handle, pointer and keyboard resizing, bounds and
  optional collapse, with caller-owned size.
- **Motion.** Only where it explains a change: copied feedback swap, floating
  sheet slide, viewer fade. Reduced motion settles at once.

## Risks

- QR adds a dependency to `content` importers only.
- Native file dialogs and drag-and-drop cannot be exercised headlessly beyond
  MyGo's tester hooks; real-window evidence is recorded separately.
