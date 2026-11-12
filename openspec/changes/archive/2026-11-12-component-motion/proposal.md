## Why

Each component's motion was reviewed against coss for entry, exit and state
changes. Most already match; a few differ or are missing: the switch thumb
does not stretch while pressed, the select chevron turns although coss shows
a static up-down icon, preview cards do not shrink on exit, the sidebar
snaps between widths, tabs and the bottom sheet use a different curve than
coss, and the indeterminate progress sweep redraws at the display's full
refresh rate.

## What Changes

- Switch: the thumb stretches 10% toward the track's centre while pressed.
- Select: a static up-down chevron replaces the turning chevron.
- PreviewCard: fades and shrinks to 98% on exit, as coss preview cards do.
- Sidebar: collapses and expands over 200 ms, linearly.
- Tabs indicator: coss ease-in-out over 200 ms.
- BottomSheet: arrives on the coss drawer curve.
- Indeterminate progress: paint-only sweep capped at 60 fps.

Reduced motion keeps its current contract: these motions settle at once.

## Impact

`form/switch.go`, `form/select.go`, `form/combo_box.go`, `overlay/popover.go`,
`overlay/preview_card.go`, `sidebar.go`, `navigation/motion.go`,
`overlay/drawer.go`, `progress.go`, their tests and `docs/components.md`.
No public API changes.
