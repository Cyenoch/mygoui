# Maintenance facts

- MyGo 0.6 applies bound input before constructors return; deferred callbacks run after configuration and before rebuilding. The pinned `docs/ui/migration.md` supersedes older descriptions of next-pass-only notices.
- `Tester.Image()` performs image acquisition work; calling it for each pixel can amplify race-instrumented tests by orders of magnitude. Acquire one frame and inspect its pixels. `conversation/surfaces_bubble_test.go` went from a five-minute race timeout to a complete package race run in about two seconds after this change.
- `Tester.Find` matches a row's accessible label before a cell's identical text when the row appears first. Clicking that label targets the whole row's center, possibly a different column. Give row labels distinct accessible names in cell-selection tests; `data/cell_selection_test.go` exercises actual row-header and rectangular selection.
