package mygoui

import (
	"strconv"
	"strings"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// NumberFlipOptions configures a value whose digits roll when it changes.
type NumberFlipOptions struct {
	// Value is shown as given: the caller formats separators, signs and
	// units. A change rolls the digits that differ and fades the rest.
	Value string
	// Style is the text role of the value; the zero value is theme.H2.
	Style theme.TextStyle
	// Color is the ink; the zero value is the foreground.
	Color ui.Color
	// Label is the accessible name; Value when empty.
	Label string
}

type numberFlipKey struct{}

// numberFlipState is kept on the root element: the value that rolled out,
// the one showing, and the animated counter that tells how far the roll is
// from the progress value it began at to its target.
type numberFlipState struct {
	shown, prev  string
	from, target float32
	rising       bool
}

// NumberFlip shows a value that rolls its digits to a new one: columns are
// matched from the end, a digit slides out and in over 300 ms in the
// direction of the change, and other characters crossfade. Reduced motion
// crossfades only. At rest it is plain text. Assistive technology reads the
// whole value once.
func NumberFlip(c *ui.Context, o NumberFlipOptions) ui.Element {
	t := theme.Colors(c)
	text := o.Style
	if text == (theme.TextStyle{}) {
		text = theme.H2()
	}
	ink := o.Color
	if ink.A == 0 {
		ink = t.Foreground
	}
	name := o.Label
	if name == "" {
		name = o.Value
	}

	root := ui.Row(c).MinWidth(0).Role(ui.RoleText).Label(name)
	st := ui.Local(root, numberFlipKey{}, func() numberFlipState { return numberFlipState{shown: o.Value} })
	m := theme.Motion()
	anim := ui.Animation{Duration: m.Roll, Ease: m.Ease, NonMotion: true}
	if st.shown != o.Value {
		st.from = root.AnimateTo(numberFlipKey{}, st.target, anim)
		st.target++
		st.prev, st.shown = st.shown, o.Value
		st.rising = numberFlipRising(st.prev, st.shown)
	}
	progress := float32(1)
	if v := root.AnimateTo(numberFlipKey{}, st.target, anim); v < st.target {
		progress = (v - st.from) / (st.target - st.from)
	}

	glyph := func(s string) ui.Element {
		return style.Text(ui.Text(c, s), c, text).TextColor(ink).FontFeatures("tnum").AccessHidden(true)
	}
	return root.Children(func() {
		if progress >= 1 {
			glyph(st.shown).MinWidth(0).SingleLine().Selectable()
			return
		}
		reduced, line := theme.ReduceMotion(c), text.LineHeight*theme.TextScale(c)
		prev, cur := []rune(st.prev), []rune(st.shown)
		n := max(len(prev), len(cur))
		for col := range n {
			p, pok := runeFromEnd(prev, n-1-col)
			q, qok := runeFromEnd(cur, n-1-col)
			switch {
			case pok && qok && p == q:
				glyph(string(q))
			case !qok:
				glyph(string(p)).Opacity(1 - progress)
			case !pok:
				glyph(string(q)).Opacity(progress)
			default:
				roll := isDigit(p) && isDigit(q) && !reduced
				numberFlipSlot(c, glyph, line, p, q, progress, st.rising, roll)
			}
		}
	})
}

// numberFlipSlot changes one column. A rolling digit stacks the two glyphs
// in a box one line high that clips them, and slides the stack by a line; any
// other change crossfades the arriving glyph over the leaving one.
func numberFlipSlot(c *ui.Context, glyph func(string) ui.Element, line float32, from, to rune, progress float32, rising, roll bool) {
	if !roll {
		ui.Box(c).Children(func() {
			glyph(string(to)).Opacity(progress)
			glyph(string(from)).Absolute().Top(0).Left(0).Opacity(1 - progress)
		})
		return
	}
	first, second, offset := from, to, -progress*line
	if !rising {
		first, second, offset = to, from, -(1-progress)*line
	}
	ui.Box(c).ClipY().Height(line).Children(func() {
		ui.Column(c).Translate(0, offset).Children(func() {
			for _, r := range [...]rune{first, second} {
				opacity := progress
				if r == from {
					opacity = 1 - progress
				}
				glyph(string(r)).Shrink(0).Opacity(opacity)
			}
		})
	})
}

func runeFromEnd(s []rune, i int) (rune, bool) {
	if i >= len(s) {
		return 0, false
	}
	return s[len(s)-1-i], true
}

func isDigit(r rune) bool { return r >= '0' && r <= '9' }

// numberFlipRising reports whether to is the larger number, reading digits,
// a sign and a decimal point; values that are not numbers roll up.
func numberFlipRising(from, to string) bool {
	a, aok := numberFlipParse(from)
	b, bok := numberFlipParse(to)
	return !aok || !bok || b >= a
}

func numberFlipParse(s string) (float64, bool) {
	s = strings.Map(func(r rune) rune {
		if isDigit(r) || r == '.' || r == '-' {
			return r
		}
		return -1
	}, s)
	f, err := strconv.ParseFloat(s, 64)
	return f, err == nil
}
