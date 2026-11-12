# MyGo 0.9.0 integration

The library and `go tool mygo` both resolve `github.com/egoist/mygo v0.9.0`.
The minimum Go version stays 1.27.3. The [0.7 integration record](mygo-0.7.md)
is historical; its ownership, flyout and automation guidance still applies
unless this page says otherwise.

Upstream release notes: [0.7.3](https://github.com/egoist/mygo/releases/tag/v0.7.3),
[0.7.4](https://github.com/egoist/mygo/releases/tag/v0.7.4),
[0.8.0](https://github.com/egoist/mygo/releases/tag/v0.8.0),
[0.8.1](https://github.com/egoist/mygo/releases/tag/v0.8.1),
[0.8.3](https://github.com/egoist/mygo/releases/tag/v0.8.3),
[0.9.0](https://github.com/egoist/mygo/releases/tag/v0.9.0).

## Toolkit behavior changes handled

| Toolkit change | Native integration |
| --- | --- |
| Focus scrolls the element into view (0.9.0) | The year picker focuses its grid with `FocusOptions.PreventScroll` on first build, so opening it keeps the selected year's row instead of jumping to the top. |
| An element's last `Transition` in a frame governs its return from an exit copy (0.9.0) | A responsive dialog reopened during its exit uses one short transition on the reopen frame. The returning backdrop carries the exit motion's share of the offset; the new panel adds only the offset it had when it closed. |
| `Tester.Click` checks the real input target (0.8.x) | Tests click the control that receives the press (a day cell, a button outside an open popup) instead of a container center. |
| Every key reaches `HandleInput` (0.8.0) | Chart, color area and color slider handlers select keys with a `switch` and return false for others, so lone modifiers and Alt keep their system meaning. |
| Inactive windows gray the accent (0.8.0) | Components keep MyGo's default; applications that want the accent in background windows set `Theme.InactiveAccent`. |
| `Opacity` fades a subtree as one group (0.8.0) | Fades on button loading, outside calendar days, drawers and toasts now composite as one layer, as CSS does; their headless paint tests pass unchanged. |
| Presses and drags follow the primary button (0.8.0) | Color area and channel slider drags already ignore other buttons. |

## Verification

Run `make check`, `make race`, `make wasm`, `golangci-lint run ./...` and
`make workflows`. `go tool mygo migrate-ui .` reports only the existing
callback-element key notices in `overlay/popover.go` and `overlay/tooltip.go`.
Current results and platform gaps belong in [VALIDATION.md](../VALIDATION.md).
