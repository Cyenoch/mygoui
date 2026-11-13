# Design

## Context

`extract-design-language` recorded the visual language in `DESIGN.md` and found
that only color, radius, the type roles and the spacing constants were public.
Motion durations lived in `internal/motion` and as literals in component files,
effect parameters were unnamed literals, and an application could not replace
the look as one value. That change's findings F1 to F17 stay open here: this
change moves values behind names and does not reconcile components that
disagree.

`theme` publishes an immutable state through an `atomic.Pointer`; every reader
does one atomic load. `Preferences` are process-wide and layered on top.

## Goals

- Motion, spacing and effect values that several components share are public,
  typed tokens in `theme`.
- One typed value describes a whole theme; the built-in look is its default.
- Components read tokens wherever a token exists.
- The default theme renders and times exactly as before.

## Non-goals

- Resolving F1 to F17. Equal values share a token; a component that differs
  keeps its own named constant.
- Per-window themes. The theme stays process-wide like the preferences.
- Tokens for control heights, icon sizes, focus ring widths and other
  per-component geometry.

## Decision: one `Theme` value published with `Use`

```go
type Theme struct {
    Light, Dark Palette
    Radius      float32
    Typography  Typography
    Spacing     Spacing
    Elevation   Elevation
    Motion      MotionTokens
    Effects     EffectTokens
}

func Default() Theme           // the built-in look, a fresh copy
func Use(t Theme) error        // publish t; preferences keep applying on top
func Current() Theme           // a copy of the published theme
func Motion() MotionTokens     // paint-path read, one atomic load
func Effects() EffectTokens    // paint-path read, one atomic load
```

An application starts from `Default()`, changes fields and calls `Use`. A
partial override inherits the rest because the starting value is complete, so
no "unset" encoding is needed for numbers. The three easings are the one
exception: `ui.Easing` is a function, a zero value is common in a hand-built
`Theme`, so a nil easing means the built-in curve.

`Use` clones the slices it is given (shadow layers) and `Current` and the
shadow functions return clones, so neither side can mutate the published
state. `Use` republishes on every call: a `Theme` holds functions and cannot
be compared, and one publish is cheap.

`Preferences` resolve over the published theme. `FontSize`, `Radius` and the
font families replace the theme's values when set. Palette preferences
(`Accent`, patches, the convenience fields) are applied to the theme's palette
by diffing the built-in derivation before and after the preference: roles the
preference moves follow it, the others keep the theme's color. For the default
theme this is the previous computation, so `Configure`, `UpdatePreferences`,
`PreferencesSnapshot` and `PalettePatch` behave as before. Nothing is
**BREAKING**.

`theme.Space(n)` returns n spacing steps of the published unit. The
`theme.SpaceN` constants remain as the built-in steps for callers that need a
constant; library code calls `Space`.

### Alternatives considered

- **More fields on `Preferences`.** Rejected: `Preferences` is a sparse,
  comparable, persistable record of user choices with zero meaning "default".
  Durations and curves are not user preferences, a zero duration is a valid
  value, and functions would make it incomparable.
- **Functional options (`theme.With...`).** Rejected: a second configuration
  style beside the existing struct-based one, and it cannot be read back as a
  value.
- **A sparse patch type with pointers for every token.** Rejected: twice the
  surface for no gain over "start from `Default()`".
- **Per-window theme carried in `ui.Context`.** Rejected: paint closures and
  package helpers read tokens without a context, and the existing preferences
  are already process-wide.

### `internal/motion`

It keeps `Progress` and `Color`, the two helpers 27 files share to animate a
state in place; they read `Motion.State` and `Motion.Ease`. Its duration
constants are gone. Deleting the package would repeat the same four lines in
every caller.

## Tokens

Defaults are the values the code used before this change.

### Motion (`theme.MotionTokens`)

| Token | Default | Readers |
| --- | --- | --- |
| `State` | 150 ms | `internal/motion` (checkbox, radio, choice card, chip, tag, toggle, menu rows, pagination, steps, sidebar rows, form field chrome, status dot, statistic, timeline, upload, OTP, hotkey input, password field, image select, tree rows, document tabs); switch thumb; slider thumb press; rate; tags input; list box; carousel pips; color area and swatch check; log viewer jump button; image viewer chrome; root fade and appear transitions (alert, avatar, image, timeline, icon swap); `internal/busy` crossfade |
| `Reveal` | 200 ms | tabs indicator and segmented thumb; slider fill; tree rows and arrow; sidebar width; resizable panel; split view; spoiler; banner exit; action bar; carousel pull-back; tour scrim |
| `Roll` | 300 ms | number flip; tour spotlight and crossfade |
| `Shift` | 150 ms | sortable list; document tabs |
| `Page` | 250 ms | carousel; navigation stack |
| `Popup` | 150 ms | `internal/anchored` (menu, popover, tooltip, select, combo box, date picker, preview card, teaching tip); floating panel |
| `Dialog` | 200 ms | dialog, alert dialog, modal, image viewer, backdrop |
| `Drawer` | 450 ms | drawer, sheet, bottom sheet, floating sheet |
| `Ease` | `ui.EaseOut` | every reader above unless listed below |
| `EaseTravel` | `ui.EaseInOut` | tabs indicator and segmented thumb; dialog and backdrop; tour spotlight |
| `EaseDrawer` | `ui.CubicBezier(0.32, 0.72, 0, 1)` | drawer, sheet, bottom sheet, floating sheet |

