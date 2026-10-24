package overlay

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// DialogOptions configures native HeroUI's dialog, the shared centered (or
// edge-aligned) modal panel with a heading, description and close control.
type DialogOptions struct {
	// Placement.Side selects the window edge; SideAuto centers the dialog
	// and Align is ignored. Edge dialogs span the edge like a sheet.
	Placement            mygoui.Placement
	Backdrop             ModalBackdrop
	Width                float32
	Title                string
	Description          string
	Header, Body, Footer func()
	BodyLabel            string
	HideClose            bool
	PreventDismiss       bool
	OnDismiss            func()
}

// Dialog shows native HeroUI's dialog: bg-overlay, 24 DIP radius, 20 DIP
// padding, an 18 DIP medium label and a 16 DIP muted description. It keeps
// and restores focus through the native modal dialog base.
func Dialog(c *ui.Context, open *bool, o DialogOptions) ui.Element {
	if open == nil {
		return ui.Element{}
	}
	wasOpen := *open
	close := func() {
		*open = false
		if o.OnDismiss != nil {
			o.OnDismiss()
		}
	}
	w, h := c.Size()
	pad := float32(theme.Space40)
	if w < 640 {
		pad = theme.Space16
	}
	width := o.Width
	if width <= 0 {
		width = 448
	}
	side := resolveSide(o.Placement, mygoui.SideAuto)
	p := ui.DialogBase(c, open, func(back, panel ui.Element) {
		back = paintBackdrop(c, back, o.Backdrop)
		t := theme.Colors(c)
		radius := overlayRadius()
		panel.Background(t.Overlay).TextColor(t.OverlayForeground).MinHeight(0)
		switch side {
		case mygoui.SideLeft:
			back.Row().Justify(ui.Start).AlignItems(ui.Stretch).Padding(0)
			panel.Clip().Width(min(width, w)).Height(h)
		case mygoui.SideRight:
			back.Row().Justify(ui.End).AlignItems(ui.Stretch).Padding(0)
			panel.Clip().Width(min(width, w)).Height(h)
		case mygoui.SideTop:
			back.Column().Justify(ui.Start).AlignItems(ui.Stretch).Padding(0)
			panel.Clip().FillWidth().MaxHeight(h*0.85).Radius(0, 0, radius, radius)
		case mygoui.SideBottom:
			back.Column().Justify(ui.End).AlignItems(ui.Stretch).Padding(0)
			panel.Clip().FillWidth().MaxHeight(h*0.85).Radius(radius, radius, 0, 0)
		default:
			back.Row().Center().Padding(pad)
			availW := max(float32(0), w-2*pad)
			availH := max(float32(0), h-2*pad)
			panel.Clip().Radius(radius).
				Width(min(width, availW)).MaxWidth(availW).
				MaxHeight(availH).MinHeight(0)
		}
		if side == mygoui.SideLeft || side == mygoui.SideRight {
			panel.Radius(0)
		}
		panel = panelElevation(c, panel)
		panel.Label(o.Title)
		if o.Description != "" {
			panel.Description(o.Description)
		}
		fillPanel(c, panel, chrome{
			title:            o.Title,
			description:      o.Description,
			header:           o.Header,
			body:             o.Body,
			footer:           o.Footer,
			bodyLabel:        o.BodyLabel,
			hideClose:        o.HideClose,
			onClose:          close,
			titleStyle:       theme.LG(),
			descriptionStyle: theme.Base(),
			padding:          theme.Space20,
			headerGap:        theme.Space12,
		})
	})
	if wasOpen && !*open && p.Valid() && p.Dismissed() {
		if o.PreventDismiss {
			*open = true
		} else if o.OnDismiss != nil {
			o.OnDismiss()
		}
	}
	return p
}
