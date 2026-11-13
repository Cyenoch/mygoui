package navigation

import (
	"time"

	"github.com/Cyenoch/mygoui/internal/motion"
	"github.com/egoist/mygo/ui"
)

// Navigation motion uses MyGo's animation channels so interrupted changes
// continue from their current position.
const (
	// tabsIndicatorMotion is the tab indicator's 200 ms slide.
	tabsIndicatorMotion = 200 * time.Millisecond
)

// slideMotion returns an indicator animation over d, ease-in-out as in coss.
func slideMotion(d time.Duration) ui.Animation {
	return ui.Animation{Duration: d, Ease: ui.EaseInOut}
}

// selectionInk is the ink of an item that is chosen or not. The change of
// choice fades; the hover ink of an unchosen item shows at once.
func selectionInk(e ui.Element, key any, selected, hovered bool, rest, on, hover ui.Color) ui.Color {
	target := rest
	if selected {
		target = on
	}
	ink := motion.Color(e, key, target)
	if hovered && !selected {
		return hover
	}
	return ink
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
