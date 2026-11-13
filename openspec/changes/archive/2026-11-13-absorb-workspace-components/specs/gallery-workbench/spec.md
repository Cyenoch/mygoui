## ADDED Requirements

### Requirement: Compositions are grouped apart from primitives

The gallery SHALL list viewers of structured content under a "Viewers" group
and finished compositions of several controls under a "Patterns" group, after
the primitive groups. A component page's column SHALL be centred in the
content area.

#### Scenario: A user looks for a finished pattern

- **WHEN** the user scans the navigation
- **THEN** the tour, command palette and settings group appear under Patterns

#### Scenario: The window is wider than the page column

- **WHEN** a component page is shown in a wide window
- **THEN** the column has equal space on both sides
