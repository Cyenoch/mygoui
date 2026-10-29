# Validation

## Bklit native charts and EmptyState alignment

2026-10-29, uncommitted native-library/Gallery changes, macOS 27.0.1
(26A434)/arm64, Go 1.27.3, MyGo 0.7.2. Bklit source is pinned at
`ec5a9705ee973f95fd0ef2d974e9d00ae3fda011`; coss remains the visual authority.
The chart and utility matrix in `docs/charts.md` states native subsets and gaps,
not full React, decorative-preset, Studio or export parity.

- Final `make check`: pass, including all headless tests, Gallery public snippet
  type checks, ordinary build, Go vet, MyGo lifetime vet, formatting and module
  integrity; `/tmp/mygoui-charts-reviewed-check.log`.
- Full `make race`: pass, no reported race; Gallery 169.945 s, overlay 248.975 s;
  `/tmp/mygoui-charts-final-race.log`. This run precedes the small Enter-action
  follow-up; the subsequent `CGO_ENABLED=1 go test -race -p 1 ./chart
  ./internal/gallery -count=1 -timeout=600s` also passes (Gallery 230.809 s),
  `/tmp/mygoui-charts-reviewed-race.log`.
- Final `golangci-lint run ./...`: zero issues;
  `/tmp/mygoui-charts-reviewed-lint.log`. Final `make wasm`: pass;
  `/tmp/mygoui-charts-reviewed-wasm.log`. No runtime module dependency is added.
- Native Tester covers keyboard/readout, pointer and disabled ancestors including
  trailing modifiers, activation counts, legend filtering, brush domain, live
  sample replacement without mutation, signed stacks, malformed numeric/geometry/
  graph/hierarchy input, huge finite angles, aggregate overflow, bounded weekend
  shading, long forward flows, theme overrides and custom plot/indicator context.
  All 35 chart/utility Gallery pages render and their Go snippets compile.
- `MYGOUI_CHART_ARTIFACT_DIR=... CGO_ENABLED=0 go test ./internal/gallery -run
  '^TestChartVisualArtifacts$' -count=1 -timeout=60s`: pass. Thirteen representative
  subjects have wide/compact, light/dark captures under `artifacts/charts/`;
  compact dark uses text scale 1.3. Headless contact sheets were inspected.
- EmptyState action alignment is checked at widths 280/640 and text scales 1/1.3,
  with exactly one callback per click. Its headless contact sheet is
  `artifacts/empty-state-fix/emptystate-action-contact.png`.
- Final `e2e` Gallery binary runs in a real macOS native window. Automation reads
  the native element tree, focuses Line Chart and presses End (value: Day 6,
  Sales 80), drills Sunburst with Home/Enter (Product breadcrumb), and activates
  Create search (one created search). Transcript:
  `/tmp/mygoui-charts-native/reviewed-verification.json`. Three inspected captures
  are `artifacts/charts/native-macos-{line-inspected,sunburst-drilled,empty-state}.png`.
  The EmptyState capture shows the button, icon, title and footer on one centerline.
- `make design-doctor`: only the existing informational
  `config-build-path-unset`; `/tmp/mygoui-charts-doctor.log`. It concerns future
  Persuade/Experience surfaces, not this native Operate workbench. No browser
  detector result is counted as Go UI validation.
- Strict `make spec-check` and `git diff --check`: pass;
  `/tmp/mygoui-charts-reviewed-spec.log`. Semantic verification maps all five
  requirements and six scenarios to source and behavior evidence. The new
  `native-charts` main spec preserves the delta Purpose and requirement blocks.
- `add-bklit-charts` is synchronized and archived as
  `2026-10-29-add-bklit-charts`, with all ten tasks complete. Post-archive
  `make spec-check` passes five main specs and both completed archives;
  `/tmp/mygoui-charts-archived-spec.log`. No active change remains.

Independent source, code and finish-design reviews were performed. Code findings
include overflow/arc hangs, calendar bounds, readout/axis displacement, ancestor
disabled routing, missing stacked gaps, angular selection and long-flow geometry;
all were resolved. Six material finish-design findings were fixed and recaptured;
the reviewer's follow-up accepts those six fixes, not all possible states.
The final Enter/disabled finding was fixed locally with a focused regression.
Read-only design documentation review found no durable system drift; DESIGN.md
and its sidecar remain unchanged. The archived change retains the detailed
requirement/scenario map and review resolutions.

Actual Windows/Linux execution, other GPU backends, IME and physical screen-reader
operation are not verified. macOS element-tree metadata is not a VoiceOver test.
No branch, commit, push or release is performed by this work.

## Impeccable and OpenSpec adoption

