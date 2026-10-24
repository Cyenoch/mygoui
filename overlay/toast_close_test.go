package overlay

import (
	"image"
	"image/color"
	"testing"
	"time"

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
			// The close button sits in flow at the card's top-end, aligned
			// with the title's line box, and inside the card column.
			if close.Y != title.Y || close.H != title.H {
				t.Fatalf("close %+v is not aligned with the title row %+v", close, title)
			}
			cardW := float32(min(460, winW-32))
			cardX := (winW - cardW) / 2
			if close.X < cardX || close.X+close.W > cardX+cardW || close.Y+close.H > winH-theme.Space16 {
				t.Fatalf("close %+v is not within the card column [%v, %v]", close, cardX, cardX+cardW)
			}

			// At rest the button paints nothing: no face, no border, no bar.
			img := tt.Image()
			assertCloseSurface(t, img, close, pal.Surface, "at rest")
			if c := img.RGBAAt(int(close.X)-2, int(close.Y+close.H/2)); !sameColor(c, pal.Surface, 2) {
				t.Fatalf("close paints a ring at rest: %v at %v", c, close.X-2)
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
			img = tt.Image()
			if c := img.RGBAAt(int(close.X)-2, int(close.Y+close.H/2)); !sameColor(c, pal.Surface, 2) {
				t.Fatalf("hover paints outside the close button: %v at %v", c, close.X-2)
			}
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

			if err := tt.Click("Close"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
			if len(state.Items()) != 0 || tt.HasText("File saved") {
				t.Fatalf("close left %d toasts", len(state.Items()))
			}
		})
	}
}

// assertCloseSurface fails when any pixel of r paints over the card surface.
func assertCloseSurface(t *testing.T, img *image.RGBA, r ui.Rect, surface ui.Color, when string) {
	t.Helper()
	for y := int(r.Y); y < int(r.Y+r.H); y++ {
		for x := int(r.X); x < int(r.X+r.W); x++ {
			if c := img.RGBAAt(x, y); !sameColor(c, surface, 2) {
				t.Fatalf("%s: close button paints %v at %d,%d", when, c, x, y)
			}
		}
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

// TestToastMotionSettlesIn250ms checks the toast card's own motion: MyGo's
// ElementTransition runs the entry and the exit for 250 ms, and reduced
// motion removes it (TestToastReducedMotionTimeoutAndClose checks that).
//
// The transition asks the card to slide 16 DIP, scale from 0.96 and fade
// from transparent. MyGo 0.6.3 lays an absolute overlay out again after its
// transitions (layoutTree's layoutAbsolute), which puts the cards back at
// their final boxes: the card's scale and fade show, but its lift is re-laid
// away and the stack position snaps when a toast appears or goes, rather
// than sliding. The test asserts the settle time, the scale it does show and
// the fade: the first frame is the page background and the settled card is
// its surface.
func TestToastMotionSettlesIn250ms(t *testing.T) {
	var state ToastState
	var ctx *ui.Context
	tt := ui.NewTester(toastTestView(&state, &ctx), 640, 480)
	tt.SetPreferences(ui.Preferences{TextScale: 1})
	tt.Frame()

	if err := tt.Click("Show"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !tt.FramePending() {
		t.Fatal("the toast entry is not animated")
	}
	entering, ok := tt.Find("File saved")
	if !ok {
		t.Fatal("the entering toast shows no title")
	}
	// Sample the card's right padding, inside the 0.96 scale but clear of
	// the title, action and revealed close face. At the first frame the
	// entering card is fully transparent, so the page background shows.
	pal := theme.ColorsForMode(false)
	cardRight := float32((640 + min(460, 640-32)) / 2)
	sample := image.Pt(int(cardRight)-20, int(entering.Y+entering.H/2))
	if c := tt.Image().RGBAAt(sample.X, sample.Y); !sameColor(c, pal.Background, 2) {
		t.Fatalf("the entering toast paints %v at %v, want the transparent page background %v", c, sample, pal.Background)
	}
	tt.Advance(125 * time.Millisecond)
	tt.Frame()
	middle, _ := tt.Find("File saved")
	tt.Advance(200 * time.Millisecond)
	tt.Frame()
	settled, _ := tt.Find("File saved")
	if tt.FramePending() {
		t.Fatal("the toast entry still animates after 250 ms")
	}
	if !(entering.W < middle.W && middle.W <= settled.W) {
		t.Fatalf("entry does not scale into place: %v -> %v -> %v", entering.W, middle.W, settled.W)
	}
	if c := tt.Image().RGBAAt(sample.X, sample.Y); !sameColor(c, pal.Surface, 2) {
		t.Fatalf("the settled toast paints %v at %v, want its surface %v", c, sample, pal.Surface)
	}

	// Exit: dismissing runs the card's exit transition, which fades and
	// scales the retained copy for 250 ms.
	closeRect, ok := tt.Find("Close")
	if !ok {
		t.Fatal("no close button")
	}
	tt.Move(closeRect.X+closeRect.W/2, closeRect.Y+closeRect.H/2)
	tt.Frame()
	if err := tt.Click("Close"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !tt.FramePending() {
		t.Fatal("the toast exit is not animated")
	}
	// The exit copy is painted, not an Element that already expired: the
	// card's surface is still visible on its first exit frame.
	if c := tt.Image().RGBAAt(sample.X, sample.Y); sameColor(c, pal.Background, 2) {
		t.Fatalf("the leaving toast paints only the background at %v: its retained copy is missing", sample)
	}
	tt.Advance(300 * time.Millisecond)
	tt.Frame()
	if tt.FramePending() {
		t.Fatal("the toast exit still animates after 250 ms")
	}
	if len(state.Items()) != 0 {
		t.Fatalf("exit left %d retained toasts", len(state.Items()))
	}
	if c := tt.Image().RGBAAt(sample.X, sample.Y); !sameColor(c, pal.Background, 2) {
		t.Fatalf("the dismissed toast still paints %v at %v", c, sample)
	}
}
