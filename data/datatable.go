package data

import (
	"maps"
	"slices"
	"strings"

	"github.com/Cyenoch/mygoui/form"
	locales "github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type DataTableState struct {
	Table          TableState
	Filter         string
	Sort           ui.SortOrder
	Page, PageSize int
	Hidden         map[ColumnID]bool

	revision       uint64
	rows           *Row
	rowCount       int
	columns        *Column
	columnCount    int
	ready          bool
	search         []string
	matched        []Row
	lastFilter     string
	lastSort       ui.SortOrder
	lastPageSize   int
	lastHidden     map[ColumnID]bool
	visible        []Column
	visibleIndices []int
	pages          int
}

type DataTableOptions struct {
	TableOptions
	// Revision must change when existing row cells, columns or comparator
	// behavior change in place. Replacing an input slice also invalidates it.
	Revision          uint64
	FilterPlaceholder string
	Height            float32
}

// DataTable adds cached local filtering and sorting, column visibility and real
// pagination to Table. The maintained native table owns virtualization, stable
// selection, column resize/reorder and keyboard navigation. Only the current
// page's rows are supplied to it. Column visibility uses a native system menu.
func DataTable(c *ui.Context, state *DataTableState, o DataTableOptions) ui.Element {
	state.query(o)
	height := o.Height
	if height <= 0 {
		height = 260
	}
	placeholder := o.FilterPlaceholder
	if placeholder == "" {
		placeholder = locales.ComponentFilterRowsPlaceholder()
	}
	t := theme.Colors(c)
	return ui.Column(c).MinWidth(0).Gap(theme.Space12).Disabled(o.Disabled).Children(func() {
		ui.Row(c).MinWidth(0).FillWidth().Wrap().Gap(theme.Space8).Children(func() {
			form.Input(c.Key("filter"), &state.Filter, form.InputOptions{
				ControlOptions: form.ControlOptions{Label: locales.ComponentFilterRows(), Disabled: o.Disabled},
				Placeholder:    placeholder,
			}).Grow(1).Basis(0).MinWidth(theme.Scale(c, 100))
			ui.MenuButton(c.Key("columns"), locales.ComponentColumns(), func(m *ui.Menu) {
				visible := len(state.visible)
				for _, column := range o.Columns {
					hidden := state.Hidden[column.ID]
					if m.Item(column.Title).Checked(!hidden).Disabled(!hidden && visible <= 1).Chosen() {
						if state.Hidden == nil {
							state.Hidden = make(map[ColumnID]bool)
						}
						state.Hidden[column.ID] = !hidden
					}
				}
			}).Disabled(o.Disabled)
		})
		start := min(state.Page*state.PageSize, len(state.matched))
		end := min(start+state.PageSize, len(state.matched))
		table := o.TableOptions
		table.Rows, table.Columns = state.matched[start:end], state.visible
		if table.RowLabel == nil {
			table.RowLabel = func(row Row) string {
				for _, original := range state.visibleIndices {
					if original < len(row.Cells) && row.Cells[original] != "" {
						return row.Cells[original]
					}
				}
				return string(row.ID)
			}
		}
		table.Cell = func(c *ui.Context, row Row, column Column) {
			if o.Cell != nil {
				o.Cell(c, row, column)
				return
			}
			for i, visible := range state.visible {
				if visible.ID == column.ID {
					original := state.visibleIndices[i]
					if original < len(row.Cells) {
						ui.Text(c, row.Cells[original]).FontSize(theme.Scale(c, theme.Body().Size))
					}
					return
				}
			}
		}
		state.Table.List.Sort = &state.Sort
		Table(c.Key("results"), &state.Table, table).Height(height)
		ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(theme.Space8).Children(func() {
			label := locales.ComponentRowsSelected(len(state.matched), state.Table.Selection.Len())
			if len(state.matched) == 0 {
				label = locales.ComponentNoMatchingRows()
			}
			ui.Text(c, label).Grow(1).FontSize(theme.Scale(c, theme.Small().Size)).TextColor(t.MutedForeground)
			navigation.Pagination(c.Key("pages"), navigation.PaginationOptions{
				Page: state.Page + 1, Pages: state.pages, Disabled: o.Disabled, Simple: true,
				OnChange: func(page int) { state.Page = page - 1 },
			})
		})
	})
}

