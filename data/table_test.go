package data

import (
	"fmt"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestTableVirtualizesAndSelectionFollowsRows(t *testing.T) {
	rows := make([]Row, 5000)
	for i := range rows {
		rows[i] = Row{ID: RowID(fmt.Sprint(i)), Cells: []string{fmt.Sprintf("Record %04d", i)}}
	}
	state := TableState{}
	built := map[RowID]bool{}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		Table(c, &state, TableOptions{Rows: rows, Columns: []Column{{ID: "name", Title: "Name"}}, Cell: func(c *ui.Context, row Row, column Column) {
			built[row.ID] = true
			ui.Text(c, row.Cells[0])
		}}).Height(180)
	}, 360, 220)
	if len(built) >= 100 || len(built) == 0 {
		t.Fatalf("built %d of 5000 rows", len(built))
	}
	if err := tt.Click("Record 0001"); err != nil {
		t.Fatal(err)
	}
	if !state.Selection.Has("1") {
		t.Fatal("native row click did not select stable identity")
	}
	rows[0], rows[1] = rows[1], rows[0]
	tt.Frame()
	if !state.Selection.Has("1") || state.Selection.Has("0") {
		t.Fatal("reorder transferred selection to an index")
	}
	tt.Key(0, ui.KeyEnd)
	if !state.Selection.Has("4999") {
		t.Fatal("native End did not reach virtual final row")
	}
}
