package mygoui

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// CloseButtonOptions configures the compact dismiss control.
type CloseButtonOptions struct {
	Label    string
	Disabled bool
	OnClick  func()
}

// CloseButton is a 24 DIP HeroUI dismiss button with a native focus ring.
func CloseButton(c *ui.Context, o CloseButtonOptions) ui.Element {
	t := theme.Colors(c)
	label := o.Label
	if label == "" {
		label = "Close"
	}
	e := ui.ButtonBase(c).Size(theme.Scale(c, theme.Space24), theme.Scale(c, theme.Space24)).
		Center().Shrink(0).Padding(theme.Space4).Radius(theme.RadiusScale().XL).
		FocusRing(true).Label(label)
	e = style.Background(e, theme.RadiusScale().XL, t.Default, t.DefaultHover, t.DefaultHover)
	e.TextColor(t.Muted)
	if o.OnClick != nil && !o.Disabled {
		e.OnClick(o.OnClick)
	}
	e.Children(func() {
		ui.Icon(c, lucide.X()).Size(theme.Scale(c, theme.IconMD), theme.Scale(c, theme.IconMD)).Shrink(0)
	})
	return style.Disabled(e, o.Disabled)
}

// LinkOptions configures a native action styled as a HeroUI link.
type LinkOptions struct {
	Label    string
	Icon     *ui.SVG
	Disabled bool
	OnClick  func()
	Tooltip  string
}

// Link is an inline, keyboard-focusable action. HeroUI's hover underline is
// painted over the native text so the static layout does not shift. URL
// opening remains application work; the library performs no I/O.
func Link(c *ui.Context, o LinkOptions) ui.Element {
	t := theme.Colors(c)
	e := ui.ButtonBase(c).Gap(theme.Space4).Padding(0).Radius(theme.RadiusScale().MD).
		AlignSelf(ui.Start).FocusRing(true).Label(o.Label)
	e = e.TextColor(t.Link)
	if o.OnClick != nil && !o.Disabled {
		e.OnClick(o.OnClick)
	}
	if o.Tooltip != "" {
		e.Tooltip(o.Tooltip)
	}
	e.Children(func() {
		ts := theme.SM()
		ts.Weight = theme.WeightMedium
		style.Text(ui.Text(c, o.Label), c, ts).TextColor(t.Link)
		if o.Icon != nil {
			ui.Icon(c, o.Icon).Size(theme.Scale(c, theme.Space12), theme.Scale(c, theme.Space12)).Shrink(0)
		}
	})
	if !o.Disabled {
		// Snapshot the states for the paint: a retained copy of a leaving
		// element paints in later passes, where this Element has expired.
		underline := e.Hovered() || e.FocusVisible()
		if underline {
			e.DrawOver(func(p *ui.Painter, r ui.Rect) {
				p.Line(r.X, r.Y+r.H-1, r.X+r.W, r.Y+r.H-1, 1.5, t.Muted.Alpha(0.5))
			})
		}
	}
	return style.Disabled(e, o.Disabled)
}

// LinkButtonOptions configures native HeroUI's link-style button.
type LinkButtonOptions struct {
	Label    string
	Icon     *ui.SVG
	Disabled bool
	OnClick  func()
}

// LinkButton is a ghost text action with no hover background, matching the
// native HeroUI variant used for inline navigation.
func LinkButton(c *ui.Context, o LinkButtonOptions) ui.Element {
	t := theme.Colors(c)
	e := ui.ButtonBase(c).Gap(theme.Space4).PaddingX(theme.Space4).PaddingY(0).
		MinHeight(theme.Scale(c, theme.HeightButtonSM)).Radius(theme.RadiusScale().LG).
		AlignSelf(ui.Start).FocusRing(true).Label(o.Label).TextColor(t.DefaultForeground)
	if o.OnClick != nil && !o.Disabled {
		e.OnClick(o.OnClick)
	}
	e.Children(func() {
		ts := theme.SM()
		ts.Weight = theme.WeightMedium
		style.Text(ui.Text(c, o.Label), c, ts)
		if o.Icon != nil {
			ui.Icon(c, o.Icon).Size(theme.Scale(c, theme.IconMD), theme.Scale(c, theme.IconMD)).Shrink(0)
		}
	})
	return style.Disabled(e, o.Disabled)
}
