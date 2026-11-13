# Native Library

## Purpose

Preserve the caller-owned state, native interaction and import boundaries
that application authors depend on when composing mygoui components.

## Requirements

### Requirement: Persistent state belongs to the caller

Components SHALL keep persistent values and widget state under caller or
window ownership. Build-scoped contexts, elements and custom parts MUST NOT
be retained across builds or used by background workers.

#### Scenario: A control is rebuilt

- **WHEN** a window rebuilds a control with its existing model and widget state
- **THEN** its current value and persistent identity remain available
- **AND** callbacks obtain geometry or services through valid native handles
  and services rather than a saved element

### Requirement: Native interaction semantics are preserved

Styled controls SHALL retain native editing, selection, focus and keyboard
behavior. Disabled actions MUST remain inert; read-only editors SHALL retain
selection. Labels and applicable help/error metadata SHALL remain associated
with their native controls.

#### Scenario: An application disables an action

- **WHEN** pointer or keyboard input reaches a disabled action
- **THEN** the action's callback is not invoked

#### Scenario: An editor is read only

- **WHEN** a user selects text in a read-only editor
- **THEN** selection remains available without changing the bound value

### Requirement: Imports do not take over application startup

Importing library packages MUST NOT initialize an application, register
application languages or fonts, perform I/O, use persistent storage or
start background jobs. Applications SHALL own explicit startup registration.

#### Scenario: A consumer imports the optional fonts package

- **WHEN** the package is imported without calling its registration function
- **THEN** no application font registration occurs

### Requirement: Background work publishes through the window

Slow external work SHALL stay outside the UI thread. Results that change
window-owned state SHALL be published through the native window update
mechanism; invalidation alone MUST NOT serve as state synchronization.

#### Scenario: An asynchronous result changes displayed state

- **WHEN** an application worker completes work for a live window
- **THEN** it publishes the model change through the window update mechanism
- **AND** the worker does not access build-scoped UI values

### Requirement: Native adaptations have explicit boundaries

The component capability documentation SHALL identify supported subsets and
toolkit limits. A catalog mapping MUST NOT imply web parity or successful
GPU, operating-system gesture, IME or screen-reader validation.

#### Scenario: A reference feature is not implemented

- **WHEN** the upstream component exposes a feature unavailable in the native API
- **THEN** the component matrix records the limitation instead of claiming parity

### Requirement: Anchored popups share one surface and popup motion

Anchored popups SHALL paint the shared popover surface. Popovers and pickers
SHALL fade in from 98% around the edge facing their anchor and fade out;
tooltips and preview cards SHALL also shrink on exit. Menus, selects and
combo boxes SHALL open and close without motion. Reduced motion SHALL remove
every popup motion.

#### Scenario: A popover closes

- **WHEN** an open popover is dismissed with motion enabled
- **THEN** its panel fades out from its full size

#### Scenario: A select opens

- **WHEN** a select popup opens
- **THEN** its list is shown at once on the shared surface

#### Scenario: A preview card closes

- **WHEN** an open preview card closes with motion enabled
- **THEN** its panel fades out while shrinking toward its anchor

### Requirement: Error messages imply the invalid state

A field that shows an error message SHALL present its invalid border and
semantics even when the caller does not set the separate invalid flag.

#### Scenario: Only an error is supplied

- **WHEN** a form field is built with an error message and no invalid flag
- **THEN** it shows the error text and the destructive border

### Requirement: Control state motion follows the recorded design language

A pressed switch SHALL stretch its thumb toward the track's centre. Select and
combo box triggers SHALL show a static up-down indicator. A sidebar SHALL move
between its widths over 200 ms. The tabs indicator SHALL use ease-in-out and
the bottom sheet the drawer curve. An indeterminate progress sweep SHALL
redraw at no more than 60 frames per second. Reduced motion SHALL settle each
of these at once.

#### Scenario: A switch is pressed

- **WHEN** the pointer holds an enabled switch with motion enabled
- **THEN** its thumb becomes wider than it is at rest

#### Scenario: A sidebar collapses

- **WHEN** a sidebar's collapsed state changes with motion enabled
- **THEN** its width passes through intermediate values before settling

#### Scenario: Motion is reduced

- **WHEN** a sidebar collapses with reduced motion
- **THEN** it reaches its collapsed width in the next frame

### Requirement: Content components copy, colour and render text

The `content` family SHALL provide a copy button that writes caller text to
the clipboard and shows copied feedback, a one-line copyable snippet, a
read-only highlighted code block with a copy action, a themed Markdown
renderer and a QR code. They SHALL use theme tokens and work in both
appearances.

#### Scenario: Text is copied

- **WHEN** the user activates a copy button
- **THEN** the clipboard holds the caller's text and the button shows a copied state that reverts

#### Scenario: An unknown language is highlighted

- **WHEN** a highlighter receives a language it does not know
- **THEN** it shows the code unstyled in the monospace face

#### Scenario: QR content cannot be encoded

- **WHEN** QR content exceeds what a QR code can hold
- **THEN** the component reports the error to the caller instead of drawing a partial code

