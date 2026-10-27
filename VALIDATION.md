# Validation

Validated on 2026-10-27 in the uncommitted `coss` worktree based on
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
