package overlay

import (
	"github.com/Cyenoch/mygoicons/lucide"
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// PopconfirmOptions configures an anchored confirmation bubble.
type PopconfirmOptions struct {
	Title, Description string
	// ConfirmLabel and CancelLabel default to the library's Confirm and
	// Cancel messages.
	ConfirmLabel, CancelLabel string
	// Destructive marks the bubble with a warning glyph, gives the confirm
	// button the danger face and puts the initial focus on Cancel, so Enter
	// does not confirm by accident; otherwise Confirm has it.
	Destructive bool
	// Loading shows the confirm button pending while the confirmed work runs,
	// and holds the bubble: Cancel is off, and Escape and an outside press do
	// not close it.
	Loading bool
	// OnConfirm runs after the confirm button closed the bubble, which it
	// may keep open by setting *open back to true and Loading. OnCancel runs
	// after the cancel button, Escape or an outside press closed it.
	OnConfirm, OnCancel func()
	// Placement anchors the bubble; SideAuto places it below the anchor.
	Placement mygoui.Placement
	// Width defaults to 280 DIPs.
	Width float32
}

// arrowOffset is the default gap of a panel with an arrow: the arrow and the
// usual 4 DIPs between its tip and the anchor.
const arrowOffset = arrowHeight + 4

// Popconfirm asks for confirmation of an inline action in a popover at its
// anchor while *open is true. Like a popover it is non-modal: Escape or an
// outside press cancels, and the press reaches what is under the pointer.
func Popconfirm(c *ui.Context, anchor ui.Element, open *bool, o PopconfirmOptions) ui.Element {
	if open == nil {
		return ui.Element{}
	}
	wasOpen := *open
	cancel := func() {
		*open = false
		if o.OnCancel != nil {
			o.OnCancel()
		}
	}
	confirm := func() {
		*open = false
		if o.OnConfirm != nil {
			o.OnConfirm()
		}
	}
	width := o.Width
	if width <= 0 {
		width = 280
	}
	placement := o.Placement
	if placement.Offset == 0 {
		placement.Offset = arrowOffset
	}
	var focusButton ui.Element
	p := Popover(c, anchor, open, PopoverOptions{
		Title: o.Title, Description: o.Description,
		Placement: placement, Arrow: true, Width: width,
		Header: func(c *ui.Context) { popconfirmHeader(c, o) },
		Footer: func(c *ui.Context) {
			cancelButton := popconfirmButton(c, o.CancelLabel, locale.M().Cancel, mygoui.ButtonOutline, o.Loading, false, cancel)
			variant := mygoui.ButtonPrimary
			if o.Destructive {
				variant = mygoui.ButtonDanger
			}
			confirmButton := popconfirmButton(c, o.ConfirmLabel, locale.M().Confirm, variant, false, o.Loading, confirm)
			focusButton = confirmButton
			if o.Destructive {
				focusButton = cancelButton
			}
		},
	})
	// The popover records the element to give the focus back to after its
	// callback; focusing the safe button before would make it that element.
	focusButton.AutoFocus()
	if wasOpen && !*open && p.Valid() && p.Dismissed() {
		if o.Loading {
			// The popover closed after it was built: the work is still
			// running, so the bubble stays.
			*open = true
		} else if o.OnCancel != nil {
			o.OnCancel()
		}
	}
	return p
}

func popconfirmHeader(c *ui.Context, o PopconfirmOptions) {
	t := theme.Colors(c)
	ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Start).Gap(theme.Space8).Children(func() {
		if o.Destructive {
			// The glyph sits on the title's first line: half a line of 20 DIPs
			// less the icon's own half.
			ui.Icon(c, lucide.TriangleAlert()).Size(theme.Scale(c, theme.IconMD), theme.Scale(c, theme.IconMD)).
				Shrink(0).TextColor(t.Destructive).Margin(theme.Space2, 0, 0, 0)
		}
		ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(theme.Space4).Children(func() {
			if o.Title != "" {
				style.Text(ui.Text(c, o.Title), c, theme.SM()).FontWeight(theme.WeightSemibold).
					Role(ui.RoleHeading).Level(3)
			}
			if o.Description != "" {
				style.Text(ui.Text(c, o.Description), c, theme.SM()).TextColor(t.MutedForeground)
			}
		})
	})
}

func popconfirmButton(c *ui.Context, label, fallback string, variant mygoui.ButtonVariant, disabled, loading bool, onClick func()) ui.Element {
	if label == "" {
		label = fallback
	}
	return mygoui.Button(c, mygoui.ButtonOptions{
		Label: label, Variant: variant, Size: theme.ButtonSmall,
		Disabled: disabled, Loading: loading, OnClick: onClick,
	})
}
