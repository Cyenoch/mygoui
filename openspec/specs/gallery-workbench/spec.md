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

The gallery SHALL expose Preview, Go source, API and Tokens views for the
selected component. Window-owned example state SHALL survive switching
between reference and preview views.

#### Scenario: A user returns to a stateful preview

- **WHEN** the user changes an example's value, views its source, then returns
  to Preview
- **THEN** the example retains its window-owned value

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
