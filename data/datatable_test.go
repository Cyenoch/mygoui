package data

import (
	"slices"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestDataTableFiltersSortsPagesAndKeepsSelectionIdentity(t *testing.T) {
	rows := []Row{
		{ID: "c", Cells: []string{"Charlie", "pending"}},
		{ID: "a", Cells: []string{"Alpha", "paid"}},
		{ID: "b", Cells: []string{"Bravo", "paid"}},
		{ID: "d", Cells: []string{"Delta", "pending"}},
	}
	state := DataTableState{PageSize: 2}
	o := DataTableOptions{TableOptions: TableOptions{Rows: rows, Columns: []Column{{ID: "name", Title: "Name"}, {ID: "status", Title: "Status"}}}}
	state.query(o)
	state.Table.Selection.Add("a")
	state.Sort = ui.SortOrder{Column: "name"}
	state.query(o)
	if got := []RowID{state.matched[0].ID, state.matched[1].ID}; !slices.Equal(got, []RowID{"a", "b"}) {
		t.Fatalf("ascending page order: %v", got)
	}
	state.Page = 1
	state.query(o)
	if state.Page != 1 || state.matched[2].ID != "c" {
		t.Fatal("second page lost its real boundary")
	}
	state.Filter = "paid"
	state.query(o)
	if state.Page != 0 || len(state.matched) != 2 || !state.Table.Selection.Has("a") {
		t.Fatalf("filter did not reset page and preserve keyed selection: page=%d rows=%v", state.Page, state.matched)
	}
	state.Filter = "pending"
	state.query(o)
	if !state.Table.Selection.Has("a") {
		t.Fatal("filtering away a selected row deleted its identity")
	}
	o.Rows = []Row{rows[0], rows[2], rows[3]}
	o.Revision++
	state.query(o)
	if state.Table.Selection.Has("a") {
		t.Fatal("deleted row retained stale selection")
	}
	state.Page = 500
	state.Filter = "nothing"
	state.query(o)
	if state.Page != 0 || len(state.matched) != 0 {
		t.Fatal("empty results retained invalid page")
	}
}

func TestDataTableCachesComparisonsAndInvalidatesVersionedCells(t *testing.T) {
	comparisons := 0
	rows := []Row{{ID: "a", Cells: []string{"z"}}, {ID: "b", Cells: []string{"a"}}}
	columns := []Column{{ID: "name", Title: "Name", Compare: func(a, b Row) int {
		comparisons++
		if a.Cells[0] < b.Cells[0] {
			return -1
		}
		if a.Cells[0] > b.Cells[0] {
			return 1
		}
		return 0
	}}}
	state := DataTableState{Sort: ui.SortOrder{Column: "name"}}
	o := DataTableOptions{TableOptions: TableOptions{Rows: rows, Columns: columns}}
	state.query(o)
	count := comparisons
	state.query(o)
	state.Page = 2
	state.query(o)
	if count == 0 || comparisons != count {
		t.Fatal("unchanged input repeated sort")
	}
	rows[0].Cells[0] = "0"
	o.Revision++
	state.query(o)
	if comparisons <= count || state.matched[0].ID != "a" {
		t.Fatal("versioned cell update did not invalidate query")
	}
	state.Hidden = map[ColumnID]bool{"name": true}
	state.query(o)
	if len(state.visible) != 1 || state.Hidden["name"] {
		t.Fatal("programmatic visibility hid every readable column")
	}
}

func TestDataTableShowsActualPageAndVisibleColumnCells(t *testing.T) {
	state := DataTableState{PageSize: 1, Page: 1, Hidden: map[ColumnID]bool{"name": true}}
	o := DataTableOptions{Height: 100, TableOptions: TableOptions{
		Label: "Invoices", Rows: []Row{{ID: "a", Cells: []string{"Alpha", "First status"}}, {ID: "b", Cells: []string{"Bravo", "Second status"}}},
		Columns: []Column{{ID: "name", Title: "Name"}, {ID: "status", Title: "Status"}},
	}}
	tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c); DataTable(c, &state, o) }, 500, 320)
	if !tt.HasText("Second status") || tt.HasText("First status") || tt.HasText("Bravo") {
		t.Fatalf("page/visible cell mapping: %v", tt.Texts())
	}
	state.Page = 0
	tt.Frame()
	if !tt.HasText("First status") || tt.HasText("Second status") {
		t.Fatal("changing page did not replace actual rows")
	}
}
