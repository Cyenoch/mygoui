# Overlay behavior

Package `overlay`. Format and evidence marks: [README](README.md). Durations
and easing: [`DESIGN.md`](../../DESIGN.md), "Overlay motion", "State motion"
and "Overlays" (inert exit paint). Native flyout windows: [`docs/flyout.md`](../flyout.md).

The toolkit owns, for the surfaces built on its bases: the modal focus trap
and the return of focus to the opener (`ui.DialogBase`), Escape and the
outside press of anchored popups and their return of focus to the anchor
(`ui.PopoverBase`), tooltip timing (`ui.TooltipBase`), toast lifetime,
timers and Escape on a toast, `FocusGroup` arrow movement, and native
flyout windows. An overlay receives the keys its focused content leaves
unhandled, the overlay built last first: a focused search field takes
Escape before its dialog, and a focused button takes Enter. [C]

Shared by the surfaces that take `PreventDismiss` (Modal, Dialog, Sheet,
Drawer, BottomSheet; AlertDialog, AboutDialog, ShortcutsDialog and
ImageViewer have no such option): it blocks Escape and the backdrop, keeps
the surface open and keeps the focus where it was. [N TestPreventDismissHoldsTheSurfaceAndItsFocus]
The test builds Modal, Dialog, Drawer and BottomSheet; Sheet passes the
option to Dialog. [C]
`OnDismiss` runs after the surface set `*open` to false. [C]
A closing surface paints its exit from a retained copy, not from rebuilt
content [T TestSurfacesPaintArrivalAndDeparture], and that copy takes no
input. [T TestEdgeSurfacesSlideAndRetainInertExits]
Under reduced motion Modal, BottomSheet and Tooltip arrive at rest
[T TestReducedMotionSurfacesArriveAtRest]; the dialog, panel and scrim
fades are skipped too. [C]

## Modal dialogs

### Modal

| | |
| --- | --- |
| Role and name | A dialog named by `Title` and described by `Description`. [C] The close button is named by the locale's Close and is absent with `HideClose`. [C] |
| Keys | Escape closes. Tab cycles inside the panel; footer actions stay reachable by Tab in a short window. [T TestModalTrapsTheFocusAndRestoresIt] [T TestShortPanelsKeepFooterActionsReachable] |
| Pointer | A backdrop press closes, even while the panel enters. [T TestModalBackdropKeepsItsHitAreaWhileEntering] |
| Focus | The focus moves into the panel, stays there and returns to the opener on close. [T TestModalTrapsTheFocusAndRestoresIt] Reopened during its exit, the panel takes the focus again. [T TestModalReopenDuringExitSettlesOnce] |
| Dismissal | `PreventDismiss` blocks Escape and the backdrop; the close button still closes and runs `OnDismiss`. [T TestModalPreventDismiss] |
| Motion | Fades; the exit copy is inert. [T TestSurfacesPaintArrivalAndDeparture] |

### Dialog

| | |
| --- | --- |
| Role and name | A dialog named by `Title`; the close button shows by default. [T TestDialogAndSheetShowCloseByDefault] |
| Keys | Escape closes and runs `OnDismiss`. [T TestResponsiveDialogMotion] |
| Pointer | A backdrop press closes and does not reach the control beneath. [T TestResponsiveDialogBackdropDuringEntry] |
| Focus | Taken into the panel and returned to the opener on close. [T TestDialogAndSheetShowCloseByDefault] |
| Dismissal | `PreventDismiss` as above. [N TestPreventDismissHoldsTheSurfaceAndItsFocus] |
| Motion | Centred on a window of at least 640 DIPs, scaling up; below that it sits at the bottom and rises; the exit copy takes no action. [T TestResponsiveDialogMotion] `Placement.Side` slides it from an edge. [T TestSurfacesPaintArrivalAndDeparture] |

### Sheet

