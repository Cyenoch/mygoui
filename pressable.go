package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// PressableFeedbackOptions configures native HeroUI's press-feedback container.
type PressableFeedbackOptions struct {
	Label     string
	Disabled  bool
	FullWidth bool
	OnClick   func()
}

// PressableFeedback provides a native press highlight. MyGo does not expose a
// live layout transform, so HeroUI Native's scale animation and ripple are not
// ported; the press and hover surfaces remain.
func PressableFeedback(c *ui.Context, o PressableFeedbackOptions, body func()) ui.Element {
	t := theme.Colors(c)
	radius := theme.RadiusScale().XL
	e := ui.ButtonBase(c).MinWidth(0).Radius(radius).Label(o.Label)
	if o.FullWidth {
		e.FillWidth()
	}
	if o.OnClick != nil && !o.Disabled {
		e.OnClick(o.OnClick)
	}
	e = style.Background(e, radius, ui.Transparent, t.Default.Alpha(0.6), t.Default)
	e.Children(body)
	return style.Disabled(e, o.Disabled)
}
