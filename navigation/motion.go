package navigation

import (
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// HeroUI's navigation motion, in native terms: values that move on MyGo's
// own animation channels, so an interrupted change continues from where it
// is instead of restarting, and reduced motion settles them at once.
const (
	// tabsIndicatorMotion is upstream's 250 ms ease-out-fluid slide of the
	// tab selection indicator.
	tabsIndicatorMotion = 250 * time.Millisecond
	// menuIndicatorMotion is the dot's 250 ms reveal; a checkmark draws
	// itself in a fifth as long, as its stroke-dashoffset does.
	menuIndicatorMotion = 250 * time.Millisecond
	menuCheckMotion     = 100 * time.Millisecond
	// panelEntryMotion is how long a floating panel fades in. Upstream has
	// no exit on these panels: closing removes them with the popover.
	panelEntryMotion = 150 * time.Millisecond
)

// slideMotion returns HeroUI's indicator animation over d.
func slideMotion(d time.Duration) ui.Animation {
	return ui.Animation{Duration: d, Ease: ui.EaseOut}
}

// fadePanel fades a collapsible's growing panel with its progress, so the
// body fades in under the clipped, growing box as upstream's opacity does.
// A fully open panel keeps its default opacity, so no layer state is set
// once it settles.
func fadePanel(panel ui.Element, progress float32) ui.Element {
	if progress >= 1 {
		return panel
	}
	return panel.Opacity(max(progress, 0))
}

// entryFade fades a floating panel in, as HeroUI's popovers arrive. Only the
// arrival fades: MyGo cannot keep an inert copy of these panels after they
// are closed, and an overlay relayout would discard a positional entry.
// Reduced motion shows the panel at once.
func entryFade(c *ui.Context, e ui.Element, d time.Duration) ui.Element {
	if theme.ReduceMotion(c) {
		return e
	}
	from := ui.Motion{Opacity: 0}
	return e.Transition(ui.ElementTransition{
		Duration: d,
		Ease:     ui.EaseOut,
		Enter:    &from,
	})
}

// keepOverlay makes sure the window's overlay layer exists this frame. An
// element's Enter applies only in a parent the previous frame had, so a
// panel whose first frame opens the layer would otherwise arrive instantly.
func keepOverlay(c *ui.Context) {
	ui.Overlay(c, func() {})
}
