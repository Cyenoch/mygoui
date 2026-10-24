package mygoui

import (
	"time"

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

// pressableHighlightDuration is native HeroUI's press-highlight fade.
const pressableHighlightDuration = 200 * time.Millisecond

// PressableFeedback provides native press feedback: HeroUI Native's 0.985
// default press scale (without its width coefficient, which needs a measured
// box) and hover/pressed surfaces that fade between faces. MyGo has no
// ripple, so Native's ripple is not ported; the scale and faces are.
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
	// Keep the transition declared in every build: omitting it reads as
	// the element leaving and paints a fading copy over the live one.
	// Reduced motion settles this frame's colors at once.
	if theme.ReduceMotion(c) {
		c.DisableAnimations()
	}
	e = e.Transition(ui.ElementTransition{
		Duration: pressableHighlightDuration, Ease: ui.EaseOut, Colors: true,
	})
	if !o.Disabled {
		e = pressedPose(e, e.Pressed(), 0.985)
	}
	e.Children(body)
	return e
}
