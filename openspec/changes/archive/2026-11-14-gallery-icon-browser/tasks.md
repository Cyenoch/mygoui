# Tasks

## 1. Planning

- [x] 1.1 Write proposal, delta spec, design and tasks; verify with
  `npm run --silent openspec -- validate gallery-icon-browser --strict --no-interactive`

## 2. Dependency

- [x] 2.1 Require mygoicons v0.3.0 and update the pinned version in
  `README.md`, `CONTRIBUTING.md` and `THIRD_PARTY_NOTICES.md`

## 3. Gallery

- [x] 3.1 Add the Icons page: search by Go name or
  Iconify ID, virtual rows, alias resolution, copy on activation with an
  announcement, empty state
- [x] 3.2 Test the three spec scenarios and row virtualization with
  `ui.NewTester`
- [x] 3.3 Measure the gallery executable and WebAssembly size with and
  without Twemoji and record the decision in `design.md`
- [x] 3.4 Use icons in the gallery's navigation groups and page actions

## 4. Verification

- [x] 4.1 `make generate`, `make check`, `golangci-lint run ./...`, `make wasm`,
  `make spec-check`
- [x] 4.2 Capture the page at desktop and phone width and inspect it
- [x] 4.3 Update `CHANGELOG.md`; record evidence in `VALIDATION.md`
