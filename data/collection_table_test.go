package data_test

import (
	"fmt"
	"image"
	"testing"

	"github.com/Cyenoch/mygoui/data"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// tablePixelNear reports whether the pixel is within tol of c.
func tablePixelNear(img *image.RGBA, x, y int, c ui.Color, tol uint8) bool {
	got := img.RGBAAt(x, y)
	near := func(got, want uint8) bool {
		d := int(got) - int(want)
		return d >= -int(tol) && d <= int(tol)
	}
	return near(got.R, c.R) && near(got.G, c.G) && near(got.B, c.B)
}

func wantTablePixel(t *testing.T, img *image.RGBA, x, y int, c ui.Color, what string) {
	t.Helper()
	if !tablePixelNear(img, x, y, c, 4) {
		got := img.RGBAAt(x, y)
		t.Fatalf("%s at (%d, %d) = %d,%d,%d, want %d,%d,%d",
			what, x, y, got.R, got.G, got.B, c.R, c.G, c.B)
	}
}

// TestTableSelectedRowPaintsOneBand checks a chosen row paints its tint across
// the whole row, including the toolkit's cell padding. Painting only the cell
// content boxes leaves seams of the native accent between them.
func TestTableSelectedRowPaintsOneBand(t *testing.T) {
	rows := fileRows()
	st := &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionSingle, Row: -1}}
	o := fileTable(rows, st, nil)
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(20).
			Background(colors.Background).Children(func() {
			data.Table(c, o).Grow(1)
		})
	}, 520, 320)
	tt.Frame()
	if err := tt.Click("Alpha"); err != nil {
		t.Fatal(err)
	}
	tt.Move(2, 2)
	tt.Frame()

	row, ok := tt.Find("Alpha")
	if !ok {
		t.Fatalf("chosen row missing: %v", tt.Texts())
	}
	head, ok := tt.Find("Name")
	if !ok {
		t.Fatalf("header missing: %v", tt.Texts())
	}
	img := tt.Image()
	y := int(row.Y + row.H/2)
	// The header's box is the first column: its left edge is the table's.
	cellRight := int(head.X + head.W)
	tableRight := cellRight + 90 // the Size column
	for _, x := range []int{
		int(head.X) + 2, // the toolkit's own left cell padding
		cellRight - 2,   // the first column's end padding
		cellRight + 5,   // the next column's start padding
		tableRight - 2,  // the table's right end padding
	} {
		wantTablePixel(t, img, x, y, ui.RGBA(0, 0, 0, 0.04).Over(colors.Background), "chosen row fill")
		if tablePixelNear(img, x, y, colors.Primary, 6) {
			t.Fatalf("chosen row shows the native primary at (%d, %d)", x, y)
		}
	}
	// The row below is not chosen: its fill is the table surface.
	other, _ := tt.Find("Gamma")
	wantTablePixel(t, img, int(head.X)+2, int(other.Y+other.H/2), colors.Background, "unchosen row fill")
}

// TestTableVariants distinguishes the plain background from the card body.
func TestTableVariants(t *testing.T) {
	rows := fileRows()
	for _, dark := range []bool{false, true} {
		for _, variant := range []data.TableVariant{data.TablePrimary, data.TableSecondary} {
			st := &data.TableState{}
			o := fileTable(rows, st, nil)
			o.Variant = variant
			var colors theme.Palette
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, dark)
				colors = theme.Colors(c)
				data.Table(c, o).Grow(1)
			}, 400, 300)
			tt.Frame()
			row, ok := tt.Find("Beta")
			if !ok {
				t.Fatalf("row missing: %v", tt.Texts())
			}
			want := colors.Background
			if variant == data.TableSecondary {
				want = colors.Card
			}
			// The row's left padding is outside every cell's content box.
			wantTablePixel(t, tt.Image(), 6, int(row.Y+row.H/2), want,
				fmt.Sprintf("variant %d, dark=%v, row fill", variant, dark))
		}
	}
}

