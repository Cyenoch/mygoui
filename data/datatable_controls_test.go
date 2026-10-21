package data_test

import (
	"testing"

	"github.com/Cyenoch/mygoui/data"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestDataTableNativeSearchSortAndColumnMenuChangeActualRows(t *testing.T) {
	state := data.DataTableState{PageSize: 1}
	options := data.DataTableOptions{Height: 120, TableOptions: data.TableOptions{
		Rows:    []data.Row{{ID: "b", Cells: []string{"Bravo", "Pending"}}, {ID: "a", Cells: []string{"Alpha", "Paid"}}},
		Columns: []data.Column{{ID: "name", Title: "Name"}, {ID: "status", Title: "Status"}},
	}}
	tt := ui.NewTester(func(c *ui.Context) { theme.Apply(c); data.DataTable(c, &state, options) }, 500, 340)
	if !tt.HasText("Bravo") || tt.HasText("Alpha") {
		t.Fatal("first page contains other page rows")
	}
	if err := tt.Click("Name"); err != nil {
		t.Fatal(err)
	}
	if !tt.HasText("Alpha") || tt.HasText("Bravo") {
		t.Fatal("native sortable header did not reorder actual page")
	}
	if err := tt.Click("Columns"); err != nil {
		t.Fatal(err)
	}
	if err := tt.ChooseMenuItem("Name"); err != nil {
		t.Fatal(err)
	}
	if !state.Hidden["name"] || tt.HasText("Alpha") || !tt.HasText("Paid") {
		t.Fatal("column menu changed wrong cell projection")
	}
	if err := tt.Click("Filter rows"); err != nil {
		t.Fatal(err)
	}
	tt.Type("pending")
	if !tt.HasText("Pending") || tt.HasText("Paid") {
		t.Fatal("native input did not filter actual rows")
	}
	tt.Type(" absent")
	if !tt.HasText("No matching rows") || tt.HasText("Pending") {
		t.Fatal("empty filter kept stale data visible")
	}
}
