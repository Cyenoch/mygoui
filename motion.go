package mygoui

import (
	"github.com/Cyenoch/mygoui/internal/motion"
	"github.com/egoist/mygo/ui"
)

// fadeTransition fades an element in as it appears and out as it leaves. It
// moves nothing, so it also runs under reduced motion. A transition animates
// everything it does not name, and Shadows is named only so that position,
// size and colors, which the view animates itself, are left alone.
var fadeTransition = ui.ElementTransition{
	Duration: motion.State, Ease: ui.EaseOut, Shadows: true,
	Enter: &ui.Motion{}, Exit: &ui.Motion{},
}

// appearTransition fades a message in as the caller starts building it, and
// its colors as the status changes. It has no exit: a leaving copy takes no
// room, so it would lie under the siblings that move up at once.
var appearTransition = ui.ElementTransition{
	Duration: motion.State, Ease: ui.EaseOut, Colors: true,
	Enter: &ui.Motion{},
}

// stateIcon is a glyph that stands for a state: one that changes crossfades
// with the one it replaces. The key is the glyph, which the icon packages
// hand out as one shared value.
func stateIcon(c *ui.Context, glyph *ui.SVG) ui.Element {
	return ui.Icon(c.Key(glyph), glyph).Transition(fadeTransition)
}

// mixState returns the colour a state change has reached: off at 0, on at 1.
func mixState(off, on ui.Color, progress float32) ui.Color {
	switch {
	case progress <= 0:
		return off
	case progress >= 1:
		return on
	}
	return off.MixOklab(on, progress)
}
