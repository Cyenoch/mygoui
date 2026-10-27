package mygoui

import (
	"time"

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

// closeButtonHoverDuration is the dismiss button's background transition.
const closeButtonHoverDuration = 100 * time.Millisecond

// CloseButton is a 24 DIP dismiss button with a native focus ring and an
// animated hover face.
func CloseButton(c *ui.Context, o CloseButtonOptions) ui.Element {
	t := theme.Colors(c)
	label := o.Label
	if label == "" {
		label = "Close"
	}
	radius := theme.RadiusScale().XL
	e := ui.ButtonBase(c).Size(theme.Scale(c, theme.Space24), theme.Scale(c, theme.Space24)).
		Center().Shrink(0).Padding(theme.Space4).Radius(radius).
		FocusRing(true).Label(label)
	// Disable before the face reads the pointer so it cannot hover.
	e = style.Disabled(e, o.Disabled)
	fill := ui.Transparent
	if e.Pressed() || e.Hovered() {
		fill = t.Accent
	}
	e.Background(fill)
	// Keep the transition declared in every build: omitting it reads as the
	// control leaving and paints a fading copy over the live one. Reduced
	// motion settles this frame's colors at once.
	if theme.ReduceMotion(c) {
		c.DisableAnimations()
	}
	e = e.Transition(ui.ElementTransition{
		Duration: closeButtonHoverDuration, Ease: ui.EaseOut, Colors: true,
	})
	e.TextColor(t.Muted)
	if o.OnClick != nil && !o.Disabled {
		e.OnClick(o.OnClick)
	}
	e.Children(func() {
		ui.Icon(c, lucide.X()).Size(theme.Scale(c, theme.IconMD), theme.Scale(c, theme.IconMD)).Shrink(0)
	})
	return e
}

// LinkOptions configures an inline native action.
type LinkOptions struct {
	Label    string
	Icon     *ui.SVG
	Disabled bool
	OnClick  func()
	Tooltip  string
}

// linkIconDuration is the opacity transition for the link icon.
const linkIconDuration = 150 * time.Millisecond

// Link is an inline, keyboard-focusable action. Hover, press and keyboard
// focus underline the label, and the trailing icon fades from 60% to full
// opacity. The underline is painted over the native text so
// the static layout does not shift. URL opening remains application work;
// the library performs no I/O.
func Link(c *ui.Context, o LinkOptions) ui.Element {
	t := theme.Colors(c)
	radius := theme.RadiusScale().XL
	e := ui.ButtonBase(c).Gap(theme.Space4).Padding(0).Radius(radius).
		AlignSelf(ui.Start).FocusRing(true).Label(o.Label)
	e = style.Disabled(e, o.Disabled)
	e.TextColor(t.Link)
	if o.OnClick != nil && !o.Disabled {
		e.OnClick(o.OnClick)
	}
	if o.Tooltip != "" {
		e.Tooltip(o.Tooltip)
	}
	// Snapshot the decoration and icon fade for this build: a retained copy
	// of a leaving element paints in later passes, where these elements have
	// expired. The ink is muted on press, muted at half alpha on hover, and
	// uses the tertiary separator color for focus.
	var underline ui.Color
	active := false
	if !o.Disabled {
		switch {
		case e.Pressed():
			underline, active = t.Muted, true
		case e.Hovered():
			underline, active = t.Muted.Alpha(0.5), true
		case e.FocusVisible():
			underline, active = t.SeparatorTertiary, true
		}
	}
	underlineWidth := theme.Scale(c, 1.5)
	underlineOffset := theme.Scale(c, 1)
	iconSize := theme.Scale(c, 10.5)
	iconTarget := float32(0.6)
	if active {
		iconTarget = 1
	}
	e.Children(func() {
		ts := theme.SM()
		ts.Weight = theme.WeightMedium
		label := style.Text(ui.Text(c, o.Label), c, ts)
		if underline.A != 0 {
			ink := underline
			label.DrawOver(func(p *ui.Painter, r ui.Rect) {
				y := r.Y + r.H - underlineOffset
				p.Line(r.X, y, r.X+r.W, y, underlineWidth, ink)
			})
		}
		if o.Icon != nil {
			icon := ui.Icon(c, o.Icon).Size(iconSize, iconSize).Shrink(0)
			opacity := icon.AnimateTo("link-icon-opacity", iconTarget,
				ui.Animation{Duration: linkIconDuration, Ease: ui.EaseOut})
			_ = icon.Opacity(opacity)
		}
	})
	return e
}

// LinkButtonOptions configures a link-style native button.
type LinkButtonOptions struct {
	Label    string
	Icon     *ui.SVG
	Disabled bool
	OnClick  func()
}

// LinkButton is a text action with no hover background, padding or minimum
// height. It uses native button activation and focus.
func LinkButton(c *ui.Context, o LinkButtonOptions) ui.Element {
	t := theme.Colors(c)
	e := ui.ButtonBase(c).Gap(theme.Space8).Padding(0).
		AlignSelf(ui.Start).FocusRing(true).Label(o.Label)
	e = style.Disabled(e, o.Disabled)
	e.TextColor(t.DefaultForeground)
	if o.OnClick != nil && !o.Disabled {
		e.OnClick(o.OnClick)
	}
	iconSize := theme.Scale(c, theme.IconMD)
	e.Children(func() {
		ts := theme.SM()
		ts.Weight = theme.WeightMedium
		style.Text(ui.Text(c, o.Label), c, ts)
		if o.Icon != nil {
			ui.Icon(c, o.Icon).Size(iconSize, iconSize).Shrink(0).Margin(0, -2)
		}
	})
	return e
}
