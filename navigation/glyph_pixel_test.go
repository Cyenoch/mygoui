package navigation_test

import (
	"image"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/navigation"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// glyphInk returns the inked rows of scan, those unlike the background
// sampled at its top-left corner, as the leftmost and rightmost column of
// each row, and the box holding them.
func glyphInk(img *image.RGBA, scan ui.Rect) (ui.Rect, [][2]int, bool) {
	bg := img.RGBAAt(int(scan.X)+1, int(scan.Y)+1)
	ink := func(x, y int) bool {
		c := img.RGBAAt(x, y)
		d := abs(int(c.R)-int(bg.R)) + abs(int(c.G)-int(bg.G)) + abs(int(c.B)-int(bg.B))
		return d > 110
	}
	minX, minY, maxX, maxY := 1<<30, 1<<30, -1, -1
	var rows [][2]int
	for y := int(scan.Y); y < int(scan.Y+scan.H); y++ {
		left, right := -1, -1
		for x := int(scan.X); x < int(scan.X+scan.W); x++ {
			if !ink(x, y) {
				continue
			}
			if left < 0 {
				left = x
			}
			right = x
			minX, maxX = min(minX, x), max(maxX, x)
			minY, maxY = min(minY, y), max(maxY, y)
		}
		rows = append(rows, [2]int{left, right})
	}
	if maxX < 0 {
		return ui.Rect{}, nil, false
	}
	return ui.Rect{X: float32(minX), Y: float32(minY), W: float32(maxX - minX + 1), H: float32(maxY - minY + 1)}, rows, true
}

// assertDownChevron checks a chevron pointing down in box: two diagonal
// strokes meeting at the bottom, centered, not the one bar the old
// Painter.Line calls drew.
func assertDownChevron(t *testing.T, img *image.RGBA, box ui.Rect, where string) {
	t.Helper()
	b, rows, ok := glyphInk(img, box)
	if !ok {
		t.Fatalf("%s: no chevron ink", where)
	}
	if b.H < 5 || b.W < 6 || b.W < 1.3*b.H {
		t.Fatalf("%s: chevron is %vx%v DIP, not a two-stroke down chevron", where, b.W, b.H)
	}
	checkCentered(t, b, box, where)
	top, bottom := int(b.Y-box.Y), int(b.Y+b.H-box.Y)-1
	if r := rows[top]; r[0] < 0 || float32(r[0]) >= b.X+b.W/2-1 || float32(r[1]) <= b.X+b.W/2+1 {
		t.Fatalf("%s: chevron top row %v has no two ends", where, r)
	}
	if r := rows[bottom]; r[0] < 0 || float32(r[0]) > b.X+b.W/2+2 || float32(r[1]) < b.X+b.W/2-2 {
		t.Fatalf("%s: chevron bottom row %v does not meet in the middle", where, r)
	}
	for i := top; i <= bottom; i++ {
		if rows[i][0] < 0 {
			t.Fatalf("%s: chevron has an empty row %d", where, i)
		}
	}
}

// assertRightChevron checks a chevron pointing right in box.
func assertRightChevron(t *testing.T, img *image.RGBA, box ui.Rect, where string) {
	t.Helper()
	b, rows, ok := glyphInk(img, box)
	if !ok {
		t.Fatalf("%s: no chevron ink", where)
	}
	if b.H < 5 || b.W < 4 || b.H < 1.3*b.W {
		t.Fatalf("%s: chevron is %vx%v DIP, not a two-stroke right chevron", where, b.W, b.H)
	}
	checkCentered(t, b, box, where)
	top, bottom := int(b.Y-box.Y), int(b.Y+b.H-box.Y)-1
	for i := top; i <= bottom; i++ {
		if rows[i][0] < 0 {
			t.Fatalf("%s: chevron has an empty row %d: not two strokes", where, i)
		}
	}
	middle := rows[top+(bottom-top)/2]
	if float32(middle[0]) < b.X+b.W/2 {
		t.Fatalf("%s: chevron joint row %v is not on the right", where, middle)
	}
}

// assertCheck checks the tick in box: a short arm down to a valley left
// of center and a long arm up to the right, not a bar.
func assertCheck(t *testing.T, img *image.RGBA, box ui.Rect, where string) {
	t.Helper()
	b, rows, ok := glyphInk(img, box)
	if !ok {
		t.Fatalf("%s: no check ink", where)
	}
	if b.H < 5 || b.W < 6 || b.W < b.H {
		t.Fatalf("%s: check is %vx%v DIP, not a two-stroke tick", where, b.W, b.H)
	}
	checkCentered(t, b, box, where)
	top, bottom := int(b.Y-box.Y), int(b.Y+b.H-box.Y)-1
	for i := top; i <= bottom; i++ {
		if rows[i][0] < 0 {
			t.Fatalf("%s: check has an empty row %d", where, i)
		}
	}
	low, high := rows[bottom], rows[top]
	if float32((low[0]+low[1])/2) >= b.X+b.W/2 || float32((high[0]+high[1])/2) <= b.X+b.W/2 {
		t.Fatalf("%s: check is not tick-shaped: low %v high %v, center %v", where, low, high, b.X+b.W/2)
	}
}

func checkCentered(t *testing.T, b, box ui.Rect, where string) {
	t.Helper()
	if d := abs32(b.X + b.W/2 - (box.X + box.W/2)); d > 3 {
		t.Fatalf("%s: glyph is %v DIP off center across", where, d)
	}
	if d := abs32(b.Y + b.H/2 - (box.Y + box.H/2)); d > 3 {
		t.Fatalf("%s: glyph is %v DIP off center down", where, d)
	}
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func abs32(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

// TestDisclosureAndAccordionChevronsAreDrawn checks the indicators of both
// disclosure headers in the rendered UI, light and dark: a down chevron of
// two diagonal strokes at the header's end.
func TestDisclosureAndAccordionChevronsAreDrawn(t *testing.T) {
	for _, dark := range []bool{false, true} {
		name := "light"
		if dark {
			name = "dark"
		}
		expanded := []string{}
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c)
			ui.Column(c).Fill().Padding(16).Gap(16).Children(func() {
				navigation.Accordion(c, navigation.AccordionOptions{
					Expanded: &expanded,
					Items: []navigation.AccordionItemOptions{
						{Key: "general", Title: "General", Body: func() { ui.Text(c, "Body") }},
					},
				})
				navigation.Disclosure(c, navigation.DisclosureOptions{
					Title: "Details", Body: func() { ui.Text(c, "Details body") },
				})
			})
		}, 640, 480)
		if dark {
			tt.SetDark(true)
		}
		tt.Frame()

		acc, ok := tt.Find("General")
		if !ok {
			t.Fatalf("%s: accordion header is not shown", name)
		}
		// The accordion indicator is the 16 DIP box at the trigger's end.
		accBox := ui.Rect{X: acc.X + acc.W - 16, Y: acc.Y + (acc.H-16)/2, W: 16, H: 16}
		assertDownChevron(t, tt.Image(), accBox, name+" accordion")

		dis, ok := tt.Find("Details")
		if !ok {
			t.Fatalf("%s: disclosure header is not shown", name)
		}
		// The disclosure trigger pads 8 DIP at its end before the indicator.
		disBox := ui.Rect{X: dis.X + dis.W - 8 - 16, Y: dis.Y + (dis.H-16)/2, W: 16, H: 16}
		assertDownChevron(t, tt.Image(), disBox, name+" disclosure")
	}
}

