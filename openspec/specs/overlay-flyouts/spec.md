# overlay-flyouts Specification

## Purpose
Let developers opt anchored popups into native flyout windows that can extend
beyond their parent window, while popups without a host stay in-window.

## Requirements

### Requirement: Anchored popups can open as native flyouts

Popover, Dropdown, Select, DatePicker, DateRangePicker and ColorPicker SHALL
accept an optional caller-owned flyout host. With a host whose window is
available, opening the popup SHALL present its content in a native flyout
window anchored to the trigger; without one, the popup SHALL open in the
window as before. The caller-owned open, value and selection state SHALL be
the same in both presentations.

#### Scenario: A select opens near the window edge

- **WHEN** a Select with a flyout host opens on a desktop platform
- **THEN** its list is presented in a flyout window anchored to the trigger
- **AND** choosing an option updates the caller's selection and closes the flyout

#### Scenario: No host is given

- **WHEN** the same Select has no flyout host
- **THEN** its list opens in the window as before

### Requirement: Flyouts fall back to the window where unavailable

Where a flyout window cannot be created, such as WebAssembly, a popup with a
host SHALL open in the window instead of failing.

#### Scenario: The browser build opens a popover

- **WHEN** a Popover with a flyout host opens in the browser build
- **THEN** it opens in the page as an in-window popover

### Requirement: Flyout state stays consistent with the caller

Closing a flyout by user action (outside press, Escape, choosing an option,
another window taking focus) SHALL set the caller's open state to false and
refresh the parent window. Setting the caller's open state to false,
removing the host, or removing the anchor from the tree SHALL close the
flyout; removing the anchor leaves the caller's open state unchanged.

#### Scenario: The user presses Escape in a flyout

- **WHEN** a popup presented as a flyout is closed with Escape
- **THEN** the caller's open state becomes false
