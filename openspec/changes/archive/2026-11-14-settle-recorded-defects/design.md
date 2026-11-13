## Context

Sources: the gap list at the end of each `docs/behavior/*.md` file, and the
table of findings F1–F17 in
`openspec/changes/archive/2026-11-13-extract-design-language/design.md`.

## Decisions

- D1. The reference states the expected behavior; the code moves to it. A fix
  lands with a test cited by the row.
- D2. For a design finding, the treatment used by the majority of sites or
  required by a rule in `DESIGN.md` wins; a finding with no majority is
  decided once, recorded in `DESIGN.md`, and expressed as a theme token when
  a design language would change it as one decision.
- D3. The default theme may change visibly where a finding is settled; each
  such change is listed in `CHANGELOG.md`.
- D4. Toolkit-owned limits are recorded, not simulated.

## Settled items

Each owner appends its section below: item, decision, files, test.

## Form

- NumberField commits instead of clamping each key. Typing reports a finite
  number within the bounds and keeps the text; Enter and blur clamp, round
  and format, or revert unparsable text. Steps start from the typed number.
  PageUp and PageDown step ten steps. The steppers are no Tab stops, keep the
  focus in the editor, step on press and repeat as AppKit's stepper (after
  400 ms, every 80 ms). Home and End keep moving the caret; the wheel and
  scrubbing are not adopted, as desktop number fields do not take them.
  `form/number_field.go`; `TestNumberFieldClampsOnCommitNotPerKey`,
  `TestNumberFieldSteppersKeepFocusInEditor`,
  `TestNumberFieldPageKeysStepTenSteps`,
  `TestNumberFieldStepperRepeatsWhileHeld`.
- Select: Escape returns the focus to the trigger. Closed-trigger typeahead
  cycles on a repeated letter and matches a longer prefix within one second,
  the toolkit list's pause. `form/select.go`;
  `TestSelectEscapeReturnsFocusToTrigger`,
  `TestSelectClosedTypeaheadCyclesAndMatchesPrefix`.
- KeyValueEditor: a removal focuses the neighbouring remove button, else the
  add button. `form/key_value_editor.go`;
  `TestKeyValueEditorRemoveMovesFocusToNeighbour`.
- SearchField, ComboBox, Autocomplete: Escape is registered as a shortcut on
  the editor while it has work to do, so it runs before an ancestor's.
  SearchField's hidden clear button is `Invisible`. ComboBox and Autocomplete
  clear a Clearable field on Escape with the popup closed and take Alt+Down
  and Alt+Up. `form/search_field.go`, `form/combo_box.go`,
  `form/autocomplete.go`; `TestSearchFieldEscapeClearsBeforeAncestor`,
  `TestSearchFieldHiddenClearIsNoTabStop`,
  `TestComboBoxAltArrowsOpenAndClose`,
  `TestComboBoxEscapeOnClosedPopupClears`,
  `TestAutocompleteAltArrowsAndEscapeOnClosedPopup`.
- Accessibility: tags are list items with their position; the hotkey
  recorder is a toggle button pressed while recording; choice cards expose
  Invalid; upload rows are a list with named status icons. New
  `locale.Messages` fields: `VerificationCode`, `Submit`, `Uploaded`,
  `UploadFailed`. `TestTagsInputTagsAreItemsOfASet`,
  `TestHotkeyInputExposesRecording`, `TestChoiceCardsExposeInvalid`,
  `TestFallbackNamesFollowMessages`, `TestUploadFilesAreAListWithStatus`.
- Rate: PageUp and PageDown move a whole star
  (`TestRatePageKeysMoveAWholeStar`). Slider already takes them from the
  toolkit (`TestSliderPageKeysStepATenth`). RangeSlider and the NumberField
  spin button role are recorded as toolkit limits; the combo box active
  descendant is settled under "Across packages".
- F2: the Switch track fades over `Motion.State` (150 ms), with its thumb.
  `form/switch.go`; `TestSwitchTrackSettlesWithTheThumb`.