// TestMenuIndicatorGlyphsAreDrawn checks the selected row's check and a
// submenu row's chevron in the rendered dropdown, light and dark.
func TestMenuIndicatorGlyphsAreDrawn(t *testing.T) {
	for _, dark := range []bool{false, true} {
		name := "light"
		if dark {
			name = "dark"
		}
		open := false
		selected := []string{"edit"}
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c)
			navigation.Dropdown(c, navigation.DropdownOptions{
				Open: &open, Label: "Actions",
				Menu: navigation.MenuOptions{
					Selection: navigation.MenuSelectMultiple, Selected: &selected,
					Sections: []navigation.MenuSectionOptions{{Items: []navigation.MenuItemOptions{
						{Key: "edit", Label: "Edit"},
						{Key: "share", Label: "Share", Submenu: []navigation.MenuItemOptions{
							{Key: "email", Label: "Email"},
						}},
					}}},
				},
			})
		}, 640, 480)
		if dark {
			tt.SetDark(true)
		}
		tt.Frame()
		if err := tt.Click("Actions"); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		// The panel fades in over 150 ms; read it once it has arrived.
		tt.Advance(300 * time.Millisecond)
		edit, ok := tt.Find("Edit")
		if !ok {
			t.Fatalf("%s: the selected row is not shown", name)
		}
		// The indicator is absolute, 8 DIP into the row, in a 16 DIP box.
		checkBox := ui.Rect{X: edit.X + 8, Y: edit.Y + (edit.H-16)/2, W: 16, H: 16}
		assertCheck(t, tt.Image(), checkBox, name+" menu check")

		share, ok := tt.Find("Share")
		if !ok {
			t.Fatalf("%s: the submenu row is not shown", name)
		}
		// A submenu row ends its content 28 DIP in, before the 14 DIP
		// chevron; its indicator moves to the end.
		chevBox := ui.Rect{X: share.X + share.W - 28 - 14, Y: share.Y + (share.H-14)/2, W: 14, H: 14}
		assertRightChevron(t, tt.Image(), chevBox, name+" menu chevron")
	}
}
