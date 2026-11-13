# Theme customization

Call `theme.Apply(c)` before constructing a view's children. Components read
`theme.Colors(c)` for the current appearance. `ApplyMode(c, dark)` selects
an explicit appearance for previews.

`theme.Configure` publishes process-wide preferences and both palettes.
It replaces the whole preference value. Use `theme.UpdatePreferences` to
change selected fields without discarding an existing palette. Apply changes
on the UI thread and invalidate the affected windows; application and
widget state remain window-owned.

## Semantic roles

Root [DESIGN.md](../DESIGN.md) records what each role looks like and where it
is used. The role names came from the MIT coss styling guide (see Token
provenance); the published CSS names are kept for orientation:

| Published CSS role | Canonical Go role | Use |
| --- | --- | --- |
| `--background`, `--foreground` | `Background`, `Foreground` | Page and ordinary text |
| `--card`, `--card-foreground` | `Card`, `CardForeground` | Contained content surfaces |
| `--popover`, `--popover-foreground` | `Popover`, `PopoverForeground` | Floating panels |
| `--primary`, `--primary-foreground` | `Primary`, `PrimaryForeground` | Main actions and their text |
| `--secondary`, `--secondary-foreground` | `Secondary`, `SecondaryForeground` | Secondary actions |
| `--muted`, `--muted-foreground` | `MutedBackground`, `MutedForeground` | Low-emphasis surface and text |
| `--accent`, `--accent-foreground` | `Accent`, `AccentForeground` | Subtle hover/selection surfaces |
| `--destructive`, `--destructive-foreground` | `Destructive`, `DestructiveForeground` | Destructive action and readable status text |
| `--info`, `--success`, `--warning` and their foregrounds | `Info`, `Success`, `Warning` and corresponding `*Foreground` | Status roles |
| `--border`, `--input`, `--ring` | `Border`, `Input`, `Ring` | Translucent boundaries, input edges and keyboard focus |
| `--sidebar` and `--sidebar-foreground` | `Sidebar`, `SidebarForeground` | Navigation surface and text |
| `--sidebar-primary`, `--sidebar-accent` and foregrounds | `SidebarPrimary`, `SidebarAccent` and corresponding `*Foreground` | Navigation actions |
| `--sidebar-border`, `--sidebar-ring` | `SidebarBorder`, `SidebarRing` | Navigation boundary and focus |
| Chart colors | `Chart1` through `Chart5` | Host data visualization colors; no chart component is implied |
| Code surface roles | `Code`, `CodeForeground`, `CodeHighlight` | Code-oriented presentation |

Every canonical Go field in this table has a matching pointer field in
`PalettePatch`. Sidebar/chart/code roles are available even where a catalog
page or constructor is absent. The sidebar/chart/code extensions must not be
mistaken for copied AGPL package implementation.

`Palette` also contains derived component conveniences such as `Surface`,
`Overlay`, `Muted`, `Danger`, soft/hover colors and field colors.
`Muted` means foreground text; prefer the explicit canonical names in new
consumer code. `PalettePatch` retains a smaller convenience override set for
those existing component families. It is not a second design authority.

## Override each appearance

```go
primary := ui.RGB(30, 90, 180)
onPrimary := ui.RGB(255, 255, 255)
darkPrimary := ui.RGB(150, 190, 255)
onDarkPrimary := ui.RGB(15, 30, 55)
transparent := ui.Transparent

theme.Configure(theme.Preferences{
    Palette: theme.PaletteOverrides{
        Light: &theme.PalettePatch{
            Primary: &primary,
            PrimaryForeground: &onPrimary,
            Border: &transparent,
        },
        Dark: &theme.PalettePatch{
            Primary: &darkPrimary,
            PrimaryForeground: &onDarkPrimary,
        },
    },
})
```

Nil patch fields inherit defaults. A pointer to a transparent color is an
intentional override. Set foreground pairs explicitly when changing a role's
contrast. Translucent `Secondary`, `MutedBackground`, `Accent`, `Border`
and `Input` composite over the actual surface; flattening them early changes
their behavior on other surfaces.

`Preferences.Accent` is a convenience that replaces the **Primary** action
color and chooses a contrasting foreground. It does not replace the
**Accent** interaction role. A canonical palette patch takes precedence.

```go
theme.UpdatePreferences(func(p *theme.Preferences) {
    p.Radius = 12
    p.FontSize = 18
})
```

The native theme installed by `Apply` derives its actions, focus and selection
from the resolved palette. A host that supplies an additional scoped native
theme must account for that installation order.