- F5: one press treatment for thumbs: 1.2 times over `Motion.State`
  ease-out, without a delay. A slider thumb grows in both axes; the switch
  thumb, which fills its track's height, widens. The factor is the theme
  token `Effects.ThumbPressScale`, which both read, as does the colour area,
  because three components share it and a design language would change it
  as one decision.
  `form/switch.go`, `form/slider.go`;
  `TestSwitchHeldThumbGrowsAsASliderThumb`.
- Waste: Rate kept its star paths and fills across frames; a frame with a
  five-star rating went from 48 to 17 allocations (105 to 100 KB) in a
  headless benchmark.

## Date, color, chart, content, layout, effect

- F3, swatch hover: the palette swatch keeps its 1.1× hover scale, shown at
  once. A swatch shows the colour being chosen, so the fill or tint other
  hovers use would change the very thing it offers; growing leaves the
  colour as it is. The chosen swatch's 0.77 inset still moves over
  `Motion.State`. `color/color_swatch.go`; `TestSwatchHoverGrowsAtOnce`.
- F5, colour-area thumb: held, it grows by `Effects.ThumbPressScale` (1.2,
  was 1.25), as a slider thumb. `color/color_area.go`;
  `TestColorAreaThumbPressed`.
- F7, calendar view swaps: the year grid, the day grid, the header arrows
  and the chevron all take `Motion.State` with `Motion.Ease` (were 200, 120
  and 150 ms). `date/calendar.go`; `TestCalendarViewsSwapOverTheStateDuration`.
- F17, loading chart: busy, and focusable while it has State, so the
  keyboard user keeps the chart's place; no inspection until data arrives.
  `chart/chart.go`; `TestLoadingChartIsBusyAndKeepsItsTabStop`.
- Date: day grids are grids of cells; arrows stop at Min and Max and mirror
  right to left; Home and End stay in the week; Shift+PageUp/PageDown move a
  year; Escape restores the year grid's year and drops a range's pending
  endpoint; a range previews under the pointer; triggers report expanded;
  TimeField gains Min, Max and Step. `date/*.go`; the eleven tests in
  `date/behavior_gaps_test.go`.
- Color: ColorArea pages brightness and mirrors; a horizontal ColorSlider
  takes Up and Down; a selected swatch is described Selected;
  ColorSwatchPicker works with a nil Value; ColorPicker focuses its first
  control and gives the focus back to its trigger itself, because MyGo
  records the autofocused control, not the opener, as the popover's focus
  return. ColorArea's brightness stays value text only (toolkit limit: one
  range per element). `color/*.go`; `color/behavior_gaps_test.go`.
- Chart, decisions: the plot is not mirrored right to left, so the arrows
  follow the drawing; one X's readout names every series, so there is no
  step between series; a heatmap stays one group whose readout names the
  column and row, not a grid of elements; Sankey's order is settled by flow
  keys (Right follows a node's first outgoing link to its target, Left goes
  back); Sunburst reuses `navigation.Breadcrumbs`. Fixes: Escape passes
  with nothing selected; the ring shows for the keyboard only; the value
  spans the data; a pointer selection ends when the pointer leaves; Gauge is
  a meter; Choropleth's wheel, pinch and drag; locale texts.
  `chart/*.go`; `chart/behavior_gaps_test.go`.
- Chart waste: line series keep their path and buffers across frames
  (`native.LineInto`, `MinMax` into kept slices), the choropleth keeps one
  path, Sankey reuses one path, the heatmap hit test is a lookup. Headless
  frame benchmarks (`chart/paint_bench_test.go`), before → after:
  Line 123 KB/99 allocs → 8.7 KB/66; Choropleth 339 KB/364 → 33 KB/54;
  Sankey 6.4 KB/69 → 4.2 KB/51. Left: `native.Area` and `native.Sector`
  allocate per call (toolkit); projecting the map runs per build, which
  happens only on change.
