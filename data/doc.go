// Package data provides a COSS-inspired native table.
//
// TablePrimary is the plain table; TableSecondary puts its body on a rounded
// card. Both keep MyGo's virtualized ui.Table, fixed scrolling header, native
// keyboard focus, row selection, sortable headers, column resizing and
// reordering. The package supplies typed columns, caller-owned state, sorting,
// hidden columns, row actions, an empty slot and a footer slot.
//
// # Limits
//
// MyGo owns the header and cell layout: its header is 32 DIP high with
// semibold text and 10 by 6 DIP cell padding rather than the COSS header's
// 40 DIP height and medium weight. Body cells add 4 DIP vertical padding to
// the native 10 by 6 inset. The card surrounds the table's available body
// viewport, including empty space, rather than sizing to the rendered rows.
// MyGo supplies row and header hover paint. Selected rows use a full-width
// COSS-like tint; when native cell selection is enabled, native cell selection
// paint and its cursor border take precedence.
//
// The header's separator is the native resize grip and appears only for
// resizable columns. MyGo ties moving a column to resizing it, so fixed-width
// columns are not movable. Restore user widths and order through
// TableState.List.Columns. Expandable rows and colspan are not implemented;
// the footer is a separate slot, not a spanning table row.
//
// Single selection tracks a row position; multiple selection uses stable row
// keys. Sorting and inserting rows keep the selected item when RowKey supplies
// stable keys. Disabled rows are corrected after native selection, so pointer
// or keyboard input can briefly pass through a disabled row. OnSelectionChange
// reports user selection changes, not remapping after a sort or insertion.
// Native cell selection and drag selection are exposed through TableState.List;
// cell selection requires multiple selection mode.
package data
