# Design

## Context

See proposal.md for the motivation. Constraints read from the source at the
pinned MyGo v0.9.0, not assumed:

- `ui/effect.go`: `Element.Effect(fx)` renders the element and its children
  into a bitmap every frame it paints and hands it to `fx.PaintEffect`. An
  element has one `effect` field, so a second `Effect` replaces the first.
  `Painter.layer` returns no canvas when a side exceeds `maxEffectLayer`
  (8192 pixels, after scale and the effect margin); `paintEffect` then
  returns false and the element paints as it is. **`PaintEffect` is not
  called in that case**, so a fallback cannot live in the effect. Web views
  inside the element do not show through the layer.
- `ui/shader/shader.go`: a shader runs on the GPU and is interpreted per
  pixel on the CPU where there is none (`pixels.Color`).
  `Painter.Accelerated()` is false in `ui.NewTester`, `Render`, with
  `MYGO_GPU=0` and on software renderers.
- `ui/motion.go`: `Element.AnimateTo` keeps frames coming only while its
  value moves, starts at its first target unless `Jump` set it elsewhere,
  and with `NonMotion` keeps animating under reduced motion.
- The dialog backdrop already blurs natively:
  `overlay/chrome.go` fills it with `glass.Blur` (`Element.Material`).
  `Element.Blur` blurs a subtree. The Tour spotlight is a path fill
  (`paintSpotlight`), not a shader. `Element.BorderConicGradient` and
  `Painter.StrokePathConicGradient` draw gradient borders natively.
- `DESIGN.md`: compact, neutral, restrained state feedback; button press
  uses colour and shadow; nothing requests frames at rest; under reduced
  motion fades stay while positions and loops settle.

The design work followed Impeccable's `animate` playbook in Operate mode
(the gallery brief is `.impeccable/surfaces/internal-gallery.md`). The
structured-question probe is unavailable to a subagent running a brief, so
direction was inferred from the brief, `PRODUCT.md` and `DESIGN.md`.

## Goals / Non-Goals

**Goals:**

- Three effects, each complete: options, tokens, fallback, motion rule,
  composition, tests, gallery, records.
- One compiled shader per effect in the process, shared by the public family
  and by `ScrollFade` and `SkeletonGroup`.
- A fallback that is designed and reachable in the gallery.

**Non-Goals:**

- A general shader toolkit or effect composition (stacking two effects on
  one element). The custom-shader path stays MyGo's `ui/shader`, documented
  in the skill.
- Decorative effects with no state to carry.
- GPU presentation evidence. The headless suite cannot give it.

## Decisions

### The set

Motion thesis (Impeccable `animate`, Operate mode): there is no focal
moment; an effect earns its place by making one state legible. Feedback
budget: one looping effect (shimmer) that stops with its state, one
one-shot (reveal), one static (mask).

| Effect | Job in an application | Why a shader |
| --- | --- | --- |
| `Mask` | Content that continues past an edge and is cut there: a preview of a longer list, a clipped line of chips, a ticker, a truncated block over an image or gradient | Only an alpha multiply on the rendered subtree lets the real background show through; overpainting a surface colour fails on anything but one flat colour |
| `Shimmer` | Loading placeholders in the application's own layout (a card grid, a chart frame, a form), animated as one region rather than per shape | One pass over the rendered subtree keeps a single band in phase across every shape and across regions; per-shape gradients drift and cost one draw each |
| `Reveal` | Content that arrives after the user asked for it: a loaded image, a generated block, a finished chart | A soft-edged directional wipe is an alpha ramp that moves; native clipping gives only a hard edge and native opacity only a uniform fade |

Rejected candidates:

- **Glow or spotlight for attention**: theme shadows and `Painter.Shadow`
  already draw a glow without a layer, and the Tour's spotlight is a path
  fill; a shader adds a layer and nothing else.
- **Grain or noise**: texture with no state to carry; the language is
  neutral surfaces, and it would run per pixel on the CPU fallback.
