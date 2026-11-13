# Proposal

## Why

mygoui holds two shaders, both private: the alpha mask behind `ScrollFade`'s
`MaskFade` and the `SkeletonGroup` shimmer. An application author who wants
either treatment on their own element writes MyGo `ui/shader` code from
scratch, repeats the GPU-less fallback and the layer limit by hand, and finds
no rule in `DESIGN.md` for when an effect is appropriate. Effects are the
missing public part of the "complete native UI solution" recorded in
`PRODUCT.md`.

## What Changes

- Add the public `effect` family package with three effects an application
  applies to its own elements: `effect.Mask` (edges fade to transparency),
  `effect.Shimmer` (one loading highlight across a placeholder subtree) and
  `effect.Reveal` (arriving content wipes in once). Each takes a typed
  options struct whose zero value reads the published theme.
- Give every effect a designed fallback, drawn where the frame has no GPU,
  where the element exceeds MyGo's 8192-pixel effect layer, and on request
  for Mask and Reveal (`Fallback`): Mask shades its edges, Shimmer rests,
  Reveal fades.
- Add the effect tokens `Effects.ShimmerPeriod`, `Effects.ShimmerPeak` and
  `Effects.ShimmerPeakDark`; `Skeleton`, `SkeletonGroup` and `effect.Shimmer`
  read them. Existing tokens and their defaults are unchanged.
- Move the two existing shaders into one shared implementation:
  `SkeletonGroup` becomes a column under `effect.Shimmer`, and `ScrollFade`'s
  mask and `effect.Mask` run the same compiled shader. Neither component
  renders differently, except that a mask `ScrollFade` larger than the layer
  limit now shows its shadow fallback on a GPU too.
- Record in `DESIGN.md` what each effect is for, when it is appropriate and
  which candidates were rejected; document the family (`docs/effects.md`,
  `docs/components.md`, `docs/behavior/effect.md`), add gallery pages, and
  lead the `mygoui-effects` skill with the public family.

## Capabilities

### New Capabilities

- `shader-effects`: the public effects an application applies to its own
  elements: what each draws, its fallback, its motion and reduced-motion
  treatment, its tokens, and what applying one leaves untouched.

### Modified Capabilities

- `native-library`: "Effects degrade without a GPU" also covers an effect
  layer above the toolkit's size limit.

## Impact

- New package `github.com/Cyenoch/mygoui/effect` and `internal/fx`; the root
  package imports `effect`. The dependency rule in `docs/design.md` gains
  one leaf.
- `theme.EffectTokens` gains three fields. A theme built from
  `theme.Default()` keeps today's look; a theme that sets
  `ShimmerPeriod` to zero keeps skeletons still.
- `scroll_fade.go` and `skeleton.go` lose their private shader sources.
- Gallery: three pages in Foundations, the API table and the token cards.
- Records: `DESIGN.md`, `PRODUCT.md`, `CHANGELOG.md`, `docs/`, the
  `mygoui-effects` skill.
- No new module dependency. GPU presentation is not established by the
  headless suite and stays a separate record in `VALIDATION.md`.
