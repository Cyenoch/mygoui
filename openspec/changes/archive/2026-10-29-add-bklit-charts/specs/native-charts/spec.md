# Native Charts

## Purpose

Let native Go applications visualize and inspect caller-owned data with chart
families and composable utilities that inherit the existing coss theme.

## ADDED Requirements

### Requirement: Chart catalog has native counterparts

The library SHALL expose native counterparts for area, bar, candlestick,
choropleth, composed, funnel, gauge, heatmap, line, profit/loss line, live line,
pie, radar, radial, ring, scatter, sankey, sparkline and sunburst charts.
Capability documentation MUST distinguish supported native behavior from
upstream web-specific styling and animation options.

#### Scenario: A consumer explores a chart

- **WHEN** a consumer selects any documented chart subject in Gallery
- **THEN** a real public native component and usable Go example are available
- **AND** unsupported upstream details are documented without a parity claim

### Requirement: Data boundaries are explicit

Charts MUST report invalid numeric data, malformed geometry and invalid flow or
hierarchy relationships as contextual errors. Empty data SHALL render an empty
state, not fabricated values. Rendering MUST NOT mutate caller data.

#### Scenario: A flow has a cycle

- **WHEN** a consumer supplies a cyclic Sankey graph
- **THEN** construction reports an error instead of hanging or drawing a false flow

#### Scenario: A series is empty

- **WHEN** a consumer supplies no samples
- **THEN** the chart displays an empty state without inventing measurements

### Requirement: Inspection preserves native input and ownership

Interactive charts SHALL support pointer inspection and keyboard focus with
arrow, Home, End and Escape navigation. Selection and hidden-series state SHALL
be caller-owned. Disabled charts MUST NOT invoke selection actions. Tooltips
and native value metadata SHALL describe the inspected measurement.

#### Scenario: Keyboard inspection

- **WHEN** a focused chart receives an arrow key
- **THEN** the selected measurement and its accessible value update together
- **AND** Escape clears the selection

### Requirement: Charts compose with semantic theme and utilities

Charts SHALL use existing coss chart, surface, typography and state tokens.
Consumers SHALL be able to compose axes, grid, background, legend, tooltip,
brush, markers, reference areas, bands, projection lines, weekend shading,
deadlines, trend lines and custom indicators. Layout SHALL follow available
native width and text scaling without introducing a web runtime.

#### Scenario: A consumer customizes a chart

- **WHEN** a host overrides chart tokens and adds an annotation or custom indicator
- **THEN** the next build uses those tokens and draws the addition in plot coordinates

### Requirement: Streaming and hierarchy actions are caller controlled

Live charts SHALL consume updated caller samples without starting background
jobs. Sunburst charts SHALL expose drill-down and return navigation. Geographic
charts SHALL consume caller-supplied region polygons without network downloads.

#### Scenario: A host publishes live data

- **WHEN** a host updates samples through its window update mechanism
- **THEN** the live chart displays the current visible sample window
- **AND** the library does not create a timer, worker or network request
