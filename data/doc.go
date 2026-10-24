// Package data provides HeroUI v3's data display components.
//
// Table is a native MyGo table: rows are virtualized by ui.ListState, the
// column header stays at the top while the rows scroll under it, columns are
// resized and reordered by the toolkit's own header gestures, and row
// selection is the toolkit's (ListState.Selected and ListState.Selection).
// The package adds HeroUI's column declaration, caller-owned table state,
// sorting, hidden columns, row actions and an empty slot.
//
// # Limits
//
// MyGo's Table owns the header and cell geometry, so a few upstream details
// are approximations rather than CSS translations: body text is 14 DIP,
// header text is 14 DIP semibold rather than 12 DIP medium, and the header
// shares the body surface because the toolkit paints the header and rows as
// one surface. TableSecondary keeps that surface instead of making its cells
// transparent. Row separators are the toolkit's, not HeroUI's per-cell
// hairlines.
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
// Native cell blocks (ListState.CellSelect) and drag selection are available
// through TableState.List but are not exercised by this package's tests.
package data
