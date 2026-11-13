# Design

## Context

See `proposal.md` for the motivation and decision D1 of
`complete-native-ui-solution` for the end state. The extraction reads the
implementation as it is in this checkout (MyGo v0.9.0), not upstream
memory. Sources read: `theme/palette.go`, `theme/tokens.go`,
`theme/theme.go`, `internal/motion`, `internal/busy`, `internal/anchored`,
`internal/style`, `internal/formstyle`, `overlay/motion.go`,
`overlay/chrome.go`, `scroll_fade.go`, `skeleton.go`, `status_dot.go`,
`spinner.go`, `progress.go` and every non-test Go file outside
`internal/gallery`, `examples` and `references` that names a duration, an
easing, an animation or a transition. MyGo's own defaults that components
inherit were read in the module cache (`ui/collapsible.go`, `ui/motion.go`,
`ui/widgets.go`, `ui/transition.go`).

Inventory method: 179 non-test library files; 51 lines name
`time.Millisecond` or `time.Second` (46 and 5) in 31 files; the shared
constants are referenced 31 times (`motion.State` 12, `motion.Reveal` 16,
`motion.Roll` 3) and their helpers 56 times (`motion.Color` 38,
`motion.Progress` 18). An animation without an `Ease` takes MyGo's default,
ease-out (`ui/motion.go` `retarget`, `ui/transition.go`), and
`Element.Animate` is ease-out; the tables below state the resolved easing.

The Impeccable `document` playbook asks the owner for a Creative North Star
and colour names. This work ran as a delegated, non-interactive task with no
question tool, so the record reuses the language already in `DESIGN.md`
and `PRODUCT.md` and names no new metaphor or identity.

## Goals / Non-Goals

**Goals:**

- `DESIGN.md` states every visual and motion rule the implementation
  follows, in its own terms, and is the authority.
- coss appears only as provenance (where a rule came from) and in the
  licensing boundary.
- The token inventory below is complete enough for `open-theme-tokens`
  to name and open each value without re-reading component files.

**Non-Goals:**

- Changing any value, Go file, test, file name, environment variable
  (`MYGOUI_COSS_*`) or branch name.
- Resolving the disagreements found. Each is a finding for a later change.
- Choosing token Go types or API shape; the names below are proposals.

## Decisions

### E1. `DESIGN.md` keeps the canonical eight sections; motion lives with what it moves

Impeccable's format allows only `colors`, `typography`, `rounded`,
`spacing` and `components` in frontmatter. Motion, shadows and effects go
in prose and in the sidecar's `extensions`. Overlay motion stays in
Elevation & Depth (it is how depth arrives); control state motion, loading
and effects go in Components. Alternative considered: a ninth "Motion"
section. Rejected: tooling parses the canonical headings, and the playbook
asks for motion to sit with the world or component it affects.

### E2. Values are recorded, not executed, in `DESIGN.md`

`DESIGN.md` states the values so a reader needs nothing else, and names
the Go file that executes each family. `theme/` and the component files
remain the executable source. When they differ, `DESIGN.md` wins: the
mismatch is a defect, settled by a change that corrects the code, or the
record if the record was wrong. Where the record states observed behavior
that may not be intended (F16), the finding names it for that change. Alternative considered: generate `DESIGN.md` from
Go. Rejected for this change: the motion values are not yet tokens;
`open-theme-tokens` makes that possible.

### E3. Disagreements are recorded as they are

Where components differ for the same kind of change, `DESIGN.md` states
the dominant rule and lists the exceptions as current behavior, not as
the rule. Findings below say which value dominates.

### E4. Provenance stays in one place

`docs/references.md` owns where each family of rules came from and the
MIT/AGPL boundary, which is kept word for word. `docs/theme.md` keeps the
token provenance and the documented/runtime dark-mix discrepancy as
history. Every other record points there instead of naming an authority.

### E5. Spec deltas restate authority, not behavior

The five requirements keep their behavior and scenarios. Names that
contained "coss" are RENAMED; the bodies are MODIFIED in full. The
elevation requirement now states the layer values, so it stands without
the coss shadow names. The `native-charts` Purpose line is edited in the
main spec. Requirements that say "upstream" about catalog parity
(`native-library` "Native adaptations have explicit boundaries",
`native-charts` "Chart catalog has native counterparts") describe
capability references, not visual authority, and stay.

## Findings: components that disagree

Each finding names the same kind of change, the values in use and where.
None was changed.

