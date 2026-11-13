# Tasks

## 1. Layout

- [x] 1.1 Create `skills/mygoui-{components,theme,design,motion,effects}/` and relative links `.agents/skills/mygoui-* -> ../../skills/mygoui-*`; verify each link resolves with `ls -L` and through `.claude/skills`

## 2. Skills

Author each with `.agents/skills/writing-for-agents/SKILL.md` and
`SKILL-MECHANICS.md`. Each task is done when `SKILL.md` names only public
APIs, `MAINTAINING.md` holds the maintainer material behind one pointer,
every named API is found by `go doc` or grep at the pinned versions, and
every Go example compiles in a scratch module that replaces the module with
this checkout.

- [x] 2.1 `mygoui-components`: families, typed options, caller-owned state, body callbacks, `c.Key` identity, `Loading` flags, `Revision`, ownership rules, common mistakes
- [x] 2.2 `mygoui-theme`: `theme.Apply`, `Colors`, `PalettePatch` per appearance, `Configure`/`UpdatePreferences`, radius, type roles, `Scale`, font registration
- [x] 2.3 `mygoui-design`: hierarchy, density, state feedback, loading, focus and keyboard, accessible names, UI copy, verification with `ui.NewTester`; its example test passes
- [x] 2.4 `mygoui-motion`: what moves and for how long, what stays instant, reduced motion, nothing at rest, `ui.Animation`, `AnimateTo`/`AnimateColor`, `ElementTransition`; its example runs under a Tester and settles
- [x] 2.5 `mygoui-effects`: `ScrollFade` `MaskFade`, `SkeletonGroup` shimmer, `Element.Effect` with `ui/shader`, the 8192-pixel layer limit, webview exclusion, CPU fallback, what headless tests do not prove; its shader compiles at run time

## 3. Records

- [x] 3.1 `README.md`: a short section with the install commands; verify the commands match the discovery evidence in design D1
- [x] 3.2 `docs/agent-workflow.md`: the published-skill layout and how to keep skills current; `AGENTS.md` Skills section: one pointer to the published skills
- [x] 3.3 `PRODUCT.md`: the Agent skills row states the delivered skills and what remains planned

## 4. Verification

- [x] 4.1 `npx skills add <checkout> --list` lists the five `mygoui-` skills, and `npx skills add <checkout>/skills --list` lists only them
- [x] 4.2 `npm test` copies `skills/` into its fixture and asserts the five links resolve through `.claude/skills` after setup and after `openspec update --force`; a fixture without `skills/` fails it. The Impeccable refresh (`make agents-update`) has not been run and is unverified
- [x] 4.3 Claude Code in this checkout lists the `mygoui-` skills. Codex loading is unverified
- [x] 4.4 `npm run --silent openspec -- validate publish-agent-skills --strict --no-interactive` and `make spec-check` pass

## Workflow follow-up

- Verify against the spec scenarios, then archive with `openspec-archive-change`.
- `extract-design-language`, `open-theme-tokens`, `component-behavior-reference`
  and `native-effects` list the affected skill update in their own tasks.