## Geometry and type

| Setting | Default native value |
| --- | --- |
| `FontSize` | 16 DIP base; control text commonly 14/20 |
| Spacing unit | 4 DIP |
| `Radius` | 10 DIP |
| Radius XS / SM / MD / LG / XL | 4 / 6 / 8 / 10 / 14 DIP |
| Radius XXL / XXXL / XXXXL / Field | 16 / 24 / 32 / 8 DIP |
| Button xs / sm / default / lg / xl | 24 / 28 / 32 / 36 / 40 DIP |

Radii are offsets from the configured base, clamped at zero where needed.
The native library uses the desktop size scale at every window width; native
window layout and `theme.Scale` handle host text scaling instead of CSS
breakpoints. (The coss web source grows button sizes by 4 px on smaller
screens; that branch was not carried over.)

`FontFamily`, `HeadingFontFamily` and `MonoFontFamily` are independent host
preferences. Heading falls back to body when omitted, and each appends native
fallbacks. `theme.Scale(c, dip)` incorporates OS text scale and `FontSize`.
Font family configuration does not download or register font data.
[fonts.md](fonts.md) records the independent OFL sources and registration
boundary.

Button feedback changes fill, border and shadow without a pressed scale.
Disabled controls use MyGo's native 50% opacity, as `DESIGN.md` records
(the coss button source this was derived from uses 64%).

`SurfaceShadows` is the card shadow with its edge, `FieldShadows` the input
variant, `OverlayShadows` the popup and dialog shadow and `TooltipShadows`
the tooltip shadow; `DESIGN.md` lists their layers. Every layer set ends
with a 1px inset edge lit from above: black along the bottom inside edge in
light appearance, white 6% along the top inside edge in dark, as the coss
`before:` edge draws it. Each call returns a copy. The layers were taken from the coss
registry classes shadow-xs/5, shadow-lg/5 and shadow-md/5 and its `before:`
edge.

## A whole theme

`theme.Theme` is the look as one value: both palettes, the base radius, the
type roles, the spacing unit, the shadow layers, motion and effects.
`theme.Default()` is the built-in look. Start from it, change what differs
and publish it with `theme.Use`; everything left alone keeps its built-in
value.

```go
custom := theme.Default()
custom.Light.Primary = ui.RGB(30, 90, 180)
custom.Light.PrimaryForeground = ui.RGB(255, 255, 255)
custom.Radius = 6
custom.Spacing.Unit = 5
custom.Motion.Dialog = 260 * time.Millisecond
custom.Effects.BackdropBlur = 8
if err := theme.Use(custom); err != nil {
    return err
}
```

`Use` is process-wide, like `Configure`: call it on the UI thread, before the
first frame or followed by an invalidation of the open windows. It copies the
value, so later edits to `custom` change nothing until the next `Use`;
`theme.Current()` returns a copy of what is published. A nil easing means
the built-in curve.

