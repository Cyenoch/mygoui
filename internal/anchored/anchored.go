// Package anchored maps typed placements onto MyGo's native anchor engine.
// It owns no state and performs no measurement.
package anchored

import (
	"github.com/Cyenoch/mygoui"
	"github.com/egoist/mygo/ui"
)

// defaultOffset matches the coss popover and picker positioners.
const defaultOffset = 4

// Apply attaches panel to anchor for a typed placement, falling back to the
// component-appropriate side when the caller leaves it unset.
func Apply(panel, anchor ui.Element, p mygoui.Placement, fallback mygoui.Side) ui.Element {
	side := p.Side
	if side == mygoui.SideAuto {
		side = fallback
	}
	offset := p.Offset
	if offset == 0 {
		offset = defaultOffset
	}
	at, self := ui.AnchorBottom, ui.AnchorTop
	switch side {
	case mygoui.SideTop:
		at, self = ui.AnchorTop, ui.AnchorBottom
		switch p.Align {
		case mygoui.AlignStart:
			at, self = ui.AnchorTopStart, ui.AnchorBottomStart
		case mygoui.AlignEnd:
			at, self = ui.AnchorTopEnd, ui.AnchorBottomEnd
		}
		panel.Margin(0, 0, offset, 0)
	case mygoui.SideRight:
		at, self = ui.AnchorRight, ui.AnchorLeft
		switch p.Align {
		case mygoui.AlignStart:
			at, self = ui.AnchorTopRight, ui.AnchorTopLeft
		case mygoui.AlignEnd:
			at, self = ui.AnchorBottomRight, ui.AnchorBottomLeft
		}
		panel.Margin(0, 0, 0, offset)
	case mygoui.SideLeft:
		at, self = ui.AnchorLeft, ui.AnchorRight
		switch p.Align {
		case mygoui.AlignStart:
			at, self = ui.AnchorTopLeft, ui.AnchorTopRight
		case mygoui.AlignEnd:
			at, self = ui.AnchorBottomLeft, ui.AnchorBottomRight
		}
		panel.Margin(0, offset, 0, 0)
	default:
		switch p.Align {
		case mygoui.AlignStart:
			at, self = ui.AnchorBottomStart, ui.AnchorTopStart
		case mygoui.AlignEnd:
			at, self = ui.AnchorBottomEnd, ui.AnchorTopEnd
		}
		panel.Margin(offset, 0, 0, 0)
	}
	return panel.AttachTo(anchor, at, self)
}