2026-10-29, macOS/arm64, adoption at `6dd9502` plus local review fixes and
Claude/Codex integration. These changes configure development workflows;
the existing Go/MyGo implementation remains outside their ownership.

- `make agents` and `make agents-update`: pass, including lockfile install,
  all twelve OpenSpec workflows, both hook manifests and verified engine
  0.1.18. The skill and agent links survive both commands. Logs:
  `/tmp/mygoui-agent-dual-setup.log` and
  `/tmp/mygoui-agent-dual-update.log`.
- `npm test`: four regressions pass after refresh. A global core profile
  cannot prune the project's workflows; both launchers bypass a stale home
  engine; conflicting Claude content remains untouched; a skill/engine
  release mismatch fails before an engine is installed. Log:
  `/tmp/mygoui-agent-dual-tests.log`.
- Claude and Codex each have four official agent profiles. Directory
  realpaths resolve to their canonical `.agents` targets; Codex TOML
  required fields and Claude frontmatter were inspected. `CLAUDE.md`
  imports the same AGENTS.md used by Codex.
- Both configured hook commands were exercised with synthetic
  `PostToolUse` and `Stop` payloads in a temporary CSS fixture. All four
  invocations exit successfully with valid JSON; Stop reports the intended
  fixture finding. This verifies command/protocol behavior, not harness
  trust or automatic events. Log: `/tmp/mygoui-agent-dual-hook-smoke.log`.
- `GOWORK=off make check`: pass after review fixes, covering the headless
  suite, ordinary build, vet/lifetime vet, formatting and module integrity.
  Log: `/tmp/mygoui-agent-dual-go-check.log`. Earlier adoption checks also
  passed `GOWORK=off make wasm workflows` and `make design-artifacts`;
  `/tmp/mygoui-agent-builds.log` and
  `/tmp/mygoui-agent-design-artifacts.log`. The button contact sheet under
  `artifacts/design-review/` was inspected as collector evidence.
- Strict current-spec/active-delta validation and `actionlint`: pass.
  Impeccable resolves PRODUCT.md, DESIGN.md and the gallery Operate brief.
  Doctor has one informational `config-build-path-unset` finding: a future
  Persuade/Experience build preference has not been chosen. No preference
  was invented. Logs: `/tmp/mygoui-agent-dual-specs.log`,
  `/tmp/mygoui-agent-dual-context.log`,
  `/tmp/mygoui-agent-dual-doctor.log`.
- `adopt-impeccable-openspec` was synchronized and archived as
  `2026-10-29-adopt-impeccable-openspec`, creating the six-requirement
  agent-workflow main spec. Final `make spec-check` passes all four main
  specs and the archive's fourteen completed tasks; no active changes
  remain. Logs: `/tmp/mygoui-agent-dual-archive.log` and
  `/tmp/mygoui-agent-dual-specs-final.log`.

Independent read-only review identified two reproducibility defects:
the npm OpenSpec entry used the user's global profile, and the design
launcher could select a stale global engine. One project-scoped CLI entry
and a checksummed sibling engine fix them; focused tests reproduce both
failure conditions. The shared-directory refresh was also exercised:
upstream Claude updates are followed by the canonical Codex variant, so
both invocation forms remain available without changing upstream skills.

The web detector and browser Live Mode do not validate native Go controls.
Actual Claude/Codex session loading, hook trust/automatic firing,
Windows/Linux execution, GPU and screen-reader operation were not exercised
by this tooling follow-up. No new race-sensitive Go behavior was introduced.
No commit, push or release was performed by this work.

## MyGo 0.7.2 integration

2026-10-29, uncommitted changes based on `674e53d`, macOS 27.0.1/arm64,
Go 1.27.3 with `GOWORK=off`. The library and module-pinned CLI resolve
MyGo v0.7.2. This section supersedes earlier current-version statements,
not their historical measurements.

- `make check`: pass; `/tmp/mygoui-072-check-final.log`. Full headless suite,
  ordinary build, Go vet, MyGo lifetime vet, formatting and module integrity.
- `make race`: pass, no reported races; `/tmp/mygoui-072-race.log`.
  Gallery 165.004 s, overlay 248.209 s; serial package execution and the
  existing 600-second package timeout remain unchanged.
- `golangci-lint run ./...`: zero issues;
  `/tmp/mygoui-072-lint-final.log`. `make wasm workflows`: pass;
  `/tmp/mygoui-072-builds-final.log`.
- Windows/amd64 and Linux/amd64 ordinary cross-builds, plus an `e2e` WASM
  cross-build: pass; `/tmp/mygoui-072-cross-builds.log`. Native execution on
  those platforms is not established by cross-compilation.
