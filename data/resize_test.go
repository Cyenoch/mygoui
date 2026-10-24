package data_test

import (
	"testing"

	"github.com/Cyenoch/mygoui/data"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestTableColumnResize proves that Resizable is wired to the toolkit's header
// drag. The column's right edge is ten DIPs past the end-aligned header text
// (the toolkit's cell padding), where its resize handle sits.
func TestTableColumnResize(t *testing.T) {
	rows := fileRows()
	st := &data.TableState{}
	o := data.TableOptions{
		Label: "Files",
		Columns: []data.TableColumnOptions{
			{
				ID: "name", Header: "Name", Resizable: true,
				Cell: func(c *ui.Context, cell data.TableCell) ui.Element {
					return ui.Text(c, rows[cell.Row.Index].name)
				},
			},
			{
				ID: "size", Header: "Size", Width: 90, MinWidth: 60, Resizable: true, Align: ui.End,
				Cell: func(c *ui.Context, cell data.TableCell) ui.Element {
					return ui.Text(c, "1234")
				},
			},
		},
		Rows: len(rows), State: st,
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		data.Table(c, o).Grow(1)
	}, 520, 320)
	rect, ok := tt.Find("Size")
	if !ok {
		t.Fatal("Size header missing")
	}
	edge := rect.X + rect.W + theme.Space(2.5) - 2
	y := rect.Y + rect.H/2
	tt.Press(edge, y)
	tt.Move(edge+30, y)
	tt.Release(edge+30, y)
	if w := st.List.Columns.Widths["size"]; w != 120 {
		t.Fatalf("resized width = %v, want 120", w)
	}
}
