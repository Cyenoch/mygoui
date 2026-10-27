package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// SurfaceVariant selects a native container surface.
type SurfaceVariant uint8

const (
	SurfaceDefault SurfaceVariant = iota
	SurfaceSecondary
	SurfaceTertiary
	SurfaceTransparent
)

// SurfaceOptions customizes a native Surface. Padding and Radius default to
// zero so callers can compose their own layout.
type SurfaceOptions struct {
	Variant   SurfaceVariant
	Padding   float32
	Radius    float32
	FillWidth bool
	Label     string
}

// Surface paints a theme surface and exposes the matching foreground.
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

// CardVariant selects a card surface, including native theme extensions.
type CardVariant uint8

const (
	CardDefault CardVariant = iota
	CardSecondary
	CardTertiary
	CardTransparent
)

// CardOptions customizes a coss ui Card. OnClick makes the whole card a real
// native action without changing its visual variants.
type CardOptions struct {
	Variant CardVariant
	Label   string
	OnClick func()
	Handle  *ui.Handle
}

// Card is a bordered rounded surface. Each section owns its padding.
func Card(c *ui.Context, o CardOptions, body func()) ui.Element {
	t := theme.Colors(c)
	radius := theme.RadiusScale().XXL
	e := ui.Column(c).FillWidth().MinWidth(0).Radius(radius).Clip().
		Border(theme.Hairline, t.Border)
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
	return ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space6).Padding(theme.Space24).Children(body)
}

// CardTitle renders coss ui's 18 DIP semibold heading.
func CardTitle(c *ui.Context, text string) ui.Element {
	ts := theme.LG()
	ts.LineHeight, ts.Weight, ts.Heading = 20, theme.WeightSemibold, true
	return style.Text(ui.Text(c, text), c, ts).TextColor(theme.Colors(c).Foreground)
}

// CardDescription renders coss ui's 14/20 muted card description.
func CardDescription(c *ui.Context, text string) ui.Element {
	ts := theme.SM()
	ts.LineHeight = 20
	return style.Text(ui.Text(c, text), c, ts).TextColor(theme.Colors(c).Muted)
}

// CardContent grows to fill the card, with horizontal padding and a 4 DIP gap.
// Standalone content can set vertical padding on the returned element.
func CardContent(c *ui.Context, body func()) ui.Element {
	return ui.Column(c).FillWidth().MinWidth(0).Grow(1).Gap(theme.Space4).
		Padding(0, theme.Space24).Children(body)
}

// CardFooter aligns actions in a wrapping row.
func CardFooter(c *ui.Context, body func()) ui.Element {
	return ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Wrap().
		Gap(theme.Space8).Padding(theme.Space24).Children(body)
}

// CardFrameOptions configures a framed group of cards and supporting content.
type CardFrameOptions struct {
	Label string
}

// CardFrame builds a muted, bordered container for cards and frame sections.
func CardFrame(c *ui.Context, o CardFrameOptions, body func()) ui.Element {
	t := theme.Colors(c)
	e := ui.Column(c).FillWidth().MinWidth(0).Radius(theme.RadiusScale().XXL).Clip().
		Border(theme.Hairline, t.Border).
		Background(t.MutedBackground.Alpha(0.72).Over(t.Card)).
		TextColor(t.CardForeground).Label(o.Label)
	return style.Shadow(e, c, theme.SurfaceShadows(c.Theme().Dark)).Children(body)
}

// CardFrameHeader groups frame metadata above its cards.
func CardFrameHeader(c *ui.Context, body func()) ui.Element {
	return ui.Column(c).FillWidth().MinWidth(0).Padding(theme.Space16, theme.Space24).Children(body)
}

// CardFrameTitle renders a compact semibold frame title.
func CardFrameTitle(c *ui.Context, title string) ui.Element {
	return FrameTitle(c, title)
}

// CardFrameDescription renders secondary frame text.
func CardFrameDescription(c *ui.Context, text string) ui.Element {
	return FrameDescription(c, text)
}

// CardFrameFooter hosts frame-level content below its cards.
func CardFrameFooter(c *ui.Context, body func()) ui.Element {
	return ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(theme.Space8).
		Padding(theme.Space16, theme.Space24).Children(body)
}