- `go test -tags e2e ./examples/gallery`, normal `e2e` native build and
  `e2e mygo_noinspector` build: pass. Dependency listing confirms normal
  Gallery builds exclude `plugins/automation`, while `e2e` includes it.
- `govulncheck@v1.9.0 ./...`: no vulnerabilities found;
  `/tmp/mygoui-072-vuln.log`. `go mod tidy -diff` is empty; all 18 modules
  have no replacements. Retained MyGo license matches v0.7.2 verbatim.
- Pinned `go tool mygo build -skip-dmg`: native compilation and ad-hoc
  signing pass; `/tmp/mygoui-072-app-final.log`. Output:
  `artifacts/coss/mygo-072/app/darwin-arm64/Coss UI Gallery.app`.
- `go tool mygo migrate-ui .` preview only reports the existing native
  callback element key notices in `overlay/popover.go` and `overlay/tooltip.go`;
  those callback elements cannot be constructed by this layer. Lifetime
  vet passes. No broad rewrite or local dependency replacement is applied.

New headless tests exercise drawn-menu keyboard selection, disabled-item
skipping, callback counts, focus restoration, submenu pointer activation and
Escape dismissal. Existing OS-menu tests remain in place. Avatar tests cover
loaded PNG/SVG paint via `DecodeImage`, plus absent, typed-nil and empty
images retaining initials. Context-menu tests also pass three shuffled runs.

Real macOS verification uses the `e2e` binary inside a separate, ad-hoc-signed
copy of the built `.app`, launched normally so the window is focused. The
private automation socket is mode `0600`. Through the 0.7.2 protocol:
component search and navigation, Shift+F10 opening, radio roles/checked state,
selection persistence, Enter/Escape focus restoration, secondary click and
outside dismissal pass. Light/dark window captures are inspected at
`artifacts/coss/mygo-072/contextmenu-<light|dark>-open.png` (2360 x 1720).
The interaction transcript is `/tmp/mygoui-072-bundle-automation.json`.
The test app quits and removes its socket. A bare, unfocused CLI launch did
not expose an accessibility focus ID; it is not used as focus evidence.

Independent read-only review found a Gallery state mismatch: its themed
menu used anchor-local selection but displayed window-owned result text.
The example now uses one window-owned selection for both. A red-first
regression confirms the selected radio mark survives Code/Preview remounts;
the final macOS automation transcript also checks its accessible radio role
and checked state after that round trip. The review's suggested radio semantic
check is covered by this real-window transcript; headless Tester does not
expose its accessibility tree publicly. The full Gallery race suite is rerun
after this fix and passes in 156.512 s;
`/tmp/mygoui-072-race-reviewed.log`.

Actual screen-reader operation, Windows/Linux desktop behavior, remote CI
and the previous arbitrary nested rounded-clip GPU limitation are not
revalidated. Existing clipping workarounds remain. No commit, push or release
is performed.

## Main CI dependency download follow-up

2026-10-28, based on main `26e320f`, with the previously validated race
budget fix cherry-picked as `3544386`. Main run `40530517066` passed all
three OS test jobs and analysis. Its race job never reached tests: the
61.2 MB `fonts-noto-cjk` download stalled at the Azure Ubuntu mirror until
the 15-minute job deadline. Check annotations confirm the job timeout;
`/tmp/mygoui-main-40530517066-race.log`.

Both Linux dependency steps now remove Azure from the hosted runner's
mirror list, retaining its official Ubuntu archive/security alternatives.
APT uses three retries and 30-second HTTP/HTTPS inactivity timeouts; each
installation step has a five-minute overall limit. Failure still stops
the job. The race package/job budgets are 600 seconds and 25 minutes.
No production code, dependencies or tests change.

- Exact workflow shell extracted for both jobs: identical; Bash syntax,
  `actionlint` and `git diff --check` pass.
- Ubuntu 24.04/amd64 Docker with the runner's mirror-list/source layout:
  the exact installation shell passes. Both official mirrors remain, APT
  reads all three settings, and all required Pango/CJK/emoji packages are
  installed. The CJK download uses `https://archive.ubuntu.com/ubuntu`;
  `/tmp/mygoui-main-ci-apt-install.log`. The minimal image needed a mounted
  CA trust bundle for this check; TLS verification was not disabled.
- `GOWORK=off make check`: pass; `/tmp/mygoui-main-ci-check.log`.
- `GOWORK=off golangci-lint run ./...`: zero issues;
  `/tmp/mygoui-main-ci-lint.log`.
- `GOWORK=off make wasm workflows`: pass;
  `/tmp/mygoui-main-ci-builds.log`.
- Independent read-only review found no actionable issue; its actionlint,
  Bash syntax, ShellCheck and mirror-filter fixture checks pass.

