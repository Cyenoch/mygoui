# Navigation behavior

Package `navigation`. Format and evidence marks: [README](README.md). Durations
and easing: [`DESIGN.md`](../../DESIGN.md), "State motion" and "Loading".
Focus-group roving (arrows, wrap, Home/End, the single Tab stop and where Tab
enters), popover Escape and outside-press dismissal, OS-native menus and
native scrolling belong to the toolkit; the rows below say which pattern a
component uses rather than repeating its keys.

## Disclosure

### Accordion

| | |
| --- | --- |
| Role and name | A group named by `Label`; each header is a disclosure named by its title, exposing Expanded. [C] |
| Keys | Down and Up move between the enabled headers, round the ends, and Home and End go to the first and last. [T TestAccordionSingleExpansionAndKeyboardFocus] [N TestHeaderArrowsLeaveTabToWhatLiesBetween] Enter toggles the focused header [T TestAccordionSingleExpansionAndKeyboardFocus], as does Space. [C] |
| Pointer | Clicking a header toggles it and focuses it. [T TestAccordionSingleExpansionAndKeyboardFocus] Hover leaves the header transparent. [T TestAccordionHeaderHoverPaint] |
| Focus | Each header is a Tab stop, and Tab goes from a header into its actions and the controls of its open panel. [N TestHeaderArrowsLeaveTabToWhatLiesBetween] |
| Selection and value | `Expanded` is caller-owned, local when nil. Single mode keeps one item open; Multiple opens several. [T TestAccordionSingleExpansionAndKeyboardFocus] [T TestAccordionMultipleExpansion] `OnChange` receives the key and the new state. [C] |
| States | Disabled headers take no clicks or focus. [C] |
| Motion | The panel grows and fades with its height, a reversal keeps its progress, and reduced motion opens at once; the chevron turns with the progress. [T TestAccordionPanelGrowsContinuously] [T TestAccordionReversalKeepsItsProgress] [T TestAccordionReducedMotionOpensAtOnce] [T TestAccordionPanelFadesWithItsHeight] [T TestDisclosureAndAccordionChevronsAreDrawn] |

### Disclosure

| | |
| --- | --- |
| Role and name | The trigger is a disclosure named by `Title`, exposing Expanded; `Actions` sit beside it, outside the trigger. [C] |
| Keys | Enter and Space toggle. [N TestDisclosureTogglesWithEnterAndSpace] |
| Pointer | A click toggles; hover paints no face. [T TestDisclosureHoverIsInert] |
| Focus | The trigger is its own Tab stop with a focus ring. [C] |
| Selection and value | `Open` is caller-owned. [N TestDisclosureTogglesWithEnterAndSpace] A nil `Open` keeps the state in the disclosure. [C] |
| States | Disabled takes no clicks or focus. [C] |
| Motion | The panel grows and fades like Accordion's; the chevron is drawn turning. [T TestDisclosureAndAccordionChevronsAreDrawn] |

### DisclosureGroup

| | |
| --- | --- |
| Role and name | A group named by `Label`. [C] |
| Keys | Up and Down move between the triggers, round the ends, and Home and End go to the first and last; each trigger stays a Tab stop. [N TestHeaderArrowsLeaveTabToWhatLiesBetween] |
| Selection and value | Without `Multiple`, opening one closes the others, and at most one stays open, preferring the one just changed. [T TestDisclosureGroupSingleExpansion] |

## Views

### Tabs

