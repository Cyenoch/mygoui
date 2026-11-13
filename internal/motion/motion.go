// Package motion holds the timing of state changes, so every control that
// changes state moves alike.
package motion

import (
	"time"

	"github.com/egoist/mygo/ui"
)

const (
	// State is a control changing state in place: a fill, a border, a mark
	// drawing in.
	State = 150 * time.Millisecond
	// Reveal is content opening or closing and an indicator moving to
	// another item.
	Reveal = 200 * time.Millisecond
	// Roll is a changing number turning its digits.
	Roll = 300 * time.Millisecond
)

// Progress returns how far e has gone from off (0) to on (1) over State.
// It moves nothing on screen by itself, so it runs under reduced motion:
// use it for an opacity, a colour mix or a mark drawing in.
func Progress(e ui.Element, key any, on bool) float32 {
	target := float32(0)
	if on {
		target = 1
	}
	return e.AnimateTo(key, target, ui.Animation{Duration: State, NonMotion: true})
}

// Color returns e's colour for key on its way to target over State.
func Color(e ui.Element, key any, target ui.Color) ui.Color {
	return e.AnimateColor(key, target, ui.Animation{Duration: State})
}
