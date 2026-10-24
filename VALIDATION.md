# Validation

Local validation used macOS arm64, Go 1.27.3, MyGo v0.6.3 and mygoicons
v0.2.1. Module metadata targets Go 1.27.1. The public repository is
`github.com/Cyenoch/mygoui`; no release has been tagged.

## Completed checks

- The overlay-motion update passed `make check MYGO_VERSION=v0.6.3`:
  `go vet ./...`, the pinned MyGo lifetime vet, the ordinary suite,
  `go build ./...`, formatting and `go mod verify`.
- All packages with tests passed: root, `theme`, `form`, `list`, `navigation`,
  `overlay`, `date`, `color`, `data`, and `internal/gallery`. Examples and
  internal helpers have no tests. Fixed test counts are not retained here;
  they become stale as behavior regressions are added or removed.
- `golangci-lint run ./...` reported zero issues and `make wasm` passed after
  the final modal identity and AlertDialog focus fixes.
- Toast's follow-up semantic-name smoke verified title and description-only
  notifications, native announcements, same-ID renaming and removal on close;
  internal IDs did not appear in the native tester's semantic labels.
- The additional checks below record the earlier rewrite validation, not reruns
  for the overlay-motion or semantic-name update.
- The last recorded `make race` run predates the overlay-motion update. It
  passed with `CGO_ENABLED=1` and no data races; `internal/gallery` took 65.2 s.
- `make workflows` (actionlint) accepted `.github/workflows/ci.yml`.
- `govulncheck v1.9.0 ./...` found no reachable vulnerabilities.
- `gofmt -l .` was empty; `go mod tidy` left `go.mod`/`go.sum` unchanged.
- The MyGo CLI build and dev paths were exercised from the repository root:
  `go run github.com/egoist/mygo/cmd/mygo@v0.6.3 build` produced a signed
  development-era `.app` and `.dmg`, and `dev -skip-dev-command` built the
  gallery, watched for changes and stopped cleanly under a timeout. The root
  `mygo.json` points `main` at `./examples/gallery`. `make dev` runs the same
  CLI at `@latest`.

## Gallery evidence

- The gallery registers 87 pages: the Overview landing page plus 86 component
  pages, each added through `RegisterComponent(ComponentPage{...})`. No family
  page registrations remain; `TestCatalogIsUniqueAndComplete` checks unique IDs
  and ordering.
- All 86 component pages use `ComponentPage.Examples`
  (`ComponentExample{Title, Description, Demo}`). `ComponentPage` has no
  page-level `Demo` field and the fallback plus its preview collector were
  deleted, so the named frames are the only demo source. The catalog contains
  260 named examples across 10 groups: Actions 19 (6 pages), Collections 5 (1),
  Color 19 (6), Data 4 (1), Dates 18 (6), Display 38 (13), Forms 81 (24),
  Foundations 28 (11), Navigation 30 (10) and Overlays 18 (8). Every page has
  2-5 named examples (minimum 2, maximum 5); none has fewer than two.
- The docs shell follows HeroUI: an 880 DIP main column (800 DIP of content),
  a 220 DIP "On this page" rail shown at 1000 DIP and hidden below it, and the
  compact sidebar branch below 760 DIP. Section links use committed
  `ui.Handle` bounds to set the content scroll and derive the active entry
  (`TestDocsTOCScrollsTheContent`, `TestDocsRailHidesOnNarrowWindows`).
- The page header shows the group label, 30/36 title, 16/28 muted description,
  rule, and Copy Markdown/Copy import pills with native clipboard feedback.
  Usage is a 44 DIP import bar with a 32 DIP icon-only copy control that grows
  for a parenthesized block. No page sets `ComponentPage.Import`: all 86
  blocks are derived from the qualifiers the page's `Code` names (module root
  always present, module subpackages and `mygoicons/lucide`, de-duplicated and
  sorted, one line or a block), and the header copy writes the same block
  (`TestDocsHeaderActionsCopyRealContent`, `TestUsageBarIsACompactImportLine`,
  `TestComponentPageImportStatement`,
  `TestComponentPageImportStatementDerivesFromCode`,
  `TestUsageAndHeaderCopyDerivedImport`, `TestUsageBarGrowsForMultipleImports`).
- Every named example renders its title and description outside its own
  clipped rounded panel; the page's full code snippet renders exactly once
  (inside a single example's panel, or in one shared `Code` panel after two or
  more). `TestComponentExampleFrames`, `TestDemoTitlesSitOutsidePreviewPanels`
  and `TestSingleExampleOwnsItsCode` pin this structure.
- Code areas carry line numbers, Go highlighting, selectable lines, horizontal
  scrolling, a tooltip-only copy control with Copied feedback, and a six-line
  fold with a surface-colored fade plus Expand code/Collapse code. Copy controls
  are 32x32, centered within 1 DIP and inset 8-12 DIP in both narrow and wide
  layouts (`TestCodeBlockCopyControl`, `TestCodeBlockLineNumbersAndHighlighting`,
  `TestCodeBlockHugsSnippet`, `TestCodeBlockScrollsLongLines`,
  `TestCopyControlsAlignAndDoNotStretchRows`).
