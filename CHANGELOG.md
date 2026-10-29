# Changelog

## Unreleased

The `coss` reconstruction replaces the previous unreleased appearance
and token contract with coss ui references and native Go APIs. Existing extra
families remain library extensions; the catalog is a capability mapping,
not a claim of pixel-perfect completion. No release is asserted by this log.

### Changed

- Upgrade MyGo to v0.7.2 and Go to 1.27.3; the CLI now resolves from the
  module tool directive instead of an unpinned latest version.
- Avatar uses native `ImageSize`, including typed-nil and empty image fallback;
  loaded bitmap/SVG bytes can be decoded with MyGo's `DecodeImage` directly.
- Single-choice context menus use native radio semantics, not checkboxes.
- Gallery supports native content zoom and Windows caption theme updates.
  Loading animations use paint-only frames, with skeleton motion capped at
  60 fps. Short-window alert autofocus reveals the safe action after layout.

- Replace previous reference submodules with `references/coss`, pinned to the
  documented source commit; retain earlier copyright notices separately.
- Native per-corner button backgrounds retain their face through exits without
  allocating persistent handles per frame. Vertical groups share member width;
  Gallery groups include explicit separators.
- Progress sweeps and table row states stay inside their rounded bounds.
  Arbitrary nested rounded clipping remains a documented toolkit limitation;
  no blanket offscreen-effect workaround changes native content composition.
- Cross-platform paint tests distinguish synthetic markers from subpixel text
  edges and use values that genuinely exceed the available API column width.

- Canonical theme roles follow coss: Primary/Secondary, separate muted
  background/foreground, Accent, Card/Popover, Destructive, Input/Ring and
  sidebar/chart roles, with per-appearance palette patches.
- Default geometry uses a 10 DIP radius base and coss desktop control sizes:
  buttons 32 DIP, small/large 28/36 and extra small/large 24/40. Body,
  heading and mono families are independently configurable. System fonts
  remain the theme default; applications own explicit font registration.
- Button states use coss color/border/shadow feedback without a pressed scale.
  Native disabled opacity remains 50% rather than the reference's 64%.
- Badge uses its own BadgeVariant: Primary, Secondary, Outline, Destructive,
  Error, Info, Success and Warning. BadgeOptions contains Label, Variant,
  Size and Icon.
- Documentation now identifies the MIT coss registry/docs, the upstream
  AGPL package boundary, dark palette provenance differences and independent
  OFL font sources. Historical references are no longer design authority.

### Added

- `navigation.PopoverMenu` opts into MyGo's themed in-window context menu
  while `ContextMenu` keeps its OS-owned default. Both share selection and
  callback wiring; the Gallery demonstrates the themed variant.
- Gallery-only `e2e` automation import for opt-in real-window tests over a
  private socket; ordinary builds and library imports do not enable it.
- Native `overlay.NewFlyout`, coss `FlyoutSurface`, zoom-aware `FlyoutAnchor`
  and a Popover Gallery example. Explicit parent/tray and dimensions keep
  window ownership outside the build; WebAssembly returns an error.
- Optional `TabsOptions.KeepMounted` preserves hidden panels' native state
  without layout space or keyboard targets.

- Root Frame with independent panels, mixed Group with explicit segment
  corners, caller-owned collapsible Sidebar, ScrollArea with
  vertical/horizontal/both axes, and reusable ScrollFade shadow, surface and
  native shader alpha-mask treatments.
- Plain-text form.Composer with caller-owned draft, persistent actions,
  Enter/Shift+Enter policy and caller-controlled submission/clearing.
- Inline editable Autocomplete with focused-input suggestion navigation.
- navigation.Command searchable actions, native ContextMenu and independent
  radio-semantics SegmentedControl with animated selection.
- CardFrame helpers, overlay.Sheet edge slide, Drawer slide/inset presentation
  and PreviewCard with hover/focus continuity.
- A 57-component coss mapping plus two native hook/service translations,
  separate extra-family inventory and explicit behavior/platform gaps.
- Optional embedded Cal Sans regular/italic from the independent official
  OFL source, with explicit `fonts.Register() error` startup registration.
  Browser registration fails explicitly; Paper Mono remains host-provided.
- Redesigned native gallery with a custom 52 DIP title bar, platform window
  control insets, grouped sidebar, keyboard search and Preview/Code/API/Tokens
  views. macOS keeps native traffic lights; other platforms reserve the
  toolkit-reported control area.
- Per-component official-browser/native visual review and an opt-in capture
  collector, with explicit missing-state and platform evidence boundaries.

### Known limits and validation

Composer rich tokens/mentions, Drawer snap points/nested gesture behavior
and automatic form validation/submission are not supplied. OS context menus
and tooltip timing remain native.
Mask effect dimensions, webview exclusions and native font rendering need
explicit attention; see the component and scrolling docs.

Current reconstruction checks are recorded in [VALIDATION.md](VALIDATION.md).
Earlier rewrite results do not certify this source state. GPU, platform
input/gesture and actual screen-reader evidence are separate from headless
layout, semantic metadata and pixel tests.