### Requirement: Form inputs cover inline editing, shortcuts, files, rating and image choice

`form` SHALL provide text that edits in place, a shortcut recorder, a file
drop zone with picker, a star rating and an image tile selection. Their
values SHALL be caller-owned, and each SHALL be operable from the keyboard
and expose its role and label.

#### Scenario: An inline edit is committed

- **WHEN** the user edits the text and presses Enter or moves focus away
- **THEN** the caller's value is updated and the text is shown again

#### Scenario: An inline edit is cancelled

- **WHEN** the user presses Escape while editing
- **THEN** the caller's value is unchanged

#### Scenario: A shortcut is recorded

- **WHEN** the recorder has focus and the user presses a key with modifiers
- **THEN** the caller's shortcut holds that combination and it is shown as keys

#### Scenario: Files are dropped

- **WHEN** files are dropped on the upload zone or chosen in the picker
- **THEN** the caller receives their paths and the library reads none of them

#### Scenario: A rating changes from the keyboard

- **WHEN** a focused rating receives an arrow key
- **THEN** the caller's value moves by one step within its bounds

### Requirement: Steps and sortable lists show and change order

`navigation.Steps` SHALL show numbered steps with completed, current and
upcoming states, horizontally or vertically. `list.SortableList` SHALL let
the user reorder caller-owned items by dragging a handle or from the
keyboard and report the new order to the caller.

#### Scenario: A step is current

- **WHEN** steps render with a current index
- **THEN** earlier steps show as completed and later steps as upcoming

#### Scenario: An item is moved from the keyboard

- **WHEN** a focused sortable item's handle receives the move-down key
- **THEN** the caller is told to move that item one place later

### Requirement: Floating and resizable surfaces are caller-sized

`ResizablePanel` SHALL let the user resize a panel from its edge within
caller bounds, by pointer and keyboard, with caller-owned size.
`overlay.FloatingPanel` SHALL float a non-modal, draggable panel over
content, and `overlay.FloatingSheet` SHALL present a non-modal sheet from the
bottom edge. Neither floating surface SHALL trap focus or block the content
beneath.

#### Scenario: A panel is resized

- **WHEN** the user drags the panel's handle past its maximum
- **THEN** the caller's size stops at the maximum

#### Scenario: Content under a floating panel is used

- **WHEN** a floating panel is open and the user presses a control outside it
- **THEN** the control is activated and the panel stays open

### Requirement: Image and result surfaces show caller-provided content

`Image` SHALL show a caller-provided image with loading and failed states.
`overlay.ImageViewer` SHALL present an image in a dismissible preview with
zoom. `Result` SHALL show a status outcome with title, description and
actions, and `GuideCard` SHALL show a dismissible tip. The library SHALL
start no image loads itself.

#### Scenario: An image failed to load

- **WHEN** the caller marks an image as failed
- **THEN** the fallback is shown in the image's box with its accessible label

#### Scenario: A guide card is dismissed

- **WHEN** the user activates the card's close control
- **THEN** the caller's dismiss callback runs

### Requirement: Filtered option lists are kept between frames

ComboBox, Autocomplete and Command SHALL keep their filtered rows until the
query, the option or item slices, or the state's caller-owned `Revision`
change. A caller that edits options in place SHALL bump `Revision` for the
edit to show.

#### Scenario: The query changes

- **WHEN** the user types a query
- **THEN** the list shows the options matching it

#### Scenario: Options are edited in place

- **WHEN** the caller changes an option's label in the same slice and bumps
  `Revision`
- **THEN** the list shows the new label

#### Scenario: A long command list grows

- **WHEN** a command list with shortcuts gains rows during layout
- **THEN** those rows draw their shortcuts

### Requirement: Long menus build the rows in view

A menu panel with 32 or more rows SHALL build only the rows near its
viewport at a fixed width, and SHALL keep keyboard movement, wrapping,
disabled and separator skipping, and each item's position for assistive
technology. The panel SHALL keep one presentation for as long as it is open.

#### Scenario: The keyboard reaches a far item

- **WHEN** the user presses End in a menu of 1,000 items
- **THEN** the last item is focused and shown

#### Scenario: An open menu grows

- **WHEN** an open menu's items grow past 32 rows
- **THEN** the focused item keeps the focus

### Requirement: Effects degrade without a GPU

Where the frame does not draw on a GPU, SkeletonGroup SHALL rest static
without requesting frames, and a mask ScrollFade SHALL show the shadow
treatment inside the viewport instead of running its shader. A mask
ScrollFade whose viewport exceeds the toolkit's effect-layer limit SHALL
show the same shadow treatment on a GPU frame.

#### Scenario: A mask fade renders in software

- **WHEN** a viewport with a mask fade has more to scroll below
- **THEN** its bottom edge is shaded and its content is otherwise unchanged

#### Scenario: A mask fade exceeds the layer limit

- **WHEN** a viewport with a mask fade is wider than the effect-layer limit
  and has more to scroll below
- **THEN** its bottom edge is shaded and its content is otherwise unchanged

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
