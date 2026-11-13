# Maintaining the mygoui theme

For work inside the mygoui repository. Paths are relative to its root.

## Ownership

- `theme/` holds the executable values: `language.go` (`Theme`, its token types, `Default`, `Use`), `palette.go` (roles, both appearances, compositing), `tokens.go` (spacing, type, radius, heights, shadows), `theme.go` (preferences resolved over the published theme, `Apply`).
- `docs/theme.md` owns role meaning, customization and token provenance, including the documented/runtime dark-mix discrepancy. `DESIGN.md` owns the visual rules and records default values in its frontmatter; refresh it when a default changes.
- `docs/references.md` owns the licensing boundary for reference values. Read it before copying any value from a reference.
- `openspec/specs/semantic-theme/spec.md` holds the theme requirements; a token or preference change is an OpenSpec change (see `docs/agent-workflow.md`).

## Rules for library code

- Library paint reads `theme.Colors(c)` each build; it never reads the context theme or caches a palette per window.
- Spacing is `theme.Space(n)`; durations, easings and effect parameters come from `theme.Motion()` and `theme.Effects()`, read in the build or paint that uses them, never in a package-level variable. Other metrics come from `theme` constants and `RadiusScale()`; scale text, line boxes and icons with `theme.Scale`.
- A new role gets a `Palette` field, a `PalettePatch` pointer field, values in both appearances, a row in `docs/theme.md`, and `make generate` for the gallery API table.
- A value becomes a token when several components share it. A value one component uses, or one that disagrees with its neighbours, stays a named constant in that component's file and is listed in `docs/theme.md`. A new token gets a field in `language.go`, its default in `Default()`, a row in `docs/theme.md`, a card in `internal/gallery/theme_tokens.go` and `make generate`.
- `Default()` must reproduce the look `DESIGN.md` records; `theme/language_test.go` pins its values.
- `Use` is the trust boundary for a theme: a value a component divides by, or that MyGo panics on, is rejected in `Theme.validate` with the field's name, and a zero motion duration is resolved to the shortest MyGo animates. A new duration token joins both.

## Keep this skill current

When a change adds, renames or removes a public theme API or token, changes `Configure`/`UpdatePreferences`/`Use` semantics or the precedence of `Palette` over `Accent`, update `SKILL.md` in the same change and list it in that change's tasks. Keep values in `DESIGN.md` and `docs/theme.md`; `SKILL.md` points at them. Compile changed examples against the checkout.
