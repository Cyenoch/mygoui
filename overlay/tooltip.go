package overlay

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TooltipOptions configures a HeroUI tooltip.
type TooltipOptions struct {
	Text    string
	Content func()
	// Placement anchors the tip; SideAuto falls back above the anchor.
	Placement mygoui.Placement
	// Arrow draws HeroUI's 12 DIP overlay arrow.
	Arrow bool
}

// Tooltip shows a HeroUI tooltip for anchor on pointer rest or keyboard
// focus, and dismisses it with Escape or on pointer press. The native tooltip
// base owns the show timing (0.6 s rest, warm reuse within 0.4 s); HeroUI's
// 1.5 s delay and 0.5 s close delay are not exposed by MyGo.
func Tooltip(c *ui.Context, anchor ui.Element, o TooltipOptions) ui.Element {
	warmOverlay(c)
	if !anchor.Valid() {
		return ui.Element{}
	}
	t := theme.Colors(c)
	radius := min(theme.Space32, theme.RadiusScale().XL)
	side := resolveSide(o.Placement, mygoui.SideTop)
	tip := ui.TooltipBase(c, anchor, func(tip ui.Element) {
		tip = popupMotion(c, tip, tooltipKey{anchor.ID()})
		tip.MaxWidth(320).Padding(theme.Space8).Radius(radius).
			Background(t.Overlay).TextColor(t.OverlayForeground)
		tip = panelElevation(c, tip)
		if o.Arrow {
			tip.Draw(func(p *ui.Painter, r ui.Rect) {
				paintArrow(p, r, side, t.Overlay)
			})
		}
		if o.Text != "" {
			tip.Label(o.Text)
		}
		tip.Children(func() {
			if o.Content != nil {
				o.Content()
				return
			}
			style.Text(ui.Text(c, o.Text), c, theme.XS()).Wrap()
		})
		placement := o.Placement
		if placement.Offset == 0 {
			if o.Arrow {
				placement.Offset = 7
			} else {
				placement.Offset = 3
			}
		}
		anchored.Apply(tip, anchor, placement, mygoui.SideTop)
	})
	return tip
}
