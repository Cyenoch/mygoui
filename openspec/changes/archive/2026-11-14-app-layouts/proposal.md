# Proposal

## Why

mygoui has components but nothing at the level of a window. An application
author who wants what every native desktop application has (a custom title
bar that keeps clear of the system window controls, a sidebar beside the
content, an inspector, a status strip, and a narrow-window adaptation) must
rebuild it from `c.TitleBar()`, `DragWindow`, `ResizablePanel` and
`overlay.Sheet`. The gallery already did exactly that, privately, in
`internal/gallery`: the insets, the shared title band height, the controls
that move below the band when they do not fit, and the compact navigation
sheet are correct there and unavailable to anyone else.

## What Changes

- Add a `layout` family package with three builders and two helpers:
  - `layout.TitleBar`: a draggable title band with leading, centre and
    trailing content that reserves the real window-control insets, covers the
    native band, grows with zoom and text scale, and moves its controls below
    the band when they no longer fit between the window controls.
  - `layout.SplitView`: a start pane, the content and an end pane (sidebar |
    content | inspector, or list | detail) with resizable or fixed panes,
    minimum sizes, collapse, and a compact mode in which a side pane is shown
    over the content as a sheet.
  - `layout.Scaffold`: the standard window shell composed from the two:
    a full-height start pane sharing the title band, the title bar over the
    content, an end pane, and a status strip.
  - `layout.PaneToggle`: the title-bar button that shows or hides a pane in
    either mode; `layout.Compact`: whether the window is compact.
- Add `overlay.SheetOptions.Flush`: a sheet that hands its whole panel to the
  caller's body, which the compact pane needs.
- Add three library messages to `locale` (`TitleBar`, `ToggleSidebar`,
  `ToggleInspector`).
- `mygoui.ResizablePanel`: a handle on the start or top edge precedes the
  content in Tab order, and `Shown` reports the size of each build so the
  scaffold's title bar can follow a folding pane.
- `overlay`: side sheets and drawers enter from the edge they rest on in
  right-to-left windows too.
- The gallery adopts `layout.Scaffold`, `layout.TitleBar` and
  `layout.PaneToggle` in place of its private shell. Its sidebar collapse now
  animates as `DESIGN.md` already prescribes for a collapsing sidebar, and its
  compact navigation shows the same pane it docks.
- New gallery pages (Title Bar, Split View, Scaffold) and an
  `examples/shell` application.
- Documentation: `docs/components.md`, `docs/design.md`,
  `docs/behavior/layout.md`, `README.md`. Text for `DESIGN.md` and the
  `mygoui-components` skill is handed to their owners in this change's
  report; other agents hold those files.

Not in this change: a navigation split view that turns list | detail into a
drill-down stack (compose `navigation.NavigationStack` when compact), a
title bar spanning above a sidebar as a scaffold variant (a Column of
`TitleBar` and `SplitView` already is one), vertical splits (use
`ResizablePanel`), persistence of pane sizes, and window-level menus.

## Capabilities

### New Capabilities

- `app-layouts`: application-level layout building blocks: the title bar,
  the split view, the scaffold, the pane toggle and the compact rule, and the
  flush sheet they rely on.

### Modified Capabilities

- `gallery-workbench`: the gallery's shell is built from the library's layout
  family, and its sidebar collapse moves.

## Impact

- New package `github.com/Cyenoch/mygoui/layout`, importing the root and
  `overlay`. No existing package imports it.
- `overlay/sheet.go`, `overlay/chrome.go`: one new option.
- `locale/locale.go`: three messages with English defaults.
- `internal/gallery`: `gallery.go`, `docs_layout.go`, new page file,
  regenerated `api_gen.go` and example sources; one gallery test advances the
  clock across the collapse.
- `examples/shell`: new consumer.
- No dependency change. No import-time work. No persistence.