`Use` publishes nothing and returns an error naming the field when a motion
duration or `Effects.LoadingDelay` is negative, or when
`Effects.SpinnerPeriod` or `Effects.SweepPeriod` is under a millisecond. A
zero motion duration is valid and means at once: the change lands on the
frame after it, with no animation (`theme.Motion()` reports it as one
nanosecond, the shortest duration MyGo animates; a zero `ui.Animation`
would select MyGo's default spring). A zero or negative
`Effects.ShimmerPeriod` keeps loading placeholders still.

A palette set on `Theme` is used as written. Changing `Primary` or `Ring`
does not recompute the roles the built-in palette derives from them:
`Selection`, `Focus`, `FieldBorderFocus`, `Link` and the hover roles
(`DefaultHover`, `AccentHover`, `SuccessHover`, `WarningHover`,
`DangerHover`, `InfoHover`) stay as given, so set them too.
`Effects.HoverOpacity` governs only the hovers components derive as they
draw (buttons, status solids, the toast action, native controls); the
palette's hover roles, which tags and breadcrumbs read, are not recomputed
from it.

`Preferences` still apply on top of the published theme, so a user's
accent, radius or font size keeps working over a custom theme, and
`Preferences.Accent` does recompute the roles derived from the primary
color. A `FontSize` or `Radius` preference equal to the published theme's
value means inherit, as zero does: `Configure(PreferencesSnapshot())` pins
neither, and a theme published later still sets them.

### Motion, spacing and effect tokens

Components read these on every build and paint. `DESIGN.md` records which
rule each one governs.

| Token | Default | Governs |
| --- | --- | --- |
| `Spacing.Unit` | 4 DIP | One step of `theme.Space(n)` |
| `Motion.State` | 150 ms | A control changing state in place |
| `Motion.Reveal` | 200 ms | Content opening; a selection mark sliding |
| `Motion.Roll` | 300 ms | A value rolling to the next one |
| `Motion.Shift` | 150 ms | Neighbours making room during a drag |
| `Motion.Page` | 250 ms | One page replacing another |
| `Motion.Popup` | 150 ms | Menus, popovers, tooltips |
| `Motion.Dialog` | 200 ms | Dialogs and their backdrop |
| `Motion.Drawer` | 450 ms | Drawers and sheets |
| `Motion.Ease` | ease-out | The default curve |
| `Motion.EaseTravel` | ease-in-out | Marks and surfaces that travel |
| `Motion.EaseDrawer` | cubic-bezier(0.32, 0.72, 0, 1) | Drawers and sheets |
| `Effects.HoverOpacity` | 0.9 | A solid fill under the pointer |
| `Effects.PopupScale` | 0.98 | The scale a surface arrives from |
| `Effects.BackdropBlur` | 4 DIP | Blur behind a modal surface |
| `Effects.FadeEdge` | 24 DIP | Depth of a scroll edge fade, a mask and a reveal's front |
| `Effects.ShimmerPeriod` | 2 s | One loop of a skeleton or shimmer; zero or less keeps them still |
| `Effects.ShimmerPeak`, `Effects.ShimmerPeakDark` | 0.64, 0.04 | How far the shimmer band lightens, per appearance |
| `Effects.SpinnerPeriod` | 1 s | One turn of a spinner |
| `Effects.SweepPeriod`, `Effects.SweepBand` | 1.5 s, 0.4 | The indeterminate bar |
| `Effects.LoadingDelay` | 150 ms | Wait before a region shows loading |
| `Effects.RefreshOpacity` | 0.55 | Content while it refreshes |
| `Effects.ThumbPressScale` | 1.2 | A slider or colour-area thumb held by the pointer |

Custom drawing reads the same values with `theme.Motion()`,
`theme.Effects()` and `theme.Space(n)`. The `theme.SpaceN` constants are the
built-in steps and do not follow `Spacing.Unit`.

A token changes a duration, never whether something moves: reduced motion
keeps its meaning under any theme.

Some timings belong to one component and are not tokens: the button press,
the hover fades of the link icon and the color swatch, the switch track, the calendar view swap, the progress
value, the toast spring, the status dot loops.

Limit: content that MyGo's `ui.CollapsibleBase` opens (accordion,
disclosure, settings expander, timeline, about dialog) keeps the toolkit's
fixed 200 ms ease-out and does not follow `Motion.Reveal`.

## Token provenance

The values are now recorded in `DESIGN.md`; this section keeps where they
came from. Origin: coss ui
`3d8ca28e7d90123e9e0f44665944dc240b026ccb`,
[`apps/ui/content/docs/(root)/styling.mdx`](https://github.com/cosscom/ui/blob/3d8ca28e7d90123e9e0f44665944dc240b026ccb/apps/ui/content/docs/%28root%29/styling.mdx).
That MIT page publishes the role system, CSS values, radius and font guidance.
Component metrics come from the MIT registry, especially
[`button.tsx`](https://github.com/cosscom/ui/blob/3d8ca28e7d90123e9e0f44665944dc240b026ccb/apps/ui/registry/default/ui/button.tsx).

The docs page and the site's runtime do not have identical dark mixes.
The currently selected native values follow the observed runtime mixes:

| Dark role | MIT styling page | Selected native/runtime mix |
| --- | --- | --- |
| Background | neutral-950 95%, white 5% | neutral-950 96%, white 4% |
| Card | background 98%, white 2% | background 98%, white 2% |
| Popover | background 98%, white 2% | background 96%, white 4% |
| Sidebar | neutral-950 97%, white 3% | neutral-950 97%, white 3% |

This records a deliberate numeric discrepancy, not an assertion that the MIT
page contains the runtime values. `apps/ui/app/globals.css` imports
`@coss/ui/globals.css`; its resolved implementation is under
`packages/ui/src/styles/globals.css`, outside the MIT subtree. Do not copy
that package source or describe it as MIT. Use the published MIT roles and
registry for the native implementation, and retain this separate observation
when comparing screenshots. See [references.md](references.md) for the
licensing boundary.
