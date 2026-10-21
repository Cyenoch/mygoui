# MyGo 0.6 upgrade

Minimum toolkit: `github.com/egoist/mygo v0.6.0`. The reviewed upstream [release notes](https://github.com/egoist/mygo/releases/tag/v0.6.0) and [full comparison](https://github.com/egoist/mygo/compare/v0.5.5...v0.6.0) are the authority. The pinned `ui/doc.go`, `docs/ui/migration.md`, widget docs and public source were inspected during integration.

## Relevant changes and integration

| Upstream change | Library integration |
| --- | --- |
| Partial GPU redraw, incremental layout, one measurement per text element, clipped-content paint culling | Components retain stable native keys, widget state and paint-only animation. No view ticker, extra rendering cache or forced full redraw is introduced. Existing pixel/layout/motion tests exercise native composition. GPU throughput remains a separate platform measurement. |
| New `plugins/loading` with 29 native indicators | Spinner uses `loading.Arc`, deleting its own stationary-ring painter and delegating visibility, animation cadence, reduced-motion poses and progress semantics to the maintained plugin. The default arc preserves the intended compact style; no new variant registry is needed. |
| Native cell blocks, drag row selection, row headers and cell menus | Table passes `Column.RowHeader` through, and retains caller-configurable `TableState.List.CellSelect`, `DragSelect`, `CellMenu` and `MenuKeepsSelection`. Native `CellSelected`, `SelectColumns`, `ExtendTo` and `ScrollToEnd` remain accessible through ListState. DataTable exposes the same state through `state.Table.List`. Gallery examples toggle real cell/drag selection. |
| Wheel scrolling glides smoothly | ScrollArea continues to use native ScrollState. A regression uses `Tester.Wheel` and the virtual clock to verify actual bounded scrolling without moving a fixed header. |
| List following its end no longer pushes short content down | Message Scroller delegates anchors/follow/end scrolling to the native ListState. Existing message/geometry regressions run against v0.6.0. |
| Long select popup stays in the window; field labels align with descriptions | Form/navigation controls retain native base widgets and their sizing constraints; existing narrow/enlarged/input tests cover their composition. |
| Rich text transition exit crash fixed; editor large-file improvements | Code Block uses native source editing/highlighting and selection, preserving exact source and streaming state. No replacement editor is introduced. |
| Windows 11 menus and native menu shortcuts on Windows/Linux | Native menu compositions inherit platform changes. Platform input/appearance acceptance remains separate from headless tests. |
| `Tester.Click`, `ClickWith` and `RightClick` reject off-window targets | The entire headless suite runs under the stricter contract. Controls intended for interaction must be visible and bounded. Use unique accessible labels when a row and its cell display the same text; a row's geometric center can be in a different column. |
| Checked Element migration and persistent identity | `go tool mygo migrate-ui .` reports no required changes. `go tool mygo vet .` checks stored build values and goroutine captures. |

## Deliberate scope

Terminal mirroring/search/graphics, WebView shortcuts, trash handling, tray icons, installers, signing, updater behavior and custom window buttons are host-application capabilities. This component library does not add application lifecycle or platform integration to expose them. Applications may use those MyGo APIs directly.

Single sliders retain horizontal/vertical behavior and ranges retain the native two-thumb horizontal engine. This release does not add arbitrary-thumb or vertical-range APIs. Drawer remains an edge modal with no swipe/velocity arbitration. Those limits are explicit in the component catalog.

## Verification

`VALIDATION.md` records the exact completed checks. Headless tests verify native layout, input, selection, paint and virtual-clock behavior. They do not establish Metal/Direct3D/OpenGL partial-presentation performance, native OS menu appearance or screen-reader compatibility.