- Content: the toolbar search fields are one helper (the log's copy is
  gone). Escape is a shortcut on the editor while it holds text, so it runs
  before a dialog's; Clear and Jump to latest keep the focus in the field
  and on the lines. Cmd+F and F7 stay viewer-scoped, as a window may hold
  several viewers; `FocusSearch` on the JSON and log states lets an
  application's Find reach one. Counts are statuses, the invalid-document
  error an alert. LogViewer is a group named by Label or the locale's Log;
  its lines take the focus through `ListState.Selected`, so the toolkit's
  context-menu key opens a line's menu, and rows read their level. Code
  blocks without a name are no groups; titled ones name their copy button
  (`locale.Messages.CopyNamed`). DiffViewer's header shows without a Title.
  `content/*.go`; `content/behavior_gaps_test.go`.
- Layout: the title bar is `RoleBanner`; a scaffold's start band stays a
  group, so a window has one banner. Before the first layout a splitter
  without Max takes the window's width as its room. The resizable sidebar
  passes `SidebarBorder` through `ResizablePanelOptions.Border`. A compact
  pane opens only through its caller-owned Open, as a Sheet; keyboard window
  moves are the system's. `layout/*.go`; `layout/behavior_gaps_test.go`.
- Effect: Shimmer marks its element busy; Mask clips its element, as the
  shader's layer does. Reveal leaves arriving content reachable at once:
  the wipe is decoration. Mask's use of the single `DrawOver` slot is a
  toolkit limit. `effect/*.go`; `effect/behavior_gaps_test.go`.

## Root, list and data

- F1 Spoiler: opens and closes with `Motion.Ease` (ease-out) over 200 ms;
  the toolkit's `CollapsibleBase` animates its own height with a fixed curve
  (`ui/collapsible.go:38`), a toolkit limit shared with navigation's
  Accordion and Disclosure. `spoiler.go`; read from the source [C]: no
  test asserts the curve or the duration.
- F2/F3 hover and press: `CloseButton` and `PressableFeedback` show their
  face at once, as buttons do. A `ResizablePanel` handle shows at once under
  the pointer and fades over `Motion.State` for keyboard focus. The link
  icon fade is the one root hover fade left, on the State token.
  `close_button.go`, `pressable_feedback.go`, `resizable_panel.go`,
  `link.go`; `TestCloseButtonShowsItsFaceAtOnce`, `TestLinkIconFadesOnHover`,
  `TestResizablePanelKeyboard`.
- F4 Sidebar: collapse eases out over 200 ms; a selected item keeps its
  weight, so the row keeps its size. `sidebar.go`;
  `TestSidebarCollapseMovesThroughWidths`,
  `TestSidebarSelectionKeepsTheRowSize`.
- F8 sweep: progress and the refresh line share `internal/busy.Sweep`, which
  enters and leaves wholly. `internal/busy`, `progress.go`;
  `TestProgressIndeterminateStaysInsideTrack`.
- F9 pulse: `Skeleton` and `SkeletonGroup` share `internal/fx.Pulse`, alpha
  to half with ease-in-out. `skeleton.go`, `internal/fx`;
  `TestSkeletonMotionIsOptIn`, `TestSkeletonStillsForHostPreferences`.
- F10: `ScrollShadow` and `Spoiler` fade over `Effects.FadeEdge` (24 DIP).
  `scroll_shadow.go`, `spoiler.go`; `TestScrollShadowFadesTheScrolledEdges`.
- F12/F15 Tree: the arrow turns with `Motion.Ease` (curve and duration read
  from the source [C]; `TestTreeArrowTurnsWithTheBranch` asserts only that
  it turns); typeahead waits 1 s, the toolkit list's pause, sharing
  `typePrefix` with ListBox. `list/tree.go`;
  `TestTreeTypeaheadMatchesAPrefixWithinTheWindow`.
- F16 elevation edge: the inset line follows the coss `before:` intent,
  light from above: `Y: -1` black 4% in light (bottom edge), `Y: 1` white
  6% in dark (top edge); fields follow. `theme/language.go`;
  `TestElevationEdgeTakesTheLightFromAbove` (failed before: light top 123
  vs 120, dark bottom 135).
