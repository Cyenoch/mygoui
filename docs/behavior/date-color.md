# Date and color behavior

Packages `date` and `color`. Format and evidence marks: [README](README.md).
Durations and easing: [`DESIGN.md`](../../DESIGN.md), "State motion" and
"Overlay motion". "Cmd" is Command on macOS and Control elsewhere.

The text editors in DateField, TimeField, DateRangePicker and ColorField are
native MyGo text inputs. Caret movement, selection, clipboard and IME
composition belong to the toolkit. So do the popover's placement, outside
click and Escape, and the radio group under ColorSwatchPicker. Every
`ui.Animation` here jumps to its end under reduced motion unless an entry
says otherwise.

## Date

### Calendar

Also covers: CalendarState, CalendarViewState (caller-owned state).

| | |
| --- | --- |
| Role and name | The day grid is a table named by Label, or by the month text when Label is empty. Its weeks are rows. The day grid is a grid and each day a grid cell named by DayLabel, which defaults to "Monday, January 2, 2006" in English; a chosen day is Selected. [T TestCalendarDaysAreGridCellsSelectedWhenChosen] The cursor day is the grid's active descendant. [N TestCalendarDayGridKeysMoveTheCursor] Out-of-bounds and DisabledDates days are disabled. [C] The month buttons are named by PrevLabel and NextLabel, which default to the locale's Previous and Next. [T TestCalendarHeaderPlacesArrowsAroundCaption] With ShowYearPicker the heading is a button named "<month>, Choose year". [T TestCalendarYearPickerDefaultPage] The year grid is a group named by the locale's ChooseYear and each year is a button that is Selected when it is the current year. [T TestDateControlsTakeTheirNamesFromTheLocale] |
| Keys | Right moves to the next day and crosses into the next month. [T TestCalendarKeyboardCrossesMonthBoundary] PageDown moves a month and clamps the day to that month. [T TestCalendarKeyboardCrossesMonthBoundary] Up and Down move a week. [N TestCalendarDayGridKeysMoveTheCursor] Left moves a day. [N TestCalendarDayGridKeysMoveTheCursor] Home goes to the first day of the week and End to the last. [T TestCalendarHomeEndKeepTheWeekAndShiftPagesAYear] PageUp goes back a month. [N TestCalendarDayGridKeysMoveTheCursor] Shift+PageUp and Shift+PageDown move a year. [T TestCalendarHomeEndKeepTheWeekAndShiftPagesAYear] Space chooses the cursor day. [N TestCalendarDayGridKeysMoveTheCursor] Enter chooses the cursor day. [T TestCalendarFocusAndSelectionAcrossDifferentWeekCounts] Enter on a disabled day does nothing. [T TestCalendarGridClickFocusesKeyboard] PageUp and PageDown stop at Min and Max. [T TestCalendarBoundsStopMonthNavigation] The day arrows stop at Min and Max. [T TestCalendarArrowsStopAtTheBounds] Right to left, Left and Right swap. [T TestCalendarArrowsMirrorRightToLeft] In the year grid, Left and Right move a year and Up and Down move three. [N TestCalendarYearGridArrowsMoveTheYear] End goes to the last offered year, and Enter chooses it and returns focus to the heading. [T TestYearPickerScrollsWithinMonthHeight] Home goes to the first offered year and Space chooses. [C] Escape closes the year grid. [T TestCalendarYearPickerEscapeCloses] The year arrows move the visible month as they go, and Escape restores the year the grid opened on. [T TestCalendarYearGridEscapeRestoresTheYear] |
| Pointer | A click on an enabled day chooses it and focuses the grid. [T TestCalendarSelectsDay] [T TestCalendarGridClickFocusesKeyboard] A disabled day ignores clicks. [T TestCalendarRejectsDisabledDate] A shown day from the next or previous month can be chosen and moves the view to its month. [T TestCalendarSelectsVisibleOutsideDay] Previous and Next change the month. [T TestCalendarMonthNavigation] They stop at Min and Max. [T TestCalendarBoundsStopMonthNavigation] While the year grid is open they ignore clicks. [T TestCalendarYearPickerHidesMonthNavigation] Hover paints plain and today cells and reverts. [T TestCalendarDayHoverPaint] A chosen cell keeps its fill. [T TestCalendarDayHoverPaint] |
| Focus | The Tab order is Previous, Next, then the day grid. [N TestCalendarDayGridKeysMoveTheCursor] When ShowYearPicker is set the heading button sits between Previous and Next. [C] With nothing chosen, entering the grid puts the cursor on today. [N TestCalendarDayGridKeysMoveTheCursor] The grid is one Tab stop. [N TestCalendarDayGridKeysMoveTheCursor] With a chosen day, the cursor starts there. [C] The focus ring is painted on the cursor day while focus is visible. [C] Changing Reduce Motion keeps the focus. [T TestCalendarGridClickFocusesKeyboard] Opening the year grid moves the focus to the current year. [N TestCalendarYearGridArrowsMoveTheYear] Choosing a year returns the focus to the heading. [T TestYearPickerScrollsWithinMonthHeight] Escape closes the year grid and returns the focus into the calendar, to the heading or the first day. [T TestYearPickerOptInAcrossDateControls] It lands on the heading. [C] |
| Selection and value | The caller owns the CalendarState: Value, Month, Cursor and YearOpen. A nil state uses one inside the component. [C] OnChange reports a chosen day. Moving the cursor or paging never reports. [T TestCalendarFocusAndSelectionAcrossDifferentWeekCounts] Unbounded, the year grid offers one page around the visible year. With bounds it offers exactly the years between them. [T TestCalendarYearPickerDefaultPage] [T TestCalendarYearPickerBounds] Choosing a year keeps the month. [T TestCalendarYearPickerDefaultPage] It clamps the day. [C] |
| States | Disabled ignores clicks, and today keeps its dot. [T TestDisabledCalendarMarksTodayAndCannotActivate] It also dims the calendar, removes it from the Tab order and exposes it as disabled. [C] Disabled days are struck through. [C] A chosen day that is disabled is not painted primary. [T TestCalendarDisabledSelectionDoesNotPaintPrimary] Days outside the month are shown at reduced opacity. [C] |
| Motion | A chosen day's fill fades in. [T TestCalendarDayFillFadesIn] A calendar built with a chosen day shows it at once. [T TestCalendarDayFillFadesIn] The year grid, the day grid and the heading swap over the state motion. [T TestCalendarViewsSwapOverTheStateDuration] The year grid fades in. [T TestCalendarYearPickerFadesIn] It paints at once under reduced motion. [T TestCalendarYearPickerFadesIn] The heading chevron turns and the month buttons fade with it. [C] Hover is instant. [C] |

