package overlay

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type alertCancelFocusKey struct{}

type alertCancelGeometry struct {
	panel ui.Rect
	set   bool
}

// AlertDialogOptions configures an explicit-choice alert dialog. Backdrop
// presses and Escape do not dismiss it.
type AlertDialogOptions struct {
	Status mygoui.StatusColor
	// ShowIcon displays the status icon. An explicit Icon also enables it.
	ShowIcon bool
	Icon     *ui.SVG
	// HideIcon overrides ShowIcon and Icon.
	HideIcon bool
	// Size selects the desktop width. Its zero value (ModalMedium) uses
	// 512 DIPs here; set Width to 448 for the standard ModalMedium width.
	Size ModalSize
	// Width overrides Size on desktop. Narrow alerts fill the window.
	Width       float32
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
	// ShowClose adds a top-right cancel control, absent by default.
	ShowClose bool
	// HideClose overrides ShowClose.
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
	motionState := responsiveDialogState(c, open)
	modalScrim(c, open, o.Backdrop, dialogMotionDuration, ui.EaseInOut)
	if *open {
		c.Key(dialogKey{open})
	}
	p := ui.DialogBase(c, open, func(back, panel ui.Element) {
		w, h := c.Size()
		availW, availH := dialogViewport(back, panel, w, h)

		t := theme.Colors(c)
		width := modalWidth(o.Size)
		if o.Size == ModalMedium {
			width = 512
		}
		if o.Size == ModalCover || o.Size == ModalFull {
			width = availW
		}
		if o.Width > 0 {
			width = o.Width
		}
		width = min(width, availW)
		if w < 640 {
			width = availW
		}
		radius := overlayRadius()
		panel.Role(ui.RoleAlertDialog).
			Width(width).MaxWidth(availW).
			MaxHeight(availH).MinHeight(0).
			Radius(radius).Clip().
			Background(t.Popover).TextColor(t.PopoverForeground).Border(theme.Hairline, t.Border)
		if w < 640 {
			radius = 0
			panel.Radius(0).BorderWidth(theme.Hairline, 0, 0, 0)
		}
		panel = panelElevation(c, panel)
		panel.Label(o.Title)
		if o.Description != "" {
			panel.Description(o.Description)
		}
		var icon func()
		if !o.HideIcon && (o.ShowIcon || o.Icon != nil) {
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
			hideClose:        !o.ShowClose || o.HideClose,
			onClose:          cancel,
			centerHeader:     w < 640,
			titleStyle:       theme.H4(),
			descriptionStyle: theme.SM(),
			padding:          theme.Space24,
			headerGap:        theme.Space8,
			maxHeight:        availH,
			footerCorners:    [2]float32{max(0, radius-theme.Hairline), max(0, radius-theme.Hairline)},
			footer: func() {
				cancelButton = mygoui.Button(c, mygoui.ButtonOptions{
					Label:     cancelLabel,
					Variant:   mygoui.ButtonGhost,
					FullWidth: w < 640,
					OnClick:   cancel,
				})
				if showAction {
					variant := mygoui.ButtonPrimary
					if o.Destructive {
						variant = mygoui.ButtonDanger
					}
					mygoui.Button(c, mygoui.ButtonOptions{
						Label:     o.ActionLabel,
						Variant:   variant,
						Disabled:  o.ActionDisabled,
						FullWidth: w < 640,
						OnClick:   accept,
					})
				}
			},
		})
		responsiveDialogMotion(c, back, panel, motionState, w, h)
	})
	// DialogBase records the opener after its callback. Focus Cancel only
	// after that capture, or the dialog would try to restore its own button.
	cancelButton.AutoFocus()
	if cancelButton.Valid() {
		last := ui.Local(cancelButton, alertCancelFocusKey{}, func() alertCancelGeometry { return alertCancelGeometry{} })
		bounds := p.LayoutBounds()
		if !last.set || last.panel != bounds && cancelButton.Focused() {
			cancelButton.ScrollIntoView()
			last.panel, last.set = bounds, true
		}
	}
	if wasOpen && !*open && p.Valid() && p.Dismissed() {
		// Alert dialogs require an explicit choice: restore the caller's
		// state instead of accepting a backdrop press or Escape.
		*open = true
	}
	motionState.observe(*open, theme.ReduceMotion(c), c.Now())
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
