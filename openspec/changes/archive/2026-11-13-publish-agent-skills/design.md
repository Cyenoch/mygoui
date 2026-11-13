# Design

## Context

See `proposal.md` for the motivation and `complete-native-ui-solution`
(D2, D3) for the scope. Constraints found in this checkout on 2026-11-13:

- `.agents/skills/` holds 26 skill directories, all maintainer-facing:
  twelve `openspec-*` written by `openspec update --force`, `impeccable`
  written by `impeccable install`, `golang-*`, `github-actions-*`,
  `create-github-action-workflow-specification`, `mygo-maintenance`,
  `ponytail` and `writing-for-agents`. `.claude/skills` links to it;
  `scripts/setup-agents.mjs` verifies three directory links only.
- `skills` CLI 1.8.1 (`npx skills`), read from its `dist/cli.mjs`
  `discoverSkills` and `docs/creating-skills.md`:
  - It searches the source root, `skills/`, `skills/.curated|.experimental|.system`
    and about 30 agent directories including `.agents/skills/` and
    `.claude/skills/`, and unions the results; the first name wins.
  - It walks with `readdir` entries that are directories, so a symlinked
    skill directory is not discovered.
  - It hides a skill whose `metadata.internal` is true, and a skill under an
    agent directory whose name is in the source's `skills-lock.json`.
    That is why `npx skills add <checkout> --list` reported 22 skills
    before this change: the four lock-listed skills were hidden.
  - A source with a subpath (`.../tree/<ref>/<path>`, or a local directory)
    starts the search there.
- Generated upstream refreshes remove only their own names: OpenSpec's
  `released-output-sweep.js` targets its workflow directories; the
  Impeccable npm launcher (`cli/bin/cli.js`) contains no removal call. The
  Impeccable engine binary was not inspected.
- An installed skill lives in the application's project, away from this
  repository; relative links to `DESIGN.md` or `docs/` would not resolve
  there. The module zip carries those records, so
  `go list -m -f '{{.Dir}}' github.com/Cyenoch/mygoui` locates the version
  the application requires, and resolves to the checkout inside it.

## Goals / Non-Goals

**Goals:**

- Five skills an application author installs with one documented command,
  and the same files loaded by this project's harnesses.
- One source per topic; maintainer material disclosed behind one pointer.
- Every API name and Go example checked against the pinned versions.

**Non-Goals:**

- Hiding the maintainer and generated skills from a whole-repository
  `npx skills add Cyenoch/mygoui` listing (see Risks).
- New rules. Skills describe the recorded system; rules found unrecorded
  are reported for `extract-design-language` and
  `component-behavior-reference`, not invented here.
- An automated check of skill API names or snippets.

## Decisions

### D1. Published skills live in root `skills/`, linked into `.agents/skills/`

`skills/mygoui-<topic>/` holds the real files. Each
`.agents/skills/mygoui-<topic>` is a relative symlink to
`../../skills/mygoui-<topic>`, so Claude (through `.claude/skills`) and Codex
load the same files.

Evidence, run against this checkout after the files were added:

- `npx skills add <checkout> --list` lists the five `mygoui-` skills (from
  `skills/`) beside the 22 skills already discoverable.
- `npx skills add <checkout>/skills --list` lists exactly the five; the
  GitHub equivalent is `npx skills add https://github.com/Cyenoch/mygoui/tree/main/skills`.
- The links are not discovered a second time, since the CLI skips symlinks.
- The CLI parses frontmatter as strict YAML and skips a skill whose plain
  `description` contains `: `; the descriptions are quoted.

Alternatives considered:

- Real directories in `.agents/skills/mygoui-*`. Discoverable with
  `npx skills add Cyenoch/mygoui`, but no source path then offers only the
  mygoui skills: the subpath `.agents/skills` lists all 27, and the authored
  files sit among generated ones that `make agents-update` rewrites.
- Real directories in `.agents/skills/` with links from `skills/`. The CLI
  ignores the links, so `skills/` would publish nothing.
- Marking maintainer skills `metadata.internal` or adding them to
  `skills-lock.json` to hide them. The first edits generated upstream skills;
  the second misstates what the lockfile installed.

The `mygoui-` prefix and the separate root keep the published skills
distinct from generated ones. OpenSpec's refresh writes neither path; the
Impeccable engine's refresh is unverified (Risks).

### D2. One skill per topic, five model-invoked skills

Names: `mygoui-components`, `mygoui-theme`, `mygoui-design` (design and
behavior), `mygoui-motion`, `mygoui-effects` (shader effects). Each topic
has its own leading word that should trigger it alone, so each keeps a
description; an application author's agent must reach them without being
told. Descriptions front-load the topic and name only the distinct branches
(choose/compose, customize, review/verify, animate, shader/fade).

Alternative considered: one `mygoui` router skill with five reference
files. One description is cheaper, but the router can only hint, and a
motion or shader question would load the whole router first.

### D3. Two audiences through `MAINTAINING.md`

`SKILL.md` states the public contract only. Its single pointer reads "for
work on the mygoui repository itself, read MAINTAINING.md", so an
application agent pays one line and never opens the file. `MAINTAINING.md`
holds `internal/motion`, `internal/busy`, the OpenSpec and Impeccable
workflow, which record owns each rule, and how to keep the skill current.
It names repository paths in backticks rather than relative links, because
the file is read both from `skills/` and through the `.agents/skills` link.

### D4. Records and APIs are looked up, not cached

Skills give `go doc <package> <symbol>` for option fields and the module
directory lookup for records. They keep what no lookup yields: when to use
which component or mode, gotchas (`Preferences.Accent` replaces Primary,
`Revision` after in-place edits, `NonMotion`), and the state-motion
durations an author needs to match the library, attributed to `DESIGN.md`.

### D5. No `agent-workflow` delta

Its requirements still hold: harness skill paths resolve to `.agents/skills`,
setup creates and verifies the same three links, and refresh keeps them.
The published-skill entries are covered by `agent-skills`.

## Risks / Trade-offs

- [`npx skills add Cyenoch/mygoui` lists maintainer and generated skills
  too] → README documents the `tree/main/skills` source as the install
  command and `--skill mygoui-...` for a selection; the full listing still
  works and names the mygoui skills.
- [Per-skill symlinks need symlink support] → The project already requires
  it for `.claude/skills`; `docs/agent-workflow.md` states the Windows
  condition. The published copy has no links.
- [Skills drift from code] → Spec requirement: a change that alters a stated
  API or rule updates the skill; `MAINTAINING.md` says how. Names and
  examples were checked once by `go doc` and a scratch module; no automated
  gate exists.
- [Rules move during `extract-design-language`] → Skills point at
  `DESIGN.md`; only the motion durations are restated, attributed, for
  authors who must match them; that change updates them.
- [Harness loading of per-skill links unverified for Codex] → Claude Code
  loading is checked once in this change; Codex is recorded as unverified
  until observed.
- [Impeccable refresh] → `make agents-update` runs the Impeccable engine,
  which was not inspected or run here. The agent setup tests assert the
  links resolve after setup and the OpenSpec refresh only; check them after
  the next real refresh.
- [Version skew] → v0.2.0 ships an older `DESIGN.md`. Skills cite a record
  and its topic rather than a heading, so the pointer holds in either version.
