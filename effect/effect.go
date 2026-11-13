// Package effect applies shader effects to an application's own elements:
// Mask fades edges to transparency, Shimmer passes a loading highlight over
// placeholders and Reveal wipes arriving content in.
//
// An effect changes only what the element paints. Layout, hit testing,
// focus and accessibility stay with the element, so an effect never is the
// only signal of a state.
//
// MyGo renders the element and its children into a layer and draws the
// layer through the effect. An element has one effect: applying another
// replaces it, whether the application or a component such as a mask
// ScrollFade installed it. Embedded web views do not show through a layer.
//
// Each effect draws a stated fallback instead of running its shader where
// the frame has no GPU (the shader would be interpreted per pixel), where a
// side of the element exceeds MyGo's 8192-pixel layer limit, device scale
// included, and where its options set Fallback.
//
// A shader is compiled the first time it paints, so an import, a build and
// a fallback compile nothing.
package effect

import (
	"github.com/Cyenoch/mygoui/internal/fx"
	"github.com/egoist/mygo/ui"
)

// gpu reports whether the frame runs shaders. Tests force it, since a
// headless frame never draws on a GPU.
func gpu(p *ui.Painter, forced bool) bool { return forced || p.Accelerated() }

// shaded reports whether MyGo draws box through an effect's shader in this
// frame: on a GPU, within the layer limit.
func shaded(p *ui.Painter, box ui.Rect, forced bool) bool {
	return gpu(p, forced) && fx.Fits(box, p.Scale())
}
