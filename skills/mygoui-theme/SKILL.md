---
name: mygoui-theme
description: "Theming a mygoui application (github.com/Cyenoch/mygoui/theme): semantic colors, per-appearance palette overrides, a whole theme as one value, radius, spacing, motion and effect tokens, type roles, text scale and fonts. Use when customizing the look, or when custom drawing needs theme values."
---

# mygoui theme

For work on the mygoui repository itself, also read [MAINTAINING.md](MAINTAINING.md).

Look up any name with `go doc github.com/Cyenoch/mygoui/theme <Name>`. The records ship with the module version your application uses (`go list -m -f '{{.Dir}}' github.com/Cyenoch/mygoui`; when vendored, the repository at that version): `docs/theme.md` owns the semantic roles and their meaning, `DESIGN.md` the visual rules and default values, `docs/fonts.md` font sources and licensing.

## Apply and read

- Call `theme.Apply(c)` at each view root before building children; `theme.ApplyMode(c, dark)` forces one appearance, as a preview does.
- Read colors with `theme.Colors(c)` in the build that uses them, every build. The palette follows the window's current appearance, so a cached `Palette` goes stale.
- Use roles by meaning (`docs/theme.md` lists them): `Primary` is the action color, `Accent` the subtle hover and selection surface, `MutedBackground` and `MutedForeground` are separate roles, status roles (`Destructive`, `Info`, `Success`, `Warning`) express state, and text on a role uses its pair (`Card` with `CardForeground`, `Background` with `Foreground`). `Palette` also carries derived conveniences (`Surface`, `Muted`, ...); prefer the canonical roles `docs/theme.md` lists.
- `Border`, `Input`, `Accent`, `Secondary` and `MutedBackground` are translucent: paint them as they are, over the real surface.

## Customize

Preferences are process-wide; widget values and focus stay per window.

- `theme.Configure(theme.Preferences{...})` before building views. It **replaces** the whole preference value, palette included.
- `theme.UpdatePreferences(func(p *theme.Preferences) {...})` edits selected fields. Call it on the UI thread, then invalidate the affected windows.
- `Preferences.Palette` holds one `*theme.PalettePatch` per appearance (`Light`, `Dark`). Every field is a pointer: nil inherits the default, a pointer to `ui.Transparent` is an intentional override. Override a role and its foreground together, in both appearances.
- `Preferences.Accent` is a shortcut that replaces the **Primary** action color and picks a contrasting foreground; it leaves the `Accent` role alone. A palette patch wins over it.
- `Radius` sets the base radius; `theme.RadiusScale()` returns the derived scale (`SM`, `MD`, `LG`, `Field`, ...). `FontSize` sets the base text size.
- Fonts: a family name selects a font, it does not load one. Register font data at startup, then name the family. The bundled Cal Sans (`github.com/Cyenoch/mygoui/fonts`) needs `fonts.Register()`, which returns an error in browser builds.

```go
func configureTheme() error {
	brand, onBrand := ui.RGB(30, 90, 180), ui.RGB(255, 255, 255)
	brandDark, onBrandDark := ui.RGB(150, 190, 255), ui.RGB(15, 30, 55)
	theme.Configure(theme.Preferences{
		Radius: 8,
		Palette: theme.PaletteOverrides{
			Light: &theme.PalettePatch{Primary: &brand, PrimaryForeground: &onBrand},
			Dark:  &theme.PalettePatch{Primary: &brandDark, PrimaryForeground: &onBrandDark},
		},
	})
	if err := fonts.Register(); err != nil {
		return err
	}
	theme.UpdatePreferences(func(p *theme.Preferences) {
		p.FontFamily, p.HeadingFontFamily = fonts.Family, fonts.Family
	})
	return nil
}
```

## Replace the whole look

`theme.Theme` holds everything the library draws with: `Light` and `Dark` palettes, `Radius`, `Typography`, `Spacing`, `Elevation`, `Motion` and `Effects`. Use it when a brand changes more than a few colors, or changes timing or spacing; use preferences for user choices (accent, radius, font size), which keep applying on top.

- Start from `theme.Default()`, change fields, publish with `theme.Use`. Fields left alone keep the built-in look; a `Theme{}` built from nothing has no colors.
- `Use` is process-wide and copies its argument. Call it before the first frame, or on the UI thread followed by invalidating the open windows.
- `Use` returns an error naming the field, and publishes nothing, for a negative duration or a spinner or sweep period under a millisecond. A zero motion duration means at once.
- A palette on `Theme` is used as written: set a role, its foreground and the roles that should follow it yourself, in both appearances (`docs/theme.md` names them). `Effects.HoverOpacity` does not recompute the palette's hover roles.
- A `FontSize` or `Radius` preference equal to the theme's value means inherit, so restoring a `PreferencesSnapshot()` pins neither.
- `docs/theme.md` lists every motion, spacing and effect token with its default and what it governs.

```go
func brandTheme() error {
	t := theme.Default()
	t.Light.Primary, t.Light.PrimaryForeground = ui.RGB(30, 90, 180), ui.RGB(255, 255, 255)
	t.Dark.Primary, t.Dark.PrimaryForeground = ui.RGB(150, 190, 255), ui.RGB(15, 30, 55)
	t.Radius = 6
	t.Spacing.Unit = 5
	t.Motion.Dialog = 260 * time.Millisecond
	return theme.Use(t)
}
```

## Sizes in custom UI

Theme metrics (spacing, type roles, control heights, shadows; `go doc github.com/Cyenoch/mygoui/theme` lists them) are unscaled DIPs.

- Space with `theme.Space(n)`, n steps of the theme's unit (`theme.Space(2)` is 8 DIP by default), so custom UI follows a theme that changes the unit. The `theme.SpaceN` constants are fixed built-in steps.

- Convert text, line-box and icon sizes with `theme.Scale(c, dip)`, which applies the OS text scale and `FontSize`.
- For ordinary text, `mygoui.Typography` and `mygoui.Heading` already apply a type role and its accessibility; read a role's `TextStyle` only for custom drawing.
- `theme.ReduceMotion(c)` reports the OS request for less motion; the `mygoui-motion` skill covers what to do with it.

```go
func caption(c *ui.Context, text string) ui.Element {
	role := theme.HelpText()
	return ui.Text(c, text).
		TextColor(theme.Colors(c).MutedForeground).
		FontSize(theme.Scale(c, role.Size)).
		FontWeight(role.Weight)
}
```

A host that installs its own native theme on top (`c.SetTheme`) does so after `theme.Apply`; mygoui paint reads `theme.Colors`, so the palette still comes from the preferences.
