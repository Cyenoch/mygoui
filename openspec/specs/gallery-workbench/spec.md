# Gallery Workbench

## Purpose

Let Go application authors discover native components, exercise examples
and inspect their public API usage within one native component workbench.

## Requirements

### Requirement: Components can be discovered by navigation and search

The gallery SHALL provide grouped component navigation and keyboard search.
Selecting a component SHALL show its associated examples and public API
material without introducing a separate gallery-only component API.

#### Scenario: A user finds a component through search

- **WHEN** the user opens Cmd/Ctrl+K search and selects a component result
- **THEN** the gallery selects the corresponding component page

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
title-bar regions. Compact layouts and native content zoom SHALL preserve
reachable navigation and content while respecting native focus.

#### Scenario: A compact window changes content scale

- **WHEN** the user changes native content zoom or resizes the window
- **THEN** content uses its scaled layout thresholds
- **AND** the system-control strip retains its window-space geometry

### Requirement: Examples demonstrate actual library behavior

Gallery examples and their Go snippets SHALL consume public components.
The component matrix SHALL distinguish native adaptations and gaps.
Headless specimens SHALL be recorded as headless evidence rather than actual
operating-system or GPU results.

#### Scenario: A gallery specimen is captured without a native window

- **WHEN** a tester renders an example into an image
- **THEN** the resulting evidence is identified as headless
- **AND** it is not reported as proof of system-menu or GPU presentation
