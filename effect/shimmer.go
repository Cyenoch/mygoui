package effect

import (
	"time"

	"github.com/Cyenoch/mygoui/internal/fx"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ShimmerOptions configures Shimmer. The zero value sweeps a highlight band.
// The loop and the band's strength are the theme's (Effects.ShimmerPeriod,
// ShimmerPeak and ShimmerPeakDark).
type ShimmerOptions struct {
	// Pulse fades the content to half and back instead of sweeping a band.
	Pulse bool
}

// Shimmer passes one highlight across e and its children while it is
// applied: use it on placeholders that stand for loading content, and stop
// applying it when the content arrives. The band is placed in the window,
// so every shimmering region moves as one.
//
// Shimmer marks e busy (Element.Busy) for assistive technology, whether
// or not the band moves; the caller names the region.
//
// Under reduced motion, high contrast or a theme period that is not
// positive, without a GPU or above the layer limit, e is drawn at rest and
// asks for no frames, as it is when Shimmer is not applied. Shimmer takes
// e's effect.
func Shimmer(c *ui.Context, e ui.Element, o ShimmerOptions) ui.Element {
	e.Busy(true)
	tokens := theme.Effects()
	if tokens.ShimmerPeriod <= 0 || theme.ReduceMotion(c) || c.Preferences().HighContrast {
		return e
	}
	peak := tokens.ShimmerPeak
	if c.Theme().Dark {
		peak = tokens.ShimmerPeakDark
	}
	s := ui.Local(e, shimmerKey{}, newShimmer)
	s.width, s.height = c.Size()
	s.period, s.peak, s.pulse = tokens.ShimmerPeriod, peak, o.Pulse
	return e.Effect(s)
}

type shimmerKey struct{}

type shimmer struct {
	shader        fx.Shimmer
	width, height float32
	period        time.Duration
	peak          float32
	pulse         bool
	// forced makes a headless frame run the shader, for tests.
	forced bool
}

func newShimmer() shimmer { return shimmer{} }

func (s *shimmer) EffectMargin() ui.Point { return ui.Point{} }

func (s *shimmer) PaintEffect(p *ui.Painter, box ui.Rect, content *ui.Bitmap, inset ui.Point) {
	if !gpu(p, s.forced) {
		fx.Plain(p, box, content, inset)
		return
	}
	s.shader.Paint(p, box, content, inset, s.width, s.height, fx.Phase(p.Now(), s.period), s.peak, s.pulse)
	p.AnimationFrameRate(60)
}
