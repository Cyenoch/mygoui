package overlay

import (
	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/internal/anchored"
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TooltipOptions configures a tooltip.
type TooltipOptions struct {
	Text    string
	Content func()
	// Placement anchors the tip; SideAuto falls back above the anchor.
	Placement mygoui.Placement
	// Arrow draws a 12 DIP overlay arrow.
	Arrow bool
}

// Tooltip shows a tooltip for anchor on pointer rest or keyboard
// focus, and dismisses it with Escape or on pointer press. The native tooltip
// base owns the show timing (0.6 s rest, warm reuse within 0.4 s). Text also
// supplies the anchor's accessible description.
func Tooltip(c *ui.Context, anchor ui.Element, o TooltipOptions) ui.Element {
	warmOverlay(c)
	if !anchor.Valid() {
		return ui.Element{}
	}
	if o.Text != "" {
		anchor.Description(o.Text)
	}
	t := theme.Colors(c)
	radius := theme.RadiusScale().MD
	w, h := c.Size()
	want := resolveSide(o.Placement, mygoui.SideTop)
	tip := ui.TooltipBase(c, anchor, func(tip ui.Element) {
		// Key the tip before reading the side it settled on: a keyed element
		// keeps its committed box from frame to frame. The native tooltip
		// base flips a tip that does not fit above its anchor, and both the
		// arrow and the zoom follow that side.
		tip = tip.Key(tooltipKey{anchor.ID()})
		side := arrowSide(tip, anchor, want)
		tip = popupMotion(c, tip, side)
		tip.MaxWidth(min(float32(320), max(float32(0), w-2*theme.Space8))).
			MaxHeight(max(float32(0), h-2*theme.Space8)).
			Padding(theme.Space4, theme.Space8).Radius(radius).
			Background(t.Popover).TextColor(t.PopoverForeground).Border(theme.Hairline, t.Border)
		tip = panelElevation(c, tip)
		if o.Arrow {
			along := arrowAlong(tip, anchor, side)
			tip.Draw(func(p *ui.Painter, r ui.Rect) {
				paintArrow(p, r, side, t.Popover, along)
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
				placement.Offset = 4
			}
		}
		anchored.Apply(tip, anchor, placement, mygoui.SideTop)
	})
	return tip
}
