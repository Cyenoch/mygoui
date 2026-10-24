package data_test

import (
	"cmp"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/data"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

type fileRow struct {
	id       string
	name     string
	size     int
	disabled bool
}

func fileRows() []fileRow {
	return []fileRow{
		{id: "b", name: "Beta", size: 30},
		{id: "a", name: "Alpha", size: 10},
		{id: "c", name: "Gamma", size: 20},
	}
}

// fileTable builds the options shared by the tests: name and size columns,
// stable IDs, and a comparator that follows the state's sort.
func fileTable(rows []fileRow, st *data.TableState, sorts *[]data.TableSort) data.TableOptions {
	return data.TableOptions{
		Label: "Files",
		Columns: []data.TableColumnOptions{
			{
				ID: "name", Header: "Name", Sortable: true, Resizable: true,
				Cell: func(c *ui.Context, cell data.TableCell) ui.Element {
					return ui.Text(c, rows[cell.Row.Index].name)
				},
			},
			{
				ID: "size", Header: "Size", Sortable: true, Align: ui.End, Width: 90,
				Cell: func(c *ui.Context, cell data.TableCell) ui.Element {
					return ui.Text(c, fmt.Sprint(rows[cell.Row.Index].size))
				},
			},
		},
		Rows:   len(rows),
		State:  st,
		RowKey: func(row int) any { return rows[row].id },
		RowDisabled: func(row int) bool {
			return rows[row].disabled
		},
		Compare: func(a, b data.TableRow) int {
			switch st.Sort.Column {
			case "name":
				return strings.Compare(rows[a.Index].name, rows[b.Index].name)
			case "size":
				return cmp.Compare(rows[a.Index].size, rows[b.Index].size)
			}
			return 0
		},
		OnSort: func(sort data.TableSort) {
			if sorts != nil {
				*sorts = append(*sorts, sort)
			}
		},
	}
}

func tableView(o *data.TableOptions) func(c *ui.Context) {
	return func(c *ui.Context) {
		theme.Apply(c)
		data.Table(c, *o).Grow(1)
	}
}

func TestTableHeaderAndBody(t *testing.T) {
	rows := fileRows()
	st := &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionMultiple, Row: -1}}
	o := fileTable(rows, st, nil)
	tt := ui.NewTester(tableView(&o), 520, 320)
	if !tt.HasText("Name") || !tt.HasText("Size") {
		t.Fatalf("headers missing: %v", tt.Texts())
	}
	for _, r := range rows {
		if !tt.HasText(r.name) {
			t.Fatalf("row %q missing: %v", r.name, tt.Texts())
		}
	}
}

func TestTableSortOrdersRowsAndNotifies(t *testing.T) {
	rows := fileRows()
	st := &data.TableState{}
	var sorts []data.TableSort
	o := fileTable(rows, st, &sorts)
	tt := ui.NewTester(tableView(&o), 520, 320)

	alpha, _ := tt.Find("Alpha")
	beta, _ := tt.Find("Beta")
	if alpha.Y < beta.Y {
		t.Fatalf("unsorted rows already in name order")
	}
	if err := tt.Click("Name"); err != nil {
		t.Fatal(err)
	}
	if len(sorts) != 1 || sorts[0] != (data.TableSort{Column: "name"}) {
		t.Fatalf("first sort callback: %+v", sorts)
	}
	alpha, _ = tt.Find("Alpha")
	beta, _ = tt.Find("Beta")
	gamma, _ := tt.Find("Gamma")
	if !(alpha.Y < beta.Y && beta.Y < gamma.Y) {
		t.Fatalf("ascending order: alpha=%v beta=%v gamma=%v", alpha.Y, beta.Y, gamma.Y)
	}

	if err := tt.Click("Name"); err != nil {
		t.Fatal(err)
	}
	if len(sorts) != 2 || sorts[1] != (data.TableSort{Column: "name", Descending: true}) {
		t.Fatalf("second sort callback: %+v", sorts)
	}
	alpha, _ = tt.Find("Alpha")
	gamma, _ = tt.Find("Gamma")
	if gamma.Y > alpha.Y {
		t.Fatalf("descending order: alpha=%v gamma=%v", alpha.Y, gamma.Y)
	}

	if err := tt.Click("Size"); err != nil {
		t.Fatal(err)
	}
	if len(sorts) != 3 || sorts[2] != (data.TableSort{Column: "size"}) {
		t.Fatalf("third sort callback: %+v", sorts)
	}
	if st.Sort != (data.TableSort{Column: "size"}) {
		t.Fatalf("state sort: %+v", st.Sort)
	}
	alpha, _ = tt.Find("Alpha") // 10
	gamma, _ = tt.Find("Gamma") // 20
	beta, _ = tt.Find("Beta")   // 30
	if !(alpha.Y < gamma.Y && gamma.Y < beta.Y) {
		t.Fatalf("size order: alpha=%v gamma=%v beta=%v", alpha.Y, gamma.Y, beta.Y)
	}
}

