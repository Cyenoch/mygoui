package mygoui

import (
	"testing"

	"github.com/egoist/mygo/ui"
)

// forcedMask paints the mask shader as a GPU frame does, which the headless
// tester never reports.
type forcedMask struct {
	*scrollMask
	top, bottom float32
}

func (f forcedMask) PaintEffect(p *ui.Painter, box ui.Rect, content *ui.Bitmap, inset ui.Point) {
	f.paintMask(p, box, content, inset, f.top, 0, f.bottom, 0)
}

// The mask shader ramps content to transparent over each edge's length and
// leaves the rest untouched.
func TestScrollMaskShaderRampsEdges(t *testing.T) {
	content, background := ui.RGB(30, 70, 230), ui.RGB(230, 40, 70)
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Box(c).Fill().Background(background).Children(func() {
			fx := newScrollMask()
			ui.Box(c).Absolute().Left(0).Top(0).Size(40, 100).Background(content).
				Effect(forcedMask{&fx, 20, 20})
		})
	}, 40, 100)
	tt.Frame()
	img := tt.Image()
	at := func(y int) uint8 { return img.RGBAAt(20, y).B }
	if at(0) >= at(5) || at(5) >= at(10) || at(10) >= at(20) {
		t.Errorf("top ramp not increasing: %d %d %d %d", at(0), at(5), at(10), at(20))
	}
	if at(20) != content.B || at(50) != content.B || at(79) != content.B {
		t.Errorf("middle altered: %d %d %d", at(20), at(50), at(79))
	}
	if at(99) >= at(94) || at(94) >= at(89) {
		t.Errorf("bottom ramp not decreasing: %d %d %d", at(89), at(94), at(99))
	}
}
