# MyGo 0.7.2 integration

The library and `go tool mygo` both resolve `github.com/egoist/mygo v0.7.2`.
The minimum Go version is 1.27.3, following upstream's 0.6.4 security update.
`make check`, `make dev` and `make app` use the module's tool directive;
they no longer resolve an unrelated `latest` CLI.

Upstream release notes: [0.6.4](https://github.com/egoist/mygo/releases/tag/v0.6.4),
[0.7.0](https://github.com/egoist/mygo/releases/tag/v0.7.0),
[0.7.1](https://github.com/egoist/mygo/releases/tag/v0.7.1),
[0.7.2](https://github.com/egoist/mygo/releases/tag/v0.7.2).
The [0.6 integration record](mygo-0.6.md) is historical.

## Adopted capabilities

| Toolkit change | Native integration |
| --- | --- |
| Themed context menus (0.7.2) | `navigation.PopoverMenu(anchor, MenuOptions)` reuses the toolkit's in-window context menu, including pointer/menu-key/Shift+F10 triggers, type-ahead, disabled-item skipping, submenus, dismissal and focus restoration. `navigation.ContextMenu` remains OS-owned by default. Both share selection/callback wiring; single choice uses `MenuItem.Radio` instead of checkbox semantics. The Gallery ContextMenu page demonstrates both. |
| Unified images (0.7.2) | Avatar uses `ui.ImageSize` instead of a bitmap/SVG type switch, preserving initials for nil, typed-nil and empty images. Hosts can decode loaded bytes with `ui.DecodeImage` and supply the resulting `ui.ImageSource` directly; the library still performs no fetching. |
| Native app automation (0.7.2) | Gallery's `e2e` build tag imports `plugins/automation`. An explicit `MYGO_AUTOMATION` socket enables real-window accessibility, input and capture tests; ordinary builds do not import it. The reusable library never installs a plugin. |
| Native fixes (0.7.2) | Existing controls inherit improved accessible row names, overlay opening-press handling, diff-view scrolling/painting and platform resize/redraw fixes without parallel implementations. Headless tests and real-window evidence are recorded separately. |
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
- Menu icons remain available on `mygo.MenuItem`, not `ui.MenuItem` in 0.7.2.
  `navigation.ContextMenu` and `PopoverMenu` preserve the latter's input
  lifecycle; neither renders item icons, descriptions, danger styling or
  section headings. `MenuItemOptions.Shortcut` remains a display string for
  the custom `Menu`/`Dropdown` path, not a typed native accelerator.
- `Font.LineHeight` is available to custom `ShapeText`/`ShapeRichText` consumers.
  Components use native text elements with `FixedLineHeight` already; there
  is no reason to replace them with manually painted text. Paragraph
  justification, indentation, baseline shifts and stable rich-document IDs
  are available directly to host editors; Composer remains plain text.
- `ui.Video`, `VideoPlayer` and `NewGIFDecoder` are available directly for
  host-owned media. The library has no media-player contract or existing
  animated avatar; no duplicate wrapper, decoder or playback loop is added.
- Audio, page sessions/security/context menus, system proxy routing,
  installers and auto-updaters are application services, not generic component
  dependencies. No import-time app initialization, extra network service or
  browser UI is added. The 0.7.2 page pinch/swipe defaults do not affect this
  native Gallery.

## External Gallery automation

Use a private socket directory, and keep the test binary out of distribution:

```sh
go build -tags e2e -o /tmp/mygoui-gallery-e2e ./examples/gallery
socket_dir=$(mktemp -d)
MYGO_AUTOMATION="$socket_dir/gallery.sock" /tmp/mygoui-gallery-e2e
```

The plugin speaks newline-delimited JSON; `windows`, `tree`, `act`, `input`,
`capture` and `idle` drive the built native app. Wait for `idle` between
dependent actions rather than sleeping. OS context menus can enter a native
menu loop; test those with `ui.NewTester`, not external right-click automation.
The in-window `PopoverMenu` can be driven through normal input and its
accessible tree. Importing the plugin alone does not listen without the
environment variable, and `mygo_noinspector` builds omit it. See the pinned
[automation protocol](https://github.com/egoist/mygo/blob/v0.7.2/docs/plugins/automation.md).

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
