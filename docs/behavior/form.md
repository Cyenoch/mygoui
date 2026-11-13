# Form behavior

Package `form`. Format and evidence marks: [README](README.md). Durations and
easing: [`DESIGN.md`](../../DESIGN.md), "State motion" and "Loading". Caret,
selection, clipboard and input methods in every editor belong to the toolkit.

Shared by every labelled control below: `Label` names the control, and
`Error` sets Invalid and is read as the error message. [T TestInputInvalidAnnouncesErrorAndStaysEditable]
While invalid the description is hidden. [T TestTextFieldInvalidHidesDescription]
`Required` appends an asterisk to the visible label. [C]
Disabled controls take no clicks or focus and render at half opacity; a
disabled container disables everything inside it. [C]

## Field structure

### Field

Also covers: Label, Description, FieldError, ErrorMessage.

| | |
| --- | --- |
| Role and name | The label names the one control built in the body; the description and error are attached to it. [T TestFieldNativeHelpTypographyAndFocus] FieldError and ErrorMessage are status regions. [C] |
| Pointer | Clicking a `Label` with a target focuses that control unless the label is disabled. [T TestFieldNativeHelpTypographyAndFocus] |

### Fieldset

Also covers: FieldGroup, FieldActions, Form.

| | |
| --- | --- |
| Role and name | Fieldset is a group named by its legend. [C] FieldGroup, FieldActions and Form only lay out; Form shows a top-level error and never submits. [C] |
| States | A disabled fieldset disables every control inside it. [C] |

### ControlField

| | |
| --- | --- |
| Role and name | The whole row is one check box or switch named by its label, with the description attached. [C] |
| Keys | Space toggles. [T TestControlFieldSwitchRow] |
| Pointer | A click anywhere on the row toggles. [T TestControlFieldCheckboxRowAndDisabled] |
| States | Disabled ignores clicks and keeps the value. [T TestControlFieldCheckboxRowAndDisabled] |

## Text entry

### Input

Also covers: TextField.

| | |
| --- | --- |
| Role and name | Text field named by AriaLabel, else Label, else Placeholder. [C] |
| Keys | Typing writes through to the caller's value. [T TestInputCommitsTypedValue] Enter calls OnSubmit when set and otherwise does nothing. [T TestTextFieldTypesAndSubmits] |
| Focus | One Tab stop; the focus ring shows on any focus, pointer included. [C] |
| States | Disabled takes no input. [T TestDisabledTextControlsTakeNoInput] ReadOnly keeps the value and stays focusable. [T TestInputReadOnlyKeepsValue] Invalid stays editable and is exposed. [T TestInputInvalidAnnouncesErrorAndStaysEditable] Loading reports busy, shows a trailing spinner and stays editable. [T TestLoadingInputStaysEditableAndReportsBusy] |
| Motion | The invalid border fades; the focus border is instant. [T TestFieldBorderFadesWhenInvalid] The spinner fades in and its frames stop once loading ends. [T TestLoadingInputStaysEditableAndReportsBusy] |

### TextArea

| | |
| --- | --- |
| Role and name | Multiline text field, named as Input. [C] |
| Keys | Enter inserts a line. [T TestTextAreaTypesMultipleLines] |
| Selection and value | Grows from the minimum to the maximum line count, then scrolls. [C] |
| States | ReadOnly keeps the value. [T TestTextAreaReadOnlyKeepsValue] No Loading. [C] |

### SearchField

| | |
| --- | --- |
| Role and name | Text field marked as search; the name falls back to the localized "Search". [C] The clear button is named "Clear". [T TestSearchFieldClearsWithButtonAndEscape] |
| Keys | Escape with text clears it and calls OnClear. [T TestSearchFieldClearsWithButtonAndEscape] It also calls OnChange and is consumed before an ancestor's Escape shortcut. [T TestSearchFieldEscapeClearsBeforeAncestor] Escape on an empty, disabled or read-only field is passed on to ancestors. [N TestSearchFieldEscapeOnEmptyValueReachesAncestor] |
| Focus | On an empty field the clear button is hidden from the pointer, Tab and assistive technology. [T TestSearchFieldHiddenClearIsNoTabStop] |
| Dismissal | The clear button empties the value. [T TestSearchFieldClearsWithButtonAndEscape] |
| States | Disabled and ReadOnly disable the clear button. [C] Loading swaps the search icon for a spinner and stays editable. [T TestLoadingSearchFieldKeepsSearchEditable] |

