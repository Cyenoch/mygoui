package layout

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TitleBarSpan says which edges of the window a title bar reaches, and so
// which window-control insets it reserves. Start and end follow the layout
// direction.
type TitleBarSpan uint8

const (
	// SpanWindow reaches both edges: the bar of a window without a
	// full-height pane.
	SpanWindow TitleBarSpan = iota
	// SpanStart reaches the start edge only: the band of a full-height
	// start pane.
	SpanStart
	// SpanEnd reaches the end edge only: the bar beside a docked start pane.
	SpanEnd
)

// TitleBarOptions configures a window's title bar. Leading and Trailing
// build controls that keep their width; Center fills the room between them
// and aligns its own content. Title is shown there when Center is nil, and
// truncates. The bar is the window's banner landmark, named by Label, else
// Title, else the locale's TitleBar.
//
// MinHeight is the native band the window was created with
// (mygo.WindowOptions.TitleBarHeight), in window DIPs: the bar covers at
// least it and the band the window reports. Borderless drops the hairline
// under the bar.
type TitleBarOptions struct {
	Leading    func()
	Center     func()
	Trailing   func()
	Title      string
	Label      string
	Span       TitleBarSpan
	MinHeight  float32
	Borderless bool
}

type (
	barKey         struct{}
	barLeadingKey  struct{}
	barCenterKey   struct{}
	barTrailingKey struct{}
)

// barMeasure is how a bar was laid out last: its own width and the widths of
// its leading and trailing groups.
type barMeasure struct{ width, leading, trailing float32 }

// barGeometry is a bar's height for one build. A stacked bar holds its
// content below the native band, top down from the window edge.
type barGeometry struct {
	height, top float32
	stacked     bool
}

// TitleBar builds a draggable title band. It reserves the insets the
// window reports for its controls on the edges Span reaches, is at least as
// tall as the native band and as its controls at the current text scale,
// and moves its content below the band when Leading and Trailing no longer
// fit between the window controls. The widths are the layout of the frame
// before, so a change settles in the next frame.
//
// The bar's own surface moves the window; controls inside it take their
// presses and are ordinary Tab stops.
func TitleBar(c *ui.Context, o TitleBarOptions) ui.Element {
	root := ui.Row(c)
	m := ui.Local(root, barKey{}, func() barMeasure { return barMeasure{} })
	left, right := barInsets(c, o.Span)
	return fillBar(c, root, o, resolveBar(c, o, m), m, left, right)
}

// barInsets returns the window-control insets on the physical left and
// right edges the span reaches.
func barInsets(c *ui.Context, span TitleBarSpan) (left, right float32) {
	bar := c.TitleBar()
	left, right = bar.Left, bar.Right
	startIsLeft := c.Direction() != ui.RightToLeft
	switch span {
	case SpanStart:
		if startIsLeft {
			right = 0
		} else {
			left = 0
		}
	case SpanEnd:
		if startIsLeft {
			left = 0
		} else {
			right = 0
		}
	}
	return left, right
}

func barPadding(c *ui.Context) float32 { return theme.Scale(c, theme.Space(3)) }

// resolveBar decides a bar's height from how it was laid out last.
func resolveBar(c *ui.Context, o TitleBarOptions, m *barMeasure) barGeometry {
	// The native band is in window DIPs; the view's are zoomed.
	band := max(o.MinHeight/c.Zoom(), c.TitleBar().Height)
	control := theme.Scale(c, theme.HeightButtonSM) + 2*theme.Scale(c, theme.Space(2))
	g := barGeometry{height: max(band, control)}
	left, right := barInsets(c, o.Span)
	if left+right == 0 {
		// Below the band the content would have the same room.
		return g
	}
	width := m.width
	if width == 0 {
		width, _ = c.Size()
	}
	need := m.leading + m.trailing
	if m.leading > 0 {
		need += theme.Space(2)
	}
	if m.trailing > 0 {
		need += theme.Space(2)
	}
	if need > width-left-right-2*barPadding(c) {
		g.stacked, g.top = true, band
		g.height = band + control - theme.Scale(c, theme.Space(2))
	}
	return g
}

// laidOut stores the width e had in the last frame's layout into w, and
// asks for another build when it differs from what the bar was resolved with.
func laidOut(c *ui.Context, e ui.Element, w *float32) float32 {
	width := e.LayoutBounds().W
	if d := *w - width; d < -0.5 || d > 0.5 {
		*w = width
		c.Invalidate()
	}
	return width
}

// rebuildOnResize asks for another build when e is laid out at another
// width than built, the width this build read.
func rebuildOnResize(c *ui.Context, e ui.Element, built float32) {
	e.DrawOver(func(_ *ui.Painter, r ui.Rect) {
		if d := built - r.W; d < -0.5 || d > 0.5 {
			c.Invalidate()
		}
	})
}

// fillBar builds a title bar on root with geometry g, clear of left and
// right at the window's edges. A nil m measures nothing: the bar then follows
// a geometry decided elsewhere. Only the bar's own parts watch their layout,
// so the returned element is the caller's to decorate.
func fillBar(c *ui.Context, root ui.Element, o TitleBarOptions, g barGeometry, m *barMeasure, left, right float32) ui.Element {
	t := theme.Colors(c)
	pad := barPadding(c)
	if g.stacked {
		left, right = 0, 0
	}
	label := o.Label
	if label == "" {
		label = o.Title
	}
	if label == "" {
		label = locale.M().TitleBar
	}
	root.FillWidth().MinWidth(0).Height(g.height).Shrink(0).DragWindow().
		Padding(g.top, right+pad, 0, left+pad).AlignItems(ui.Center).Gap(theme.Space(2)).
		Role(ui.RoleBanner).Label(label)
	if !o.Borderless {
		root.BorderWidth(0, 0, theme.Hairline, 0).BorderColor(t.Border)
	}
	if m != nil {
		laidOut(c, root, &m.width)
	}
	group := func(key any, build func(), w *float32) {
		if build == nil {
			if w != nil {
				*w = 0
			}
			return
		}
		e := ui.Row(c.Key(key)).Shrink(0).AlignItems(ui.Center).Gap(theme.Space(2))
		if w != nil {
			rebuildOnResize(c, e, laidOut(c, e, w))
		}
		e.Children(build)
	}
	var leading, trailing *float32
	if m != nil {
		leading, trailing = &m.leading, &m.trailing
	}
	root.Children(func() {
		group(barLeadingKey{}, o.Leading, leading)
		center := ui.Row(c.Key(barCenterKey{})).Grow(1).Basis(0).MinWidth(0).AlignItems(ui.Center)
		if m != nil {
			// The center takes what is left, so it resizes with the bar.
			rebuildOnResize(c, center, center.LayoutBounds().W)
		}
		center.Children(func() {
			switch {
			case o.Center != nil:
				o.Center()
			case o.Title != "":
				ts := theme.SM()
				ts.Weight = theme.WeightMedium
				style.Text(ui.Text(c, o.Title), c, ts).SingleLine().MinWidth(0)
			}
		})
		group(barTrailingKey{}, o.Trailing, trailing)
	})
	return root
}
