package overlay

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ModalSize is HeroUI's modal width scale.
type ModalSize uint8

const (
	// ModalMedium is max-w-md (448 DIPs), HeroUI's default.
	ModalMedium ModalSize = iota
	// ModalSmall is max-w-sm (384 DIPs).
	ModalSmall
	// ModalLarge is max-w-lg (512 DIPs).
	ModalLarge
	// ModalExtraLarge extends the scale to max-w-xl (576 DIPs).
	ModalExtraLarge
	// ModalFull fills the window without radius or shadow.
	ModalFull
)

// ModalOptions configures a HeroUI modal.
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

// Modal shows a HeroUI modal over a dimmed window while *open is true. The
// native dialog base keeps the focus inside, restores it to the opener, and
// closes on backdrop press or Escape unless PreventDismiss is set.
func Modal(c *ui.Context, open *bool, o ModalOptions) ui.Element {
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
		back.Padding(pad).Row().Center()

		t := theme.Colors(c)
		radius := overlayRadius()
		if full {
			radius = 0
		}
		panel.Radius(radius).Clip().Background(t.Overlay).TextColor(t.OverlayForeground)
		panel = panelElevation(c, panel)
		if full {
			panel.Width(w).Height(h)
		} else {
			width := modalWidth(o.Size)
			avail := max(float32(0), w-2*pad)
			panel.Width(min(width, avail)).MaxWidth(avail).
				MaxHeight(max(float32(0), h-2*pad)).MinHeight(0)
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
			titleStyle:       theme.Base(),
			descriptionStyle: theme.SM(),
			bodyInk:          t.Muted,
			padding:          theme.Space24,
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

func modalWidth(size ModalSize) float32 {
	switch size {
	case ModalSmall:
		return 384
	case ModalLarge:
		return 512
	case ModalExtraLarge:
		return 576
	}
	return 448
}
