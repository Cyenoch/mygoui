# Validation

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

Actual Windows execution remains pending remote CI; cross-compilation is not
a Windows runtime result.

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
