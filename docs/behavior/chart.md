# Chart behavior

Package `chart`. Format and evidence marks: [README](README.md). Loading:
[`DESIGN.md`](../../DESIGN.md), "Loading". Capabilities and data rules:
[`charts.md`](../charts.md).

## Inspection

Shared by every chart below that takes `Options`; the entries list only what
differs.

| | |
| --- | --- |
| Role and name | The plot is a group named by Label and described by Description. [C] While a measurement is selected the group's value is its readout [N TestChartArrowsStepWithoutWrappingAndValueFollows], also shown as a Tooltip status under the plot. [T TestInspectionKeyboardDisabledAndDataReplacement] Without a selection, and on a chart without State, the value spans the data: the first and last readouts joined by an en dash. [T TestChartValueSpansTheDataWithoutASelection] |
| Keys | With State and focus: Right and Down select the next measurement, Left and Up the previous, stopping at the ends; Left or Up with nothing selected takes the last. [N TestChartArrowsStepWithoutWrappingAndValueFollows] Home and End take the first and last; Escape clears the selection. [T TestInspectionKeyboardDisabledAndDataReplacement] PageUp and PageDown move a tenth of the measurements. Horizontal bars run up the plot, so Up moves forward. The plot is not mirrored right to left, so neither are the arrows. One X reads every visible series at it, so there is no separate step between series. [T TestChartKeysFollowTheDrawing] Enter calls OnActivate on the selection [T TestInspectionRespectsTrailingDisabledAncestor], and does nothing without one. [N TestChartArrowsStepWithoutWrappingAndValueFollows] |
| Pointer | Moving the pointer selects the measurement under it. [T TestInspectionRespectsTrailingDisabledAncestor] Leaving the plot clears a selection the pointer made; one made with the keys stays. [T TestChartSelectionFromThePointerEndsWhenItLeaves] A click calls OnActivate on the selection. [C] |
| Focus | With State the plot is one Tab stop; without State it is skipped. [N TestChartArrowsStepWithoutWrappingAndValueFollows] The focus ring shows only for keyboard focus. [T TestChartFocusRingShowsOnlyForKeyboardFocus] |
| Dismissal | Escape clears the selection. With nothing selected it is not taken, so it reaches the dialog or popover around the chart. [T TestChartEscapeReachesTheDialogWithNothingSelected] |
| Selection and value | State.Selected is caller-owned; OnSelect reports each change, including -1. [T TestInspectionKeyboardDisabledAndDataReplacement] A selection beyond new data resets to -1. [T TestInspectionKeyboardDisabledAndDataReplacement] Showing the readout does not move the axes. [T TestSelectionDoesNotMoveAxesAndInheritsDisabled] |
| States | Disabled, or a disabled ancestor, takes no keys, pointer or activation. [T TestInspectionKeyboardDisabledAndDataReplacement] [T TestSelectionDoesNotMoveAxesAndInheritsDisabled] [T TestInspectionRespectsTrailingDisabledAncestor] Loading shows a skeleton and is a busy region that keeps its Tab stop. [T TestLoadingChartIsBusyAndKeepsItsTabStop] Empty data shows EmptyText, else the locale's NoData. [T TestChartTextsComeFromTheLocale] |
| Motion | Selection marks are drawn at once; only the loading skeleton moves. [C] |

## Cartesian

### Line

Also covers: Area, Bar, Scatter, ProfitLoss, Composed, LiveLine, Sparkline,
Axes, Grid, Background, SeriesColor (paint and scale helpers, no behavior of
their own).

| | |
| --- | --- |
| Role and name | Axis tick labels are plain text. [C] The readout reads "label · series: value" [T TestInspectionKeyboardDisabledAndDataReplacement], joined by "; " for each visible series at that X. [C] |
| Keys | The arrows step through the measurements without wrapping. [N TestChartArrowsStepWithoutWrappingAndValueFollows] The measurements are every distinct X of the visible series, once each. [C] |
| Pointer | The nearest X is selected; Scatter takes the nearest point in two dimensions. [C] |
| Selection and value | Hidden series and a Brush narrow what is inspected. [T TestLegendAndBrushDriveVisibleData] LiveLine inspects only the latest samples and follows the caller's updates. [T TestSignedStackAndCallerDataRemainAccurate] Sparkline is inspectable only with State. [C] |

### Candlestick

