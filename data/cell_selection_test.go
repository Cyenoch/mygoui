package data_test

import (
	"slices"
	"testing"

	"github.com/egoist/mygo/ui"

	"github.com/Cyenoch/mygoui/data"
	"github.com/Cyenoch/mygoui/theme"
)

func TestTableCellSelectionAndRowHeaders(t *testing.T) {
	state := data.TableState{List: ui.ListState{CellSelect: true}}
	rows := []data.Row{
		{ID: "ada", Cells: []string{"1", "Ada", "Ready"}},
		{ID: "grace", Cells: []string{"2", "Grace", "Review"}},
		{ID: "alan", Cells: []string{"3", "Alan", "Draft"}},
	}
	columns := []data.Column{
		{ID: "number", Title: "#", Width: 48, Fixed: true, RowHeader: true},
		{ID: "name", Title: "Name", Width: 140},
		{ID: "status", Title: "Status", Width: 140},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		data.Table(c, &state, data.TableOptions{
			Rows: rows, Columns: columns,
			RowLabel: func(row data.Row) string { return "Person " + row.Cells[1] },
		}).Height(200)
	}, 360, 240)
	if err := tt.Click("Ada"); err != nil {
		t.Fatal(err)
	}
	if err := tt.ClickWith(ui.Shift, "Review"); err != nil {
		t.Fatal(err)
	}
	if state.Selection.Len() != 2 || !state.Selection.Has("ada") || !state.Selection.Has("grace") {
		t.Fatalf("cell block lost row identities: %v", slices.Collect(state.Selection.All()))
	}
	if !state.List.CellSelected(0, 1) || !state.List.CellSelected(1, 2) || state.List.CellSelected(2, 2) {
		t.Fatal("Shift-click did not choose the actual rectangular cell block")
	}
	if err := tt.Click("3"); err != nil {
		t.Fatal(err)
	}
	if !state.List.Cells.All || state.Selection.Len() != 1 || !state.Selection.Has("alan") {
		t.Fatalf("row header did not choose the whole row: cells=%+v selection=%v", state.List.Cells, slices.Collect(state.Selection.All()))
	}
	if err := tt.Click("#"); err != nil {
		t.Fatal(err)
	}
	if state.Selection.Len() != len(rows) || !state.List.Cells.All {
		t.Fatal("row-header heading did not choose every row")
	}
}
