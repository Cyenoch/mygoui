# Component capability matrix

The additional [chart matrix](charts.md) maps all 19 Bklit chart subjects and
the utility/axis catalog to native APIs. Bklit supplies chart capabilities;
the visual system remains coss. The matrix states missing variants explicitly,
not React feature parity.

This maps the coss ui catalog to the native Go APIs observed in this checkout.
The reference is commit `3d8ca28e7d90123e9e0f44665944dc240b026ccb`;
see [sources and licensing](references.md).

The downloaded site links and component metadata list 56 component pages.
`llms.txt` lists 55: it adds Sidebar and omits Composer and Segmented Control.
Their union is **57 component subjects**. The two documented hooks below make
59 upstream subjects, not 59 fully implemented or pixel-perfect components.
ScrollFade is an additional native capability.

Every row states a Go mapping and its current scope. An existing API has not
thereby passed behavioral, visual or platform validation. Explicit partial or
pending rows remain gaps; [VALIDATION.md](../VALIDATION.md) records evidence.
Names link to upstream documentation. Source paths after each API are relative
to this repository; upstream implementations normally live at
`apps/ui/registry/default/ui/<slug>.tsx`, with MDX at
`apps/ui/content/docs/components/<slug>.mdx`. Date Picker is a documented
composition rather than a standalone registry UI file. Sidebar has a registry
file even though its page is missing from the downloaded link list. On
2026-10-27 the live Sidebar component URL returned 404; visual review uses the
actual docs sidebar and that registry source instead. See the per-component
[visual review](visual-review.md) for inspected appearances and limitations.

## Components