### Effects (`theme.EffectTokens`)

| Token | Default | Readers |
| --- | --- | --- |
| `HoverOpacity` | 0.9 | hovers derived at draw time: native control hover, button, every status solid, toast action; and the palette roles a preference re-derives. Not the hover roles of a theme's own palettes |
| `PopupScale` | 0.98 | anchored surfaces, dialog |
| `BackdropBlur` | 4 | overlay backdrop |
| `FadeEdge` | 24 | scroll fade default depth, document tabs strip |
| `SpinnerPeriod` | 1 s | spinner, busy spinner |
| `SweepPeriod` | 1.5 s | indeterminate progress bar, busy line |
| `SweepBand` | 0.4 | indeterminate progress bar, busy line |
| `LoadingDelay` | 150 ms | `internal/busy` (button, input, select, table, list, tree, loading overlay) |
| `RefreshOpacity` | 0.55 | `internal/busy` refreshing content |

### Other layers

| Token | Default | Readers |
| --- | --- | --- |
| `Spacing.Unit` | 4 | `theme.Space(n)`, used by every component package |
| `Radius` | 10 | `theme.RadiusBase`, `theme.RadiusScale` |
| `Typography.Size` and roles | 16; the existing role table | the type role functions |
| `Elevation.Light`, `Elevation.Dark` | the existing shadow layers | `SurfaceShadows`, `FieldShadows`, `OverlayShadows`, `TooltipShadows` |
| `Light`, `Dark` | the built-in palettes | `theme.Colors`, `theme.ColorsForMode` |

## Values that stay in the component

Each is a named constant or a documented literal in its file. A value stays
when one component uses it, when it disagrees with its neighbours and the
disagreement is an open finding, or when it is behavior rather than look.

| Where | Value | Why |
| --- | --- | --- |
| `button.go` `buttonPressDuration` | 250 ms | Only the button press uses it |
| `pressable.go` `pressableHighlightDuration` | 200 ms | Hover fade differs from link and swatch (F2, F3) |
| `link.go` `closeButtonHoverDuration`, `linkIconDuration` | 100 ms, 150 ms | Same disagreement (F2, F3) |
| `color/color_swatch.go` `swatchHoverDuration` | 100 ms | F3 |
| `resizable_panel.go` `resizeHandleFade` | 100 ms | F2 |
| `form/switch.go` `switchTrackFade`, `switchStretch` | 200 ms, 100 ms | Track and thumb disagree (F2, F5) |
| `date/calendar.go` `calendarYearFade`, `calendarDayFade`, `calendarHeaderFade` | 200, 120, 150 ms | Three durations in one component (F7) |
| `progress.go` `progressDuration` | 500 ms | Only the determinate bar uses it |
| `progress.go`, `internal/busy` sweep shape | `ui.EaseInOut(phase)`, travel 4.5 and 4 bands | Shape of the sweep; the two differ (F8) |
| `spoiler.go` | `ui.EaseInOut` | Differs from other reveals (F1) |
| `sidebar.go` | `ui.Linear` | Width follows linearly (F4) |
| `list/tree.go` arrow | `ui.EaseInOut` | Differs from the row it belongs to (F12) |
| `banner.go` exit | `ui.EaseInOut` | Only the banner exit uses it |
| `overlay/toast.go` `toastSpring`, content fade | 320 ms spring, 200 ms | Only the toast stack uses a spring |
| `overlay/tour.go` crossfade | `ui.Linear` | A crossfade must be linear to keep total opacity |
| `skeleton.go`, `status_dot.go` | 2 s, 1.6 s periods and their phase shapes | Each loop belongs to one component |
| `form/input_otp.go` | caret blink phase | Follows the native caret |
| `content/copy_button.go` `copiedFeedback` | 1.5 s | Behavior: how long "copied" shows |
| `list/tree.go` `typeWindow`, `revealSlack` | 700 ms, 100 ms | Behavior: type-ahead window, reveal tolerance |
| `navigation/menubar.go` `altTapLimit` | 600 ms | Behavior: Alt tap recognition |
| `navigation/menu_ahead.go` `menuAheadQuiet` | 1 s | Behavior: submenu aim grace |
| `navigation/carousel.go` `carouselWheelQuiet` | 150 ms | Behavior: wheel gesture end |
| `overlay/preview_card.go` | 150 ms leave | Behavior: hover leave grace |
| Field vertical padding (`form/input.go`, `input_group.go`, `number_field.go`, `password_field.go`, `search_field.go`, `editable_text.go`, `tags_input.go`, `content/search_field.go`, `log_toolbar.go`) | 5, 3 and 2 DIP | Centres a 20 DIP line in a fixed control height; follows the height, not the grid |
| `navigation/menu.go` indicator inset | 28 DIP | Clears a fixed-size check indicator |
| Control heights, icon sizes, focus ring widths, dialog rise and maximum width, spoiler cut, frame rates | per file | Per-component geometry, outside this change |

