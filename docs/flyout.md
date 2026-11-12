# Native flyout

## Anchored popups as flyouts

Popover, Dropdown, Select, DatePicker, DateRangePicker and ColorPicker stay
inside the window by default. A caller that wants a popup to extend beyond
the window opts in per component with `Flyout`:

```go
// Window-owned state, set once the window exists:
state.Flyouts = mygoui.FlyoutHost{Window: window}

// In the view:
form.Select(c, form.SelectOptions{Label: "Fruit", Items: items,
    State: &state.Fruit, Flyout: &state.Flyouts})
```

The component opens a fitted native flyout after the build that sees it open,
re-anchors it as the trigger moves, closes it when its open state turns false
or the trigger leaves the view, and clears that state when the user dismisses
it. The flyout follows the parent's light or dark mode. Popup content builds
in the flyout's own context, so Popover's `Header`, `Body` and `Footer` take
one. A nil host, a host without a window, and WebAssembly keep the in-window
popup. A Dropdown whose menu has submenus stays in the window. Flyouts have no
popup motion or Popover arrow, do not follow the parent's content zoom, and
macOS and Windows draw their outline. Dialogs, sheets and other modal layers
have no flyout form.

## Custom flyout windows

`overlay.NewFlyout` creates a separate native MyGo window that can extend
outside its parent. It is not the in-window `overlay.Popover`, a WebView, or
a browser portal. Create it in an input action on the main thread after
`App.WhenReady`, not unconditionally while building a view. The caller owns
the returned `*mygo.Window`, its model and any widget handles. The parent
or tray must be explicit, and the window needs a positive width and height in
window DIPs.

```go
// Fields belonging to the parent window's state:
var trigger ui.Handle
var flyout *mygo.Window

// In the parent view, after theme.Apply(c):
ui.Button(c, "Details").Bind(&trigger).OnClick(func() {
    if flyout != nil {
        return
    }
    w, err := overlay.NewFlyout(overlay.FlyoutOptions{
        Parent: parentWindow,
        Anchor: overlay.FlyoutAnchor(c, trigger),
        Placement: mygo.PlacementBottomStart,
        Gap: 4, Width: 320, Height: 240, Focusable: true,
        Surface: overlay.FlyoutSurfaceOptions{
            Title: "Details",
            Body: func(c *ui.Context) { ui.Text(c, "Native content") },
        },
    })
    if err != nil {
        // Report the error through the application's normal error handling.
        reportError(err)
        return
    }
    flyout = w
    w.OnClosed(func() {
        flyout = nil
        parentWindow.Invalidate()
    })
})
```

Keep these fields per parent window, not as package globals. `FlyoutAnchor`
reads the handle's committed bounds in the current parent context, converts
content DIPs using its UI zoom and rounds outward to window DIPs. It never
retains the handle's element or context. For a point or other known window
rectangle, pass `mygo.Rectangle` directly. If the target moves or zoom changes
while the flyout is open, the owner must recompute the anchor from a current
parent context and call `flyout.SetAnchor`; MyGo follows parent window moves,
but cannot infer changes to the element's layout.

For tray use set `Tray` instead of `Parent`; the anchor is ignored. The
callback for each `Header`, `Body`, and `Footer` receives the flyout's own
build-scoped context. Do not capture the parent context, elements or custom
parts in those callbacks. Per-flyout widget state and handles belong to the
caller and must outlive the window; background results should use the
flyout's `Window.Update`. Native window creation failures and pre-ready calls
can still panic in MyGo; option errors are returned before window creation.

The independently rendered window installs `theme.Apply` each build and
`FlyoutSurface` draws the coss popover colors, border, rounded surface,
header, body and footer inside a native scroll viewport. The header callback
replaces the visual title and description; `Label` overrides the title as the
surface's accessible label. The entire content scrolls, including the footer,
when it exceeds the window. `FlyoutSurface` is also usable independently
inside an already themed view and returns a build-scoped element. The
window requests the system shadow; its appearance and corner treatment vary
by platform, including system outlines on macOS and Windows 11.

Set `FlyoutOptions.Dark` to a caller-owned `*bool` to follow an explicit
application appearance instead of the OS. After changing it, invalidate the
flyout too. A nil pointer follows the system. The Gallery's Popover page demonstrates a
Popover opened with `Flyout`.

## Native behavior and limits

- MyGo places the flyout relative to the anchor and screen work area, flipping,
  sliding or shrinking it when needed. It owns platform focus and lifecycle.
- With `Focusable`, MyGo takes keyboard focus and dismisses on outside press,
  another window taking focus, or app deactivation. An unhandled Escape closes
  either focusable or non-focusable flyouts; non-focusable content cannot take
  keyboard focus. An `OnClose` listener can normally prevent dismissal, but
  Wayland's compositor can dismiss it regardless. Use `OnClosed` to reconcile
  caller state, and do not assume focus is restored identically on every OS.
- On Linux and Windows the outside press dismisses without activating the
  target; on macOS it generally reaches the target, except the anchor. Linux
  tray placement is centered in the primary work area because tray icons do
  not provide position/click information. See MyGo `FlyoutOptions` for the
  platform contract.
- MyGo's WebAssembly backend does not implement native flyout positioning;
  `NewFlyout` returns an error there rather than presenting a fake in-page
  flyout. Use an in-window `Popover` explicitly if the application needs one.
- Headless surface tests do not verify OS placement, window shadows, native
  dismissal, keyboard focus transfer, GPU clipping or screen readers. Validate
  these separately on the target OS and desktop environment.
