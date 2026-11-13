# Spec Delta

## Purpose

Give MyGo applications the window-level layout every native desktop
application needs: a custom title bar, side panes beside the content and the
shell that composes them, with caller-owned state and native focus.

## ADDED Requirements

### Requirement: The title bar keeps clear of the window controls

`layout.TitleBar` SHALL place its leading, centre and trailing content
outside the window-control insets the window reports on the edges the bar
reaches, and SHALL release an inset when the window stops reporting it. The
library SHALL NOT assume a platform's inset: only the reported values apply.

#### Scenario: A window reports leading and trailing controls

- **WHEN** a window reports control insets on its left and right edges
- **THEN** every title-bar control lies between the two insets

#### Scenario: The window enters full screen

- **WHEN** the window stops reporting an inset
- **THEN** the title-bar content moves into the room the inset held

#### Scenario: A bar reaches one window edge

- **WHEN** a bar is declared to reach only the end edge
- **THEN** it reserves the end inset and no start inset

### Requirement: The title bar covers the native band and grows with scale

`layout.TitleBar` SHALL be at least as tall as the native title band and the
caller's minimum, and tall enough for its controls at the current text scale
and content zoom. When its leading and trailing content no longer fit
between the reserved insets, the content SHALL move below the native band,
where it has the bar's full width.

#### Scenario: Text is enlarged in a narrow window

- **WHEN** the leading and trailing controls are wider than the room between
  the window controls
- **THEN** the bar grows and its controls sit below the native band, inside
  the bar

#### Scenario: A window has no window controls

- **WHEN** the window reports no insets
- **THEN** the controls stay in the band, since moving them gains no room

### Requirement: The title bar moves the window without taking its controls

The title bar's own surface SHALL be a window drag handle. Controls inside
it SHALL receive their presses and keyboard focus as they do elsewhere.
Whether the operating system then moves or maximizes the window is the
toolkit's behavior and is not asserted by headless tests.

#### Scenario: A control in the bar is pressed

- **WHEN** the user presses a button inside the title bar
- **THEN** the button is activated

#### Scenario: The empty part of the bar is pressed

- **WHEN** the user presses the bar where it has no control
- **THEN** the bar itself is the input target

### Requirement: A split view sizes its panes within caller bounds

`layout.SplitView` SHALL show a start pane, the content and an end pane in
layout direction. A resizable pane's divider SHALL move by pointer and by
keyboard within the pane's minimum and maximum, SHALL leave the content its
minimum width while the view has the room, and SHALL be a splitter named
after the pane. A fixed pane SHALL keep its width and offer no divider. Pane
size is caller-owned when the caller passes it.

#### Scenario: A divider is dragged past the maximum

- **WHEN** the user drags a pane's divider beyond its maximum
- **THEN** the caller's size stops at the maximum

#### Scenario: A divider is moved from the keyboard

- **WHEN** the divider has focus and the user presses the arrow that grows
  the pane
- **THEN** the caller's size grows by one step

#### Scenario: The layout direction is right to left

- **WHEN** the layout direction is right to left
- **THEN** the start pane is on the right and its divider grows it leftwards

### Requirement: A pane collapses and returns under caller-owned state

A pane with a caller-owned collapsed flag SHALL hide while the flag is set
and SHALL keep its size for its return. A resizable pane's divider SHALL
stay reachable by Tab while the pane is collapsed and SHALL expand it from
the keyboard. Collapsing SHALL move the width over the shared reveal
duration and SHALL settle at once under reduced motion.

#### Scenario: A collapsed pane is reached from the keyboard

- **WHEN** a resizable pane is collapsed and the user tabs to its divider and
  presses Enter
- **THEN** the pane expands at the size it had

#### Scenario: Motion is reduced

- **WHEN** reduced motion is on and a pane collapses
- **THEN** the pane is gone in the next frame

### Requirement: A compact view shows side panes over the content

In compact mode a split view SHALL show the content alone and SHALL present
a side pane as a modal sheet from the pane's edge while the caller's open
flag is set. The sheet SHALL be dismissed by Escape or a press outside it,
clearing the flag and returning focus. Leaving compact mode SHALL clear the
flag.

#### Scenario: A pane is opened in a compact window

- **WHEN** the caller sets a pane's open flag in compact mode
- **THEN** the pane's body is shown in a sheet at the pane's edge

#### Scenario: Escape is pressed in the sheet

- **WHEN** the user presses Escape while the pane's sheet is open
- **THEN** the open flag is cleared

#### Scenario: The window becomes wide again

- **WHEN** the window widens past the threshold while a pane's sheet is open
- **THEN** the open flag is cleared and the pane is docked as its collapsed
  flag says

