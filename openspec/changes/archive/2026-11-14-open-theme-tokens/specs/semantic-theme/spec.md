# Spec Delta

## ADDED Requirements

### Requirement: Motion, spacing and effect values are theme tokens

The theme SHALL publish the durations and easings of shared motion, the
spacing unit and the shared effect parameters as tokens beside colour,
radius, type and elevation. Every component that performs a kind of change
governed by a token SHALL read that token, so one override changes them
together. A value peculiar to one component MAY stay with that component.
The built-in theme SHALL reproduce the values recorded in root `DESIGN.md`.

#### Scenario: A host shortens the state duration

- **WHEN** a host publishes a theme whose state duration is shorter and a
  control changes state
- **THEN** the control reaches its new state within the supplied duration

#### Scenario: No theme is supplied

- **WHEN** an application uses the library without publishing a theme
- **THEN** rendering and timing equal the recorded defaults

### Requirement: A theme is one value

A host SHALL be able to supply a complete design language as one typed
value holding both palettes, the radius base, type, spacing, elevation,
motion and effects, and SHALL be able to obtain the built-in theme as such a
value. A host that changes part of the built-in value SHALL inherit the rest.
The theme MUST copy caller-owned storage at publication and when returned.

#### Scenario: A host supplies palette and geometry together

- **WHEN** a host publishes a theme with its own palettes and radius base
- **THEN** components paint the supplied colours and derive radii from the
  supplied base in the same build

#### Scenario: A host overrides one token

- **WHEN** a host takes the built-in theme, changes one motion token and
  publishes it
- **THEN** every other token keeps its built-in value

A theme is application input. Publication MUST refuse, with an error naming
the field and without changing the published theme, a negative duration and
a spinner or sweep period shorter than a millisecond. A zero motion duration
SHALL mean that the change shows at once. A theme's palettes SHALL be used
as written: the hover opacity token governs only hovers that components
derive from a fill as they draw, and every such hover SHALL read it.

#### Scenario: A theme carries a duration components cannot run

- **WHEN** a host publishes a theme whose sweep period is zero or whose state
  duration is negative
- **THEN** publication returns an error naming that field and the previously
  published theme stays in effect

#### Scenario: A theme sets a motion duration to zero

- **WHEN** a host publishes a theme whose state duration is zero and a
  control changes state
- **THEN** the control is at rest on the frame after the change

#### Scenario: A theme changes the hover opacity

- **WHEN** a host publishes a theme with another hover opacity
- **THEN** the solid of every status hovers at that opacity, and the hover
  roles of the theme's palettes keep the colours the theme gave them

#### Scenario: A caller edits a value after publishing it

- **WHEN** a caller mutates the shadow layers of a theme it published, or of
  a theme the library returned
- **THEN** the published theme is unchanged

### Requirement: Preferences resolve over the published theme

Process-wide preferences SHALL keep their meaning over any published theme:
a preference the host sets replaces the theme's value, and an unset
preference inherits it. A palette preference SHALL also move the tokens
derived from the roles it changes. Publishing a theme MUST NOT discard
published preferences, and updating preferences MUST NOT discard the
published theme.

#### Scenario: A palette override meets a custom theme

- **WHEN** a host publishes a theme and overrides one light role through
  preferences
- **THEN** that role uses the override and the other roles use the theme

A font size or radius preference equal to the published theme's value SHALL
mean inherit, so republishing the preferences the library returned pins
neither.

#### Scenario: Preferences are saved and restored

- **WHEN** a host republishes the preferences snapshot the library returned
  and later publishes a theme with another radius base and base text size
- **THEN** radii and text size follow the new theme

#### Scenario: A preference update leaves theme values inherited

- **WHEN** a host changes only the font family through the preference
  update API and later publishes a theme with another radius base
- **THEN** radii follow the new theme

### Requirement: Reduced motion does not depend on token values

The native reduced-motion preference SHALL keep its treatment under any
theme: control fades play over the theme's durations, while positions,
sizes, scales and loops settle at once and overlays appear and disappear at
once.

#### Scenario: A custom duration under reduced motion

- **WHEN** a theme lengthens overlay durations and reduced motion is on
- **THEN** an overlay still appears and disappears at once
