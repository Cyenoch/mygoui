package effect

import (
	"github.com/Cyenoch/mygoui/internal/fx"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// RevealDirection is the way a reveal travels.
type RevealDirection uint8

const (
	// RevealDown starts at the top edge.
	RevealDown RevealDirection = iota
	// RevealUp starts at the bottom edge.
	RevealUp
	// RevealForward starts at the leading edge of the layout direction.
	RevealForward
	// RevealBackward starts at the trailing edge.
	RevealBackward
)

// RevealOptions configures Reveal. The zero value wipes downwards. The
// duration, the curve and the depth of the fading front are the theme's
// (Motion.Reveal, Motion.Ease and Effects.FadeEdge).
type RevealOptions struct {
	Direction RevealDirection
	// Fallback draws the fallback where the shader could run.
	Fallback bool
}

// Reveal wipes e and its children in behind a soft front, once, when e is
// first built: use it on content that arrives after the user asked for it,
// not on what a window opens with. Afterwards e is drawn as it is and asks
// for no frames. Building e under another key plays it again.
//
// Under reduced motion, without a GPU or with Fallback, e fades in over the
// same duration instead; above the layer limit it shows at once. While it
// plays, Reveal takes e's effect.
func Reveal(c *ui.Context, e ui.Element, o RevealOptions) ui.Element {
	if o.Direction > RevealBackward {
		panic("effect.Reveal: unsupported direction")
	}
	motion := theme.Motion()
	r := ui.Local(e, revealKey{}, newReveal)
	if !r.started {
		r.started = true
		e.Jump(revealKey{}, 0)
	}
	r.progress = e.AnimateTo(revealKey{}, 1, ui.Animation{Duration: motion.Reveal, Ease: motion.Ease, NonMotion: true})
	if r.progress >= 1 {
		return e
	}
	r.soft = theme.Effects().FadeEdge
	r.fade = o.Fallback || theme.ReduceMotion(c)
	r.dx, r.dy = 0, 1
	switch o.Direction {
	case RevealUp:
		r.dy = -1
	case RevealForward, RevealBackward:
		r.dx, r.dy = 1, 0
		if (o.Direction == RevealBackward) != (c.Direction() == ui.RightToLeft) {
			r.dx = -1
		}
	}
	return e.Effect(r)
}

type revealKey struct{}

type reveal struct {
	shader   fx.Wipe
	started  bool
	progress float32
	dx, dy   float32
	soft     float32
	fade     bool
	// forced makes a headless frame run the shader, for tests.
	forced bool
}

func newReveal() reveal { return reveal{} }

func (r *reveal) EffectMargin() ui.Point { return ui.Point{} }

func (r *reveal) PaintEffect(p *ui.Painter, box ui.Rect, content *ui.Bitmap, inset ui.Point) {
	if r.fade || !gpu(p, r.forced) {
		fx.Faded(p, box, content, inset, r.progress)
		return
	}
	r.shader.Paint(p, box, content, inset, r.dx, r.dy, r.progress, r.soft)
}