### PasswordField

| | |
| --- | --- |
| Role and name | Password text field; the reveal toggle is a toggle button named "Show password" (plus the subject when set), checked while revealed. [T TestPasswordFieldToggleRevealsTypedText] |
| Pointer | The toggle reveals and conceals. [T TestPasswordFieldToggleRevealsTypedText] |
| Focus | Focus leaving the field, including the window losing the keyboard, conceals it again unless KeepRevealed. [T TestPasswordFieldConcealsAgainWhenFocusLeaves] |
| Selection and value | Revealed is caller-owned when given, else kept per window. [T TestPasswordFieldFollowsCallerRevealedFlag] Strength is reported as typed. [T TestPasswordFieldReportsStrengthAsTyped] |
| States | Disabled makes the toggle do nothing. [T TestPasswordFieldDisabledToggleDoesNothing] |
| Motion | The eye icons crossfade; strength segments fade. [T TestPasswordRevealCrossfades] |

### InputGroup

| | |
| --- | --- |
| Pointer | Clicking a non-interactive prefix, suffix or block focuses the editor. [T TestInputGroupAddonFocusesEditor] |
| Focus | Addon actions keep their own focus; the shell's ring follows the editor only. [T TestInputGroupAddonActionKeepsItsFocus] A block prefix leaves the editor as the field's control. [T TestInputGroupBlockPrefixKeepsEditorAsFieldControl] |
| States | Disabled stops the editor and the addon actions. [T TestInputGroupDisabledStopsEditorAndAddonActions] |

### NumberField

| | |
| --- | --- |
| Role and name | A group holding Decrease, the editor and Increase. The editor is a text field, given a range only when Min or Max is set. [C] |
| Keys | Up and Down step by Step [T TestNumberFieldFractionalStep], PageUp and PageDown by ten steps, from the typed number. [T TestNumberFieldPageKeysStepTenSteps] A key step clamps at Min and Max. [C] Typing reports each finite number within the bounds and keeps the text as typed. [T TestNumberFieldClampsOnCommitNotPerKey] Enter commits: an out-of-range number clamps and text that does not parse reverts to the last value. [T TestNumberFieldClampsStepsAndRejectsText] Home and End move the caret. [C] |
| Pointer | The steppers step by Step and stop at Min and Max. [T TestNumberFieldClampsStepsAndRejectsText] A held stepper steps at once, then every 80 ms after 400 ms. [T TestNumberFieldStepperRepeatsWhileHeld] |
| Focus | The editor is the one Tab stop; the steppers leave the focus in it. [T TestNumberFieldSteppersKeepFocusInEditor] Leaving the editor commits as Enter does. [T TestNumberFieldClampsOnCommitNotPerKey] A Value the application sets replaces the text even while the editor has the focus, and leaving it then reports no change. [N TestNumberFieldAdoptsProgrammaticValueWhileFocused] |
| States | Steppers disable at their bound, when Disabled and when ReadOnly. [C] |
| Motion | Stepper press feedback keeps the hit box. [T TestNumberFieldStepperKeepsHitBoxOnPress] |

### InputOTP

| | |
| --- | --- |
| Role and name | One hidden native editor over painted slots; the name falls back to the localized "Verification code". [T TestFallbackNamesFollowMessages] |
| Keys | Characters outside the filter are dropped; OnComplete runs once per completion. [T TestInputOTPFiltersASCIIAndCompletesOnce] |
| States | Disabled takes no input. [T TestDisabledTextControlsTakeNoInput] ReadOnly keeps the code. [T TestInputOTPReadOnlyPreservesCode] |
| Motion | The painted caret blinks; under reduced motion it stays solid. [C] |

