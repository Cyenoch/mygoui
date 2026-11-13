# behavior-reference Specification

## Purpose
Give application authors and maintainers one checkable record of how each
public mygoui component behaves, with the evidence behind every statement.

## Requirements

### Requirement: Every public component has a behavior entry

The repository SHALL keep a behavior reference under `docs/behavior/`, one
document per package family, whose README defines a single entry format.
Each exported component builder SHALL have one entry in its family's
document, using the README's headings in the README's order, or be named by
the entry of the component it serves.

#### Scenario: An author looks up a component

- **WHEN** an author opens the family document of an exported component
- **THEN** it contains an entry for that component with the README's headings
  in the README's order
- **AND** rows that do not apply are omitted rather than filled

#### Scenario: A family document is complete

- **WHEN** a family's task in the change is marked done
- **THEN** every exported builder of that family's packages has an entry or is
  named in another entry
- **AND** the document ends with a gap list measured against platform
  conventions and the WAI-ARIA Authoring Practices

### Requirement: Behavior statements carry their evidence

Every statement in an entry SHALL carry one evidence mark defined in the
README: protected by a named existing test, protected by a test added for the
reference, or observed in the code only. A test mark SHALL name a test that
exists in the repository and exercises the statement.

#### Scenario: A statement names its test

- **WHEN** a statement is marked as protected by a test
- **THEN** a test with that name exists and asserts the stated behavior

#### Scenario: A statement has no test

- **WHEN** behavior was read from the source and no test protects it
- **THEN** the statement is marked as observed in the code only

### Requirement: The reference attributes evidence to its source

Behavior the MyGo toolkit owns, such as the caret, text selection and input
methods, SHALL be attributed to the toolkit instead of enumerated. Headless
test results SHALL NOT be presented as screen-reader, operating-system or GPU
evidence.

#### Scenario: An entry covers a text editor

- **WHEN** an entry describes a component built on a native editor
- **THEN** caret, selection, clipboard and input-method behavior are named as
  the toolkit's and not listed key by key

#### Scenario: A reader weighs a test mark

- **WHEN** a reader checks what a test mark proves
- **THEN** the README states that headless tests prove the accessibility
  tree, focus and key handling only

### Requirement: Behavior changes update the reference

A change that alters a component's observable keys, focus, dismissal, state
handling, motion or assistive-technology output SHALL update that component's
entry in the same change. A defect found while writing the reference SHALL be
recorded in the family's gap list rather than described as intended behavior.

#### Scenario: A key binding changes

- **WHEN** a change adds, removes or alters a key a component answers
- **THEN** the same change updates the component's entry and its evidence mark

#### Scenario: Writing an entry reveals a defect

- **WHEN** the code contradicts a platform convention or a documented promise
- **THEN** the entry states what the code does
- **AND** the gap list records the difference
