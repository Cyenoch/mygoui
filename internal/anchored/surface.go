package anchored

import (
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// MotionDuration is the coss popover and tooltip transition.
const MotionDuration = 150 * time.Millisecond

// Surface paints the coss popup face: popover fill and text, hairline
// border, radius and overlay elevation.
func Surface(c *ui.Context, panel ui.Element, radius float32) ui.Element {
	t := theme.Colors(c)
	panel.Radius(radius).Background(t.Popover).TextColor(t.PopoverForeground).Border(theme.Hairline, t.Border)
	return style.Shadow(panel, c, theme.OverlayShadows(c.Theme().Dark))
}

// Warm builds the overlay layer now. An element's Enter applies only in a
// parent the frame before had, and the layer is rebuilt each frame: a
// component with an animated popup asks for it every frame, open or not, or
// its first popup after a frame without overlays arrives without motion.
func Warm(c *ui.Context) {
	ui.Overlay(c, func() {})
}

// Motion fades a popup in from 98% around the edge facing its anchor and
// fades it out, as coss popovers do; zoomOut also shrinks it on exit, as coss
// tooltips and preview cards do. The caller keys the panel before building its content and
// calls Warm every frame. Reduced motion shows and hides it at once.
func Motion(c *ui.Context, panel ui.Element, side mygoui.Side, zoomOut bool) ui.Element {
	if theme.ReduceMotion(c) {
		return panel
	}
	panel.TransformOrigin(origin(side))
	enter := ui.Motion{Opacity: 0, Scale: 0.98}
	exit := ui.Motion{Opacity: 0}
	if zoomOut {
		exit.Scale = 0.98
	}
	return panel.Transition(ui.ElementTransition{
		Duration: MotionDuration, Ease: ui.EaseOut, Enter: &enter, Exit: &exit,
	})
}

// origin is the point of a panel on side that faces its anchor.
func origin(side mygoui.Side) (x, y float32) {
	switch side {
	case mygoui.SideTop:
		return 0.5, 1
	case mygoui.SideLeft:
		return 1, 0.5
	case mygoui.SideRight:
		return 0, 0.5
	}
	return 0.5, 0
}