### MaskInput

Also covers: MaskFormat (a pure formatter).

| | |
| --- | --- |
| Role and name | The placeholder underscores are hidden from assistive technology. [C] |
| Keys | Cmd+Z undoes and Cmd+Shift+Z or Cmd+Y redoes through the raw values, a run of typing being one step. [T TestMaskInputUndoAndRedoStepThroughRawValues] Backspace or Delete over a literal removes the character beyond it. [T TestMaskInputDeletesAcrossLiterals] |
| Selection and value | The caller's value holds the raw characters; literals are inserted as typed. [T TestMaskInputWritesLiteralsAndStoresRawValue] Pasting at least as many characters as slots replaces the value and keeps the last ones. [T TestMaskInputPasteKeepsTheLastDigitsOfALongNumber] [T TestMaskInputLongPasteReplacesAPartlyFilledValue] |
| States | ReadOnly and Disabled answer no keys. [C] |

### TagsInput

| | |
| --- | --- |
| Role and name | The editor is a text field; each tag is a list item named by its text, with its position in the set [T TestTagsInputTagsAreItemsOfASet], and the focus marks the current tag; each remove button is named "Remove" plus the tag. [C] |
| Keys | Enter or a comma commits the draft. [T TestTagsInputCommitsOnEnterAndComma] Pasted text splits at commas, tabs and lines. [T TestTagsInputPasteSplitsAtCommasTabsAndLines] Escape clears a non-empty draft. [N TestTagsInputEscapeClearsDraft] Backspace in an empty editor removes the last tag. [T TestTagsInputBackspaceRemovesLastOnlyWhenEmpty] Left enters the tags, the arrows move between them, and Right past the last returns to the editor; Backspace or Delete on a tag removes it. [T TestTagsInputKeyboardWalksTheTokens] The arrows mirror in RTL; after a removal the focus goes to the previous tag, the next, or the editor when none is left. [C] |
| Focus | One Tab stop: tags and remove buttons are skipped by Tab, and a remove button keeps the focus where it was. [C] |
| Dismissal | Leaving the field commits the draft unless KeepDraftOnBlur. [T TestTagsInputCommitsDraftWhenFocusLeaves] |
| States | At Max no tag is added until one is removed. [T TestTagsInputStopsAtMax] The editor turns read-only and shows the count. [C] A refused duplicate stays in the editor [T TestTagsInputKeepsAndFlagsARefusedDraft], marks the field invalid and tints the matching tag. [C] |
| Motion | A removed tag fades out. [T TestTagsInputChipFadesOnRemoval] |

### KeyValueEditor

| | |
| --- | --- |
| Role and name | A group; editors named "Key N" and "Value N"; remove buttons "Remove KEY"; secret toggles "Show password KEY". [T TestKeyValueEditorSecretConcealsUntilRevealed] |
| Keys | Enter in a key moves to its value; Enter in the last non-empty value asks for a row and focuses its key. [T TestKeyValueEditorEnterMovesOnAndAddsRow] Pasted KEY=VALUE lines go to OnPaste. [T TestKeyValueEditorPasteReportsEnvLines] |
| Pointer | Add focuses the new row's key. [T TestKeyValueEditorAddFocusesNewRowAndRemoveDropsItsRow] |
| Focus | Rows with an ID keep focus and caret when earlier rows go. [T TestKeyValueEditorRowIDKeepsEditingStateAcrossRemoval] Removing a row with its button moves the focus to the remove button now in its place, else the previous one, else the add button. [T TestKeyValueEditorRemoveMovesFocusToNeighbour] |
| Selection and value | Pairs are never written; edits are reports the caller applies, and a declined edit is dropped. [T TestKeyValueEditorNeverWritesCallerPairs] A removal never reveals the secret that moves into its place. [T TestKeyValueEditorRemovingARowKeepsLaterSecretsConcealed] |
| States | Nil OnChange makes every row read-only. [C] Duplicate keys are invalid with a message. [T TestKeyValueEditorMarksDuplicateKeys] |

