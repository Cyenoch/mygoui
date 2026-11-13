package data_test

import (
	"testing"

	"github.com/Cyenoch/mygoui/data"
	"github.com/egoist/mygo/ui"
)

// TestTableArrowsPassADisabledRow pins the arrows over a disabled row: they
// carry the choice on to the next enabled row, in both modes.
func TestTableArrowsPassADisabledRow(t *testing.T) {
	for _, mode := range []data.TableSelectionMode{data.TableSelectionSingle, data.TableSelectionMultiple} {
		rows := fileRows()
		rows[1].disabled = true // Alpha
		st := &data.TableState{Selection: data.TableSelection{Mode: mode, Row: -1}}
		var reported []data.TableSelection
		o := fileTable(rows, st, nil)
		o.OnSelectionChange = func(s data.TableSelection) { reported = append(reported, s) }
		tt := ui.NewTester(tableView(&o), 520, 320)
		if err := tt.Click("Beta"); err != nil {
			t.Fatal(err)
		}
		tt.Key(0, ui.KeyDown)
		tt.Frame()
		if st.Selection.Row != 2 || mode == data.TableSelectionMultiple && (!st.Selection.HasKey("c") || st.Selection.Rows.Len() != 1) {
			t.Fatalf("mode %d: Down from Beta chose row %d, keys %v; want Gamma", mode, st.Selection.Row, keys(st.Selection))
		}
		if last := reported[len(reported)-1]; last.Row != 2 {
			t.Fatalf("mode %d: the move to Gamma was reported as row %d", mode, last.Row)
		}
		tt.Key(0, ui.KeyUp)
		tt.Frame()
		if st.Selection.Row != 0 {
			t.Fatalf("mode %d: Up from Gamma chose row %d, want Beta", mode, st.Selection.Row)
		}
	}
}

// TestTableHomeStopsOnTheNearestEnabledRow pins Home on a table whose first
// row is disabled: the choice goes to the first enabled row.
func TestTableHomeStopsOnTheNearestEnabledRow(t *testing.T) {
	rows := fileRows()
	rows[0].disabled = true // Beta
	st := &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionSingle, Row: -1}}
	o := fileTable(rows, st, nil)
	tt := ui.NewTester(tableView(&o), 520, 320)
	if err := tt.Click("Gamma"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyHome)
	tt.Frame()
	if st.Selection.Row != 1 {
		t.Fatalf("Home chose row %d, want Alpha", st.Selection.Row)
	}
}

// TestTableReportsBusyAndCanBeDisabled pins the table's own states: busy
// while loading, and a whole-table Disabled that takes no focus or click.
func TestTableReportsBusyAndCanBeDisabled(t *testing.T) {
	st := &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionSingle, Row: -1}}
	o := fileTable(fileRows(), st, nil)
	o.Loading = true
	tt := ui.NewTester(tableView(&o), 520, 320)
	if table, _ := tt.AccessByRole(ui.RoleTable, "Files"); !table.Busy {
		t.Fatal("a loading table is not busy")
	}
	o.Loading, o.Disabled = false, true
	tt.Frame()
	if table, _ := tt.AccessByRole(ui.RoleTable, "Files"); !table.Disabled {
		t.Fatal("a disabled table is not disabled")
	}
	tt.Key(0, ui.KeyTab)
	if table, _ := tt.AccessByRole(ui.RoleTable, "Files"); table.Focused {
		t.Fatal("a disabled table takes the focus")
	}
	_ = tt.Click("Alpha")
	if st.Selection.Row != -1 {
		t.Fatalf("a disabled table took a click: row %d", st.Selection.Row)
	}
}

// TestTableShiftArrowKeepsTheRangeAtADisabledEnd pins Shift+arrow toward a
// disabled last row: with no enabled row past it the range stays as it was,
// rather than collapsing to one row.
func TestTableShiftArrowKeepsTheRangeAtADisabledEnd(t *testing.T) {
	rows := fileRows()
	rows[2].disabled = true // Gamma
	st := &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionMultiple, Row: -1}}
	o := fileTable(rows, st, nil)
	tt := ui.NewTester(tableView(&o), 520, 320)
	if err := tt.Click("Beta"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Shift, ui.KeyDown)
	tt.Frame()
	tt.Key(ui.Shift, ui.KeyDown)
	tt.Frame()
	if st.Selection.Row != 1 || st.Selection.Rows.Len() != 2 || !st.Selection.HasKey("b") || !st.Selection.HasKey("a") {
		t.Fatalf("Shift+Down to the disabled end chose row %d, keys %v; want row 1, [b a]", st.Selection.Row, keys(st.Selection))
	}

	// Over a disabled middle row the range extends past it.
	rows = fileRows()
	rows[1].disabled = true // Alpha
	st = &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionMultiple, Row: -1}}
	o = fileTable(rows, st, nil)
	tt = ui.NewTester(tableView(&o), 520, 320)
	if err := tt.Click("Beta"); err != nil {
		t.Fatal(err)
	}
	tt.Key(ui.Shift, ui.KeyDown)
	tt.Frame()
	if st.Selection.Row != 2 || st.Selection.Rows.Len() != 2 || !st.Selection.HasKey("b") || !st.Selection.HasKey("c") {
		t.Fatalf("Shift+Down over a disabled row chose row %d, keys %v; want row 2, [b c]", st.Selection.Row, keys(st.Selection))
	}
}
