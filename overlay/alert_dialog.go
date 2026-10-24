package overlay

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// AlertDialogOptions configures a HeroUI alert dialog. Both choices are
// explicit: the backdrop and Escape never accept or dismiss it, matching
// HeroUI's isDismissable and isKeyboardDismissDisabled defaults.
type AlertDialogOptions struct {
	Status      mygoui.StatusColor
	Icon        *ui.SVG
	HideIcon    bool
	Size        ModalSize
	Backdrop    ModalBackdrop
	Title       string
	Description string
	Body        func()
	BodyLabel   string
	// CancelLabel defaults to the library's Close message; an empty
	// ActionLabel with no OnAction leaves the confirm action out rather
	// than fabricating one.
	CancelLabel, ActionLabel string
	OnCancel, OnAction       func()
	ActionDisabled           bool
	Destructive              bool
	// HideClose removes the top-right explicit dismiss control.
	HideClose bool
}

// AlertDialog shows an alert dialog asking for an explicit choice. Backdrop
// presses and Escape are ignored; Cancel receives the initial focus.
func AlertDialog(c *ui.Context, open *bool, o AlertDialogOptions) ui.Element {
	if open == nil {
		return ui.Element{}
	}
	warmOverlay(c)
	wasOpen := *open
	cancel := func() {
		*open = false
		if o.OnCancel != nil {
			o.OnCancel()
		}
	}
	accept := func() {
		if o.ActionDisabled {
			return
		}
		*open = false
		if o.OnAction != nil {
			o.OnAction()
		}
	}
	cancelLabel := o.CancelLabel
	if cancelLabel == "" {
		cancelLabel = locale.M().Close
	}
	showAction := o.ActionLabel != "" && o.OnAction != nil
	var cancelButton ui.Element
	if *open {
		c.Key(dialogKey{open})
	}
	p := ui.DialogBase(c, open, func(back, panel ui.Element) {
		w, h := c.Size()
		pad := float32(theme.Space40)
		if w < 640 {
			pad = theme.Space16
		}
		back = paintBackdrop(c, back, o.Backdrop)
		back = dialogMotion(c, back, dialogMotionDuration)
		back.Padding(pad).Row().Center()

		t := theme.Colors(c)
		availW := max(float32(0), w-2*pad)
		availH := max(float32(0), h-2*pad)
		// Alert dialogs have no full-screen size of their own: a size that
		// fills its container, ModalCover or ModalFull, takes the space the
		// padded container leaves.
		width := min(modalWidth(o.Size), availW)
		if o.Size == ModalCover || o.Size == ModalFull {
			width = availW
		}
		radius := overlayRadius()
		panel.Role(ui.RoleAlertDialog).
			Width(width).MaxWidth(availW).
			MaxHeight(availH).MinHeight(0).
			Radius(radius).Clip().
			Background(t.Overlay).TextColor(t.OverlayForeground)
		panel = panelElevation(c, panel)
		panel.Label(o.Title)
		if o.Description != "" {
			panel.Description(o.Description)
		}
		var icon func()
		if !o.HideIcon {
			glyph := o.Icon
			if glyph == nil {
				glyph = statusGlyph(o.Status)
			}
			icon = statusIcon(c, o.Status, glyph)
		}
		fillPanel(c, panel, chrome{
			title:            o.Title,
			description:      o.Description,
			body:             o.Body,
			bodyLabel:        o.BodyLabel,
			icon:             icon,
			hideClose:        o.HideClose,
			onClose:          cancel,
			titleStyle:       theme.Base(),
			descriptionStyle: theme.SM(),
			bodyInk:          t.Muted,
			padding:          theme.Space24,
			headerGap:        theme.Space12,
			footer: func() {
				cancelButton = mygoui.Button(c, mygoui.ButtonOptions{
					Label:   cancelLabel,
					Variant: mygoui.ButtonTertiary,
					OnClick: cancel,
				})
				if showAction {
					variant := mygoui.ButtonPrimary
					if o.Destructive {
						variant = mygoui.ButtonDanger
					}
					mygoui.Button(c, mygoui.ButtonOptions{
						Label:    o.ActionLabel,
						Variant:  variant,
						Disabled: o.ActionDisabled,
						OnClick:  accept,
					})
				}
			},
		})
	})
	// DialogBase records the opener after its callback. Focus Cancel only
	// after that capture, or the dialog would try to restore its own button.
	cancelButton.AutoFocus()
	if wasOpen && !*open && p.Valid() && p.Dismissed() {
		// Alert dialogs require an explicit choice: restore the caller's
		// state instead of accepting a backdrop press or Escape.
		*open = true
	}
	return p
}

func statusGlyph(status mygoui.StatusColor) *ui.SVG {
	switch status {
	case mygoui.ColorSuccess:
		return lucide.CircleCheck()
	case mygoui.ColorWarning:
		return lucide.TriangleAlert()
	case mygoui.ColorDanger:
		return lucide.CircleX()
	}
	return lucide.Info()
}
