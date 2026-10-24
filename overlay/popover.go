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
	// back to below the anchor; a zero Offset uses HeroUI's 8 DIPs.
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
	want := resolveSide(o.Placement, mygoui.SideBottom)
	label := o.Label
	if label == "" {
		label = o.Title
	}
	minWidth := o.MinWidth
	if minWidth == 0 && o.Width <= 0 {
		minWidth = 180
	}
	p := ui.PopoverBase(c, anchor, open, func(panel ui.Element) {
		// Key the panel before reading the side it settled on: a keyed
		// element keeps its state, and so its committed box, from frame to
		// frame. The attach engine flips a panel that no longer fits below
		// its anchor, and both the arrow and the zoom follow that side.
		panel = panel.Key(popoverKey{open})
		side := arrowSide(panel, anchor, want)
		panel = popupMotion(c, panel, side)
		panel.Radius(radius).Background(t.Overlay).TextColor(t.OverlayForeground).
			Padding(theme.Space16)
		if o.Width > 0 {
			panel.Width(o.Width)
		}
		if minWidth > 0 {
			panel.MinWidth(minWidth)
		}
		if o.Arrow {
			along := arrowAlong(panel, anchor, side)
			panel.Draw(func(p *ui.Painter, r ui.Rect) {
				paintArrow(p, r, side, t.Overlay, along)
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
		// A zero Offset keeps HeroUI's React Aria default of 8 DIPs, which
		// the anchor placement applies, with or without an arrow.
		anchored.Apply(panel, anchor, o.Placement, mygoui.SideBottom)
	})
	return p
}
