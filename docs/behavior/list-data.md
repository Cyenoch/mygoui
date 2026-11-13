# List and data behavior

Packages `list` and `data`. Format and evidence marks: [README](README.md).
Durations and easing: [`DESIGN.md`](../../DESIGN.md), "State motion" and
"Loading". "Cmd" is Command on macOS and Control elsewhere.

Every component here is built on MyGo's virtual list or table. Scrolling,
virtualization, row measurement and wheel handling belong to the toolkit;
so do the Table's keyboard, row and cell selection and header drags.

## List

### ListBox

| | |
| --- | --- |
| Role and name | A list box named by Label, multiselectable in multiple mode; each row is an option named by its label, Selected when chosen and Disabled when disabled. A section label is a level-3 heading. [T TestListBoxIsAListBoxOfOptions] The highlighted row is the active descendant. [T TestListBoxIsAListBoxOfOptions] |
| Keys | Down and Up move the highlight, skipping headers and disabled rows; Home and End go to the first and last enabled row; Enter chooses the highlighted row. [T TestListBoxSectionHeaderKeyboardNavigation] [N TestListBoxKeysActivateTheHighlightedRow] Space and Enter run OnAction on the highlighted row; an action-only list chooses nothing. [N TestListBoxKeysActivateTheHighlightedRow] Space toggles the highlighted row of a multiple selection. [N TestListBoxSpaceTogglesInMultipleSelection] A letter A–Z or a digit moves to the next enabled row starting with it, skipping headers [T TestListBoxTypeaheadSkipsHeadersAndDisabled]; the same letter again cycles and wraps. [N TestListBoxTypeaheadCyclesThroughMatchingRows] Letters typed within a second build a prefix; a pause starts anew. [T TestListBoxTypeaheadBuildsAPrefix] Other letters are a toolkit limit (see Gaps). PageUp and PageDown move the highlight by the rows in view, except on macOS, where they scroll as the toolkit's lists do. [T TestListBoxPagesFollowThePlatform] Multiple: Shift+Down and Shift+Up choose from the anchor to the highlight; Cmd+A chooses every enabled row. [T TestListBoxExtendsAMultipleSelectionFromTheKeyboard] The arrows stop at the ends; Escape is not handled. [C] |
| Pointer | A click on an enabled row activates it as Enter does; a disabled row ignores it. [T TestListBoxSingleSelectionAndAction] A click toggles a row of a multiple selection. [T TestListBoxMultipleSections] Hover paints the row and reverts. [T TestListBoxRowHoverPaint] With KeyboardHighlight an outside editor drives the highlight, and the pointer moves the same highlight. [T TestListBoxEditorHighlightPaintsOnceAndMoves] |
| Focus | One Tab stop; rows never take the focus. [N TestListBoxKeysActivateTheHighlightedRow] A click focuses the list unless NoFocusOnClick. [C] The highlighted row is painted and ringed while the focus is visible, and is the active descendant. [T TestListBoxIsAListBoxOfOptions] With Editor it is the editor's active descendant too, as the editor keeps the focus. [T TestTextPopupsPointTheirInputAtTheHighlightedOption] [T TestCommandInputPointsAtTheHighlightedRow] |
| Selection and value | The caller owns State. Single mode holds SelectedKey, multiple mode Selected; OnSelectionChange reports a user change. [T TestListBoxSingleSelectionAndAction] Choosing the chosen row again reports no change and still runs OnAction. [C] The choice is marked by a check at the end, or at the start with IndicatorStart, and the row keeps its surface. [T TestListBoxChosenRowKeepsSurfaceAndShowsIndicator] [T TestListBoxLeadingIndicatorReservesSpaceAndMarksChoice] Section headers scroll with their rows. [T TestListBoxSectionHeaderScrollsAway] With no rows the Empty text shows and the keys do nothing. [T TestListBoxEmptyShowsMessage] |
| States | Disabled dims the list, removes it from the Tab order and binds no keys. [C] Loading with no rows shows skeleton rows and no empty text [T TestListBoxLoadingFirstLoadShowsNoEmptyMessage]; with rows it keeps them and ignores clicks. [T TestListBoxLoadingRefreshKeepsRowsAndIgnoresInput] The list reports busy while loading. [C] |
| Motion | The check scales in while fading in and leaves by fading; under reduced motion only the fade remains; a row built chosen shows its check at once. [T TestListBoxCheckFadesIn] Hover and highlight are instant. [C] |

### Tree