| | |
| --- | --- |
| Role and name | A dialog named by `Title`, on the end side by default; the close button shows by default. [T TestDialogAndSheetShowCloseByDefault] `Flush` gives the panel to `Body` with no header or close button; `Title` still names it. [T TestFlushSheetHandsThePanelToItsBody] It suits a left or right sheet without `Inset`. [C] |
| Keys | Escape closes. [T TestSheetAndDrawerReducedMotionAndDismissal] |
| Focus | As Dialog. [T TestDialogAndSheetShowCloseByDefault] |
| Dismissal | `PreventDismiss` blocks Escape and the backdrop; the close button still closes. [T TestSheetAndDrawerReducedMotionAndDismissal] |
| Motion | Slides from its edge, the one its side names in the layout direction [T TestSideSheetSlidesFromItsEdgeInRightToLeft]; the exit copy is inert; a dismissal during entry leaves no scrim behind. [T TestEdgeSurfacesSlideAndRetainInertExits] [T TestSheetDoesNotStrandScrimWhenDismissedDuringEntry] |

### AlertDialog

| | |
| --- | --- |
| Role and name | An alert dialog named by `Title`. [C] The close button is hidden unless `ShowClose`. [T TestAlertDialogOptionalIconAndClose] |
| Keys | Enter on the focused Cancel cancels. [T TestResponsiveDialogMotion] Escape cancels and runs `OnCancel`. [N TestAlertDialogEscapeCancels] |
| Pointer | Backdrop presses do nothing. [T TestAlertDialogKeepsItsChoiceAndIgnoresDismissal] |
| Focus | Cancel takes the focus on open and keeps it through ignored backdrop presses; the focus returns to the opener after the choice. [T TestAlertDialogKeepsItsChoiceAndIgnoresDismissal] [N TestAlertDialogEscapeCancels] |
| Selection and value | Cancel runs once. The action shows only with `ActionLabel` and `OnAction`, and `ActionDisabled` blocks it. [T TestAlertDialogKeepsItsChoiceAndIgnoresDismissal] [T TestAlertDialogActionRequiresCallbackAndHonorsDisabled] |

### AboutDialog

| | |
| --- | --- |
| Role and name | A dialog named by `Name`; links are named by their labels. [T TestAboutDialogKeepsIdentityInPlaceWhileSectionsScroll] |
| Keys | Escape closes. [T TestAboutDialogKeepsIdentityInPlaceWhileSectionsScroll] A collapsible section is a disclosure, closed until opened. [T TestAboutDialogCopiesInfoAndFoldsLongSections] |
| Focus | Returned to the opener on close. [T TestAboutDialogKeepsIdentityInPlaceWhileSectionsScroll] |
| Selection and value | `CopyInfo` adds a button that copies the information. [T TestAboutDialogCopiesInfoAndFoldsLongSections] |
| States | Icon, name and version stay in place while the sections scroll. [T TestAboutDialogKeepsIdentityInPlaceWhileSectionsScroll] |

### ShortcutsDialog

| | |
| --- | --- |
| Role and name | A dialog named by `Title`; rows are list items inside a list. [C] |
| Keys | Typing narrows rows and groups by every word, ignoring case, and matches keys as well as labels. [T TestShortcutsDialogSearchNarrowsRowsAndGroups] [T TestShortcutsDialogSearchFindsKeysAndBuildsChordsFromHotkeys] The first Escape empties the search; the next closes. [N TestShortcutsDialogEscapeClearsTheSearchBeforeClosing] |
| Focus | The search field takes the focus on open; the focus returns to the opener on close. [T TestShortcutsDialogSearchNarrowsRowsAndGroups] [N TestShortcutsDialogEscapeClearsTheSearchBeforeClosing] |
| States | No match shows "No results"; reopening clears the query; the title and field stay while the rows scroll. [T TestShortcutsDialogSearchNarrowsRowsAndGroups] [T TestShortcutsDialogKeepsSearchInPlaceWhileRowsScroll] |

### ImageViewer

