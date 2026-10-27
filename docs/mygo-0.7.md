# MyGo 0.7.1 integration

The library and `go tool mygo` both resolve `github.com/egoist/mygo v0.7.1`.
The minimum Go version is 1.27.3, following upstream's 0.6.4 security update.
`make check`, `make dev` and `make app` use the module's tool directive;
they no longer resolve an unrelated `latest` CLI.

Upstream release notes: [0.6.4](https://github.com/egoist/mygo/releases/tag/v0.6.4),
[0.7.0](https://github.com/egoist/mygo/releases/tag/v0.7.0),
[0.7.1](https://github.com/egoist/mygo/releases/tag/v0.7.1).
The [0.6 integration record](mygo-0.6.md) is historical.

## Adopted capabilities

| Toolkit change | Native integration |
| --- | --- |
| Content zoom (0.7.0) | Gallery Cmd/Ctrl `+`, `-`, `0` uses `Context.SetZoom`, with window-owned steps from 50% to 300%. Native text, input, IME, accessibility and menus scale together. Responsive breakpoints use content DIPs, not the physical window width. |
| Zoom-aware geometry | Gallery keeps its system title-bar strip 52 window DIPs high and respects the zoom-adjusted native button insets. `overlay.FlyoutAnchor` converts committed handle bounds back to parent-window DIPs. |
| Native flyout windows | `overlay.NewFlyout` composes a coss popover surface in a separate native window. Explicit parent/tray, dimensions, focusability, theme and lifecycle replace browser portals. The Popover Gallery page includes a native example. See [flyout.md](flyout.md). |
| Paint-only frames / GPU damage (0.7.0) | Existing Progress, Spinner and Skeleton draw callbacks use native paint scheduling instead of view tickers. Standalone skeletons cap slow shimmer/pulse at 60 fps using `Painter.AnimationFrameRate`; grouped effects use the same cadence but may repaint their layer. Tests check that animation changes without rebuilding the view, and ordinary draw animations do not repaint an unrelated still callback. |
| `Element.Hidden` (0.7.0) | `TabsOptions.KeepMounted` opt-in retains inactive panels' editor/local/scroll state. Hidden panels take no room and do not intercept keyboard focus. The default still builds only the selected panel. Retained callbacks run during builds: they must not perform unconditional side effects. |
| Per-window caption appearance (0.7.0) | Gallery forwards explicit Light/Dark changes to `Window.SetWindowButtonsTheme`. Windows custom caption buttons follow the specimen theme; macOS and Linux keep the toolkit's app-appearance policy. |
| CJK editing, route-return colors, short labels (0.6.4) | Form controls remain native. No duplicate editor, ellipsis algorithm or route color cache is introduced. Existing editing, navigation and retained paint regressions run against the new toolkit. |
| Clipped/scrolling animation and macOS drag regions (0.7.x) | Overlay chrome, scrolling surfaces and the draggable title bar continue to use native layout/motion. These fixes do not by themselves establish arbitrary nested rounded-clip correctness on every GPU. |

## Deliberately not applied globally

- `Element.Cache` is not added to every container: cached layers omit embedded
  WebViews, have an 8192-pixel side limit and trade GPU memory/text rendering
  behavior for drawing work. There is no measured expensive static subtree
  requiring a new blanket cache in this component library.
- Menu icons are available on `mygo.MenuItem`, not `ui.MenuItem` in 0.7.1.
  `navigation.ContextMenu` preserves the latter's native input lifecycle;
  its SVG item icon cannot honestly be wired to an unavailable method.
- `Font.LineHeight` is available to custom `ShapeText`/`ShapeRichText` consumers.
  Components use native text elements with `FixedLineHeight` already; there
  is no reason to replace them with manually painted text.
- Audio, web-page windows, installers and auto-updaters are application services,
  not generic component dependencies. No import-time app initialization,
  extra network service or browser UI is added.

## Layout and ownership

`ui.Column` stretches children unless controls specify their own alignment;
`ui.Row` centers them on the cross axis. Intrinsic groups explicitly own their
members' width; `FullWidth` remains opt-in. A bold label must be measured with
its actual font weight. Keep native focus/editing/accessibility behavior.

Elements, custom parts and contexts remain build-scoped. Retain handles and
application state per window. Flyout content receives its own context; do not
capture an element or context from the parent in another window's view.

## Verification

Run `make check`, `make race`, `make wasm`, `golangci-lint run ./...` and
`make workflows`. Migration preview and lifetime vet use the module-pinned CLI.
Current results, GPU evidence and platform gaps belong in
[VALIDATION.md](../VALIDATION.md), not in an assumed release-note guarantee.
