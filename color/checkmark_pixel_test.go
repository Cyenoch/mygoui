package color

import (
	"image"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// swatchCheckInk returns the box and rows of the check ink in scan, observed
// against the swatch color sampled at the top of scan.
func swatchCheckInk(img *image.RGBA, scan ui.Rect) (ui.Rect, [][2]int, bool) {
	bg := img.RGBAAt(int(scan.X+scan.W/2), int(scan.Y))
	ink := func(x, y int) bool {
		c := img.RGBAAt(x, y)
		d := abs(int(c.R)-int(bg.R)) + abs(int(c.G)-int(bg.G)) + abs(int(c.B)-int(bg.B))
		return d > 200
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

// assertSwatchCheck checks the selected swatch's tick: two diagonal strokes,
// wider than tall, with its valley left of center and its long arm right of
// it, centered in swatch. It used to be two Painter.Line calls, which paint
// only horizontal and vertical strokes: a short bar and an upright one.
func assertSwatchCheck(t *testing.T, img *image.RGBA, swatch ui.Rect, where string) {
	t.Helper()
	scan := ui.Rect{
		X: swatch.X + swatch.W*0.25, Y: swatch.Y + swatch.H*0.25,
		W: swatch.W * 0.5, H: swatch.H * 0.5,
	}
	b, rows, ok := swatchCheckInk(img, scan)
	if !ok {
		t.Fatalf("%s: no check ink", where)
	}
	if b.H < 5 || b.W < 6 || b.W < b.H {
		t.Fatalf("%s: check is %vx%v DIP, not a two-stroke tick", where, b.W, b.H)
	}
	if d := abs32(b.X + b.W/2 - (swatch.X + swatch.W/2)); d > 3 {
		t.Fatalf("%s: check is %v DIP off center across", where, d)
	}
	if d := abs32(b.Y + b.H/2 - (swatch.Y + swatch.H/2)); d > 3 {
		t.Fatalf("%s: check is %v DIP off center down", where, d)
	}
	top := int(b.Y) - int(scan.Y)
	bottom := int(b.Y+b.H) - int(scan.Y) - 1
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

// TestColorSwatchSelectedCheckIsDrawn checks the tick on a selected swatch,
// in both inks: white on the dark red, black on the light amber.
func TestColorSwatchSelectedCheckIsDrawn(t *testing.T) {
	for _, tc := range []struct {
		name  string
		color ui.Color
	}{
		{"white ink", ui.Hex("#e0010f")},
		{"dark ink", ui.Hex("#ffd60a")},
	} {
		tt := ui.NewTester(func(c *ui.Context) {
			theme.Apply(c)
			ui.Column(c).Fill().Padding(16).Children(func() {
				ColorSwatch(c, ColorSwatchOptions{
					Color: tc.color, Label: "chosen", Selected: true, Shape: SwatchSquare,
				})
			})
		}, 200, 200)
		tt.Frame()
		sw, ok := tt.Find("chosen")
		if !ok {
			t.Fatalf("%s: selected swatch is not shown", tc.name)
		}
		assertSwatchCheck(t, tt.Image(), sw, tc.name)
	}
}

// TestColorSwatchPickerSelectedCheckIsDrawn checks the tick on the chosen
// swatch of a palette.
func TestColorSwatchPickerSelectedCheckIsDrawn(t *testing.T) {
	value := ui.Hex("#e0010f")
	colors := []ui.Color{ui.Hex("#e0010f"), ui.Hex("#0a64b5")}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ui.Column(c).Fill().Padding(16).Children(func() {
			ColorSwatchPicker(c, ColorSwatchPickerOptions{
				Label: "Palette", Colors: colors, Names: []string{"Crimson", "Ocean"},
				Shape: SwatchSquare, Value: &value,
			})
		})
	}, 240, 200)
	tt.Frame()
	sw, ok := tt.Find("Crimson")
	if !ok {
		t.Fatal("the chosen swatch of the palette is not shown")
	}
	assertSwatchCheck(t, tt.Image(), sw, "palette")
}