| | |
| --- | --- |
| Role and name | The readout reads "label: open, high, low, close". [C] |

### Heatmap

| | |
| --- | --- |
| Role and name | The plot is one group with an active readout, as every chart is, rather than a grid of cell elements: the cells are painted, not built, so a large heatmap adds no element per cell and keeps the keyboard model every chart shares. The readout reads "column · row: value", naming the column and row a grid cell would. [C] |
| Keys | Left and Right move across columns, Up and Down across rows (row 0 at the bottom), passing over missing cells and stopping at the edges. [T TestHeatmapArrowsMoveAcrossRowsAndColumns] Home, End and the pages follow Cells in slice order. [C] |

## Polar and proportional

### Pie

Also covers: Ring, Radar, Radial, Funnel.

| | |
| --- | --- |
| Role and name | The readout reads "label: value", with the share in a Pie and every series' value in a Radar. [C] |
| Pointer | A Pie takes the sector under the pointer in a circle. [T TestInspectionRespectsTrailingDisabledAncestor] The selected sector grows. [C] |
| Selection and value | Values hidden through State.Hidden are left out. [C] |

### Gauge

| | |
| --- | --- |
| Role and name | A meter named by Label whose range is Min to Max and whose value is Value, read without a selection. [T TestGaugeIsAMeter] |

## Flow and maps

### Sankey

| | |
| --- | --- |
| Role and name | Nodes are read by label, then links as "source → target: value" by node label. [T TestSankeyRightFollowsTheFlow] |
| Keys | Up and Down step through every node, then every link. Right follows the flow, from a node to its first outgoing link and from a link to its target; Left goes back, to a node's first incoming link or a link's source. [T TestSankeyRightFollowsTheFlow] |

### Choropleth

| | |
| --- | --- |
| Role and name | Regions are read as "label: value". [C] With a Viewport, buttons named by the locale's ZoomIn, ZoomOut and ResetZoom sit above the map. [T TestChartTextsComeFromTheLocale] |
| Pointer | The region under the pointer is selected. [C] The buttons zoom by 1.5 between 1 and 16 and reset. [C] With a Viewport the wheel and a pinch zoom around the pointer, and a drag pans; a drag activates nothing. [T TestChoroplethWheelPinchAndDragMoveTheMap] |
| Selection and value | MapState is caller-owned; the gestures write it. [T TestChoroplethWheelPinchAndDragMoveTheMap] |

### Sunburst

| | |
| --- | --- |
| Role and name | Segments are read as "label: value". [T TestSunburstDrillAndBreadcrumb] With Drill, the path from the root to the current node is a navigation Breadcrumbs: earlier nodes are links named by label and the current node is its last item, marked current, not a link. [T TestSunburstEscapeGoesUpALevel] |
| Keys | Enter on a branch drills into it [T TestSunburstDrillAndBreadcrumb], clears the selection and keeps the focus on the chart. [N TestSunburstEnterDrillKeepsFocusOnChart] Escape clears the selection, and with nothing selected goes up a level. [T TestSunburstEscapeGoesUpALevel] |
| Pointer | A click on a branch drills. [C] A path link returns to that node. [T TestSunburstDrillAndBreadcrumb] |
| Selection and value | SunburstState.FocusID is caller-owned; an ID not in the tree is an error. [C] |

## Utilities

### Tooltip

| | |
| --- | --- |
| Role and name | A status region in the layout, read when its text changes. [C] |

### Legend

| | |
| --- | --- |
| Role and name | With State, each item is a toggle button named by its label, pressed while the series shows; with nil State, plain text. [C] |
| Pointer | A toggle hides or restores its series [T TestLegendAndBrushDriveVisibleData] and clears the selection. [C] |
| Focus | Each toggle is a Tab stop. [C] |
| States | Disabled disables every toggle. [C] |

### Brush

| | |
| --- | --- |
| Role and name | The form range slider named by the locale's Range [T TestChartTextsComeFromTheLocale]; form names its thumbs. [C] |
| Keys | As form's Slider (toolkit stepping). [C] |
| Selection and value | BrushState is caller-owned and inclusive; points on either bound show. [T TestBrushKeepsPointsWithinItsBounds] A nil state or a non-increasing domain is an error. [C] |
| States | Disabled ignores keys and drags. [C] |

## Gaps

Measured against platform conventions and the
[APG patterns](https://www.w3.org/WAI/ARIA/apg/patterns/). None remains open.
