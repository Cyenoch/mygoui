# Root behavior

Package `mygoui`. Format and evidence marks: [README](README.md). Durations
and easing: [`DESIGN.md`](../../DESIGN.md), [State motion](../../DESIGN.md#state-motion)
and [Loading](../../DESIGN.md#loading). Native scrolling, text selection and
the activation of a focused button by Enter and Space belong to the toolkit.

Shared by every action below: a disabled control takes no clicks or focus and
renders at half opacity. [C] A loading control keeps its size, swaps its
content for a spinner, is busy and takes no activation; loading is not
disabled. [T TestButtonActivationAndPendingState]

## Actions

### Button

| | |
| --- | --- |
| Role and name | A button named by `Label`; an icon-only button shows its label as a tooltip. [C] |
| Keys | Enter and Space activate the focused button. [T TestButtonActivationAndPendingState] |
| Pointer | Activation happens on release over the button; releasing outside does not activate. [T TestButtonPressKeepsGeometryAndActivatesOnRelease] Hover and press change the fill; the link variant underlines instead. [T TestButtonHoverAndPressedPaint] |
| Focus | Click and Tab focus it and draw the focus ring. [T TestButtonActivationAndPendingState] |
| States | Loading keeps the natural width and the position of siblings. [T TestButtonLoadingPreservesNaturalWidthAndSiblingPosition] While pending the face does not react to hover or press and activation is blocked. [T TestButtonPendingKeepsFaceAndBlocksActivation] Disabled and loading both block activation. [T TestButtonActivationAndPendingState] |
| Motion | Press keeps the geometry, with or without reduced motion. [T TestButtonPressKeepsGeometryAndActivatesOnRelease] Swapping the icon crossfades. [T TestStateChangesFadeThroughAnIntermediateFrame] |

### ButtonGroup

| | |
| --- | --- |
| Role and name | A row of separate buttons with joined corners; labelled, the row is a group named by its label. [T TestButtonGroupsAreGroupsAndJoinAToolbar] |
| Keys | Each button is its own Tab stop; inside a `ui.Toolbar` the group joins the toolbar's one Tab stop and the arrows walk it. [T TestButtonGroupsAreGroupsAndJoinAToolbar] |
| States | Separators paint seams without changing layout. [T TestButtonGroupSeparatorsKeepLayoutAndPaintSeams] Right-to-left rounds the start edge. [T TestButtonGroupRightToLeftRoundsStartAndPaintsSeam] A full-width group aligns its labels. [T TestFullWidthButtonGroupAlignsLabels] |

### ToggleButton

| | |
| --- | --- |
| Role and name | A toggle button named by `Label`, exposed as checked while selected. [C] |
| Keys | Space toggles the focused button. [T TestCOSSToggleGroupOutlineKeyboardFocus] Enter also toggles (toolkit). [C] |
| Pointer | A click flips the caller's `Selected` without a callback. [T TestToggleButtonUpdatesCallerStateWithoutCallback] `OnChange` is called when one is set. [C] |
| Selection and value | The caller owns `Selected`; with a nil pointer the button takes no click but stays enabled. [C] |
| States | Loading keeps the selection, is busy and blocks presses. [T TestToggleButtonLoadingKeepsStateAndBlocksPress] Disabled paints the rest face, takes no focus and ignores Space. [T TestDisabledTogglePaintsRestFace] |
| Motion | The selected face fades in. [T TestStateChangesFadeThroughAnIntermediateFrame] |

### ToggleButtonGroup

| | |
| --- | --- |
| Role and name | A row of toggle buttons; joined only when all are outline and not detached. Labelled, the row is a group named by its label. [T TestButtonGroupsAreGroupsAndJoinAToolbar] |
| Keys | Each toggle is its own Tab stop; Tab skips disabled toggles. [T TestCOSSToggleGroupHoverFocusAndSelection] Inside a `ui.Toolbar` the group joins the toolbar's one Tab stop and the arrows walk it. [T TestButtonGroupsAreGroupsAndJoinAToolbar] |
| Focus | The grouped focus ring is inset. [T TestGroupedFocusRingIsInset] |
| Selection and value | Each toggle owns its own flag; selecting one calls its `OnChange` once. [T TestCOSSToggleGroupHoverFocusAndSelection] |

### CloseButton

| | |
| --- | --- |
| Role and name | A button named by `Label`, defaulting to the locale's Close. [T TestDefaultNamesFollowTheLocale] |
| Keys | Enter and Space activate (toolkit). [C] |
| Pointer | Pressing does not activate; releasing does. [T TestCloseButtonShowsItsFaceAtOnce] |
| Motion | Hover and press show their fill at once, without a frame in between. [T TestCloseButtonShowsItsFaceAtOnce] |

### Link

| | |
| --- | --- |
| Role and name | A link named by `Label`. [T TestLinkIsALinkThatTakesEnterNotSpace] |
| Keys | Enter follows it; Space does not. [T TestLinkIsALinkThatTakesEnterNotSpace] |
| Pointer | Hover, press and keyboard focus underline it at once; the trailing icon fades from dim to full on hover over `Motion.State`. [T TestLinkIconFadesOnHover] The library opens no URL; `OnClick` does the navigation. [C] |

### LinkButton

| | |
| --- | --- |
| Role and name | A plain text button named by `Label`. [C] |
| Keys | Enter and Space activate (toolkit). [C] |

### PressableFeedback

| | |
| --- | --- |
| Role and name | A button named by `Label` around caller content. [C] |
| Keys | Enter and Space activate (toolkit). [C] |
| Pointer | Hover and press show a fill at once; press scales the content down slightly. A disabled one shows no feedback. [C] |
| Motion | The press scale settles at once under reduced motion. [C] |

## Containers and layout

### Card

Also covers: Surface, CardHeader, CardTitle, CardDescription, CardContent, CardFooter.

| | |
| --- | --- |
| Role and name | A static card has no role; `Label` names it. With `OnClick` it is one button named by `Label`. [N TestClickableCardIsAKeyboardButton] CardTitle is a level-3 heading; the caller may set another level. [T TestTitlesAreHeadings] Surface and the sections only lay out. [C] |
| Keys | A clickable card is a Tab stop; Enter and Space activate it (toolkit). [N TestClickableCardIsAKeyboardButton] |
| Pointer | Actions inside a static card keep their own clicks. [T TestCardCompositionKeepsActionsIndependent] A clickable card shows the pointer cursor. [C] |
| Focus | A clickable card draws the focus ring and binds the caller's `Handle`. [C] |

### CardFrame

Also covers: CardFrameHeader, CardFrameTitle, CardFrameDescription, CardFrameFooter.

| | |
| --- | --- |
| Role and name | A layout container named by `Label`; the title is a level-3 heading. [T TestTitlesAreHeadings] |

### Frame

Also covers: FramePanel, FrameHeader, FrameTitle, FrameDescription, FrameFooter.

| | |
| --- | --- |
| Role and name | Layout only; no role. FrameTitle is a level-3 heading. [T TestTitlesAreHeadings] |
| Pointer | Actions in separate panels stay independent. [T TestFramePanelsKeepInsetsAndIndependentActions] |

### Group

| | |
| --- | --- |
| Role and name | A group named by `Label` that joins mixed controls. [C] |
| Keys | Each joined control keeps its own Tab stop. [T TestGroupJoinsMixedControlsAndPreservesFocusAcrossReorder] |
| Focus | Keyed items keep their focus when the group is reordered. [T TestGroupJoinsMixedControlsAndPreservesFocusAcrossReorder] Unkeyed items keep distinct identities. [T TestGroupUnkeyedControlsKeepDistinctNativeIdentity] |

### GuideCard

| | |
| --- | --- |
| Role and name | A card named by its title; the cover is decorative. [C] |
| Dismissal | With `OnDismiss` a close button, named by the locale's Close, calls the caller; the caller stops building the card. [T TestGuideCardDismissCallsTheCaller] Focus on the close button is dropped when the card goes; the caller moves it with a `ui.Handle` in `OnDismiss` (toolkit limit, see Gaps). [C] |

### ListGroup

Also covers: ListGroupItem.

| | |
| --- | --- |
| Role and name | The group is a list named by `Label`; each item is a list item named by its title, with its description attached. [T TestListGroupIsAListOfItems] |
| Keys | An item with `OnClick` holds a button and a Tab stop; Enter and Space activate it (toolkit). [T TestListGroupIsAListOfItems] |
| Pointer | A clickable item shows a hover fill and the pointer cursor; a disabled item takes no click. [C] Its chevron points with the reading direction. [T TestListGroupChevronPointsWithTheReadingDirection] |

### SettingsGroup

| | |
| --- | --- |
| Role and name | A group named by its title. A plain row names its bound control by the row title and attaches the subtitle as its description. [T TestSettingsRowTitleNamesItsControl] A row with `OnClick` is one button named by its title. [T TestSettingsPressableRowActivates] An expander row exposes Expanded. [T TestSettingsExpanderRowRevealsNestedRows] |
| Keys | Tab reaches a pressable row and Enter activates it. [T TestSettingsPressableRowActivates] Enter and Space open and close an expander row. [N TestSettingsExpanderRowTogglesFromTheKeyboard] |
| Pointer | Clicking a plain row's title toggles its control. [T TestSettingsRowTitleNamesItsControl] The whole pressable row is the click target; a disabled row is inert. [T TestSettingsPressableRowActivates] |
| Selection and value | The expander follows the caller's `Expanded` flag. [T TestSettingsExpanderRowRevealsNestedRows] It stays open when a row is inserted above it. [T TestSettingsExpanderKeepsItsStateWhenARowIsInsertedAbove] |
| States | Loading with no rows shows skeleton rows under the title; loading with rows keeps them readable but takes no click; the group is busy. [T TestSettingsGroupLoadingFirstLoadAndRefresh] When narrow, the control wraps under the text. [T TestSettingsRowControlWrapsUnderItsTextWhenNarrow] |
| Motion | Nested rows reveal with height, then settle. [T TestBlocksRevealWithHeightThenSettle] The chevron rotates; in right-to-left the pressable row's chevron is mirrored. [C] |

### ResizablePanel

| | |
| --- | --- |
| Role and name | The handle is a splitter named by `Label` or the locale's Resize, with its size as the value. Expanded is exposed when the panel can collapse. [C] |
| Keys | Tab reaches the handle. Left and Right step the size by 16 and stop at the bounds; Home goes to Min; End goes to Max; Enter collapses; Right on a collapsed panel expands it. [T TestResizablePanelKeyboard] PageUp and PageDown step by 64; without a Max, End is not taken: there is no end to go to. [T TestResizablePanelPagesAndEndsWithoutAMaximum] Arrows follow physical direction in right-to-left. [C] |
| Pointer | Dragging clamps to Min and Max. [T TestResizablePanelDragClampsToBounds] Dragging below half of Min collapses. [T TestResizablePanelDragCollapses] Double click toggles collapse. [C] |
| Focus | The handle precedes the panel's content in Tab order on the start or top edge and follows it otherwise; on its own strip it takes the pointer over the content. [T TestResizablePanelHandleFollowsItsEdgeInTabOrder] |
| States | The hairline on the edge takes `Border`, or the theme's. Under the pointer, or held, the handle shows a 2 DIP ring line at once; keyboard focus fades to it over `Motion.State`. [C] |
| Motion | Collapse and expand animate the width; a drag follows the pointer at once. `Shown` reports the size of each build. [C] |

### Sidebar

| | |
| --- | --- |
| Role and name | A navigation landmark named by `Label`, defaulting to the locale's Sidebar. The toggle is a toggle button that is checked while expanded, named by default by the locale's Toggle sidebar. [T TestDefaultNamesFollowTheLocale] Items are buttons; the selected one is the current page. [T TestSidebarMarksTheCurrentPageOfANavigationLandmark] |
| Keys | Each item is its own Tab stop; disabled items are skipped; Enter and Space select the focused item. [T TestSidebarSelectionAndDisabledItems] Space on the toggle collapses the sidebar. [T TestSidebarCollapseAndSlots] Tab scrolls an offscreen item into view. [T TestSidebarTabRevealsOffscreenItem] |
| Focus | Focus stays on the same item across rebuilds. [T TestSidebarFocusFollowsStableItems] It also stays there across collapse. [T TestSidebarCollapsePreservesFocusedItem] |
| Selection and value | The caller owns `Selected` and `Collapsed`; a nil pointer panics. [C] A disabled item takes no click and changes nothing. [T TestSidebarSelectionAndDisabledItems] Collapsed items show their label as a tooltip. [C] |
| States | The selected item keeps its text weight; only its fill and ink change, so the row keeps its size. [T TestSidebarSelectionKeepsTheRowSize] |
| Motion | Collapse moves through intermediate widths over `Motion.Reveal` with `Motion.Ease`, and settles at once under reduced motion. [T TestSidebarCollapseMovesThroughWidths] The selected fill fades. [T TestStateChangesFadeThroughAnIntermediateFrame] |

### ScrollArea

| | |
| --- | --- |
| Role and name | A native scroll viewport named by `Label`, with a focus ring. [C] |
| Keys | Wheel and keyboard scroll natively (toolkit). [T TestScrollAreaKeepsNativeAxesAndWheel] |
| Focus | Children keep their activation and focus under the fade. [T TestScrollAreaFadePreservesChildActivationAndFocus] |
| States | A nil state or an unsupported axis panics. [C] |

### ScrollFade

| | |
| --- | --- |
| Role and name | A decorative layer that passes input through. [C] |
| States | Only the requested edges fade. [T TestScrollFadeDecoratesOnlyRequestedEdges] Rounded viewport corners and borders are kept. [T TestScrollFadeKeepsRoundedViewportCornersAndBorder] A nil state, an unknown mode or a translucent surface fade panics. [T TestScrollFadeRejectsInvalidBoundaryArguments] A mask viewport too large for an effect layer shows the shadow. [T TestScrollMaskFadeShadesAboveTheLayerLimit] |

### ScrollShadow

| | |
| --- | --- |
| Role and name | A native scroll viewport named by `Label`. [C] |
| States | Each edge fades in proportion to the distance left to scroll. [T TestScrollShadowFadesTheScrolledEdges] Horizontal scrolling fades both inline edges. [T TestScrollShadowHorizontalFadesBothInlineEdges] Content that fits shows no fade. [T TestScrollShadowWithoutOverflowPaintsNoFade] The fade is `Effects.FadeEdge` deep by default. [C] A nil state panics, as ScrollFade's and ScrollArea's do. [T TestScrollShadowRejectsANilState] |

### Separator

| | |
| --- | --- |
| Role and name | Decorative; not exposed. [C] |

### Spoiler

| | |
| --- | --- |
| Role and name | The control is a button named by the locale's Show more or Show less. It exposes Expanded. [T TestSpoilerClampsAndTogglesTheCallerFlag] It names the region it controls. [C] It appears only when the content overflows. [T TestSpoilerHidesItsControlWhenContentFits] |
| Keys | Enter on the control toggles the caller's flag. [T TestSpoilerClampsAndTogglesTheCallerFlag] |
| Focus | Focus reaching clamped content opens it. [T TestSpoilerOpensWhenFocusReachesHiddenContent] |
| Motion | The height reveals, then settles. [T TestBlocksRevealWithHeightThenSettle] It runs over `Motion.Reveal` with `Motion.Ease`; the cut fades the last `Effects.FadeEdge` of the content. [C] |

## Status and feedback

### Alert

Also covers: StatusColors, StatusIcon.

| | |
| --- | --- |
| Role and name | A status region named by its title, with the description attached; a danger alert is an alert, announced assertively. [T TestDangerMessagesAreAnnouncedAssertively] StatusColors and StatusIcon are lookups. [C] |
| Keys | The action slot keeps its own button; Enter activates it. [T TestAlertKeepsActionWorkingWithAndWithoutIcon] |
| States | The status icon is opt-in; an explicit icon wins over it. [T TestAlertStatusIconOptInAndExplicitIconPrecedence] |
| Motion | It fades in when it appears. [T TestStateChangesFadeThroughAnIntermediateFrame] It has no exit motion. [C] |

### Banner

| | |
| --- | --- |
| Role and name | A status region across its pane, named by its text; a danger banner is an alert. [T TestBannerSpansItsPaneAndRunsItsControls] [T TestDangerMessagesAreAnnouncedAssertively] |
| Pointer | The action button works. [T TestBannerSpansItsPaneAndRunsItsControls] |
| Dismissal | A close button, named by the locale's Close, appears only with `OnDismiss` and calls it; the caller stops building the banner. [T TestBannerSpansItsPaneAndRunsItsControls] Focus on the close button is dropped when the banner goes, so the next Tab starts from the top of the window; the caller moves it with a `ui.Handle` in `OnDismiss` (toolkit limit, see Gaps). [C] |
| Motion | It fades in. [T TestStateChangesFadeThroughAnIntermediateFrame] It fades out on exit while siblings move at once. [C] |

### EmptyState

Also covers: Header.

| | |
| --- | --- |
| Role and name | No role; the title is a level-3 heading. [T TestTitlesAreHeadings] With no copy it shows the locale's empty placeholder. [T TestEmptyStateDefaultsToTheLocalePlaceholder] Header is muted section text, not a heading. [C] |
| Pointer | The action slot keeps its own button. [T TestEmptyStateCentersRecoveryButton] |

### Result

| | |
| --- | --- |
| Role and name | An outcome icon over an empty state; the icon tile is an image named by the status (Success, Warning, Error, Info). [T TestStatusIsNamedNotOnlyColoured] |
| Pointer | Its actions activate. [T TestResultShowsOutcomeAndActions] |

### StatusDot

| | |
| --- | --- |
| Role and name | The dot is an image named by its status; a default or accent dot is decorative. The label is visible text. [T TestStatusIsNamedNotOnlyColoured] |
| Motion | Pulse and breathe run only when asked and stop under reduced motion. [T TestStatusDotPulsesOnlyWhenAskedAndAllowed] Colour changes fade. [T TestStateChangesFadeThroughAnIntermediateFrame] |

### Chip

| | |
| --- | --- |
| Role and name | Named by its label. [T TestChipAndBadgeKeepAccessibleText] |
| Motion | Colour changes fade. [T TestStateChangesFadeThroughAnIntermediateFrame] |

### Badge

Also covers: BadgeAnchor.

| | |
| --- | --- |
| Role and name | Named by its label. [T TestChipAndBadgeKeepAccessibleText] |
| Pointer | An anchored badge passes clicks through to the content under it. [T TestBadgeAnchorProtrudesAndKeepsContentClickable] |
| Motion | Colour changes fade. [T TestStateChangesFadeThroughAnIntermediateFrame] |

### TagGroup

Also covers: Tag.

| | |
| --- | --- |
| Role and name | The tag row is a list named by `Label`; each tag is a list item, selected when chosen. The error text is a status. [C] The remove button is named by the locale's Remove. [T TestTagGroupSelectsAndRemoves] |
| Keys | The group is one Tab stop; Left and Right walk the tags and their remove buttons, Home and End go to the ends, and Tab leaves the group. Space selects; Backspace removes. [N TestTagGroupIsOneTabStopWalkedByArrows] Delete removes. [T TestTagRemovesWithDeleteKey] Enter also toggles the selection. [C] |
| Pointer | Clicking a tag selects it; clicking its remove button calls `OnRemove`. [T TestTagGroupSelectsAndRemoves] |
| Selection and value | The caller owns the state; multiple selection adds, single selection replaces. [T TestTagGroupSingleSelectionReplaces] `OnChange` receives the selected keys. [T TestTagGroupSelectsAndRemoves] |
| Motion | Selection fades from the face on screen; a removed tag leaves a copy that fades out. [T TestTagSelectionFadesAndRemovalExits] Under reduced motion the selected face shows at once. [T TestTagSelectionHonorsReducedMotion] The exit also shrinks slightly, except under reduced motion. [C] |

## Loading and progress

### Spinner

| | |
| --- | --- |
| Role and name | An indeterminate progress indicator named by `Label` or the locale's Loading. [C] |
| Motion | Still under reduced motion (toolkit). [C] |

### ProgressBar

| | |
| --- | --- |
| Role and name | A progress indicator with its range and value; indeterminate reports no value. [C] |
| Selection and value | The value clamps at the bounds. [T TestProgressClampsAtItsBounds] |
| States | The indeterminate sweep keeps painting. [T TestProgressIndeterminateKeepsPainting] It stays inside the track. [T TestProgressIndeterminateStaysInsideTrack] |
| Motion | A value change eases the fill. [T TestProgressValueChangeEasesTheFill] Reduced motion jumps the fill. [T TestReducedMotionJumpsTheFill] Colour changes fade. [T TestStateChangesFadeThroughAnIntermediateFrame] |

### Meter

| | |
| --- | --- |
| Role and name | A meter with its range and value. [C] Its track geometry differs from a progress bar. [T TestProgressAndMeterUseDistinctTrackGeometry] |

### ProgressCircle

| | |
| --- | --- |
| Role and name | A progress indicator named by `Label`, defaulting to the percentage, or Loading when indeterminate. [C] |
| Motion | A value change eases and covers its bounds. [T TestProgressCircleEasesAndCoversItsBounds] |

### Skeleton

Also covers: SkeletonGroup.

| | |
| --- | --- |
| Role and name | A skeleton is decorative; the group is a busy group named by the locale's Loading. [T TestSkeletonGroupIsBusy] |
| Motion | The pulse lowers the fill to half its alpha and back, eased in and out, one shape for a skeleton and a group. [C] The sheen is opt-in. [T TestSkeletonMotionIsOptIn] It stops under reduced motion and high contrast. [T TestSkeletonStillsForHostPreferences] It repaints without rebuilding still content. [T TestLoadingPaintFramesDoNotRebuildStillContent] |

### LoadingOverlay

| | |
| --- | --- |
| Role and name | While visible the content is busy and the cover is a status named by `Label` or the locale's Loading. [T TestLoadingOverlayBlocksAndAnnouncesWhileVisible] |
| Focus | Showing it takes the focus out of the content. [T TestLoadingOverlayBlocksAndAnnouncesWhileVisible] Hiding it gives the focus back to where it was. [T TestLoadingOverlayGivesTheFocusBackWhenItHides] |
| States | The content takes no input from the first frame, while the cover shows only after the loading delay. [T TestLoadingOverlayLooksLateButHoldsInputAtOnce] |

## Data display

### Avatar

Also covers: AvatarGroup.

| | |
| --- | --- |
| Role and name | An image named by `Name`. [C] Initials stand in until the image arrives. [T TestAvatarImageSourcesAndFallback] The group is a group; its "+N" count is an image named by the locale's count of the rest ("3 more"). [T TestAvatarGroupNamesItsOverflow] |
| States | The count follows the group size. [T TestAvatarGroupCountFollowsTheGroupSize] |
| Motion | A late image fades in. [T TestStateChangesFadeThroughAnIntermediateFrame] |

### Descriptions

| | |
| --- | --- |
| Role and name | A group named by its title; each item is a group named by its label. Values are selectable; an empty value shows a dash; a truncated value shows the full text as a tooltip. [C] |
| States | Loading keeps the labels, hides the values behind skeletons and is busy. [T TestDescriptionsLoadingKeepsLabelsAndHidesValues] It stays readable when narrow. [T TestDescriptionsStayReadableWhenNarrow] |

### Highlight

| | |
| --- | --- |
| Role and name | One paragraph read as one text [T TestHighlightKeepsTheTextWhole]; matches are marked by colour only (toolkit limit, see Gaps). [C] |
| States | Matches cover whole multi-byte characters. [T TestHighlightMarksWholeMultiByteCharacters] |

### Image

| | |
| --- | --- |
| Role and name | An image named by `Label`; with `OnClick` and a ready image it is a button named by `Label` instead. [C] |
| Keys | A clickable image is a Tab stop; Enter and Space activate it (toolkit). [C] |
| Pointer | Only a ready image takes the click. [T TestImageShowsTheCallerStatusInItsBox] |
| States | Loading shows a skeleton and failure shows a broken-image icon or the caller's fallback, in the image's box. [T TestImageShowsTheCallerStatusInItsBox] |
| Motion | The image fades in when it arrives. [T TestStateChangesFadeThroughAnIntermediateFrame] |

### Kbd

Also covers: KbdGroup.

| | |
| --- | --- |
| Role and name | Named by its keys joined with spaces. [C] |

### NumberFlip

| | |
| --- | --- |
| Role and name | One text named by `Label` or the value, read once rather than per glyph. [T TestNumberFlipRestsAndReadsOnce] |
| Motion | Digits roll to the new value; other glyphs crossfade. [T TestNumberFlipRollsToTheNewValue] Reduced motion keeps glyphs in place. [T TestNumberFlipReducedMotionKeepsGlyphsInPlace] |

### Statistic

| | |
| --- | --- |
| Role and name | A group named by `Label`. The trend icon is named Increase or Decrease. [T TestStatisticShowsValueAffixesAndNamedTrend] |
| States | Loading holds the value's height and is busy. [T TestStatisticLoadingHoldsTheValueHeight] |
| Motion | The value rolls as a NumberFlip. [C] The delta colour fades. [T TestStateChangesFadeThroughAnIntermediateFrame] |

### Timeline

| | |
| --- | --- |
| Role and name | A list named by `Label`; each event is a list item named by its title with its description attached. [T TestTimelineListsEventsInOrder] Each marker is an image named by its status; a default or accent marker is decorative. [T TestStatusIsNamedNotOnlyColoured] |
| Keys | The control folding earlier events is a disclosure; Tab reaches it and Enter opens it. [N TestTimelineEarlierControlOpensByKeyboard] An event with `OnClick` has a button; Enter activates it (toolkit). [C] |
| Pointer | Clicking the control shows the earlier events. [T TestTimelineFoldsEarlierItemsBehindAControl] |
| States | Only the last item is busy while pending, with a loading marker. [T TestTimelinePendingMarksOnlyTheLastItem] |
| Motion | Marker colour changes fade. [T TestStateChangesFadeThroughAnIntermediateFrame] |

## Typography

### Typography

Also covers: Heading, Paragraph, Code, Prose.

| | |
| --- | --- |
| Role and name | Heading exposes a heading with its level, clamped to 1 to 6. [C] Paragraph and Code are selectable text; Prose only lays out. [C] |

## Gaps

Measured against platform conventions and the
[APG patterns](https://www.w3.org/WAI/ARIA/apg/patterns/). Each remaining
gap is a toolkit limit.

- Banner, GuideCard: dismissing from the close button leaves no focus; the
  next Tab starts at the top of the window instead of after the removed
  region (probed for Banner). Toolkit limit: MyGo drops the focus of an
  element that is no longer built and moves Tab from the window's start
  (`ui/input.go`, `moveFocus`); it offers no call to focus the element after
  a removed one. The caller, who knows what follows, focuses it through a
  `ui.Handle` in `OnDismiss`.
- Highlight: matches are marked by colour only. Toolkit limit: a `ui.Span`
  carries no role, and MyGo's roles (`ui/access.go`) have no mark; splitting
  the paragraph into elements would break it into separate texts.
