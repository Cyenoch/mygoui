package overlay

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// PopoverOptions configures an anchored HeroUI popover.
type PopoverOptions struct {
	Title, Description string
	// Placement anchors the panel with mygoui.Placement. SideAuto falls
	// back to below the anchor; zero Offset uses 4 DIPs, or 7 with Arrow.
	Placement mygoui.Placement
	// Arrow draws HeroUI's 12 DIP arrow at the anchor-facing edge.
	Arrow bool
	// Width fixes the panel width; zero sizes it to its content.
	Width float32
	// MinWidth keeps small popovers usable; zero uses 180 DIPs.
	MinWidth             float32
	Header, Body, Footer func()
	// Label names the popup; Title is the fallback.
	Label string
}

// Popover shows an anchored HeroUI popover while *open is true. It is
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
	radius := overlayRadius()
	side := resolveSide(o.Placement, mygoui.SideBottom)
	label := o.Label
	if label == "" {
		label = o.Title
	}
	p := ui.PopoverBase(c, anchor, open, func(panel ui.Element) {
		panel = popupMotion(c, panel, popoverKey{open})
		panel.Radius(radius).Background(t.Overlay).TextColor(t.OverlayForeground).
			Padding(theme.Space16).MinWidth(o.MinWidth)
		if o.Width > 0 {
			panel.Width(o.Width)
		}
		if o.MinWidth > 0 {
			panel.MinWidth(o.MinWidth)
		} else if o.Width <= 0 {
			panel.MinWidth(180)
		}
		if o.Arrow {
			panel.Draw(func(p *ui.Painter, r ui.Rect) {
				paintArrow(p, r, side, t.Overlay)
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
		ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space12).Children(func() {
			if hasHeader {
				h := ui.Column(c).FillWidth().MinWidth(0).Gap(theme.Space4)
				h.Children(func() {
					if o.Header != nil {
						o.Header()
						return
					}
					if o.Title != "" {
						title := style.Text(ui.Text(c, o.Title), c, theme.SM())
						title.FontWeight(theme.WeightMedium).Role(ui.RoleHeading).Level(3)
					}
					if o.Description != "" {
						style.Text(ui.Text(c, o.Description), c, theme.SM()).TextColor(t.Muted)
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
		placement := o.Placement
		if placement.Offset == 0 {
			if o.Arrow {
				placement.Offset = 7
			} else {
				placement.Offset = 4
			}
		}
		anchored.Apply(panel, anchor, placement, mygoui.SideBottom)
	})
	return p
}
