package navigation

import (
	"image"
	"image/color"
	"testing"

	"github.com/egoist/mygo/ui"
)

// glyphFrame draws one glyph in a box at a fixed place and returns the frame
// and the box, so a test can read its ink back.
func glyphFrame(t *testing.T, draw func(p *ui.Painter, r ui.Rect)) (*ui.Tester, ui.Rect) {
	t.Helper()
	box := ui.Rect{X: 24, Y: 24, W: 16, H: 16}
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Box(c).Absolute().Left(box.X).Top(box.Y).Size(box.W, box.H).
			Draw(func(p *ui.Painter, r ui.Rect) { draw(p, r) })
	}, 64, 64)
	tt.Frame()
	return tt, box
}

// inkBounds returns the box of the ink pixels of r, those unlike the page.
func inkBounds(img *image.RGBA, r ui.Rect) (bounds ui.Rect, ok bool) {
	minX, minY, maxX, maxY := 1<<30, 1<<30, -1, -1
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if !isInk(img.RGBAAt(x, y)) {
				continue
			}
			minX, minY = min(minX, x), min(minY, y)
			maxX, maxY = max(maxX, x), max(maxY, y)
		}
	}
	if maxX < 0 {
		return ui.Rect{}, false
	}
	return ui.Rect{X: float32(minX), Y: float32(minY), W: float32(maxX - minX + 1), H: float32(maxY - minY + 1)}, true
}

func isInk(c color.RGBA) bool {
	return int(c.R)+int(c.G)+int(c.B) < 3*150
}

// inkAt tells whether the pixel at x, y in row band [y0, y1] of r is inked.
func inkXs(img *image.RGBA, r ui.Rect, y0, y1 int) []int {
	var xs []int
	for y := y0; y <= y1; y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if isInk(img.RGBAAt(x, y)) {
				xs = append(xs, x)
			}
		}
	}
	return xs
}

// TestChevronsAreTwoDiagonalStrokes is the regression for the drawn
// indicators: both chevrons were two Painter.Line calls, which paint only
// horizontal and vertical strokes, so each drew one bar. They now stroke a
// ui.Path: a right chevron is taller than wide with its ends at the left of
// its top and bottom thirds and the joint at the right of its middle third;
// a down chevron is wider than tall with both ends at its top and the joint
// at its bottom.
func TestChevronsAreTwoDiagonalStrokes(t *testing.T) {
	for _, tc := range []struct {
		name  string
		draw  func(p *ui.Painter, r ui.Rect)
		tall  bool   // taller than wide, as a right or left chevron is
		joint string // the side of the middle row the joint points to
	}{
		{"right", func(p *ui.Painter, r ui.Rect) { drawChevron(p, &ui.Path{}, r, 0, 1.5, colorInk) }, true, "right"},
		{"down", func(p *ui.Painter, r ui.Rect) { drawChevronDown(p, &ui.Path{}, r, 0, 1.5, colorInk) }, false, "down"},
		{"turned down", func(p *ui.Painter, r ui.Rect) { drawChevron(p, &ui.Path{}, r, 90, 1.5, colorInk) }, false, "down"},
		{"turned left", func(p *ui.Painter, r ui.Rect) { drawChevron(p, &ui.Path{}, r, 180, 1.5, colorInk) }, true, "left"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tt, box := glyphFrame(t, tc.draw)
			img := tt.Image()
			b, ok := inkBounds(img, box)
			if !ok {
				t.Fatal("the chevron paints no ink")
			}
			// A bar is one row: this glyph covers its whole height.
			if b.H < 6 {
				t.Fatalf("chevron is %vx%v DIP: a bar, not two strokes", b.W, b.H)
			}
			if tc.tall && b.H < 1.3*b.W || !tc.tall && b.W < 1.3*b.H {
				t.Fatalf("chevron is %vx%v: not a two-stroke chevron at this angle", b.W, b.H)
			}
			// Centered in its box.
			if d := abs32(b.X + b.W/2 - (box.X + box.W/2)); d > 1.5 {
				t.Fatalf("chevron is %v DIP off center across", d)
			}
			if d := abs32(b.Y + b.H/2 - (box.Y + box.H/2)); d > 1.5 {
				t.Fatalf("chevron is %v DIP off center down", d)
			}
			// Every row between the ends is inked: the two strokes meet.
			for y := int(b.Y); y < int(b.Y+b.H); y++ {
				if len(inkXs(img, b, y, y)) == 0 {
					t.Fatalf("chevron has an empty row at %d", y)
				}
			}
			// A wide chevron has both stroke ends on its top row, on both
			// sides of the middle, and its joint below the middle; a tall one
			// has one end at the top, the other at the bottom, and its joint
			// in the middle, on the side away from them.
			if !tc.tall {
				if !hasInkBothSides(img, b, int(b.Y)) {
					t.Fatalf("chevron top row does not hold two stroke ends")
				}
				if !hasInkMiddle(img, b, int(b.Y+b.H)-1) {
					t.Fatalf("chevron bottom row does not meet in the middle")
				}
				return
			}
			for _, y := range []int{int(b.Y), int(b.Y+b.H) - 1} {
				side := avg(inkXs(img, b, y, y))
				if tc.joint == "right" && side > b.X+b.W/2 {
					t.Fatalf("chevron row %d is on the joint's side: %v", y, side)
				}
				if tc.joint == "left" && side < b.X+b.W/2 {
					t.Fatalf("chevron row %d is on the joint's side: %v", y, side)
				}
			}
			mid := b.Y + b.H/2
			joint := inkXs(img, b, int(mid), int(mid))
			if len(joint) == 0 {
				t.Fatal("chevron does not cross its middle row")
			}
			if tc.joint == "right" && avg(joint) <= b.X+b.W/2 {
				t.Fatalf("chevron joint at %v, want the right of the middle", avg(joint))
			}
			if tc.joint == "left" && avg(joint) >= b.X+b.W/2 {
				t.Fatalf("chevron joint at %v, want the left of the middle", avg(joint))
			}
		})
	}
}

