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

The MIT coss styling guide is the semantic reference. Go names preserve its
roles rather than reusing the previous design system's meanings:

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
color and chooses a contrasting foreground. It does not replace the coss
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
The coss web button sizes grow by 4 px on its smaller-screen branch. The
native library uses the desktop size scale; native window layout and
`theme.Scale` handle host text scaling instead of CSS breakpoints.

`FontFamily`, `HeadingFontFamily` and `MonoFontFamily` are independent host
preferences. Heading falls back to body when omitted, and each appends native
fallbacks. `theme.Scale(c, dip)` incorporates OS text scale and `FontSize`.
Font family configuration does not download or register font data.
[fonts.md](fonts.md) records the independent OFL sources and registration
boundary.

Button feedback follows the coss registry's fills, borders and shadow changes
without a pressed scale. Disabled controls retain MyGo's native 50% opacity;
the coss button source uses 64%. These are explicit native choices, not
pixel-parity claims.

## Token provenance

Reference: coss ui
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
