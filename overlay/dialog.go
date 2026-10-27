package overlay

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// DialogOptions configures a responsive or edge-aligned native modal panel.
type DialogOptions struct {
	// Placement.Side selects the window edge. SideAuto places the dialog
	// above center on desktop and against the bottom on narrow windows.
	// Align is ignored. Edge dialogs span the edge like a sheet.
	Placement mygoui.Placement
	Backdrop  ModalBackdrop
	Width     float32
	// Inset floats an edge-aligned dialog 16 DIPs from the window on widths
	// of 640 DIPs or more.
	Inset       bool
	Title       string
	Description string
	// Footer builds actions directly; use ButtonOptions.FullWidth on narrow windows.
	Header, Body, Footer func()
	BodyLabel            string
	HideClose            bool
	PreventDismiss       bool
	OnDismiss            func()
}

// Dialog shows a modal panel with heading, body, footer and close control.
// It keeps and restores focus through the native modal dialog base.
func Dialog(c *ui.Context, open *bool, o DialogOptions) ui.Element {
	return dialog(c, open, o, false)
}

func dialog(c *ui.Context, open *bool, o DialogOptions, slide bool) ui.Element {
	if open == nil {
		return ui.Element{}
	}
	warmOverlay(c)
	wasOpen := *open
	close := func() {
		*open = false
		if o.OnDismiss != nil {
			o.OnDismiss()
		}
	}
	w, h := c.Size()
	width := o.Width
	if width <= 0 {
		width = 512
	}
	side := resolveSide(o.Placement, mygoui.SideAuto)
	if side == mygoui.SideLeft || side == mygoui.SideRight {
		if o.Width <= 0 {
			width = 448
		}
		width = min(width, max(float32(0), w-theme.Space48))
	}
	inset := float32(0)
	if o.Inset && w >= 640 && side != mygoui.SideAuto {
		inset = theme.Space16
	}
	var motionState *dialogState
	if !slide && side == mygoui.SideAuto {
		motionState = responsiveDialogState(c, open)
	}
	if slide || side == mygoui.SideAuto {
		modalScrim(c, open, o.Backdrop, dialogMotionDuration, ui.EaseInOut)
	}
	if *open {
		c.Key(dialogKey{open})
	}
	p := ui.DialogBase(c, open, func(back, panel ui.Element) {
		if slide {
			back = edgeExit(c, back, side, theme.Space32, dialogMotionDuration, ui.EaseInOut)
		} else if side != mygoui.SideAuto {
			back = paintBackdrop(c, back, o.Backdrop)
			back = dialogMotion(c, back, dialogMotionDuration)
		}
		t := theme.Colors(c)
		radius := overlayRadius()
		footerCorners := [2]float32{}
		contentHeight := max(float32(0), h-theme.Space48)
		panel.Background(t.Popover).TextColor(t.PopoverForeground).Border(theme.Hairline, t.Border).MinHeight(0)
		switch side {
		case mygoui.SideLeft:
			contentHeight = max(0, h-2*inset)
			back.Row().Justify(ui.Start).AlignItems(ui.Stretch).Padding(inset)
			panel.Clip().Width(min(width, max(float32(0), w-2*inset))).Height(max(float32(0), h-2*inset))
		case mygoui.SideRight:
			contentHeight = max(0, h-2*inset)
			back.Row().Justify(ui.End).AlignItems(ui.Stretch).Padding(inset)
			panel.Clip().Width(min(width, max(float32(0), w-2*inset))).Height(max(float32(0), h-2*inset))
		case mygoui.SideTop:
			back.Column().Justify(ui.Start).AlignItems(ui.Stretch).Padding(inset)
			panel.Clip().FillWidth().MaxHeight(max(float32(0), h-theme.Space48)).Radius(0)
		case mygoui.SideBottom:
			back.Column().Justify(ui.End).AlignItems(ui.Stretch).Padding(inset)
			panel.Clip().FillWidth().MaxHeight(max(float32(0), h-theme.Space48)).Radius(0)
		default:
			availW, availH := dialogViewport(back, panel, w, h)
			contentHeight = availH
			panel.Clip().Radius(radius).
				Width(min(width, availW)).MaxWidth(availW).
				MaxHeight(availH).MinHeight(0)
			if w < 640 {
				panel.Width(availW).Radius(0).BorderWidth(theme.Hairline, 0, 0, 0)
			}
		}
		if side == mygoui.SideLeft || side == mygoui.SideRight {
			panel.Radius(0)
		}
		if inset > 0 {
			panel.Radius(radius)
		}
		if inset > 0 || side == mygoui.SideAuto && w >= 640 {
			footerCorners = [2]float32{radius - theme.Hairline, radius - theme.Hairline}
		}
		if slide {
			edgeEnter(c, panel, side, theme.Space32, dialogMotionDuration, ui.EaseInOut)
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
			titleStyle:       theme.H4(),
			descriptionStyle: theme.SM(),
			padding:          theme.Space24,
			headerGap:        theme.Space8,
			maxHeight:        contentHeight,
			fillHeight:       side == mygoui.SideLeft || side == mygoui.SideRight,
			footerCorners:    footerCorners,
		})
		if !slide && side == mygoui.SideAuto {
			responsiveDialogMotion(c, back, panel, motionState, w, h)
		}
	})
	if wasOpen && !*open && p.Valid() && p.Dismissed() {
		if o.PreventDismiss {
			*open = true
		} else if o.OnDismiss != nil {
			o.OnDismiss()
		}
	}
	if motionState != nil {
		motionState.observe(*open, theme.ReduceMotion(c), c.Now())
	}
	return p
}