### RangeCalendar

Also covers: DateRange (Empty, Ordered, Contains), RangeCalendarState.Clear.

| | |
| --- | --- |
| Role and name | The same as Calendar. Every day from the start of the range to its end is Checked. [C] |
| Keys | The same keys as Calendar. Enter on the cursor day sets or completes an endpoint. [T TestCalendarGridClickFocusesKeyboard] Enter on a disabled day leaves the pending endpoint as it is. [T TestRangeCalendarSelectionAcrossDifferentWeekCounts] Escape cancels a pending endpoint and restores the range it replaced. [T TestRangeCalendarEscapeCancelsThePendingEndpoint] |
| Pointer | The first activation sets the start, records the anchor and does not report. [T TestRangeCalendarOrdersEndpoints] The second completes the range in order and reports it once. [T TestRangeCalendarOrdersEndpoints] A range can end in another month. [T TestRangeCalendarAcrossMonths] An out-of-bounds day cannot start a range. [T TestRangeCalendarRejectsOutOfBoundsEndpoint] A shown day from another month starts a new range. [T TestRangeCalendarShowsOutsideDaysInRange] While an endpoint is pending, the days between it and the pointer are painted as the range at once. [T TestRangeCalendarPreviewsTheRangeUnderThePointer] |
| Focus | The same as Calendar. [T TestCalendarGridClickFocusesKeyboard] |
| Selection and value | The caller owns RangeCalendarState: Value, Anchor and the view. OnChange reports only a completed range, and completing one clears the anchor. [T TestRangeCalendarOrdersEndpoints] A same-day range is allowed and has rounded ends. [T TestRangeCalendarSameDayHasRoundedEnds] |
| States | While only the start is chosen, no band is painted. [T TestRangeCalendarStartOnlyPaintsNoBand] A disabled day breaks the band. [T TestRangeCalendarDisabledDayStopsBand] The band is rounded at week edges. [T TestRangeCalendarBandRoundsAtWeekEdges] A disabled today keeps its marker. [T TestRangeCalendarDisabledTodayKeepsMarker] |
| Motion | The same as Calendar. [C] |