| # | Kind of change | Values in use | Where | Dominant |
| --- | --- | --- | --- | --- |
| F1 | Content opening/closing | 200 ms **ease-out** (MyGo `CollapsibleBase`, toolkit-owned); 200 ms **ease-in-out** | Accordion, Disclosure, Settings expander row, Timeline, AboutDialog sections; Spoiler (`spoiler.go:50`) | Ease-out (5 of 6). The previous `DESIGN.md` claimed ease-in-out for all |
| F2 | Colour fade of a state | 150 ms ease-out (`motion.Color`, `motion.State`); 200 ms ease-out; 100 ms ease-out | 40+ sites; Switch track (`form/switch.go:160`), PressableFeedback (`pressable.go:20`); close button (`link.go:20`), ResizablePanel handle (`resizable_panel.go:245`) | 150 ms |
| F3 | Hover | At once (the recorded rule); fades over 100 ms / 200 ms; scales over 100 ms | Every control; close button, PressableFeedback; ColorSwatch hover scale 1.1 (`color/color_swatch.go:209`) | At once |
| F4 | Panel collapsing to a narrower width | 200 ms **linear**; 200 ms **ease-out** | Sidebar (`sidebar.go:105`); ResizablePanel (`resizable_panel.go:164`) | None (one each) |
| F5 | Thumb grows under a held pointer | ×1.2; ×1.25 (both 150 ms ease-out); Switch stretches 10% after 100 ms instead | Slider (`form/slider.go:330`); ColorArea (`color/color_area.go:157`); Switch (`form/switch.go:164`) | None |
| F6 | Modal surface and scrim entering and leaving | 200 ms ease-in-out (scale/rise + fade, and its scrim); 200 ms ease-out fade of surface and backdrop together; 450 ms drawer-curve scrim; 450 ms ease-out backdrop fade under a 450 ms drawer-curve slide | Dialog (responsive and slide) and AlertDialog; Modal (`overlay/modal.go:82`), ImageViewer (`overlay/image_viewer.go:123`) and side Dialog without slide, all through `surfaceMotion` (`overlay/motion.go:201-209`); Drawer (`overlay/drawer.go:94`); BottomSheet (`overlay/drawer.go:328`) | Ease-in-out for Dialog, ease-out for the rest |
| F7 | Swapping a view in place | 200 ms ease-out in; 120 ms ease-out in and out; 150 ms ease-out | Calendar year grid; Calendar day grid; Calendar header arrows and chevron (`date/calendar.go:279,293,320,348,364`) | None within one component |
| F8 | Indeterminate sweep | 40% band from −100% to +350% of its width; 40% band from −50% over 400% (both 1.5 s ease-in-out, 60 fps cap) | Progress (`progress.go:109`); refresh line (`internal/busy` `Line`) | None |
| F9 | Skeleton pulse | Fills towards the Surface colour up to 50%; halves the content's alpha | `Skeleton` with `SkeletonPulse`; `SkeletonGroup` with `SkeletonPulse` (shader) | None |
| F10 | Fading content at a clipped edge | 24 DIP alpha mask or shadow (ScrollFade default); 40 DIP surface fade (ScrollShadow); 24 DIP alpha mask (DocumentTabs strip, own gradient); 32 DIP alpha mask (Spoiler cut) | `scroll_fade.go`, `scroll_shadow.go`, `navigation/document_tabs.go:497`, `spoiler.go:27` | 24 DIP |
| F11 | Easing chosen implicitly | No `Ease` given, so MyGo's ease-out applies | Switch thumb, Slider thumb press, Toast content fade, NumberFlip, Tree rows, `motion.Progress`, `motion.Color`, `busy` fades | Ease-out; a toolkit default change would move these silently |
| F12 | Branch reveal | Rows 200 ms ease-out (implicit); arrow 200 ms ease-in-out | Tree (`list/tree.go:833,917`) | None |
| F13 | Pager indicator | The previous record listed "pager" as a 200 ms slide; PipsPager grows the current pip over 150 ms ease-out | `navigation/carousel.go:151` | Recorded as it is |
| F14 | One value, several sources | 200 ms indicator slide is `motion.Reveal` (Tour dots), `tabsIndicatorMotion` (Tabs) and a literal (SegmentedControl); 150 ms popup equals `motion.State` but is `anchored.MotionDuration` | `navigation/motion.go:14`, `navigation/segmented_control.go:102`, `internal/anchored/surface.go:13` | Values agree; sources do not |
| F15 | Type-ahead window (behavior, for `component-behavior-reference`) | 700 ms; 1 s | Tree (`list/tree.go:163`); menus (`navigation/menu_ahead.go:11`) | None |
| F16 | Elevation edge direction (possible defect) | A headless `ui.NewTester` render on a mid-grey ground shows the inset edge of `OverlayShadows`, `FieldShadows` and `SurfaceShadows` along the **top inside** row in light appearance (128 → 123/120) and the **bottom inside** row in dark (128 → 135). The previous record and `docs/theme.md` said "dark lower edge in light, light upper edge in dark", the coss intent this was copied from | `theme/tokens.go:241-242`, MyGo `ui/paint.go:1663` `insetShadow` | Record states the observed edge; a later change decides which is intended |
| F17 | Busy region while loading | A loading chart shows a skeleton but does not mark its region busy, unlike the other data regions (`settings.go:67` `Busy`) | `chart/chart.go:85` | Busy region (the recorded rule) |

