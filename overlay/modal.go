package overlay

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ModalSize is the modal width scale.
type ModalSize uint8

const (
	// ModalMedium is 448 DIPs wide before viewport clamping.
	ModalMedium ModalSize = iota
	// ModalSmall is max-w-sm (384 DIPs).
	ModalSmall
	// ModalLarge is max-w-lg (512 DIPs).
	ModalLarge
	// ModalExtraLarge extends the scale to max-w-xl (576 DIPs).
	ModalExtraLarge
	// ModalFull fills the window without radius or shadow.
	ModalFull
	// ModalExtraSmall is 320 DIPs wide before viewport clamping.
	ModalExtraSmall
	// ModalCover fills the space its container leaves, keeping the floating
	// surface's radius and shadow.
	ModalCover
)

// ModalOptions configures a modal.
type ModalOptions struct {
	Size        ModalSize
	Backdrop    ModalBackdrop
	Title       string
	Description string
	// Header replaces the generated title/description block.
	Header func()
	// Body is the scrollable content; Footer is the action row.
	Body   func()
	Footer func()
	// BodyLabel names the scroll region for assistive technology.
	BodyLabel string
	// HideClose removes the top-right close button.
	HideClose bool
	// PreventDismiss keeps the modal open on backdrop press and Escape.
	PreventDismiss bool
	// OnDismiss runs when the backdrop, Escape or the close button closed
	// the modal, after Open is set false.
	OnDismiss func()
}

// Modal shows a modal over a dimmed window while *open is true. The
// native dialog base keeps the focus inside, restores it to the opener, and
// closes on backdrop press or Escape unless PreventDismiss is set.
func Modal(c *ui.Context, open *bool, o ModalOptions) ui.Element {
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
	if *open {
		c.Key(dialogKey{open})
	}
	p := ui.DialogBase(c, open, func(back, panel ui.Element) {
		w, h := c.Size()
		full := o.Size == ModalFull
		pad := float32(theme.Space40)
		if w < 640 {
			pad = theme.Space16
		}
		if full {
			pad = 0
		}
		back = paintBackdrop(c, back, o.Backdrop)
		back = dialogMotion(c, back, dialogMotionDuration)
		back.Padding(pad).Row().Center()

		t := theme.Colors(c)
		radius := overlayRadius()
		if full {
			radius = 0
		}
		panel.Radius(radius).Clip().Background(t.Popover).TextColor(t.PopoverForeground).Border(theme.Hairline, t.Border)
		if !full {
			panel = panelElevation(c, panel)
		}
		availW := max(float32(0), w-2*pad)
		availH := max(float32(0), h-2*pad)
		switch {
		case full:
			panel.Width(w).Height(h)
		case o.Size == ModalCover:
			panel.Width(availW).MaxWidth(availW).Height(availH).
				MaxHeight(availH).MinHeight(0)
		default:
			panel.Width(min(modalWidth(o.Size), availW)).MaxWidth(availW).
				MaxHeight(availH).MinHeight(0)
		}
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

// modalWidth is the maximum width for a floating size. ModalFull and ModalCover
// take the space their container leaves instead.
func modalWidth(size ModalSize) float32 {
	switch size {
	case ModalExtraSmall:
		return 320
	case ModalSmall:
		return 384
	case ModalLarge:
		return 512
	case ModalExtraLarge:
		return 576
	case ModalFull, ModalCover:
		return 0
	}
	return 448
}