### Composer

| | |
| --- | --- |
| Role and name | Multiline text field named by Label, else Placeholder; the submit button defaults to the localized "Submit". [T TestFallbackNamesFollowMessages] |
| Keys | Enter submits; Shift+Enter inserts a line. [T TestComposerSubmitsCallerDraftAndKeepsNewlines] An empty draft does not submit unless AllowEmpty. [T TestComposerDoesNotSubmitDisabledOrEmptyDraft] |
| Focus | The submit button keeps the focus in the editor. [T TestComposerSubmitButtonPreservesEditorFocus] |
| States | ReadOnly and Disabled neither change nor submit the draft. [T TestComposerReadOnlyPreservesDraft] [T TestComposerDisabledDoesNotAcceptInput] |

### EditableText

| | |
| --- | --- |
| Role and name | At rest a button named by AriaLabel or the value, valued with the value. [C] |
| Keys | Click, Enter, Space or F2 starts editing in an editor that takes the focus. [T TestEditableTextCommitsOnEnter] Enter commits and returns the focus to the text. [T TestEditableTextCommitsOnEnter] |
| Dismissal | Escape drops the draft [T TestEditableTextEscapeCancels] and returns the focus to the text. [N TestEditableTextEscapeReturnsFocusToText] Leaving the editor commits. [T TestEditableTextCommitsOnBlur] |
| States | Disabled does not start editing. [C] |

### PathField

| | |
| --- | --- |
| Role and name | The text field reads the whole path even while it is drawn truncated. [T TestPathFieldEditsAndReadsFullPathWhileTruncated] |
| Pointer | Browse calls OnBrowse and leaves the value to the caller. [T TestPathFieldBrowseCallsCallerAndKeepsValue] A drop sets the first path and reports it. [T TestPathFieldDropSetsPathAndReports] |
| Focus | At rest a long path is drawn middle-truncated with the full path in a tooltip; focus shows the editable text. [C] |
| States | ReadOnly still browses [T TestPathFieldReadOnlyStillBrowses] and takes no drop. [C] Disabled does not browse. [T TestPathFieldDisabledDoesNotBrowse] |

### HotkeyInput

| | |
| --- | --- |
| Role and name | A toggle button named by AriaLabel, valued with the shortcut as "Ctrl+K"-style text, pressed while recording. [T TestHotkeyInputExposesRecording] |
| Keys | Click or Enter starts recording [N TestHotkeyInputEnterStartsRecording], as does Space. [C] While recording a modifier plus a key records and stops; bare keys other than Tab are swallowed. [T TestHotkeyInputRecordsChord] Backspace or Delete clears when Clearable. [T TestHotkeyInputBackspaceClears] |
| Dismissal | Escape stops recording and keeps the shortcut. [T TestHotkeyInputEscapeCancelsRecording] Escape at rest is passed on. [C] Focus leaving stops recording. [C] |
| Motion | The key caps dim while recording. [T TestHotkeyCapsFadeWhileRecording] |

## Choices

### Checkbox

| | |
| --- | --- |
| Role and name | Check box named by its label, Mixed while Indeterminate. [T TestCheckboxIndeterminateAndDisabled] |
| Keys | Space toggles. [T TestCheckboxTogglesAndReports] |
| Pointer | The whole row toggles and takes the focus. [T TestCheckboxTogglesAndReports] |
| Focus | The ring shows on keyboard focus only. [C] |
| Selection and value | Activating an indeterminate box checks it. Nil Checked renders disabled. [T TestCheckboxIndeterminateAndDisabled] |
| Motion | The check draws in and fades out; the dash grows from its centre. [T TestCheckboxFadesAndDrawsTheMark] [T TestCheckboxDashGrowsFromItsCentre] |

### CheckboxGroup

| | |
| --- | --- |
| Role and name | A group named by its label; the select-all row is mixed when some are chosen. [T TestCheckboxGroupSelectAll] |
| States | Disabled disables every box. [C] |

