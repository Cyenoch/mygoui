package mygoui

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// fadeIn and fadeOut are the zero motions the fades share: MyGo only reads
// them, and a value shared across builds keeps a build from allocating them.
var fadeIn, fadeOut = &ui.Motion{}, &ui.Motion{}

// fadeTransition fades an element in as it appears and out as it leaves. It
// moves nothing, so it also runs under reduced motion. A transition animates
// everything it does not name, and Shadows is named only so that position,
// size and colors, which the view animates itself, are left alone.
func fadeTransition() ui.ElementTransition {
	m := theme.Motion()
	return ui.ElementTransition{
		Duration: m.State, Ease: m.Ease, Shadows: true,
		Enter: fadeIn, Exit: fadeOut,
	}
}

// appearTransition fades a message in as the caller starts building it, and
// its colors as the status changes. It has no exit: a leaving copy takes no
// room, so it would lie under the siblings that move up at once.
func appearTransition() ui.ElementTransition {
	m := theme.Motion()
	return ui.ElementTransition{
		Duration: m.State, Ease: m.Ease, Colors: true,
		Enter: fadeIn,
	}
}

// stateIcon is a glyph that stands for a state: one that changes crossfades
// with the one it replaces. The key is the glyph, which the icon packages
// hand out as one shared value.
func stateIcon(c *ui.Context, glyph *ui.SVG) ui.Element {
	return ui.Icon(c.Key(glyph), glyph).Transition(fadeTransition())
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
