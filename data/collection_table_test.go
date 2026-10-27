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

func TestTableSecondaryRowsKeepCardBoundary(t *testing.T) {
	rows := fileRows()
	st := &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionSingle, Row: -1}}
	o := fileTable(rows, st, nil)
	o.Variant = data.TableSecondary
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(20).Background(ui.RGB(180, 200, 220)).Children(func() {
			data.Table(c, o).Height(153)
		})
	}, 440, 220)
	tt.Move(2, 2)
	tt.Frame()
	table, ok := tt.Find("Files")
	if !ok {
		t.Fatal("table missing")
	}
	before := tt.Image()
	wantTablePixel(t, before, int(table.X)+1, int(table.Y+table.H)-2, ui.RGB(180, 200, 220), "rounded corner")
	if before.RGBAAt(int(table.X+table.W/2), int(table.Y+table.H)-1) ==
		before.RGBAAt(int(table.X+table.W/2), int(table.Y+table.H)-3) {
		t.Fatal("card bottom border is missing")
	}
	for _, tc := range []struct {
		name   string
		points [][2]int
	}{
		{"Beta", [][2]int{
			{int(table.X) + 1, int(table.Y) + 34},
			{int(table.X+table.W) - 2, int(table.Y) + 34},
		}},
		{"Gamma", [][2]int{
			{int(table.X) + 1, int(table.Y+table.H) - 2},
			{int(table.X+table.W) - 2, int(table.Y+table.H) - 2},
			{int(table.X + table.W/2), int(table.Y+table.H) - 1},
		}},
	} {
		row, ok := tt.Find(tc.name)
		if !ok {
			t.Fatalf("row %q missing", tc.name)
		}
		tt.Move(row.X+row.W/2, row.Y+row.H/2)
		tt.Frame()
		hover := tt.Image()
		for _, point := range tc.points {
			if got, want := hover.RGBAAt(point[0], point[1]), before.RGBAAt(point[0], point[1]); got != want {
				t.Errorf("%s hover changed card boundary at %v: got %v, want %v", tc.name, point, got, want)
			}
		}
		x, y := int(table.X+table.W/2), int(row.Y+row.H/2)
		if hover.RGBAAt(x, y) == before.RGBAAt(x, y) {
			t.Errorf("%s was not visibly hovered", tc.name)
		}
		if err := tt.Click(tc.name); err != nil {
			t.Fatal(err)
		}
		tt.Move(2, 2)
		tt.Frame()
		selected := tt.Image()
		for _, point := range tc.points {
			if got, want := selected.RGBAAt(point[0], point[1]), before.RGBAAt(point[0], point[1]); got != want {
				t.Errorf("%s selection changed card boundary at %v: got %v, want %v", tc.name, point, got, want)
			}
		}
		if selected.RGBAAt(x, y) == before.RGBAAt(x, y) {
			t.Errorf("%s was not visibly selected", tc.name)
		}
	}
}

func TestTableSecondaryLeavesHeaderVisible(t *testing.T) {
	rows := fileRows()
	st := &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionSingle, Row: 0}}
	o := fileTable(rows, st, nil)
	o.Variant = data.TableSecondary
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Column(c).Fill().Padding(20).Children(func() { data.Table(c, o).Height(153) })
	}, 440, 220)
	tt.Move(2, 2)
	tt.Frame()
	header, ok := tt.Find("Name")
	if !ok {
		t.Fatal("header missing")
	}
	img := tt.Image()
	background := img.RGBAAt(int(header.X+header.W)-10, int(header.Y+header.H/2))
	ink := 0
	for y := int(header.Y); y < int(header.Y+header.H); y++ {
		for x := int(header.X); x < int(header.X+80); x++ {
			if img.RGBAAt(x, y) != background {
				ink++
			}
		}
	}
	if ink == 0 {
		t.Fatal("card clipping hid the header text")
	}
	if err := tt.Click("Name"); err != nil {
		t.Fatal(err)
	}
	if st.Sort.Column != "name" {
		t.Fatalf("header sorting stopped working: %+v", st.Sort)
	}
}