// TestCheckIsTwoDiagonalStrokes checks HeroUI's tick: a short arm down to a
// valley left of center and a long arm up to the right, not a bar.
func TestCheckIsTwoDiagonalStrokes(t *testing.T) {
	tt, box := glyphFrame(t, func(p *ui.Painter, r ui.Rect) { drawCheck(p, &ui.Path{}, r, 1, colorInk) })
	img := tt.Image()
	b, ok := inkBounds(img, box)
	if !ok {
		t.Fatal("the checkmark paints no ink")
	}
	if b.H < 5 || b.W < 6 {
		t.Fatalf("checkmark is %vx%v DIP: a bar, not two strokes", b.W, b.H)
	}
	if d := abs32(b.X + b.W/2 - (box.X + box.W/2)); d > 1.5 {
		t.Fatalf("checkmark is %v DIP off center across", d)
	}
	if d := abs32(b.Y + b.H/2 - (box.Y + box.H/2)); d > 1.5 {
		t.Fatalf("checkmark is %v DIP off center down", d)
	}
	// Its lowest row is the valley, left of center; its highest row is the
	// long arm's end, right of center.
	low := inkXs(img, b, int(b.Y+b.H)-1, int(b.Y+b.H)-1)
	high := inkXs(img, b, int(b.Y), int(b.Y))
	if len(low) == 0 || len(high) == 0 {
		t.Fatalf("checkmark ends missing: low %v high %v", low, high)
	}
	if avg(low) >= b.X+b.W/2 || avg(high) <= b.X+b.W/2 {
		t.Fatalf("checkmark is not tick-shaped: low x %v high x %v, center %v", avg(low), avg(high), b.X+b.W/2)
	}
	for y := int(b.Y); y < int(b.Y+b.H); y++ {
		if len(inkXs(img, b, y, y)) == 0 {
			t.Fatalf("checkmark has an empty row at %d", y)
		}
	}
}

// hasInkBothSides tells whether row y of b has ink left and right of the
// b's middle.
func hasInkBothSides(img *image.RGBA, b ui.Rect, y int) bool {
	left, right := false, false
	for _, x := range inkXs(img, b, y, y) {
		switch {
		case float32(x) < b.X+b.W/2-1:
			left = true
		case float32(x) > b.X+b.W/2+1:
			right = true
		}
	}
	return left && right
}

// hasInkMiddle tells whether row y of b has ink near its middle.
func hasInkMiddle(img *image.RGBA, b ui.Rect, y int) bool {
	for _, x := range inkXs(img, b, y, y) {
		if abs32(float32(x)-(b.X+b.W/2)) <= 2 {
			return true
		}
	}
	return false
}

func avg(xs []int) float32 {
	var sum float32
	for _, x := range xs {
		sum += float32(x)
	}
	return sum / float32(len(xs))
}

func abs32(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

var colorInk = ui.RGB(0, 0, 0)
