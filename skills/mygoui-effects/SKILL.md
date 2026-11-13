---
name: mygoui-effects
description: "Shader effects and edge fades in mygoui. Use when adding, choosing or debugging an effect."
---

# mygoui shader effects

For work on the mygoui repository itself, also read [MAINTAINING.md](MAINTAINING.md).

`docs/effects.md` owns the `effect` package's options, limits and examples, `docs/scrollfade.md` the scroll fade modes, `DESIGN.md` ("Effects") when an effect is appropriate; all ship with the module version your application uses (`go list -m -f '{{.Dir}}' github.com/Cyenoch/mygoui`; when vendored, the repository at that version). Signatures: `go doc github.com/Cyenoch/mygoui/effect`.

## Choose

| The interface needs | Use |
| --- | --- |
| A scroll viewport that fades where more can be scrolled | `mygoui.ScrollArea` (option `Fade`) or `mygoui.ScrollFade`; modes below |
| Static content cut at an edge, fading into a varied background | `effect.Mask` |
| Placeholders in your own layout that wait as one region | `effect.Shimmer`, or `mygoui.SkeletonGroup` for a plain column of them |
| Content that arrives after the user asked for it | `effect.Reveal` |
| Blur or a translucent material | MyGo's `Element.Material(glass.Blur{})` or `Element.Blur`; not an effect |
| A glow, or a gradient border | A theme shadow; `Element.BorderConicGradient`; not an effect |
| Anything else | A custom shader: [CUSTOM-EFFECT.md](CUSTOM-EFFECT.md) |

An effect makes one state legible and is never the only signal: say the state in a label, role or text too. The design language has no press ripple, grain, success sheen or animated border; do not add one with a custom shader.

## The effect family

Build the element with its children and styling, then pass it to the effect, which returns it:

```go
func results(c *ui.Context, loading bool, rows []string) {
	if loading {
		placeholder := ui.Column(c).Gap(8).Label("Loading results")
		placeholder.Children(func() {
			for i := range 3 {
				mygoui.Skeleton(c.Key(i), mygoui.SkeletonOptions{Height: 12, Animation: mygoui.SkeletonNone})
			}
		})
		effect.Shimmer(c, placeholder, effect.ShimmerOptions{})
		return
	}
	arriving := ui.Column(c)
	arriving.Children(func() {
		list := ui.Column(c).Gap(8).Height(120).Clip()
		list.Children(func() {
			for i, row := range rows {
				ui.Text(c.Key(i), row)
			}
		})
		effect.Mask(c, list, effect.MaskOptions{Edges: effect.EdgeBottom})
	})
	effect.Reveal(c, arriving, effect.RevealOptions{})
}
```

- **Zero options are the designed default.** Mask depth comes from `theme.Effects()` unless `Size` is set. Shimmer's period and strength and Reveal's duration, curve and front come only from the theme; change the theme to change them.
- **Fallback.** The shader runs only on a GPU and below the layer limit. Otherwise Mask shades its edges (or blends them into `Surface`), Shimmer rests and Reveal fades uniformly. `Fallback: true` on Mask and Reveal draws that everywhere; wire it to your own "reduce effects" setting, and skip the `Shimmer` call under the same setting.
- **Mask** cuts overflowing content at the element's box on the shader path only: give the element `Clip()` so the fallback cuts it too. Pass the element's corner radius as `Radius`. It uses the element's `DrawOver`.
- **Shimmer** lightens everything the element paints: apply it to the element holding the placeholders, not the card around them. It tells assistive technology nothing: mark the region busy (`Element.Busy`) and name it. Stop applying it when the content arrives. It rests under reduced motion and high contrast.
- **Reveal** plays once, when the element is first built. Do not put it on what a window opens with, or on virtualized rows, which are rebuilt as they scroll into view.

## Limits every effect shares

- **One slot.** An element has one effect; a second replaces the first, including the one `MaskFade` or `SkeletonGroup` installed. Wrap the element to combine two, as the example does.
- **8192-pixel layers.** A layer side larger than 8192 pixels, counting device scale and zoom, paints without the effect; the family draws its fallback there (Reveal shows the content at once).
- **Webviews** inside the element do not show through the layer.
- **Cost.** The subtree renders into a layer every frame it paints.

## Edge fades on a scroll viewport

| Mode | Use when |
| --- | --- |
| `ShadowFade` (default) | Content scrolls over anything; a soft shadow marks the edge |
| `SurfaceFade` | The viewport sits on one known opaque color, passed as `Surface` |
| `MaskFade` | The content must really fade to transparency over a varied background; `ShadowFade`'s overlay stands in without a GPU or above the layer limit |

A fade shows only at an edge with more to scroll and adds no timer. Apply `MaskFade` to the layer whose paint should fade: an opaque background on that layer fades with its children.

## What a test proves

`ui.NewTester` paints without a GPU, so it exercises the fallback, not the shader: shaded mask edges, a resting shimmer, a fading reveal, unchanged hit testing, focus and accessibility, and a view that rests. A real transparency check places the faded content over a patterned background on a GPU window; record that, per platform, apart from headless results.