func TestTableSecondaryCellSelectionKeepsCardBoundary(t *testing.T) {
	rows := fileRows()
	st := &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionMultiple, Row: -1}}
	st.List.CellSelect = true
	o := fileTable(rows, st, nil)
	o.Variant = data.TableSecondary
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).Fill().Padding(20).Background(ui.RGB(180, 200, 220)).Children(func() {
			data.Table(c, o).Height(153)
		})
	}, 440, 220)
	tt.Move(2, 2)
	tt.Frame()
	table, _ := tt.Find("Files")
	before := tt.Image()
	if err := tt.Click("Gamma"); err != nil {
		t.Fatal(err)
	}
	tt.Move(2, 2)
	tt.Frame()
	after := tt.Image()
	points := [][2]int{
		{int(table.X) + 1, int(table.Y+table.H) - 2},
		{int(table.X+table.W) - 2, int(table.Y+table.H) - 2},
		{int(table.X + table.W/2), int(table.Y+table.H) - 1},
	}
	for _, point := range points {
		if got, want := after.RGBAAt(point[0], point[1]), before.RGBAAt(point[0], point[1]); got != want {
			t.Errorf("cell selection changed card boundary at %v: got %v, want %v", point, got, want)
		}
	}
	row, _ := tt.Find("Gamma")
	wantTablePixel(t, after, int(table.X)+100, int(row.Y+row.H/2), colors.Primary, "selected cell")
	if !st.List.CellSelected(2, 0) || st.List.CellSelected(2, 1) {
		t.Fatalf("native cell selection = %+v", st.List.Cells)
	}

	name, _ := tt.Find("Name")
	if got := after.RGBAAt(int(name.X+name.W/2), int(row.Y)-11); got.R < 200 {
		t.Fatalf("selected cell cursor border not visible: %v", got)
	}
	tt.Key(ui.Shift, ui.KeyRight)
	if st.List.Cells.Cursor != "size" || !st.List.CellSelected(2, 1) {
		t.Fatalf("keyboard cell selection = %+v", st.List.Cells)
	}
	selection := tt.Image()
	size, _ := tt.Find("Size")
	wantTablePixel(t, selection, int(name.X+name.W/2), int(row.Y+row.H/2), colors.Primary, "first selected cell")
	wantTablePixel(t, selection, int(size.X+size.W/2), int(row.Y+row.H/2), colors.Primary, "second selected cell")
	for _, point := range points {
		if got, want := selection.RGBAAt(point[0], point[1]), before.RGBAAt(point[0], point[1]); got != want {
			t.Errorf("cell selection changed card boundary at %v: got %v, want %v", point, got, want)
		}
	}
}

func TestTableSecondaryKeepsNativeCellControlsAndVirtualizedRows(t *testing.T) {
	st := &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionSingle, Row: -1}}
	var original ui.Theme
	var wrongTheme bool
	o := data.TableOptions{
		Label:   "Rows",
		Variant: data.TableSecondary,
		Columns: []data.TableColumnOptions{{ID: "name", Header: "Name"}},
		Rows:    100,
		State:   st,
		Cell: func(c *ui.Context, cell data.TableCell) ui.Element {
			if *c.Theme() != original {
				wrongTheme = true
			}
			return ui.Button(c, fmt.Sprintf("Row %d", cell.Row.Index))
		},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		original = *c.Theme()
		ui.Column(c).Fill().Padding(20).Children(func() {
			data.Table(c, o).Height(153)
		})
	}, 440, 220)
	if err := tt.Click("Row 0"); err != nil {
		t.Fatal(err)
	}
	if !tt.Focused("Row 0") {
		t.Fatal("native button did not receive focus")
	}
	tt.Scroll(200, 130, 0, 5000)
	tt.Frame()
	if wrongTheme {
		t.Fatal("a cell was built with the internal row theme")
	}
	if !tt.HasText("Row 90") {
		t.Fatalf("scroll did not build later rows: %v", tt.Texts())
	}
}

