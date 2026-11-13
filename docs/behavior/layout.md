# Layout behavior

Package `layout`. Format and evidence marks: [README](README.md). Durations
and easing: [`DESIGN.md`](../../DESIGN.md), "State motion" and "Overlay
motion". The toolkit moves and maximizes the window from a drag region and
reports the window-control insets; a headless test proves which element takes
a press, not that the window moves. The divider of a resizable pane is
`mygoui.ResizablePanel` ([root](root.md)); a compact pane is `overlay.Sheet`
([overlay](overlay.md)). The rows below say what the family adds.

### TitleBar

| | |
| --- | --- |
| Role and name | The window's banner landmark, named by `Label`, then `Title`, then the locale's `TitleBar`. [T TestTitleBarIsTheWindowsBanner] |
| Keys | None of its own; controls inside keep theirs. [T TestTitleBarControlsTakePressesAndTheBarTakesTheRest] |
| Pointer | A press on a control reaches the control; a press anywhere else reaches the bar, which is the window's drag region. [T TestTitleBarControlsTakePressesAndTheBarTakesTheRest] Moving and double-click maximizing are the toolkit's. [C] |
| Focus | The bar is not a Tab stop. Leading controls come before trailing ones. [T TestTitleBarControlsTakePressesAndTheBarTakesTheRest] |
| States | Content stays clear of the window-control insets on the edges `Span` reaches, mirrored right to left. [T TestTitleBarReservesTheInsetsItsSpanReaches] The bar is at least `MinHeight` window DIPs and the reported band, under zoom too; full screen, which reports no insets, releases them. [T TestTitleBarFollowsTheWindowBandAndFullScreen] [T TestTitleBarCoversTheNativeBandUnderZoom] When `Leading` and `Trailing` do not fit between the window controls they move below the band, one frame after the measurement, and return when they fit; a window without controls never stacks. [T TestTitleBarMovesControlsBelowTheBandWhenTheyDoNotFit] |
| Motion | None: a height change applies at once. [C] |

### SplitView

| | |
| --- | --- |
| Role and name | Each pane is a group named by its `Label`; a resizable pane's divider is a splitter with the same name. [T TestSplitViewDividerDragsWithinPaneAndContentBounds] A compact pane is a dialog named by `Label`, with no header and no close button. [T TestSplitViewCompactShowsThePaneOverTheContent] [T TestFlushSheetHandsThePanelToItsBody] |
| Keys | On a divider: arrows move it 16 DIPs toward and away from the content, Home and End go to `Min` and `Max`, Enter collapses a pane that has `Collapsed` and Enter again returns it. [T TestSplitViewDividerKeysResizeCollapseAndReturn] Escape closes a compact pane. [T TestSplitViewCompactShowsThePaneOverTheContent] |
| Pointer | Dragging a divider resizes within `Min` and `Max` and leaves the content `ContentMin` while the view has the room. [T TestSplitViewDividerDragsWithinPaneAndContentBounds] A double click, or a drag below half of `Min`, collapses a pane that has `Collapsed`. [C] A press on the scrim closes a compact pane. [C] |
| Focus | Start pane, its divider, content, end pane's divider, end pane: the order on screen. [T TestResizablePanelHandleFollowsItsEdgeInTabOrder] A collapsed pane's divider stays a Tab stop. [T TestSplitViewDividerKeysResizeCollapseAndReturn] A compact pane takes focus as a Sheet does and returns it to the opener. [T TestPaneToggleFlipsTheFlagOfTheCurrentMode] |
| Selection and value | `Size`, `Collapsed` and `Open` are caller-owned; a zero `Size` shows `Width` until the user moves the divider. [T TestSplitViewDividerDragsWithinPaneAndContentBounds] Without `Max` the splitter's maximum is the room the view had in its last layout, the window's before the first. [T TestSplitViewDividerRangeBeforeTheFirstLayout] A nil `Size` is kept under the view's own identity, so it is lost when elements before an unkeyed view come and go; key the view or pass `Size`. [C] The view clears `Open` whenever it is not compact. [T TestSplitViewCompactShowsThePaneOverTheContent] |
| States | `CompactAuto` is compact below `CompactWidth` scaled with the text size; `CompactNever` and `CompactAlways` force a mode. [T TestSplitViewCompactModeCanBeForced] A compact pane shows only through its caller-owned `Open`, as a Sheet does: a `PaneToggle` or the application's own shortcut sets it; the view adds no key. [C] A `Fixed` pane has no divider. [T TestFixedPaneHasNoDividerAndFoldsOverTime] Start and end swap sides right to left. [T TestSplitViewMirrorsInRightToLeft] The content keeps its focus and state across modes. [T TestSplitViewCompactShowsThePaneOverTheContent] |
| Hairline | On a scaffold's sidebar a pane's edge uses `SidebarBorder`, resizable or fixed; elsewhere `Border`. [T TestScaffoldSidebarHairlineUsesTheSidebarBorder] |
| Motion | A fixed pane folds its width over the reveal duration and is at once folded under reduced motion. [T TestFixedPaneHasNoDividerAndFoldsOverTime] A resizable pane collapses as ResizablePanel does, a compact pane slides as Sheet does. [C] |

### PaneToggle

| | |
| --- | --- |
| Role and name | A toggle button named by `Label`, defaulting to the locale's `ToggleSidebar` or `ToggleInspector`; checked while the pane shows. [T TestPaneToggleFlipsTheFlagOfTheCurrentMode] |
| Keys | Space toggles [T TestPaneToggleFlipsTheFlagOfTheCurrentMode], as does Enter. [C] |
| Selection and value | Flips `Collapsed` while the view docks its panes and `Open` while it is compact. [T TestPaneToggleFlipsTheFlagOfTheCurrentMode] |
| States | Disabled when the flag of the current mode is nil. [C] The icon mirrors right to left. [C] |

### Scaffold

| | |
| --- | --- |
| Role and name | The parts keep their own roles: the title bar is the banner and the start band a group beside it, so the window has one banner [T TestTitleBarIsTheWindowsBanner]; pane groups and splitters. [T TestScaffoldSharesTheTitleBandBetweenPaneAndBar] |
| Focus | Start band, start pane, its divider, title bar, content, end pane's divider, end pane, status bar. [C] |
| Pointer | Dragging the start divider leaves the content its minimum beside a docked end pane, which keeps its size. [T TestScaffoldStartDividerLeavesTheContentItsMinimumBesideTheEndPane] |
| States | The start band reserves the start window controls; the title bar reserves the end ones, and the start ones too while the start pane is collapsed or the view is compact. Both have one height, a stacked bar's included. [T TestScaffoldSharesTheTitleBandBetweenPaneAndBar] [T TestScaffoldBandGrowsWithAStackedBar] The start pane's sheet holds its band, clear of the end window controls too where the sheet reaches them. [T TestScaffoldBandGrowsWithAStackedBar] |
| Motion | While the start pane folds, the title bar reserves what the pane leaves of the start inset in every frame: its controls move with the pane's edge and never pass under the window controls. The bar's height and stacking change at once. [T TestScaffoldTitleBarClearsTheWindowControlsWhileThePaneFolds] |

## Gaps

- TitleBar, platform limit: moving the window from the keyboard is the
  operating system's window menu; the bar adds no key for it.
- SplitView: whether each platform passes the splitter's range on to its
  screen reader is unverified.
