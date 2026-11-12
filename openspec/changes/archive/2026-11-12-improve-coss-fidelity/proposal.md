# Proposal

## Why

Audits against the pinned coss registry and MyGo 0.9.0 found popups that
drift from coss elevation and motion, duplicated component code whose copies
already disagree (wrong slider decimals, error-only fields without an invalid
border, an OTP face without dark styling), and hand-written helpers that MyGo
now supplies natively.

## What Changes

- Popup and dialog elevation uses coss `shadow-lg/5` plus a 1px edge in both
  appearances; tooltips use `shadow-md/5` through a new `theme.TooltipShadows`.
- Every anchored popup (select, combo box, autocomplete, date and color
  pickers, menus, command, popover, tooltip) shares one surface and one
  anchor-facing fade/zoom with an exit.
- Fields show their invalid treatment whenever they show an error message.
- Slider value text uses the step's real decimal count.
- InputOTP uses the shared field face.
- Duplicated helpers are replaced by MyGo facilities (`Painter.FillRadii`,
  `Color.Mix`) or by one internal helper: overlay warm-up, field error text
  and status glyphs.
- Native context menus show item descriptions, shortcuts and danger styling.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `semantic-theme`: elevation tokens follow the coss popup and tooltip scales
  in both appearances.
- `native-library`: anchored popups share a surface and motion contract; form
  error state is consistent.

## Impact

Root, `theme`, `form`, `date`, `color`, `navigation`, `overlay`, `data` and
`internal/*` packages; docs/theme.md, docs/components.md, DESIGN.md and the
changelog. Adds `theme.TooltipShadows`; no public API is removed.