func TestTableSecondaryHoveredRowFollowsScroll(t *testing.T) {
	st := &data.TableState{Selection: data.TableSelection{Mode: data.TableSelectionSingle, Row: -1}}
	o := data.TableOptions{
		Label:   "Rows",
		Variant: data.TableSecondary,
		Columns: []data.TableColumnOptions{{ID: "name", Header: "Name"}},
		Rows:    100,
		State:   st,
		Cell: func(c *ui.Context, cell data.TableCell) ui.Element {
			return ui.Text(c, fmt.Sprintf("Row %d", cell.Row.Index))
		},
	}
	var hover ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		hover = c.Theme().SurfaceHover
		ui.Column(c).Fill().Padding(20).Children(func() { data.Table(c, o).Height(153) })
	}, 440, 220)
	table, _ := tt.Find("Rows")
	x, y := table.X+table.W/2, table.Y+100
	tt.Move(x, y)
	tt.Scroll(x, y, 0, 1000)
	tt.Frame()
	wantTablePixel(t, tt.Image(), int(x), int(y), hover, "scrolled hovered row")
	if !tt.HasText("Row 20") && !tt.HasText("Row 21") && !tt.HasText("Row 22") && !tt.HasText("Row 23") {
		t.Fatalf("table did not scroll to later rows: %v", tt.Texts())
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

func TestCardTableHasOneTopBoundaryBelowItsHeader(t *testing.T) {
	for _, dark := range []bool{false, true} {
		t.Run(fmt.Sprint(dark), func(t *testing.T) {
			o := fileTable(fileRows(), &data.TableState{}, nil)
			o.Variant = data.TableSecondary
			background := ui.RGB(180, 200, 220)
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, dark)
				ui.Column(c).Fill().Padding(20).Background(background).Children(func() {
					data.Table(c, o).Grow(1)
				})
			}, 440, 240)
			head, ok := tt.Find("Name")
			if !ok {
				t.Fatal("header missing")
			}
			img := tt.Image()
			// A full-width native divider above the rounded card would make a
			// second top edge outside its corner, as in the reported screenshot.
			wantTablePixel(t, img, int(head.X)+2, int(head.Y+head.H), background, "space above card corner")
		})
	}
}

func TestCardTableKeepsHeaderFeedback(t *testing.T) {
	for _, rows := range []int{0, 3} {
		t.Run(fmt.Sprint(rows), func(t *testing.T) {
			o := fileTable(fileRows(), &data.TableState{}, nil)
			o.Rows, o.Variant = rows, data.TableSecondary
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				data.Table(c, o).Grow(1)
			}, 440, 240)
			head, _ := tt.Find("Name")
			x, y := int(head.X)+4, int(head.Y+head.H/2)
			before := tt.Image().RGBAAt(x, y)
			tt.Move(head.X+head.W/2, head.Y+head.H/2)
			if got := tt.Image().RGBAAt(x, y); got == before {
				t.Fatal("sortable header has no hover feedback")
			}
			// The native resize grip is on the first column's trailing edge.
			edge := int(head.X+head.W) - 1
			tt.Move(2, 100)
			before = tt.Image().RGBAAt(edge, y)
			if tablePixelNear(tt.Image(), edge, y, theme.ColorsForMode(false).Background, 4) {
				t.Fatal("resize grip is missing at rest")
			}
			tt.Move(float32(edge), float32(y))
			if got := tt.Image().RGBAAt(edge, y); got == before {
				t.Fatal("resize grip has no hover feedback")
			}
			tt.Press(float32(edge), float32(y))
			tt.Move(float32(edge+24), float32(y))
			tt.Release(float32(edge+24), float32(y))
			if width := o.State.List.Columns.Widths["name"]; width != head.W+24 {
				t.Fatalf("restored grip no longer resizes its native column: %g", width)
			}
			if err := tt.Click("Name"); err != nil {
				t.Fatal(err)
			}
			if o.State.Sort.Column != "name" {
				t.Fatal("restored feedback no longer sorts its native column")
			}
		})
	}
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