## Token inventory for `open-theme-tokens`

Names are proposals in a portable dotted form; `open-theme-tokens` picks
the Go spelling. "Shared" means one source today.

### Durations and easings

| Proposed token | Value | Easing | Used for | Where |
| --- | --- | --- | --- | --- |
| `motion.duration.state` | 150 ms | ease-out | Colour, fill, border and ink of a state; a mark drawing in; icon crossfade; message appear; tag, rate star, chip, list check, carousel pip, tags-input refusal, image viewer chrome, log jump button, document tab face and fade, document tab scroll | Shared: `internal/motion.State`, `motion.Color`, `motion.Progress`, `motion.go` |
| `motion.duration.reveal` | 200 ms | ease-in-out (content, indicator) or ease-out (layout, value) | Spoiler; Tree rows and arrow; slider value set by the application; document tab reflow and scroll-into-view; carousel pull-back; image viewer stage; Tour scrim, height and dots; ActionBar entry; Banner exit | Shared: `internal/motion.Reveal` |
| `motion.duration.roll` | 300 ms | ease-out (digits), ease-in-out (travel), linear (crossfade) | NumberFlip digits; Tour highlight travel and step crossfade | Shared: `internal/motion.Roll` |
| `motion.duration.collapsible` | 200 ms | ease-out | Accordion, Disclosure, Settings expander, Timeline, AboutDialog | MyGo `ui/collapsible.go:38` (toolkit-owned; F1) |
| `motion.duration.indicator` | 200 ms | ease-in-out | Tabs indicator, SegmentedControl thumb | `navigation/motion.go:14`, `navigation/segmented_control.go:102` (F14) |
| `motion.duration.popup` | 150 ms | ease-out | Popover, pickers, TeachingTip, tooltip (enter 98% → 100%, exit fade, tooltip/preview also 98%); FloatingPanel fade | `internal/anchored/surface.go:13` |
| `motion.duration.dialog` | 200 ms | ease-in-out (Dialog, AlertDialog) / ease-out (Modal, ImageViewer, side Dialog without slide) | Dialog scale 0.98 (≥ 640 DIP) or 16 DIP rise; edge sheet 32 DIP slide; Dialog and AlertDialog scrim; Modal and ImageViewer surface-and-backdrop fade (F6) | `overlay/motion.go:13` |
| `motion.duration.drawer` | 450 ms | `motion.ease.drawer` | Drawer slide and scrim; BottomSheet entry and drag settle; FloatingSheet | `overlay/motion.go:14`, `overlay/drawer.go:10` |
| `motion.ease.drawer` | — | cubic-bezier(0.32, 0.72, 0, 1) | Drawer family | `overlay/drawer.go:10` |
| `motion.duration.progress` | 500 ms | ease-out | Progress bar width and ring value | `progress.go:38` |
| `motion.duration.page` | 250 ms | ease-out | Carousel slide; NavigationStack page slide and colours | `navigation/carousel.go:71`, `navigation/navigation_stack.go:56` |
| `motion.duration.shift` | 150 ms | ease-out | Rows and tabs moving aside during a drag (SortableList, DocumentTabs) | `list/sortable_list.go:50`, `navigation/document_tabs.go:129` |
| `motion.duration.panel` | 200 ms | linear (Sidebar) / ease-out (ResizablePanel) | Panel width collapse (F4) | `sidebar.go:105`, `resizable_panel.go:164` |
| `motion.duration.press` | 250 ms | ease-out | PressableFeedback scale to 0.985 | `button.go:63` |
| `motion.duration.highlight` | 200 ms | ease-out | PressableFeedback hover/press fill (F2, F3) | `pressable.go:20` |
| `motion.duration.hover` | 100 ms | ease-out | Close button fill; ResizablePanel handle colour; ColorSwatch scale (F2, F3) | `link.go:20`, `resizable_panel.go:245`, `color/color_swatch.go:209` |
| `motion.duration.thumb` | 150 ms | ease-out | Switch thumb position; Slider thumb press ×1.2; ColorArea thumb press ×1.25; ColorSwatch check | `form/switch.go:147`, `form/slider.go:330`, `color/color_area.go:157`, `color/color_swatch.go:227` |
| `motion.duration.switchTrack` | 200 ms | ease-out | Switch track colour (F2) | `form/switch.go:160` |
| `motion.switch.stretch` | 100 ms after 100 ms | ease-out | Held switch thumb widens 10% | `form/switch.go:164` |
| `motion.duration.linkIcon` | 150 ms | ease-out | Link trailing icon opacity | `link.go:70` |
| `motion.duration.calendarSwap` | 200 ms / 120 ms / 150 ms | ease-out | Year grid in; day grid in/out; header arrows fade and chevron turn (F7) | `date/calendar.go` |
| `motion.spring.toast` | spring, 320 ms | spring | Toast deck position and size, swipe | `overlay/toast.go:252` |
| `motion.duration.toastContent` | 200 ms | ease-out (implicit) | Toast content fade as the deck expands | `overlay/toast.go:256` |
| `loading.duration.fade` | 150 ms | ease-out | Content ↔ spinner crossfade; region indicator fade | `internal/busy` `fade` |
| `loading.delay` | 150 ms | — | Region indicator wait before showing | `internal/busy` `wait` |
| `loading.period.spinner` | 1 s per turn | linear | Spinner, busy spinner, painted field spinner | `spinner.go:49`, `internal/busy` |
| `loading.period.sweep` | 1.5 s | ease-in-out | Indeterminate progress and refresh line (F8) | `progress.go:109`, `internal/busy` `Line` |
| `effect.skeleton.period` | 2 s | linear phase | Skeleton shimmer and pulse, SkeletonGroup shader | `skeleton.go:84` |
| `effect.statusDot.pulse` | 2 s | ease-out ring growth | StatusDot pulse | `status_dot.go:29` |
| `effect.statusDot.breath` | 1.6 s | cosine | StatusDot breathe | `status_dot.go:31` |
| `effect.caret.blink` | 500 ms cycle, on for the first half | — | InputOTP caret | `form/input_otp.go:135` |