### RadioGroup

Also covers: Radio.

| | |
| --- | --- |
| Role and name | Radio group named by its label; each radio named by its label. [T TestRadioGroupSelectsOne] |
| Keys | Arrows, Home and End choose; the arrows wrap (toolkit). [N TestRadioGroupTabEntersAtChoiceAndKeysChoose] |
| Pointer | A click chooses and focuses. A disabled radio never chooses. [T TestRadioGroupSelectsOne] |
| Focus | One Tab stop entering at the chosen radio. [N TestRadioGroupTabEntersAtChoiceAndKeysChoose] |
| States | Nil State renders the radio disabled. [T TestRadioWithoutStateIsDisabled] A disabled group disables every radio. [T TestDisabledRadioGroup] |
| Motion | The dot scales in and fades out. [T TestRadioDotScalesInAndFadesOut] |

### Switch

Also covers: SwitchGroup.

| | |
| --- | --- |
| Role and name | Switch named by its label. [T TestSwitchTogglesAndReports] SwitchGroup is a labelled group of independent switches. [T TestSwitchGroup] |
| Keys | Space toggles. [T TestSwitchTogglesAndReports] |
| States | Disabled ignores presses. [T TestSwitchLabelPositionAndDisabled] A disabled SwitchGroup disables every switch. [C] Loading reports busy, keeps the value and ignores presses. [T TestSwitchLoadingKeepsValueAndBlocksPress] |
| Motion | The thumb slides and the track colour fades together over the State duration; reduced motion moves at once. [T TestSwitchThumbSlidesAndSettles] [T TestSwitchTrackSettlesWithTheThumb] A held thumb takes the slider thumb's press: it widens by `Effects.ThumbPressScale` (1.2) over the State duration, without a delay. [T TestSwitchThumbStretchesWhilePressed] [T TestSwitchHeldThumbGrowsAsASliderThumb] |

### CheckboxCard

Also covers: RadioCard.

| | |
| --- | --- |
| Role and name | The whole card is one check box (or radio) named by its title, the description read after it; Content is display-only. [T TestChoiceCardContentShowsAndKeepsTheCardOneControl] |
| Keys | Space toggles a CheckboxCard; RadioCards in a RadioGroup move with the arrows. [T TestRadioCardsShareOneChoiceInRadioGroup] |
| Pointer | A click anywhere on the card toggles. [T TestCheckboxCardTogglesFromAnywhereOnTheCard] |
| States | Disabled keeps the value. [T TestCheckboxCardDisabledKeepsValue] Invalid is painted and exposed. [T TestChoiceCardsExposeInvalid] |
| Motion | The border fades. [T TestChoiceCardBorderFades] |

### ImageSelect

| | |
| --- | --- |
| Role and name | With Value, a radio group; with Values, a group of check boxes. [C] |
| Keys | Single: one Tab stop, the arrows choose. [T TestImageSelectSingleChoosesWithClickAndArrows] Multiple: each tile is a Tab stop and Space toggles. [T TestImageSelectMultipleTogglesTiles] |
| Selection and value | Multiple reports a new slice and never edits the caller's. [C] Neither Value nor Values renders disabled. [C] |
| Motion | The chosen ring fades in. [T TestImageSelectRingFadesIn] |

### Rate

| | |
| --- | --- |
| Role and name | Slider from 0 to Count, valued "v / Count". [C] |
| Keys | Arrows step within bounds and Left/Right follow the direction; Home is none, End is all. [T TestRateArrowsStepWithinBounds] [T TestRateArrowsFollowDirection] PageUp and PageDown move a whole star. [T TestRatePageKeysMoveAWholeStar] Enter and Space do not pick the hovered star. [T TestRateKeyboardActivationIgnoresRestingPointer] |
| Pointer | Hover previews; a click sets whole or half stars; AllowClear clears on a click of the current value. [T TestRateClickChoosesStarAndHalf] |
| States | ReadOnly is not focusable and ignores input. [T TestRateReadOnlyIgnoresInput] |
| Motion | The pointer preview is instant; a set rating fills with motion. [T TestRateStarsFillAsTheRatingChanges] |

