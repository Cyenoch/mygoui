## MODIFIED Requirements

### Requirement: Preview and reference views share the selected component

A component page SHALL be one scrolling document with, in order, the lead
example, installation, usage, API reference, further examples and theme
tokens. Each example SHALL switch between its preview and its own Go source.
Window-owned example state SHALL survive that switch.

#### Scenario: A user returns to a stateful preview

- **WHEN** the user changes an example's value, views its source, then returns
  to Preview
- **THEN** the example retains its window-owned value

#### Scenario: A user reads an example's source

- **WHEN** the user switches an example to Code
- **THEN** the Go source of that example is shown with a copy control

## ADDED Requirements

### Requirement: Reference content is derived from source

Example code SHALL be the example's own source. The API reference SHALL list
every exported field of the option types the page's usage snippet names,
with its Go type, derived from the library source, and the build SHALL fail when that derived data is
stale.

#### Scenario: An option field is added to a component

- **WHEN** a field is added to an option type and the reference is regenerated
- **THEN** the component's page lists the field without a hand-written row

#### Scenario: The generated reference is stale

- **WHEN** the library source changes and the reference is not regenerated
- **THEN** the test suite fails

### Requirement: A page index follows the reader

Wide windows SHALL show the page's sections beside the content, mark the
section being read and scroll to a section when it is chosen. Narrow windows
SHALL omit the index without losing content.

#### Scenario: A section is chosen

- **WHEN** the user activates a section in the index
- **THEN** the content scrolls to that section
