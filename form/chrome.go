package form

import (
	"github.com/Cyenoch/mygoui/internal/formstyle"
	"github.com/egoist/mygo/ui"
)

// ControlVariant selects the form-field surface. Both variants use the
// same Coss field face.
type ControlVariant uint8

const (
	ControlPrimary ControlVariant = iota
	ControlSecondary
)

// fieldSurface paints the field background, border and invalid/focus
// states, including the field ring, around a native editor.
func fieldSurface(e ui.Element, c *ui.Context, variant ControlVariant, invalid, disabled bool) ui.Element {
	return formstyle.Surface(e, c, variant == ControlSecondary, invalid, disabled)
}

// triggerSurface is fieldSurface for a button-like popup trigger, whose
// upstream CSS rings on :focus-visible rather than every focus.
func triggerSurface(e ui.Element, c *ui.Context, variant ControlVariant, invalid, disabled bool) ui.Element {
	return formstyle.SurfaceMode(e, c, formstyle.RingOnKeyboardFocus,
		variant == ControlSecondary, invalid, disabled)
}

// fieldText applies field typography and foreground.
func fieldText(e ui.Element, c *ui.Context) ui.Element {
	return formstyle.Text(e, c)
}

// fieldLabel associates a visible label with an already-built control handle
// without rebuilding the control.
func fieldLabel(c *ui.Context, text string, required, invalid, disabled bool, target *ui.Handle) ui.Element {
	return Label(c, text, LabelOptions{
		Required: required, Invalid: invalid, Disabled: disabled, Target: target,
	})
}
