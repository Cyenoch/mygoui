# Maintaining mygoui motion

For work inside the mygoui repository. Paths are relative to its root.

## Internal contract

- `theme.Motion()` and `theme.Effects()` are the single source of shared timing (`theme/language.go`; `docs/theme.md` lists the tokens). A component reads them in the build or paint that animates, never in a package-level variable, and passes `Ease` explicitly: MyGo's default curve ignores the theme.
- `internal/motion` keeps two helpers over `Motion.State` and `Motion.Ease`: `Progress` (a `NonMotion` 0–1 value, for opacity, a color mix or a mark drawing in) and `Color`.
- `internal/busy` holds the loading timings (crossfade, the region's wait, the dim, the top line); see the `mygoui-components` skill's MAINTAINING.md.
- The root `motion.go` holds the root package's shared transitions: `fadeTransition` (fade only, so it also runs under reduced motion; it names `Shadows` so position, size and colors stay with the view) and `stateIcon` (a glyph keyed by itself, crossfading on change).
- A timing one component uses, or one that disagrees with its neighbours, is a named constant in that component's file; `docs/theme.md` lists them. Do not fold one into a token without a change that settles the disagreement in `DESIGN.md`.
- MyGo's `ui.CollapsibleBase` (accordion, disclosure, settings expander, timeline, about dialog) animates at the toolkit's fixed 200 ms ease-out and does not follow `Motion.Reveal`.

## Records

- `DESIGN.md` owns the rules in its state motion and loading sections; `.impeccable/design.json` records the native motion extensions. A changed rule goes there first, through an OpenSpec change (`docs/agent-workflow.md`), with Impeccable's `animate` playbook for design work.
- `openspec/specs/native-library/spec.md` holds the motion requirements.
- Tests pin the rules with `ui.NewTester`: an intermediate frame during a state change, none on a first-frame value, `FramePending` false at rest, and reduced motion. `motion_state_test.go` and `actions_motion_test.go` show the patterns; `theme_tokens_test.go` shows a test that publishes a theme and restores `theme.Default()` in `t.Cleanup`.

## Keep this skill current

`SKILL.md` names the token for each kind of change; `docs/theme.md` and `DESIGN.md` own the values. When `DESIGN.md` changes which token a change uses, an instant case or the reduced-motion rule, or when a motion token is added, renamed or removed, update `SKILL.md` in the same change and list it in that change's tasks. Recompile the example against the checkout.