| | |
| --- | --- |
| Role and name | A tab list with tabs exposing Selected; the panel is a tab panel named and controlled by the selected tab. [N TestTabsPanelIsATabPanelOfTheSelectedTab] |
| Keys | Toolkit tab list: Left/Right and Up/Down move, wrap and skip disabled tabs; Home/End go to the ends; selection follows focus. [T TestTabsSelectionAndPanelSwitch] [T TestTabsKeyboardSkipsDisabledTab] |
| Pointer | A click selects and focuses the tab. [T TestTabsSelectionAndPanelSwitch] Hover dims the label. [T TestTabsHoverDimsLabel] |
| Focus | Tab enters at the selected tab and goes on into the panel (probed). [C] With `KeepMounted`, hidden panels keep their local state and are not Tab targets. [T TestTabsKeepMountedPreservesLocalStateWithoutHiddenTargets] |
| Selection and value | `Selected` is caller-owned. A disabled or missing selection falls back to the nearest enabled tab without `OnChange`, and is restored when the tab returns or becomes enabled. [T TestTabsReconcilesDisabledSelection] [T TestTabsKeepsSelectionUsableWhenItemDisappears] [T TestTabsRestoresSelectionWhenItemsReturn] [T TestTabsRestoresSelectionWhenAnItemBecomesEnabled] |
| States | Disabled tabs take no clicks and are skipped. [T TestTabsKeyboardSkipsDisabledTab] |
| Motion | The indicator slides and settles, keeps its course on reversal, jumps under reduced motion and follows relayout; label ink fades to the chosen tab. [T TestTabsIndicatorSlidesAndSettles] [T TestTabsIndicatorSettlesAfterReversal] [T TestTabsIndicatorReducedMotion] [T TestSelectionIndicatorsFollowRelayout] [T TestTabLabelInkFadesToTheChosenTab] |

### SegmentedControl

| | |
| --- | --- |
| Role and name | A radio group named by `Label`; segments are radios named by their labels. [C] |
| Keys | Toolkit radio group: the arrows move and select, skipping disabled segments. [T TestSegmentedControlChangesSelectionWithoutBuildingPanels] |
| Pointer | Selection changes on release; a disabled segment ignores clicks. [T TestSegmentedControlUpdatesSelectionBeforeReturningOnRelease] [T TestSegmentedControlChangesSelectionWithoutBuildingPanels] |
| Selection and value | `Selected` is caller-owned; item content is never built. [T TestSegmentedControlChangesSelectionWithoutBuildingPanels] |
| Motion | The indicator slides, jumps under reduced motion and remeasures changed labels; label ink fades. [T TestSegmentedControlIndicatorMovesAndRespectsReducedMotion] [T TestSegmentedIndicatorRemeasuresChangedLabels] [T TestSegmentLabelInkFadesToTheChosenSegment] |

### DocumentTabs

| | |
| --- | --- |
| Role and name | A tab list named by `Label`; tabs expose Selected and a dirty tab is described "Modified". [T TestDocumentTabsExposeTabsAndReportSelection] Each close slot is a button "Close <title>". [T TestDocumentTabsCloseByButtonKeysAndMiddleClick] |
| Keys | The arrows move focus only (manual activation). [T TestDocumentTabsCloseKeyActsOnTheActiveTab] Enter or Space selects. [C] Delete closes the focused unpinned tab and Cmd/Ctrl+W the active one, never a pinned one. [T TestDocumentTabsCloseByButtonKeysAndMiddleClick] [T TestDocumentTabsCloseKeyActsOnTheActiveTab] Ctrl+Tab and Ctrl+Shift+Tab step with wrap, from the ends when no tab is active [T TestDocumentTabsKeysChooseTabsAndCloseTheActiveOne] [T TestDocumentTabsWithoutAnActiveTabStepFromTheEnds]; Ctrl+PageDown and Ctrl+PageUp do the same. [C] With `WindowShortcuts` these keys work window-wide and leave focus where it is. [T TestDocumentTabsWindowShortcutsKeepTheEditorFocused] Alt+Left/Right move a tab within its pinned or unpinned group, and focus follows it. [T TestDocumentTabsReorderByKeysAndDrag] |
| Pointer | A click selects; the close button closes without selecting; a middle click released on the same tab closes it. [T TestDocumentTabsCloseByButtonKeysAndMiddleClick] Dragging reorders. [T TestDocumentTabsReorderByKeysAndDrag] The wheel scrolls the strip sideways. [T TestDocumentTabsWheelScrollsTheStripSideways] The context menu offers Close, Close others and Close to the right, never closing pinned tabs. [T TestDocumentTabsContextMenuClosesTheOthersAndTheTabsAfterButNotPinned] |
| Focus | Close slots are not Tab stops. [C] When the caller removes the focused tab, focus moves to the tab that takes its place. [T TestDocumentTabsFocusMovesToTheNeighbourOfAClosedTab] |
| Selection and value | `Tabs` and `Active` are caller-owned and never modified; callbacks report each request. [T TestDocumentTabsExposeTabsAndReportSelection] |
| States | Overflow shrinks tabs before scrolling and shows an "All tabs" menu; the active tab scrolls into view. [T TestDocumentTabsShrinkToFitBeforeScrolling] [T TestDocumentTabsMenuListsEveryTabAndScrollsToTheChosenOne] [T TestDocumentTabsKeepTheActiveTabInView] Pinned tabs show an icon without a close slot and stay first. [T TestDocumentTabsPinnedTabsAreIconsWithoutCloseAndStayFirst] |
| Motion | The active face fades; tabs slide, grow on arrival and leave when closed; reduced motion keeps only the fades. [T TestDocumentTabActiveFaceFadesAndClosedTabLeaves] [T TestDocumentTabsCloseAndAddAnimate] [C] |

