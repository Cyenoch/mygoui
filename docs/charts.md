# Native charts

`github.com/Cyenoch/mygoui/chart` covers the 19 Bklit chart subjects and utility
subjects in native MyGo, with coss appearance. Series use Chart1 through Chart5;
axes use Border/MutedForeground; surfaces, controls and typography reuse mygoui.
This is a native capability mapping, **not full Bklit React/export parity**.
See [reference provenance](references.md) for pinned source and licensing.

Chart constructors return `(ui.Element, error)`; render or propagate errors.
Call `theme.Apply(c)` before charts. Imports do not register fonts, download
geography or start timers. Rendering never mutates caller data.

## Ownership and inspection

Retain one `chart.State` per chart per window. Initialize `Selected: -1` for no
initial selection. A native handle cannot bind two charts in one build. Nil
state makes a chart noninteractive. Retain values and handles, never contexts or
elements. Live samples arrive through the host's `Window.Update`; the host owns
its clock, pause/resume, jobs and cancellation.

Pointer movement inspects measurements. Tab focuses the chart; arrows, Home,
End and Escape change/clear selection. `OnSelect` reports changes including -1;
`OnActivate` handles click/Enter. Readouts are wrapping native text and accessible
values. Disabled/loading charts do not inspect or activate. OS screen-reader
operation requires separate verification.

```go
type dashboard struct {
    selection chart.State // initialize Selected to -1
    samples []chart.Point
}

func (s *dashboard) view(c *ui.Context) {
    theme.Apply(c)
    _, err := chart.Line(c, chart.CartesianOptions{
        Options: chart.Options{Label: "Revenue", State: &s.selection},
        Series: []chart.Series{{ID: "revenue", Label: "Revenue", Points: s.samples}},
    })
    if err != nil {
        ui.Text(c, err.Error())
    }
}
```

Height uses unscaled DIPs (240 default, 48 Sparkline); layout fills native width
and follows text scaling. Skeleton loading honors reduced motion; there is no
entrance delay. React providers and SVG/CSS effects are not Go APIs.

## Chart matrix

| Bklit subject | Native API and supported encoding | Native boundary |
| --- | --- | --- |
| Area | `Area`, multiple series, baseline fill, curves, signed stacks | Solid translucent fill; no pattern/gradient presets or independent Y scales |
| Bar | `Bar`, vertical/horizontal, grouped/stacked, signed baseline | Numeric/category-index X via `Axis.Format`; no square/depth modes |
| Candlestick | `Candlestick`, OHLC wick/body, direction colors, flat bodies | Solid bodies; no volume/pattern layers |
| Choropleth | `Choropleth`, polygon/multipolygon rings, Mercator, sequential values, inspection, zoom/reset, caller pan | No bundled map, GeoJSON parser, graticule or pointer-pan gesture; split antimeridian geometry; latitude within +/-85.0511 degrees |
| Composed | `Composed`, line/area/bar/scatter on shared numeric/time X | Stacking requires aligned bar/area layers only; no independent Y scales |
| Funnel | `Funnel`, stage trapezoids proportional to largest quantity | Vertical straight edges; no horizontal/curved/halo presets |
| Gauge | `Gauge`, arc/linear notches and explicit domain | Fixed 270-degree arc; no separate center statistic or patterns |
| Heatmap | `Heatmap`, sparse category bins, five sequential levels, thresholds, axes, inspection | Host computes calendar bins; no ghost dates, weekday rotation or level legend |
| Line | `Line`, multi-series, linear/monotone/step curves, missing gaps, inspection | Shared Y; no glow, gradient or dashed-tail presets |
| Profit/Loss Line | `ProfitLoss`, interpolated zero crossings, semantic sign colors | No dedicated synchronized sign legend |
| Live Line | `LiveLine`, latest-N supplied samples, dense-line reduction, adaptive domain | Host clock/updates; no autonomous scroll, interpolation or live-tip badge |
| Pie | `Pie`, pie/donut, proportional sectors, gaps, circular hit, share readout | Full circle, solid fill; no partial arc/custom center/patterns |
| Radar | `Radar`, multi-series metric polygons, spokes/grid, fixed/auto maximum | No custom metric label callback or animated polygons |
| Radial | `Radial`, category bars/dots/line/area | Single layer/full circle; no angular ring bars, stacks, partial arc, ghost/now or scrubber |
| Ring | `Ring`, independent value/max progress tracks and arcs | Full circle; no custom center/labels/grid or partial arcs |
| Scatter | `Scatter`, multi-series dots, 2D nearest X inspection | Unique increasing X per series; no jitter/rings; coincident X across series share readout |
| Sankey | `Sankey`, weighted acyclic graph, proportional nodes/ribbons, node/link inspection, connected-flow emphasis | Deterministic layers, no crossing minimization, gradient or drag editing |
| Sparkline | `Sparkline`, line/area/bar, missing gaps, signed bars, optional inspection | Common plot, not special shared-observer engine; native readout when state supplied |
| Sunburst | `Sunburst`, descendant sums, hierarchy partition, click/Enter drill and ancestors | Solid segments; no animated ancestor growth, focus legend or segment labels |