| coss reference | Go mapping and source | Native capability and limits |
| --- | --- | --- |
| [Accordion](https://coss.com/ui/docs/components/accordion) | `navigation.Accordion` - `navigation/accordion.go` | Typed sections, single/multiple expansion and keyboard movement. Native layout/motion replaces CSS state selectors. |
| [Alert](https://coss.com/ui/docs/components/alert) | `mygoui.Alert` - `alert.go` | Callout with icon/title/description/action and status colors. Compose additional content with native controls. |
| [Alert Dialog](https://coss.com/ui/docs/components/alert-dialog) | `overlay.AlertDialog` - `overlay/alert_dialog.go` | Explicit action/cancel modal with native focus ownership. Backdrop/Escape do not choose the action. |
| [Autocomplete](https://coss.com/ui/docs/components/autocomplete) | `form.Autocomplete / AutocompleteState` - `form/autocomplete.go` | Inline native editable input with filtered virtualized suggestions; editor keeps focus during arrows/Enter selection. Single/multiple keyed selection, clear/loading and query callbacks; the application owns async requests. SearchLabel/SearchPlaceholder are editor fallbacks, not a second input. |
| [Avatar](https://coss.com/ui/docs/components/avatar) | `mygoui.Avatar` - `avatar.go` | Caller-provided bitmap/SVG or initials, sizes and optional badge. `ui.DecodeImage` accepts loaded bitmap/SVG bytes; nil/empty sources keep initials. Image fetching/loading policy belongs to the application. |
| [Badge](https://coss.com/ui/docs/components/badge) | `mygoui.Badge` - `chip.go` | Independent BadgeVariant: BadgePrimary (default), BadgeSecondary, BadgeOutline, BadgeDestructive, BadgeError, BadgeInfo, BadgeSuccess, BadgeWarning. Options: Label, Variant, Size, Icon. |
| [Breadcrumb](https://coss.com/ui/docs/components/breadcrumb) | `navigation.Breadcrumbs / Breadcrumb` - `navigation/breadcrumbs.go` | Typed path actions, current item and overflow. URL navigation is an application/native service action. |
| [Button](https://coss.com/ui/docs/components/button) | `mygoui.Button` - `button.go` | Primary/default, secondary, outline, ghost, link, destructive and destructive-outline mappings; loading, disabled, icons and xs through xl. Native disabled opacity is 50%, coss reference 64%; no pressed scale. |
| [Calendar](https://coss.com/ui/docs/components/calendar) | `date.Calendar / RangeCalendar` - `date/calendar.go, date/range_calendar.go` | Single day or start/end range, Gregorian month navigation, bounds and disabled days. Calendar's month heading is static by default; ShowYearPicker enables the year selector and HideYearPicker overrides it. No arbitrary multiple-date selection, RTL date ordering or alternate calendar system. |
| [Card](https://coss.com/ui/docs/components/card) | `mygoui.Card and CardHeader/Title/Description/Content/Footer; CardFrame and CardFrameHeader/Title/Description/Footer` - `card.go` | Composable bordered card and muted card-frame helpers. Sections own their padding. Additional Surface variants are library extensions, not additional coss contracts. |
| [Checkbox](https://coss.com/ui/docs/components/checkbox) | `form.Checkbox` - `form/checkbox.go` | Native checkbox with labels, disabled/invalid and mixed state. MyGo owns activation, focus and semantics. |
| [Checkbox Group](https://coss.com/ui/docs/components/checkbox-group) | `form.CheckboxGroup / CheckboxGroupState` - `form/checkbox_group.go` | Caller-owned grouped selection and native checkboxes in a body callback; explicit group validation. |
| [Collapsible](https://coss.com/ui/docs/components/collapsible) | `navigation.Disclosure` - `navigation/disclosure.go` | DisclosureOptions supplies caller-owned expansion and body; DisclosureGroup coordinates sections. No React Collapsible parts/provider API. |
| [Combobox](https://coss.com/ui/docs/components/combobox) | `form.ComboBox / ComboBoxState` - `form/combo_box.go` | Editable filtering, single/multiple selection, removable selected items, custom values and custom filter. Async queries are callbacks; application owns requests/results. |
| [Composer](https://coss.com/ui/docs/components/composer) | `form.Composer` - `form/composer.go` | Plain-text multiline draft with persistent actions, Enter submit/Shift+Enter newline and caller-controlled clearing. No rich structured draft, mention/token/slash popup or adaptive inline action layout. |
| [Command](https://coss.com/ui/docs/components/command) | `navigation.Command / CommandState` - `navigation/command.go` | Inline searchable action list with label/keyword matching, named groups, disabled items, shared pointer/keyboard highlight, arrows/Enter and callbacks. Compose inside Dialog for a palette. Shortcut strings are display only; global shortcut registration belongs to the application. |
| [Context Menu](https://coss.com/ui/docs/components/context-menu) | `navigation.ContextMenu / PopoverMenu(anchor, MenuOptions)` - `navigation/context_menu.go` | OS menu by default; PopoverMenu opts into MyGo's current-theme in-window menu with focus restoration, type-ahead and submenu navigation. Both support secondary click/menu-key/Shift+F10, actions, radio/checkbox selection, disabled items and separators. Section headings, item icons, descriptions, danger styling and shortcut strings are not rendered on these paths. |
| [Date Picker](https://coss.com/ui/docs/components/date-picker) | `date.DatePicker / DateRangePicker` - `date/picker.go, date/range_picker.go` | Native ISO date editors plus anchored single/range calendar. Strict Gregorian date validation; no localized free-form parsing or timezone policy. |
| [Dialog](https://coss.com/ui/docs/components/dialog) | `overlay.Dialog` - `overlay/dialog.go` | Caller-owned open flag with native modal focus, header/body/footer and dismissal. The default backdrop dims black at 32% and blurs with 4 DIP sigma; explicit dim-only and transparent modes are available. Keep constructor installed for retained exit paint. Footer is full-bleed; build actions directly and choose FullWidth explicitly on narrow layouts. Short panels scroll chrome as needed to keep actions reachable. |
| [Drawer](https://coss.com/ui/docs/components/drawer) | `overlay.Drawer` - `overlay/drawer.go` | Four-edge attached modal with 450 ms native slide, separately fading scrim and optional wider-window Inset. Drag the middle of the inward-facing edge outward to dismiss; native body scrolling remains independent. Handle is optional, ShowClose adds an explicit close button, and PreventDismiss blocks swipe/backdrop/Escape dismissal. No snap points or nested drawer gesture coordination. BottomSheet is a separate floating handle-swipe API extension, demonstrated within the Drawer Gallery page. |
| [Empty](https://coss.com/ui/docs/components/empty) | `mygoui.EmptyState` - `empty_state.go` | Empty result with icon, title/description and explicit recovery action; compose richer content in native containers. |
| [Field](https://coss.com/ui/docs/components/field) | `form.Field / Description / FieldError` - `form/field.go` | Native label/help/error association around exactly one primary control. Application supplies validation and control disabled state; no HTML constraint-validation engine. |
| [Fieldset](https://coss.com/ui/docs/components/fieldset) | `form.Fieldset` - `form/field.go` | Legend, description, disabled group and native body layout. No DOM legend/fieldset tree is created. |
| [Form](https://coss.com/ui/docs/components/form) | `form.Form` - `form/field.go` | Adapted: vertical grouping and top-level error only. Submit actions, validation, reset, field collection and persistence belong to the caller; no FormData or automatic validation/submission. |
| [Frame](https://coss.com/ui/docs/components/frame) | `mygoui.Frame and FramePanel/Header/Title/Description/Footer` - `frame.go` | Muted outer frame with independently composable bordered panels. Options expose padding, gap and label; no compound provider is needed. |
| [Group](https://coss.com/ui/docs/components/group) | `mygoui.Group / GroupItem` - `frame.go` | Mixed controls receive explicit outer corners through Build; caller applies corners. Vertical members share intrinsic width. Group removes adjacent returned-control borders; separators are explicit GroupItem builders. No descendant inspection or GroupText API. |
| [Input](https://coss.com/ui/docs/components/input) | `form.Input` - `form/input.go` | Native bound editor, labels/help/errors, read-only, disabled and invalid state. Browser input-type/autofill contracts are not implied. |
| [Input Group](https://coss.com/ui/docs/components/input-group) | `form.InputGroup` - `form/input_group.go` | One native editor fills the remaining inline track, with Prefix/Suffix and horizontal BlockPrefix/BlockSuffix callbacks. Addon actions retain their own focus without lighting the editor shell; Error also marks the shared shell invalid. |
| [Kbd](https://coss.com/ui/docs/components/kbd) | `mygoui.Kbd` - `kbd.go` | Visual key chord/keycap with accessible text. It displays a shortcut; it does not register the command. |
| [Label](https://coss.com/ui/docs/components/label) | `form.Label` - `form/field.go` | Standalone label can focus a ui.Handle target. Prefer Field for native control/help/error association. |
| [Menu](https://coss.com/ui/docs/components/menu) | `navigation.Menu / Dropdown / SubMenu` - `navigation/menu.go` | Typed sections and item actions, checked single/multiple selection, disabled items, nested submenus and keyboard traversal. Dropdown supplies the trigger; context-menu chrome is separately platform-owned. |
| [Meter](https://coss.com/ui/docs/components/meter) | `mygoui.Meter` - `progress.go` | Determinate quantity with explicit range, label/output and indicator color. Not a form value or task lifecycle manager. |
| [Number Field](https://coss.com/ui/docs/components/number-field) | `form.NumberField / NumberFieldState` - `form/number_field.go` | Native numeric editor and steppers with bounds/step and display formatting. Parsing uses strconv.ParseFloat, independent of locale; no upstream scrub-area interaction. |
| [OTP Field](https://coss.com/ui/docs/components/otp-field) | `form.InputOTP / OTPState` - `form/input_otp.go` | One native editor with painted slots, filtering, completion and selection. Length defaults to 6 and is capped at 12; sanitizing input may move the caret. No platform OTP delivery/autofill guarantee. |
| [Pagination](https://coss.com/ui/docs/components/pagination) | `navigation.Pagination` - `navigation/pagination.go` | One-based page control, previous/next, boundaries, siblings and ellipsis. Siblings zero means no adjacent pages, not a default of one. Application owns fetching and result/page state. |
| [Popover](https://coss.com/ui/docs/components/popover) | `overlay.Popover` - `overlay/popover.go` | Same-build anchor plus caller-owned open flag, placement, optional arrow, body and native dismissal. No DOM portal or detached React handle API. |
| [Preview Card](https://coss.com/ui/docs/components/preview-card) | `overlay.PreviewCard / PreviewCardState` - `overlay/preview_card.go` | Rich anchored body opens on hover/focus and stays open over the panel. Current leave grace is 150 ms; no configurable shared delay/provider or fetched link metadata. |
| [Progress](https://coss.com/ui/docs/components/progress) | `mygoui.ProgressBar` - `progress.go` | Determinate or indeterminate task indicator with caller value/range and optional output. ProgressCircle is an extra native family. |
| [Radio Group](https://coss.com/ui/docs/components/radio-group) | `form.RadioGroup / Radio` - `form/radio_group.go, form/radio.go` | Typed mutually exclusive values with native roving focus and explicit group state. Compose choices in the body. |
| [Scroll Area](https://coss.com/ui/docs/components/scroll-area) | `mygoui.ScrollArea / ScrollAreaOptions` - `scroll_fade.go` | Axis: ScrollVertical (default), ScrollHorizontal or ScrollBoth; caller-owned ui.ScrollState. Height, Label and optional Fade. Native scrollbars replace web scrollbar parts. |
| [Select](https://coss.com/ui/docs/components/select) | `form.Select / SelectState` - `form/select.go` | Single/multiple keyed selection with sections, disabled items and clear. Native popup/list semantics and placement replace browser select/form integration. |
| [Segmented Control](https://coss.com/ui/docs/components/segmented-control) | `navigation.SegmentedControl` - `navigation/segmented_control.go` | Independent native SegmentedBase radio group with Selected, []TabOptions and OnChange. Disabled choices and arrow selection use native semantics; the selected surface moves/resizes over 200 ms and settles for reduced motion. Item Content is ignored; panels belong to Tabs. |
| [Separator](https://coss.com/ui/docs/components/separator) | `mygoui.Separator` - `separator.go` | Horizontal/vertical native rule, with additional library emphasis levels. |
| [Sheet](https://coss.com/ui/docs/components/sheet) | `overlay.Sheet` - `overlay/sheet.go` | Edge modal, right by default or left/top/bottom, with 200 ms native slide and separate scrim fade. Backdrop and optional wider-window Inset join header/body/footer, focus restoration and explicit dismissal policy. |
| [Sidebar](https://coss.com/ui/docs/components/sidebar) | `mygoui.Sidebar / SidebarOptions` - `sidebar.go` | Required caller-owned Selected and Collapsed, typed keyed groups/items, native collapse toggle, scrolling content and header/footer. Default expanded/collapsed widths are 256/48 DIP. No automatic mobile sheet, resize rail, nested menu parts or persisted collapse policy. |
| [Skeleton](https://coss.com/ui/docs/components/skeleton) | `mygoui.Skeleton / SkeletonGroup` - `skeleton.go` | Native shimmer/pulse/static placeholders with reduced-motion handling; composed dimensions replace CSS utilities. |
| [Slider](https://coss.com/ui/docs/components/slider) | `form.Slider / SliderState` - `form/slider.go` | Single-thumb or horizontal two-thumb range, bounds, steps, ticks/output and keyboard input. Vertical is single-thumb only; no arbitrary multi-thumb set. |
| [Spinner](https://coss.com/ui/docs/components/spinner) | `mygoui.Spinner` - `spinner.go` | MyGo loading plugin animation with native size/status options. Geometry is a native adaptation, not the same web SVG animation asset. |
| [Switch](https://coss.com/ui/docs/components/switch) | `form.Switch` - `form/switch.go` | Native binary switch with explicit value, label, disabled/invalid and change callback. |
| [Table](https://coss.com/ui/docs/components/table) | `data.Table / TableState` - `data/table.go` | Native virtualized data table with sort/selection/actions, column resize/hide and Footer callback. Upstream Table is a composable HTML table; native header metrics, row tinting and resize grips differ. No colspan/rowspan or expandable rows. |
| [Tabs](https://coss.com/ui/docs/components/tabs) | `navigation.Tabs` - `navigation/tabs.go` | Typed items and selected index, panels, disabled entries, horizontal/vertical layout, roving focus and selected-indicator motion. KeepMounted optionally preserves inactive native state using Hidden. No arbitrary React tab part tree. |
| [Textarea](https://coss.com/ui/docs/components/textarea) | `form.TextArea` - `form/textarea.go` | Native multiline editor with line bounds, selection, read-only/disabled/invalid and field metadata. |
| [Toast](https://coss.com/ui/docs/components/toast) | `overlay.Toast / ToastState.View` - `overlay/toast.go` | Window-owned notification identities, action, timeout, updates, native queue and animated deck. Call Toast from an action, install View in the build. MyGo caps visible cards at three; no promise helper or arbitrary browser portal. |
| [Toggle](https://coss.com/ui/docs/components/toggle) | `mygoui.ToggleButton` - `toggle.go` | Native pressed boolean with typed options and explicit change callback. Additional button appearances are library extensions. |
| [Toggle Group](https://coss.com/ui/docs/components/toggle-group) | `mygoui.ToggleButtonGroup` - `toggle.go` | Typed button options with single/multiple selection. Default items keep individual rounded corners and 2 DIP gaps; outline-only groups join, while Detached uses 4 DIP gaps. State is explicit rather than provided through React context. |
| [Toolbar](https://coss.com/ui/docs/components/toolbar) | `navigation.Toolbar / ToolbarSeparator` - `navigation/toolbar.go` | Body-based native focus group with horizontal/vertical arrows and explicit separator orientation. |
| [Tooltip](https://coss.com/ui/docs/components/tooltip) | `overlay.Tooltip` - `overlay/tooltip.go` | Native pointer-rest/keyboard-focus tooltip with placement and optional arrow. MyGo owns 600 ms rest and 400 ms warm reuse; public options do not expose upstream provider timing or shared viewport transitions. |

## Hooks and application services

`overlay.NewFlyout` is an additional native window capability, not another
coss component clone. It can extend outside the parent window, unlike the
in-window Popover. Explicit dimensions and Parent/Tray, focusability,
zoom-aware handle anchors and independently rendered coss chrome are covered
in [flyout.md](flyout.md). The host owns the window and closes/repositions it
when its anchor disappears or changes; native content has no FitContent.

| coss reference | Native mapping | Boundary |
| --- | --- | --- |
| [useMediaQuery](https://coss.com/ui/docs/hooks/use-media-query) | `c.Size()`, `c.Preferences()` and explicit native layout branches | No CSS media-query parser or React hook. Use window DIP size and available native preferences; pointer/CSS feature queries have no blanket equivalent. |
| [useCopyToClipboard](https://coss.com/ui/docs/hooks/use-copy-to-clipboard) | Capture `services := c.Services()`, call `services.WriteClipboard(text)` in an action | Native clipboard has no browser permission promise. The application owns copied-feedback state and its expiry. Do not retain a build context or perform clipboard writes during unconditional rendering. |

Persistent services also expose native URL opening. Network requests, image
loading, filesystem access, async search and persistence stay with the host.
Publish worker results through `Window.Update`.

## ScrollFade and customization

`ScrollFade(c, viewport, state, options)` decorates an existing native scroll
viewport. `MaskFade` uses a native shader effect to fade the subtree to
transparency; `ShadowFade` paints shadow ink and `SurfaceFade` blends into the
supplied opaque surface. Flags select physical edges and fades follow actual
remaining scroll distance. Native mask effect layers cap each side at 8192 pixels
and omit embedded webviews. Shadow/surface paint stays inside the viewport
clip and does not cover its border. Scoped Metal rounded-viewport evidence
is recorded in `VALIDATION.md`; arbitrary nested clips remain a toolkit limit.
See [scrollfade.md](scrollfade.md) for signatures, errors and examples.

Use [theme.md](theme.md) for canonical coss roles, per-appearance palette
patches, radius offsets and type metrics. Body/heading/mono families are
independent, but native system fonts are the default. The optional `fonts`
package embeds independently sourced Cal Sans and registers it only when
`fonts.Register()` is called. Paper Mono is not bundled. [fonts.md](fonts.md)
records OFL sources, browser limits and unresolved rendering equivalence.

## Additional library families

These remain useful native components outside the 57-component coss catalog
mapping; their presence is not an upstream fidelity guarantee:

| Family | Public surface |
| --- | --- |
| Root composition and text | `Surface`, `Typography`, `Heading`, `Paragraph`, `Code`, `Prose`, `Header`, `Link`, `LinkButton` |
| Root actions and feedback | `ButtonGroup`, `CloseButton`, `PressableFeedback`, `Chip`, `BadgeAnchor`, `AvatarGroup`, `ProgressCircle`, `ScrollShadow`, `StatusColors` |
| Root collections | `ListGroup`/`ListGroupItem`, `Tag`/`TagGroup`, `SkeletonGroup` |
| Form helpers | `TextField`, `SearchField`, `ControlField`, `SwitchGroup`, `FieldGroup`, `FieldActions`, `ErrorMessage` |
| List | `list.ListBox` with typed `ListBoxItem`/`ListBoxSection` and caller-owned `ListBoxState` |
| Overlays | `overlay.Modal`, `overlay.BottomSheet`; the latter has handle-only swipe dismissal, without coss Drawer snap points |
| Dates | `date.DateField`, `date.TimeField`, `DateRange` and calendar view state |
| Color | `ColorArea`, `ColorSlider`, `ColorField`, `ColorSwatch`, `ColorSwatchPicker`, `ColorPicker`, `ColorState`, hex parse/format helpers |

## Shared limits

- **Toolkit/platform:** OS context-menu chrome, tooltip timing, scrollbars,
  clipboard, IME and actual accessibility integration remain native. The
  presence of accessible metadata in a Tester is not screen-reader evidence.
- **Paint:** disabled nodes use native 50% opacity instead of the coss button's
  64%. Focus rings and joined-group clipping follow native paint order;
  CSS z-index and descendant selectors have no direct equivalent.
- **Dates/color:** Gregorian/ISO date behavior and sRGB color editing are the
  current contracts. Locale labels can be supplied; arbitrary calendar systems,
  locale date parsing and wide-gamut round-tripping are not supplied.
- **Motion/gestures:** reduced-motion preferences must be respected. Drawer
  snap points, nested gesture arbitration and unrestricted body swipe are
  absent. Existing slide/fade/zoom/spring transitions do not certify every coss
  interaction curve.
- **API shape:** HTML form submission, DOM refs, React providers/render props,
  arbitrary Tailwind classes and browser portals do not become native APIs.
  Native options, handles, services and composition carry the supported intent.
- **Evidence:** headless paint and native event tests must be recorded separately
  from real OS/GPU/gesture/IME/screen-reader checks. A gallery page or this
  catalog is not an acceptance result.
