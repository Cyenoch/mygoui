package data

import (
	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TableVariant is HeroUI's table surface variant.
type TableVariant uint8

const (
	// TablePrimary puts the table on a padded secondary surface, as HeroUI's
	// default: a white table card inside a gray rounded frame.
	TablePrimary TableVariant = iota
	// TableSecondary paints no frame, leaving the table flush with its parent.
	TableSecondary
)

// TableSelectionMode says whether and how rows can be chosen.
type TableSelectionMode uint8

const (
	// TableSelectionNone does not let the user choose rows.
	TableSelectionNone TableSelectionMode = iota
	// TableSelectionSingle chooses one row, by its position.
	TableSelectionSingle
	// TableSelectionMultiple chooses rows by their stable keys.
	TableSelectionMultiple
)

// TableSort is the column a table's rows are sorted by and its direction.
// Column is a column ID: TableColumnOptions.ID, or its Header when ID is
// empty. The zero value means unsorted.
type TableSort struct {
	Column     string
	Descending bool
}

// TableSelection is the rows chosen in a table. Mode is the selection mode;
// Row is the chosen row (single mode) or the row the keys move from (multiple
// mode), -1 for none; Rows holds the chosen keys in multiple mode.
//
// A caller-owned TableSelection keeps its selection across frames. The zero
// value chooses row 0 in single mode because Row is an int; set Row to -1 to
// start with no row chosen.
type TableSelection struct {
	Mode TableSelectionMode
	Row  int
	Rows ui.Selection[any]
}

// Clear chooses no row.
func (s *TableSelection) Clear() {
	s.Row = -1
	s.Rows.Clear()
}

// HasKey reports whether the row of key k is chosen in multiple mode.
func (s *TableSelection) HasKey(k any) bool { return s.Rows.Has(k) }

// TableState is a table's caller-owned state: its native scroll and selection
// state, its sort, the columns the caller hid, and a revision for changes to
// the rows that a frame cannot see.
//
// Revision must change when row data, a comparator's behavior or RowDisabled
// changes in place. Rows appearing or disappearing are noticed without it.
type TableState struct {
	// List is the native table state: its position, the widths and order the
	// user gave the columns (List.Columns), and native additions such as cell
	// selection and row dragging. Table sets List.Sort, List.Selected,
	// List.Selection and List.Key on every build.
	List ui.ListState
	// Selection says whether and how rows are chosen.
	Selection TableSelection
	// Sort is the column the rows are sorted by. Header clicks update it;
	// set it to sort or clear the sort from the app.
	Sort TableSort
	// Hidden holds the IDs of columns the caller hid, including columns that
	// are not Hideable. Persist and restore it with the app's settings.
	Hidden map[string]bool
	// Revision tells the table that row data or a comparator changed in
	// place. Increment it when a cell, a key, a comparator result or a
	// disabled flag changes without the row count changing.
	Revision uint64

	order      []int
	orderRows  int
	orderRev   uint64
	orderSort  TableSort
	ordered    bool
	nativeSort ui.SortOrder
	lastSort   TableSort

	disabled     []int
	disabledSet  map[int]struct{}
	disabledRows int
	disabledRev  uint64
	disabledHas  bool
	disabledDone bool

	leadUnder   int
	snapshot    []any
	snapshotLen int
	snapshotOK  bool
}

// HideColumn hides the column of id, if it is shown.
func (s *TableState) HideColumn(id string) {
	if s.Hidden == nil {
		s.Hidden = map[string]bool{}
	}
	s.Hidden[id] = true
}

// ShowColumn shows the column of id again.
func (s *TableState) ShowColumn(id string) {
	delete(s.Hidden, id)
}

// ToggleColumn hides the column of id, or shows it when it is hidden.
func (s *TableState) ToggleColumn(id string) {
	if s.Hidden[id] {
		s.ShowColumn(id)
		return
	}
	s.HideColumn(id)
}

// ColumnHidden reports whether the caller hid the column of id.
func (s *TableState) ColumnHidden(id string) bool { return s.Hidden[id] }

// TableColumnOptions declares a column. ID is its stable identity, for the
// sort and for the hidden set; Header is its title and the fallback ID. Width
// 0 shares the room the other columns leave, MinWidth and MaxWidth bound it,
// and Cell builds a body cell. A nil Cell falls back to TableOptions.Cell.
type TableColumnOptions struct {
	ID       string
	Header   string
	Align    ui.Align
	Width    float32
	MinWidth float32
	MaxWidth float32
	// Sortable makes a header click sort by this column.
	Sortable bool
	// Resizable lets the user drag the column's edge. It also keeps the
	// column movable; a non-resizable column is fixed, because the native
	// table expresses both with one flag.
	Resizable bool
	// Hideable marks the column as one the caller may hide. The table does not
	// add a visibility menu; callers use TableState.HideColumn and friends.
	Hideable bool
	// Hidden hides the column regardless of TableState.Hidden.
	Hidden bool
	Cell   func(c *ui.Context, cell TableCell) ui.Element
}

// TableColumn is the resolved form of a column, passed to cell builders. Index
// is its position among the visible columns.
type TableColumn struct {
	Index      int
	ID         string
	Header     string
	Align      ui.Align
	Width      float32
	MinWidth   float32
	MaxWidth   float32
	Sortable   bool
	Resizable  bool
	Hideable   bool
	Sorted     bool
	Descending bool
}

// TableRow identifies a body row while building. Index is the caller's row
// index, Key is TableOptions.RowKey's value (the index when RowKey is nil),
// and Selected and Disabled are the row's current state.
type TableRow struct {
	Index    int
	Key      any
	Selected bool
	Disabled bool
}

// TableCell is one body cell passed to a cell builder.
type TableCell struct {
	Row    TableRow
	Column TableColumn
}

// TableOptions configures a table.
type TableOptions struct {
	// Label is the table's accessible name.
	Label string
	// Variant selects the framed (default) or flush surface.
	Variant TableVariant
	// Columns declares the columns, left to right. Hidden columns are
	// skipped. IDs must be unique.
	Columns []TableColumnOptions
	// Rows is the number of caller rows, built only as they scroll into view.
	Rows int
	// State is the caller-owned state. When nil the table keeps a private
	// state for one window: it still scrolls and its headers still sort, but
	// the app cannot read or set its selection or sort.
	State *TableState
	// Cell builds a cell for a column whose Cell is nil.
	Cell func(c *ui.Context, cell TableCell) ui.Element
	// RowKey returns a stable, comparable identity for caller row i. Without
	// it, the row's index is its key and the selection moves with positions
	// when rows are inserted above it.
	RowKey func(row int) any
	// RowDisabled reports whether caller row i is disabled. Disabled rows are
	// dimmed, their cell content is disabled, and a selection that lands on
	// one is corrected after the frame.
	RowDisabled func(row int) bool
	// Compare orders rows for State.Sort, returning a negative number when a
	// comes first. It is stable and is not called when there is no sort. Bump
	// State.Revision when its results change without the row count changing.
	Compare func(a, b TableRow) int
	// OnSort reports a sort the user chose by clicking a header.
	OnSort func(sort TableSort)
	// OnSelectionChange reports a selection the user changed. State.Selection
	// is authoritative and is updated before the callback.
	OnSelectionChange func(selection TableSelection)
	// OnRowActivate reports a double click or Enter on the chosen row. The
	// native table reports this only while rows are selectable, so it requires
	// TableSelectionSingle or TableSelectionMultiple.
	OnRowActivate func(row TableRow)
	// RowActions builds a trailing, fixed actions cell. The table adds the
	// column itself; ActionsHeader and ActionsWidth configure it.
	RowActions func(c *ui.Context, row TableRow) ui.Element
	// RowMenu builds a row's context menu, as ui.ListState.CellMenu does.
	RowMenu func(m *ui.Menu, row TableRow)
	// ActionsHeader titles the actions column.
	ActionsHeader string
	// ActionsWidth is the width of the actions column in DIPs; 0 uses 64.
	ActionsWidth float32
	// Empty builds the empty state shown when Rows is 0 or every column is
	// hidden. Without it, the table shows the localized empty message.
	Empty func(c *ui.Context) ui.Element
}

// Table creates a HeroUI table. Rows are virtualized and the header is fixed
// while the rows scroll, through ui.Table and ui.ListState. The caller sizes
// the table or grows it in its parent:
//
//	state := data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionMultiple, Row: -1}}
//	data.Table(c, data.TableOptions{
//		Label:   "Files",
//		Columns: []data.TableColumnOptions{{ID: "name", Header: "Name", Sortable: true, Cell: ...}},
//		Rows:    len(files),
//		State:   &state,
//	}).Grow(1)
func Table(c *ui.Context, o TableOptions) ui.Element {
	root := ui.Column(c).FillWidth().MinWidth(0)
	st := o.State
	if st == nil {
		st = ui.Local[TableState](root, "table", func() TableState {
			return TableState{Selection: TableSelection{Row: -1}}
		})
	}

	t := theme.Colors(c)
	switch o.Variant {
	case TableSecondary:
	case TablePrimary:
		radius := min(theme.Space32, theme.RadiusBase()*2.5)
		root.Padding(0, theme.Space4, theme.Space4, theme.Space4).
			Radius(radius).Clip().
			Background(t.SurfaceSecondary).TextColor(t.SurfaceSecondaryForeground)
	}

	// The caller may have changed the sort; a header click before this frame
	// changed the native sort. Reconcile them before the rows are built.
	sortChanged := st.syncSort()
	st.syncDisabled(&o)
	if sortChanged && o.OnSort != nil {
		o.OnSort(st.Sort)
	}
	if o.Rows > 0 {
		st.syncOrder(&o)
	}

	tl := newTableLayout(c, o, st)
	if o.Rows <= 0 || len(tl.columns) == 0 {
		root.Children(func() { buildEmpty(c, o, st, tl) })
		return root
	}

	st.wire(&o)
	var table ui.Element
	root.Children(func() {
		table = ui.Table(c, &st.List, tl.nativeColumns(), o.Rows, func(row, col int) { tl.buildCell(c, row, col) }).
			Grow(1).MinWidth(0).
			Label(o.Label).
			Background(t.Surface).
			FontSize(theme.Scale(c, 14))
	})

	changed := st.afterBuild(&o, table.Changed(), sortChanged)
	if changed && o.OnSelectionChange != nil {
		o.OnSelectionChange(st.Selection)
	}
	if table.Submitted() && o.OnRowActivate != nil {
		if row, ok := st.activeRow(&o); ok {
			o.OnRowActivate(row)
		}
	}
	return root
}

// columnBuild is one visible column: its public description and its cell
// builder, or the actions column.
type columnBuild struct {
	column TableColumn
	cell   func(c *ui.Context, cell TableCell) ui.Element
}

type tableLayout struct {
	o       TableOptions
	st      *TableState
	columns []columnBuild
	actions int // visible index of the actions column, -1 for none
}

func newTableLayout(c *ui.Context, o TableOptions, st *TableState) *tableLayout {
	tl := &tableLayout{o: o, st: st, actions: -1}
	for _, co := range o.Columns {
		id := co.ID
		if id == "" {
			id = co.Header
		}
		if co.Hidden || st.Hidden[id] {
			continue
		}
		live := TableColumn{
			Index:      len(tl.columns),
			ID:         id,
			Header:     co.Header,
			Align:      co.Align,
			Width:      co.Width,
			MinWidth:   co.MinWidth,
			MaxWidth:   co.MaxWidth,
			Sortable:   co.Sortable,
			Resizable:  co.Resizable,
			Hideable:   co.Hideable,
			Sorted:     st.Sort.Column == id,
			Descending: st.Sort.Column == id && st.Sort.Descending,
		}
		tl.columns = append(tl.columns, columnBuild{column: live, cell: co.Cell})
	}
	if o.RowActions != nil {
		tl.actions = len(tl.columns)
	}
	return tl
}

func (tl *tableLayout) nativeColumns() []ui.TableColumn {
	cols := make([]ui.TableColumn, 0, len(tl.columns)+1)
	for _, cb := range tl.columns {
		col := cb.column
		cols = append(cols, ui.TableColumn{
			Title:    col.Header,
			ID:       col.ID,
			Width:    col.Width,
			MinWidth: col.MinWidth,
			MaxWidth: col.MaxWidth,
			Align:    col.Align,
			Sortable: col.Sortable,
			Fixed:    !col.Resizable,
		})
	}
	if tl.actions >= 0 {
		w := tl.o.ActionsWidth
		if w <= 0 {
			w = 64
		}
		cols = append(cols, ui.TableColumn{
			Title:    tl.o.ActionsHeader,
			ID:       actionsColumnID,
			Width:    w,
			MinWidth: w,
			Align:    ui.End,
			Fixed:    true,
		})
	}
	return cols
}

const actionsColumnID = "\x00data.table.actions"

// buildCell builds body cell (row, col). It runs inside ui.Table's row build,
// so constructors called by the caller's builder attach to the cell.
func (tl *tableLayout) buildCell(c *ui.Context, row, col int) {
	st := tl.st
	if row < 0 || row >= len(st.order) {
		return
	}
	r := st.rowAt(&tl.o, row)
	t := theme.Colors(c)
	background := ui.Color{}
	ink := t.SurfaceForeground
	if r.Selected {
		background = t.SurfaceSecondary
	}
	cell := ui.Box(c).FillWidth().FillHeight().MinWidth(0).Background(background).TextColor(ink)
	cell.Children(func() {
		var e ui.Element
		switch {
		case col == tl.actions:
			if tl.o.RowActions != nil {
				e = tl.o.RowActions(c, r)
			}
		case col >= 0 && col < len(tl.columns):
			cb := tl.columns[col]
			build := cb.cell
			if build == nil {
				build = tl.o.Cell
			}
			if build != nil {
				e = build(c, TableCell{Row: r, Column: cb.column})
			}
		}
		if r.Disabled {
			style.Disabled(e, true)
		}
	})
}