// TestTableCellAlignment checks column alignment and the vertical centering of
// cell content in a row made taller by another cell.
func TestTableCellAlignment(t *testing.T) {
	st := &data.TableState{}
	o := data.TableOptions{
		Label: "Mixed",
		Columns: []data.TableColumnOptions{
			{ID: "tall", Header: "Tall", Sortable: true, Cell: func(c *ui.Context, cell data.TableCell) ui.Element {
				return ui.Box(c).Width(20).Height(40).Background(ui.Hex("#dddddd")).Label("TALL")
			}},
			{ID: "end", Header: "End", Width: 90, Align: ui.End, Sortable: true,
				Cell: func(c *ui.Context, cell data.TableCell) ui.Element { return ui.Text(c, "END") }},
			{ID: "mid", Header: "Mid", Width: 90, Align: ui.Center, Sortable: true,
				Cell: func(c *ui.Context, cell data.TableCell) ui.Element { return ui.Text(c, "MID") }},
		},
		Rows: 1, State: st,
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		data.Table(c, o).Grow(1)
	}, 500, 220)
	tt.Frame()
	tall, ok := tt.Find("TALL")
	if !ok {
		t.Fatalf("tall cell missing: %v", tt.Texts())
	}
	end, ok := tt.Find("END")
	if !ok {
		t.Fatalf("end cell missing: %v", tt.Texts())
	}
	mid, ok := tt.Find("MID")
	if !ok {
		t.Fatalf("center cell missing: %v", tt.Texts())
	}
	endHead, ok := tt.Find("End")
	if !ok {
		t.Fatalf("end header missing: %v", tt.Texts())
	}
	midHead, ok := tt.Find("Mid")
	if !ok {
		t.Fatalf("center header missing: %v", tt.Texts())
	}
	// The end-aligned text stops at the toolkit's 10 DIP cell padding.
	wantEnd := endHead.X + endHead.W - 10
	if got := end.X + end.W; got < wantEnd-3 || got > wantEnd+3 {
		t.Fatalf("end-aligned text right edge = %v, want %v", got, wantEnd)
	}
	if got, want := mid.X+mid.W/2, midHead.X+midHead.W/2; got < want-3 || got > want+3 {
		t.Fatalf("centered text center = %v, want %v", got, want)
	}
	for _, r := range []ui.Rect{end, mid} {
		if got, want := r.Y+r.H/2, tall.Y+tall.H/2; got < want-1 || got > want+1 {
			t.Fatalf("cell content center = %v, want %v", got, want)
		}
	}
}

// The default Coss table has a transparent header and a border below it.
func TestTableHeaderChrome(t *testing.T) {
	rows := fileRows()
	var colors theme.Palette
	o := fileTable(rows, &data.TableState{}, nil)
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		data.Table(c, o).Grow(1)
	}, 400, 240)
	tt.Frame()
	wantTablePixel(t, tt.Image(), 200, 16, colors.Background, "header face")
	wantTablePixel(t, tt.Image(), 200, 32, colors.Border.Over(colors.Background), "header divider")
}

func TestEmptyCardWithoutColumnsStartsAtTop(t *testing.T) {
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, true)
		colors = theme.Colors(c)
		data.Table(c, data.TableOptions{
			Variant: data.TableSecondary,
			Empty:   func(c *ui.Context) ui.Element { return ui.Text(c, "No rows") },
		}).Grow(1)
	}, 400, 240)
	wantTablePixel(t, tt.Image(), 200, 16, colors.Card, "empty card body")
}

func TestTableKeepsNativeFocusRing(t *testing.T) {
	rows := fileRows()
	st := &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionSingle, Row: -1}}
	o := fileTable(rows, st, nil)
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).Fill().Padding(20).Background(colors.Background).Children(func() {
			data.Table(c, o).Grow(1)
		})
	}, 420, 260)
	table, ok := tt.Find("Files")
	if !ok {
		t.Fatal("table missing")
	}
	x, y := int(table.X+table.W/2), int(table.Y)-2
	wantTablePixel(t, tt.Image(), x, y, colors.Background, "unfocused table edge")
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Files") {
		t.Fatal("tab did not focus the table")
	}
	if tablePixelNear(tt.Image(), x, y, colors.Background, 4) {
		t.Fatal("keyboard focus ring is not visible")
	}
	tt.Key(0, ui.KeyDown)
	if st.Selection.Row < 0 {
		t.Fatal("keyboard navigation did not choose a row")
	}
}

func TestTableNativeCellSelectionKeepsItsPaint(t *testing.T) {
	rows := fileRows()
	st := &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionMultiple, Row: -1}}
	st.List.CellSelect = true
	o := fileTable(rows, st, nil)
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).Fill().Padding(20).Background(colors.Background).Children(func() {
			data.Table(c, o).Grow(1)
		})
	}, 520, 260)
	if err := tt.Click("Beta"); err != nil {
		t.Fatal(err)
	}
	row, _ := tt.Find("Beta")
	head, _ := tt.Find("Name")
	y := int(row.Y + row.H/2)
	img := tt.Image()
	wantTablePixel(t, img, int(head.X)+3, y, colors.Primary, "selected cell")
	wantTablePixel(t, img, int(head.X+head.W)+5, y, colors.Background, "unselected cell")
	if !st.List.CellSelected(0, 0) || st.List.CellSelected(0, 1) {
		t.Fatalf("native cell selection = %+v", st.List.Cells)
	}
	tt.Key(0, ui.KeyRight)
	if st.List.Cells.Cursor != "size" || !st.List.CellSelected(0, 1) {
		t.Fatalf("keyboard cell selection = %+v", st.List.Cells)
	}
}

// The footer remains in the table and follows its body.
func TestTableFooterSitsInFrame(t *testing.T) {
	rows := fileRows()
	o := fileTable(rows, &data.TableState{}, nil)
	o.Footer = func(c *ui.Context) ui.Element { return ui.Text(c, "1 to 3 of 3 results") }
	tt := ui.NewTester(tableView(&o), 400, 320)
	foot, ok := tt.Find("1 to 3 of 3 results")
	if !ok {
		t.Fatalf("footer missing: %v", tt.Texts())
	}
	row, ok := tt.Find("Gamma")
	if !ok || foot.Y <= row.Y {
		t.Fatalf("footer=%v row=%v", foot, row)
	}
	if foot.X != theme.Space16 {
		t.Fatalf("footer inset = %v, want %v", foot.X, theme.Space16)
	}
}