### Requirement: The compact mode follows the window or the caller

By default a view SHALL be compact while the window is narrower than 760 DIP
scaled with the text size; the caller SHALL be able to force either mode.
`layout.Compact` SHALL report the same decision. Controls in the content
SHALL keep their focus and state across the change of mode.

#### Scenario: The window becomes narrow

- **WHEN** the window narrows below the compact threshold
- **THEN** the side panes leave the row and the content keeps its focus

#### Scenario: The caller forces compact mode

- **WHEN** the caller forces compact mode in a wide window
- **THEN** the side panes are shown only as sheets

### Requirement: A pane toggle reaches its pane in either mode

`layout.PaneToggle` SHALL be a named toggle button that flips the pane's
collapsed flag while the pane is docked and its open flag while the view is
compact, and SHALL report as checked while the pane shows.

#### Scenario: The toggle is used in a wide window

- **WHEN** the user activates the toggle while the pane is docked
- **THEN** the collapsed flag flips and the open flag is unchanged

#### Scenario: The toggle is used in a compact window

- **WHEN** the user activates the toggle in compact mode
- **THEN** the open flag is set and the collapsed flag is unchanged

### Requirement: The scaffold composes the standard window shell

`layout.Scaffold` SHALL compose a full-height start pane whose top band
shares the title band, the title bar over the content, an end pane beside
the content and a strip under everything. The start pane's band and the
title bar SHALL have one height. The start pane's band SHALL reserve the
start window-control inset; the title bar SHALL reserve the end inset, and
the start inset as well whenever the start pane is not docked beside it.

#### Scenario: The start pane is docked

- **WHEN** the start pane is shown beside the content
- **THEN** its band and the title bar have the same height, and the title
  bar's controls begin after the pane

#### Scenario: The start pane is collapsed

- **WHEN** the start pane is collapsed
- **THEN** the title bar's leading controls clear the start window controls

### Requirement: The scaffold keeps its parts clear while panes move

While the start pane collapses or expands, the title bar of a
`layout.Scaffold` SHALL reserve whatever of the start window-control inset
the pane does not cover, in every frame. The band of a start pane shown as a
sheet SHALL also clear the end window controls the sheet reaches. The start
pane's divider SHALL leave the content its minimum beside a docked end pane.

#### Scenario: The start pane is folding

- **WHEN** the start pane is part way through collapsing or expanding
- **THEN** the title bar's leading controls are clear of the start window
  controls and move with the pane's edge, without a jump

#### Scenario: The start divider is dragged towards a docked end pane

- **WHEN** the start pane's divider is dragged as far as it goes while an
  end pane is docked
- **THEN** the content between the panes keeps its minimum width and the end
  pane keeps its size

### Requirement: A side sheet can hand its panel to the caller

`overlay.Sheet` with `Flush` SHALL build the caller's body alone, edge to
edge and at the panel's full height, without the sheet's padding, header,
footer, close button or own scrolling. It SHALL keep the sheet's modal
focus, dismissal, motion and accessible name.

#### Scenario: A flush sheet shows a full-height body

- **WHEN** a flush side sheet opens with a body that fills its parent
- **THEN** the body spans the panel from edge to edge and top to bottom

### Requirement: Layout state stays with the caller

The layout family SHALL keep no state outside the window's widget state and
the pointers the caller passes, SHALL do no work on import, and SHALL NOT
store pane sizes or flags. A pane without a caller-owned size SHALL keep its
size for the life of the window only.

#### Scenario: An application restores a saved layout

- **WHEN** the application sets the size and flags it saved before building
- **THEN** the view shows them without consulting any library state

### Requirement: Edge surfaces enter from the edge they rest on

A side `overlay.Sheet` or `overlay.Drawer` rests on the edge its side names
in the layout direction. Its entry and exit motion SHALL cross that same
edge in both layout directions.

#### Scenario: A start sheet opens right to left

- **WHEN** a sheet on the start side opens in a right-to-left window
- **THEN** it rests on the right edge and arrives from the right

### Requirement: A resizable panel's handle keeps its screen order

`mygoui.ResizablePanel` SHALL place its handle before the panel's content in
the Tab order when the handle is on the start or top edge, and after it
otherwise, and the handle SHALL take the pointer over the content on its
own strip. With `Shown` set, the panel SHALL report the size it has in each
build, which trails the caller's size while it collapses or expands.

#### Scenario: Tab enters a panel resized from its start edge

- **WHEN** focus moves by Tab into a panel whose handle is on its start edge
- **THEN** the handle is focused before any control in the panel