The remaining `time.Millisecond`, `time.Second` and easing literals in
non-test library code outside `internal/gallery`, `examples` and
`theme/language.go` are exactly the rows above.

## Decisions from review

- **`Use` returns an error.** A theme is application input. `Use` refuses a
  negative motion duration or loading delay and a spinner or sweep period
  under a millisecond, names the field and publishes nothing: components
  divide by the periods, and MyGo panics on a negative duration. An error
  rather than a panic follows the project's explicit-error rule.
- **Zero motion means at once.** MyGo has no at-once value for `AnimateTo`
  and `AnimateColor`: a zero `Duration` selects its default spring and a
  negative one panics (`ui/motion.go`); only `ElementTransition` takes a
  negative duration as instant. `Motion()` therefore reports a zero token as
  one nanosecond, the shortest curve MyGo runs, so every reader lands on the
  frame after the change without each call site branching. `Current()`
  returns the value as written. `Effects.ShimmerPeriod` keeps the effect
  family's meaning: zero or less keeps placeholders still.
- **Equal to the theme means inherit.** `PreferencesSnapshot` returns the
  font size and radius in effect, so `Configure` stores a value equal to the
  published theme's as unset. `Configure(PreferencesSnapshot())` is again a
  no-op, as in 0.2.0, and `UpdatePreferences` needs no rule of its own. A
  host that wants the theme's current value pinned against a later theme
  cannot express it; none does today.
- **A theme's palette is used as written.** `HoverOpacity` does not
  re-derive the palette's hover roles (`DefaultHover`, `AccentHover`, the
  status hovers); a theme author owns them. It governs hovers computed at
  draw time, and `statusColors` now derives all six solids from it instead
  of mixing the token with palette roles. With the default theme the result
  is identical. Tags and breadcrumbs read `DefaultHover` and therefore stay
  at the palette's value. Re-deriving was rejected: it would overwrite
  colours the theme supplied.
- **A narrow sweep band leaves the track.** The band travels
  `max(4.5 × band, track + band)` (progress) and
  `max(4 × band, track + band/2)` (busy line). At the default 0.4 the first
  term is larger, so the default output is unchanged; below about 0.29 the
  second term carries the band past the far edge instead of jumping back.
  Starting at `-band` with travel `1 + band` for every band was rejected
  because it changes the default theme's positions.
- **Grid spacing.** The tabs, segmented control and breadcrumbs gaps and
  paddings that were bare 2 and 10 DIP now read `theme.Space`. `layout/` and
  `effect/` already do.

## Risks

- **MyGo's collapsible is not tokenizable.** `ui.CollapsibleBase` animates at a
  fixed 200 ms ease-out inside the toolkit. Accordion, disclosure, the
  settings expander, the timeline and the about dialog therefore do not follow
  `Motion.Reveal`. They match the default; a theme that changes `Reveal`
  leaves them at 200 ms. Fixing it needs a MyGo option or an own
  implementation and belongs with F1.
- **Folding equal values.** Slider and switch thumbs and the busy crossfade
  were separate 150 ms literals and now read `Motion.State`. They are equal
  today and all describe a state changing in place; a theme cannot set them
  apart.
- **Derived palette under a custom theme.** A preference applied over a custom
  palette moves the roles the built-in derivation would move. A theme whose
  derived roles deliberately break that derivation loses those roles when the
  matching preference is set. This is the documented meaning of the
  preference.
- **Process-wide state in tests.** A test that calls `Use` must restore
  `Default()`; the new tests do so with `t.Cleanup`.