| | |
| --- | --- |
| Role and name | A dialog named by the current image's label; the controls form a toolbar; Previous and Next are disabled at the ends. [C] |
| Keys | `=` and Shift+`=` zoom in, `-` zooms out, `0` resets, Left and Right step, Escape closes. [T TestImageViewerZoomsAndStepsFromTheKeyboard] Right to left, Left steps forward and Previous and Next change sides. [N TestImageViewerStepsTheWayImagesRunInRightToLeft] |
| Pointer | Dragging pans while zoomed; a press outside the image closes, a drag does not. [T TestImageViewerPansWhenZoomedAndClosesOutsideTheImage] |
| Focus | The close button takes the focus on open; Escape and the close button return it to the opener. [N TestImageViewerFocusesCloseAndReturnsTheFocus] |
| Motion | Zoom glides; the next image crossfades. [T TestImageViewerZoomGlidesAndNextImageCrossfades] |

## Edge panels

### Drawer

| | |
| --- | --- |
| Role and name | A dialog named by `Title`; the close button is absent unless `ShowClose`. [T TestDrawerCloseButtonOptions] |
| Keys | Escape closes. [T TestDrawerCloseButtonOptions] No key resizes or swipes it. [C] |
| Pointer | An outward swipe from the grab strip in the middle of the inner edge dismisses; a short one springs back and an inward one does nothing; the body's scroll and controls do not swipe. [T TestDrawerSwipeDirectionThresholdAndFocus] [T TestDrawerBodyScrollAndControlsDoNotSwipe] The threshold is min(96 DIPs, 30% of the panel). [C] |
| Focus | Taken into the panel; the grab keeps it during a swipe; returned to the opener on close or swipe. [T TestDrawerSwipeDirectionThresholdAndFocus] [T TestDrawerCloseButtonOptions] |
| Dismissal | `PreventDismiss` also blocks the swipe, which animates back and can be caught again. [T TestDrawerSwipePreventionAndReducedMotion] [T TestDrawerSwipePreventionAnimatesBackAndCanBeCaught] |
| Motion | Slides from its side; a swipe exit leaves an inert copy. [T TestEdgeSurfacesSlideAndRetainInertExits] [T TestDrawerSwipeExitSnapshotIsInert] Right to left, a start drawer sits on the right edge and swipes away to the right. [N TestDrawerFacesTheWindowInRightToLeft] Its scrim darkens like BottomSheet's; dialogs and sheets share another. [N TestScrimsOfOneKindDarkenAlike] |

### BottomSheet

| | |
| --- | --- |
| Role and name | A dialog named by `Title`, with a close button. [C] |
| Keys | Escape closes. [T TestBottomSheetPaintsItsSlideAndDeparture] No key resizes or swipes it. [C] |
| Pointer | A long downward drag on the handle dismisses; a short one springs back; a drag outside the grab strip does nothing; the grab catches a springing sheet. `HideHandle` removes the handle. [T TestBottomSheetLongDragSlidesOutAndDismisses] [T TestBottomSheetShortDragSpringsBack] [T TestBottomSheetDragOutsideTheGrabStripDoesNotDismiss] [T TestBottomSheetGrabCatchesTheSpringingSheet] |
| Focus | With no body the close button takes the focus; the focus returns to the opener on close. [T TestBottomSheetLongDragSlidesOutAndDismisses] [T TestBottomSheetPaintsItsSlideAndDeparture] |
| Dismissal | The close button slides the sheet out like a swipe before the focus returns; with `HideHandle` it closes at once. [T TestBottomSheetCloseButtonSlidesOut] [T TestBottomSheetHideHandleCloseStillDismisses] `PreventDismiss` also blocks the drag. [T TestBottomSheetLongDragHonoursPreventDismiss] |
| Motion | Slides up and down with its exit painted. [T TestBottomSheetPaintsItsSlideAndDeparture] |

## Anchored popups

### Popover

