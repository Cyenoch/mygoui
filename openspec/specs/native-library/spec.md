# Native Library

## Purpose

Preserve the caller-owned state, native interaction and import boundaries
that application authors depend on when composing mygoui components.

## Requirements

### Requirement: Persistent state belongs to the caller

Components SHALL keep persistent values and widget state under caller or
window ownership. Build-scoped contexts, elements and custom parts MUST NOT
be retained across builds or used by background workers.

#### Scenario: A control is rebuilt

- **WHEN** a window rebuilds a control with its existing model and widget state
- **THEN** its current value and persistent identity remain available
- **AND** callbacks obtain geometry or services through valid native handles
  and services rather than a saved element

### Requirement: Native interaction semantics are preserved

Styled controls SHALL retain native editing, selection, focus and keyboard
behavior. Disabled actions MUST remain inert; read-only editors SHALL retain
selection. Labels and applicable help/error metadata SHALL remain associated
with their native controls.

#### Scenario: An application disables an action

- **WHEN** pointer or keyboard input reaches a disabled action
- **THEN** the action's callback is not invoked

#### Scenario: An editor is read only

- **WHEN** a user selects text in a read-only editor
- **THEN** selection remains available without changing the bound value

### Requirement: Imports do not take over application startup

Importing library packages MUST NOT initialize an application, register
application languages or fonts, perform I/O, use persistent storage or
start background jobs. Applications SHALL own explicit startup registration.

#### Scenario: A consumer imports the optional fonts package

- **WHEN** the package is imported without calling its registration function
- **THEN** no application font registration occurs

### Requirement: Background work publishes through the window

Slow external work SHALL stay outside the UI thread. Results that change
window-owned state SHALL be published through the native window update
mechanism; invalidation alone MUST NOT serve as state synchronization.

#### Scenario: An asynchronous result changes displayed state

- **WHEN** an application worker completes work for a live window
- **THEN** it publishes the model change through the window update mechanism
- **AND** the worker does not access build-scoped UI values

### Requirement: Native adaptations have explicit boundaries

The component capability documentation SHALL identify supported subsets and
toolkit limits. A catalog mapping MUST NOT imply web parity or successful
GPU, operating-system gesture, IME or screen-reader validation.

#### Scenario: A reference feature is not implemented

- **WHEN** the upstream component exposes a feature unavailable in the native API
- **THEN** the component matrix records the limitation instead of claiming parity