Interaction timing, not motion, recorded for `component-behavior-reference`:
copied feedback 1.5 s (`content/copy_button.go:36`); Tree type-ahead
700 ms and reveal slack 100 ms (`list/tree.go:163,777`); menu type-ahead
1 s (`navigation/menu_ahead.go:11`); Alt tap 600 ms
(`navigation/menubar.go:48`); carousel wheel quiet 150 ms
(`navigation/carousel.go:82`); PreviewCard leave delay 150 ms
(`overlay/preview_card.go:71`).

### Spacing and geometry literals

| Proposed token | Value | Used for | Where |
| --- | --- | --- | --- |
| `space.unit` | 4 DIP | Already public (`theme.SpaceUnit`, `Space(n)`) | `theme/tokens.go` |
| `motion.distance.dialogRise` | 16 DIP | Compact dialog rises on entry | `overlay/motion.go` |
| `motion.distance.sheetSlide` | 32 DIP | Edge-anchored dialog slide | `overlay/dialog.go:77,122` |
| `motion.scale.popupEnter` | 0.98 | Popup and wide dialog scale on entry; tooltip/preview on exit | `internal/anchored/surface.go`, `overlay/motion.go` |
| `layout.dialogCompactBelow` | 640 DIP | Dialog switches from scale to rise; inset sheets | `overlay/motion.go`, `overlay/dialog.go` |
| `focus.ring.fieldWidth` | 3 DIP outside | Field focus ring | `internal/formstyle` `FieldRing` |
| `focus.ring.inset` | 2 DIP, 1 DIP inside | Joined segment focus ring | `internal/style` `insetRing` |
| `overlay.backdropBlur` | 4 DIP | Default dialog backdrop blur | `overlay/chrome.go:29` |
| `overlay.drawerWidth` | 448 DIP | Side drawer default width | `overlay/drawer.go:88` |
| `overlay.sheetMaxHeight` | 85% of window | BottomSheet cap | `overlay/drawer.go` |
| `overlay.tooltipMaxWidth` | 320 DIP, padding 4 × 8 | Tooltip | `overlay/tooltip.go` |
| `loading.lineHeight` | 2 DIP | Refresh line | `internal/busy` `LineHeight` |
| `fade.edge` | 24 DIP | ScrollFade default, DocumentTabs strip (F10) | `scroll_fade.go`, `navigation/document_tabs.go:484` |
| `fade.shadowEdge` | 40 DIP | ScrollShadow default | `scroll_shadow.go:39` |
| `fade.spoilerCut` | 32 DIP | Spoiler cut | `spoiler.go:27` |
| `button.paddingX` | 7 / 9 / 11 / 13 / 15 DIP (xs–xl); 0 icon-only | Button | `button.go:216` |
| `button.gap` | 8 DIP; icon 16 DIP at 80% opacity, −2 DIP side margin | Button content | `button.go` |
| `scrollbar.width` | 6 DIP | Native scrollbar | `theme/theme.go` |
| `statusDot.size` | 8 DIP | StatusDot | `status_dot.go` |
| `skeleton.groupGap` | 8 DIP | SkeletonGroup rows | `skeleton.go` |

