# Native Library

## ADDED Requirements

### Requirement: Anchored popups share one surface and coss motion

Anchored popups SHALL paint the shared popover surface. Popovers and pickers
SHALL fade in from 98% around the edge facing their anchor and fade out;
tooltips SHALL also shrink on exit. Menus, selects and combo boxes SHALL open
and close without motion. Reduced motion SHALL remove every popup motion.

#### Scenario: A popover closes

- **WHEN** an open popover is dismissed with motion enabled
- **THEN** its panel fades out from its full size

#### Scenario: A select opens

- **WHEN** a select popup opens
- **THEN** its list is shown at once on the shared surface

### Requirement: Error messages imply the invalid state

A field that shows an error message SHALL present its invalid border and
semantics even when the caller does not set the separate invalid flag.

#### Scenario: Only an error is supplied

- **WHEN** a form field is built with an error message and no invalid flag
- **THEN** it shows the error text and the destructive border
