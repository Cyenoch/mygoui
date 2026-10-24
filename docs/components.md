# Components

The library implements the union of HeroUI v3's catalogs as one native library.
The sections below are functional families, not web/native targets. Every
component below is constructible from Go with typed options; none requires
JavaScript, CSS, a WebView or generated assets.

Run `go run ./examples/gallery` for live specimens. The gallery is the
authoritative usage reference for prototypes and state wiring.

## Root package (`mygoui`)

| Component | Behavior |
| --- | --- |
| `Alert` | Surface callout with default/accent/info/success/warning/danger status, icon, title, description and an optional real action. |
| `Avatar` | Caller-loaded image or initials, sm/md/lg, soft fallback colors, optional badge. |
| `AvatarGroup` | Overlapping stack or grid with 8 DIP overlap, optional background ring and overflow count. |
| `Badge` | Standalone badge with color/variant/size and a background-colored ring. |
| `BadgeAnchor` | Places a badge at a content corner through the native overlay layer. |
| `Button` | primary/secondary/tertiary/outline/ghost/danger/danger-soft, sm/md/lg, icon-only, full width, pending and disabled states. |
| `ButtonGroup` | Joined buttons with outer-only corners, horizontal/vertical orientation and optional separators. |
| `Card` + `CardHeader`/`CardTitle`/`CardDescription`/`CardContent`/`CardFooter` | HeroUI's 16 DIP padded, 24 DIP rounded, 12 DIP gapped surfaces in default/secondary/tertiary/transparent variants. Native HeroUI's `Card.Body` is `CardContent`. |
| `Chip` | Compact capsule in default/accent/success/warning/danger colors and primary/secondary/soft/tertiary variants. |
| `CloseButton` | 24 DIP dismiss control with native focus and hover. |
| `EmptyState` | Low-emphasis empty result with optional icon and recovery action. |
| `Header` | Small muted section heading. |
| `Kbd` | One keycap holding a key chord with concatenated accessible text. |
| `Link` | Inline keyboard action with hover underline; URL opening stays application work. |
| `LinkButton` | Native ghost/link-style text action without a hover surface. |
| `ListGroup` + `ListGroupItem` | Native surface list with prefix/suffix content and real row activation. |
| `Meter` | Determinate quantity with label, output and size; `IndicatorColor` defaults to accent. |
| `PressableFeedback` | Native press-highlight container. |
| `ProgressBar` | Determinate or indeterminate bar with label/output and sizes; `IndicatorColor` defaults to accent. |
| `ProgressCircle` | Determinate segmented ring or indeterminate native spinner with optional output; `IndicatorColor` defaults to accent. |
| `ScrollShadow` | Native scroll viewport with surface-colored edge fades driven by actual scroll position. |
| `Separator` | Horizontal/vertical rule with default/secondary/tertiary levels. |
| `Skeleton`/`SkeletonGroup` | Shimmer, pulse or static placeholders; motion respects reduced motion. |
| `Spinner` | Native loading arc in four sizes and status colors. |
| `Surface` | default/secondary/tertiary/transparent container surface. |
| `Tag`/`TagGroup` | Selectable, removable tags with single/multiple selection and group layout. |
| `ToggleButton`/`ToggleButtonGroup` | Persistent pressed state, default/ghost, attached or detached groups. |
| `Typography` | HeroUI h1–h6, body, body-sm, body-xs and code roles; `Heading`, `Paragraph`, `Code`, `Prose` helpers. |

## `form`