| | |
| --- | --- |
| Role and name | A tree named by Label; each row a tree item with its level and position in its set, Expanded on an open branch and Selected when chosen. [T TestTreeExposesHierarchyToAssistiveTechnology] A Checkable row reports Checked or Mixed in place of Selected. [T TestTreeChecksCascadeAndMixParents] The keyboard row is the active descendant. [T TestTreeReportsTheKeyboardRowToAssistiveTechnology] A loading node is busy. [T TestTreeLazyBranchesAskForTheirChildrenAndReportProgress] A disabled row is exposed as disabled [C]; a multiple selection marks the tree multiselectable. [T TestTreeMarksAMultipleSelectionAndCanBeDisabled] |
| Keys | Down and Up skip disabled rows; Right opens a branch, then moves to its first child; Left closes it, then moves to the parent; Home goes to the first row; Enter activates. [T TestTreeKeyboardExpandsMovesAndSelects] End goes to the last enabled row; Right and Left swap in RTL. [C] Space checks a Checkable row and its subtree, a second press clearing what it can reach. [T TestTreeChecksCascadeAndMixParents] Otherwise Space chooses the row [C] without running OnAction [T TestTreeSpaceSelectsWithoutRunningOnAction], and with no selection opens or closes a branch. [C] Letters and digits typed within a second, the toolkit's pause, build a prefix; one letter repeated steps on. [T TestTreeTypeaheadMatchesAPrefixWithinTheWindow] Shift+F10 and the menu key open the keyboard row's menu, choosing it first. [T TestTreeMenuKeyOpensTheKeyboardRowsMenu] PageUp and PageDown move a page except on macOS, where they scroll. [T TestTreePagesFollowThePlatform] Multiple: Shift+Down and Shift+Up extend from the anchor, Cmd+A chooses every enabled row in view, Escape keeps the selection. [T TestTreeMultipleSelectionFollowsDesktopRules] |
| Pointer | A click chooses the row; with no selection it checks a Checkable row or opens and closes a branch. [T TestTreeKeyboardExpandsMovesAndSelects] [T TestTreeChecksCascadeAndMixParents] [T TestTreeLazyBranchesAskForTheirChildrenAndReportProgress] Cmd-click toggles a row and Shift-click selects the range, skipping disabled rows. [T TestTreeMultipleSelectionFollowsDesktopRules] A double click opens a branch while choosing. [T TestTreeDoubleClickOpensABranchWhileSelecting] The arrow opens and closes and the check box checks without choosing. [C] A secondary click chooses an unchosen row before OnContextMenu fills the menu; disabled rows open none. [T TestTreeContextMenuChoosesTheRowFirst] Disabled rows take no click. [T TestTreeKeyboardExpandsMovesAndSelects] A truncated label shows its full text in a tooltip. [C] |
| Focus | Tab focuses the tree and the keys then act on it [T TestTreeTypeaheadMatchesAPrefixWithinTheWindow]; it is one Tab stop. [C] A click on a row, arrow or check box focuses the tree and moves the keyboard row there. [C] Closing a branch that holds the keyboard row moves it to the branch. [C] |
| Selection and value | The caller owns State: expansion, selection and checks. OnExpandedChange reports each open and close. [T TestTreeKeyboardExpandsMovesAndSelects] A branch is checked when every enabled leaf is and mixed when some are; disabled leaves keep their check. [T TestTreeChecksCascadeAndMixParents] An unloaded Lazy branch checks as one unit. [T TestTreeUnloadedLazyBranchChecksAsAUnit] Only the rows in view are built; Revision shows edits made in place. [T TestTreeBuildsOnlyTheViewportAndFollowsRevision] |
| States | A disabled node and its subtree are never chosen, checked or toggled. [C] Disabled dims the whole tree, takes it out of the Tab order and stops its keys and clicks. [T TestTreeMarksAMultipleSelectionAndCanBeDisabled] Loading with no rows shows indented skeleton rows [T TestTreeLoadingFirstLoadShowsNoNodes]; with rows it keeps them and ignores clicks. [T TestTreeLoadingRefreshKeepsRowsAndIgnoresInput] A node's Loading swaps its arrow for a spinner. [C] |
| Motion | The arrow turns with the branch [T TestTreeArrowTurnsWithTheBranch] over `Motion.Reveal` with `Motion.Ease`. [C] Opening fades the new rows in and slides the rows below; reduced motion keeps only the fade; closing drops the rows at once. [C] The check draws in; a row built checked shows it at once. [T TestTreeCheckDrawsIn] |

### SortableList

