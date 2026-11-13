// Package motion holds the two state-change helpers components share, so
// every control that changes state in place moves alike. The timing is the
// theme's: theme.Motion().
package motion

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// Progress returns how far e has gone from off (0) to on (1) over the state
// duration. It moves nothing on screen by itself, so it runs under reduced
// motion: use it for an opacity, a colour mix or a mark drawing in.
func Progress(e ui.Element, key any, on bool) float32 {
	target := float32(0)
	if on {
		target = 1
	}
	m := theme.Motion()
	return e.AnimateTo(key, target, ui.Animation{Duration: m.State, Ease: m.Ease, NonMotion: true})
}

// Color returns e's colour for key on its way to target over the state
// duration.
func Color(e ui.Element, key any, target ui.Color) ui.Color {
	m := theme.Motion()
	return e.AnimateColor(key, target, ui.Animation{Duration: m.State, Ease: m.Ease})
}
