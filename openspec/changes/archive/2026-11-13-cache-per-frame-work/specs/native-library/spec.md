## ADDED Requirements

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
treatment inside the viewport instead of running its shader.

#### Scenario: A mask fade renders in software

- **WHEN** a viewport with a mask fade has more to scroll below
- **THEN** its bottom edge is shaded and its content is otherwise unchanged