- `Effects.ThumbPressScale` (1.2) is the size a held thumb grows to, for
  form's Switch and Slider to read. `theme/language.go`.
- Root assistive technology: status marks are named images, titles are
  level-3 headings, the sidebar is a navigation landmark with a current
  page, danger messages are assertive alerts, `ListGroup` is a list, the
  avatar overflow count is named, `SkeletonGroup` is busy, `ButtonGroup`
  joins a toolbar, `LoadingOverlay` gives the focus back when it hides,
  `ScrollShadow` panics on a nil state, `ResizablePanel` pages by 64 and
  ignores End without Max. Tests are marked in `docs/behavior/root.md`.
  Toolkit limits: focus after a Banner or GuideCard dismissal
  (`ui/input.go` moveFocus; the caller moves it with a `ui.Handle`), and
  Highlight has no mark role (`ui/access.go`).
- ListBox: listbox and option roles, multiselectable, active descendant (so
  ComboBox and Autocomplete can point at it), prefix typeahead with digits,
  Shift+arrows and Cmd+A in multiple mode, PageUp/PageDown move off darwin
  and scroll on darwin. `list/list_box.go`; `TestListBoxIsAListBoxOfOptions`,
  `TestListBoxTypeaheadBuildsAPrefix`,
  `TestListBoxExtendsAMultipleSelectionFromTheKeyboard`,
  `TestListBoxPagesFollowThePlatform`.
- Tree: multiselectable, `TreeOptions.Disabled`, the menu key and Shift+F10
  open the keyboard row's menu, paging off darwin. `list/tree.go`;
  `TestTreeMarksAMultipleSelectionAndCanBeDisabled`,
  `TestTreeMenuKeyOpensTheKeyboardRowsMenu`, `TestTreePagesFollowThePlatform`.
- SortableList: Escape cancels a keyboard or pointer drag; rows report
  position in set; grips declare Alt+Up/Down; a move is announced with
  `locale.Messages.MovedTo`. `list/sortable_list.go`;
  `TestSortableListEscapeCancelsADrag`,
  `TestSortableListTellsAssistiveTechnology`.
- Table: arrows, Home and End pass disabled rows. The native list moves the
  choice before any caller shortcut sees the key, so the table reads the
  held navigation keys (only when it has disabled rows) and moves the choice
  on in the key's direction. The table itself reports Busy while loading;
  `TableOptions.Disabled` takes it out of input. A row cannot be reported
  disabled (`ui/access.go` RoleRow, no per-row state in `ListState`): a
  toolkit limit. `data/table.go`, `data/state.go`;
  `TestTableArrowsPassADisabledRow`,
  `TestTableHomeStopsOnTheNearestEnabledRow`,
  `TestTableReportsBusyAndCanBeDisabled`.
- List typeahead limits: non-ASCII letters (`ui/input.go` typedRune) and
  `*` to expand siblings (`ui/keys.go` has no such key) are toolkit limits.

## Overlay and navigation

Findings:

- F6: one token pair per kind of modal surface, and the scrim follows its
  surface. Dialogs and sheets (Dialog, Sheet, AlertDialog, Modal,
  ImageViewer) move and fade over `Motion.Dialog` with `EaseTravel`; drawers
  and BottomSheet over `Motion.Drawer` with `EaseDrawer`. Modal, ImageViewer
  and a side Dialog without slide drop their ease-out. `overlay/motion.go`,
  `modal.go`, `image_viewer.go`, `dialog.go`, `drawer.go`;
  `TestScrimsOfOneKindDarkenAlike`.
- F10: the DocumentTabs strip already fades over `Effects.FadeEdge`
  (`navigation/document_tabs.go`); no change.
- F13: `DESIGN.md` records PipsPager as it is; confirmed, no change.
- F15 (menus): keep 1 s; the window is `menuTypeWindow`
  (`navigation/menu_ahead.go`).

Decisions:

