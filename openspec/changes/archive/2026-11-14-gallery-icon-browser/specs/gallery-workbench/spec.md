# Spec Delta

## ADDED Requirements

### Requirement: Icons can be found and copied

The gallery SHALL list the icons of the Lucide collection, which the library
draws with, on one page. The reader SHALL be able to narrow the list by Go
name or Iconify ID. Each drawing SHALL appear once: an alias is found
by its name but shown as the drawing it repeats. Activating an icon SHALL
copy the Go expression that returns it and say so. The page SHALL build only
the rows near its viewport.

#### Scenario: An icon is searched by its Iconify ID

- **WHEN** the reader types `arrow-right` into the icon search
- **THEN** the list narrows to the matching icons and reports their count

#### Scenario: An icon is activated

- **WHEN** the reader activates the icon named ArrowRight
- **THEN** the clipboard holds `lucide.ArrowRight()`
- **AND** the page reports the copy

#### Scenario: Nothing matches

- **WHEN** the search matches no icon
- **THEN** the page shows an empty state instead of the list
