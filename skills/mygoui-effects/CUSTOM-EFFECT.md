# A custom mygoui effect

Use this only when `effect.Mask`, `effect.Shimmer` and `effect.Reveal` do not cover the job and `DESIGN.md` ("Effects") does not reject the treatment.

Write a `ui/shader` source with a `ColorEffect`, `DistortionEffect` or `LayerEffect` entry function (`go doc github.com/egoist/mygo/ui/shader` describes the language). Compile it once per process, keep per-element uniforms with `ui.Local`, and draw the content plain where `Painter.Accelerated()` is false or the effect has nothing to change:

```go
var desaturateShader = sync.OnceValue(func() *shader.Shader {
	return shader.MustCompile(`package main

var Amount float

func ColorEffect(pos vec2, color vec4) vec4 {
	gray := dot(color.rgb, vec3(0.2126, 0.7152, 0.0722))
	return vec4(mix(color.rgb, vec3(gray), Amount), color.a)
}
`)
})

type desaturate struct {
	values *shader.Values
	amount shader.Uniform
	effect shader.Effect
	level  float32
}

func newDesaturate() desaturate {
	s := desaturateShader()
	values := s.NewValues()
	return desaturate{values: values, amount: s.Uniform("Amount"), effect: s.Effect(shader.Options{Values: values})}
}

func (fx *desaturate) EffectMargin() ui.Point { return ui.Point{} }

func (fx *desaturate) PaintEffect(p *ui.Painter, box ui.Rect, content *ui.Bitmap, inset ui.Point) {
	if !p.Accelerated() || fx.level == 0 {
		s := p.Scale()
		p.ImageRegion(content, ui.Rect{X: inset.X, Y: inset.Y, W: box.W * s, H: box.H * s}, box, ui.ImageOptions{})
		return
	}
	fx.values.Float(fx.amount, fx.level)
	fx.effect.PaintEffect(p, box, content, inset)
}

type desaturateKey struct{}

// Desaturate draws e and its children through the effect.
func Desaturate(e ui.Element, level float32) ui.Element {
	fx := ui.Local(e, desaturateKey{}, newDesaturate)
	fx.level = level
	return e.Effect(fx)
}
```

- `MustCompile` reports a bad source only when it first runs, so build the effect once in a test.
- Colors in a shader are premultiplied.
- An effect that reads beyond the box sets `shader.Options.MaxSampleOffset`, and its wrapper returns the inner effect's `EffectMargin()`.
- MyGo does not call `PaintEffect` when a layer side exceeds 8192 pixels; the element paints without the effect. A fallback that must show there is painted by the element itself (`DrawOver`).
- An element has one effect: this one replaces a family effect or `MaskFade` on the same element.
- An animated effect asks for frames from `PaintEffect` (`p.AnimationFrameRate(60)`) only while it moves, and keeps still under `theme.ReduceMotion(c)`.