### NavigationStack

| | |
| --- | --- |
| Role and name | A group named by `Label`; the header holds a Back button whose tooltip names the parent page and the title as a heading. Pages are groups named by their titles. [C] |
| Keys | Alt+Left (Alt+Right in right-to-left) and the mouse Back button pop. Escape is left to the focused editor. [T TestNavigationStackPushesAndPopsWithFocus] [T TestNavigationStackEscapeLeavesTheEditorItsEditing] |
| Pointer | The Back button pops. [T TestNavigationStackPushesAndPopsWithFocus] |
| Focus | When focus was inside the stack it moves to the new page. [T TestNavigationStackPushesAndPopsWithFocus] Otherwise the new title is announced. Covered pages are hidden from assistive technology and input. [C] |
| Selection and value | `Pages` is caller-owned; `OnPop` asks the caller to remove the top page. `KeepMounted` keeps covered pages' state. [T TestNavigationStackKeepMountedPreservesCoveredPageState] |
| Motion | Pages slide; reduced motion changes them at once. [T TestNavigationStackSlidesUnlessMotionIsReduced] |

## Trails and paging

### Breadcrumbs

Also covers: Breadcrumb.

| | |
| --- | --- |
| Role and name | A group named by `Label`, or by the locale's Breadcrumbs without one. [T TestUnnamedBreadcrumbsTakeTheLocalesName] Crumbs before the last are links; the last is plain text marked current page. [N TestBreadcrumbsFollowLinksFromTheKeyboard] |
| Keys | Each link is a Tab stop followed with Enter. [N TestBreadcrumbsFollowLinksFromTheKeyboard] |
| Pointer | A click follows an enabled link; ink darkens on hover and focus. [C] |
| States | With `MaxVisible` of 3 or more, middle crumbs move into a "More pages" menu button whose OS menu reports their own index. [N TestBreadcrumbsFollowLinksFromTheKeyboard] Separators mirror in right-to-left. [C] Disabled crumbs take no clicks. [C] |

### Pagination

| | |
| --- | --- |
| Role and name | A group named by `Label`; Previous and Next carry locale labels, pages are named by number, the current one marked current page [T TestPaginationSelectedPageHasVisibleDarkSurface], the ellipsis is text named "More pages". [C] |
| Keys | Every button is a Tab stop, pressed with Enter or Space. [C] |
| Focus | Previous or Next disabling while it holds the focus hands it to the current page. [N TestPaginationHandsTheFocusToTheCurrentPageAtAnEnd] |
| Pointer | A click chooses a page; Previous and Next are disabled at the ends. [T TestPaginationBoundariesAndEllipsis] |
| States | In a window narrower than 640 DIP, Previous and Next show icons only and keep their names; the summary stacks below. [T TestPaginationNarrowControlsKeepAccessibleLabels] [T TestPaginationSummaryLayout] The current page has a dark surface. [T TestPaginationSelectedPageHasVisibleDarkSurface] |
| Motion | The active face fades. [T TestPaginationActivePageFaceFades] |

### PipsPager

| | |
| --- | --- |
| Role and name | A group named by `Label`; pips are buttons "Page N" with their position in the set, the current one marked current page. [T TestCarouselBuildsOnlyTheSlidesInView] |
| Keys | One horizontal focus group: the arrows move focus without choosing; Enter chooses. [T TestPipsPagerOwnsItsIndexWithoutACaller] [T TestCarouselBuildsOnlyTheSlidesInView] |
| Pointer | A click chooses a pip. [T TestPipsPagerOwnsItsIndexWithoutACaller] |
| Selection and value | The index is caller-owned, local when nil. At most `MaxVisible` pips (7) show, windowed round the current one, with shrunken edge pips. [T TestPipsPagerOwnsItsIndexWithoutACaller] [C] |
| Motion | Pip size and ink change over state motion; the current pip darkens. [T TestPipsPagerCurrentPipDarkens] |