### Slider

| | |
| --- | --- |
| Role and name | Slider named by its label, with format and description exposed. [T TestSliderFormatAndDescription] |
| Keys | Arrows step, Home and End go to the bounds (toolkit). [T TestSliderDragAndKeyboard] PageUp and PageDown move a tenth of the range (toolkit). [N TestSliderPageKeysStepATenth] |
| Pointer | Drags change the value live; a range press moves the nearest thumb and the thumbs never cross. [T TestSliderRangeKeepsThumbsOrdered] |
| Selection and value | A range build orders and clamps the caller's low and high without calling OnChange. [C] |
| States | Disabled ignores drags and keys. [T TestSliderDisabled] |
| Motion | The fill eases to a value the application sets and follows a held pointer at once. [T TestSliderFillEasesToAnApplicationValue] |

### RangeSlider

| | |
| --- | --- |
| Role and name | Each thumb is a slider named "<Label> minimum" and "<Label> maximum". [T TestRangeSliderKeyboardMovesFocusedThumb] |
| Keys | Each thumb is a Tab stop; the keys move the focused thumb. [T TestRangeSliderKeyboardMovesFocusedThumb] A thumb stops at MinGap. [T TestRangeSliderMinGapStopsEitherThumb] |
| Focus | Either thumb's focus rings both thumbs (toolkit limit, see Gaps). [C] |
| States | Disabled ignores keys. [T TestRangeSliderDisabledIgnoresKeys] |

## Popup choices

### Select

| | |
| --- | --- |
| Role and name | Pop-up button named by Label, else Placeholder, expanded while open, valued with the chosen labels. [C] |
| Keys | Closed: Down, Enter and Space open at the chosen or first option. [C] Up opens at the last option, and a letter chooses the first option starting with it without opening. [N TestSelectClosedKeysOpenAtEndsAndLettersChoose] Within a second the same letter cycles through the options starting with it, and other letters extend the prefix. [T TestSelectClosedTypeaheadCyclesAndMatchesPrefix] Backspace and Delete clear when Clearable. [C] Open: the list's own navigation and typeahead apply. [C] |
| Focus | The trigger is the Tab stop; opening moves the focus into the list; a single choice closes and returns it to the trigger. [T TestSelectKeyboardCommitReturnsFocus] |
| Dismissal | Escape and an outside click close the popup (toolkit). Escape returns the focus to the trigger. [T TestSelectEscapeReturnsFocusToTrigger] |
| Selection and value | Multiple keeps the popup open. [T TestSelectMultipleKeepsPopupOpen] A row without a key is never chosen. [T TestSelectRowWithoutKeyIsNeverChosen] The clear button empties the choice. [T TestSelectSingleChoiceAndClear] |
| States | Disabled does not open. [T TestSelectDisabledDoesNotOpen] Loading swaps the chevron for a spinner and shows a Loading row in an empty popup. [T TestLoadingPopupRowsFollowOptions] |

### ComboBox

| | |
| --- | --- |
| Role and name | Combo box text field, expanded while open. [C] The option the arrows are on is the field's active descendant, so the field keeps the focus while assistive technology follows the arrows. [T TestTextPopupsPointTheirInputAtTheHighlightedOption] |
| Keys | Typing filters and opens. [T TestComboBoxListFollowsQueryAndNewOptions] Down and Up open or move the highlight, skipping disabled rows; Enter chooses it and keeps the focus. [T TestComboBoxKeyboardChoosesAndKeepsFocus] [T TestEditablePopupKeyboardFiltersGroupedOptions] Enter with no match keeps the text as Custom when AllowsCustomValue, else only closes. [C] Alt+Down opens and Alt+Up closes, keeping the text. [T TestComboBoxAltArrowsOpenAndClose] |
| Focus | The editor keeps the focus and the Tab stop; focusing it opens the popup. [T TestComboBoxEscapeClosesAndStaysClosed] Chip remove buttons are Tab stops. [C] |
| Dismissal | Escape closes and the popup stays closed until typing. [T TestComboBoxEscapeClosesAndStaysClosed] An outside click closes it for good. [T TestComboboxesOutsideClickClosesAndStaysClosed] Escape is taken before an ancestor's shortcut; on a closed popup it clears a Clearable field and is otherwise passed on. [T TestComboBoxEscapeOnClosedPopupClears] |
| Selection and value | A pointer choice applies once. [T TestComboBoxPointerChoosesOnce] OnChange may rewrite the text. [T TestComboBoxSelectionCallbackCanEditText] |
| States | Disabled closes and ignores keys. [C] Loading as Select. [T TestLoadingPopupRowsFollowOptions] |

