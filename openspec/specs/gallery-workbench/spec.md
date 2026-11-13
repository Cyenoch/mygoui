# Gallery Workbench

## Purpose

Let Go application authors discover native components, exercise examples
and inspect their public API usage within one native component workbench.

## Requirements

### Requirement: Components can be discovered by navigation and search

The gallery SHALL list every component under its group heading without
requiring a group to be expanded, and SHALL provide keyboard search. Selecting
a component SHALL show its associated examples and public API material without
introducing a separate gallery-only component API.

#### Scenario: A user finds a component through search

- **WHEN** the user opens Cmd/Ctrl+K search and selects a component result
- **THEN** the gallery selects the corresponding component page

#### Scenario: A user scans the navigation

- **WHEN** the gallery opens in a wide window
- **THEN** each component's navigation entry is reachable without expanding a group

### Requirement: Preview and reference views share the selected component

A component page SHALL be one scrolling document with, in order, the lead
example, installation, usage, API reference, further examples and theme
tokens. Each example SHALL switch between its preview and its own Go source.
Window-owned example state SHALL survive that switch.

#### Scenario: A user returns to a stateful preview

- **WHEN** the user changes an example's value, views its source, then returns
  to Preview
- **THEN** the example retains its window-owned value

#### Scenario: A user reads an example's source

- **WHEN** the user switches an example to Code
- **THEN** the Go source of that example is shown with a copy control

### Requirement: Window chrome and resizing preserve native operation

The gallery SHALL reserve actual system window-control insets and draggable
title-bar regions on desktop. Compact layouts and native content zoom SHALL
preserve reachable navigation and content while respecting native focus. The
browser build SHALL start with system fonts.

#### Scenario: A compact window changes content scale

- **WHEN** the user changes native content zoom or resizes the window
- **THEN** content uses its scaled layout thresholds
- **AND** the title bar covers the native title band and keeps its controls
  clear of the system window controls

#### Scenario: The gallery runs in a browser

- **WHEN** the gallery starts where bundled fonts cannot be registered
- **THEN** it starts with the system font and offers no font choice

### Requirement: Examples demonstrate actual library behavior

Gallery examples and their Go snippets SHALL consume public components.
The component matrix SHALL distinguish native adaptations and gaps.
Headless specimens SHALL be recorded as headless evidence rather than actual
operating-system or GPU results.

#### Scenario: A gallery specimen is captured without a native window

- **WHEN** a tester renders an example into an image
- **THEN** the resulting evidence is identified as headless
- **AND** it is not reported as proof of system-menu or GPU presentation

### Requirement: Opening gallery controls does not move the page

Appearance controls and compact navigation SHALL open as overlays. The header
SHALL keep one height on every page, and copy feedback SHALL not change an
action's width.

#### Scenario: A user opens appearance controls

- **WHEN** the user opens appearance controls on a component page
- **THEN** the page content keeps its position

### Requirement: The home view introduces and distributes the gallery

The gallery's home view SHALL present the library, its install command, live
components, the component catalog and download links for each published
platform package.

#### Scenario: A visitor wants the desktop gallery

- **WHEN** a visitor chooses a platform's download on the home view
- **THEN** the gallery opens that platform's latest release asset externally

### Requirement: Builds publish the gallery

A release tag SHALL attach macOS, Windows, Linux and web gallery packages
with stable names to its GitHub release. The default branch SHALL publish the
browser gallery to GitHub Pages.

#### Scenario: A version is tagged

- **WHEN** a `v*` tag is pushed
- **THEN** the release contains each platform's gallery package under its stable name

### Requirement: Reference content is derived from source

Example code SHALL be the example's own source. The API reference SHALL list
every exported field of the option types the page's usage snippet names,
with its Go type, derived from the library source, and the build SHALL fail when that derived data is
stale.

#### Scenario: An option field is added to a component

- **WHEN** a field is added to an option type and the reference is regenerated
- **THEN** the component's page lists the field without a hand-written row

#### Scenario: The generated reference is stale

- **WHEN** the library source changes and the reference is not regenerated
- **THEN** the test suite fails

### Requirement: A page index follows the reader

Wide windows SHALL show the page's sections beside the content, mark the
section being read and scroll to a section when it is chosen. Narrow windows
SHALL omit the index without losing content.

#### Scenario: A section is chosen

- **WHEN** the user activates a section in the index
- **THEN** the content scrolls to that section

### Requirement: The sidebar builds the rows near its viewport

The sidebar SHALL build the navigation rows near its scrolled viewport and
keep its scroll range, and SHALL keep building a row that holds the keyboard
focus.

#### Scenario: A focused row scrolls out of view

- **WHEN** a navigation row holds the focus and the sidebar scrolls it away
- **THEN** the row keeps the focus

### Requirement: Compositions are grouped apart from primitives

The gallery SHALL list viewers of structured content under a "Viewers" group
and finished compositions of several controls under a "Patterns" group, after
the primitive groups. A component page's column SHALL be centred in the
content area.

#### Scenario: A user looks for a finished pattern

- **WHEN** the user scans the navigation
- **THEN** the tour, command palette and settings group appear under Patterns

#### Scenario: The window is wider than the page column

- **WHEN** a component page is shown in a wide window
- **THEN** the column has equal space on both sides

### Requirement: The gallery shell is built from the layout family

The gallery SHALL build its title bar, its full-height sidebar and its
compact navigation from the library's public layout family, not from a
private shell. Collapsing the sidebar SHALL move its width as a collapsing
pane does and SHALL settle at once under reduced motion. The compact
navigation SHALL show the pane the wide window docks.

#### Scenario: The sidebar is collapsed in a wide window

- **WHEN** the user activates the navigation toggle in a wide window
- **THEN** the sidebar leaves once its width has settled and the toggle keeps
  the focus

#### Scenario: Navigation is opened in a compact window

- **WHEN** the user activates the navigation toggle in a compact window
- **THEN** search and the page list open in a sheet clear of the window
  controls
