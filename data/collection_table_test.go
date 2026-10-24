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
		wantTablePixel(t, img, x, y, colors.SurfaceSecondary, "chosen row fill")
		if tablePixelNear(img, x, y, colors.Accent, 6) {
			t.Fatalf("chosen row shows the native accent at (%d, %d)", x, y)
		}
	}
	// The row below is not chosen: its fill is the table surface.
	other, _ := tt.Find("Gamma")
	wantTablePixel(t, img, int(head.X)+2, int(other.Y+other.H/2), colors.Surface, "unchosen row fill")
}

// TestTableSecondaryVariantIsFlush checks the secondary variant paints no
// surface card, while the primary one does.
func TestTableSecondaryVariantIsFlush(t *testing.T) {
	rows := fileRows()
	for _, variant := range []data.TableVariant{data.TablePrimary, data.TableSecondary} {
		st := &data.TableState{}
		o := fileTable(rows, st, nil)
		o.Variant = variant
		var colors theme.Palette
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			colors = theme.Colors(c)
			data.Table(c, o).Grow(1)
		}, 400, 300)
		tt.Frame()
		row, ok := tt.Find("Beta")
		if !ok {
			t.Fatalf("row missing: %v", tt.Texts())
		}
		want := colors.Background
		if variant == data.TablePrimary {
			want = colors.Surface
		}
		// The row's left padding is outside every cell's content box.
		wantTablePixel(t, tt.Image(), 6, int(row.Y+row.H/2), want,
			fmt.Sprintf("variant %d row fill", variant))
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
	// The header box is the column; the end-aligned cell text stops at the
	// toolkit's 10 DIP cell padding plus the table's 6 DIP inset.
	wantEnd := endHead.X + endHead.W - 16
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

// TestTableHeaderChrome keeps HeroUI's header chrome: type smaller than the
// cells, and a separator line at half strength, lighter than the toolkit's
// divider.
func TestTableHeaderChrome(t *testing.T) {
	rows := fileRows()
	var colors theme.Palette
	st := &data.TableState{}
	o := data.TableOptions{
		Label: "Files",
		Columns: []data.TableColumnOptions{{
			ID: "name", Header: "Name",
			Cell: func(c *ui.Context, cell data.TableCell) ui.Element {
				return ui.Text(c, rows[cell.Row.Index].name)
			},
		}},
		Rows: len(rows), State: st,
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		data.Table(c, o).Grow(1)
	}, 400, 240)
	tt.Frame()
	head, ok := tt.Find("Name")
	if !ok {
		t.Fatalf("header missing: %v", tt.Texts())
	}
	body, ok := tt.Find("Beta")
	if !ok {
		t.Fatalf("row missing: %v", tt.Texts())
	}
	if head.H >= body.H {
		t.Fatalf("header type %v is not smaller than cell type %v", head.H, body.H)
	}

	img := tt.Image()
	wantTablePixel(t, img, 200, 16, colors.SurfaceSecondary, "header band")
	line := img.RGBAAt(200, 32)
	if tablePixelNear(img, 200, 32, colors.Border, 3) {
		t.Fatalf("header border uses the toolkit's divider %d,%d,%d", line.R, line.G, line.B)
	}
	if line.R >= colors.SurfaceSecondary.R || line.R <= colors.Border.R {
		t.Fatalf("header border %d is not the separator at half strength", line.R)
	}
}

// TestTableFooterSitsInFrame checks the footer keeps HeroUI's inset, inside
// the frame, under the body card.
func TestTableFooterSitsInFrame(t *testing.T) {
	rows := fileRows()
	var colors theme.Palette
	st := &data.TableState{}
	o := data.TableOptions{
		Label: "Files",
		Columns: []data.TableColumnOptions{
			{ID: "name", Header: "Name", Cell: func(c *ui.Context, cell data.TableCell) ui.Element {
				return ui.Text(c, rows[cell.Row.Index].name)
			}},
		},
		Rows: len(rows), State: st,
		Footer: func(c *ui.Context) ui.Element {
			return ui.Row(c).FillWidth().Children(func() {
				ui.Text(c, "1 to 3 of 3 results")
			})
		},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		data.Table(c, o).Grow(1)
	}, 400, 320)
	tt.Frame()
	foot, ok := tt.Find("1 to 3 of 3 results")
	if !ok {
		t.Fatalf("footer missing: %v", tt.Texts())
	}
	if foot.X != 20 {
		t.Fatalf("footer text starts at %v, want the 4 DIP frame plus 16 DIP inset (20)", foot.X)
	}
	img := tt.Image()
	// The body card ends above the footer; the footer sits on the frame.
	wantTablePixel(t, img, 200, int(foot.Y-16), colors.Surface, "body card above the footer")
	wantTablePixel(t, img, 200, int(foot.Y+foot.H/2), colors.SurfaceSecondary, "footer frame")
	// The frame's own 4 DIP padding runs down the sides, and the body card
	// starts inside it.
	wantTablePixel(t, img, 2, int(foot.Y-60), colors.SurfaceSecondary, "frame beside the card")
	wantTablePixel(t, img, 6, int(foot.Y-60), colors.Surface, "card inside the frame")
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
