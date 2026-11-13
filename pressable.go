package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// PressableFeedbackOptions configures a native press-feedback container.
type PressableFeedbackOptions struct {
	Label     string
	Disabled  bool
	FullWidth bool
	OnClick   func()
}

// PressableFeedback animates a 0.985 press scale and shows its hover and
// pressed surfaces at once. It does not draw a ripple.
func PressableFeedback(c *ui.Context, o PressableFeedbackOptions, body func()) ui.Element {
	t := theme.Colors(c)
	radius := theme.RadiusScale().XL
	e := ui.ButtonBase(c).MinWidth(0).Radius(radius).Clip().Label(o.Label)
	// Disable before the face reads its pointer state, so disabled content
	// stays at rest and takes no feedback.
	e = style.Disabled(e, o.Disabled)
	if o.FullWidth {
		e.FillWidth()
	}
	if o.OnClick != nil && !o.Disabled {
		e.OnClick(o.OnClick)
	}
	fill := ui.Transparent
	switch {
	case e.Pressed():
		fill = t.Default
	case e.Hovered():
		fill = t.Default.Alpha(0.6)
	}
	e.Background(fill)
	if !o.Disabled {
		e = pressedPose(c, e, e.Pressed(), 0.985)
	}
	e.Children(body)
	return e
}