### Carousel

| | |
| --- | --- |
| Role and name | The viewport is a focusable group named by `Label`; slides are groups with their position in the set. [T TestCarouselJumpsWithHomeAndEndAndSaysItsPlaceAmongManyPages] The viewport's role reads as the locale's Carousel ("carousel") and each slide's as its Slide ("slide"). [T TestCarouselTellsItsRolesInTheLocale] An "i / n" counter, hidden from assistive technology, shows past 7 pages. [T TestCarouselJumpsWithHomeAndEndAndSaysItsPlaceAmongManyPages] |
| Keys | Left/Right step and Home/End jump while focus is in the carousel; right-to-left mirrors the arrows [C]. [T TestCarouselStepsWithButtonsAndArrowKeys] [T TestCarouselJumpsWithHomeAndEndAndSaysItsPlaceAmongManyPages] |
| Pointer | Previous and Next step and are disabled at the ends unless `Loop`. A drag past the threshold or a fling turns one page, a shorter drag springs back, and vertical drags go to the page; a sideways wheel or swipe turns one page per gesture. [T TestCarouselStepsWithButtonsAndArrowKeys] [T TestCarouselDragTurnsOnePageOnlyPastTheThreshold] [T TestCarouselSidewaysSwipeTurnsOnePage] |
| Focus | Next disabling at the last slide hands the focus to the slides, whose arrows step on; each step announces the page reached. [N TestCarouselKeepsTheFocusAndAnnouncesThePageAStepReaches] |
| States | Only slides in view are built. [T TestCarouselBuildsOnlyTheSlidesInView] The arrows are centred on the slides. [T TestCarouselArrowsAreCenteredOnTheSlides] |
| Motion | Slides translate; reduced motion changes them at once. [C] |

### Steps

| | |
| --- | --- |
| Role and name | A list named by `Label`; each step is a list item named by its title with its description, the current one marked current step. [T TestStepsShowCompletedCurrentAndUpcoming] |
| Keys | With `OnSelect`, each step head is a button and Tab stop pressed with Space [T TestStepsSelectReportsThePressedStep] or Enter. [C] The arrows along its direction, mirrored right to left, and Home and End move between the steps. [N TestStepsArrowsMoveAlongPressableSteps] |
| States | Completed steps show a check, failed ones a cross. [C] |
| Motion | The face fades and the tick draws in. [T TestStepsFaceFadesAndTickDrawsIn] |

## Menus

### Menu

Also covers: SubMenu.

| | |
| --- | --- |
| Role and name | A menu; items are menu items, radio items in single selection, or check items when checkable, with description and shortcut text. Sections are groups named by their titles; separators are hidden. [C] |
| Keys | One vertical focus group (arrows wrap, Home/End). Type-ahead focuses the next enabled item starting with the letters typed and cycles on a repeated letter. Right opens a submenu, which focuses its first item; Left closes it. [T TestMenuTypeAheadFocusesTheItemTypedAndSkipsDisabledOnes] [T TestMenubarSubmenuKeepsItsArrows] [T TestDropdownKeyboardSubmenu] Menus of 32 rows or more are virtual lists that keep Up/Down/Home/End with wrap and reach every item. [T TestDropdownLongMenuKeyboardReachesEveryItem] [T TestDropdownLongMenuScrollsToLastItem] Type-ahead reaches their rows not yet built; letters typed within a second build a prefix. [N TestLongMenuTypeAheadReachesRowsNotYetBuilt] |
| Pointer | A click chooses an enabled item and closes the menu that holds it. [T TestDropdownActivationAndSubmenu] |
| States | Disabled items take no focus. [T TestMenuTypeAheadFocusesTheItemTypedAndSkipsDisabledOnes] |
| Motion | The check draws in and fades out. [T TestMenuCheckDrawsInAndFadesOut] |

### Dropdown

