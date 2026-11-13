# Content behavior

Package `content`. Format and evidence marks: [README](README.md). Durations
and easing: [`DESIGN.md`](../../DESIGN.md), "State motion". Text selection
and its copy shortcuts, scrolling, the native Markdown renderer, the native
diff viewer's editors and separators, and the outline's arrow keys belong to
the toolkit.

Shared by the code surfaces below (Highlighter, DiffViewer, JSONViewer,
LogViewer): code reads left to right in a right-to-left interface too. [C]
Every copy action reports a clipboard failure to OnError when it is set and
otherwise drops it. [T TestCopyButtonReportsClipboardFailure]

The toolbar search fields of JSONViewer and LogViewer behave alike: a text
field named "Search"; Escape with text clears it [T TestLogViewerToolbarSearches]
before a shortcut around the viewer, as a dialog's, sees the key, and keeps
the field focused [T TestViewerSearchEscapeClearsBeforeTheDialog];
Escape on an empty field is passed on to ancestors. [N TestViewerSearchEscapeOnEmptyQueryReachesAncestor]
A Clear button shows while the field holds text, empties it and leaves the
focus in the field. [T TestViewerSearchClearKeepsTheFocus] Cmd+F (Ctrl+F)
and F7 belong to the viewer, since a window can hold several; the state's
FocusSearch lets the application's Find focus one. [T TestViewerSearchFocusesForTheWindowsFind]
The match, filter and change counts are statuses, read out as they change. [T TestViewerCountsAreAnnounced]

## Copy

### CopyButton

| | |
| --- | --- |
| Role and name | Icon-only button named by Label, else the localized "Copy"; the same text is its tooltip. [C] The copy buttons of a titled Highlighter, DiffViewer or JSONViewer and of a labelled Snippet are named "Copy" and that name. [T TestCodeBlocksNameTheirGroupAndCopyButton] For the feedback time the name reads "Copied" and "Copied" is announced once. [T TestCopyButtonCopiesAndRevertsFeedback] |
| Keys | Enter and Space activate it (toolkit button). [C] |
| Pointer | A click writes Text, or what TextFunc returns at that moment, to the clipboard. [T TestCopyButtonCopiesAndRevertsFeedback] |
| Focus | One Tab stop; disabled, it is skipped. [C] |
| Selection and value | The check and "Copied" show only once the clipboard took the text, and revert by themselves. [T TestCopyButtonCopiesAndRevertsFeedback] A failed write shows no copied state. [T TestCopyButtonReportsClipboardFailure] |
| States | Disabled takes no press. [C] |
| Motion | The copy icon crossfades to the check and back. [C] |

### Snippet

| | |
| --- | --- |
| Role and name | A group named by Label; without it, no group. [T TestCodeBlocksNameTheirGroupAndCopyButton] The text is selectable, the prefix is not. [C] |
| Pointer | The copy button copies the text without the prefix. [T TestSnippetCopiesTextWithoutPrefix] |
| Selection and value | Text that does not fit ends in an ellipsis; the whole text is still copied. [C] NoCopy leaves out the button. [C] |

## Code

### Highlighter

| | |
| --- | --- |
| Role and name | A group named by Title, else the language's name, else Language; with none of them, no group. [T TestCodeBlocksNameTheirGroupAndCopyButton] The code is one selectable text. [C] |
| Pointer | The copy button, in the header or floating in the top corner without one, copies the code without its surrounding blank lines. [T TestHighlighterUnknownLanguageRendersPlainAndCopies] |
| Selection and value | An unknown or empty language renders plain. [T TestHighlighterUnknownLanguageRendersPlainAndCopies] Without Wrap long lines scroll sideways (toolkit). [C] |

### DiffViewer

| | |
| --- | --- |
| Role and name | A group named by Title; without it, no group, but the same header. [T TestDiffViewerWithoutTitleKeepsItsControls] The header shows the added and deleted counts. [T TestDiffViewerCountsFollowTheTexts] A unified diff exposes the new text only; a split diff exposes "Old text" and "New text" (toolkit). [T TestDiffViewerCountsFollowTheTexts] The layout toggle is a toggle button named "Side by side", pressed while split. [C] |
| Keys | With Height, F7 steps to the next change [N TestDiffViewerF7StepsForward] and Shift+F7 to the previous, wrapping. [T TestDiffViewerStepsThroughChanges] Both answer only while the focus is inside the viewer. [C] |
| Pointer | Previous and Next step through the changes and show "n / total". [T TestDiffViewerStepsThroughChanges] The layout toggle switches unified and split; Copy copies New. [C] Collapsed unchanged lines expand from their separators (toolkit). [C] |
| Focus | The header buttons are Tab stops before the text. [C] |
| Selection and value | Identical texts show "No changes" in place of the lines. [T TestDiffViewerCountsFollowTheTexts] The layout toggle is kept per viewer and resets when the Layout option changes. [C] |

### JSONViewer

Also covers: JSONViewerState.ExpandAll, JSONViewerState.CollapseAll.

