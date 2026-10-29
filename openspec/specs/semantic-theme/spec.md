# Semantic Theme

## Purpose

Keep library and native widget appearance consistent while allowing hosts
to customize colors, type and geometry without replacing component behavior.

## Requirements

### Requirement: Appearance uses live semantic roles

Components SHALL resolve current appearance through the shared semantic
theme. Primary and Accent, and muted background and foreground, SHALL remain
distinct roles. Translucent tokens SHALL composite on their actual surface.

#### Scenario: A host supplies a translucent border

- **WHEN** a component paints that border over its current surface
- **THEN** the supplied transparency is preserved until composition

### Requirement: Hosts can override each appearance independently

Hosts SHALL be able to override light and dark semantic tokens independently.
An omitted override SHALL inherit the default; explicit transparency SHALL
remain an override. Applying the native theme SHALL preserve resolved host
customization for native widgets and library paint.

#### Scenario: Only one appearance is customized

- **WHEN** a host overrides a light token without a dark override
- **THEN** light components use the supplied value and dark components inherit
  their defaults

#### Scenario: A partial preference update preserves customization

- **WHEN** a host changes font size through the preference update API
- **THEN** its existing palette overrides remain present

### Requirement: Font preferences do not load assets

Body, heading and mono family preferences SHALL be independently configurable
with native fallbacks. Setting family names MUST NOT download or register
font data; host startup SHALL own font registration.

#### Scenario: A heading family is omitted

- **WHEN** a host configures a body family without a heading family
- **THEN** headings use the configured body family with native fallbacks

### Requirement: Text scaling and reduced motion remain native preferences

Text-relative component metrics SHALL account for the configured base font
size and OS text scaling. Components with motion SHALL honor the native
reduced-motion preference.

#### Scenario: A user increases OS text scale

- **WHEN** the native preference changes and the view rebuilds
- **THEN** text-relative metrics use the updated scale rather than fixed
  physical-pixel sizes

#### Scenario: Reduced motion is requested

- **WHEN** a motion-enabled control renders with reduced motion enabled
- **THEN** its implementation uses its reduced-motion treatment