### DateField

| | |
| --- | --- |
| Role and name | A form field whose native text input is named by Label. [C] The placeholder defaults to "YYYY-MM-DD", which is not localized. [C] |
| Keys | The toolkit's text editing. [C] |
| Focus | One Tab stop. The input group shows the shared focus-within ring. [C] |
| Selection and value | Value is caller-owned; a nil Value keeps the date inside the field. [C] A complete in-bounds ISO date commits as it is typed. An out-of-bounds date shows the localized invalid message and does not commit. [T TestDateFieldCommitsISOAndRejectsOutOfBounds] A date that does not exist, such as February 30, is invalid and is not rolled over. [T TestDateFieldRejectsImpossibleDate] A date in DisabledDates is invalid. Partial text shows no error. Clearing the text commits the zero time and reports it. A Value changed by the caller replaces the text. [C] |
| States | Hovering does not change the surface. [T TestDateFieldHoverPaint] Disabled and Invalid come from Field. [C] |

### TimeField

| | |
| --- | --- |
| Role and name | A form field whose native text input is named by Label. [C] The placeholder is "HH:MM", or "hh:mm AM" with Hour12. [C] |
| Keys | The toolkit's text editing. [C] |
| Focus | One Tab stop. [C] |
| Selection and value | The caller owns TimeFieldState; a nil state keeps the time inside the field. [C] A complete HH:MM commits as it is typed and is reported. [T TestTimeFieldCommits] With Hour12, 12:00 AM is midnight and 12:30 PM is after noon. [T TestTimeFieldNoonAndMidnight] Changing Hour12 rewrites the text. [T TestTimeFieldUpdatesWhenHourCycleChanges] The date part of the state is kept, or is today's date when the state is empty. Text that does not parse shows the invalid message. Clearing the text reports the zero time. [C] A time before Min, after Max or off Step (counted from Min) is invalid and is not committed. [T TestTimeFieldRejectsTimesOutsideItsBoundsAndStep] |

### DatePicker

| | |
| --- | --- |
| Role and name | A DateField with a trigger button. The trigger is named by OpenLabel, which defaults to the locale's OpenCalendar. [T TestDateControlsTakeTheirNamesFromTheLocale] The popup's day grid is named by Label. [N TestDatePickerOutsideClickRestoresTrigger] The trigger reports Expanded while the popup is open. [T TestDatePickerTriggersReportExpanded] |
| Keys | Tab moves from the text input to the trigger. Enter or Space on the trigger opens the popup. [C] Escape closes the popup and reports OnOpenChange(false). [T TestDatePickerNotifiesEscapeDismiss] Escape in the popup's year grid closes only the year grid. [T TestYearPickerOptInAcrossDateControls] Inside the popup the keys are the same as Calendar's. [C] |
| Pointer | The trigger opens the popup. Choosing a day commits it and closes the popup. [T TestDatePickerOpensAndCommits] Hovering the trigger paints it. [T TestDatePickerTriggerHoverAndFocus] |
| Focus | Opening the popup moves focus into the day grid. [N TestDatePickerOutsideClickRestoresTrigger] Escape returns focus to the trigger, and the popup reopens after its exit transition. [T TestDatePickerTriggerHoverAndFocus] Choosing a day also returns focus to the trigger. [C] |
| Dismissal | Escape [T TestDatePickerNotifiesEscapeDismiss] or a click outside closes the popup without choosing. Either way OnOpenChange(false) is reported and focus returns to the trigger. [N TestDatePickerOutsideClickRestoresTrigger] With Flyout the calendar opens as a native flyout, and dismissing the flyout takes the same path. [C] |
| Selection and value | Value and Open are caller-owned; nil keeps them inside the component. Changing the value from outside moves the calendar to it. [C] Typing in the field follows DateField. [C] |
| States | Disabled disables the trigger as well, so the popup cannot open. [C] |
| Motion | The popup uses the overlay motion. [C] |

### DateRangePicker