`Point.X` is finite and strictly increasing per series. Unix seconds plus
`Axis.Time` gives dates; `Axis.Format` gives category/host labels. NaN Y is a
missing sample; infinity is rejected. Stacks need aligned X. Empty data displays
`No data` or `EmptyText`. Stable nonempty IDs connect legends to series/values.
Polygon holes use winding opposite to the exterior. Gallery geography is
synthetic, not a copied dataset. Sunburst leaves define sums; branch Value is
ignored with children. IDs are globally unique; cycles/shared children fail.
Stale FocusID is an explicit error; the host decides whether to reset it.

## Utility matrix

| Bklit utility | Native API | Native boundary |
| --- | --- | --- |
| Legend | `Legend`, coss toggles, `State.Hidden` | No percentage/progress subcomponents or coordinated hover |
| Grid | `Grid`, shared Y ticks, `HideGrid` | No shimmer/fade or custom/emphasized ticks |
| Background | `Background`, solid/dots/horizontal lines, Cartesian Background | No SVG pattern/fade presets |
| Reference Area | `Annotation{Kind: ReferenceArea}`, data rectangle | Explicit endpoints, clipped; no bracket/overflow modes |
| Projection Line | `Annotation{Kind: ProjectionLine}`, manual dashed line, domain inclusion | No regression/Bezier/terminal-marker modes |
| Band | `Annotation{Kind: Band}`, varying lower/upper envelope, domain inclusion | No Bollinger helper, independent edges or range tooltip |
| Weekends | `Annotation{Kind: Weekends}`, Saturday/Sunday shading, optional location | Unix X, years 1..9999 and at most 100 years; no categories/horizontal bars/custom weekdays |
| Deadline | `Annotation{Kind: Deadline}`, vertical marker and label | No fading range/category date lookup |
| Trend Line | `Annotation{Kind: TrendLine}`, least squares over finite ordered points | No independent axes or robust regression |
| Markers | `Annotation{Kind: Marker}`, point and label | No grouped dates/fan-out/links/rich popovers |
| Tooltip | `Tooltip`, automatic measurement readout | In-flow native text, not floating/date pill/custom rows |
| Brush | `Brush`, ordered native two-thumb slider, `BrushState` domain filter | No mini-chart, draggable selection body or track overlays |
| Axis | `Axis`, `Axes`, time/format/hidden/ticks, chart-specific matrix/radar geometry | No React wrappers, independent Y or moving tick animations |
| useChart | `PlotContext`, scales/current selection and Build callback | Explicit build-scoped Go callback, no ambient provider |
| Custom Indicator | `Options.Indicator` replaces default selection paint | Paint-only; chart owns input |

Plot fractions run X left-to-right, Y bottom-to-top. Use MyGo ui/chart marks in
Build or painter geometry in Indicator. Do not retain build-scoped UI values.
Draw callbacks only paint; use handled input for actions. Custom marks need
native labels/semantics of their own.

## Scope

Coverage means website chart/utility **subjects**, not every upstream export or
preset. Typed series/options represent supported layers; missing variants are
explicit above. Coss Button/Card/Code already cover shared roles. Marketing fan
circles/shimmer, example blocks, Studio/editor controls, registry installation
and video/export are not a second native application. Proprietary Studio is not
reused. See the [source review and verification](../openspec/changes/archive/2026-10-29-add-bklit-charts/verification.md)
for the upstream inventory and actual evidence.