### Effect parameters

| Proposed token | Value | Used for | Where |
| --- | --- | --- | --- |
| `loading.refreshDim` | 55% opacity | Content under a refresh | `internal/busy` `dim` |
| `loading.sweepBand` | 40% of track | Indeterminate band (F8) | `progress.go`, `internal/busy` |
| `effect.spinner.arc` | 1.8 rad, stroke size × 2.5/24 | Painted spinner | `internal/busy` `PaintSpinner` |
| `effect.frameRate.continuous` | 60 fps | Skeleton, indeterminate sweep, refresh line, spinner | `skeleton.go`, `progress.go`, `internal/busy` |
| `effect.frameRate.ambient` | 30 fps | StatusDot pulse and breath | `status_dot.go` |
| `effect.skeleton.angle` | 120° | Shimmer direction | `skeleton.go` |
| `effect.skeleton.band` | 0.4 / 0.5 / 0.6 stops; shader ±0.1 | Shimmer width | `skeleton.go` |
| `effect.skeleton.peak` | white 64% light, 4% dark | Shimmer highlight | `skeleton.go` `skeletonHighlight` |
| `effect.skeleton.pulse` | 50% | Pulse depth (F9) | `skeleton.go` |
| `effect.scrollFade.ink` | Foreground at 18% | Shadow edge and mask fallback | `scroll_fade.go` |
| `effect.statusDot.ring` | 40% → 0%, grows to 2× | Pulse ring | `status_dot.go` |
| `effect.statusDot.breathRange` | 65%–100% | Breathe | `status_dot.go` |
| `effect.press.scale` | 0.985 | PressableFeedback | `pressable.go:55` |
| `effect.thumb.pressScale` | 1.2 / 1.25 (F5) | Slider, ColorArea | `form/slider.go`, `color/color_area.go` |
| `effect.swatch.scale` | hover 1.1, chosen 0.77; check 0.8 → 1 | ColorSwatch | `color/color_swatch.go` |
| `effect.switch.stretch` | 10% | Held switch thumb | `form/switch.go` |
| `color.hoverOpacity` | 90% of the fill | Solid button and status hover | `theme/palette.go`, `button.go` |
| `color.softAlpha` | 8% light, 16% dark; hover +4% | Soft status fills | `theme/palette.go` |

## Risks / Trade-offs

- [`DESIGN.md` drifts from Go values] → It names the executing file for
  each family; `open-theme-tokens` makes motion and effects tokens so the
  record can cite them. Until then, a value change updates `DESIGN.md` in
  the same change.
- [coss-named tests read as the authority] → They keep running and keep
  their names (non-goal); `docs/visual-review.md` and the tests compare
  rendering with the pinned reference as provenance.
- [Removing coss as authority loses a tie-breaker] → `DESIGN.md` now
  states the dominant rule per finding; a missing rule is added to the
  record through a change.
- [Concurrent edits to shared records] → Only the Visual Authority,
  Capabilities, Native Gallery and Development Context parts of
  `PRODUCT.md`, the Design Work and Engineering authority sentences of
  `AGENTS.md`, the design paragraphs of `docs/agent-workflow.md` and the
  authority lines of `README.md` and `CONTRIBUTING.md` are edited, each
  re-read before a small edit.

## Open Questions

- Whether F1–F14 are settled to one value each is a design decision for a
  later change (or for `open-theme-tokens` when it names the tokens).
