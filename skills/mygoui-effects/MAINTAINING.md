# Maintaining mygoui shader effects

For work inside the mygoui repository. Paths are relative to its root.

## Internal contract

- `internal/fx` holds the three shaders (`Mask`, `Shimmer`, `Wipe`), each compiled once behind `sync.OnceValue` the first time it paints (the zero value of each type is ready, so a build or a fallback compiles nothing), plus `Plain` and `Faded` content drawing, `PaintEdges` (the fallback gradients) and `Fits` (the layer-limit test, mirroring MyGo's `maxEffectLayer`).
- `effect` (`effect/mask.go`, `shimmer.go`, `reveal.go`) is the public family. It imports only `theme` and `internal/fx`; the root imports it (`SkeletonGroup` is a column under `effect.Shimmer`).
- `ScrollFade`'s `MaskFade` keeps its own `ui.Effect` (`scroll_fade.go`, `scrollMask`): its mask lengths come from scroll state at paint time. It shares `fx.Mask`, `fx.PaintEdges` and `fx.Fits`.
- MyGo does not call `PaintEffect` above the layer limit. A fallback that must show there is painted by the element: Mask through `DrawOver`, `MaskFade` through its shadow layer. Both paint when there is no GPU, or the box does not fit a layer.
- Per-element state lives in `ui.Local` and carries an unexported `forced` flag. Internal tests set it after the public call so the frame takes the shader path through MyGo's CPU interpreter.
- An effect that has nothing to draw returns the element without `Effect`, so no layer is rendered: Shimmer under reduced motion, high contrast or a non-positive `ShimmerPeriod`; Reveal once it has finished; Mask under `Fallback`.

## Records

- `docs/effects.md` owns the family's API and limits; `DESIGN.md` ("Effects") when each effect is appropriate and what was rejected; `docs/behavior/effect.md` the behavior contract; `docs/scrollfade.md` the fade modes; `docs/mygo-vs-css.md` the paint, clip and effect model; `docs/theme.md` the tokens.
- The MyGo source at the pinned version (`go list -m -f '{{.Dir}}' github.com/egoist/mygo`) defines the layer limit (`ui/effect.go`, `maxEffectLayer`) and the shader entry points (`ui/shader/shader.go`). Check them again on every MyGo upgrade (`docs/mygo-0.9.md`), and `fx.Fits` with them.
- Tests: `effect/effect_test.go` covers the fallbacks through the public calls; `effect/shader_test.go` the shader path with `forced`; `scroll_fade_mask_test.go`, `scroll_fade_shader_test.go` and `scroll_fade_limit_test.go` the scroll mask. `MYGOUI_EFFECT_SPECIMEN_DIR=<dir> go test ./effect -run TestEffectSpecimens` writes a PNG of every effect, state and fallback in both appearances; look at them after changing a shader or a fallback.
- None of these tests shows GPU presentation. GPU and platform observations go in `VALIDATION.md`, separate from headless results.

## Keep this skill current

When an effect, an option, a token, a fallback, the layer limit or the webview rule changes, update `SKILL.md` and `CUSTOM-EFFECT.md` in the same change and list them in that change's tasks. Compile every Go block against the checkout, and run the `CUSTOM-EFFECT.md` shader once in a Tester after editing it: `MustCompile` fails only at run time.
