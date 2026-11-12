# Design

## Context

MyGo's `NewFlyout` creates a frameless, non-resizable top-level window owned
by a parent window or tray. It anchors to a rectangle in parent DIPs, flips,
slides or shrinks inside the screen work area, can take keyboard focus,
dismisses on outside press, focus loss and Escape, and has no WebAssembly
implementation. It has no backdrop, modality, focus trap or accessibility
relation to the parent. Its content is a separate `ui.View` root with its own
build-scoped context. Components receive only a `*ui.Context`, which exposes
no window.

## Decisions

### Caller-owned host

`mygoui.FlyoutHost` holds the parent `*mygo.Window`. The caller keeps one per
window, like other widget state, and passes a pointer in the component
options. The root package owns the type because every family imports it;
only `overlay` imported MyGo's root package before. A nil host keeps the
in-window popup, so existing callers are unchanged.

### One shared presenter

`internal/anchored` gains the one implementation every component calls:

- the build that sees the caller's open state posts the window creation with
  `Window.Update`, since MyGo creates windows on the main thread outside
  builds; the window is created hidden with `FitContent` and then shown;
- presenter state is keyed by the caller's open pointer, which is stable,
  rather than the trigger's tree position, which can change between builds;
  a `ui.Local` closer on the trigger, keyed by that pointer, closes the
  flyout when the trigger leaves the tree, and a build that presents the
  popup in-window closes one it presented before; size bounds are content
  units converted with the parent's zoom;
- each parent build re-anchors from the trigger's committed bounds and
  invalidates the flyout so it rebuilds from the same caller state; a flyout
  build of its own invalidates the parent, and one that follows a parent
  build does not, which ends the exchange;
- `OnClosed` sets the caller's open state to false, runs the component's
  dismissal callback and invalidates the parent; a closed caller state closes
  the window;
- a nil host or window, and WebAssembly, fall back to the in-window path.

Each component splits its popup body into a function of a context, so the
same body builds in the parent popover or the flyout window. Bodies may not
capture parent elements, handles bound in the parent, or `ui.Local` state
keyed to parent elements; per-popup widget state moves into the caller-owned
component state.

### Scope

Popover, Dropdown (top level), Select, DatePicker, DateRangePicker and
ColorPicker are focusable panels whose state is caller-owned, so they move
cleanly. A Dropdown with submenus stays in-window: a submenu inside a fitted
flyout window would be clipped, and chaining flyouts is left out. Dialog-family overlays need modality and a backdrop; ComboBox and
Autocomplete need the input to keep the keyboard while the list is in another
window, losing the active-descendant relation; Tooltip and PreviewCard depend
on hover tracking across windows; Toast belongs to the window viewport.
These stay in-window.

### Motion and appearance

The flyout uses the coss popover fill and radius with the system shadow.
macOS and Windows outline it themselves, so it draws the coss border only on
Linux. In-window popup motion does not cross windows; flyouts use the
platform's window presentation.

## Risks

- Placement, dismissal, focus transfer and screen-reader behavior are
  platform-owned and not verifiable headlessly; they are recorded per OS in
  VALIDATION.md, not promised by the spec.
- macOS passes an outside press through to its target; Linux and Windows
  swallow it. Wayland can dismiss regardless of `OnClose`.
- Two windows must be invalidated; a missed refresh leaves stale content.
  The presenter owns both refreshes so components do not.
