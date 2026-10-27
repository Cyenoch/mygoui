# Native scrolling and edge fades

`ScrollArea` wraps a native viewport. `ScrollFade` can decorate an existing
tracked viewport without creating a second scroll owner. Both keep
`ui.ScrollState` with the application.

## ScrollArea

`ScrollArea(c, state, options, body)` takes a non-nil `*ui.ScrollState`.
`ScrollAreaOptions` contains `Axis`, `Height`, `Label` and an optional
`*ScrollFadeOptions`. `Axis` is `ScrollVertical` by default, with
`ScrollHorizontal` and `ScrollBoth` available. A positive `Height` sets
the viewport extent; otherwise its parent must constrain it.

```go
type pane struct {
    scroll ui.ScrollState
}

func (s *pane) View(c *ui.Context) {
    theme.Apply(c)
    mygoui.ScrollArea(c, &s.scroll, mygoui.ScrollAreaOptions{
        Axis: mygoui.ScrollVertical,
        Height: 240,
        Label: "Activity",
        Fade: &mygoui.ScrollFadeOptions{Mode: mygoui.ShadowFade},
    }, func() {
        for _, title := range []string{"Queued", "Building", "Published"} {
            ui.Text(c, title)
        }
    })
}
```

The sample uses the module root, `theme` and MyGo `ui` imports. Short content
has no overflow and therefore no edge fade. Scrollbars, wheel input, focus
and keyboard routing remain native; there is no copied DOM scrollbar skin.

## Existing viewports

Use the same scroll state on the viewport and its fade:

```go
viewport := ui.Scroll(c).TrackScroll(&s.scroll).Height(240)
viewport.Children(func() {
    // Build the scrolling content here.
})
mygoui.ScrollFade(c, viewport, &s.scroll, mygoui.ScrollFadeOptions{
    Mode: mygoui.MaskFade,
    Edges: mygoui.FadeTop | mygoui.FadeBottom,
    Size: 24,
})
```

The viewport must belong to the current build. Shadow and surface fades paint
in a non-interactive absolute layer inside its rounded content clip, before
its border and native scrollbars. Build scrolling content before decorating
the viewport. Radius and border can still be set on the returned element.

`MaskFade` reduces the
viewport's own paint and its native descendants to transparency, revealing
the actual background behind it. Apply it to the layer whose paint should
fade; an opaque background on that layer fades with the children.

## Modes and contracts

| Option | Meaning |
| --- | --- |
| `ShadowFade` | Default. Paints shadow gradients; `Color` customizes the ink. Zero-alpha color selects the theme-derived default. |
| `SurfaceFade` | Paints gradients into `Surface`, which must be opaque. Images and mixed backgrounds cannot be reconstructed by one surface color. |
| `MaskFade` | Fades the viewport and its native descendants to transparency with a `ui/shader` ColorEffect. Needs neither Surface nor Color. |
| `Edges` | Physical `FadeTop`, `FadeRight`, `FadeBottom`, `FadeLeft` flags. Zero enables all four. |
| `Size` | Fade depth in DIP. Non-positive selects 24. Painted shadow/surface depth clamps to half the viewport; mask depth follows the selected edge's remaining distance up to Size. |

Fades reflect remaining scroll distance independently at each edge and
disappear where that edge has no overflow. They add paint, not hit targets
or a background timer. Nil state, an unsupported mode/axis and a translucent
`SurfaceFade.Surface` fail explicitly.

`MaskFade` is a true alpha treatment, unlike shadow or surface overpainting.
It uses `ui.Element.Effect` and `ui/shader`, which execute on GPU renderers
and have a CPU fallback. The effect changes paint, preserving native hit
testing, focus and accessibility. The viewport has one effect slot:
`MaskFade` replaces an existing effect installed on that same element.

The pinned toolkit caps each effect-layer side at **8192 pixels**, including
device scale and margins. Oversized layers paint without the effect.
Embedded webviews do not appear through an effect layer. Rendering a native
subtree into a layer also adds render work; this documentation does not claim
measured GPU cost or completed platform verification.

A meaningful transparency test places the viewport above a patterned
background and checks that the pattern shows through faded content. Testing
only against a matching solid color cannot distinguish an alpha mask from
surface overpainting. Headless fallback tests and real GPU evidence belong
in separate entries in `VALIDATION.md`.

Native implementation references at the pinned MyGo v0.7.1:
[`ui/effect.go`](https://github.com/egoist/mygo/blob/v0.7.1/ui/effect.go)
defines subtree rendering and layer limits;
[`ui/shader/shader.go`](https://github.com/egoist/mygo/blob/v0.7.1/ui/shader/shader.go)
defines ColorEffect and native shader execution.

`ScrollShadow` remains an extra convenience viewport over `SurfaceFade`.
Its default depth is 40 DIP and its zero surface selects the page background.
It is not an additional coss catalog component or a true mask implementation.
