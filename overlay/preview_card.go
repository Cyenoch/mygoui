package overlay

import (
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// PreviewCardState is the caller-owned hover and focus state of a preview.
// Keep one state per anchor and window.
type PreviewCardState struct {
	Open      bool
	panel     ui.Handle
	leave     time.Time
	dismissed bool
}

// PreviewCardOptions describes rich content shown beside an anchor.
type PreviewCardOptions struct {
	Placement mygoui.Placement
	Width     float32
	Label     string
	Body      func()
}

// PreviewCard opens for pointer hover or focus and keeps its interactive
// content open while the pointer or focus moves into it. Escape and an
// outside press dismiss it through MyGo's native popover.
func PreviewCard(c *ui.Context, anchor ui.Element, state *PreviewCardState, o PreviewCardOptions) ui.Element {
	if state == nil {
		panic("overlay: PreviewCard requires a caller-owned state")
	}
	if !anchor.Valid() {
		state.Open = false
		state.leave = time.Time{}
		state.dismissed = false
		return ui.Element{}
	}
	_, _, overAnchor := anchor.PointerPosition()
	anchorActive := overAnchor || anchor.FocusWithin()
	if anchorActive && !state.dismissed {
		state.Open = true
	}
	width := o.Width
	if width <= 0 {
		width = 256
	}
	panel := Popover(c, anchor, &state.Open, PopoverOptions{
		Placement: o.Placement, Width: width, Label: o.Label, Body: o.Body,
	})
	overPanel := false
	if panel.Valid() {
		panel.Bind(&state.panel).FontSize(theme.Scale(c, theme.SM().Size))
		_, _, overPanel = panel.PointerPosition()
		if panel.Dismissed() {
			state.dismissed = true
		}
	}
	active := anchorActive || overPanel || state.panel.FocusWithin(c)
	if !active {
		state.dismissed = false
	}
	if active && !state.dismissed {
		state.Open = true
		state.leave = time.Time{}
	} else if state.Open {
		if state.leave.IsZero() {
			state.leave = c.Now().Add(150 * time.Millisecond)
		}
		if !c.Now().Before(state.leave) {
			state.Open = false
			c.AnimationFrame()
		} else {
			c.After(state.leave.Sub(c.Now()))
		}
	}
	return panel
}
