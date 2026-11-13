# Changelog

## 0.2.0 - 2026-11-13

### Added

- Components absorbed from other libraries, with caller-owned state and no
  I/O; capabilities are in [docs/components.md](docs/components.md).
  - `mygoui`: `NumberFlip` (rolls changed digits; `Statistic` uses it),
    `Statistic`, `Descriptions`, `Timeline`, `StatusDot`, `Highlight`,
    `Spoiler`, `LoadingOverlay`, `Banner`, `SettingsGroup`.
  - `form`: `PasswordField`, `TagsInput`, `RangeSlider`, `CheckboxCard`,
    `RadioCard`, `PathField`, `KeyValueEditor`, `MaskInput`.
  - `list`: `Tree` with single and multiple selection, checks, type-ahead,
    lazy and loading nodes, guides and a row context menu.
  - `navigation`: `Menubar`, `SplitButton`, `OverflowList`, `DocumentTabs`,
    `NavigationStack`, `Carousel`, `PipsPager`, `StatusBar`.
  - `overlay`: `Popconfirm`, `TeachingTip`, `Tour`, `ShortcutsDialog`,
    `ActionBar`, `AboutDialog`.
  - `content`: `DiffViewer`, `JSONViewer`, `LogViewer`, each with a toolbar
    for search and copy; the log adds level filter and wrapping, the diff
    change navigation and a layout toggle.
- `Loading` on Button, ToggleButton, SplitButton, Switch, Input, TextField,
  SearchField, Select, ComboBox, Autocomplete, `data.Table`, `list.ListBox`,
  `list.Tree`, `navigation.Command`, Descriptions, Statistic and
  SettingsGroup. Controls keep their size and swap to a spinner, fields stay
  editable, and collections show skeleton rows or dim and show a top line
  while refreshing.
- State-change motion in shared durations (150, 200 and 300 ms): checkbox,
  radio, switch, chip, tag, toggle, sidebar item, calendar day, steps, menu
  and list checks, rating stars, image select, upload zone, OTP slot and
  colour swatch fade or draw in instead of switching at once.
- `navigation.MenuItemOptions.Accelerator`: shown in the platform's glyphs,
  and registered for the window by a Menubar. Short menus type ahead.
- Gallery groups "Viewers" and "Patterns".
- `navigation.MenuItemOptions.Checkable`: a check box row that is off while
  `Checked` is false. `locale` messages `CloseNamed` and `RemoveNamed`.

### Changed

- A pending Button no longer looks disabled: its content crossfades to a
  centred spinner, it keeps its face and takes no press.
- A range Slider paints its own track and thumbs, so the native ones no
  longer show around them. Slider fills ease to values the application sets
  and still follow a held pointer.
- Progress and ProgressCircle fills fade between status colours.
- Popover arrows carry the panel's hairline, so they stay visible on a light
  surface.
- ListBox selection checks scale in over 150 ms (was 200 ms) and still fade
  under reduced motion.
- The ImageViewer glides to a zoom chosen by key or button and fades the next
  image in over the one leaving; an Avatar portrait fades in over its
  initials.
- Gallery page column is centred between the sidebar and the page index.

## 0.1.0 - 2026-11-13

First tagged release. The `coss` reconstruction replaces the previous
unreleased appearance and token contract with coss ui references and native
Go APIs. Existing extra families remain library extensions; the catalog is a
capability mapping, not a claim of pixel-perfect completion.

### Changed

- ComboBox, Autocomplete and Command keep their filtered rows between
  frames instead of refiltering every option each frame. They refilter when
  the query, the option slices or the new `ComboBoxState.Revision`,
  `AutocompleteState.Revision` and `CommandState.Revision` change: bump
  Revision after editing options in place, as with `data.TableState`.
- Select no longer walks its options each frame while nothing is chosen and
  builds its typeahead handlers once. A row without a key is never the
  chosen row.
- Menus with 32 or more rows build only the rows in view, at a fixed width
  (320 DIP, submenus 360), with the same keyboard movement and positions for
  assistive technology.
- Dialog, AlertDialog, Drawer, Sheet and Modal round their scroll viewport
  instead of clipping the panel, so a sliding or fading panel no longer
  rasterizes a clip mask on the CPU. In a browser a Drawer and an
  AlertDialog enter and exit with no frame over 16 ms.
- Without a GPU, SkeletonGroup rests static and a mask ScrollFade shows the
  shadow treatment instead of interpreting its shader per pixel.
- The browser gallery leaves the browser its own shortcuts, such as reload
  and find: MyGo's web surface cancels every key pressed with Ctrl or ⌘, so
  the gallery keeps all but its own (editing and ⌘K search) from it.
- The gallery sidebar builds the rows near its viewport: a quarter of the
  memory and half the allocations per frame.
- Frame cost, from a review of every component. Drawer, BottomSheet, inset
  Sheet and PressableFeedback move one cached bitmap while they slide or
  scale instead of rasterizing their rounded clip on the CPU each frame.
  A Table re-sort resolves each row's key once (100,000 rows: 86 ms to
  16 ms) and a multiple selection no longer walks every disabled row each
  frame. Stacked areas thin like unstacked ones, ProfitLoss strokes two
  paths, series sharing their X values skip the per-frame sort, Choropleth
  projects once per build and hit-tests by bounding box without
  allocating, and Radar and Sunburst keep their shaped labels. Status
  inks are resolved once per palette; Skeleton and ProgressCircle stop
  copying the palette per instance; a single Select stops at its chosen
  row.
