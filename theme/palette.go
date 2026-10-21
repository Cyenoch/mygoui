package theme

import "github.com/egoist/mygo/ui"

func nativeTheme(p ColorTokens, dark bool) ui.Theme {
	return ui.Theme{
		Dark: dark, Background: p.Background, Surface: p.Card,
		SurfaceHover: p.AccentHover, SurfacePressed: p.AccentPressed,
		Border: p.Border, Text: p.Foreground, TextMuted: p.MutedForeground,
		Accent: p.Primary, AccentHover: p.PrimaryHover, AccentPressed: p.PrimaryPressed,
		AccentText: p.PrimaryForeground, Danger: p.Destructive,
		Warning: p.Warning, Success: p.Success, Selection: p.Selection, Focus: p.Ring,
		Inverse: p.Foreground, InverseText: p.Background,
		Scrollbar: p.Foreground.Alpha(0.32), ScrollbarWidth: Space6,
		Radius: RadiusCard, Spacing: SpaceUnit, FontSize: Body().Size,
	}
}

// Over composites c over the actual painted base, then attenuates that result
// by f. It consumes c's existing alpha exactly once. Use an opaque base for an
// opaque resolved paint; never multiply by c.A/255 before calling this helper.
func Over(c, base ui.Color, f float32) ui.Color { return mix(base, c.Over(base), f) }

// mix performs an sRGB mix using premultiplied alpha and returns straight alpha.
func mix(a, b ui.Color, t float32) ui.Color {
	t = max(0, min(t, 1))
	aw, bw := float32(a.A)*(1-t), float32(b.A)*t
	alpha := aw + bw
	if alpha == 0 {
		return ui.Color{}
	}
	channel := func(x, y uint8) uint8 { return uint8((float32(x)*aw+float32(y)*bw)/alpha + 0.5) }
	return ui.Color{R: channel(a.R, b.R), G: channel(a.G, b.G), B: channel(a.B, b.B), A: uint8(alpha + 0.5)}
}
