# Design

## Context

See `proposal.md` for the motivation. The state each pillar starts from, as
read in this checkout on 2026-11-13:

| Pillar | Exists | Missing |
| --- | --- | --- |
| Theme tokens | `theme`: `Palette` and per-appearance `PalettePatch`, radius scale from one base, type roles, four shadow sets, `Space(n)` on a 4 DIP unit, text scale and reduced motion | Motion is private: `internal/motion` holds three durations and about 46 other duration literals sit in component files. No effect tokens. No single value that names a whole theme, so coss is the theme rather than one theme. `.impeccable/design.json` lists three motion entries by hand |
| Component library | 57 coss subjects, the Lobe and survey absorptions, charts, colour, dates (`docs/components.md`) | Gaps the matrix already names, such as `FloatingPanel` keyboard move and resize, `Drawer` snap points, `NumberField` scrub |
| Design and behavior reference | `DESIGN.md` (visual rules, state motion, loading), `docs/design.md` (ownership), `docs/mygo-vs-css.md`, one capability line per component | `DESIGN.md` states it "does not establish a new visual identity" and defers to coss. Keyboard, focus, selection and dismissal are described per component only in passing; there is no contract a reader can check a component against |
| Agent skills | Maintainer skills in `.agents/skills`: OpenSpec, Impeccable, `mygo-maintenance`, Go skills | Nothing for an application author's agent. Nothing specific to mygoui's theme, motion or shader effects for either audience |

Constraints that carry over unchanged: caller-owned state, native focus,
editing and accessibility, no I/O or registration on import, and evidence
kept separate for headless and real OS/GPU results.

## Goals / Non-Goals

**Goals:**

- One product statement that names the pillars and who owns each record.
- A fixed order of follow-up changes with boundaries that do not overlap.
- A stated end state for the coss relationship and the rule that holds
  until it is reached.

**Non-Goals:**

- Designing mygoui's own visual language. That needs its own design
  direction and is not scheduled here.
- Any API, token, skill or behavior work. Those belong to the follow-ups.
- Renaming coss-named tests, environment variables or the `coss` branch.

## Decisions

### D1. coss becomes provenance by extraction, not by rewrite

`extract-design-language` uses Impeccable's `extract` and `document`
playbooks to move every visual and motion rule the implementation follows
into `DESIGN.md`, stated in mygoui's own terms. After it, `DESIGN.md` is
the authority, `docs/references.md` records where each rule came from, and
the coss pin is consulted for provenance and licensing only. The five
requirements listed in the proposal get their deltas there.

Alternative considered: declare `DESIGN.md` the authority now. Rejected:
the file currently defers to coss for anything it does not record, so the
declaration would leave gaps with no authority at all.

### D2. Skills come first and describe the system as it is

The user chose skills as the first follow-up. They can be written from the
existing records and code, and they expose which rules are not written down
anywhere, which feeds the extraction and behavior changes. Skills that
describe design, motion and effects are revised as those changes land; each
of those changes lists the skill update among its tasks.

Alternative considered: extraction first, so skills never quote a rule that
later moves. Rejected: it delays the pillar the user asked for, and a
skill that points at `DESIGN.md` and `theme` survives the move.

### D3. One skill set, two audiences

Both audiences are served from one source per topic: components, theme,
design and behavior, motion, and shader effects. The application-author
part states the public contract; a maintainer section adds internal
contracts (`internal/motion`, `internal/busy`, the OpenSpec and Impeccable
workflow). `publish-agent-skills` decides the directory that is both
installable from the repository and loaded by this project's harnesses,
and how generated upstream skills stay separate from it.

Alternative considered: two parallel skill sets. Rejected: the shared
rules would be written twice and drift.

### D4. Tokens open in layers, with coss as the default theme

`open-theme-tokens` makes motion (durations, easings), spacing, elevation
and effect parameters public tokens beside the existing colour, radius and
type tokens, and lets a host supply a whole theme as one value. Component
files read tokens instead of literals. The coss values become the default
theme. This is what lets a later design language replace coss without
component edits.

### D5. Behavior is a checkable contract per component

`component-behavior-reference` records, per component, the keys, focus
movement, selection, dismissal, loading and motion a consumer can rely on,
in a form a `ui.NewTester` scenario can check. Its sources are the existing
native implementation, platform conventions and the WAI-ARIA authoring
patterns; headless behavior libraries such as UI SDK serve as a checklist
of what to cover, not as a dependency. Closing a gap it finds is a separate
change per component.

### D6. Effects become a public family last

Two private shaders exist today (the `ScrollFade` mask and the
`SkeletonGroup` shimmer). `native-effects` makes shader effects a public,
tokenised family with the GPU-less fallback and the 8192-pixel layer limit
stated as contract. It follows the token change because its parameters are
tokens.

### Follow-up changes, in order

| # | Change | Delivers | Depends on |
| --- | --- | --- | --- |
| 1 | `publish-agent-skills` | Skills for components, theme, design and behavior, motion, shader effects; both audiences | this change |
| 2 | `extract-design-language` | `DESIGN.md` as authority, coss as provenance, deltas for the five coss-named requirements | this change |
| 3 | `open-theme-tokens` | Public motion, spacing, elevation and effect tokens; a theme as one value | 2 |
| 4 | `component-behavior-reference` | Per-component behavior contract and its checks | 2 |
| 5 | `native-effects` | Public shader effect family | 3 |

Changes 3 and 4 are independent of each other.

## Risks / Trade-offs

- [Skills quote rules that move during extraction] → Skills link to the
  owning record instead of copying values; each later change updates the
  affected skill in its own tasks.
- [The product record promises pillars that do not exist yet] →
  `PRODUCT.md` states which pillars are planned and names the change that
  delivers each; capability claims stay in `docs/components.md`.
- [Decoupling from coss loses the fidelity checks] → The coss-fidelity
  tests keep running through extraction; they are re-described as checks
  of the recorded design language only when `DESIGN.md` covers what they
  assert.
- [Licensing boundary blurs once coss is "only provenance"] →
  `docs/references.md` keeps the MIT/AGPL boundary unchanged.

## Open Questions

- Where installable skills live and how `npx skills add` finds them is
  decided in `publish-agent-skills`.
- Whether the behavior contract is recorded as OpenSpec capabilities per
  family or as documents under `docs/` is decided in
  `component-behavior-reference`.
