## MODIFIED Requirements

### Requirement: Anchored popups share one surface and coss motion

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

## ADDED Requirements

### Requirement: Control state motion follows coss

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