The earlier full Linux race replay below covers the unchanged production
and test source; it is not rerun for the apt-only change. Hosted CI on the
new main revision remains pending. No native desktop verification is claimed.

## Race timeout budget follow-up

2026-10-27, uncommitted follow-up to `cfb0515`. CI run `40529194388`
passed analysis and all three OS test jobs. Its serial race job still
timed out in overlay at 300.026 s; Gallery passed in 283.182 s. The active
AlertDialog test was only one second old when the package timer expired,
and raster workers were evaluating effect coverage. No data-race report
appears in the failed log, `/tmp/mygoui-ci-cfb0515-race.log`.

The race package timeout is now 600 seconds in both CI and `make race`.
The CI race job allows 25 minutes for serial packages, setup and compilation.
Race instrumentation, package-local concurrency, CI shuffle and all tests
remain enabled. No production code or test assertions change.

- The active `TestAlertDialogOptionalIconAndClose` passes alone with race
  detection and the failed shuffle seed `1793114605814973555` in 2.793 s
  on macOS/arm64; `/tmp/mygoui-ci-alert-race.log`.
- `GOWORK=off make check`: pass; `/tmp/mygoui-ci-budget-check.log`.
- `GOWORK=off golangci-lint run ./...`: zero issues;
  `/tmp/mygoui-ci-budget-lint.log`.
- `GOWORK=off make wasm workflows` and `git diff --check`: pass;
  `/tmp/mygoui-ci-budget-builds.log`.
- Full Linux/arm64 Docker race suite capped at two CPUs: pass with
  `GOWORK=off CGO_ENABLED=1 GOMAXPROCS=2 go test -race -p 1 -shuffle=1793114605814973555 ./... -count=1 -timeout=600s`.
  All 17 packages complete with no reported races; Gallery takes 164.871 s
  and overlay 218.606 s. First package start to final result is 441.429 s;
  `/tmp/mygoui-ci-race-linux.jsonl`. This local run also fits the old package
  timeout; it does not reproduce hosted Linux/amd64 runner speed.
- Independent read-only review found no blocking issue in the timeout
  alignment or job budget; its workflow validation also passes. It did not
  rerun the test suite.

No commit, push, remote CI rerun or native desktop verification is claimed.
Hosted-runner timing remains unverified for this follow-up.

## PR #3 race scheduling follow-up

2026-10-27, follow-up to `4502b4e`. Run `40523761606` passed analysis and
all three OS test jobs, including Windows job `122257584865`. Race job
`122257584206` hit the 300-second package timeout in Gallery and overlay;
its log contains no data-race report. The stacks show CPU shader evaluation
and backdrop blur, not a deadlock.

Race test binaries now run serially (`-p 1`) in CI and `make race`, avoiding
competition between internally parallel software rasterizers. Package-local
goroutines, race instrumentation, CI shuffle and the 300-second package
timeout remain enabled. The 15-minute job limit is unchanged. Gallery's
bulk tests no longer force an extra frame after `NewTester`, `Click` or
`SetPreferences`, which already settle frames in pinned MyGo v0.7.1.
Explicit frames after direct state mutations and all existing assertions
remain. Production component behavior is unchanged.

- In Linux Docker capped at two CPUs, old overlay code passes its failing
  shuffle seed `1793112409978699930` in 230.577 s in isolation. CPU profiles
  identify glass blur and software effect rasterization as hot paths.
- In the same two-CPU-capped setup, Gallery's full race suite with seed
  `1793112401198632472` passes before and after removing duplicate frames:
  284.118 s and 178.718 s respectively. These are shared-host observations,
  not a controlled performance benchmark. Logs:
  `/tmp/mygoui-pr3-race-gallery-full-baseline.jsonl` and
  `/tmp/mygoui-pr3-race-gallery-no-duplicate-frames.jsonl`.
- `make check`: pass; `/tmp/mygoui-pr3-race-check.log`.
- Full Linux Docker race suite capped at two CPUs: pass using
  `CGO_ENABLED=1 go test -race -p 1 -shuffle=1793112401198632472 ./... -count=1 -timeout=300s`.
  Gallery 177.329 s, overlay 228.003 s; from the first package start to the
  final result is about 7 min 42 s. Log:
  `/tmp/mygoui-pr3-race-serial-linux.jsonl`.
- Full lint: zero issues; `/tmp/mygoui-pr3-race-lint.log`.
- `make workflows`, `make wasm` and Windows/amd64 Gallery test-binary
  cross-compilation: pass.
- Independent read-only review found no blocking issue in frame removal or
  race scheduling; it did not run tests. It noted that serial execution
  changes scheduling and must still fit the overall job limit.