| | |
| --- | --- |
| Role and name | Two native text inputs inside one field. They are named by StartLabel and EndLabel, which default to the locale's StartDate and EndDate. [T TestDateControlsTakeTheirNamesFromTheLocale] An en dash sits between them. [C] A trigger button is named by OpenLabel, which defaults to the locale's OpenCalendar, and reports Expanded while the popup is open. [T TestDatePickerTriggersReportExpanded] |
| Keys | The same as DatePicker. [C] |
| Pointer | In the popup, the first day does not commit or close. The second day commits the ordered range, reports it once and closes the popup. [T TestDateRangePickerCompletesAndRestoresFocus] |
| Focus | Completing a range returns focus to the trigger. [T TestDateRangePickerCompletesAndRestoresFocus] Escape and a click outside also return focus to the trigger. [C] |
| Dismissal | The same as DatePicker. [C] A pending endpoint survives the dismissal and is still pending when the popup reopens. [C] |
| Selection and value | Typed text commits once both endpoints are complete and valid. Endpoints typed in reverse are ordered and their text is rewritten. An empty or partial endpoint waits. Clearing both endpoints reports an empty range. [C] The calendar follows Value while no endpoint is pending. [C] |
| States | Disabled disables both inputs and the trigger. [C] |
| Motion | The same as DatePicker. [C] |

## Color

The color controls share one caller-owned ColorState. ColorArea, ColorSlider
and ColorField run OnChange only when the color changes. [C]

### ColorArea

| | |
| --- | --- |
| Role and name | Named by Label [T TestColorAreaDragChangesColor], which defaults to "Color area". It is a slider. [C] Its value text is "<saturation>%, <brightness>%". Its numeric range is saturation alone, from 0 to 100: a MyGo element exposes one range. [C] |
| Keys | The arrows step saturation (Left and Right) and brightness (Down and Up) by one percent, or by ten with Shift. [N TestColorAreaArrowsStepSaturationAndBrightness] Home sets saturation to zero and End sets it to full. [N TestColorAreaArrowsStepSaturationAndBrightness] Each change is reported once. [N TestColorAreaArrowsStepSaturationAndBrightness] PageUp and PageDown step brightness by ten. Right to left, saturation grows to the left: the gradient, the thumb, the pointer and Left and Right all mirror. [T TestColorAreaPagesBrightnessAndMirrorsRightToLeft] |
| Pointer | Pressing and dragging with the primary button set saturation and brightness from the pointer, clamped to the area. [T TestColorAreaDragChangesColor] |
| Focus | One Tab stop. [C] The focus ring is drawn around the thumb while focus is visible. [C] |
| Selection and value | The hue survives a move to black. [T TestColorStateKeepsHueThroughBlack] |
| States | Disabled dims the area, removes it from the Tab order, exposes it as disabled and ignores input. [C] There is no hover face. [T TestColorAreaThumbPressed] |
| Motion | A press grows the thumb by Effects.ThumbPressScale over the state motion, and it settles back on release. [T TestColorAreaThumbPressed] |

### ColorSlider

Also covers: ColorChannel (String, Range, Percent).

| | |
| --- | --- |
| Role and name | A slider named by Label, which defaults to the channel's name. Its range is the channel's range. Its value text is the formatted value, such as "120°" or "50%". [C] ShowOutput prints that value beside the label. [C] |
| Keys | On a horizontal slider, Right and Left step by one unit, or by ten with Shift. PageUp and PageDown step by ten, and Home and End go to the ends. [N TestColorSliderKeysStepTheChannel] Up and Down step every slider, Up raising the value. [T TestColorSliderUpAndDownStepAHorizontalSlider] RTL swaps Left and Right on a horizontal slider. A vertical slider takes Left and Right as well. [C] |
| Pointer | Pressing and dragging set the channel. A press beyond the end caps reaches the bounds. [T TestColorSliderChannelUpdates] [T TestColorSliderPointerClampsToEnds] While pressed, the cursor changes from grab to grabbing. [C] |
| Focus | One Tab stop. The focus ring is drawn around the thumb. [C] |
| Selection and value | A hue slider changes only the hue. [T TestColorSliderChannelUpdates] Every channel clamps to its range. [T TestColorStateClampsChannels] |
| States | Disabled dims the slider, removes it from the Tab order, exposes it as disabled and ignores input. [C] The end caps show the channel's extreme colors. [T TestColorSliderCapsPaintChannelEnds] |
| Motion | None. The thumb follows the value at once. [C] |

### ColorField

Also covers: ParseHex, FormatHex, ColorState.SetHex and ColorState.Hex.

