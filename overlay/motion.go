package overlay

import (
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// HeroUI's surface motion: an overlay fades as it comes and goes.
const (
	dialogMotionDuration = 200 * time.Millisecond
	popupMotionDuration  = 150 * time.Millisecond
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
// after a frame without overlays would appear without its arrival.
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
func dialogMotion(c *ui.Context, back ui.Element) ui.Element {
	return surfaceMotion(c, back, dialogMotionDuration)
}

// popupMotion fades an anchored panel. It is called before the panel builds
// its content, so that the content follows the panel's identity.
func popupMotion(c *ui.Context, panel ui.Element, key any) ui.Element {
	return surfaceMotion(c, panel.Key(key), popupMotionDuration)
}
