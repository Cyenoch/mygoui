# Tasks

## 1. Planning

- [x] 1.1 Write proposal, delta specs (`app-layouts`, `gallery-workbench`),
  design and tasks; verify with
  `npm run --silent openspec -- validate app-layouts --strict --no-interactive`

## 2. Library

- [x] 2.1 Add `overlay.SheetOptions.Flush` and a test that a flush side sheet
  hands the whole panel to its body; verify `CGO_ENABLED=0 go test ./overlay/...`
- [x] 2.2 Add the `TitleBar`, `ToggleSidebar` and `ToggleInspector` messages
  to `locale` with English defaults
- [x] 2.3 Add `layout.TitleBar` (slots, insets by span, band height, measured
  stacking, drag surface) and `layout.Compact`
- [x] 2.4 Add `layout.SplitView` (resizable and fixed panes, content minimum,
  collapse, compact sheet, right-to-left) and `layout.PaneToggle`
- [x] 2.5 Add `layout.Scaffold` (shared band height, span by pane state,
  strip)
- [x] 2.6 `ui.NewTester` scenarios in `layout/`: insets reserved and released,
  stacking, controls keep their presses and the bar is the target elsewhere,
  divider pointer and keyboard resizing with clamping, content minimum,
  collapse and keyboard return, compact sheet open and dismissal, flag reset
  on widening, content focus across modes, right-to-left, toggle in both
  modes, scaffold band height and span, reduced motion

## 3. Consumers

- [x] 3.1 Build the gallery shell from `layout.Scaffold`, `layout.TitleBar`
  and `layout.PaneToggle`; keep every existing gallery test passing, editing
  only the one that asserts an instant sidebar collapse
- [x] 3.2 Add gallery pages for Title Bar, Split View and Scaffold with live
  examples, following the component page pattern
- [x] 3.3 Add `examples/shell`, a small application using the scaffold
- [x] 3.4 Run `make generate`

## 4. Documentation

- [x] 4.1 `docs/components.md`: a layout family section; the Sheet row
  mentions `Flush`
- [x] 4.2 `docs/design.md`: package table, dependency rule, the title-bar
  contract moved to the library's description
- [x] 4.3 `docs/behavior/layout.md` with one entry per builder and a gap
  list; its row in `docs/behavior/README.md`
- [x] 4.4 `README.md` usage note
- [x] 4.5 Hand the `DESIGN.md` layout rules and the `mygoui-components`
  skill text to their owners in the report (those files are held by other
  agents during this change)

## 5. Verification

- [x] 5.1 Render light and dark, compact and wide output of the gallery
  shell and the new pages (`make design-artifacts` and Tester images),
  inspect it and fix what is wrong
- [x] 5.2 `make check`, `make spec-check`,
  `GOOS=js GOARCH=wasm CGO_ENABLED=0 go build ./...`
- [x] 5.3 Record the real-window checks (window drag and maximize from the
  bar, window-control clearance, divider cursor, screen readers) as pending
  evidence in `VALIDATION.md`

## 6. Review fixes

- [x] 6.1 Scaffold: the start divider leaves the content its minimum beside
  a docked end pane; test
- [x] 6.2 Scaffold: the title bar reserves what the start pane leaves of the
  start inset in every frame of the fold (`ResizablePanelOptions.Shown`);
  frame-sampling test without reduced motion
- [x] 6.3 Title bar: measure with `LayoutBounds` at build; paint callbacks on
  inner rows only invalidate; a decorated bar still stacks
- [x] 6.4 overlay: side sheets and drawers enter from their own edge right to
  left; test
- [x] 6.5 `ResizablePanel`: a start- or top-edge handle precedes the content
  in Tab order; test
- [x] 6.6 Gallery: sidebar hairline in `SidebarBorder`, brand offset, page
  area follows the fold; the start sheet's band clears the end controls
- [x] 6.7 Docs: splitter range, nil `Size` identity, `Flush` limits, focus
  order, motion tokens; the zoom assertion can fail
