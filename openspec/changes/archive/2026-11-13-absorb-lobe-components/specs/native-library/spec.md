## ADDED Requirements

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