Remote CI on this follow-up remains pending; the local container run does not
establish hosted-runner wall time or native desktop behavior.

## PR #3 Windows marker follow-up

2026-10-27, follow-up to `215767b`. Windows job `122248524857` in
run `40521219792` failed overlay motion tests: the detected residual bounds
matched the Open label rather than the solid motion marker. The old
equal-green/blue color filter still accepts some red ClearType fringes.
A synthetic equal-G/B fringe regression fails with that filter.

The shared test sampler now requires filled 5x5 red patches and unions their
full extents, rejecting thin glyph fringes while retaining the bounds of the
larger solid markers, including their faded exits. Production components,
motion/focus assertions and CI configuration are unchanged; no test is skipped.

- The fringe regression passes, as do the responsive dialog, edge surface
  and drawer exit tests with the failed job's shuffle seed
  `1793111328520692000`;
  `/tmp/mygoui-pr3-windows-marker-regression.log`.
- `make check`: pass; `/tmp/mygoui-pr3-marker-check.log`.
- Linux Docker full overlay suite with the same shuffle seed: pass;
  `/tmp/mygoui-pr3-marker-linux.log`.
- Targeted overlay race tests with the same shuffle seed: pass;
  `/tmp/mygoui-pr3-marker-race.log` (111.204 s).
- Full lint: zero issues; `/tmp/mygoui-pr3-marker-lint.log`.
- Windows/amd64 overlay test-binary cross-compilation and WASM build: pass.

Windows execution subsequently passed on `4502b4e` in job `122257584865`, run
`40523761606`. This is a hosted headless Windows test result, not native GPU
or desktop input verification.

## MyGo 0.7.1 follow-up

2026-10-27, `coss`, based on `a8cd38e99fa564e4b162620d6313085e5476de3a`.
MyGo is now v0.7.1 and the minimum Go version is 1.27.3. This section
supersedes the older pinned-version statements below, not their historical
measurements. Remote CI and merge are tracked separately by PR #3.

- `make check`: pass. Full headless suite, build, Go vet, MyGo lifetime vet,
  formatting and module integrity. Final review-fix run:
  `/tmp/mygoui-071-check-reviewed.log`.
- `make race`: pass. Gallery 261.531 s and overlay 247.393 s;
  `/tmp/mygoui-071-race.log`. Subsequent data/navigation/overlay race passes
  are in `/tmp/mygoui-071-race-final.log`; final header-fix data race is in
  `/tmp/mygoui-071-data-race-reviewed.log`.
- `golangci-lint run ./...`: pass, zero issues;
  `/tmp/mygoui-071-lint-reviewed.log`.
- WASM and Windows/amd64 headless cross-builds: pass. `make workflows`: pass.
- Linux Docker full shuffled headless suite: pass;
  `/tmp/mygoui-071-linux-final.log`. The final header-fix data/Gallery
  shuffled recheck is in `/tmp/mygoui-071-linux-reviewed.log`.
- `govulncheck@v1.9.0 ./...`: no reachable vulnerabilities found;
  `/tmp/mygoui-071-vuln.log`.
- Native app build/sign with the pinned CLI: pass;
  `/tmp/mygoui-071-app-final.log`. Output:
  `artifacts/coss/mygo-071/app/darwin-arm64/Coss UI Gallery.app`.

New tests cover zoomed search/input and compact layout, native title-bar
insets, per-window zoom, inactive tab state/focus, loading paint without view
rebuilds, flyout option boundaries/anchors/surface, and initially visible
Cancel in a short 2x-text AlertDialog. `go tool mygo migrate-ui -write .`
moves safe constructor keys; its remaining Popover/Tooltip notices refer to
native base callback elements, which cannot be constructed by this layer.
MyGo lifetime vet passes without stored frame values.

Actual Metal offscreen scenes were rendered using the v0.7.1 module, not the
headless rasterizer: `artifacts/coss/mygo-071/`. Light/dark ScrollFade,
Modal and BottomSheet corners pass scoped pixel checks. Light ScrollFade,
dark BottomSheet, light FlyoutSurface and light card Table were inspected.
The card Table no longer draws a square header bottom border across the
rounded body's top corners. A red-first regression and selected-row/card
boundary tests cover that change. Gallery labels now call this the card
variant. The plain variant keeps its separate header rule and square body.

Independent pre-merge review identified that the card's transparent row
theme also hid header hover and resize grips. Card headers now restore those
affordances from committed native cell geometry while native sorting,
resizing/reordering and keyboard input remain untouched. Empty card headers
clip only their bottom divider and retain the native theme. A red-first
regression covers hover, resting/hovered grips, resizing and sorting in empty
and populated tables. No backdrop repaint or whole-table bitmap effect is
used to erase the rule.
Reviewed Metal captures include `table-grip-hover-false-metal.png` and
`table-card-true-metal.png`: header feedback remains outside the rounded body.

