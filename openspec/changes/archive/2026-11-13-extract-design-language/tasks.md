## 1. Extraction

- [x] 1.1 Inventory colour roles and derived colours, type, spacing, radius, elevation, control sizing, focus and disabled states from `theme/` and the shared styling packages; verify each value in `DESIGN.md` against its Go source
- [x] 1.2 Inventory every motion duration, easing, loading and shader parameter in non-test library Go files; verify with a search for `time.Millisecond`, `time.Second`, `motion.` and `Ease` that each hit appears in the `design.md` token inventory
- [x] 1.3 Record components that disagree as findings in `design.md` without changing Go; verify `git status` shows no `.go` file changed by this work

## 2. Design record

- [x] 2.1 Rewrite root `DESIGN.md` in the eight canonical sections with valid frontmatter, standalone rules and a provenance pointer; verify `make design-doctor` reports no `DESIGN.md` schema finding
- [x] 2.2 Refresh `.impeccable/design.json` from the new record; verify it parses as JSON with `schemaVersion` 2 and `make design-doctor` reports no sidecar finding

## 3. Provenance

- [x] 3.1 Demote coss to provenance in `docs/references.md`, `docs/theme.md`, `docs/design.md`, `docs/components.md` and `docs/mygo-vs-css.md`, keeping the MIT/AGPL boundary unchanged; verify a search for `coss` in these files finds only provenance, licensing and per-row upstream links
- [x] 3.2 Update authority wording in `PRODUCT.md` (Visual Authority, Capabilities, Native Gallery), `AGENTS.md` Design Work, `docs/agent-workflow.md` design paragraphs, `openspec/config.yaml` context, `README.md` references line and `.impeccable/surfaces/internal-gallery.md`, re-reading each shared file before a targeted edit; verify the same search and that other sections are untouched in `git diff`
- [x] 3.3 Edit the `native-charts` main spec Purpose to name the semantic theme; verify `make spec-check` passes

## 4. Verification

- [x] 4.1 Validate the change with `npm run --silent openspec -- validate extract-design-language --strict --no-interactive` and run `make spec-check`
- [x] 4.2 Run `CGO_ENABLED=0 go test ./theme/... -count=1` and `make test`; record results and any failure owned by concurrent work

## Workflow follow-up

- Archive after review; the delta specs then rename the five requirements in the main specs.