| | |
| --- | --- |
| Role and name | A group and a tree named by Title, else "JSON"; each row is a tree item read as key, value or preview, and child count, expandable and expanded when it holds members. [C] The header's Expand all, Collapse all and Copy are icon buttons. [C] |
| Keys | The outline's arrows choose rows, Right opens and Left closes (toolkit). [C] Cmd+C (Ctrl+C) copies the chosen row's value: a string's characters, any other value's JSON. [T TestJSONViewerSearchRevealsMatchesAndCopies] With Toolbar, Cmd+F (Ctrl+F) focuses the search field while the focus is inside the viewer [N TestJSONViewerSearchFieldKeysStepAndClear]; there Enter and Shift+Enter step to the next and previous match, and Cmd+C does not copy the chosen row. [N TestJSONViewerSearchFieldKeysStepAndClear] |
| Pointer | A double click opens or closes a container. [T TestJSONViewerDoubleClickTogglesAContainer] A row's context menu copies its value or its path. [T TestJSONViewerCopiesValueAndPath] Expand all and Collapse all open every node or close all below the root. [T TestJSONViewerKeepsOrderAndExpansion] Previous and Next step through matches. [T TestJSONViewerSearchRevealsMatchesAndCopies] |
| Selection and value | Members keep the document's order; open nodes survive a new document by path. [T TestJSONViewerKeepsOrderAndExpansion] A query opens the containers above its first match, chooses it and shows "n / total". [T TestJSONViewerSearchRevealsMatchesAndCopies] A new document under the same query keeps the match the user is on. [T TestJSONViewerSearchKeepsItsMatchWhenTheSourceChanges] State is caller-owned when given, else kept per window. [C] Only the rows in view are built. [T TestJSONViewerBuildsOnlyRowsInView] |
| States | Source that does not parse returns the error and shows it in place of the tree [T TestJSONViewerReportsInvalidDocument], as an alert. [T TestJSONViewerAnnouncesAnInvalidDocument] The header buttons are disabled meanwhile. [C] Previous and Next are disabled without matches. [C] |

### LogViewer

Also covers: LevelSet, LogLevels.Shows.

| | |
| --- | --- |
| Role and name | A group and a list named by Label, else the locale's Log; each line is a list item read as its text, then its level for Debug, Warning and Error lines. [T TestLogLinesAreReachableByKeyboard] The level and wrap toggles are toggle buttons named by level ("Debug", "Info", "Warning", "Error") and "Wrap lines". [C] |
| Keys | With Toolbar, Cmd+F (Ctrl+F) focuses the search field while the focus is inside the viewer [T TestLogViewerToolbarSearches], also when nothing matches. [T TestLogViewerSearchShortcutWorksWithoutResults] |
| Pointer | A line's context menu copies the line or every line listed. [C] Jump to latest scrolls back to the end [T TestLogViewerFollowsUntilScrolledAway] and leaves the focus on the lines. [T TestLogJumpToLatestLeavesTheFocusOnTheLines] |
| Focus | The lines are one Tab stop: Up and Down choose a line and the context-menu key (Shift+F10) opens its menu. [T TestLogLinesAreReachableByKeyboard] The text of each is selectable and a selection does not span lines. [C] |
| Selection and value | Follow keeps the newest line in view while the user is at the end, stops once they scroll away and shows Jump to latest. [T TestLogViewerFollowsUntilScrolledAway] Query and Levels list only matching lines, case-insensitively, under their own numbers; nothing listed shows "No results". [T TestLogViewerQueryAndLevelsListMatchingLines] Appended lines are filtered as they arrive. [T TestLogViewerFilterFollowsAppendedLines] A plain line counts as Info. [T TestLogViewerFilterKeepsScrollOnRevisionAndListsPlainAsInfo] Raising Revision lists again without moving the scroll [T TestLogViewerFilterKeepsScrollOnRevisionAndListsPlainAsInfo]; changing the query or levels returns to the start, or to the end while following. [C] The toolbar shows "listed / total" while a filter is on. [T TestLogViewerToolbarSearches] |
| States | Level tints Debug, Warning and Error lines, whose names read the level. [T TestLogLinesAreReachableByKeyboard] |
| Motion | Jump to latest fades in and out. [T TestLogViewerJumpButtonFadesIn] The toolbar at rest asks for no frames. [T TestLogViewerToolbarSearches] |

## Rendered content

### Markdown

| | |
| --- | --- |
| Role and name | Headings, links, lists and tables are exposed by the native renderer (toolkit). [C] |
| Pointer | With OnLink a link click is reported and not opened. [T TestMarkdownReportsLinksAndHighlightsCodeBlocks] Without it web and mail links open and heading anchors scroll (toolkit). [C] |
| Selection and value | Code blocks render as a Highlighter with its copy button. [T TestMarkdownReportsLinksAndHighlightsCodeBlocks] Selection copies as text, HTML or Markdown (toolkit). [C] |

### QRCode

| | |
| --- | --- |
| Role and name | An image named by Label, else Value. [T TestQRCodePaintsModulesInsideQuietZone] |
| Selection and value | A value that cannot be encoded returns the error and draws nothing. [T TestQRCodeReportsUnencodableValue] An Icon raises the level to QRHigh. [C] Colours default to the light palette in both appearances. [C] |

## Gaps

Measured against platform conventions and the
[APG patterns](https://www.w3.org/WAI/ARIA/apg/patterns/).

- The counts and the invalid-document alert are verified through the
  toolkit's live-region announcements headlessly; not yet with a real screen
  reader.