- The gallery's example sources are generated by `make generate` instead
  of parsed in the browser on the first component page.
- `make race` and the CI race job pass `-short`, which skips the
  rendering sweeps: about one minute instead of seventeen.
- Dialog and AlertDialog enter with their fade and scale in a browser. MyGo
  rasterizes a scaled rounded clip on the CPU each frame, which took longer
  than the 200 ms motion in WebAssembly, so the panel appeared late and at
  once; the panel now scales as one cached bitmap while it enters.
- Cartesian charts read the caller's points without copying them, find a
  brush range by binary search and thin dense unstacked areas as lines
  already were: three series of 5,000 points build in about half the time
  with a quarter of the allocations. Heatmap validation no longer builds a
  map per frame.
- Gallery component pages are one document in the coss docs order: lead
  example, Installation, Usage, API Reference, Examples and Theme tokens,
  with an "On this page" index on wide windows. Each example shows its own
  source behind a Code tab, and the reference lists every field of the
  option types from the library source (`make generate`; a test fails when
  it is stale). Listings share one surface with a corner copy control and
  colour Go keywords, strings, numbers, comments, calls and qualifiers.
- Gallery home page: a full-window landing page without the sidebar, with
  a centered hero, a stage of live components, catalog
  numbers, what the library is built on, families and downloads.
- Motion follows coss more closely: a held switch thumb stretches after
  100 ms, Select and ComboBox show a static up-down chevron,
  preview cards shrink as they fade out, Sidebar collapses over 200 ms, the
  tabs indicator uses ease-in-out and BottomSheet the drawer curve. The
  indeterminate progress sweep redraws at most 60 times a second.
- `fonts.Family` is now `"Cal Sans (mygoui)"`, distinct from the font's own
  name, so it never resolves to an installed Cal Sans and renders with the
  same metrics on Windows as elsewhere.
- Gallery: one title bar beside a full-height sidebar, flat grouped
  navigation, an appearance popover and a compact navigation sheet, so opening
  controls no longer moves the page. Pages lead with their title and view
  switch on one edge; copy feedback keeps its width. The home view is a
  landing page with the install command, downloads, a live preview and the
  catalog. The browser build starts with the system font instead of exiting.
- ComboBox and Autocomplete no longer reopen when a dismissed popup returns
  focus to the input. Slider mark labels sit under their ticks, overlapping
  AvatarGroup members are opaque, Sidebar content fades at its scroll edges,
  date picker triggers sit closer to the field edge and the calendar's year
  toggle uses a down chevron.
- Release tags attach the gallery for macOS, Windows, Linux and the web to
  their GitHub release; `main` publishes the browser gallery to GitHub Pages.
- EmptyState action rows center recovery buttons with the icon and title,
  including narrow layouts and increased text scale.
- Popups, menus and dialogs use coss shadow-lg/5 with its appearance edge;
  tooltips use the new `theme.TooltipShadows` (shadow-md/5). Anchored popups
  share one surface. Popovers and pickers fade in from 98% and fade out,
  tooltips also shrink on exit, and menus, selects and combo boxes no longer
  animate, as in coss. Toolbars drop their overlay shadow.
- Fields show their invalid border whenever they show an error message.
  Slider values show fractional steps exactly (0.25, not 0.2). InputOTP uses
  the shared field face, including dark fill and the focused invalid border.
- New `StatusIcon` maps a status to its coss icon; danger toasts and alert
  dialogs now use the coss circle alert. Context menus show item
  descriptions and danger styling, and PopoverMenu shows shortcut strings.
- Upgrade MyGo to v0.9.0. The year picker keeps the selected year in view
  when focus now scrolls into view, and a dialog reopened during its exit no
  longer jumps for one frame. Development tools move to OpenSpec 2.0.0 and
  Impeccable 4.8.2.
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

- Sixteen components modelled on Lobe MyGo's capabilities and rebuilt with
  coss anatomy and this theme: a new `content` package (`CopyButton`,
  `Snippet`, `Highlighter`, `Markdown`, `QRCode`), `form.EditableText`,
  `HotkeyInput`, `Upload` with `PickFiles`, `Rate` and `ImageSelect`,
  `navigation.Steps`, `list.SortableList`, `overlay.FloatingPanel`,
  `FloatingSheet` and `ImageViewer`, and root `Image`, `Result`, `GuideCard`
  and `ResizablePanel`. `content.QRCode` adds the `rsc.io/qr` dependency for
  importers of `content`.
- Optional `Flyout` host on Popover, Dropdown, Select, DatePicker,
  DateRangePicker and ColorPicker opens their popup as a fitted native
  flyout that can extend beyond the window; without one, and on the web,
  they stay in the window. Popover's `Header`, `Body` and `Footer` now
  receive the context they build in.
- Native `chart` family covering the 19 documented Bklit chart subjects with
  coss chart roles, caller-owned inspection, legends, brush, axes and annotations.
  Gallery adds Charts and Chart Utilities with public Go examples. The source
  reference and MIT notice are pinned; `docs/charts.md` records native subsets
  and missing upstream variants. Proprietary Studio is not reused.

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