### Autocomplete

| | |
| --- | --- |
| Role and name | As ComboBox. [C] The highlighted option is the editor's active descendant. [T TestTextPopupsPointTheirInputAtTheHighlightedOption] |
| Keys | Typing filters and opens; Enter chooses the highlight, writes its label and reports query and selection. [T TestAutocompleteEditsFiltersAndChoosesWithKeyboard] The arrows skip disabled rows. [T TestAutocompleteArrowsSkipDisabledAndPreserveFocus] Alt+Down and Alt+Up as ComboBox. [T TestAutocompleteAltArrowsAndEscapeOnClosedPopup] |
| Focus | The editor keeps the focus throughout; a chip's remove button and the clear button are reachable with Tab and return the focus to the editor. [T TestAutocompleteMultipleChipCanBeRemovedWithKeyboard] [T TestAutocompleteClearIsKeyboardAccessible] |
| Dismissal | Escape closes and typing reopens; focus moving away closes. [T TestAutocompleteEditsFiltersAndChoosesWithKeyboard] [T TestAutocompleteArrowsSkipDisabledAndPreserveFocus] Escape on a closed popup as ComboBox. [T TestAutocompleteAltArrowsAndEscapeOnClosedPopup] |
| Selection and value | Multiple keeps the popup open and the editor focused. [T TestAutocompleteMultipleSelectionKeepsEditorAndOptions] The query follows the caller's selection unless the caller set it. [T TestAutocompleteReflectsCallerSelection] |
| States | Disabled while open stops selection. [T TestAutocompleteDisabledWhileOpenStopsSelection] Loading as Select. [T TestLoadingPopupRowsFollowOptions] |

## Files

### Upload

Also covers: PickFiles (opens the system dialog; blocking, no UI of its own).

| | |
| --- | --- |
| Role and name | The drop zone is a group named by Title. The files are a list of rows named by the file; the done and failed icons are named. [T TestUploadFilesAreAListWithStatus] |
| Pointer | Dropped files are filtered by Accept, the first only without Multiple. [T TestUploadReportsAcceptedDroppedPaths] Choose and per-file Remove buttons call their handlers. [T TestUploadListRemovesAndChooses] |
| Focus | The zone is not focusable; the choose button is the keyboard path. [C] |
| States | A disabled zone takes no drop. [T TestUploadReportsAcceptedDroppedPaths] |
| Motion | The zone's border and fill fade while files hover. [C] |

## Gaps

Measured against platform conventions and the
[APG patterns](https://www.w3.org/WAI/ARIA/apg/patterns/). Each remaining gap
is a toolkit limit.

- NumberField: the editor is a text field with a range, not a spin button.
  Toolkit limit: MyGo's roles (`ui/access.go`) have no editable spin button;
  `RoleStepper` is AppKit's arrows-only stepper, and the editor must stay a
  text field to be edited.
- RangeSlider: no PageUp/PageDown, and focusing one thumb rings both.
  Toolkit limit: `ui.RangeSlider` (`coreRangeSlider` in `ui/indicators.go`)
  keeps its knobs internal and handles only the arrows, Home and End on them,
  so `form` can neither tell which knob has the focus nor add keys to it.
