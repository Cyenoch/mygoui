# Effects

`github.com/Cyenoch/mygoui/effect` applies three shader effects to an
application's own elements. Root [DESIGN.md](../DESIGN.md) ("Effects") says
when each is appropriate; this document owns the API and its limits.
[scrollfade.md](scrollfade.md) covers the scroll-driven edge fade and
[behavior/effect.md](behavior/effect.md) the behavior contract.

| Effect | Use it for | Fallback |
| --- | --- | --- |
| `effect.Mask` | Content cut at an edge that should fade into whatever is behind it | The edge is shaded instead |
| `effect.Shimmer` | Loading placeholders in your own layout, animated as one region | The placeholders rest |
| `effect.Reveal` | Content that arrives after the user asked for it | A plain fade |

For a blurred or translucent material use MyGo's own
`Element.Material(glass.Blur{})` (the dialog backdrop does) and
`Element.Blur`; for a glow use a theme shadow; for a gradient border use
`Element.BorderConicGradient`. None of them needs a shader.

## How an effect works

MyGo renders the element and its children into a layer every frame the
element paints and draws the layer through the effect. Only paint changes:
layout, hit testing, focus, focus rings and the accessibility tree stay with
the element. Say a state in a label, role or text as well; an effect is not
a signal on its own.

Each call takes the element after its children and styling are set and
returns it:

```go
func masked(c *ui.Context) {
	preview := ui.Column(c).Height(184).Radius(10).Clip()
	preview.Children(func() {
		for i := range 8 {
			ui.Text(c.Key(i), fmt.Sprintf("Release 2.%d", 8-i))
		}
	})
	effect.Mask(c, preview, effect.MaskOptions{Edges: effect.EdgeBottom, Radius: 10})
}
```

Limits every effect shares, read from MyGo v0.9.0 `ui/effect.go`:

- **One effect per element.** Applying an effect replaces the one the
  element had, including the one a mask `ScrollFade` or a `SkeletonGroup`
  installed. Wrap the element to combine two.
- **Layer size.** A layer side above 8192 pixels, device scale and zoom
  included, is not rendered: MyGo paints the element without its effect.
  Each effect states what it draws then.
- **Web views** inside the element do not show through a layer.
- **Cost.** The subtree is rendered into the layer on every painted frame.
  Shimmer and Reveal remove the effect when they have nothing to draw.

## The fallback

Without a GPU (`Painter.Accelerated()` is false: software renderers, most
virtual machines, `MYGO_GPU=0`, `ui.NewTester`) MyGo would interpret a
shader per pixel on the CPU, so the family does not run it. The same
fallback is drawn above the layer limit. `MaskOptions` and `RevealOptions`
have `Fallback: true` to draw it on a GPU too, which suits an application's
own "reduce effects" setting; for Shimmer, do not apply it. A shader is
compiled the first time it paints, so a fallback compiles nothing.

| Effect | No GPU, or `Fallback` | Above the layer limit |
| --- | --- | --- |
| Mask | Content whole; each selected edge shaded over the mask depth with the foreground at 18%, or blended into `Surface` | The same |
| Shimmer | Content at rest; no frames requested | The same |
| Reveal | Content fades in uniformly over the same duration | Content shows at once |

## Mask

```go
func Mask(c *ui.Context, e ui.Element, o MaskOptions) ui.Element
```

| Option | Zero value | Meaning |
| --- | --- | --- |
| `Edges` | every edge | Physical `EdgeTop`, `EdgeRight`, `EdgeBottom`, `EdgeLeft` flags |
| `Size` | `theme.Effects().FadeEdge` (24 DIP) | Fade depth in DIPs, at most half of the element |
| `Surface` | shadow fallback | The opaque colour behind the element; the fallback blends the edge into it. A translucent colour panics |
| `Radius` | 0 | The element's corner radius, which the fallback's shading follows |
| `Fallback` | false | Draw the fallback on a GPU too |

The mask multiplies what the element rendered, so rounded corners and
translucent fills survive it and the real background shows through: a
gradient, an image, another window material. The layer holds the element's
box, so content overflowing the box is cut there; Mask clips the element so
the fallback cuts it too. Mask takes the element's `DrawOver` callback for
the fallback's shading. Edges are physical and do not mirror in
right-to-left layouts. A scroll viewport uses
[`ScrollFade`](scrollfade.md), which follows the scroll position.