- Rounded gallery surfaces clip their children to the corner: example panels,
  the shared code panel, the Usage bar, the API panel, token cards and overview
  group cards. `TestRoundedSurfacesClipCorners` samples the four outside-corner
  pixels against the page background and one inside pixel against the surface,
  in light and dark modes, while keeping copy focus rings clear of the clip.
- Long API type/default values wrap instead of clipping and stack at narrow
  widths (`TestAPILongValuesWrapInsteadOfClipping`,
  `TestAPILongValuesWrapInStackedLayout`, `TestAPIShowsFullTypeValues`,
  `TestComponentPageHeaderAndAPIPills`); theme-token cards show light/dark
  swatches, hex values and roles (`TestThemeTokenCardsShowValuesAndRoles`,
  `TestThemeTokenCardsMarkCustomTokens`).
- The Overview page is a landing experience: hero, feature points, and one card
  per component group with page counts and links. The sidebar shows the MyGo UI
  brand, lists every component page by group, and uses ScrollShadow for
  overflow (`TestSidebarShowsEveryComponentPage`, `TestSidebarUsesScrollShadow`,
  `TestCompactDrawerUsesScrollShadow`).
- `TestEveryPageRendersAtDesktopAndCompactWidths` renders every registered page
  at 640x900 and 1100x900 in light mode and 1100x900 in dark mode.
  `TestPagesRenderWithEnlargedText` renders every page at 1.75x text scale.
  `TestCompactNavigationOpensAndSelects` exercises the sub-760 DIP navigation.
- Toast exit-paint regressions cover the close-button crash
  (`TestToastPageDismissesWithoutExpiredPaint`,
  `TestToastTopPlacementDismissDuringExit`,
  `TestToastDismissDuringExitKeepsPaintSafe`).
  `TestToastCloseGlyphIsOneCenteredX` checks the floating upper-right close
  button, normal action right padding, centered X and keyboard focus ring;
  `TestToastCloseRevealsOnHoverAndDismisses` checks hover reveal and dismissal.
- Toast motion uses MyGo's spring with a 320 ms response, not a fixed-duration
  250 ms transition. `TestToastEntryTravelsAndFadesFromWindowEdge` samples
  entry translation and fade at both window edges. Deck tests cover continuous
  insertion/removal, unequal-height hover/focus expansion, stable gap hover and
  timeout pause, rapid reversal, same-ID focus/drag continuity and swipe
  cancel/commit. `TestToastReducedMotionAndRetainedExit` checks immediate
  reduced-motion settlement and painted but inert exit snapshots.
- Modal-family and popup tests sample first, intermediate and settled paint
  and require pending animation frames on entry and exit in isolated Testers
  without a Toast viewport. They also cover immediate reopen, backdrop
  dismissal, AlertDialog Cancel focus/restoration, inert departing controls,
  reduced motion and all four Drawer edges with 16 DIP gutters.
- A throwaway `ui.NewTester` smoke exercised actual entry/exit rendering for
  Modal, Dialog, AlertDialog, BottomSheet, Popover, Tooltip and all four Drawer
  directions. Every path painted intermediate frames, removed closed semantic
  content and stopped requesting frames after exit. Separate Toast smokes
  observed collapsed/intermediate/expanded geometry and verified Undo's 16 DIP
  right padding with a floating close button that dismissed independently.
- The diagonal indicator glyphs stroke a `ui.Path` instead of
  `Painter.Line`, which only draws horizontal and vertical strokes: the
  disclosure and accordion down chevrons and the menu check/right chevron
  (`navigation/glyph_pixel_test.go`: `TestDisclosureAndAccordionChevronsAreDrawn`,
  `TestMenuIndicatorGlyphsAreDrawn`) and the selected color swatch tick
  (`color/checkmark_pixel_test.go`: `TestColorSwatchSelectedCheckIsDrawn`,
  `TestColorSwatchPickerSelectedCheckIsDrawn`) are each asserted as two
  diagonal strokes that meet at the right vertex.
- Hover faces repaint from build-time snapshots when the pointer moves,
  checked in the headless tester: Button primary/ghost, standalone
  ToggleButton default/selected/ghost, joined ghost ButtonGroup, Checkbox,
  Radio, Switch, InputGroup, Tabs, Accordion header, ListBox rows, calendar
  days and DateField (`TestButtonHoverAndPressedPaint`,
  `TestToggleButtonHoverAndPressedPaint`, `TestButtonGroupGhostHoverPaint`,
  `TestCheckboxHoverAndPressedPaint`, `TestRadioHoverAndPressedPaint`,
  `TestSwitchHoverPaint`, `TestInputGroupHoverPaint`, `TestTabsHoverPaint`,
  `TestAccordionHeaderHoverPaint`, `TestListBoxRowHoverPaint`,
  `TestCalendarDayHoverPaint`, `TestDateFieldHoverPaint`). Surfaces with no
  upstream hover rule stay inert or pressed-only (`TestDisclosureHoverIsInert`,
  `TestColorSwatchPickerHoverGrowsSwatch`, `TestColorAreaThumbPressed`,
  `TestSliderPressedShrinksKnob`).
