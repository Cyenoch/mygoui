package gallery

import (
	"github.com/egoist/mygo/ui"
)

func init() {
	RegisterComponent(ComponentPage{
		ID: "table", Title: "Table", Group: "Data", Order: 10,
		Description: "Virtualized rows with a fixed sortable header, resizable columns and row selection.",
		Examples: []ComponentExample{
			{Title: "Pagination", Description: "Eight deployments, four per page, with the summary and page controls in the footer.",
				Demo: func(c *ui.Context, s *State) {
					st := StateOf[tableState](s, "table")
					st.ensure()
					tablePaginatedExample(c, st)
				}},
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
			{Title: "Card variant",
				Demo: func(c *ui.Context, s *State) {
					st := StateOf[tableState](s, "table")
					st.ensure()
					tableSecondaryExample(c, st)
				}},
			{Title: "Column visibility", Description: "Hidden columns keep their width and sort order.",
				Demo: func(c *ui.Context, s *State) {
					st := StateOf[tableState](s, "table")
					st.ensure()
					tableVisibilityExample(c, st)
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
    Label: "Deployments", Variant: data.TablePrimary, Columns: columns,
    Rows: len(pageRows), State: &state,
    RowKey: func(row int) any { return pageRows[row].ID },
    Footer: func(c *ui.Context) ui.Element {
        return navigation.Pagination(c, navigation.PaginationOptions{
            Page: &page, Pages: pages, Size: navigation.PaginationSmall,
            Summary: func() { ui.Text(c, summary) },
        })
    },
}).Grow(1)`,
		Params: []ComponentParam{
			{"Label", "string", "required", "Names the table for assistive technology."},
			{"Variant", "data.TableVariant", "TablePrimary", "Framed card or a flush table under a header card."},
			{"Columns", "[]data.TableColumnOptions", "required", "Columns left to right; sortable and resizable."},
			{"Rows", "int", "0", "Caller row count; rows build only as they scroll in."},
			{"State", "*data.TableState", "nil", "Caller-owned selection, sort, widths and hidden set."},
			{"RowKey", "func(int) any", "nil", "Stable row identity for selection."},
			{"RowDisabled", "func(int) bool", "nil", "Marks rows that cannot be chosen or activated."},
			{"Compare", "func(a, b data.TableRow) int", "nil", "Orders rows for the active sort."},
			{"RowActions", "func(*ui.Context, data.TableRow) ui.Element", "nil", "Trailing fixed actions cell."},
			{"Empty", "func(*ui.Context) ui.Element", "nil", "Replaces the row area when there are no rows."},
			{"Footer", "func(*ui.Context) ui.Element", "nil", "Row inside the frame for a summary and pagination."},
		},
		Tokens: []string{"Surface", "SurfaceForeground", "SurfaceSecondary", "SurfaceSecondaryForeground", "SeparatorTertiary", "Border"},
	})
}