A separate nested-clip probe **still fails on Metal in 0.7.1**:
`rounded ancestor -> rectangular Clip -> red fill` paints a red corner where
the CPU image is white (`nested-clip-metal.png`). The explicit final-fill
clips/corner fills in changed components remain necessary. This is not a
claim that arbitrary rounded-container descendants are fixed; blanket
Effect/Cache remains rejected because it hides embedded WebViews and adds
texture-size/fade regressions.

macOS live accessibility observed the rebuilt Gallery and a separate native
flyout dialog with its heading, Name editor and Done action. The automation
then returned to the parent; no completed editor/dismissal sequence is
claimed. Windows caption appearance, Linux desktop flyout placement, native
gestures and actual screen-reader operation remain unverified locally.

## Clipping and CI follow-up

Follow-up on 2026-10-27 is based on `a8cd38e99fa564e4b162620d6313085e5476de3a`
on `coss`. Changes remain uncommitted; no push or remote CI rerun is asserted.
The older reconstruction record follows this section and does not override it.

- `references/coss` is checked out at
  `3d8ca28e7d90123e9e0f44665944dc240b026ccb`. Previous reference submodules
  are removed. Repository text outside the upstream checkout has no old-brand
  matches; original NextUI copyright remains in the retained Apache notice.
- Added red-first regressions for exiting per-corner button faces, vertical
  intrinsic group widths and the indeterminate progress band's paint bounds.
  The old CornerFill probe grew window references every rebuild and lost its
  fill on exit; the revised native background probe retains zero new handles
  and fades the face through intermediate frames.
- TableSecondary row/hover/cell-selection fills now respect the body's rounded
  boundary without clipping the header. Focus, virtualized controls and
  scroll-hover tests pass. Group examples use explicit separators; Vertical
  reproduces upstream `p-group-9.tsx`, with full-width actions a separate example.
- Actual Metal offscreen captures in
  `artifacts/coss/rounding-final-<false|true>-metal.png` cover light/dark Group,
  hovered Billing, vertical zoom controls and selected Table last row. Both
  images were inspected. These are GPU-rendered scenes, not window/OS input.
- The native audit was re-collected into
  `artifacts/coss/native-audit-clipping-final/manifest.json`: 57 subjects,
  114 appearance entries, 550 captures and 102 contact sheets. Capture count
  is evidence inventory, not a claim that every state matches the website.
- Browser reopened the official Group page and read live segment radius,
  borders, separators and size; two screenshot operations timed out. The
  earlier saved official vertical zoom example was inspected alongside the
  pinned source. Native app automation also timed out; no new OS-input claim.

### Known renderer boundary

An actual Metal probe confirms MyGo v0.6.3 can lose an outer rounded clip under
a rectangular descendant clip (`Card -> Clip -> fill`), while the CPU renderer
clips correctly. Native per-corner button fills and innermost progress/table
clips fix those concrete paths, not arbitrary consumer subtrees.

A blanket offscreen-effect workaround was tried and rejected: independent
review reproduced unsynchronized arrow/shadow fades, hidden embedded WebViews
and missing shadows beyond the effect size limit. Its race suite also timed
out. That workaround is **not in the final source**. The generic nested-clip
backend issue remains unresolved and must not be reported as fixed.

### Follow-up checks

- `make check MYGO_VERSION=v0.6.3`: pass;
  `/tmp/mygoui-coss-check-final-followup.log`.
- `make race`: pass, all packages, no reported races;
  `/tmp/mygoui-coss-race-final-followup.log` (Gallery 260.755 s,
  overlay 246.017 s). The earlier rejected-workaround timeout is not this run.
- Linux/arm64 Docker, Go 1.27.3, full
  `CGO_ENABLED=0 go test -shuffle=on ./... -count=1 -timeout=180s`: pass;
  `/tmp/mygoui-coss-linux-final-followup.log`. This is actual Linux test
  execution, not a cross-build.
- `make wasm`, Windows/amd64 cross-build, `golangci-lint run ./...`: pass.
- Remote run `40493219365` failed on the base commit. Its Linux/Windows logs
  identified a font-width-dependent fixture and paint samplers that can pick
  up glyph ink/subpixel color fringes. Tests now use a genuinely overwide
  fixture, ink-free sample and shared synthetic-marker detector. Windows
  runtime verification still requires remote CI; no platform skip was added.
- Native Gallery app/DMG rebuilt locally; no release, notarization or fresh
  full-window interaction acceptance is claimed.