- Exits keep painting retained copies instead of expired Elements:
  `TestToastDismissDuringExitKeepsPaintSafe`,
  `TestToastPageDismissesWithoutExpiredPaint`,
  `TestToastTopPlacementDismissDuringExit` (the gallery Toast page),
  `TestModalWithButtonClosesWithoutExpiredPaint` and
  `TestButtonPaintsThroughExitTransition` (Button, ToggleButton and
  ButtonGroup inside an `ElementTransition` exit, which used to panic in
  `style.Background` and `style.CornerFill`).
- Page-local tests open and dismiss modals, drawers, bottom sheets, dialogs,
  popovers and menus; commit dates; select options; sort and select table rows;
  run toasts; and mutate form, toggle, tag and progress state.
- Pixel tests verify Button/Card fills, disabled-state single-dimming, calendar
  range endpoints (solid accent) versus middle days (soft accent), joined-group
  clipping and inset focus (`TestToggleButtonGroupFillStaysInsideRoundedBox`,
  `TestButtonGroupFillStaysInsideRoundedBox`,
  `TestGroupedFocusRingIsInset`), and grouped button width sharing. Select
  single/multiple/clear, ComboBox/Autocomplete pointer selection, DatePicker
  dismissal and progress's default 0–100 range have direct regressions.

## Corrections made during acceptance

- Calendar range days are no longer marked with native `Selected`, which painted
  every middle day in solid accent; only endpoints paint solid, middle days use
  the HeroUI soft band, and a pending range with only a start day paints no band
  (`TestRangeCalendarStartOnlyPaintsNoBand`). Calendar cells are square and use
  HeroUI's clamped 24 DIP corner.
- The calendar year picker pages 12 years at a time and clamps the page to the
  configured bounds (`TestCalendarYearPickerDefaultPage`,
  `TestCalendarYearPickerBounds`).
- The per-component migration left the old family-page builders unreachable;
  26 unused functions and the `actionsOpen` field (27 `unused` findings) were
  deleted along with an orphaned import. `golangci-lint run ./...` reports zero
  issues.
- Root intrinsic components (Button, ToggleButton, Chip, Badge, Tag, Kbd, Link)
  opt out of the default column stretch; full-width variants still fill.
- Field surfaces paint a 2 DIP focus ring; disabled nodes are dimmed once by the
  toolkit rather than by both an outer and inner node.
- Popover/Tooltip and all form/date/color popups share the dark inset hairline;
  top/left anchor offsets move panels away from, not into, their anchor.
- Selection in ComboBox/Autocomplete is applied once on pointer input; the
  custom filter receives raw query text; section items no longer duplicate as
  ungrouped rows.
- DatePicker notifies `OnOpenChange(false)` on Escape/outside dismissal; the
  DateRangePicker popup autofocuses its grid.
- The page-level `Demo` fallback and its preview collector were removed;
  `ComponentPage.Examples` is the only demo source, and two or more frames
  share one Code panel (`TestSingleExampleOwnsItsCode`,
  `TestComponentExampleFrames`).
- The Usage bar and the header Copy import action derive their block from the
  packages the page's snippet names instead of a fixed module-root line, so a
  page that uses `form`, `overlay`, `theme` or `lucide` no longer copies an
  import it cannot build with (`TestUsageAndHeaderCopyDerivedImport`).
- Hover, pressed and focus faces snapshot state during the build; pointer-move
  tests cover the families above, and no paint callback reads an Element from
  an earlier pass (`TestButtonPaintsThroughExitTransition`).

## Platform limits

The native Gallery window was inspected through screenshots and accessibility
data, and its Toast triggers were exercised. Isolated headless render captures
also inspected Toast stacking, floating close layout, dialogs and Drawer gutters.
No personal application configuration was read or changed. These checks do not
establish screen-reader operation, platform IME/gesture behavior, OS menu
appearance or GPU partial-presentation performance; those remain separate
platform checks.

Documented approximations are listed in `docs/components.md`: no pressed scale
on buttons, no swipe dismissal for drawers/bottom sheets, surface-colored
ScrollShadow instead of an alpha mask, and the native list/menu geometry owned
by MyGo. Table selection is tinted in the cells while the toolkit owns the
selection semantics.

MyGo 0.6.3 re-lays absolute overlays after element transitions. Toast avoids
the previous vertical-motion and stack-reflow snap by spring-animating its
declared `Top` and `Height` through `AnimateTo`; scale and fade use native
motion/transitions, and reduced motion settles immediately. Modal, Drawer,
BottomSheet, Popover and Tooltip use native fades without changing their layout
or hit area. Installed overlay constructors keep the native overlay layer alive
while closed so first-open entry and last-close exit can run without a Toast
viewport. Native exit snapshots are inert and never rebuild consumer content.
