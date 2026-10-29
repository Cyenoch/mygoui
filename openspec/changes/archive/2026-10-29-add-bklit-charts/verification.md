# Verification: add-bklit-charts

2026-10-29. Native capability mapping, not full Bklit React/export parity.
The small EmptyState alignment correction preserves the existing centered
empty-state contract and is verified alongside this change.

## Source inventory and boundary

- `references/bklit` pins `bklit/bklit-ui` at
  `ec5a9705ee973f95fd0ef2d974e9d00ae3fda011`. Catalog discovery used public MDX,
  `apps/web/content/docs` metadata and `packages/ui/src/charts/index.ts` together;
  README and llms lists alone omit subjects.
- There are 19 documented chart subjects: area, bar, candlestick, choropleth,
  composed, funnel, gauge, heatmap, line, profit/loss line, live line, pie, radar,
  radial, ring, scatter, Sankey, sparkline and sunburst.
- Utilities comprise legend, grid, background, reference area, projection line,
  band, weekends, deadline, trend line, markers, tooltip, brush, useChart and
  custom indicator. Axis has overview, X and Y documentation. Gallery maps this
  to 19 chart pages and 16 utility/axis pages, all with public Go snippets.
- The chart library is MIT, credited to 2026 uixmat; exact retained text is in
  `licenses/bklit-ui-NOTICE.md`. Studio is proprietary and `packages/icons` is
  UNLICENSED. Neither implementation is reused. No geographic, font or vendored
  pattern asset is copied. Gallery geography is synthetic caller geometry.
- `docs/charts.md` is the per-subject subset/gap matrix. React providers, CSS/SVG
  presets, registry installation, marketing effects and Studio/editor/export
  applications are not implemented as native components. Shared Button/Card/Code
  roles reuse the incumbent coss components.

## Requirement and scenario mapping

| Requirement | Implementation | Scenario evidence |
| --- | --- | --- |
| Chart catalog has native counterparts | `chart/cartesian.go`, `chart/candlestick.go`, `chart/polar.go`, `chart/heatmap.go`, `chart/choropleth.go`, `chart/sankey.go`, `chart/sunburst.go`; Gallery catalog and `docs/charts.md` | `TestChartCatalogRendersEveryPublicExample` renders every chart/utility page; `TestComponentSnippetsTypeCheck` compiles public snippets |
| Data boundaries are explicit | Chart constructors validate numbers, OHLC, polygon coordinates, DAGs and single-parent hierarchies; shared empty surface | `TestConstructorsRejectMalformedData` includes cyclic flow; `TestInspectionKeyboardDisabledAndDataReplacement` checks empty/replaced samples; `TestSignedStackAndCallerDataRemainAccurate` checks caller data is not mutated |
| Inspection preserves native input and ownership | `chart/chart.go` caller State/Handle, committed native pointer/shortcut actions, readout and value metadata | Keyboard/disabled/data-replacement tests; stable-axis test; trailing-disabled-ancestor test includes Enter and pointer; real macOS line focus/End exposes `Day 6` and Sales 80 |
| Charts compose with semantic theme and utilities | Shared coss tokens, native axes, Legend/Brush, annotations and build-scoped PlotContext/Indicator | `TestCustomPlotUsesHostThemeAndAnnotationDomain` checks a changed Chart1 token reaches native paint, Band upper domain, selection and custom indicator; legend/brush test and utility Gallery renders |
| Streaming and hierarchy actions are caller controlled | `LiveLine` trims caller samples without a worker; `Sunburst` partitions/drills/navigates; `Choropleth` consumes polygons and caller viewport | Live test updates caller samples and inspects the new window; `TestSunburstDrillAndBreadcrumb` verifies sums, Enter and ancestor return; malformed polygon test plus Gallery geometry rendering; source inspection finds no chart I/O or job startup |

## Review findings and resolution

Independent source inventory review confirmed catalog and licensing scope.
Independent code review found the following substantive issues; fixes and focused
regressions are in `chart/chart_test.go`:

- Huge finite Ring/Sunburst quantities overflowed when multiplied before division,
  potentially hanging arc construction. Ratios now precede angular multiplication.
  Aggregate Pie/Sankey overflow fails explicitly.
- Weekend shading could iterate an unbounded calendar domain. Domains now fit
  years 1..9999 and at most 100 years, with an error regression.
- An inspection readout inside Plot's center cell shifted the axes. `plotFrame`
  keeps readouts outside that grid cell; the zero tick stays fixed on selection.
- Ancestor-disabled charts could inspect via polled pointer state. Native committed
  event routing now handles pointer and keyboard actions, including modifiers set
  after Children. The final local follow-up also defers Enter activation, with a
  regression confirming exactly one enabled activation and none while disabled.
- Other resolved findings: NaN gaps in stacked areas, nearest angular selection
  for radar/dots and forward ribbons in long Sankey chains.

Independent finish-design review found six actionable issues: endpoint bar/candle
clipping, radar label bounds, sunburst center measurement, missing pie/radar
legends and heatmap key, Band upper domain exclusion, and ambiguous example data.
Automatic half-step domain padding, measured text, native legends/key, annotation
extent inclusion and `Sample data` example titles resolve them. The follow-up
verdict is ship for those six fixes, not certification of every upstream preset.

Read-only design documentation review found coss palette, typography, layout and
control inheritance consistent. Root DESIGN.md and its sidecar remain unchanged.
The EmptyState action slot now wraps content in a centered row: the Button's
AlignSelf(Start) no longer defeats column centering. Its focused regression uses
280/640 DIP widths, text scales 1/1.3 and one callback per click.

## Evidence

Engineering commands, final outcomes and platform limits are recorded in the
`Bklit native charts and EmptyState alignment` section of `VALIDATION.md`.
Final check/lint/WASM, full race and the changed chart/Gallery race follow-up
pass. Strict spec validation and whitespace checks pass.

- Native headless captures: `artifacts/charts/*-contact.png` for 13 representative
  chart/utility subjects, wide/compact, light/dark; compact dark uses 1.3 text scale.
- EmptyState captures: `artifacts/empty-state-fix/emptystate-action-contact.png`.
- Real macOS captures: `artifacts/charts/native-macos-line-inspected.png`,
  `artifacts/charts/native-macos-sunburst-drilled.png` and
  `artifacts/charts/native-macos-empty-state.png`. All three were visually inspected.
  The latter shows centered action and one created search after activation.
- Native captures and headless paint do not establish other OS/GPU backends,
  physical screen-reader operation or IME behavior. No such certification is made.

## Final assessment

| Dimension | Result |
| --- | --- |
| Completeness | Ten tasks completed; all five ADDED requirements implemented and six scenarios mapped |
| Correctness | Focused behaviors, 35 Gallery renders, snippet compilation, engineering checks and real macOS smoke checks pass |
| Coherence | Typed caller-owned native APIs, coss inheritance and documented subset scope follow the design; no runtime dependencies or chart import-side jobs/I/O |

No unresolved critical findings, warnings or code-pattern suggestions remain in
the checks performed. Windows/Linux execution, other GPU backends, IME and actual
screen-reader operation are explicitly outside this verification, not inferred
from headless or metadata results. These platform limits are not skipped semantic
requirements: the spec does not assert platform certification. The change is
ready for synchronization/archive within the documented native subset.