- Menu buttons (Dropdown, SplitButton, Menubar) give the focus back to the
  button that opened their menu after a choice or Escape, never after an
  outside press; Down and Up open the menu at its first and last item; Tab
  out of a Menubar menu closes it. The menu panel is keyed by its menu so
  the toolkit's opener record survives rebuilds. `navigation/menu.go`,
  `menubar.go`, `split_button.go`; `TestMenuButtonsTakeTheFocusBackAsTheirMenuCloses`,
  `TestMenuButtonsOpenFromTheArrows`, `TestMenubarTabClosesTheOpenMenu`,
  `TestSplitButtonEscapeReturnsFocusToTheArrow`.
- One roving helper (`rove` in `navigation/state.go`) gives Accordion
  headers, DisclosureGroup triggers and pressable Steps arrows, Home and End
  while each stays a Tab stop, so Tab reaches what lies between them.
  `TestHeaderArrowsLeaveTabToWhatLiesBetween`,
  `TestStepsArrowsMoveAlongPressableSteps`.
- Long menus build type-ahead from the item labels, not from built rows, so
  rows not yet built are reachable. `navigation/menu_ahead.go`,
  `menu_list.go`; `TestLongMenuTypeAheadReachesRowsNotYetBuilt`.
- Current place: Pagination's current page and Breadcrumbs' last crumb are
  marked current page instead of a "Selected" suffix; a Previous or Next
  that disables while focused hands the focus to the current page; Carousel
  hands it to the slides and announces the page reached. Tabs' panel is a
  tab panel controlled by its tab. Toolbar exposes its orientation and its
  separators. Command's input is the one Tab stop, controls the list, and
  Escape empties a query first. `TestPaginationHandsTheFocusToTheCurrentPageAtAnEnd`,
  `TestBreadcrumbsFollowLinksFromTheKeyboard`,
  `TestCarouselKeepsTheFocusAndAnnouncesThePageAStepReaches`,
  `TestTabsPanelIsATabPanelOfTheSelectedTab`,
  `TestToolbarExposesItsOrientationAndSeparators`,
  `TestCommandInputIsTheOneTabStopAndEscapeClearsTheQuery`.
- Non-modal surfaces take no focus as they open, so the toolkit may hold no
  opener for them. One `focusReturn` (`overlay/focus_return.go`) records
  whether the focus is inside while open and hands it on once on close:
  TeachingTip to its anchor, FloatingPanel and FloatingSheet to the new
  `ReturnFocus` (the panel falls back to its opener), ActionBar to its
  existing `ReturnFocus`. `TestNonModalSurfacesReturnTheFocusAsTheyClose`.
- Reduced motion keeps an instant transition (`anchored.Still`, a negative
  Duration) instead of dropping it. MyGo treats an element that drops its
  transition while shown as gone (`ui/transition.go`, animateExits) and
  paints an inert copy over the live one that takes its focus; a negative
  Duration shows every change, Enter and Exit at once and keeps the element
  (`transition.animation`, `exitsAtOnce`). ActionBar now appears and leaves
  at once. `internal/anchored/surface.go`, `overlay/motion.go`,
  `action_bar.go`; `TestActionBarAppearsAndLeavesAtOnceUnderReducedMotion`,
  and `TestActionBarStaysWholeAndStillAsReducedMotionTurnsOn`, which fails
  when Still is dropped.
- AlertDialog: Escape cancels and runs `OnCancel`; the backdrop stays inert.
  `TestAlertDialogEscapeCancels`.
- Popover takes the focus into its panel on open and closes when the focus
  leaves it; PreviewCard, opened by hover, does not. `TestPopoverTakesTheFocusAndClosesAsTheFocusMovesOn`.
- ImageViewer: in right-to-left, Left steps forward and Previous and Next
  change sides. `TestImageViewerStepsTheWayImagesRunInRightToLeft`.
- FloatingPanel: Alt and an arrow move it 24 DIPs within the window.
  FloatingSheet: the handle reports its range and value when `MaxHeight` is
  finite, and Home and End reach the limits. `floatingStep` names the shared
  24 DIPs. `TestFloatingPanelMovesByKeyWithinTheWindow`,
  `TestFloatingSheetSplitterExposesItsRangeAndReachesItsEnds`.
