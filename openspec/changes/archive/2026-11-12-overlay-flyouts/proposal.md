# Proposal

## Why

Anchored popups draw inside the window. Near a window edge or in a small
window, a Select list, date calendar or popover menu is squeezed, flipped or
clipped by the window bounds. MyGo 0.9 offers native flyout windows that can
extend past the parent window within the screen work area, and the library
already wraps them as `overlay.NewFlyout`, but no component can use them.

## What Changes

- A caller-owned `mygoui.FlyoutHost` names the parent window. Popover,
  Dropdown, Select, DatePicker, DateRangePicker and ColorPicker accept an
  optional host; with one, their popup opens as a native flyout.
- Without a host, or where a flyout cannot be created (WebAssembly), the same
  popup opens in the window as today.
- Dialogs, sheets, drawers, toasts, tooltips, preview cards, ComboBox and
  Autocomplete stay in-window: they need a backdrop, modality, focus trap,
  hover tracking or an input that keeps the keyboard, which a flyout window
  does not provide.
- `overlay.FlyoutOptions` exposes MyGo's content fitting with size bounds;
  stale documentation that says native content cannot fit is corrected.

## Capabilities

### New Capabilities

- `overlay-flyouts`: opt-in native flyout presentation for anchored popups.

### Modified Capabilities

None.

## Impact

Root package (`FlyoutHost`), `internal/anchored`, `overlay`, `navigation`,
`form`, `date`, `color`, the gallery's examples, docs/flyout.md,
docs/components.md and CHANGELOG. In-window behavior is unchanged when no
host is given.
