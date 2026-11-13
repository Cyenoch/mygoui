# Effect behavior

Package `effect`. Format and evidence marks: [README](README.md). When each
effect is appropriate and its tokens: [`DESIGN.md`](../../DESIGN.md),
"Effects". API and limits: [`effects.md`](../effects.md).

Shared by every effect below: it changes paint only, except that Shimmer
marks its element busy. The element and its children keep their role, name,
Tab order, focus and activation.
[T TestEffectsLeaveFocusActivationAndAccessibilityAlone] An element holds one
effect; a later one replaces it. [T TestALaterEffectReplacesAnEarlierOne]
Tests marked here that name a shader run it through MyGo's CPU interpreter;
none shows what a GPU presents.

### Mask

| | |
| --- | --- |
| Role and name | None of its own; the element's are unchanged. [T TestEffectsLeaveFocusActivationAndAccessibilityAlone] |
| Pointer | Faded content still takes clicks; the mask is not a hit target. [C] |
| States | On the shader path the background shows through at the selected edges only. [T TestMaskLetsTheBackgroundThroughAtSelectedEdges] Content overflowing the element is cut at its edges on every path. [T TestMaskFallbackKeepsOverflowInside] Without a GPU the content stays whole and the selected edges are shaded, towards `Surface` when one is given. [T TestMaskWithoutGPUShadesSelectedEdges] `Fallback` draws the same where the shader could run. [T TestMaskFallbackShadesWhereTheShaderCouldRun] Opposite fades reach at most half of the element, so a small one keeps an opaque middle. [T TestMaskKeepsTheMiddleOfASmallElementOpaque] An element too large for a layer is shaded the same way. [T TestMaskShadesItsEdgeAboveTheLayerLimit] The depth follows `Effects.FadeEdge` unless `Size` is set. [T TestMaskDepthFollowsTheFadeEdgeToken] A translucent `Surface` panics. [T TestMaskRejectsATranslucentSurface] Edges are physical and do not mirror in right-to-left. [C] |
| Motion | None; it requests no frames. [C] |

### Shimmer

| | |
| --- | --- |
| Role and name | None of its own; the element is busy while Shimmer applies, under reduced motion too. [T TestShimmerMarksItsRegionBusy] The caller names it. [C] |
| States | Without a GPU the content rests. [T TestShimmerWithoutGPURests] The band's strength follows `Effects.ShimmerPeak`. [T TestShimmerStrengthFollowsThePeakToken] A `ShimmerPeriod` of zero rests. [T TestShimmerRestsForAZeroPeriod] |
| Motion | While applied on the shader path a band sweeps the content and frames are requested; under reduced motion or high contrast it rests and requests none. [T TestShimmerSweepsWhileAppliedAndRestsWhenAskedTo] It stops in the build that stops applying it. [T TestShimmerStopsWithItsState] `Pulse` fades the content to half and back instead. [C] |

### Reveal

| | |
| --- | --- |
| Role and name | None of its own; the content is in the accessibility tree and focusable from its first build, before the reveal ends: the wipe is decoration, and arriving content is not held back from the keyboard or assistive technology for it. [T TestEffectsLeaveFocusActivationAndAccessibilityAlone] |
| States | Without a GPU the content fades in uniformly, then rests. [T TestRevealWithoutGPUFadesInThenRests] `Fallback` draws the same fade where the shader could run. [T TestRevealFadesUnderReducedMotionAndFallback] Above the layer limit it shows at once. [C] `RevealForward` and `RevealBackward` mirror in right-to-left. [C] |
| Motion | Plays once when the element is first built: the entering edge shows before the far edge, then the element rests with no effect and no frames; a later build does not replay it. [T TestRevealWipesInFromItsEdgeOnceThenRests] Under reduced motion it fades with no travelling front. [T TestRevealFadesUnderReducedMotionAndFallback] An element built again under a new key reveals again. [C] |

## Gaps

- Mask, toolkit limit: the fallback shading takes the element's `DrawOver`
  slot, so a caller's own `DrawOver` on that element is replaced. MyGo keeps
  one callback per element (`paintAfterFn`, ui/element.go) with no getter to
  chain it.
- GPU presentation of the three shaders is not verified on any platform.