func TestTableSelectionSingleAndMultiple(t *testing.T) {
	rows := fileRows()

	single := &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionSingle, Row: -1}}
	o := fileTable(rows, single, nil)
	tt := ui.NewTester(tableView(&o), 520, 320)
	if err := tt.Click("Gamma"); err != nil {
		t.Fatal(err)
	}
	if single.Selection.Row != 2 {
		t.Fatalf("single selection row = %d, want 2", single.Selection.Row)
	}

	multiple := &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionMultiple, Row: -1}}
	var calls int
	o = fileTable(rows, multiple, nil)
	o.OnSelectionChange = func(data.TableSelection) { calls++ }
	tt = ui.NewTester(tableView(&o), 520, 320)
	if err := tt.Click("Beta"); err != nil {
		t.Fatal(err)
	}
	if !multiple.Selection.HasKey("b") || multiple.Selection.Row != 0 {
		t.Fatalf("beta chosen: row=%d keys=%v", multiple.Selection.Row, keys(multiple.Selection))
	}
	if err := tt.ClickWith(ui.Cmd, "Alpha"); err != nil {
		t.Fatal(err)
	}
	if !multiple.Selection.HasKey("b") || !multiple.Selection.HasKey("a") {
		t.Fatalf("cmd-click added alpha: %v", keys(multiple.Selection))
	}
	if calls != 2 {
		t.Fatalf("selection callbacks = %d, want 2", calls)
	}

	// A sort keeps a keyed multiple selection with its rows, remaps a single
	// selection to the same row, and is not itself a selection change.
	if err := tt.Click("Name"); err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("sort reported a selection change: %d callbacks", calls)
	}
	if !multiple.Selection.HasKey("b") || !multiple.Selection.HasKey("a") {
		t.Fatalf("sort lost the multiple selection: %v", keys(multiple.Selection))
	}
	if err := tt.Click("Beta"); err != nil {
		t.Fatal(err)
	}
	if multiple.Selection.Row != 1 || calls != 3 {
		t.Fatalf("beta after sort: row = %d, callbacks = %d, want 1 and 3", multiple.Selection.Row, calls)
	}

	single = &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionSingle, Row: -1}}
	o = fileTable(rows, single, nil)
	tt = ui.NewTester(tableView(&o), 520, 320)
	if err := tt.Click("Beta"); err != nil {
		t.Fatal(err)
	}
	if err := tt.Click("Name"); err != nil {
		t.Fatal(err)
	}
	if single.Selection.Row != 1 {
		t.Fatalf("single selection did not follow its row: %d", single.Selection.Row)
	}
}

func keys(s data.TableSelection) []any {
	var out []any
	for k := range s.Rows.All() {
		out = append(out, k)
	}
	return out
}

func TestTableDisabledRowCannotBeChosen(t *testing.T) {
	rows := fileRows()
	rows[0].disabled = true // Beta

	st := &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionMultiple, Row: -1}}
	o := fileTable(rows, st, nil)
	tt := ui.NewTester(tableView(&o), 520, 320)
	if err := tt.Click("Beta"); err != nil {
		t.Fatal(err)
	}
	if st.Selection.HasKey("b") || st.Selection.Row != -1 {
		t.Fatalf("disabled row chosen: row=%d keys=%v", st.Selection.Row, keys(st.Selection))
	}

	single := &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionSingle, Row: -1}}
	o = fileTable(rows, single, nil)
	tt = ui.NewTester(tableView(&o), 520, 320)
	if err := tt.Click("Beta"); err != nil {
		t.Fatal(err)
	}
	if single.Selection.Row != -1 {
		t.Fatalf("disabled row chosen singly: %d", single.Selection.Row)
	}
	if err := tt.Click("Alpha"); err != nil {
		t.Fatal(err)
	}
	if single.Selection.Row != 1 {
		t.Fatalf("enabled row after disabled: %d", single.Selection.Row)
	}
}

