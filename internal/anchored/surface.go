package anchored

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

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

// Still gives e a transition that shows every change at once, for reduced
// motion, where e would otherwise move: the toolkit takes an element that
// drops its transition while shown, as the preference turns on, for one that
// went, and leaves an inert copy of it over the live one that takes its
// focus. A negative Duration also keeps a theme change from fading its
// colours and shadows over the toolkit's default 200 ms.
func Still(e ui.Element) ui.Element { return e.Transition(stillTransition) }

var stillTransition = ui.ElementTransition{Duration: -1}

// Motion fades a popup in from the popup scale around the edge facing its anchor and
// fades it out, as coss popovers do; zoomOut also shrinks it on exit, as coss
// tooltips and preview cards do. The caller keys the panel before building its content and
// calls Warm every frame. Reduced motion shows and hides it at once.
func Motion(c *ui.Context, panel ui.Element, side mygoui.Side, zoomOut bool) ui.Element {
	if theme.ReduceMotion(c) {
		return Still(panel)
	}
	panel.TransformOrigin(origin(side))
	m, scale := theme.Motion(), theme.Effects().PopupScale
	enter := ui.Motion{Opacity: 0, Scale: scale}
	exit := ui.Motion{Opacity: 0}
	if zoomOut {
		exit.Scale = scale
	}
	return panel.Transition(ui.ElementTransition{
		Duration: m.Popup, Ease: m.Ease, Enter: &enter, Exit: &exit,
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
