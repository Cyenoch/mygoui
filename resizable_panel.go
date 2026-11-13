package mygoui

import (
	"time"

	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// PanelSide selects the edge of a ResizablePanel that carries its handle.
// Start and end follow the layout direction.
type PanelSide uint8

const (
	PanelEnd PanelSide = iota
	PanelStart
	PanelTop
	PanelBottom
)

// ResizablePanelOptions configures a panel resized from one edge. Size, Min
// and Max are DIPs along the resized axis before the text scale; Size is
// always shown within [Min, Max], and a nonpositive Max sets no upper bound.
// A nil Size is kept by the panel. A non-nil Collapsed makes the panel
// collapsible. Label names the handle and defaults to the locale's Resize.
type ResizablePanelOptions struct {
	Label            string
	Side             PanelSide
	Size             *float32
	Min, Max         float32
	Collapsed        *bool
	OnResize         func(size float32)
	OnCollapseChange func(collapsed bool)
}

const (
	resizablePanelHit  = 8
	resizablePanelStep = 16
)

type resizablePanelKey struct{}

// resizablePanelState is the transient drag and the size shown last frame.
type resizablePanelState struct {
	size     float32 // used when the caller passes no Size
	dragging bool
	raw      float32 // unclamped size under the pointer
	start    float32 // size when the drag began
	shown    float32
}

// resizablePanel resolves the options for one build.
type resizablePanel struct {
	o    ResizablePanelOptions
	st   *resizablePanelState
	size *float32
}

func (p resizablePanel) clamp(v float32) float32 {
	v = max(v, p.o.Min)
	if p.o.Max > 0 {
		v = min(v, max(p.o.Max, p.o.Min))
	}
	return v
}

func (p resizablePanel) collapsed() bool { return p.o.Collapsed != nil && *p.o.Collapsed }

func (p resizablePanel) resize(v float32) {
	v = p.clamp(v)
	if v == *p.size {
		return
	}
	*p.size = v
	if p.o.OnResize != nil {
		p.o.OnResize(v)
	}
}

func (p resizablePanel) collapse(collapsed bool) {
	if p.o.Collapsed == nil || *p.o.Collapsed == collapsed {
		return
	}
	*p.o.Collapsed = collapsed
	if p.o.OnCollapseChange != nil {
		p.o.OnCollapseChange(collapsed)
	}
}

// step moves the edge from the keyboard; growing a collapsed panel expands it.
func (p resizablePanel) step(delta float32) {
	if p.collapsed() {
		if delta > 0 {
			p.collapse(false)
		}
		return
	}
	p.resize(p.clamp(*p.size) + delta)
}

// drag follows the pointer by delta. Below half of Min a collapsible panel
// collapses and keeps the size it had before the drag.
func (p resizablePanel) drag(delta float32, held bool) {
	st := p.st
	if !held {
		st.dragging = false
		return
	}
	if !st.dragging {
		// Start from the size shown, which may be clamped above the stored one.
		shown := p.clamp(*p.size)
		st.dragging, st.start, st.raw = true, shown, shown
		if p.collapsed() {
			st.raw = 0
		}
	}
	if delta == 0 {
		return
	}
	st.raw += delta
	if p.o.Collapsed != nil && st.raw < p.o.Min/2 {
		p.collapse(true)
		p.resize(st.start)
		return
	}
	p.collapse(false)
	p.resize(st.raw)
}

// ResizablePanel sizes body along one axis and lets the user move the edge
// on Side by pointer or keyboard: arrows step 16 DIPs, Home and End go to Min
// and Max, and on a collapsible panel Enter or a double click toggles it. The
// cross axis stretches to the parent.
func ResizablePanel(c *ui.Context, o ResizablePanelOptions, body func()) ui.Element {
	root := ui.Box(c).Shrink(0).AlignSelf(ui.Stretch).ZIndex(1)
	st := ui.Local(root, resizablePanelKey{}, func() resizablePanelState { return resizablePanelState{} })
	p := resizablePanel{o: o, st: st, size: o.Size}
	if p.size == nil {
		p.size = &st.size
	}
	vertical := o.Side == PanelTop || o.Side == PanelBottom
	// far is whether the handle sits on the right or bottom edge, where
	// moving it along the axis grows the panel.
	far := o.Side == PanelBottom
	if !vertical {
		far = (o.Side == PanelEnd) != (c.Direction() == ui.RightToLeft)
	}
	root.Children(func() {
		content := ui.Column(c).Fill().MinWidth(0).MinHeight(0).Clip()
		if !p.collapsed() || st.shown > 0 {
			content.Children(body)
		}
		p.handle(c, vertical, far)
	})
	target := theme.Scale(c, p.clamp(*p.size))
	if p.collapsed() {
		target = 0
	}
	if st.dragging {
		st.shown = root.Jump(resizablePanelKey{}, target)
	} else {
		st.shown = root.AnimateTo(resizablePanelKey{}, target, ui.Animation{
			Duration: 200 * time.Millisecond, Ease: ui.EaseOut,
		})
	}
	if vertical {
		return root.Height(st.shown)
	}
	return root.Width(st.shown)
}

// handle builds the splitter straddling the panel's edge and applies its
// pointer and keyboard input.
func (p resizablePanel) handle(c *ui.Context, vertical, far bool) {
	t := theme.Colors(c)
	label := p.o.Label
	if label == "" {
		label = locale.M().Resize
	}
	h := ui.Box(c).Absolute().Focusable().FocusRing(false).
		Role(ui.RoleSplitter).Label(label).
		Range(float64(p.o.Min), float64(max(p.o.Max, *p.size)), float64(*p.size)).
		Step(resizablePanelStep)
	const out = -resizablePanelHit / 2
	grow, shrink := ui.KeyRight, ui.KeyLeft
	switch {
	case vertical:
		h.Left(0).Right(0).Height(resizablePanelHit).Cursor(ui.CursorResizeNS).
			Orientation(ui.Horizontal)
		grow, shrink = ui.KeyDown, ui.KeyUp
		if far {
			h.Bottom(out)
		} else {
			h.Top(out)
		}
	default:
		h.Top(0).Bottom(0).Width(resizablePanelHit).Cursor(ui.CursorResizeEW).
			Orientation(ui.Vertical)
		if far {
			h.Right(out)
		} else {
			h.Left(out)
		}
	}
	if !far {
		grow, shrink = shrink, grow
	}

	dx, dy, held := h.Dragged()
	delta := dx
	if vertical {
		delta = dy
	}
	if !far {
		delta = -delta
	}
	p.drag(delta/theme.TextScale(c), held)
	if h.DoubleClicked() {
		p.collapse(!p.collapsed())
		// The press still held starts over from the toggled state.
		p.st.dragging = false
	}
	h.OnShortcut(0, grow, func() { p.step(resizablePanelStep) })
	h.OnShortcut(0, shrink, func() { p.step(-resizablePanelStep) })
	h.OnShortcut(0, ui.KeyHome, func() { p.step(p.o.Min - *p.size) })
	if p.o.Max > 0 {
		h.OnShortcut(0, ui.KeyEnd, func() { p.step(p.o.Max - *p.size) })
	}
	if p.o.Collapsed != nil {
		h.Expanded(!p.collapsed())
		h.OnShortcut(0, ui.KeyEnter, func() { p.collapse(!p.collapsed()) })
	}

	// A hairline on the edge; it thickens in the ring color while in use.
	active := held || h.Hovered() || h.FocusVisible()
	color, thick, at := t.Border, float32(theme.Hairline), float32(0)
	switch {
	case active:
		color, thick, at = t.Ring, 2, -1
	case far:
		at = -theme.Hairline
	}
	color = h.AnimateColor(resizablePanelKey{}, color, ui.Animation{
		Duration: 100 * time.Millisecond, Ease: ui.EaseOut,
	})
	// At rest the line lies inside the panel; in use it centers on the edge.
	h.Draw(func(pt *ui.Painter, r ui.Rect) {
		if vertical {
			pt.Fill(ui.Rect{X: r.X, Y: r.Y + r.H/2 + at, W: r.W, H: thick}, color, 0)
			return
		}
		pt.Fill(ui.Rect{X: r.X + r.W/2 + at, Y: r.Y, W: thick, H: r.H}, color, 0)
	})
}