// TestTableOverflowScrollsSideways checks columns keep their widths when the
// table is narrower than their sum, and a precise scroll reveals them.
func TestTableOverflowScrollsSideways(t *testing.T) {
	st := &data.TableState{}
	o := data.TableOptions{
		Label: "Wide",
		Columns: []data.TableColumnOptions{
			{ID: "first", Header: "First", Width: 200, MinWidth: 200,
				Cell: func(c *ui.Context, cell data.TableCell) ui.Element { return ui.Text(c, "1") }},
			{ID: "second", Header: "Second", Width: 200, MinWidth: 200,
				Cell: func(c *ui.Context, cell data.TableCell) ui.Element { return ui.Text(c, "2") }},
		},
		Rows: 4, State: st,
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		data.Table(c, o).Grow(1)
	}, 300, 240)
	tt.Frame()
	before, ok := tt.Find("Second")
	if !ok {
		t.Fatalf("second header missing: %v", tt.Texts())
	}
	// A press inside the first column, then a precise sideways scroll.
	tt.Move(100, 100)
	tt.Scroll(100, 100, 120, 0)
	tt.Frame()
	after, ok := tt.Find("Second")
	if !ok {
		t.Fatalf("second header missing after scrolling: %v", tt.Texts())
	}
	if after.X > before.X-80 {
		t.Fatalf("table did not scroll sideways: second header %v -> %v", before.X, after.X)
	}
}

// TestTableEmptyHeaderMatchesPopulatedHeader keeps the empty table's header
// the same native header a populated table shows.
func TestTableEmptyHeaderMatchesPopulatedHeader(t *testing.T) {
	rows := fileRows()
	build := func(rowCount int) (ui.Rect, ui.Rect) {
		st := &data.TableState{}
		o := data.TableOptions{
			Label: "Files",
			Columns: []data.TableColumnOptions{
				{ID: "name", Header: "Name", Sortable: true, Resizable: true},
				{ID: "size", Header: "Size", Sortable: true, Width: 90, Align: ui.End},
			},
			Rows:  rowCount,
			State: st,
		}
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			data.Table(c, o).Grow(1)
		}, 400, 300)
		tt.Frame()
		name, ok := tt.Find("Name")
		if !ok {
			t.Fatalf("Name header missing with %d rows: %v", rowCount, tt.Texts())
		}
		size, ok := tt.Find("Size")
		if !ok {
			t.Fatalf("Size header missing with %d rows: %v", rowCount, tt.Texts())
		}
		return name, size
	}
	emptyName, emptySize := build(0)
	fullName, fullSize := build(len(rows))
	if emptyName != fullName || emptySize != fullSize {
		t.Fatalf("empty header %v/%v differs from populated %v/%v",
			emptyName, emptySize, fullName, fullSize)
	}
}

// TestTableInsertionKeepsKeyedSelection checks a stable row key keeps the
// choice on its row across an insertion, without reporting a change.
func TestTableInsertionKeepsKeyedSelection(t *testing.T) {
	rows := []fileRow{
		{id: "b", name: "Beta"},
		{id: "a", name: "Alpha"},
		{id: "c", name: "Gamma"},
	}
	st := &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionSingle, Row: -1}}
	calls := 0
	build := func(c *ui.Context) {
		theme.ApplyMode(c, false)
		data.Table(c, data.TableOptions{
			Label: "Files",
			Columns: []data.TableColumnOptions{{
				ID: "name", Header: "Name",
				Cell: func(c *ui.Context, cell data.TableCell) ui.Element {
					return ui.Text(c, rows[cell.Row.Index].name)
				},
			}},
			Rows:              len(rows),
			State:             st,
			RowKey:            func(row int) any { return rows[row].id },
			OnSelectionChange: func(data.TableSelection) { calls++ },
		}).Grow(1)
	}
	tt := ui.NewTester(build, 400, 320)
	tt.Frame()
	if err := tt.Click("Beta"); err != nil {
		t.Fatal(err)
	}
	if st.Selection.Row != 0 || calls != 1 {
		t.Fatalf("chosen row = %d, callbacks %d, want 0 and 1", st.Selection.Row, calls)
	}
	rows = append([]fileRow{{id: "z", name: "Zulu"}}, rows...)
	tt.Frame()
	if st.Selection.Row != 1 {
		t.Fatalf("selection did not follow its key: row = %d", st.Selection.Row)
	}
	if calls != 1 {
		t.Fatalf("insertion reported a selection change: %d callbacks", calls)
	}
	zulu, _ := tt.Find("Zulu")
	beta, _ := tt.Find("Beta")
	if !(zulu.Y < beta.Y) {
		t.Fatalf("inserted row not first: zulu=%v beta=%v", zulu, beta)
	}
}
