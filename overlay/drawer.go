package overlay

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// DrawerSide is the window edge a drawer slides from.
type DrawerSide uint8

const (
	// DrawerBottom is HeroUI's default drawer placement.
	DrawerBottom DrawerSide = iota
	DrawerTop
	DrawerLeft
	DrawerRight
)

// DrawerOptions configures an edge drawer.
type DrawerOptions struct {
	Side                 DrawerSide
	Backdrop             ModalBackdrop
	Title                string
	Description          string
	Header, Body, Footer func()
	BodyLabel            string
	// Width overrides the panel width of left/right drawers (default 384,
	// capped at 85% of the window). Height overrides the panel extent of
	// top/bottom drawers (default content height, capped at 85%).
	Width, Height float32
	// Handle shows the drag-handle bar on top/bottom drawers. It is
	// visual; pointer swipe dismissal is not implemented.
	Handle         bool
	HideClose      bool
	PreventDismiss bool
	OnDismiss      func()
}

// Drawer shows a modal panel from one window edge. The native dialog base
// traps and restores the focus and dismisses on backdrop press or Escape.
// Swipe-to-dismiss is not provided: MyGo has no maintained drawer gesture
// arbitration with its scroll containers.
func Drawer(c *ui.Context, open *bool, o DrawerOptions) ui.Element {
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
	side := o.Side
	p := ui.DialogBase(c, open, func(back, panel ui.Element) {
		back = paintBackdrop(c, back, o.Backdrop)
		layoutEdge(back, side)
		t := theme.Colors(c)
		radius := min(theme.Space32, theme.RadiusScale().XXL)
		panel.Column().Clip().Background(t.Overlay).TextColor(t.OverlayForeground).MinHeight(0)
		panel = panelElevation(c, panel)
		switch side {
		case DrawerLeft, DrawerRight:
			width := o.Width
			if width <= 0 {
				width = 384
			}
			panel.Width(min(width, w*0.85)).Height(h).MaxWidth(w)
		default:
			panel.FillWidth().MaxHeight(h * 0.85)
			if o.Height > 0 {
				panel.Height(min(o.Height, h*0.85))
			}
			if side == DrawerTop {
				panel.Radius(0, 0, radius, radius)
			} else {
				panel.Radius(radius, radius, 0, 0)
			}
		}
		panel.Label(o.Title)
		if o.Description != "" {
			panel.Description(o.Description)
		}
		if o.Handle && side == DrawerBottom {
			drawerHandle(c)
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
			titleStyle:       theme.Base(),
			descriptionStyle: theme.SM(),
			bodyInk:          t.Muted,
			padding:          theme.Space24,
			headerGap:        theme.Space12,
		})
		if o.Handle && side == DrawerTop {
			drawerHandle(c)
		}
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

// BottomSheetOptions configures native HeroUI's bottom sheet.
type BottomSheetOptions struct {
	Backdrop             ModalBackdrop
	Title                string
	Description          string
	Header, Body, Footer func()
	BodyLabel            string
	// Height fixes the panel height; zero sizes it to its content. Both
	// are capped at 85% of the window.
	Height float32
	// HideHandle removes the visual drag handle. Swipe dismissal is not
	// implemented; the backdrop and Escape dismiss the sheet.
	HideHandle     bool
	HideClose      bool
	PreventDismiss bool
	OnDismiss      func()
}

// BottomSheet shows native HeroUI's bottom panel: overlay surface, 32 DIP
// top corners, a centered separator handle and 20 DIP padding. It is modal,
// so the native dialog base keeps and restores the focus.
func BottomSheet(c *ui.Context, open *bool, o BottomSheetOptions) ui.Element {
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
	_, h := c.Size()
	p := ui.DialogBase(c, open, func(back, panel ui.Element) {
		back = paintBackdrop(c, back, o.Backdrop)
		back.Column().Justify(ui.End).AlignItems(ui.Stretch).Padding(0)
		t := theme.Colors(c)
		radius := min(theme.Space32, theme.RadiusScale().XXXXL)
		panel.Column().FillWidth().MaxHeight(h*0.85).MinHeight(0).Clip().
			Radius(radius, radius, 0, 0).
			Background(t.Overlay).TextColor(t.OverlayForeground)
		if o.Height > 0 {
			panel.Height(min(o.Height, h*0.85))
		}
		panel = panelElevation(c, panel)
		panel.Label(o.Title)
		if o.Description != "" {
			panel.Description(o.Description)
		}
		if !o.HideHandle {
			drawerHandle(c)
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

// layoutEdge overrides DialogBase's centered backdrop for an edge panel.
func layoutEdge(back ui.Element, side DrawerSide) {
	switch side {
	case DrawerLeft:
		back.Row().Justify(ui.Start).AlignItems(ui.Stretch).Padding(0)
	case DrawerRight:
		back.Row().Justify(ui.End).AlignItems(ui.Stretch).Padding(0)
	case DrawerTop:
		back.Column().Justify(ui.Start).AlignItems(ui.Stretch).Padding(0)
	default:
		back.Column().Justify(ui.End).AlignItems(ui.Stretch).Padding(0)
	}
}

// drawerHandle paints the 36×4 separator bar in the 8 DIP handle gutter.
func drawerHandle(c *ui.Context) {
	t := theme.Colors(c)
	gutter := ui.Box(c).FillWidth().Center().PaddingY(theme.Space4).Shrink(0)
	gutter.Children(func() {
		ui.Box(c).Width(theme.Space32 + theme.Space4).Height(theme.Space4).
			Radius(theme.RadiusScale().XS).Background(t.Separator).Shrink(0)
	})
}
