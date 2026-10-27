package overlay

import (
	"github.com/Cyenoch/mygoui"
	"github.com/egoist/mygo/ui"
)

// SheetOptions configures a modal edge sheet. The default side is right.
type SheetOptions struct {
	Side     mygoui.Side
	Width    float32
	Backdrop ModalBackdrop
	// Inset floats the sheet 16 DIPs from the window on wider screens.
	Inset                bool
	Title, Description   string
	Header, Body, Footer func()
	BodyLabel            string
	HideClose            bool
	PreventDismiss       bool
	OnDismiss            func()
}

// Sheet builds an edge-aligned native modal. It slides 32 DIPs over 200 ms
// and restores focus immediately when dismissed. A zero Width uses 448 DIPs
// for side sheets, capped to leave 48 DIPs of the window visible.
func Sheet(c *ui.Context, open *bool, o SheetOptions) ui.Element {
	side := o.Side
	if side == mygoui.SideAuto {
		side = mygoui.SideRight
	}
	return dialog(c, open, DialogOptions{
		Placement: mygoui.Placement{Side: side}, Width: o.Width, Inset: o.Inset, Backdrop: o.Backdrop,
		Title: o.Title, Description: o.Description, Header: o.Header,
		Body: o.Body, Footer: o.Footer, BodyLabel: o.BodyLabel,
		HideClose: o.HideClose, PreventDismiss: o.PreventDismiss, OnDismiss: o.OnDismiss,
	}, true)
}
