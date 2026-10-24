# Changelog

## Unreleased

Complete native rewrite from the shadcn-inspired API to a HeroUI v3 component
library. The previous component set is not part of this API and no compatibility
shims are provided.

### Added

- HeroUI v3 semantic theme: light/dark palettes, accent/default/status roles,
  derived hover/soft/separator levels, configured radius and system typography.
- Host palette injection: `theme.Preferences.Palette` overlays semantic roles
  per appearance through `theme.PalettePatch` — canvas, surface, ink, muted
  ink, border (transparent is valid), accent and status pairs, hover and
  selection fills, overlay and backdrop. Unset roles keep HeroUI's defaults;
  derived hover, soft and foreground levels rebuild from the injected base.
  `theme.Apply` composes its native theme from the resolved palette, and
  `theme.UpdatePreferences` edits typography without dropping it. `ColorInfo`
  joins the status colors, defaulting to the accent family.
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

- Toast: restored spring entry/exit, stacked cards, hover/focus expansion,
  continuous insertion/removal and horizontal swipe dismissal. Animated
  declared geometry survives MyGo's absolute-overlay relayout. Reduced motion
  settles immediately; native timers, focus pause and inert exit snapshots
  remain authoritative.
- Toast: the close button floats across the upper-right corner without
  reserving message or action width. Undo keeps the normal right padding;
  the centered lucide X reveals on hover or focus and remains keyboard-accessible.
- Gallery Toast triggers create a fresh notification on every press; only
  the explicit Update example reuses its ID.
- Modal, AlertDialog, Dialog, Drawer, BottomSheet, Popover and Tooltip now
  animate entry and exit with native inert snapshots. Their constructors
  keep the overlay layer alive while closed, so animation works without a
  Toast viewport. Reduced motion skips the transitions.
- Modal identity is assigned before native input is read, preserving backdrop
  dismissal on immediate reopen. AlertDialog focuses Cancel only after the
  native base records its opener, so closing restores focus to the trigger.
- Drawer panels keep 16 DIP screen gutters on every side and rounded corners
  instead of touching the window edges.
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