| | |
| --- | --- |
| Role and name | A form field whose native text input is named by Label, which defaults to "Hex". A small swatch of the current color comes before the input. [T TestColorFieldParsesAndRejects] [C] |
| Keys | The toolkit's text editing. [C] |
| Focus | One Tab stop. [C] |
| Selection and value | The text is parsed as it is typed. #rgb, #rgba, #rrggbb and #rrggbbaa are accepted, with or without the #. Text that does not parse keeps the last valid color and shows the localized invalid message. [T TestColorFieldParsesAndRejects] [T TestParseAndFormatHex] The error is announced as a status. A color changed by another control replaces the text. [C] |
| States | Disabled and Invalid come from Field. [C] |

### ColorSwatch

| | |
| --- | --- |
| Role and name | An image named by Label, which defaults to the hex value. [C] Selected is drawn as a tick and is the image's description, the locale's Selected. [T TestColorSwatchSaysWhenItIsSelected] |
| States | The tick is white on a dark color and black on a light one. [T TestColorSwatchSelectedCheckIsDrawn] A translucent color is drawn over a checkerboard. [C] |
| Motion | The tick draws in when the swatch becomes selected [T TestSwatchTickDrawsInWhenSelected] and fades out when it is unselected. [C] |

### ColorSwatchPicker

| | |
| --- | --- |
| Role and name | A native radio group named by Label, which defaults to "Swatches". Each swatch is a radio button named by its Names entry, or by its hex value, and the same name is its tooltip. [C] |
| Keys | The toolkit's radio group. Tab reaches the chosen swatch, and the arrows move and choose. [T TestColorSwatchPickerArrowRoving] |
| Pointer | A click chooses a swatch. [T TestColorSwatchPickerSelects] Hover grows the swatch inside its border at once, and leaving shrinks it at once. [T TestSwatchHoverGrowsAtOnce] |
| Selection and value | Value is caller-owned. A choice writes Value and reports it. A nil Value keeps the choice inside the picker, and a choice is still reported. [T TestColorSwatchPickerChoosesWithoutAValue] |
| States | Disabled disables the group and stops the hover growth. [C] |
| Motion | The chosen swatch shrinks over the state motion; hover does not animate. [C] The chosen swatch's border fades in. [T TestSwatchPickerChoiceBorderFades] Its tick scales and fades in. [T TestSwatchPickerCheckmarkAnimatesIn] [T TestColorSwatchPickerSelectedCheckIsDrawn] |

### ColorPicker

Also covers: NewColorState and the rest of ColorState's API; DefaultSwatches,
the palette shown when Swatches is nil.

| | |
| --- | --- |
| Role and name | The trigger is a color well named by Label, which defaults to "Color picker". Its value is the hex text, and it reports Expanded while the panel shows. The panel is named by Label. [N TestColorPickerDismissalReturnsToTrigger] [C] |
| Keys | Enter or Space on the trigger toggles the panel. [C] Escape closes it. [N TestColorPickerDismissalReturnsToTrigger] |
| Pointer | A click on the trigger opens the panel and a second click closes it. [N TestColorPickerDismissalReturnsToTrigger] [C] The panel's area, sliders, field and swatches all edit the same state. [T TestColorPickerSyncsValue] |
| Focus | Opening moves focus to the panel's first control. [T TestColorPickerMovesFocusIntoThePanelAndReportsOpening] Escape returns focus to the trigger. [N TestColorPickerDismissalReturnsToTrigger] |
| Dismissal | Escape or a click outside closes the panel without changing the color. [N TestColorPickerDismissalReturnsToTrigger] [C] OnOpenChange reports opening and every way of closing. [T TestColorPickerMovesFocusIntoThePanelAndReportsOpening] With a nil Open, only the trigger is shown. [C] |
| States | Disabled disables the trigger and every control in the panel. [C] |
| Motion | The panel fades in with the overlay motion. [T TestColorPickerPopoverFadesIn] |

## Gaps

Measured against platform conventions and the
[APG patterns](https://www.w3.org/WAI/ARIA/apg/patterns/): the date picker
dialog, grid and slider.

- ColorArea, toolkit limit: assistive technology gets the brightness only as
  value text. The numeric value is saturation alone, because a MyGo element
  exposes one range: `Element.Range(lo, hi, value)` and `AccessInfo`'s single
  Min, Max and Now (ui/api_gen.go).
- With Flyout, the pickers' focus return and reports go through the native
  flyout window; that path is unverified in a real window.
