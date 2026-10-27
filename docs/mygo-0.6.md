# MyGo 0.6.3 integration (historical)

Current integration: [MyGo 0.7.1](mygo-0.7.md). The results below describe
the earlier version, not current validation.

Pinned toolkit: `github.com/egoist/mygo v0.6.3`. The earlier dependency upgrade
from v0.6.0 is retained by the coss reconstruction. Upstream release notes:
[v0.6.1](https://github.com/egoist/mygo/releases/tag/v0.6.1),
[v0.6.2](https://github.com/egoist/mygo/releases/tag/v0.6.2) and
[v0.6.3](https://github.com/egoist/mygo/releases/tag/v0.6.3).

The patch releases kept the public API stable. The library builds unchanged
against 0.6.3; no migration shims were added.

## Changes the components rely on

| Upstream change | Integration |
| --- | --- |
| Faster text layout for large paragraphs and inline transforms checked once per text (0.6.3) | Typography and form help text benefit without a library-side cache. No copied text layout is kept. |
| Table resizing leaves no edge after the last room-sharing column (0.6.3) | The Table component keeps native resize behavior and its regression test exercises header dragging. |
| Presses pass through a hidden scrollbar's track (0.6.3) | Overlay bodies and menus continue to use native scroll views; no custom hit testing is added. |
| Touchpad pinches reported as `InputZoom` (0.6.2) and GUI pinch tests corrected on macOS (0.6.3) | ScrollShadow and table scrolling stay native; diagram/zoom gestures remain application behavior. |
| Windows ConPTY, DWM caption, WebView2, task dialog and Metal atlas fixes | No component changes; platform behavior is inherited from the toolkit. |
| 0.6.0 partial GPU redraw, incremental layout, `plugins/loading`, native cell/drag selection, smooth wheel scrolling and strict `Tester.Click` bounds | Components keep stable keys and native state; Spinner delegates to `loading.Arc`; no view ticker or parallel render cache exists. |

## Native layout defaults

Two native layout defaults matter when translating the coss reference:

- `ui.Row` centers children across the row by default; CSS flex rows stretch.
  Rows that need top alignment set `AlignItems(ui.Start)`. Button content is
  centered explicitly and can be left/right aligned with `ButtonAlign`.
- `ui.Column` stretches children to its width, matching CSS's align-items:
  stretch. Containers that should hug content set `Shrink(0)` or
  `AlignSelf(ui.Start)`.

`ui.Start` is the zero `ui.Align` value. Components therefore do not use
`ui.Align != 0` as a presence check; they define their own enum where the zero
value must mean centered (for example `ButtonAlignCenter`).

## Verification

`make check MYGO_VERSION=v0.6.3` runs the headless tests, build, Go vet and
pinned MyGo lifetime vet. Also run `make race`, `make wasm` and
`golangci-lint run ./...` after integration. Current results are recorded in
[VALIDATION.md](../VALIDATION.md), separately from historical dependency work.
Headless tests exercise editing, focus, selection, scrolling, overlays, tables
and paint states. GPU presentation, OS gestures, native menu appearance and
screen-reader output require separate platform evidence.
