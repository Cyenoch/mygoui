# Tasks

## 1. Shared shaders and tokens

- [x] 1.1 Add `internal/fx` with the mask, shimmer and wipe shaders (each compiled when it first paints), plain and faded content drawing, the edge-gradient painter and the layer-limit test; verified by `go build ./...` and by the effect tests in group 2 that run each shader
- [x] 1.2 Add `Effects.ShimmerPeriod`, `ShimmerPeak` and `ShimmerPeakDark` to `theme.EffectTokens` and `theme.Default()`; list them in `docs/theme.md` and the gallery token cards; verified by `go test ./theme ./internal/gallery`

## 2. The effect family

- [x] 2.1 Implement `effect.Mask` with `MaskOptions`, the shadow and surface fallbacks and the layer-limit guard; tests: the forced shader path lets the background through at a selected edge only, the software path and an oversized element show the shaded edge, a `FadeEdge` override changes the depth, a translucent `Surface` panics
- [x] 2.2 Implement `effect.Shimmer` with `ShimmerOptions`; tests: the forced shader path requests frames and lightens a band, reduced motion and the software path rest without a pending frame, a token override changes the result
- [x] 2.3 Implement `effect.Reveal` with `RevealOptions`; tests: the forced path shows the entering edge before the far edge and rests after the duration, the software path and reduced motion fade uniformly and rest, a later build does not replay
- [x] 2.4 Test that an effect leaves focus, activation and the accessible role and name of the element and its children untouched, and that a second effect replaces the first
- [x] 2.5 Render every effect headless in light and dark, at rest and mid-animation, plus each fallback, into the scratch directory; inspect the images and fix what looks wrong

## 3. Existing components on the shared implementation

- [x] 3.1 Move `ScrollFade`'s mask onto `internal/fx` and guard its shadow layer by the layer limit; the existing `scroll_fade*_test.go` pass unchanged and a new test covers the oversized viewport
- [x] 3.2 Build `SkeletonGroup` on `effect.Shimmer` and read the shimmer tokens in `Skeleton`; the existing `display_skeleton_test.go` passes unchanged

## 4. Gallery

- [x] 4.1 Add Mask, Shimmer and Reveal pages in Foundations with live examples and a Fallback switch; register the `effect` qualifier for snippets and the API table; run `make generate`; verified by `go test ./internal/gallery`

## 5. Records

- [x] 5.1 `DESIGN.md` Effects section: when each effect is appropriate, its tokens, its fallback, the rejected treatments
- [x] 5.2 `docs/effects.md` (API, limits, compiling examples), `docs/components.md` section, `docs/design.md` package table, dependency rule and effect boundary, `docs/scrollfade.md`, `docs/mygo-vs-css.md`, `docs/theme.md`
- [x] 5.3 `docs/behavior/effect.md` in the reference format, with its row in `docs/behavior/README.md`
- [x] 5.4 `PRODUCT.md` pillar table and capabilities, `CHANGELOG.md` Unreleased
- [x] 5.5 `skills/mygoui-effects/` (SKILL.md, CUSTOM-EFFECT.md, MAINTAINING.md) lead with the public family; every Go block in the docs and the skill compiles in the scratch module

## 6. Verification

- [x] 6.1 `make generate`, `make check`, `make race`, `golangci-lint run ./...`, the WebAssembly build, `make spec-check`, `make design-doctor` (findings read)
- [x] 6.2 GPU presentation: run the gallery on a GPU window and capture each effect and its fallback, or record that it stays unverified. Recorded as unverified: MyGo's automation `capture` rasterises the frame on the CPU (`ui/window.go`, `windowHost.capture`), so it is not GPU evidence, and an OS screen capture was not available to this run

## 7. Review findings

- [x] 7.1 Mask caps each depth at half of the element on the shader path as on the fallback; `TestMaskKeepsTheMiddleOfASmallElementOpaque` (fails without the cap)
- [x] 7.2 Gallery switch rows wrap; Shimmer and Reveal pages rendered at 400 px and text scale 1.75 and inspected
- [x] 7.3 Shimmer tests use a 96 by 8 scene and apply the shimmer after the clock has moved
- [x] 7.4 Remove `ShimmerOptions.Period`, `Peak` and `Fallback` and `RevealOptions.Duration` and `Soft` (theme tokens, no caller); docs, gallery table, skill and spec follow
- [x] 7.5 Forced-path tests for `Fallback` on Mask and Reveal; `docs/behavior/effect.md` cites them
- [x] 7.6 Shaders compile when they first paint; docs and spec say so
- [x] 7.7 Rename `fx.Sweep` to `fx.Shimmer`; `TestShimmerRestsForAZeroPeriod`; changelog and docs state what `theme.Use` accepts
- [x] 7.8 Shimmer examples mark the region busy; the doc comment says the caller does
- [x] 7.9 Gallery example descriptions that explained the treatment removed

## Workflow follow-up

- Archive the change after review (`npm run openspec -- archive native-effects`).
- Record GPU observations per platform in `VALIDATION.md` when they exist.