## Shimmer

```go
func Shimmer(c *ui.Context, e ui.Element, o ShimmerOptions) ui.Element
```

| Option | Zero value | Meaning |
| --- | --- | --- |
| `Pulse` | false | Fade the content to half and back instead of sweeping a band |

The loop and the band's strength are the theme's (see Tokens); there is no
per-call override.

```go
func profile(c *ui.Context, loading bool) {
	placeholder := ui.Column(c).Gap(8).Label("Loading profile")
	placeholder.Children(func() {
		mygoui.Skeleton(c, mygoui.SkeletonOptions{Height: 12, Animation: mygoui.SkeletonNone})
		mygoui.Skeleton(c, mygoui.SkeletonOptions{Width: 120, Height: 12, Animation: mygoui.SkeletonNone})
	})
	if loading {
		effect.Shimmer(c, placeholder, effect.ShimmerOptions{})
	}
}
```

The shimmer runs while it is applied; stop applying it when the content
arrives. It marks the element busy (`Element.Busy`), also when the band
rests; name the region, as the example does. It requests frames (60 a
second at most) only on the shader path. Under reduced motion, high
contrast or a `ShimmerPeriod` of zero or less it returns the element
busy but with no layer. The band is positioned in the
window, so every shimmering region moves as one. It lightens everything the
element paints: apply it to the element that holds the placeholder shapes,
not to the surface around them. `mygoui.SkeletonGroup` is a column under
this effect; a component with a `Loading` flag draws its own placeholders.

## Reveal

```go
func Reveal(c *ui.Context, e ui.Element, o RevealOptions) ui.Element
```

| Option | Zero value | Meaning |
| --- | --- | --- |
| `Direction` | `RevealDown` | `RevealDown`, `RevealUp`, `RevealForward` (from the leading edge, mirrored in right-to-left) or `RevealBackward` |
| `Fallback` | false | Fade in without a travelling front |

The duration, the curve and the depth of the fading front are the theme's
(`Motion.Reveal`, `Motion.Ease`, `Effects.FadeEdge`).

```go
func summary(c *ui.Context, arrived bool, text string) {
	if !arrived {
		return
	}
	block := ui.Column(c).Gap(4)
	block.Children(func() { ui.Text(c, text) })
	effect.Reveal(c, block, effect.RevealOptions{})
}
```

The reveal plays once, when the element is first built, and then returns
the element untouched: no layer and no frames at rest. Apply it to content
that arrives, not to what a window opens with. An element built under a new
`c.Key` reveals again, so do not put it on virtualized rows, which are
rebuilt as they scroll into view. Under reduced motion the element fades in
over the same duration with no travelling front.

## Tokens

| Token | Default | Read by |
| --- | --- | --- |
| `Effects.FadeEdge` | 24 DIP | Mask depth, Reveal front, `ScrollFade` |
| `Effects.ShimmerPeriod` | 2 s | Shimmer, `Skeleton`, `SkeletonGroup`; zero or less keeps them still (`theme.Use` accepts it, unlike the spinner and sweep periods) |
| `Effects.ShimmerPeak`, `Effects.ShimmerPeakDark` | 0.64, 0.04 | The band's strength in each appearance |
| `Motion.Reveal`, `Motion.Ease` | 200 ms, ease-out | Reveal |

Depths are DIPs: they follow content zoom and device scale, not the text
scale.

## Testing

`ui.NewTester` has no GPU, so a test through the public calls sees the
fallback: shaded mask edges, a resting shimmer, a fading reveal, and
unchanged focus and accessibility. The library's own tests also run each
shader through MyGo's CPU interpreter. Neither establishes what a GPU
presents; record that per platform in [VALIDATION.md](../VALIDATION.md).

For an effect the family does not have, write one with MyGo's `ui/shader`;
the `mygoui-effects` skill (`skills/mygoui-effects/CUSTOM-EFFECT.md`)
carries a compiled example that draws its content plain without a GPU.
