package overlay

import (
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// HeroUI's surface motion, read from the component stylesheets: a backdrop
// fades as it comes and goes and a panel fades in place.
//
//	A transition carries one duration for both directions, so each surface
//	takes the enter duration of its reference: the modal backdrop's 150 ms
//	(it leaves over 100 ms), the drawer backdrop's 250 ms (it leaves over
//	200 ms), and the popover's and tooltip's 150 ms (they leave over 100 ms).
const (
	dialogMotionDuration = 150 * time.Millisecond
	drawerMotionDuration = 250 * time.Millisecond
	popupMotionDuration  = 150 * time.Millisecond
)

// HeroUI's popover and tooltip zooms: zoom-in-90 arrives from 90% and
// zoom-out-95 leaves at 95% of the panel's size, around the anchor point the
// panel points at.
const (
	popupEnterScale = 0.90
	popupExitScale  = 0.95
)

// Allocate modal keys before DialogBase: its backdrop reads click/shortcut
// state before invoking the styling callback. Rekeying there loses input and
// focus restoration. Anchored panels key before building their content.
type (
	dialogKey  struct{ open *bool }
	popoverKey struct{ open *bool }
	tooltipKey struct{ anchor uint64 }
)

// warmOverlay builds the overlay layer now. An element's Enter applies only
// in a parent the frame before had, and the layer is rebuilt each frame: a
// component must ask for it every frame, or the first surface a window shows
// after a frame without overlays would appear without its arrival. It adds
// no element: an empty overlay only makes sure the layer has been built.
func warmOverlay(c *ui.Context) {
	ui.Overlay(c, func() {})
}

// surfaceMotion fades an overlay surface in and out.
//
// An element's Exit paints an inert copy of the frame it last showed, under
// its parent, so no surface rebuilds its content to leave. MyGo lays an
// absolute overlay out again after its transitions, which discards a
// positional motion; and the whole modal backdrop, scrim included, is one
// element, so the surfaces fade as one layer rather than move. Under reduced
// motion they settle at once.
func surfaceMotion(c *ui.Context, e ui.Element, d time.Duration) ui.Element {
	if theme.ReduceMotion(c) {
		return e
	}
	motion := ui.Motion{Opacity: 0}
	return e.Transition(ui.ElementTransition{
		Duration: d,
		Ease:     ui.EaseOut,
		Enter:    &motion,
		Exit:     &motion,
	})
}

// dialogMotion fades the native dialog backdrop, and so the panel it holds.
// DialogBase builds neither once *open is false, so the backdrop is the one
// element a Transition can keep as the dialog goes: its exit copy holds the
// last paint of the caller's body inside the panel.
//
// The panel itself is a child of that backdrop, which is born in the frame the
// dialog opens: a Transition's Enter needs a parent the frame before had, so a
// panel scale would have to run on frames the library schedules itself. This
// package keeps the fade instead: a per-frame rebuild of the view for a 4%
// zoom is the cost a native window must not pay, and scaling the one backdrop
// element would shrink the scrim and the hit area outside the panel.
func dialogMotion(c *ui.Context, back ui.Element, d time.Duration) ui.Element {
	return surfaceMotion(c, back, d)
}

// popupMotion fades an anchored panel and zooms it from HeroUI's scale around
// the anchor point. The caller keys the panel first, before building its
// content, so that the content follows the panel's identity and the side it
// settled on can be read from its committed box. The panel is a child of the
// overlay layer, which warmOverlay keeps alive, so its Enter applies: the
// transition drives its own frames and nothing is left at a partial scale.
func popupMotion(c *ui.Context, panel ui.Element, side mygoui.Side) ui.Element {
	if theme.ReduceMotion(c) {
		return panel
	}
	panel.TransformOrigin(popupOrigin(side))
	enter := ui.Motion{Opacity: 0, Scale: popupEnterScale}
	exit := ui.Motion{Opacity: 0, Scale: popupExitScale}
	return panel.Transition(ui.ElementTransition{
		Duration: popupMotionDuration,
		Ease:     ui.EaseOut,
		Enter:    &enter,
		Exit:     &exit,
	})
}

// popupOrigin returns the transform origin HeroUI's origin-(--trigger-anchor-
// point) uses for a panel on side: the point of the panel that faces the
// anchor grows and shrinks in place.
func popupOrigin(side mygoui.Side) (x, y float32) {
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