## Earlier Reconstruction

Validated on 2026-10-27 in the then-uncommitted `coss` worktree based on
`8d1d2a336da159998edbedf2d12bd505b5f8cebe`. No commit, push or release is
asserted. These results apply to this reconstruction, not the previous design.

Environment: Go 1.27.3, darwin/arm64, macOS 27.0.1 (26A434), Apple M5 Pro
(Metal 4 capable). MyGo is pinned at v0.6.3 and mygoicons at v0.2.1;
`go.mod` minimum Go version is 1.27.1. `GOWORK` is unset and release module
metadata has no local replacement.

## Required checks

| Check | Status | Evidence to record |
| --- | --- | --- |
| `make check MYGO_VERSION=v0.6.3` | Pass | Full headless suite, build, Go vet, MyGo lifetime vet, gofmt and module verify. `/tmp/mygoui-coss-check-final.log`: Gallery 13.616 s, overlay 31.655 s. |
| `make race` | Pass | Final full run exited 0. `/tmp/mygoui-coss-race-verified.log`: root 28.798 s, Gallery 279.625 s, overlay 259.130 s; all packages pass without reported data races. |
| `golangci-lint run ./...` | Pass | Final invocation: `0 issues.` |
| `make wasm` | Pass: build only | Complete `GOOS=js GOARCH=wasm CGO_ENABLED=0 go build ./...`. No browser runtime claim. |
| `make workflows` | Not run; workflows unchanged | Run actionlint before release if workflow files change. |
| Dependency/module integrity | Pass | `go mod tidy`, `go mod verify`, no version upgrades, added direct dependencies or local replacements. |
| Documentation/API review | Source reconciliation performed | Re-read current inline Autocomplete, SegmentedBase radio control, Sidebar, sliding Sheet/Drawer, CardFrame, ScrollFade modes and optional font registration; no build/runtime result inferred |

Windows amd64 and Linux amd64 cross-builds also pass:
`GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build ./...` and the equivalent
`GOOS=linux` command. `git diff --check` passes. These are build/whitespace
checks, not native platform runtime evidence.

An initial Gallery check exposed lost Calendar keyboard focus with reduced
motion. Native transition registration was corrected and compact/wide tests
now pass. Five root joined-corner/seam/focus tests now explicitly use outline
groups; separate new tests protect default rounded items and 2 DIP gaps.
Assertions were not removed to manufacture a passing suite.

## Behavior and headless rendering

| Area | Status | Consumer-visible evidence needed |
| --- | --- | --- |
| Theme roles and patches | Pass: headless | Both appearances, all canonical patch roles, transparency, derived alpha, detached snapshots and native-theme installation. |
| Controls and surfaces | Pass: headless | State faces, geometry, native disabled/loading/invalid, no Button press scale, Group seams/corners and toggle focus. |
| Forms | Pass: headless | Editing, selection, bounds, filtering, OTP, native label/help/error association and Composer submit policy. Final help-text and number centering also visually rechecked. |
| ScrollArea / ScrollFade | Pass: headless; scoped native GPU check below | Axis/edge/RTL/endpoints/no-overflow, patterned alpha transparency and child activation/focus. |
| Command / Tabs / segments | Pass: headless plus queued-event probes | Filtering/activation, pointer/keyboard highlight, disabled or removed selection, moving-parent indicator and reduced motion. |
| Sidebar | Pass: headless | Selected/collapsed state, callbacks, disabled items, overflow and labels. |
| Overlays / dates / data | Pass: headless | Dismissal/focus, retained exit safety, placement/motion, PreviewCard handoff, date boundaries/keyboard and table selection/paging. OS context menu popup is not rasterized by Tester. |
| Gallery/examples | Pass: headless plus scoped native window | Snippets typechecked, compact/wide/enlarged-text interaction, state retention, search, appearance and titlebar inset cases. |
| Optional fonts | Pass: native registration/headless shaping | Concurrent/idempotent Register, regular/italic glyph advances at 400/700 versus original faces, visible text and no import-side initialization. |

`ui.NewTester` can validate native event routing, semantic metadata, layout
and headless pixel output. A recorded Tester test or screenshot proves only
what it actually exercises. Catalog row counts and gallery page counts are
not acceptance gates or evidence of pixel-perfect coss behavior.

## Platform evidence

