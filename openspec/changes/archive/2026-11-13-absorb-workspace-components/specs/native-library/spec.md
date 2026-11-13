## ADDED Requirements

### Requirement: Loading has one presentation

A component with a `Loading` option SHALL keep its layout while loading. A
region that already shows content SHALL dim it and show a progress line; a
region without content SHALL show a placeholder. An indicator SHALL wait
150 ms before it appears, and a loading control SHALL not act.

#### Scenario: Short work finishes

- **WHEN** loading ends within 150 ms
- **THEN** no indicator was shown

#### Scenario: A loading command list receives Enter

- **WHEN** the user presses Enter while a command list is loading
- **THEN** no item is activated

### Requirement: State changes move without shifting layout

A change of a component's state SHALL be shown by motion of at most 300 ms
that settles at once under reduced motion, SHALL not animate on the first
build and SHALL request no frame at rest. A state SHALL not change a text's
font weight.

#### Scenario: A value changes

- **WHEN** a number flip's value changes with motion enabled
- **THEN** its digits pass through intermediate positions before settling

#### Scenario: A row is selected

- **WHEN** a tree row becomes selected
- **THEN** its label keeps its width

### Requirement: Workspace navigation components keep focus with the user

`navigation` SHALL provide document tabs, a status bar, an overflowing action
row, a menu bar, a navigation stack, a split button and a carousel, and `list`
a tree. Their selection, open and expanded state SHALL be caller-owned. A
window-level shortcut SHALL not take focus from the element that has it.

#### Scenario: Tabs are switched from the document

- **WHEN** focus is in the document and the user presses Ctrl+Tab
- **THEN** the next tab becomes active and focus stays in the document

#### Scenario: The close shortcut meets a pinned tab

- **WHEN** the active tab is pinned and the user presses Cmd/Ctrl+W
- **THEN** no tab is closed

#### Scenario: Escape is pressed in an editor on a pushed page

- **WHEN** the user presses Escape in a text field on a pushed page
- **THEN** the page stays

#### Scenario: Space is pressed on a tree row

- **WHEN** the user presses Space on a row of a selectable tree
- **THEN** the row's selection changes and its action does not run

### Requirement: Viewers search, filter and copy what they show

`content` SHALL provide JSON, log and diff viewers with an optional toolbar.
Search and filter state SHALL be caller-owned. New data SHALL not move the
reader's position unless the query changed.

#### Scenario: A log is filtered by level

- **WHEN** the user presses the Error level
- **THEN** only error lines are listed, and releasing it lists every line again

#### Scenario: A log grows under a filter

- **WHEN** lines are appended while a filter is on
- **THEN** the scroll position is kept

#### Scenario: A JSON container is double-clicked

- **WHEN** the user double-clicks a container row
- **THEN** it toggles between open and closed

### Requirement: Structured inputs keep refused and secret text safe

`form` SHALL provide tag, key-value, masked, path and password inputs with
caller-owned values. A refused tag SHALL stay in the editor with feedback. A
secret value SHALL be shown only after the user asks for that row.

#### Scenario: A duplicate tag is pasted

- **WHEN** pasted text holds a tag the input refuses
- **THEN** the accepted tags are added and the refused one stays in the editor

#### Scenario: A row above a secret is removed

- **WHEN** a revealed row is removed and a secret row takes its place
- **THEN** the secret row is concealed

#### Scenario: A complete value is pasted into a partly filled mask

- **WHEN** the pasted text holds enough characters to fill the mask
- **THEN** it replaces the whole value

### Requirement: A tour shows that it is running

`overlay.Tour` SHALL dim the window around a spotlight on the current step's
target and show one panel with the step's position. Between steps the
spotlight and panel SHALL travel continuously. Its open state and step SHALL
be caller-owned, and closing SHALL return focus.

#### Scenario: The next step is requested

- **WHEN** the user activates Next with motion enabled
- **THEN** the panel passes through positions between the two targets

#### Scenario: A step has no laid-out target

- **WHEN** a step's anchor is not in the window
- **THEN** the panel is centred without a spotlight