| Component | Behavior |
| --- | --- |
| `Field`, `Label`, `Description`, `FieldError`, `ErrorMessage`, `Fieldset`, `FieldGroup`, `FieldActions`, `Form` | Native field association, grouping, help and error text. |
| `Input`, `TextArea`, `TextField` | Native single/multiline editing, primary/secondary variants, labels, descriptions, validation, disabled states, and values that stay selectable without changing. |
| `InputGroup` | One native editor with inline or block prefix/suffix regions. |
| `InputOTP` | Slot-guided one-time code entry over one native editor, committed ASCII filtering and completion callback. |
| `SearchField` | Native search input with icon, clear action, Escape and Enter behavior. |
| `NumberField` | Locale-independent numeric input with min/max/step, increment/decrement and formatted output. |
| `Checkbox`, `CheckboxGroup` | Native checked/mixed states with labels, descriptions, errors and orientation. |
| `Radio`, `RadioGroup` | Native single choice with roving focus and disabled options. |
| `Switch`, `SwitchGroup` | Native boolean control in sm/md/lg sizes with label placement and group orientation. |
| `ControlField` | Label/description plus Switch or Checkbox as one pressable row. |
| `Slider` | Single value or two-thumb range with min/max/step, marks, output, horizontal/vertical single-thumb support. |
| `Select` | Field trigger with virtualized single/multiple list, sections, typeahead, clear and keyboard commit. |
| `ComboBox` | Filtering text input with a virtualized list, single/multiple choice and optional custom value. |
| `Autocomplete` | Field trigger opening a search field above a filtered virtualized list, with a caller-owned loading indicator. |

## `list`

| Component | Behavior |
| --- | --- |
| `ListBox` | Native virtualized list with stable keys, hover/highlight paint, single/multiple selection, sections, disabled rows, typeahead and an empty state. `ListBoxItem` and `ListBoxSection` are typed data literals passed through options, not constructors. |

## `navigation`

| Component | Behavior |
| --- | --- |
| `Accordion` | Single or multiple disclosure over a slice of `AccordionItemOptions`, with keyboard focus movement, surface variant and separators. |
| `Disclosure`, `DisclosureGroup` | Caller-owned collapsible panels; the group takes a slice of `DisclosureOptions` and coordinates single/multiple expansion. |
| `Tabs` | Primary pill or secondary underline tabs over `TabOptions`, horizontal/vertical, disabled tabs and native roving focus. |
| `Breadcrumbs`/`Breadcrumb` | Native path actions, current-page semantics and overflow menu. |
| `Pagination` | One-based pages with previous/next, boundaries, siblings and ellipsis. |
| `Toolbar`, `ToolbarSeparator` | Native toolbar focus group with horizontal/vertical arrows and attached surface; separators take an explicit orientation. |
| `Menu` | Action menu over `MenuSectionOptions` and `MenuItemOptions`, with selection, disabled items, nested submenus and keyboard traversal. |
| `Dropdown` | Button trigger for a `Menu`, with placement and outside-dismissal. |
| `SubMenu` | Nested anchored submenu. |

## `overlay`

| Component | Behavior |
| --- | --- |
| `Modal` | Focus-contained dialog in seven sizes, including cover and extra small; opaque/blur/transparent backdrop, scrollable body, footer and close button. |
| `AlertDialog` | Explicit action/cancel dialog; backdrop and Escape do not accept. |
| `Drawer` | Floating edge panel in four directions with 16 DIP screen gutters, rounded corners and the native dialog focus contract. |
| `BottomSheet` | Floating bottom panel with 16 DIP gutters, a 250 ms slide-in, handle swipe dismissal (short drags ease back, long drags slide out) and a scrim that clears with the sheet. The close button runs the same slide-out; backdrop and Escape use the native fade. |
| `Dialog` | Native center or edge modal with header/body/footer. |
| `Popover` | Anchored rich panel with HeroUI placement, arrow, zoom entry and dismissal. |
| `Tooltip` | Native hover/focus description with placement, zoom entry and Escape dismissal. |
| `Toast` | Spring-animated deck with hover/focus expansion, swipe dismissal, variants, action, timeout and MyGo's three-visible cap. Covered cards hide their content. The close button floats across the top-right corner without reserving content width. `ToastState` holds caller-owned identities and installs the region with `View`; reusing an ID updates in place without replaying entry. |

Toast status controls use the notification title as their accessible name,
falling back to the description for a title-less notification. IDs remain
update/dismiss identities, not user-facing labels.

