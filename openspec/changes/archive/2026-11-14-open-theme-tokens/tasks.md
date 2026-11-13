# Tasks

## 1. Theme value and tokens

- [x] 1.1 Add `Theme`, `Typography`, `Spacing`, `Elevation`, `Shadows`, `MotionTokens`, `EffectTokens`, `Default`, `Use`, `Current`, `Motion` and `Effects` to `theme`; `go build ./theme` passes
- [x] 1.2 Resolve `Preferences` over the published theme without changing `Configure`, `UpdatePreferences`, `PreferencesSnapshot` or `PalettePatch`; the existing `theme` tests pass unedited
- [x] 1.3 `Space(n)`, the type role functions and the shadow functions read the published theme; reads stay one atomic load
- [x] 1.4 Tests in `theme`: default theme equals the previous defaults, a whole theme changes palette and geometry together, a partial override inherits the rest, supplied and returned values are copies, preferences resolve over a custom theme

## 2. Components read tokens

- [x] 2.1 `internal/motion`, `internal/busy`, `internal/anchored`, `internal/style`, `internal/formstyle` read motion, effect and spacing tokens; `go build ./...` passes
- [x] 2.2 Root package reads tokens; package-level transitions that named a duration become functions
- [x] 2.3 `form`, `list`, `content` and `data` read tokens
- [x] 2.4 `navigation` reads tokens; the tabs indicator and segmented thumb read the reveal token (F14 sources)
- [x] 2.5 `overlay` and `layout` read tokens
- [x] 2.6 `date`, `color` and `chart` read tokens
- [x] 2.7 Test with `ui.NewTester`: overriding a motion token changes a component's timing; reduced motion still settles an overlay at once under a longer duration
- [x] 2.8 List every remaining `time.Millisecond`, `time.Second` and easing literal in non-test library code outside `internal/gallery` and `examples`, each with its reason, in `design.md`

## 3. Gallery

- [x] 3.1 Token cards show motion, spacing and effect tokens with their live values; pages whose components read them list them; the appearance controls still work; a test pins that a card shows the published value (`go test ./internal/gallery`)
- [x] 3.2 `make generate` after the exported struct changes

## 4. Records

- [x] 4.1 `DESIGN.md` names the token beside each rule it governs, in targeted edits
- [x] 4.2 `docs/theme.md` documents the tokens and the one-value theme with an example compiled against the checkout
- [x] 4.3 `docs/design.md` package table and `internal/motion` paragraph; the "Theme tokens" row of `PRODUCT.md`; `.impeccable/design.json` motion purposes; `CHANGELOG.md`
- [x] 4.4 `skills/mygoui-theme` and `skills/mygoui-motion` (`SKILL.md` and `MAINTAINING.md`) follow `writing-for-agents`; every Go block compiles in a scratch module that replaces the module with the checkout

## 6. Review fixes

- [x] 6.1 `Configure` stores a font size or radius equal to the theme's as unset; test that `Configure(PreferencesSnapshot())` followed by `Use` follows the new theme
- [x] 6.2 `Use` returns an error naming the field for a negative duration or a period under a millisecond; a zero motion duration settles at once; tests for both
- [x] 6.3 A theme's palette is used as written; every status solid hovers at `Effects.HoverOpacity`; test
- [x] 6.4 A narrow `SweepBand` leaves the track; the default output is unchanged
- [x] 6.5 Tabs, segmented control and breadcrumbs read `theme.Space`; field vertical padding and the menu indicator inset are listed as local values
- [x] 6.6 The default values are pinned in a table and the fallback curves are asserted; skills and `docs/theme.md` aligned

## 5. Verification

- [x] 5.1 `make generate` and `make check` pass with no test expectation edited
- [x] 5.2 `make race` passes
- [x] 5.3 `golangci-lint run ./...` reports no new finding
- [x] 5.4 `GOOS=js GOARCH=wasm CGO_ENABLED=0 go build ./...` passes
- [x] 5.5 `make spec-check` passes; `make design-doctor` findings read: one `mention`, `config-build-path-unset`, unrelated to this change
- [x] 5.6 Record OS, GPU and screen-reader evidence as pending in `VALIDATION.md`
