# MyGo 0.6.3 integration

Minimum toolkit: `github.com/egoist/mygo v0.6.3`, upgraded from v0.6.0 during
this rewrite. The reviewed upstream release notes are the authority:
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

MyGo is not CSS. Two defaults matter for HeroUI parity:

- `ui.Row` centers children across the row by default; CSS flex rows stretch.
  Rows that need top alignment set `AlignItems(ui.Start)`. Button content is
  centered explicitly and can be left/right aligned with `ButtonAlign`.
- `ui.Column` stretches children to its width, matching CSS's align-items:
  stretch. Containers that should hug content set `Shrink(0)` or
  `AlignSelf(ui.Start)`.

HeroUI's `Start` is the zero `ui.Align` value. Components therefore do not use
`ui.Align != 0` as a presence check; they define their own enum where the zero
value must mean centered (for example `ButtonAlignCenter`).

## Verification

`go build ./...`, `go vet ./...`, the ordinary and race suites, the WebAssembly
build, lint and `make vet` (MyGo CLI at `@latest`) pass against v0.6.3. Headless tests cover
native editing, focus, selection, scrolling, overlays, tables and paint states.
GPU presentation, OS gestures, native menu appearance and screen-reader output
remain platform checks.