Modal, AlertDialog and Dialog fade in/out. Drawer fades; BottomSheet slides up
and down. Popover and Tooltip fade and zoom from the anchor-facing edge. Every
installed constructor keeps the overlay layer alive while closed, so the first
arrival and final exit animate even without a Toast region. Closed content is an
inert native snapshot, not a second live build. Reduced motion settles all
surfaces immediately.

## `date`

| Component | Behavior |
| --- | --- |
| `Calendar` | Month grid, month/year navigation, year picker, bounds, disabled dates, today marker and keyboard movement. |
| `RangeCalendar` | Start/end range selection with ordered endpoints and range paint. |
| `DateField` | ISO `YYYY-MM-DD` native editor with bounds, disabled dates, label, description and error. |
| `DatePicker` | Date field plus anchored calendar popup. |
| `DateRangePicker` | Two endpoint editors plus anchored range calendar. |
| `TimeField` | `HH:MM` or 12-hour time editor with validation. |

## `color`

| Component | Behavior |
| --- | --- |
| `ColorArea` | 2D saturation/value surface with draggable, keyboard-adjustable thumb. |
| `ColorSlider` | Hue, alpha or arbitrary channel slider with gradient track. |
| `ColorField` | Hex color input with validation. |
| `ColorSwatch` | Color preview with accessible name and optional selected mark. |
| `ColorSwatchPicker` | Keyboard-roving palette selection. |
| `ColorPicker` | Trigger and popover combining area, sliders, field and swatches over one `ColorState`. |

## `data`

| Component | Behavior |
| --- | --- |
| `Table` | Virtualized native table with HeroUI's primary/secondary chrome, 44 DIP rows, header band, column separators, sortable columns, single/multiple selection, disabled rows, row actions, column resizing/hiding, an empty slot and an in-frame `Footer` slot for pagination. |

## Deliberate limits

These boundaries are native platform decisions, not missing glue:

- **Pointer feedback:** Button, ToggleButton, CloseButton, PressableFeedback and
  NumberField steppers animate HeroUI's pressed scale. MyGo does not expose a
  live layout transform, so there is no ripple.
- **Gestures:** BottomSheet dismisses from its handle swipe; the scrollable body
  keeps its own pointer handling. Drawer does not implement swipe-to-dismiss or
  velocity arbitration; backdrop press, Escape and explicit close actions work.
- **ScrollShadow:** MyGo has no alpha mask, so the fade is painted in the
  caller-supplied opaque surface color rather than masking the content.
- **Dates:** arithmetic is Gregorian. Date fields use strict ISO; alternate
  calendars, locale-specific month-name parsing and Calendar week/day views
  are not implemented. The month grid, year picker, range selection and
  keyboard movement are.
- **Color:** values are sRGB hex; wide-gamut round-tripping and arbitrary color
  spaces beyond MyGo's `ui.Color` are not exposed.
- **Table:** no expandable/colspan rows. Column resize and hide are native;
  saved width/order must be restored through `TableState.List.Columns`.
  Selected rows paint a full-bleed tint while native selection semantics and
  keyboard behavior remain intact. The toolkit owns the header row's exact 32
  DIP height, 10/6 padding and 600 weight; the library sets the 12 DIP muted
  type and 16 DIP column insets around that chrome. Column separators are the
  native resize grip: a column shows one only when `Resizable`, and the last
  column never does.
- **Menus:** native popup placement and OS menu appearance remain
  platform-owned.
- **Tooltips:** timing follows MyGo's native hover/focus policy rather than
  CSS delays.
- **Toast:** MyGo displays at most three toasts; additional identities remain in
  `ToastState` and appear as earlier ones close.
- **OTP:** rejected characters rewrite the native editor text, which may move
  the caret to the end for that edit.
- **NumberField:** parsing is locale-independent (`strconv.ParseFloat`); the
  `Format` hook changes display only.
- **Slider:** range sliders are horizontal; vertical supports one thumb.
- **Typography/Prose:** semantic roles and inline code are native; arbitrary
  CSS prose (nested lists, blockquotes, images) is composed by the caller.
- **Web-only API shape:** render props, `asChild`, DOM polymorphism, Tailwind
  classes and React context are intentionally absent. Extend components with
  typed Go options and composition instead.