| | |
| --- | --- |
| Role and name | The trigger is a menu button named by `Label`, exposing Expanded. [C] |
| Keys | Enter or Space on the trigger opens the menu at its first item [C]; Down opens it at its first item and Up at its last. [N TestMenuButtonsOpenFromTheArrows] |
| Pointer | A click toggles; choosing an item closes. [T TestDropdownActivationAndSubmenu] |
| Focus | The menu keeps focus while it grows long. [T TestDropdownKeepsFocusWhenItGrowsLongWhileOpen] |
| Dismissal | Escape closes. [T TestDropdownEscapeCloses] An outside press closes and leaves the focus where it landed; Escape and choosing return it to the trigger. [N TestMenuButtonsTakeTheFocusBackAsTheirMenuCloses] |
| Motion | The menu is open and usable once the entry settles, with and without reduced motion. [T TestDropdownEntryFade] |

### ContextMenu

| | |
| --- | --- |
| Role and name | An OS-native context menu; section titles and icons are not shown. [C] |
| Keys | Shift+F10 and the menu key open it on the focused child, which keeps its focus. [T TestContextMenuOnInteractiveChild] |
| Pointer | A secondary click opens it at the pointer and still reaches an interactive child; text selection is kept. [T TestContextMenuOpensAtNativeTriggerAndActivates] [T TestContextMenuOnInteractiveChild] [T TestContextMenuKeepsTextSelection] |
| Selection and value | Single and multiple selection toggle the caller's state, local when nil. [T TestContextMenuTogglesCallerSelection] |
| States | Disabled items cannot be chosen. [T TestContextMenuOnInteractiveChild] The anchor can leave the tree and return. [T TestContextMenuAnchorCanLeaveAndReturn] |

### PopoverMenu

| | |
| --- | --- |
| Role and name | A drawn in-window menu with descriptions and shortcuts. [T TestPopoverMenuShowsDescriptionAndShortcut] |
| Keys | Shift+F10 and the menu key open it; Down skips disabled items; Enter chooses. [T TestPopoverMenuKeyboardSelectionAndFocus] |
| Pointer | A secondary click opens it; disabled items ignore clicks; submenus open on click. [T TestPopoverMenuPointerSubmenuAndDismissal] |
| Dismissal | Escape dismisses. [T TestPopoverMenuPointerSubmenuAndDismissal] Choosing closes and returns focus to the anchor. [T TestPopoverMenuKeyboardSelectionAndFocus] |

### Menubar

| | |
| --- | --- |
| Role and name | A menu bar named by `Label`; titles are menu items exposing Expanded. [T TestMenubarRolesForAssistiveTechnology] |
| Keys | One horizontal focus group. Enter, Space or Down open the focused title at its first item (probed). [C] With a menu open, Left/Right switch to the neighbouring enabled menu with wrap; Escape closes and refocuses the title. [T TestMenubarArrowsSwitchMenusAndEscapeRestoresFocus] Accelerators run enabled items with the menus closed. Off macOS, F10 or a lone Alt tap focuses the bar and Alt+letter opens a menu, with mnemonics underlined while the keyboard drives the bar. [T TestMenubarAcceleratorActivatesItsItemWithTheMenuClosed] [T TestMenubarMnemonicsAndF10FollowWindowsConventions] [T TestMenubarAltTapFocusesTheBarButAltWithAKeyDoesNot] [T TestMenubarAltWithAnotherModifierDoesNotFocusTheBar] |
| Pointer | A click toggles a menu; moving over another title switches an open menu; choosing closes and refocuses the title. [T TestMenubarHoverSwitchesAnOpenMenuAndChoosingCloses] |
| Dismissal | An outside press closes without refocusing. [C] Tab from an open menu closes it. [N TestMenubarTabClosesTheOpenMenu] |
| States | Disabled titles take no clicks or focus. [C] |

### SplitButton

| | |
| --- | --- |
| Role and name | A group named by `Label`, holding the action button and an arrow menu button named by `MenuLabel` (locale "More" by default) exposing Expanded. [C] |
| Keys | Each half is a Tab stop. Enter on the arrow opens the menu at its first item. [N TestSplitButtonEscapeReturnsFocusToTheArrow] Down on the closed arrow opens the menu at its first item and Up at its last. [N TestMenuButtonsOpenFromTheArrows] |
| Pointer | The action runs on click; the arrow toggles the menu. [T TestSplitButtonActsAndOffersAlternatives] |
| Dismissal | Escape closes and returns focus to the arrow. [N TestSplitButtonEscapeReturnsFocusToTheArrow] Choosing an item does too. [N TestMenuButtonsTakeTheFocusBackAsTheirMenuCloses] |
| Selection and value | `Selected` remembers the alternative chosen last. [T TestSplitButtonRemembersTheAlternativeChosenLast] |
| States | `Loading` stops the action and keeps the menu usable. [T TestSplitButtonLoadingStopsActionAndKeepsMenu] |

