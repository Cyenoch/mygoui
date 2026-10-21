// Package layout provides scoped direction, resizable panes and tracked native
// scrolling. Widget state belongs to the caller and is independent per preview.
package layout

import (
	"math"

	"github.com/Cyenoch/mygoui"
	"github.com/egoist/mygo/ui"
)

// Direction establishes native direction before building its descendants.
// Native controls, overlays and logical edges inherit this scope.
func Direction(c *ui.Context, direction ui.Direction, body func()) ui.Element {
	return ui.Column(c).Direction(direction).MinWidth(0).Children(body)
}

type ResizableState struct {
	Percent     float32
	initialized bool
}

type ResizableOptions struct {
	Vertical bool
	Label    string
	Disabled bool
	OnChange func(float32)
}

// Resizable preserves the caller's proportion when its available space changes.
// Native Split owns pointer capture, RTL, focusable keyboard divider and minimum
// pane boundaries, including temporary compression in a small viewport.
func Resizable(c *ui.Context, state *ResizableState, o ResizableOptions, first, second func()) ui.Element {
	if !state.initialized {
		if state.Percent == 0 {
			state.Percent = 50
		}
		state.initialized = true
	}
	if math.IsNaN(float64(state.Percent)) || math.IsInf(float64(state.Percent), 0) {
		state.Percent = 50
	}
	state.Percent = min(100, max(0, state.Percent))
	// Parent Disabled must be applied before the native divider is constructed.
	return ui.Box(c).MinWidth(0).MinHeight(0).Disabled(o.Disabled).Label(o.Label).Children(func() {
		var split ui.Element
		if o.Vertical {
			split = ui.SplitVerticalPercent(c, &state.Percent, first, second)
		} else {
			split = ui.SplitPercent(c, &state.Percent, first, second)
		}
		split.Fill().OnChange(func() {
			if o.OnChange != nil {
				o.OnChange(state.Percent)
			}
		})
	})
}

type ScrollOrientation uint8

const (
	ScrollVertical ScrollOrientation = iota
	ScrollHorizontal
	ScrollBoth
)

type ScrollAreaOptions struct {
	Orientation ScrollOrientation
	Label       string
	Edges       *mygoui.ScrollEdgesOptions
}

// ScrollArea uses native wheel, touchpad, scrollbar and keyboard scrolling with
// caller-owned offsets. Bound its returned element with Height, Size or Grow.
func ScrollArea(c *ui.Context, state *ui.ScrollState, o ScrollAreaOptions, body func()) ui.Element {
	var area ui.Element
	switch o.Orientation {
	case ScrollHorizontal:
		area = ui.ScrollHorizontal(c)
	case ScrollBoth:
		area = ui.ScrollBoth(c)
	default:
		area = ui.Scroll(c)
	}
	area.MinWidth(0).MinHeight(0).Label(o.Label)
	if state != nil {
		area.TrackScroll(state)
	}
	if o.Edges != nil {
		edges := *o.Edges
		switch o.Orientation {
		case ScrollHorizontal:
			edges.Axes = ui.Horizontal
		case ScrollBoth:
			edges.Axes = ui.Horizontal | ui.Vertical
		default:
			edges.Axes = ui.Vertical
		}
		mygoui.ScrollEdges(c, area, state, edges)
	}
	return area.Children(body)
}
