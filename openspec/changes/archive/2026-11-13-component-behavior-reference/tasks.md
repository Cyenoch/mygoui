# Tasks

Each family task follows `docs/behavior/README.md` and uses
`docs/behavior/form.md` as the model: read every builder's source and tests,
write one entry per component with an evidence mark on every statement, probe
doubtful behavior with a temporary `ui.NewTester` test, add `behavior_*_test.go`
tests only for documented key, focus and dismissal behavior no test protects,
never edit component code, and end with a gap list. Verify each family with
`CGO_ENABLED=0 go test ./<pkg>/... -count=1`, `go vet`, `gofmt -l`, a check
that every cited test name exists, and a check that every exported builder
has an entry.

## 1. Format and form family

- [x] 1.1 Write proposal, `behavior-reference` delta spec, design and tasks;
  verify with `npm run --silent openspec -- validate component-behavior-reference --strict --no-interactive`
- [x] 1.2 Write `docs/behavior/README.md` with the entry format, the evidence
  marks and the headless-evidence caveat; verify the template rows match the
  headings used in `form.md`
- [x] 1.3 Write `docs/behavior/form.md` covering every exported `form`
  builder, with gaps; verify every cited test exists and every builder has an
  entry or is named under "Also covers"
- [x] 1.4 Add `form/behavior_keyboard_test.go` and
  `form/behavior_dismissal_test.go`; verify `CGO_ENABLED=0 go test ./form/... -count=1`,
  `CGO_ENABLED=0 go vet ./form/...` and `gofmt -l form`
- [x] 1.5 Link the reference from `docs/components.md` and the PRODUCT.md
  pillar table; verify the links resolve

## 2. Root family

- [x] 2.1 Write `docs/behavior/root.md` for the `mygoui` package (buttons,
  toggles, menus, tabs, dialogs, toasts, workspace components and the rest of
  its 72 builders), with tests and gaps as above; mark it in the README table

## 3. Navigation family

- [x] 3.1 Write `docs/behavior/navigation.md` for `navigation`, with tests
  and gaps as above; mark it in the README table

## 4. Overlay family

- [x] 4.1 Write `docs/behavior/overlay.md` for `overlay`, including focus on
  open and return on close for every popup and flyout, with tests and gaps as
  above; mark it in the README table

## 5. List and data family

- [x] 5.1 Write `docs/behavior/list-data.md` for `list` and `data`, including
  ListBox navigation, typeahead and active-descendant exposure, with tests and
  gaps as above; mark it in the README table

## 6. Date and color family

- [x] 6.1 Write `docs/behavior/date-color.md` for `date` and `color`, with
  tests and gaps as above; mark it in the README table

## 7. Content family

- [x] 7.1 Write `docs/behavior/content.md` for `content`, with tests and gaps
  as above; mark it in the README table

## 8. Chart family

- [x] 8.1 Write `docs/behavior/chart.md` for `chart` (keyboard exploration,
  tooltips, what assistive technology reads), with tests and gaps as above;
  mark it in the README table

## 9. Close

- [x] 9.1 Verify every exported builder in the repository has an entry, run
  `make check` and `make spec-check`, then verify and archive the change

## Workflow follow-up

- Gaps listed in each family document become their own changes, one per
  component, as `complete-native-ui-solution` D5 requires.
