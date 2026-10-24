// Package data provides HeroUI v3's data display components.
//
// Table is a native MyGo table: rows are virtualized by ui.ListState, the
// column header stays at the top while the rows scroll under it, columns are
// resized and reordered by the toolkit's own header gestures, and row
// selection is the toolkit's (ListState.Selected and ListState.Selection).
// The package adds HeroUI's column declaration, caller-owned table state,
// sorting, hidden columns, row actions, a footer slot and an empty slot.
//
// # Limits
//
// MyGo's Table owns the header geometry, so a few upstream details are
// approximations rather than CSS translations: the header text is 12 DIP
// semibold rather than 12 DIP medium, the header row is 32 DIP rather than
// 36, and its cell padding is the toolkit's 10 by 6. In the primary variant
// the header sits on the frame's secondary surface above the white body card;
// the secondary variant paints HeroUI's standalone header card over
// transparent cells. The toolkit paints a column separator only as the
// resize grip, so a column shows one when TableColumnOptions.Resizable is
// set. Body cells carry HeroUI's 16 by 12 inset and a hairline separator on
// each row boundary. A chosen row paints one tint across the full row,
// covering the toolkit's own accent face.
//
// Columns are resizable only when TableColumnOptions.Resizable is set. Because
// the toolkit expresses "not resizable" with TableColumn.Fixed, a column that
// is not resizable is also not movable. A table cannot show the column widths
// the user chose before the first frame; restore them through
// TableState.List.Columns.
//
// Expandable or tree rows are not implemented. MyGo's Table has no colspan,
// and its section headers pin, so a detail row cannot span the columns or be
// inserted without changing section semantics. A caller that needs details
// can build them inside a cell.
//
// Selection mode single tracks a row position, as ui.ListState.Selected does;
// selection mode multiple tracks stable row keys. A sort therefore keeps a
// multiple selection with its rows but remaps a single selection to the same
// row. Disabled rows are made non-selectable by correcting the toolkit's
// choice after the frame: a selection that lands on a disabled row is
// restored, and disabled keys are removed from a multiple selection. Pointer
// and keyboard input may therefore move through a disabled row for one frame.
// OnSelectionChange reports a change the user made: sorting or inserting rows
// moves the chosen row with its key without reporting it.
// Native cell blocks (ListState.CellSelect) and drag selection are available
// through TableState.List but are not exercised by this package's tests.
package data
