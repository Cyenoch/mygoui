package overlay

import "github.com/egoist/mygo/ui"

// focusReturn hands the focus on when a non-modal surface closes with the
// focus inside it. Such a surface takes no focus as it opens, so the
// toolkit may hold no opener to return it to, and a surface built in its
// container has none at all: the surface names its own target.
type focusReturn struct{ inside bool }

// note records whether the focus is inside surface, built open this frame.
func (r *focusReturn) note(surface ui.Element) { r.inside = surface.FocusWithin() }

// closed calls give once for the surface closed after holding the focus.
func (r *focusReturn) closed(give func()) {
	if r.inside {
		r.inside = false
		give()
	}
}

type focusReturnKey struct{ open *bool }

// overlayFocusReturn keeps the focusReturn of an overlay surface open while
// *open is true on a zero-size owner built in an overlay of its own every
// frame, which also keeps the overlays after it in place (anchored.Warm).
func overlayFocusReturn(c *ui.Context, open *bool) *focusReturn {
	var r *focusReturn
	ui.Overlay(c, func() {
		owner := ui.Box(c.Key(focusReturnKey{open})).Absolute().Size(0, 0).PassThrough().Role(ui.RoleNone)
		r = ui.Local(owner, focusReturnKey{}, func() focusReturn { return focusReturn{} })
	})
	return r
}
