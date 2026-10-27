package overlay

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// PopoverOptions configures an anchored popover.
type PopoverOptions struct {
	Title, Description string
	// Placement anchors the panel with mygoui.Placement. SideAuto falls
	// back to below the anchor; a zero Offset uses 4 DIPs.
	Placement mygoui.Placement
	// Arrow draws a 12 DIP arrow at the anchor-facing edge.
	Arrow bool
	// Width fixes the panel width; zero sizes it to its content. The panel
	// is capped to the window and its content scrolls when it is too tall.
	Width float32
	// MinWidth keeps small popovers usable; zero uses 180 DIPs.
	MinWidth             float32
	Header, Body, Footer func()
	// Label names the popup; Title is the fallback.
	Label string
}

// Popover shows an anchored popover while *open is true. It is
// non-modal: the press that dismisses it reaches what is under the pointer,
// and Escape or an outside press sets *open to false.
func Popover(c *ui.Context, anchor ui.Element, open *bool, o PopoverOptions) ui.Element {
	if open == nil {
		return ui.Element{}
	}
	warmOverlay(c)
	if !anchor.Valid() {
		return ui.Element{}
	}
	t := theme.Colors(c)
	radius := theme.RadiusScale().LG
	want := resolveSide(o.Placement, mygoui.SideBottom)
	label := o.Label
	if label == "" {
		label = o.Title
	}
	minWidth := o.MinWidth
	if minWidth == 0 && o.Width <= 0 {
		minWidth = 180
	}
	w, h := c.Size()
	availableWidth := max(float32(0), w-2*theme.Space8)
	availableHeight := max(float32(0), h-2*theme.Space8)
	p := ui.PopoverBase(c, anchor, open, func(panel ui.Element) {
		// Key the panel before reading the side it settled on: a keyed
		// element keeps its state, and so its committed box, from frame to
		// frame. The attach engine flips a panel that no longer fits below
		// its anchor, and both the arrow and the zoom follow that side.
		panel = panel.Key(popoverKey{open})
		side := arrowSide(panel, anchor, want)
		panel = popupMotion(c, panel, side)
		panel.Radius(radius).Background(t.Popover).TextColor(t.PopoverForeground).Border(theme.Hairline, t.Border).
			MaxWidth(availableWidth).MaxHeight(availableHeight).MinHeight(0)
		if o.Width > 0 {
			panel.Width(min(o.Width, availableWidth))
		}
		if minWidth > 0 {
			panel.MinWidth(min(minWidth, availableWidth))
		}
		if o.Arrow {
			along := arrowAlong(panel, anchor, side)
			panel.Draw(func(p *ui.Painter, r ui.Rect) {
				paintArrow(p, r, side, t.Popover, along)
			})
		}
		panel = panelElevation(c, panel)
		if label != "" {
			panel.Label(label)
		}
		if o.Description != "" {
			panel.Description(o.Description)
		}
		hasHeader := o.Header != nil || o.Title != "" || o.Description != ""
		hasBody := o.Body != nil
		hasFooter := o.Footer != nil
		ui.Scroll(c).FillWidth().MinWidth(0).MinHeight(0).Padding(theme.Space16).Children(func() {
			ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).Children(func() {
				if hasHeader {
					h := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space4)
					h.Children(func() {
						if o.Header != nil {
							o.Header()
							return
						}
						if o.Title != "" {
							title := style.Text(ui.Text(c, o.Title), c, theme.LG())
							title.FontWeight(theme.WeightSemibold).Role(ui.RoleHeading).Level(3)
						}
						if o.Description != "" {
							style.Text(ui.Text(c, o.Description), c, theme.SM()).TextColor(t.MutedForeground)
						}
					})
				}
				if hasBody {
					b := ui.Column(c).FillWidth().MinWidth(0)
					b.Children(o.Body)
				}
				if hasFooter {
					ui.Row(c).FillWidth().Justify(ui.End).Gap(theme.Space8).Children(o.Footer)
				}
			})
		})
		placement := o.Placement
		if placement.Offset == 0 {
			placement.Offset = 4
		}
		anchored.Apply(panel, anchor, placement, mygoui.SideBottom)
	})
	return p
}
