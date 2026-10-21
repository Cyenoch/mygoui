// Package overlay provides native modal, anchored, menu, and transient surfaces.
package overlay

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	myui "github.com/egoist/mygo/ui"
)

// DialogOptions keeps default title, description and actions together. Header
// replaces the default heading; Footer composes arbitrary native action controls.
type DialogOptions struct {
	Title, Description                      string
	Width                                   float32
	Small, HideClose, PreventOutsideDismiss bool
	Header, Footer                          func()
	// BodyScroll retains the native body viewport without replacing its Scroll.
	BodyScroll *myui.ScrollState
	BodyLabel  string
}

func Dialog(c *myui.Context, open *bool, o DialogOptions, body func()) myui.Element {
	return edgePanel(c, open, o, SideDefault, false, 0, body)
}

func styleDialog(c *myui.Context, panel myui.Element, o DialogOptions) {
	w, h := c.Size()
	p := theme.Colors(c)
	width := o.Width
	if width <= 0 {
		width = 448
		if o.Small {
			width = 320
		}
	}
	panel.Column().Width(width).MaxWidth(max(0, w-32)).MaxHeight(max(0, h-32)).
		Padding(theme.Space24).Gap(theme.Space24).Radius(theme.Radii().XL).Clip().
		Background(p.Popover).TextColor(p.PopoverForeground).
		Border(theme.HairlineWidth, p.Border).Label(o.Title).Description(o.Description)
}

func dialogContent(c *myui.Context, open *bool, o DialogOptions, body func()) {
	if o.Header != nil {
		o.Header()
	} else if o.HideClose {
		heading(c, o.Title, o.Description)
	} else {
		myui.Column(c).Padding(0, theme.Metrics(c).CompactHeight+theme.Space8, 0, 0).
			Children(func() { heading(c, o.Title, o.Description) })
	}
	if body != nil {
		viewport := myui.Scroll(c).MinHeight(0).Shrink(1)
		if o.BodyScroll != nil {
			viewport.TrackScroll(o.BodyScroll)
		}
		if o.BodyLabel != "" {
			viewport.Label(o.BodyLabel)
		}
		viewport.Children(body)
	}
	if o.Footer != nil {
		myui.Row(c).Wrap().Gap(8).Justify(myui.End).Children(o.Footer)
	}
	if !o.HideClose {
		mygoui.Button(c.Key("close"), mygoui.ButtonOptions{Label: locales.ComponentClose(), Icon: lucide.X(),
			Variant: mygoui.ButtonGhost, Size: mygoui.ButtonIconSM, OnClick: func() { *open = false }}).
			Absolute().Top(theme.Space16).Right(theme.Space16).Background(theme.Colors(c).Secondary)
	}
}

func heading(c *myui.Context, title, description string) {
	if title == "" && description == "" {
		return
	}
	myui.Column(c).Gap(6).Children(func() {
		if title != "" {
			myui.Text(c, title).Bold().FontSize(theme.Scale(c, 16)).Role(myui.RoleHeading).Level(2)
		}
		if description != "" {
			myui.Text(c, description).FontSize(theme.Scale(c, 14)).TextColor(theme.Colors(c).MutedForeground)
		}
	})
}

type AlertDialogOptions struct {
	Title, Description, CancelLabel, ActionLabel string
	Small, Destructive, ActionDisabled           bool
	OnAction, OnCancel                           func()
	Media                                        func()
}

// AlertDialog never dismisses on an outside press. Cancel receives initial
// focus; Enter activates only the focused control, never an implicit action.
func AlertDialog(c *myui.Context, open *bool, o AlertDialogOptions, body func()) myui.Element {
	cancel := func() {
		*open = false
		if o.OnCancel != nil {
			o.OnCancel()
		}
	}
	d := DialogOptions{Title: o.Title, Description: o.Description, Small: o.Small, HideClose: true, PreventOutsideDismiss: true}
	d.Header = func() {
		if o.Media != nil {
			myui.Column(c).Size(64, 64).Radius(32).Background(theme.Colors(c).Muted).Center().Children(o.Media)
		}
		heading(c, o.Title, o.Description)
	}
	d.Footer = func() {
		label := o.CancelLabel
		if label == "" {
			label = locales.ComponentCancel()
		}
		mygoui.Button(c.Key("cancel"), mygoui.ButtonOptions{Label: label, Variant: mygoui.ButtonOutline, OnClick: cancel}).AutoFocus()
		label = o.ActionLabel
		if label == "" {
			label = locales.ComponentContinue()
		}
		variant := mygoui.ButtonDefault
		if o.Destructive {
			variant = mygoui.ButtonDestructive
		}
		mygoui.Button(c.Key("action"), mygoui.ButtonOptions{Label: label, Variant: variant, Disabled: o.ActionDisabled || o.OnAction == nil,
			OnClick: func() {
				*open = false
				if o.OnAction != nil {
					o.OnAction()
				}
			}})
	}
	wasOpen := open != nil && *open
	panel := edgePanel(c, open, d, SideDefault, false, 0, body)
	if panel.Valid() {
		panel.Role(myui.RoleAlertDialog)
	}
	if wasOpen && !*open && o.OnCancel != nil {
		o.OnCancel()
	}
	return panel
}

