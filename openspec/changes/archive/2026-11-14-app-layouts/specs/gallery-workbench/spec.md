# Spec Delta

## ADDED Requirements

### Requirement: The gallery shell is built from the layout family

The gallery SHALL build its title bar, its full-height sidebar and its
compact navigation from the library's public layout family, not from a
private shell. Collapsing the sidebar SHALL move its width as a collapsing
pane does and SHALL settle at once under reduced motion. The compact
navigation SHALL show the pane the wide window docks.

#### Scenario: The sidebar is collapsed in a wide window

- **WHEN** the user activates the navigation toggle in a wide window
- **THEN** the sidebar leaves once its width has settled and the toggle keeps
  the focus

#### Scenario: Navigation is opened in a compact window

- **WHEN** the user activates the navigation toggle in a compact window
- **THEN** search and the page list open in a sheet clear of the window
  controls
