## ADDED Requirements

### Requirement: The sidebar builds the rows near its viewport

The sidebar SHALL build the navigation rows near its scrolled viewport and
keep its scroll range, and SHALL keep building a row that holds the keyboard
focus.

#### Scenario: A focused row scrolls out of view

- **WHEN** a navigation row holds the focus and the sidebar scrolls it away
- **THEN** the row keeps the focus
