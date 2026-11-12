# Design

## Context

Four read-only audits (form; navigation, list and data; root and theme;
overlay, date and color) and a cross-package duplication audit compared each
family with coss commit 3d8ca28 and MyGo v0.9.0. This change takes the shared
findings and the duplicate-caused defects; larger rebuilds are deferred below.

## Decisions

- **Elevation.** `OverlayShadows` returns shadow-lg/5 plus the same edge
  layer surfaces use (`0 1px black/4%` light, `0 -1px white/6%` dark), drawn
  after the drop shadow. `TooltipShadows` is the md/5 variant used by
  `overlay.Tooltip`; coss tooltips are the only md/5 popup surface.
- **Toolbar.** coss toolbars are bordered cards without shadow; the overlay
  shadow is removed.
- **Popup surface and motion live in `internal/anchored`.** That package
  already owns placement and its 4 DIP default offset, so per-package offset
  defaults go away. `popupMotion` keeps the anchor-facing transform origin and
  MyGo exit copies. coss gives only popovers and tooltips motion: popovers
  enter from 98% and fade out, tooltips also shrink on exit. Menus, selects,
  combo boxes and autocomplete have no starting or ending style, so their
  fades are removed; date and color pickers are popovers and take the
  popover motion.
- **Error implies invalid.** Fields compute `invalid || error != ""` once in
  the shared field helper, matching coss `aria-invalid` usage.
- **Native replacements.** `Painter.FillRadii` replaces the two cubic corner
  fills (and fixes mixed continuous/circular corners on macOS);
  `ui.Color.Mix` replaces table `mixColor`; context menus map Description to
  `Subtitle`, Shortcut to `Trailing` and the danger variant to `Danger`.
- **Status icon.** coss toasts use circle-alert for errors; `StatusIcon`
  is exported so new status surfaces reuse the one mapping.
- **Disabled opacity stays 50%.** It is MyGo's native dimming and the
  documented deviation; moving to coss 64% needs native checks of nested
  disabled widgets first.

## Deferred

Recorded for later changes, not part of this one: a shared ComboBox and
Autocomplete core with one filtering/typeahead model across Select, ListBox
and Command; Menu/SubMenu/Dropdown on `ui.MenuBase`; Accordion on
`ui.AccordionBase`; `FocusRingOffset` for inset rings; native mask fades for
ScrollFade; one public `Orientation` (breaking); per-size button gaps, sheet
inward borders, toast metrics and calendar band corners.

## Risks

Headless paint tests cover shadow layers and popup geometry, not GPU blur
quality. Exit motion on menus that previously had none needs a native check.
