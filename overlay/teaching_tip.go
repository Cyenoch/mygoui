package overlay

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TeachingTipOptions configures an anchored coach mark.
type TeachingTipOptions struct {
	Title, Body string
	// Footer builds the action row under the body.
	Footer func()
	// Placement anchors the tip; SideAuto places it below the anchor.
	Placement mygoui.Placement
	// Width defaults to 320 DIPs.
	Width float32
	// HideClose removes the close button, leaving dismissal to Escape and
	// the footer's controls.
	HideClose bool
	// OnDismiss runs after the close button or Escape closed the tip.
	OnDismiss func()
}

type teachingTipKey struct{ open *bool }

// TeachingTip points a note with an arrow at its anchor while *open is true.
// Without a built anchor it floats at the bottom center of the window
// instead, without an arrow. It is non-modal and takes no focus: the window
// stays usable, an outside press leaves it open, and its controls follow the
// anchor in the Tab order. Its close button, or Escape while the focus is
// inside it, sets *open to false.
func TeachingTip(c *ui.Context, anchor ui.Element, open *bool, o TeachingTipOptions) ui.Element {
	if open == nil {
		return ui.Element{}
	}
	anchored.Warm(c)
	if !*open {
		return ui.Element{}
	}
	dismiss := func() {
		*open = false
		if o.OnDismiss != nil {
			o.OnDismiss()
		}
	}
	t := theme.Colors(c)
	w, _ := c.Size()
	width := o.Width
	if width <= 0 {
		width = 320
	}
	width = min(width, max(float32(0), w-2*theme.Space8))
	var panel ui.Element
	ui.Overlay(c, func() {
		panel = ui.Column(c.Key(teachingTipKey{open})).Width(width).Gap(theme.Space12).
			Padding(theme.Space16).Role(ui.RoleDialog)
		switch {
		case o.Title == "":
			panel.Label(o.Body)
		case o.Body == "":
			panel.Label(o.Title)
		default:
			panel.Label(o.Title).Description(o.Body)
		}
		if anchor.Valid() {
			teachingTipAttach(c, panel, anchor, o.Placement)
		} else {
			panel.Absolute().Left((w - width) / 2).Bottom(theme.Space16)
			panel = anchored.Motion(c, panel, mygoui.SideTop, false)
		}
		panel = anchored.Surface(c, panel, theme.RadiusScale().LG)
		if panel.Shortcut(0, ui.KeyEscape) {
			dismiss()
		}
		panel.Children(func() {
			teachingTipHeader(c, o, dismiss)
			if o.Body != "" {
				style.Text(ui.Text(c, o.Body), c, theme.SM()).TextColor(t.MutedForeground)
			}
			if o.Footer != nil {
				ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Center).Justify(ui.End).
					Gap(theme.Space8).Children(o.Footer)
			}
		})
	})
	return panel
}

// teachingTipAttach places the tip beside its anchor and draws its arrow on
// the side the panel settled on.
func teachingTipAttach(c *ui.Context, panel, anchor ui.Element, placement mygoui.Placement) {
	if placement.Offset == 0 {
		placement.Offset = arrowOffset
	}
	side := arrowSide(panel, anchor, resolveSide(placement, mygoui.SideBottom))
	anchored.Motion(c, panel, side, false)
	along := arrowAlong(panel, anchor, side)
	t := theme.Colors(c)
	panel.Draw(func(p *ui.Painter, r ui.Rect) {
		paintOutlinedArrow(p, r, side, t.Popover, t.Border, along)
	})
	anchored.Apply(panel, anchor, placement, mygoui.SideBottom)
}

func teachingTipHeader(c *ui.Context, o TeachingTipOptions, dismiss func()) {
	if o.Title == "" && o.HideClose {
		return
	}
	ui.Row(c).FillWidth().MinWidth(0).AlignItems(ui.Start).Justify(ui.End).Gap(theme.Space8).Children(func() {
		if o.Title != "" {
			style.Text(ui.Text(c, o.Title), c, theme.SM()).FontWeight(theme.WeightSemibold).
				Grow(1).Basis(0).MinWidth(0).Role(ui.RoleHeading).Level(3)
		}
		if !o.HideClose {
			// Pull the button into the padding so its glyph lines up with
			// the text edge.
			mygoui.CloseButton(c, mygoui.CloseButtonOptions{Label: locale.M().Close, OnClick: dismiss}).
				Margin(-theme.Space4, 0, 0, 0).MarginEnd(-theme.Space4)
		}
	})
}