// Side names physical window edges; it does not reverse under RTL.
type Side uint8

const (
	SideDefault Side = iota
	SideRight
	SideBottom
	SideLeft
	SideTop
)

type SheetOptions struct {
	DialogOptions
	Side Side
}
type DrawerOptions struct {
	DialogOptions
	Side     Side
	NonModal bool
	Height   float32
}

func Sheet(c *myui.Context, open *bool, o SheetOptions, body func()) myui.Element {
	side := o.Side
	if side == SideDefault {
		side = SideRight
	}
	return edgePanel(c, open, o.DialogOptions, side, false, 0, body)
}

// Drawer uses native modality, placement and scroll containment. Swipe-to-close,
// velocity snap points and nested swipe progress are not provided: pinned MyGo
// exposes pointer events but no maintained drawer gesture arbitration with Scroll.
func Drawer(c *myui.Context, open *bool, o DrawerOptions, body func()) myui.Element {
	side := o.Side
	if side == SideDefault {
		side = SideBottom
	}
	return edgePanel(c, open, o.DialogOptions, side, o.NonModal, o.Height, body)
}

func edgePanel(c *myui.Context, open *bool, o DialogOptions, side Side, nonModal bool, height float32, body func()) myui.Element {
	// Consume caller identity in its own scope, including while closed. Overlay's
	// shared root must survive closure so MyGo can paint its inert exit copy.
	identity := myui.Box(c).Absolute().Size(0, 0).PassThrough().ID()
	w, h := c.Size()
	var back, panel myui.Element
	myui.Overlay(c, func() {
		if open == nil || !*open {
			return
		}
		back = myui.Column(c.Key(identity)).Absolute().Left(0).Top(0).Right(0).Bottom(0)
		if nonModal {
			back.PassThrough()
		} else {
			back.Background(theme.Colors(c).Overlay).Modal().OnClick(func() {})
		}
		if !theme.ReduceMotion(c) {
			// One composite fade preserves the last painted layout without moving
			// the scrim or retaining a caller's body/data for an exit rebuild.
			back.Transition(myui.ElementTransition{Colors: true, Enter: &myui.Motion{}, Exit: &myui.Motion{}})
		}
	})
	if !back.Valid() {
		return panel
	}
	// Overlay captures its opener on return. Build AutoFocus content only after
	// that capture, otherwise it records the new field as its own opener.
	switch side {
	case SideRight:
		back.AlignItems(myui.End)
	case SideLeft:
		back.AlignItems(myui.Start)
	case SideBottom:
		back.Justify(myui.End)
	case SideDefault:
		back.Center()
	}
	back.Children(func() {
		panel = myui.Column(c).Role(myui.RoleDialog).OnClick(func() {})
		styleDialog(c, panel, o)
		switch side {
		case SideRight, SideLeft:
			width := o.Width
			if width <= 0 {
				width = 384
			}
			panel.Width(width).MaxWidth(w * 0.75).FillHeight().MaxHeight(h).Radius(0)
			if side == SideRight {
				panel.BorderWidth(0, 0, 0, theme.HairlineWidth)
			} else {
				panel.BorderWidth(0, theme.HairlineWidth, 0, 0)
			}
		case SideTop, SideBottom:
			panel.FillWidth().MaxWidth(w).MaxHeight(max(0, h-96))
			radius := theme.Radii().XL
			if side == SideBottom {
				panel.Radius(radius, radius, 0, 0).BorderWidth(theme.HairlineWidth, 0, 0, 0)
			} else {
				panel.Radius(0, 0, radius, radius).BorderWidth(0, 0, theme.HairlineWidth, 0)
			}
			if height > 0 {
				panel.Height(height)
			}
		}
		panel.Children(func() { dialogContent(c, open, o, body) })
		if panel.OverlayShortcut(0, myui.KeyEscape) || (!o.PreventOutsideDismiss && panel.PressedOutside()) {
			*open = false
		}
	})
	return panel
}