func TestTableWheelScrollKeepsHeader(t *testing.T) {
	const n = 200
	rows := make([]fileRow, n)
	for i := range rows {
		rows[i] = fileRow{id: fmt.Sprintf("r%d", i), name: fmt.Sprintf("Row %03d", i), size: i}
	}
	st := &data.TableState{}
	o := data.TableOptions{
		Label: "Rows",
		Columns: []data.TableColumnOptions{{
			ID: "name", Header: "Header", Sortable: true,
			Cell: func(c *ui.Context, cell data.TableCell) ui.Element {
				return ui.Text(c, rows[cell.Row.Index].name)
			},
		}},
		Rows: n, State: st,
		RowKey: func(row int) any { return rows[row].id },
	}
	tt := ui.NewTester(tableView(&o), 520, 300)
	head, ok := tt.Find("Header")
	if !ok {
		t.Fatal("header missing")
	}
	tt.Move(250, 200)
	tt.Wheel(250, 200, 0, 300)
	tt.Advance(time.Second)
	after, ok := tt.Find("Header")
	if !ok {
		t.Fatal("header gone after scrolling")
	}
	if after.Y != head.Y {
		t.Fatalf("header moved from %v to %v", head.Y, after.Y)
	}
	first, last := st.List.Visible()
	if first == 0 {
		t.Fatalf("rows did not scroll: visible %d..%d", first, last)
	}
	if !tt.HasText(fmt.Sprintf("Row %03d", first)) {
		t.Fatalf("row %d not built: %v", first, tt.Texts())
	}
}

func TestTableEmptyStateAndHiddenColumn(t *testing.T) {
	st := &data.TableState{}
	custom := false
	o := data.TableOptions{
		Label: "Projects",
		Columns: []data.TableColumnOptions{
			{ID: "name", Header: "Name"},
			{ID: "size", Header: "Size", Hideable: true},
		},
		Rows: 0, State: st,
		Empty: func(c *ui.Context) ui.Element {
			custom = true
			return ui.Text(c, "No projects yet")
		},
	}
	tt := ui.NewTester(tableView(&o), 520, 320)
	if !custom || !tt.HasText("No projects yet") {
		t.Fatalf("empty slot missing: %v", tt.Texts())
	}
	if !tt.HasText("Name") || !tt.HasText("Size") {
		t.Fatalf("empty header: %v", tt.Texts())
	}

	st.HideColumn("size")
	tt.Frame()
	if tt.HasText("Size") || !tt.HasText("Name") {
		t.Fatalf("hidden column: %v", tt.Texts())
	}
	st.HideColumn("name")
	tt.Frame()
	if tt.HasText("Name") || !tt.HasText("No projects yet") {
		t.Fatalf("all columns hidden: %v", tt.Texts())
	}
	if !st.ColumnHidden("name") {
		t.Fatal("hidden state not recorded")
	}
	st.ToggleColumn("name")
	tt.Frame()
	if !tt.HasText("Name") {
		t.Fatalf("shown column: %v", tt.Texts())
	}
}

func TestTableRowActions(t *testing.T) {
	rows := fileRows()
	st := &data.TableState{}
	var activated string
	o := fileTable(rows, st, nil)
	o.RowActions = func(c *ui.Context, row data.TableRow) ui.Element {
		return ui.Text(c, "Edit "+rows[row.Index].name)
	}
	o.OnRowActivate = func(row data.TableRow) { activated = rows[row.Index].name }
	tt := ui.NewTester(tableView(&o), 640, 320)
	for _, r := range rows {
		if !tt.HasText("Edit " + r.name) {
			t.Fatalf("actions cell for %q missing: %v", r.name, tt.Texts())
		}
	}
	// Activation requires a selection in the native table.
	o.State.Selection.Mode = data.TableSelectionSingle
	tt.Frame()
	if err := tt.Click("Gamma"); err != nil {
		t.Fatal(err)
	}
	tt.Key(0, ui.KeyEnter)
	if activated != "Gamma" {
		t.Fatalf("activated %q", activated)
	}
}
