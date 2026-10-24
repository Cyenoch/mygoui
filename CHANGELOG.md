# Changelog

## Unreleased

Complete native rewrite from the shadcn-inspired API to a HeroUI v3 component
library. The previous component set is not part of this API and no compatibility
shims are provided.

### Added

- HeroUI v3 semantic theme: light/dark palettes, accent/default/status roles,
  derived hover/soft/separator levels, configured radius and system typography.
- Root components: Alert, Avatar, AvatarGroup, Badge, BadgeAnchor, Button,
  ButtonGroup, Card and parts, Chip, CloseButton, EmptyState, Header, Kbd, Link,
  LinkButton, ListGroup, Meter, PressableFeedback, ProgressBar, ProgressCircle,
  ScrollShadow, Separator, Skeleton, SkeletonGroup, Spinner, Surface, Tag,
  TagGroup, ToggleButton, ToggleButtonGroup and Typography.
- Form package: field primitives, Input, TextArea, TextField, InputGroup,
  InputOTP, SearchField, NumberField, Checkbox/Group, Radio/Group, Switch/Group,
  ControlField, Slider, Select, ComboBox and Autocomplete.
- ListBox package with virtualized rows, sections, typeahead and selection.
- Navigation: Accordion, Disclosure/Group, Tabs, Breadcrumbs, Pagination,
  Toolbar and Menu/Dropdown/SubMenu.
- Overlays: Modal, AlertDialog, Drawer, BottomSheet, Dialog, Popover, Tooltip
  and a Toast manager.
- Date: Calendar, RangeCalendar, DateField, DatePicker, DateRangePicker and
  TimeField.
- Color: ColorArea, ColorSlider, ColorField, ColorSwatch, ColorSwatchPicker and
  ColorPicker.
- Data: virtualized Table with sorting, selection, row actions, resizing,
  hiding and an empty state.
- Per-component gallery in the HeroUI docs layout: a 880 DIP main column and
  a right-hand "On this page" rail (hidden below 1000 DIP) whose links scroll
  to Usage, Examples, API Reference and Theme tokens. The page header carries
  the group label, title, description, Copy Markdown and Copy import actions
  with native clipboard feedback. Usage is a compact import bar: with no
  explicit `ComponentPage.Import`, its block is derived from the package
  qualifiers the page's snippet names (module root plus the subpackages and
  `mygoicons/lucide` it uses, de-duplicated and sorted, as one line or a
  parenthesized block), and both copy actions write that block.
- Gallery named examples: `ComponentPage.Examples` holds
  `ComponentExample{Title, Description, Demo}` frames. Each frame renders its
  title and description outside its own clipped variant panel. All 86
  component pages ship 2-5 named examples (260 in total); there is no
  page-level `Demo` fallback.
- Gallery code areas: the page's full Go snippet renders once with line
  numbers, syntax highlighting, selectable lines, horizontal scrolling and a
  tooltip-only copy control; snippets longer than six lines fold under a
  surface-colored fade with Expand code/Collapse code. API Reference tables
  wrap long code pills and stack at narrow widths; theme-token cards show
  light/dark swatches, hex values and role descriptions.
- Gallery overview: a landing hero with feature points and one card per
  component group, with page counts and links. The sidebar groups every
  component under the MyGo UI brand and uses ScrollShadow for overflow;
  light/dark modes, accent and radius customization, and search remain live.
- MyGo UI brand across the README, documentation and gallery shell.

### Changed

- Upgraded MyGo to v0.6.3 and mygoicons to v0.2.1.
- Theme application now installs a HeroUI palette derived from upstream
  variables; `theme.Configure` publishes process-wide accent, radius,
  typography and motion preferences.
- Coordinated collections (accordion sections, disclosures, tabs, menus) take
  typed slices instead of ambient build scopes. No context-keyed globals or
  build values are stored outside the current frame.
- Documentation, examples, notices and validation were rewritten for the
  HeroUI reference and pinned upstream submodules.
- Hover, pressed and focus faces are pinned by headless pointer-move pixel
  tests: Button, ToggleButton, ButtonGroup, Checkbox, Radio, Switch,
  InputGroup, Tabs, Accordion header, ListBox rows, calendar days, DateField
  and the color swatch picker, plus inert/pressed-only behavior for Disclosure
  and ColorArea. A face that stops repainting on hover now fails a test
  instead of being found by eye.

### Fixed

- Toast: the close button now sits inside the card's top-end corner instead
  of painting on the page behind the rounded corner; it still reveals on
  hover or keyboard focus and dismisses the toast.
- Toast: the revealed close button is one centered 12 DIP X on a 20 DIP
  button. It used to paint a pale disc with two `Painter.Line` calls for the
  X, but `Painter.Line` draws only horizontal and vertical strokes, so both
  diagonals collapsed into a dark bar near the disc's top; the disc and its
  border also showed a ring this button does not have. It paints nothing at
  rest and rings only on keyboard focus.
- Toast: entering and exiting cards now fade from `Opacity: 0` as well as
  scaling 0.96 over the same 250 ms, so the toast appears and disappears
  instead of popping in at full opacity. Reduced motion still skips the
  transition.
- Navigation and color indicators: the menu's check mark and submenu
  chevron, the disclosure and accordion chevrons, and the selected color
  swatch tick used `Painter.Line` too, so their diagonals collapsed into
  bars. They now stroke a `ui.Path` with the same sizes, angles and
  right-to-left direction.
- Toast: dismissing a toast no longer panics while its exit transition paints
  the retained copy, which read Elements from the build pass that made it.
- Component state faces (buttons, toggles, tags, links, form controls, date,
  color, navigation and list rows) snapshot hover, pressed and focus state
  during the build instead of reading Elements from a paint callback, so any
  component placed inside an `ElementTransition` exit keeps painting without
  an expired-element panic.

No version has been released from this checkout.