- **Progressive blur or tinted material**: MyGo provides the material
  (`Element.Material(glass.Blur{})`, used by the dialog backdrop) and
  `Element.Blur`; `docs/effects.md` points there. A progressive variant
  needs the backdrop-reading shader path, which the headless suite cannot
  exercise, and no component asks for it.
- **Gradient or animated border**: `Element.BorderConicGradient` and
  `Painter.StrokePathConicGradient` draw it natively; as a focus signal it
  would compete with the native focus ring.
- **Ripple**: contradicts the recorded press rule (colour and shadow, no
  pressed scale) and imports another platform's idiom.
- **Sheen for success**: success is already an icon swap and a colour fade
  (State motion); a sweep would be a second signal for one state.
- **Dissolve (noise threshold)**: the same job as Reveal with a louder
  material; one reveal is enough.

### Package: a new `effect` family, shared code in `internal/fx`

`effect` is a leaf: it imports `theme` and `internal/fx` only. The root
imports `effect` (`SkeletonGroup` is a column under `effect.Shimmer`). The
dependency rule gains one line and stays acyclic.

Alternative considered: the module root. Rejected: the root would grow by
three functions, four types and their constants named `Effect*`, and other
families that want an effect would still import the root. A family package
gives `effect.Mask`, `effect.Shimmer`, `effect.Reveal`.

`internal/fx` holds what both the root and `effect` draw with: the three
shaders (each behind `sync.OnceValue`), plain and faded content drawing, the
edge-gradient painter and the layer-limit test. `ScrollFade` needs it
directly: its mask lengths come from scroll state at paint time, which a
build-time `effect.Mask` call cannot express, so it keeps its own
`ui.Effect` and shares the shader and the edge painter. Adding a paint-time
callback to `MaskOptions` for one internal caller was rejected.

### API shape

```go
func Mask(c *ui.Context, e ui.Element, o MaskOptions) ui.Element
func Shimmer(c *ui.Context, e ui.Element, o ShimmerOptions) ui.Element
func Reveal(c *ui.Context, e ui.Element, o RevealOptions) ui.Element
```

Each takes the element last, after its children and styling are set, as
`ScrollFade` does, and returns it. Nothing build-scoped is retained: per
element state (uniform values, the forced-GPU test flag) lives in
`ui.Local`; paint callbacks capture that state and plain values.

`MaskOptions` and `RevealOptions` carry `Fallback bool`: draw the fallback
even on a GPU. It serves an application's "reduce effects" setting and lets
the gallery show the fallback on a GPU machine. Shimmer has none: its
fallback is the element at rest, which is what not applying it draws.

Options exist only where an application case needs a per-call value:
Mask's edges, depth, surface and radius, Shimmer's `Pulse`, Reveal's
direction. Shimmer's period and strength and Reveal's duration and front
are theme tokens with no per-call override (review finding: no caller or
example needed one).

Mask caps each depth at half of the element on both paths, so opposite
fades meet in the middle of a small element instead of multiplying into a
translucent whole (review finding: a 100 by 32 DIP element kept 44% alpha
at its centre). `ScrollFade` keeps its own depths.

The `fx` shader types are ready as zero values and compile when they first
paint, so a build, a fallback and a frame without a GPU compile nothing. A process-wide switch was
rejected: it would be library-owned global state.

### Fallbacks

The no-GPU fallback keeps the existing rule: a shader interpreted per pixel
every frame costs more than the effect is worth, so the shader does not run.

| Effect | No GPU, or `Fallback` | Above the layer limit |
| --- | --- | --- |
| `Mask` | The content, with each selected edge shaded over the mask depth: the foreground at 18% fading inwards, or a blend into `Surface` when the call names an opaque one | The same |
| `Shimmer` | The content at rest; no frames | The same |
| `Reveal` | The content fading in uniformly over the same duration | The content at once |