| | |
| --- | --- |
| Role and name | A popup named by `Label`, else `Title`. [C] |
| Keys | Escape closes. [T TestPopoverAnchorsAndDismisses] |
| Pointer | An outside press closes and reaches the control under the pointer. [N TestPopoverOutsidePressReachesItsTarget] Tall content scrolls. [T TestPopoverKeepsTallContentScrollable] |
| Focus | The panel takes the focus on open; moving the focus out of it closes it. [N TestPopoverTakesTheFocusAndClosesAsTheFocusMovesOn] Opened while a text field anchor holds the focus, it leaves the focus there, so typing goes on; the focus in its anchor does not close it. [N TestPopoverOpenedWhileTypingInItsAnchorLeavesTheFocus] The window ceasing to be active does not close it. [N TestPopoverAndPopconfirmSurviveWindowDeactivation] Escape returns it to the anchor. [T TestPopoverAnchorsAndDismisses] [N TestPopoverOutsidePressReachesItsTarget] |
| Motion | Scales up from 98% beside its anchor, the arrow following a flipped side; reopened during its exit it shows whole. [T TestPopoverEntranceZoomsToFullSize] [T TestPopoverArrowFollowsFlippedPanel] [T TestPopoverReopenDuringExitIsWhole] |

### Popconfirm

| | |
| --- | --- |
| Role and name | A popup with the question and Cancel and Confirm buttons. [C] |
| Keys | Enter activates the focused button; Escape cancels. [T TestPopconfirmFocusesCancelSoEnterDoesNotConfirm] [T TestPopconfirmConfirmsOnceAndDismissalCancels] |
| Pointer | An outside press cancels. [T TestPopconfirmConfirmsOnceAndDismissalCancels] |
| Focus | Cancel takes the focus when `Destructive`, Confirm otherwise; it returns to the anchor. [T TestPopconfirmFocusesCancelSoEnterDoesNotConfirm] [T TestPopconfirmFocusesConfirmWhenNotDestructive] The window ceasing to be active neither closes nor cancels it. [N TestPopoverAndPopconfirmSurviveWindowDeactivation] |
| Selection and value | Confirm runs once; Escape and an outside press run `OnCancel`. [T TestPopconfirmConfirmsOnceAndDismissalCancels] |
| States | `Loading` holds it open and disables Cancel. [T TestPopconfirmStaysOpenWhileLoading] |

### PreviewCard

| | |
| --- | --- |
| Role and name | A popup beside its anchor; it requires a caller-owned state and panics without one. [C] |
| Keys | Tab focus on the anchor opens it; Tab moves into it and it stays open while the focus is inside; Escape dismisses. [T TestPreviewCardHoverFocusAndDismiss] [T TestPreviewCardStaysOpenWhilePanelHasFocus] |
| Pointer | Hover opens it; it stays while the pointer moves into it, and closes after a short grace when the pointer leaves the card, the anchor or the window. An outside press dismisses; an action inside works. [T TestPreviewCardPointerReentersAndActionWorks] [T TestPreviewCardClosesWhenPointerLeavesWindow] |
| States | It closes when its anchor disappears. [T TestPreviewCardClosesWhenAnchorDisappears] |

### Tooltip

| | |
| --- | --- |
| Role and name | `Text` becomes the anchor's description; the tip itself is a tooltip the toolkit keeps from assistive technology. [C] |
| Keys | Keyboard focus shows it; Escape hides it and the focus stays on the anchor. [T TestSurfacesPaintArrivalAndDeparture] |
| Pointer | Pointer rest shows it after the toolkit's delay; a press hides it. [C] |
| Motion | Paints its arrival and departure. [T TestSurfacesPaintArrivalAndDeparture] |

## Guided surfaces

### TeachingTip

| | |
| --- | --- |
| Role and name | A dialog named by `Title` and described by `Body`. [T TestTeachingTipLeavesTheWindowUsableUntilDismissed] |
| Keys | Escape closes only while the focus is inside; its controls follow the anchor in the Tab order. [T TestTeachingTipLeavesTheWindowUsableUntilDismissed] |
| Pointer | Non-modal: an outside press reaches its target and leaves the tip open. [T TestTeachingTipLeavesTheWindowUsableUntilDismissed] |
| Focus | Takes none on open. [T TestTeachingTipLeavesTheWindowUsableUntilDismissed] Closing with the focus inside returns it to the anchor. [N TestNonModalSurfacesReturnTheFocusAsTheyClose] |
| Dismissal | The close button or Escape inside sets `*open` false and runs `OnDismiss`. [T TestTeachingTipLeavesTheWindowUsableUntilDismissed] |
| States | Without a built anchor it floats at the bottom centre without an arrow. [C] |