| | |
| --- | --- |
| Role and name | A list named by Label of list items, each named by its label with its position in the set; each row's grip is a button named "Reorder" plus the item's label (localized), which tells its Alt+Up and Alt+Down keys. [T TestSortableListMovesFromTheKeyboard] [T TestSortableListTellsAssistiveTechnology] A move is read out as the item's label and its new place ("Alpha, 2 of 3", localized). [T TestSortableListTellsAssistiveTechnology] |
| Keys | Alt+Down and Alt+Up on a focused grip move its item one place; the grip keeps the focus as its row moves. [T TestSortableListMovesFromTheKeyboard] At either end the key does nothing. [C] |
| Pointer | Dragging a grip moves the row with the pointer while the others make room; the drop reports the move. [T TestSortableListMovesByDraggingTheHandle] The held row is held between the first and last rows and lifted with a shadow. [C] |
| Focus | Each grip is a Tab stop. [C] |
| Dismissal | Escape cancels a drag: the row goes back and the release moves nothing. [T TestSortableListEscapeCancelsADrag] A new item list under a drag cancels it. [C] |
| Selection and value | Items are never changed: OnMove reports from and to, counted after the removal, and the caller reorders. [T TestSortableListMovesByDraggingTheHandle] |
| States | Disabled, or a nil OnMove, disables the grips: no drag and no keys. [T TestSortableListMovesByDraggingTheHandle] |
| Motion | Rows moving aside slide; after a drop the rows settle at once; reduced motion settles them at once. See "A row or tab moving aside during a drag". [C] The held row follows the pointer at once. [C] |

## Data

### Table

| | |
| --- | --- |
| Role and name | A table named by Label, with column headers; the sorted header reports its direction; rows are selectable and the table is multiselectable in multiple mode (toolkit). [C] A disabled row disables its cells' content; the row itself is not exposed as disabled (toolkit limit, see Gaps). [C] While loading, the table reports busy. [T TestTableReportsBusyAndCanBeDisabled] The default empty text is a status region. [C] |
| Keys | The toolkit's: arrows choose a row once the table has the focus [T TestTableKeepsNativeFocusRing], and with CellSelect the arrows move the cell cursor. [T TestTableNativeCellSelectionKeepsItsPaint] Enter reports OnRowActivate only while rows are selectable. [T TestTableRowActions] The keys that move the choice carry it past a disabled row to the next enabled one; Home and End stop on the nearest enabled row. [T TestTableArrowsPassADisabledRow] [T TestTableHomeStopsOnTheNearestEnabledRow] With Shift the range extends past a disabled row, and stays as it was when only disabled rows lie further. [N TestTableShiftArrowKeepsTheRangeAtADisabledEnd] |
| Pointer | A sortable header click sorts ascending, then descending, and reports OnSort. [T TestTableSortOrdersRowsAndNotifies] A Resizable column's edge drags its width. [T TestTableColumnResize] A click chooses a row and Cmd-click adds one. [T TestTableSelectionSingleAndMultiple] A disabled row is never chosen. [T TestTableDisabledRowCannotBeChosen] A double click reports OnRowActivate; RowMenu fills a row's context menu. [C] The header stays fixed while the wheel scrolls rows [T TestTableWheelScrollKeepsHeader]; narrow tables scroll sideways. [T TestTableOverflowScrollsSideways] |
| Focus | One Tab stop with the native focus ring. [T TestTableKeepsNativeFocusRing] Controls in cells take their own focus. [T TestTableSecondaryKeepsNativeCellControlsAndVirtualizedRows] |
| Selection and value | The caller owns State; a nil State keeps a private one that still scrolls and sorts. [C] Single mode chooses a position, multiple mode keys; a sort or an insertion keeps a keyed choice on its row without reporting it. [T TestTableSelectionSingleAndMultiple] [T TestTableInsertionKeepsKeyedSelection] Disabled keys never stay in a multiple selection. [T TestTableMultipleSelectionDropsDisabledKeys] Hidden columns and a table with every column hidden show the empty slot. [T TestTableEmptyStateAndHiddenColumn] |
| States | Loading with no rows shows skeleton rows under the real header and builds no caller cell [T TestTableLoadingFirstLoadKeepsHeaderAndBuildsNoCaller]; with rows it keeps them and ignores clicks. [T TestTableLoadingRefreshKeepsRowsAndIgnoresInput] Disabled dims the whole table, takes it out of the Tab order and stops its keys and clicks. [T TestTableReportsBusyAndCanBeDisabled] |
| Motion | None of its own; loading follows "Loading". Hover is instant. [C] |

## Gaps

Measured against platform conventions and the
[APG patterns](https://www.w3.org/WAI/ARIA/apg/patterns/) (listbox, tree
view, grid). Each remaining gap is a toolkit limit.

- ListBox, Tree: typeahead takes A–Z, the digits and space only. Toolkit
  limit: a list that does not edit text receives keys, not text, and MyGo
  types a row's letters from those keys alone (`ui/input.go`, `typedRune`).
- Tree: `*` does not open the sibling branches. Toolkit limit: MyGo has no
  key for `*` (`ui/keys.go`), and shortcuts match a key and its modifiers,
  so the character `*` cannot be bound apart from a keyboard layout.
- Table: a disabled row is not exposed as disabled; only its cells' content
  is (probed). Toolkit limit: the table's rows are the toolkit's
  (`ui/access.go` gives them `RoleRow`), and `ui.ListState` has no per-row
  state the table could set.
