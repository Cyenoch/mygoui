package overlay

import (
	"image"
	"image/color"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestToastCloseGlyphIsOneCenteredX is the regression for the reported toast
// close artifact: a pale disc with a dark horizontal bar near its top. The
// close button painted its own face in Draw with two Painter.Line calls for
// the X; Painter.Line paints only horizontal and vertical strokes, so both
// diagonals collapsed into one bar. It now holds a single lucide X, sized 12
// DIP in a 20 DIP button and centered by the layout.
func TestToastCloseGlyphIsOneCenteredX(t *testing.T) {
	for _, mode := range []struct {
		name string
		dark bool
	}{{"light", false}, {"dark", true}} {
		t.Run(mode.name, func(t *testing.T) {
			const winW, winH = 640, 480
			var state ToastState
			var ctx *ui.Context
			tt := ui.NewTester(toastTestView(&state, &ctx), winW, winH)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			if mode.dark {
				tt.SetDark(true)
			}
			tt.Frame()
			if err := tt.Click("Show"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()

			pal := theme.ColorsForMode(mode.dark)
			close, ok := tt.Find("Close")
			if !ok {
				t.Fatal("toast shows no close button")
			}
			title, ok := tt.Find("File saved")
			if !ok {
				t.Fatal("toast shows no title")
			}
			action, ok := tt.Find("Undo")
			if !ok {
				t.Fatal("toast shows no action")
			}
			card, ok := tt.Find("File saved")
			if !ok {
				t.Fatal("toast shows no card")
			}
			if gap := card.X + card.W - action.X - action.W; abs32(gap-theme.Space16) > 1 {
				t.Fatalf("action reserves space for the close button: right gap=%v", gap)
			}
			if close.Y >= card.Y || close.X+close.W <= card.X+card.W {
				t.Fatalf("close does not float across the upper-right corner: close=%+v card=%+v", close, card)
			}
			// HeroUI's -end-1/-top-1: the 20 DIP close button overlaps the
			// card's corner by 4 DIPs on both axes.
			if card.Y-close.Y != 4 || close.X+close.W-(card.X+card.W) != 4 {
				t.Fatalf("close sits off the card corner: card=%+v close=%+v", card, close)
			}
			if close.X < action.X+action.W && close.X+close.W > action.X &&
				close.Y < action.Y+action.H && close.Y+close.H > action.Y {
				t.Fatalf("floating close overlaps the action: close=%+v action=%+v", close, action)
			}

			// Hover reveals one centered X on a pale face.
			tt.Move(close.X+close.W/2, close.Y+close.H/2)
			tt.Frame()
			revealedClose, _ := tt.Find("Close")
			revealedTitle, _ := tt.Find("File saved")
			if revealedClose != close || revealedTitle != title {
				t.Fatalf("revealing the close button shifted the layout: close %+v -> %+v, title %+v -> %+v",
					close, revealedClose, title, revealedTitle)
			}
			img := tt.Image()
			bg := img.RGBAAt(int(close.X+close.W/2), int(close.Y+2))
			if sameColor(bg, pal.Surface, 2) {
				t.Fatalf("hover did not reveal the close face: %v", bg)
			}
			assertSingleX(t, img, close, bg, 0)

			// The focus ring belongs to keyboard focus alone.
			tt.Move(2, 2)
			tt.Frame()
			focused := false
			for range 8 {
				tt.Key(0, ui.KeyTab)
				tt.Frame()
				if tt.Focused("Close") {
					focused = true
					break
				}
			}
			if !focused {
				t.Fatal("Tab never reached the close button")
			}
			img = tt.Image()
			ring := img.RGBAAt(int(close.X)-2, int(close.Y+close.H/2))
			if !sameColor(ring, pal.Focus, 12) {
				t.Fatalf("keyboard focus does not ring the close button: %v, want %v", ring, pal.Focus)
			}
			assertSingleX(t, img, close, img.RGBAAt(int(close.X+close.W/2), int(close.Y+2)), 5)
		})
	}
}

// assertSingleX checks that the ink in the close button is one X, centered:
// square around the button's center, every row inked, both ends of the top
// and bottom rows inked away from the middle, and the strokes crossing at the
// middle row. A horizontal bar is one row, and the old two-Line paint made
// one.
func assertSingleX(t *testing.T, img *image.RGBA, r ui.Rect, bg color.RGBA, pad int) {
	t.Helper()
	ink := func(x, y int) bool {
		c := img.RGBAAt(x, y)
		diff := abs(int(c.R)-int(bg.R)) + abs(int(c.G)-int(bg.G)) + abs(int(c.B)-int(bg.B))
		return diff > 120
	}
	minX, minY, maxX, maxY := 1<<30, 1<<30, -1, -1
	for y := int(r.Y) + pad; y < int(r.Y+r.H)-pad; y++ {
		for x := int(r.X) + pad; x < int(r.X+r.W)-pad; x++ {
			if ink(x, y) {
				minX, minY = min(minX, x), min(minY, y)
				maxX, maxY = max(maxX, x), max(maxY, y)
			}
		}
	}
	if maxX < 0 {
		t.Fatal("the revealed close button paints no X")
	}
	w, h := maxX-minX+1, maxY-minY+1
	if w < 5 || w > 11 || h < 5 || h > 11 || abs(w-h) > 1 {
		t.Fatalf("close glyph is %dx%d, want one square X about 7 DIP across", w, h)
	}
	// The glyph centers on the button within 1.5 DIP.
	wantX, wantY := r.X+r.W/2, r.Y+r.H/2
	gotX, gotY := float32(minX+maxX+1)/2, float32(minY+maxY+1)/2
	if abs32(gotX-wantX) > 1.5 || abs32(gotY-wantY) > 1.5 {
		t.Fatalf("close glyph center = (%v, %v), want (%v, %v) within 1.5 DIP", gotX, gotY, wantX, wantY)
	}
	// Every row of the glyph is inked: an X covers its whole height.
	for y := minY; y <= maxY; y++ {
		row := false
		for x := minX; x <= maxX; x++ {
			row = row || ink(x, y)
		}
		if !row {
			t.Fatalf("close glyph has an empty row at %d: not an X", y)
		}
	}
	// The top and bottom rows hold the two stroke ends, away from the middle;
	// the middle row holds their crossing.
	for _, y := range []int{minY, maxY} {
		left, right, mid := false, false, false
		for x := minX; x <= maxX; x++ {
			if !ink(x, y) {
				continue
			}
			switch {
			case x <= minX+w/4:
				left = true
			case x >= maxX-w/4:
				right = true
			default:
				mid = true
			}
		}
		if !left || !right || mid {
			t.Fatalf("close glyph row %d = ends %v/%v middle %v: not one X", y, left, right, mid)
		}
	}
	crossed := false
	for x := minX; x <= maxX; x++ {
		if ink(x, (minY+maxY)/2) && abs(x-int(wantX)) <= 2 {
			crossed = true
			break
		}
	}
	if !crossed {
		t.Fatal("close glyph does not cross at its middle: not an X")
	}
}

func sameColor(c color.RGBA, want ui.Color, tol int) bool {
	return abs(int(c.R)-int(want.R)) <= tol && abs(int(c.G)-int(want.G)) <= tol && abs(int(c.B)-int(want.B)) <= tol
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