### Tour

| | |
| --- | --- |
| Role and name | Each step's panel is a dialog named by its title and described by its body; the progress is an image named "Step N of M", announced on each step. [C] |
| Keys | Right or Enter goes forward, Left back, mirrored right to left; Escape ends the tour; Done on the last step finishes it. [T TestTourIsModalAndDrivenByTheKeyboard] |
| Pointer | Presses outside the panel, the target's included, do nothing. [T TestTourIsModalAndDrivenByTheKeyboard] |
| Focus | Next takes the focus on start and keeps it across steps; Back is absent on the first step; Escape returns the focus to where it was. [T TestTourIsModalAndDrivenByTheKeyboard] |
| States | A step without a laid-out target centres its panel without a spotlight. [T TestTourStepWithoutTargetIsCentred] |
| Motion | Spotlight and panel travel together and the text crossfades; under reduced motion they move at once. [T TestTourTravelsContinuouslyToTheNextStep] [T TestTourStepCrossfadesItsText] [T TestTourStepsInPlaceUnderReducedMotion] |

## Non-modal surfaces

### FloatingPanel

| | |
| --- | --- |
| Role and name | A dialog named by `Title`. [T TestFloatingPanelLeavesContentUsableAndClosesFromInside] |
| Keys | Escape closes only while the focus is inside. [T TestFloatingPanelLeavesContentUsableAndClosesFromInside] With the focus inside, Alt and an arrow move it 24 DIPs within the window. [N TestFloatingPanelMovesByKeyWithinTheWindow] |
| Pointer | The content beneath stays usable; dragging the title bar moves the panel within the window. [T TestFloatingPanelLeavesContentUsableAndClosesFromInside] [T TestFloatingPanelDragsByTitleWithinTheWindow] |
| Focus | Takes none on open. [T TestFloatingPanelLeavesContentUsableAndClosesFromInside] Closing with the focus inside moves it to `ReturnFocus` [N TestNonModalSurfacesReturnTheFocusAsTheyClose], else to the opener, if any. [C] |

### FloatingSheet

| | |
| --- | --- |
| Role and name | A dialog named by `Title`; the handle is a focusable splitter named by the locale's Resize. [C] With a finite `MaxHeight` it reports the height within its limits. [N TestFloatingSheetSplitterExposesItsRangeAndReachesItsEnds] |
| Keys | On the handle Up grows and Down shrinks the height within `MinHeight` and `MaxHeight`. [T TestFloatingSheetResizesWithinLimitsAndStaysNonModal] [N TestFloatingSheetGrowsByKeyAndClosesOnEscapeInside] Each press moves 24 DIPs. [C] Home goes to `MinHeight` and End to a finite `MaxHeight`. [N TestFloatingSheetSplitterExposesItsRangeAndReachesItsEnds] Escape closes only while the focus is inside. [N TestFloatingSheetGrowsByKeyAndClosesOnEscapeInside] |
| Pointer | The content beneath stays usable; dragging the handle resizes; the close button closes. [T TestFloatingSheetResizesWithinLimitsAndStaysNonModal] |
| Focus | Takes none on open. Closing with the focus inside moves it to `ReturnFocus`. [N TestNonModalSurfacesReturnTheFocusAsTheyClose] |
| Selection and value | `Height` is caller-owned; nil keeps it in the sheet. [C] |

### ActionBar

