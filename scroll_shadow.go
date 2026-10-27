package mygoui

import (
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// ScrollShadowOrientation selects the scroll axis; vertical is the default.
type ScrollShadowOrientation uint8

const (
	ScrollShadowVertical ScrollShadowOrientation = iota
	ScrollShadowHorizontal
)

// ScrollShadowOptions describes an opaque scrolling surface. The edge fade
// is painted in Surface, which should be the actual opaque background.
// ScrollArea with MaskFade instead fades content to transparency.
type ScrollShadowOptions struct {
	Orientation ScrollShadowOrientation
	Surface     ui.Color
	Size        float32
	Height      float32
	Label       string
}

// ScrollShadow wraps caller-owned content in a native scroll viewport and
// fades both edges in proportion to the remaining scroll distance.
func ScrollShadow(c *ui.Context, s *ui.ScrollState, o ScrollShadowOptions, body func()) ui.Element {
	if s == nil {
		return ui.Element{}
	}
	surface := o.Surface
	if surface.A == 0 {
		surface = theme.Colors(c).Background
	}
	depth := o.Size
	if depth <= 0 {
		depth = theme.Space40
	}
	var e ui.Element
	if o.Orientation == ScrollShadowHorizontal {
		e = ui.ScrollHorizontal(c).TrackScroll(s)
	} else {
		e = ui.Scroll(c).TrackScroll(s)
	}
	if o.Height > 0 {
		e.Height(o.Height)
	}
	e.MinHeight(0).Label(o.Label).Background(surface)
	e.Children(body)
	return ScrollFade(c, e, s, ScrollFadeOptions{Mode: SurfaceFade, Surface: surface, Size: depth})
}