- Toast: F6 from outside moves the focus to the newest toast's close button;
  closing the toast that holds the focus hands it to the newest one left.
  `TestToastRegionTakesTheFocusAndKeepsItAsToastsClose`.

Toolkit limit: when the last toast closes with the focus inside, nothing
holds the focus. MyGo owns the toast lifetime and gives no handle to the
element focused before (`ui.Context` resolves known handles only).

Needed other packages, settled under "Across packages": Carousel role
descriptions and Command's active descendant.

Waste removed: a second dialog motion path duplicating `surfaceMotion`, and
Dropdown's duplicate popover branch.

API, additive: `FloatingPanelOptions.ReturnFocus`,
`FloatingSheetOptions.ReturnFocus`.

## Across packages

- Active descendant: `list.ListBoxOptions.Editor` is a `*ui.Handle` bound
  to the element that keeps the focus while it drives the list; the
  highlighted row becomes its active descendant as well as the list's
  (MyGo's `Element.ActiveDescendant`, on the element `Context.Resolve` gives
  for this pass, after the rows are built). A handle, not a `ui.Element`,
  because `mygo vet` rejects a build value stored in a struct. ComboBox,
  Autocomplete and Command pass their input's handle. `list/list_box.go`, `form/combo_box.go`,
  `form/autocomplete.go`, `navigation/command.go`;
  `TestTextPopupsPointTheirInputAtTheHighlightedOption`,
  `TestCommandInputPointsAtTheHighlightedRow`.
- Carousel: the viewport's role reads as `locale.Messages.Carousel`
  ("carousel") and each slide's as `Slide` ("slide"), through
  `RoleDescription`. `navigation/carousel.go`;
  `TestCarouselTellsItsRolesInTheLocale`.
- Breadcrumbs: without Label the group is named by
  `locale.Messages.Breadcrumbs`, as Sidebar takes `Sidebar`; Sunburst's path
  takes the same name. `navigation/breadcrumbs.go`;
  `TestUnnamedBreadcrumbsTakeTheLocalesName`.

## Review

An independent review against v0.2.0 found nine issues; each fix has a test
that failed before it.

- NumberField committed stale text on blur after the application set Value
  while the editor had the focus. A changed Value now replaces the text
  before the focus-change branch, focused or not.
  `TestNumberFieldAdoptsProgrammaticValueWhileFocused`.
- Popover opened while the user typed in its text field anchor took the
  focus. It no longer auto-focuses when the anchor holds the focus and the
  focus takes text (a button anchor still hands it to the panel); MyGo's
  `FocusLeft` already counts the anchor as inside.
  `TestPopoverOpenedWhileTypingInItsAnchorLeavesTheFocus`.
- Popover and Popconfirm closed when the window was deactivated: the focus
  leaving is ignored while `Context.WindowActive` is false.
  `TestPopoverAndPopconfirmSurviveWindowDeactivation`.
- Table Shift+arrow toward a disabled end collapsed the range. With Shift
  held, passing a disabled row adds to the range, and with no enabled row
  further the snapshot is restored. `TestTableShiftArrowKeepsTheRangeAtADisabledEnd`.
- `fadeTransition` and `appearTransition` allocated their zero motions per
  call (3 allocs per pair, 0 after); they share package values and still
  read the theme's tokens. Tree bound its keys once per build already; its
  36 type-ahead closures are now made once per state (86 to 50 allocations
  per build of the test tree; v0.2.0 measured 76).
- `anchored.Still` uses a negative Duration: changes show at once and the
  element keeps its transition, so it is not ghost-copied, and a theme
  change no longer fades it. `TestActionBarStaysWholeAndStillAsReducedMotionTurnsOn`
  fails both with the old Still and with Still dropped.
- F1 and F12 curve and duration evidence is now [C].
- CHANGELOG: NumberField text and commit, mirrored edge slides, the true
  elevation edge users, popover focus rules.
- Heatmap's single group is its stated behavior, not a gap.