Because `PaintEffect` is skipped above the limit, Mask's shading is painted
by the element itself (`DrawOver`), guarded by "no GPU, `Fallback`, or the
box does not fit a layer". It is therefore inside the layer when the shader
does not run and absent when it does. Mask takes `Radius` so that shading
follows rounded corners; it does not call `Clip`, which would cut child
focus rings and re-enter the nested-clip limit recorded in `docs/design.md`.
On the shader path corners need nothing: the mask multiplies what the
element already rendered.

Reveal fades through the effect (the layer drawn with a tint) rather than
through `Element.Opacity`, so it does not take over an opacity the caller
set. Above the limit there is no layer to tint, so the content shows at
once; that is stated, not hidden.

`ScrollFade`'s mask gets the same layer-limit guard on its existing shadow
layer. This is the one rendering change to an existing component and only
for viewports above 8192 pixels a side.

### Motion

- **Shimmer** is a paint loop (`Painter.Now`, `AnimationFrameRate(60)`),
  asked for only inside `PaintEffect` on the shader path. Reduced motion,
  high contrast and a non-positive period return the element without an
  effect, so no layer is rendered either. It stops when the caller stops
  applying it.
- **Reveal** is a build-level value: the first build jumps it to 0, then
  `AnimateTo(1)` over `Motion.Reveal` with `Motion.Ease` and `NonMotion`.
  MyGo keeps frames coming only while it moves. At 1 the element is
  returned without an effect. Under reduced motion the value still runs and
  the effect draws a uniform fade: the Reduced Motion Rule keeps fades and
  settles positions.
- Reveal plays when the element is first built, as a native enter
  transition does. Content the window opens with should not be given one
  ("a value set by the application on the first frame does not animate
  in"); that is the caller's choice and the documentation says so. A
  two-way `Hidden` flag was rejected: a hidden but built element stays
  focusable and in the accessibility tree.

### Tokens

`Effects.ShimmerPeriod` (2 s), `Effects.ShimmerPeak` (0.64) and
`Effects.ShimmerPeakDark` (0.04) replace the constants in `skeleton.go`;
`Skeleton`, `SkeletonGroup` and `effect.Shimmer` read them, so one decision
changes every placeholder. Mask depth and the reveal's soft front read the
existing `Effects.FadeEdge`; the reveal reads `Motion.Reveal` and
`Motion.Ease`. The 18% fallback shadow stays where `ScrollFade` has it.
Sizes are DIPs: they follow content zoom through the painter's scale and do
not follow text scale, as `ScrollFade`'s depth does not.

### Testing the shader path headless

`ui.NewTester` never reports a GPU, so the public calls reach only the
fallback. Each effect's `ui.Local` state has an unexported `forced` flag;
package-internal tests fetch that state for the element under test and set
it, which makes the frame behave as a GPU frame and runs the shader through
MyGo's CPU interpreter. This is the existing `forcedMask` technique made to
go through the public entry point. It shows the shader's output under the
interpreter; it does not show GPU presentation.

## Risks / Trade-offs

- [GPU output differs from the interpreter] → Not mitigated by this change;
  recorded as unverified. `VALIDATION.md` receives GPU evidence separately.
- [Mask uses the `DrawOver` slot] → Stated in the contract and docs; the
  caller composes by wrapping the element.
- [A theme sets `ShimmerPeriod` to zero or less; `theme.Use` accepts it,
  unlike the spinner and sweep periods] → It keeps placeholders still rather than
  dividing by zero; documented in `docs/theme.md` and the changelog.
- [Reveal replays when a keyed element leaves and re-enters, such as a
  virtualized row] → Documented: apply it to the arriving region, not to
  recycled rows.
- [The shader path cuts overflowing content at the element's box and the
  fallback does not] → Documented: give the element `Clip()`; Mask does not
  call it for the caller (see Fallbacks).
- [One effect per element] → Applying a family effect to a `MaskFade`
  viewport or a `SkeletonGroup` replaces theirs; stated in the docs.
- [An effect layer re-renders the subtree every painted frame] → Shimmer
  and Reveal drop the effect when they have nothing to draw; Mask is static
  and paints only when the element does.