| Target | Status | Evidence needed |
| --- | --- | --- |
| macOS / Metal | Scoped pass | Packaged app renders, actual traffic lights/titlebar alignment, Cmd+K search, Button navigation/action, Light/Dark, MaskFade gradient transparency and child activation/Page Down. Not exhaustive OS integration. |
| Windows / Direct3D | Build only | Actual native behavior/effects/fonts not verified. |
| Linux / available renderer | Build only | Actual native behavior/effects/fonts not verified. |
| WebAssembly / WebGL | Build pass; runtime not verified | Optional fonts explicitly reject browser registration; no browser runtime claim. |
| Screen reader | Not verified | Semantic metadata tests do not establish actual reader announcements/navigation. |
| Reduced motion and enlarged text | Headless pass | Native Calendar focus and Gallery resize/keyboard cases, animated control settlement. No exhaustive OS preference matrix. |

Native effects render subtrees through an intermediate layer, have an
8192-pixel side limit and omit embedded webviews. A CPU/headless mask test
cannot establish GPU effect correctness or performance. Fonts selected by
family name alone do not establish matching coss optical sizing/shaping.

## Visual comparison and real window

The per-component inspection and differences are recorded in
[docs/visual-review.md](docs/visual-review.md), not inferred from catalog counts.
Official rendered light/dark reference pages were captured for all 56 currently
available component URLs. The Sidebar URL returned 404; the actual docs sidebar
and MIT registry source supply its reference instead.

Inspected every mapped native component's light/dark baseline against rendered
reference samples. Repaired the three concrete visual findings: Field help/error
typography, NumberField centering and default ToggleGroup corner/gap behavior.
After production freeze the parent re-ran the collector successfully: 57
subjects, 114 appearance entries, 550 captures and 102 contact sheets. Final
manifest: `artifacts/coss/native-audit-final/manifest.json`; command log:
`/tmp/mygoui-coss-native-audit-final-parent.log`. Inspected corrected contact
sheets and the updated default Calendar/Card examples. The manifest records
unexercised states explicitly. No pixel-perfect or all-state parity is claimed.

Built and signed the macOS app with the pinned CLI:
`go run github.com/egoist/mygo/cmd/mygo@v0.6.3 build -platform darwin/arm64 -skip-dmg -skip-notarize -o artifacts/coss/app`.
Final log: `/tmp/mygoui-coss-app-final.log`; the 13.0 MB app is at
`artifacts/coss/app/darwin-arm64/Coss UI Gallery.app`. No notarization or
distribution claim. Opened the native app and inspected real traffic lights
at (16, 19) aligned with the 52 DIP title bar. Operated Cmd+K search, navigated
to Button, activated Default (counter advanced), opened Appearance and selected
Dark. Both appearances rendered in this scoped check.

On the real macOS Metal window, MaskFade revealed the blue/green gradient
beneath fading content. Activated a masked child button (counter advanced),
then Page Down moved from builds 01-04 to 05-08 with both overflowing edges
fading. The task's native screenshots hold this evidence. This is a scoped
render/input check, not a GPU benchmark or every-backend certification.

## Supplemental independent probes

Final local probe reruns passed queued Command filtering/Enter and repeated
activation, actual PreviewCard PointerLeave, SegmentedControl following parent
movement, mask background/RTL/no-overflow/input, and Tabs after removing the
selected item or receiving a disabled selected index.

Probes use a temporary Go overlay against pinned MyGo solely to inject queued
events and PointerLeave that ordinary Tester does not expose. Release module
metadata is unchanged. Files: `/tmp/mygoui_final_integration_review_test.go`,
`/tmp/mygoui_review_native_events_test.go`, `/tmp/mygoui_review_segment_test.go`,
`/tmp/mygoui_review_mask_test.go`; overlay/modfile:
`/tmp/mygoui_review_overlay.json`, `/tmp/mygoui_review.mod`. They supplement,
not replace, the portable repository suite.

No physical Windows/Linux window, browser runtime, screen reader, exhaustive
IME/clipboard/gesture behavior or per-component all-state OS matrix was checked.
Native menus/data-grid geometry/scrollbars/loading arcs/provider timing are
intentional adaptations. Composer is plain text and Drawer lacks snap-point
and nested-gesture coordination; see the component capability matrix.

## Documentation and provenance observations

The documentation pass read the downloaded coss site/LLM catalogs, pinned
`apps/ui` registry/MDX, mixed-license declarations and current Go sources.
It mapped 57 distinct component subjects and two hooks, recording partial
adaptations and missing upstream behavior rather than certifying completion.

Official font repository metadata and license text were inspected on
2026-10-27: Cal Sans at `9fe03d44f00806aadeb6d3653e37dce358554de0` and
Paper Mono at `e6eaeceaef02e77e3db997711e07a16378de2bd7`. Cal Sans upright/italic binaries were subsequently embedded unmodified with
the OFL notice and explicitly registered/tested; Paper Mono was not embedded. See [docs/fonts.md](docs/fonts.md)
and [docs/references.md](docs/references.md).