## Command and bars

### Command

| | |
| --- | --- |
| Role and name | A search field named by `Label` (locale "Search" by default) with its placeholder, controlling the list below it. [N TestCommandInputIsTheOneTabStopAndEscapeClearsTheQuery] The highlighted row is the field's active descendant. [T TestCommandInputPointsAtTheHighlightedRow] |
| Keys | In the field, Down/Up move the highlight, clamped at the ends and skipping disabled items; Enter activates the highlighted item. [T TestCommandArrowsSkipDisabledActions] [T TestCommandFiltersAndActivatesFromKeyboardAndPointer] Escape empties a query; with none it goes on to the surface holding the command. [N TestCommandInputIsTheOneTabStopAndEscapeClearsTheQuery] |
| Pointer | A click activates an item; pointer and keyboard share one highlight. [T TestCommandFiltersAndActivatesFromKeyboardAndPointer] [T TestCommandPointerAndKeyboardShareHighlight] |
| Focus | The field is the one Tab stop; Tab passes over the list. [N TestCommandInputIsTheOneTabStopAndEscapeClearsTheQuery] |
| Selection and value | `State` is required and caller-owned. Filtering is a case-insensitive substring match on label and keywords; groups become sections. [T TestCommandFiltersAndActivatesFromKeyboardAndPointer] [T TestCommandFiltersGroupedActions] Selection can clear the query, and the highlight follows new items. [T TestCommandSelectionCanClearQuery] [T TestCommandFollowsQueryNewItemsAndRevision] |
| States | `Loading` marks the root busy, keeps the field editable and makes results inert. [T TestCommandLoadingKeepsInputEditableAndResultsInert] |

### Toolbar

Also covers: ToolbarSeparator.

| | |
| --- | --- |
| Role and name | A toolbar named by `Label`, exposing its orientation; separators are separators running across it. [N TestToolbarExposesItsOrientationAndSeparators] |
| Keys | One focus group along its orientation: the arrows move and wrap, Home/End go to the ends, the cross-axis arrows do nothing. [N TestToolbarIsOneTabStopWithWrappingArrows] |
| Focus | Tab enters and leaves once and returns to the control used last. [N TestToolbarIsOneTabStopWithWrappingArrows] |

### OverflowList

| | |
| --- | --- |
| Role and name | A toolbar named by `Label`; items are buttons, or toggle buttons when `Checked` is set; the "More" button is a menu button. [C] |
| Keys | One horizontal focus group including the "More" button; in its menu, End and Enter choose. [T TestOverflowListCollapsesWhatDoesNotFitIntoItsMenu] Escape in the menu returns focus to "More" (probed). [C] |
| States | Items that do not fit first drop their labels (kept as name and tooltip), then move into the menu by priority from the collapsing end; pinned items stay. Collapsed items leave the accessibility tree; toggles become check items. [T TestOverflowListCollapsesWhatDoesNotFitIntoItsMenu] [T TestOverflowListCollapsesFromTheStart] [T TestOverflowListKeepsPinnedAndHighPriorityItemsAndShowsToggles] [T TestOverflowListMenuChecksOnlyItsToggles] |

### StatusBar

| | |
| --- | --- |
| Role and name | A group named by `Label`; the message is a status region announced on change. Items with an action are buttons named by label and text; others are plain; progress items are progress bars. [T TestStatusBarAnnouncesMessageChanges] [T TestStatusBarItemsPressOnlyWithAnAction] |
| Keys | Action items form one horizontal focus group; Enter presses. [T TestStatusBarItemsPressOnlyWithAnAction] |
| Focus | Focus stays with its item when an earlier one leaves. [T TestStatusBarFocusStaysWithItsItemWhenAnEarlierOneLeaves] |
| States | Items drop by priority while the message keeps its minimum width. [T TestStatusBarDropsLowPriorityItemsAndKeepsTheMessage] |
| Motion | The message crossfades. [T TestStatusBarMessageCrossfades] |

## Gaps

Measured against platform conventions and the
[APG patterns](https://www.w3.org/WAI/ARIA/apg/patterns/). None remains open.
