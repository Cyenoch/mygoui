# Proposal

## Why

Theme tokens are one pillar of the native UI solution, and root `DESIGN.md`
is now the design authority. A later design language must replace the
built-in one without editing component code. Today only colour, radius, type
and elevation are public in `theme`; motion lives in the private
`internal/motion` and in duration and easing literals across component
files, the spacing unit is a constant, effect parameters are literals, and
no value names a whole theme, so the built-in look is "the theme" instead of
"the default theme".

## What Changes

- `theme` gains public motion tokens (durations and easings), a spacing
  scale that follows the theme, and effect parameters, beside the existing
  colour, radius, type and elevation tokens.
- `theme.Theme` names a complete design language as one typed value: both
  palettes, the radius base, type roles and families, the spacing unit,
  elevation, motion and effects. `theme.Default()` returns the built-in
  theme, `theme.Use` publishes one and `theme.Current` returns the published
  one.
- `Preferences`, `Configure`, `UpdatePreferences`, `PreferencesSnapshot` and
  `PalettePatch` keep their signatures and behavior. Preferences resolve
  over the published theme. No breaking change.
- Component files in the root package and every family read the tokens
  instead of literals wherever a token exists. `internal/motion` keeps only
  its two helpers, which read the tokens; its three duration constants go.
- The default theme reproduces the current rendering and timing. Findings
  F1–F17 of `extract-design-language` stay open: a component that disagrees
  keeps its own named constant.
- The gallery's token cards show motion, spacing and effect tokens with
  their live values.
- Records follow: `DESIGN.md` names the token beside each rule,
  `docs/theme.md`, `docs/design.md`, `PRODUCT.md`, `.impeccable/design.json`,
  `CHANGELOG.md` and the `mygoui-theme` and `mygoui-motion` skills.

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `semantic-theme`: adds requirements for motion, spacing and effect tokens,
  for a theme supplied as one value, and for preferences resolving over it.

## Impact

- Public API added in `theme`: `Theme`, `Typography`, `Spacing`,
  `Elevation`, `Shadows`, `MotionTokens`, `EffectTokens`, `Default`, `Use`,
  `Current`, `Motion`, `Effects`. `theme.Space(n)` and the type role
  functions (`XS` … `HelpText`) read the published theme.
- Library code: every package that animates or spaces content (root, `form`,
  `navigation`, `overlay`, `list`, `data`, `date`, `color`, `chart`,
  `content`, `layout`, `internal/*`).
- No dependency change. `theme` still imports no family.
- Process-wide state: one more published value. Window-local widget state is
  unchanged.
