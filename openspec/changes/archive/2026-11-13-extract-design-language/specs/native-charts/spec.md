## MODIFIED Requirements

### Requirement: Charts compose with semantic theme and utilities

Charts SHALL use the semantic theme's chart, surface, typography and state
tokens.
Consumers SHALL be able to compose axes, grid, background, legend, tooltip,
brush, markers, reference areas, bands, projection lines, weekend shading,
deadlines, trend lines and custom indicators. Layout SHALL follow available
native width and text scaling without introducing a web runtime.

#### Scenario: A consumer customizes a chart

- **WHEN** a host overrides chart tokens and adds an annotation or custom indicator
- **THEN** the next build uses those tokens and draws the addition in plot coordinates
