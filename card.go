package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// SurfaceVariant is HeroUI's container elevation token.
type SurfaceVariant uint8

const (
	SurfaceDefault SurfaceVariant = iota
	SurfaceSecondary
	SurfaceTertiary
	SurfaceTransparent
)

// SurfaceOptions customizes a HeroUI Surface. Padding and Radius default to
// zero because HeroUI's Surface is a style boundary, not a layout container.
type SurfaceOptions struct {
	Variant   SurfaceVariant
	Padding   float32
	Radius    float32
	FillWidth bool
	Label     string
}

// Surface paints a HeroUI surface and exposes the matching foreground.
func Surface(c *ui.Context, o SurfaceOptions, body func()) ui.Element {
	t := theme.Colors(c)
	e := ui.Column(c).MinWidth(0)
	if o.FillWidth {
		e.FillWidth()
	}
	switch o.Variant {
	case SurfaceSecondary:
		e.Background(t.SurfaceSecondary).TextColor(t.SurfaceSecondaryForeground)
	case SurfaceTertiary:
		e.Background(t.SurfaceTertiary).TextColor(t.SurfaceTertiaryForeground)
	case SurfaceTransparent:
		e.Background(ui.Transparent).TextColor(t.Foreground)
	default:
		e.Background(t.Surface).TextColor(t.SurfaceForeground)
	}
	if o.Padding > 0 {
		e.Padding(o.Padding)
	}
	if o.Radius > 0 {
		e.Radius(o.Radius).Clip()
	}
	if o.Label != "" {
		e.Label(o.Label)
	}
	return e.Children(body)
}

// CardVariant is HeroUI's card surface variant.
type CardVariant uint8

const (
	CardDefault CardVariant = iota
	CardSecondary
	CardTertiary
	CardTransparent
)

// CardOptions customizes a HeroUI Card. OnClick makes the whole card a real
// native action without changing its visual variants.
type CardOptions struct {
	Variant CardVariant
	Label   string
	OnClick func()
	Handle  *ui.Handle
}

// Card is a 16 DIP padded, 24 DIP rounded surface with a 12 DIP section gap.
func Card(c *ui.Context, o CardOptions, body func()) ui.Element {
	t := theme.Colors(c)
	radius := min(theme.Space32, theme.RadiusScale().XXXL)
	e := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).
		Padding(theme.Space16).Radius(radius).Clip()
	switch o.Variant {
	case CardSecondary:
		e.Background(t.SurfaceSecondary).TextColor(t.SurfaceSecondaryForeground)
		e = style.Shadow(e, c, theme.SurfaceShadows(c.Theme().Dark))
	case CardTertiary:
		e.Background(t.SurfaceTertiary).TextColor(t.SurfaceTertiaryForeground)
		e = style.Shadow(e, c, theme.SurfaceShadows(c.Theme().Dark))
	case CardTransparent:
		e.Background(ui.Transparent).TextColor(t.Foreground)
	default:
		e.Background(t.Surface).TextColor(t.SurfaceForeground)
		e = style.Shadow(e, c, theme.SurfaceShadows(c.Theme().Dark))
	}
	if o.OnClick != nil {
		e.Focusable().FocusRing(true).Cursor(ui.CursorPointer).OnClick(o.OnClick)
		if o.Handle != nil {
			e.Bind(o.Handle)
		}
	}
	if o.Label != "" {
		e.Label(o.Label)
	}
	return e.Children(body)
}

// CardHeader groups a title and description.
func CardHeader(c *ui.Context, body func()) ui.Element {
	return ui.Column(c).FillWidth().MinWidth(0).Children(body)
}

// CardTitle renders HeroUI's 14/24 medium card title.
func CardTitle(c *ui.Context, text string) ui.Element {
	ts := theme.SM()
	ts.LineHeight, ts.Weight = 24, theme.WeightMedium
	return style.Text(ui.Text(c, text), c, ts).TextColor(theme.Colors(c).Foreground)
}

// CardDescription renders HeroUI's 14/20 muted card description.
func CardDescription(c *ui.Context, text string) ui.Element {
	ts := theme.SM()
	ts.LineHeight = 20
	return style.Text(ui.Text(c, text), c, ts).TextColor(theme.Colors(c).Muted)
}

// CardContent grows to fill the card and uses HeroUI's 4 DIP inner gap.
func CardContent(c *ui.Context, body func()) ui.Element {
	return ui.Column(c).FillWidth().MinWidth(0).Grow(1).Gap(theme.Space4).Children(body)
}

// CardFooter aligns actions in a wrapping row.
func CardFooter(c *ui.Context, body func()) ui.Element {
	return ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Wrap().
		Gap(theme.Space8).Children(body)
}