// query is the data derivation seam, not a table engine: rows are filtered and
// sorted only when their version, search or native header sort changes. Row
// cells, selection and column geometry remain in their existing native owners.
func (s *DataTableState) query(o DataTableOptions) {
	var rows *Row
	var columns *Column
	if len(o.Rows) > 0 {
		rows = &o.Rows[0]
	}
	if len(o.Columns) > 0 {
		columns = &o.Columns[0]
	}
	dataChanged := !s.ready || s.revision != o.Revision || s.rows != rows || s.rowCount != len(o.Rows) || s.columns != columns || s.columnCount != len(o.Columns)
	if s.PageSize <= 0 {
		s.PageSize = 10
	}
	if dataChanged {
		s.search = slices.Grow(s.search[:0], len(o.Rows))
		present := make(map[RowID]struct{}, len(o.Rows))
		for _, row := range o.Rows {
			s.search = append(s.search, strings.ToLower(strings.Join(row.Cells, "\x00")))
			present[row.ID] = struct{}{}
		}
		for id := range s.Table.Selection.All() {
			if _, exists := present[id]; !exists {
				s.Table.Selection.Remove(id)
			}
		}
	}
	if dataChanged || s.lastFilter != s.Filter || s.lastSort != s.Sort {
		if s.ready && s.lastFilter != s.Filter {
			s.Page = 0
		}
		s.matched = s.matched[:0]
		filter := strings.ToLower(strings.TrimSpace(s.Filter))
		for i, row := range o.Rows {
			if filter == "" || strings.Contains(s.search[i], filter) {
				s.matched = append(s.matched, row)
			}
		}
		for i, column := range o.Columns {
			if string(column.ID) != s.Sort.Column {
				continue
			}
			compare := column.Compare
			if compare == nil {
				compare = func(a, b Row) int {
					var av, bv string
					if i < len(a.Cells) {
						av = a.Cells[i]
					}
					if i < len(b.Cells) {
						bv = b.Cells[i]
					}
					return strings.Compare(av, bv)
				}
			}
			slices.SortStableFunc(s.matched, func(a, b Row) int {
				if s.Sort.Descending {
					return compare(b, a)
				}
				return compare(a, b)
			})
			break
		}
	}
	if dataChanged || !maps.Equal(s.lastHidden, s.Hidden) {
		s.visible, s.visibleIndices = s.visible[:0], s.visibleIndices[:0]
		for i, column := range o.Columns {
			if !s.Hidden[column.ID] {
				s.visible = append(s.visible, column)
				s.visibleIndices = append(s.visibleIndices, i)
			}
		}
		// Even programmatic visibility changes retain one readable column.
		if len(s.visible) == 0 && len(o.Columns) > 0 {
			s.visible, s.visibleIndices = append(s.visible, o.Columns[0]), append(s.visibleIndices, 0)
			delete(s.Hidden, o.Columns[0].ID)
		}
		s.lastHidden = maps.Clone(s.Hidden)
	}
	if s.ready && s.lastPageSize != s.PageSize {
		s.Page = 0
	}
	s.pages = 0
	if len(s.matched) > 0 {
		s.pages = 1 + (len(s.matched)-1)/s.PageSize
	}
	s.Page = min(max(0, s.Page), max(0, s.pages-1))
	s.revision, s.rows, s.rowCount, s.columns, s.columnCount = o.Revision, rows, len(o.Rows), columns, len(o.Columns)
	s.lastFilter, s.lastSort, s.lastPageSize, s.ready = s.Filter, s.Sort, s.PageSize, true
}
