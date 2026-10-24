package gallery

import (
	"github.com/egoist/mygo/ui"
)

func init() {
	RegisterComponent(ComponentPage{
		ID: "table", Title: "Table", Group: "Data", Order: 10,
		Description: "Virtualized rows with a fixed sortable header, resizable columns and row selection.",
		Examples: []ComponentExample{
			{Title: "Selection", Description: "Multiple row selection with an action on every row.",
				Demo: func(c *ui.Context, s *State) {
					st := StateOf[tableState](s, "table")
					st.ensure()
					tableSelectionExample(c, st)
				}},
			{Title: "Sorting", Description: "Click a sortable header to reorder the rows.",
				Demo: func(c *ui.Context, s *State) {
					st := StateOf[tableState](s, "table")
					st.ensure()
					tableSortingExample(c, st)
				}},
			{Title: "Empty state", Description: "The empty builder replaces the row area when there are no rows.",
				Demo: func(c *ui.Context, s *State) {
					st := StateOf[tableState](s, "table")
					st.ensure()
					tableEmptyExample(c, st)
				}},
			{Title: "Disabled row", Description: "One row cannot be focused or selected.",
				Demo: func(c *ui.Context, s *State) {
					st := StateOf[tableState](s, "table")
					st.ensure()
					tableDisabledExample(c, st)
				}},
		},
		Code: `data.Table(c, data.TableOptions{
    Label: "Deployments", Columns: columns, Rows: len(rows),
    State: &state, RowKey: func(row int) any { return rows[row].ID },
    Compare: compareRows,
}).Grow(1)`,
		Params: []ComponentParam{
			{"Label", "string", "required", "Names the table for assistive technology."},
			{"Columns", "[]data.TableColumnOptions", "required", "Columns left to right; sortable and resizable."},
			{"Rows", "int", "0", "Caller row count; rows build only as they scroll in."},
			{"State", "*data.TableState", "nil", "Caller-owned selection, sort and widths."},
			{"RowKey", "func(int) any", "nil", "Stable row identity for selection."},
			{"Compare", "func(a, b data.TableRow) int", "nil", "Orders rows for the active sort."},
		},
		Tokens: []string{"Surface", "SurfaceForeground", "SurfaceSecondary", "SurfaceSecondaryForeground", "Border"},
	})
}