| | |
| --- | --- |
| Role and name | A toolbar named by the count and `Label`, such as "2 selected"; count changes are announced. [T TestActionBarFloatsAtTheBottomWithoutTakingTheContainer] [C] |
| Keys | One Tab stop; Left and Right move between its controls. [N TestActionBarIsOneTabStopWithArrowsInside] Escape inside the bar runs `OnClose`; Escape inside its More menu only closes the menu. [T TestActionBarFoldsActionsInANarrowWindowAndEscapeCloses] |
| Pointer | The bar's frame passes presses through to the content it floats over; actions that do not fit fold into More. [T TestActionBarFloatsAtTheBottomWithoutTakingTheContainer] [T TestActionBarFoldsActionsInANarrowWindowAndEscapeCloses] |
| Focus | Takes none on open; closing with the focus inside moves it to `ReturnFocus`. [T TestActionBarFloatsAtTheBottomWithoutTakingTheContainer] |
| Motion | Rises in and out; under reduced motion it appears and leaves at once [N TestActionBarAppearsAndLeavesAtOnceUnderReducedMotion], stays one live bar as the preference turns on while it shows, and takes a theme change at once. [N TestActionBarStaysWholeAndStillAsReducedMotionTurnsOn] |

## Notifications

### Toast

Also covers: ToastState.View, ToastState.Dismiss, ToastState.Items.

| | |
| --- | --- |
| Role and name | Each toast is a status region; Escape on a focused toast and its timers belong to the toolkit. [C] |
| Keys | Tab reaches the close button and its ring shows; Enter on the action runs it. [T TestToastCloseGlyphIsOneCenteredX] [T TestToastDeckFocusExpansionAndSameIDUpdate] F6 from outside the region moves the focus to the newest toast's close button. [N TestToastRegionTakesTheFocusAndKeepsItAsToastsClose] |
| Pointer | Hover reveals the close button; a long horizontal drag dismisses, a short one cancels, and controls inside keep their own presses. [T TestToastCloseRevealsOnHoverAndDismisses] [T TestToastDragCancelsCommitsAndControlsOwnTheGesture] The threshold is max(56 DIPs, 28% of the width). [C] |
| Focus | Focus or hover expands the deck; an update with the same ID keeps the focus and pauses expiry. [T TestToastDeckFocusExpansionAndSameIDUpdate] When the toast holding the focus closes, the newest toast left takes it [N TestToastRegionTakesTheFocusAndKeepsItAsToastsClose]; after the last one closes no element holds it (toolkit limit, see Gaps). [C] |
| Selection and value | `Toast` adds or updates by ID and three show at once; `Dismiss` closes by ID and runs `OnClose`; the action shows only with `OnAction` and closes the toast. [T TestToastUpdatesCapsAndDismisses] [T TestToastActionRequiresCallback] |
| Motion | Under reduced motion it moves at once; a closing toast leaves an inert copy. [T TestToastReducedMotionAndRetainedExit] [T TestToastDismissDuringExitKeepsPaintSafe] |

## Native flyouts

### NewFlyout

Also covers: FlyoutAnchor, FlyoutSurface.

| | |
| --- | --- |
| Role and name | A native window; its focus, outside-press and Escape dismissal belong to the toolkit (see [`docs/flyout.md`](../flyout.md)). Invalid options return an error, and so does WebAssembly. [T TestFlyoutOptions] |
| Pointer | `FlyoutAnchor` gives the anchor's screen rectangle, following the zoom. [T TestFlyoutAnchorZoom] |
| Focus | `FlyoutSurface` fills the window and scrolls; a scrolled footer action is pressable and takes the focus. [T TestFlyoutSurfaceScrollingAndFocus] |
| States | `Header` replaces the title; the surface follows the theme. [T TestFlyoutSurfaceHeaderAndTheme] |

## Gaps

Measured against platform conventions and the
[APG patterns](https://www.w3.org/WAI/ARIA/apg/patterns/). Each remaining gap
is a toolkit limit.

- Toast: when the last toast closes with the focus inside, no element holds
  the focus. Toolkit limit: the toolkit owns the toast's lifetime, and MyGo
  gives no handle to the element focused before the region took the focus;
  `ui.Context` resolves a known handle (`Focused`, `FocusWithin`) but cannot
  name the focused element.
