package mygoui_test

import (
	"image/color"
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/egoist/mygo/ui"
)

func TestScrollEdgesFollowCurrentLayoutAndDistance(t *testing.T) {
	var scroll ui.ScrollState
	short, disabled := false, false
	white, black := ui.RGB(255, 255, 255), ui.RGB(0, 0, 0)
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Fill().Background(white).Children(func() {
			ui.Box(c).FillWidth().Height(32).Background(black).Label("Header")
			viewport := ui.Scroll(c).FillWidth().Height(100).TrackScroll(&scroll).Label("Viewport")
			mygoui.ScrollEdges(c, viewport, &scroll, mygoui.ScrollEdgesOptions{Surface: white, Depth: 24, Disabled: disabled})
			viewport.Children(func() {
				height := float32(300)
				if short {
					height = 80
				}
				ui.Box(c).FillWidth().Height(height).Shrink(0).Background(black)
			})
		})
	}, 240, 180)
	pixel := func(x, y int) color.RGBA { return tt.Image().RGBAAt(x, y) }
	if scroll.MaxY <= 0 {
		t.Fatal("fixture did not overflow")
	}
	if top, bottom := pixel(40, 33), pixel(40, 130); top.R != 0 || bottom.R < 220 {
		t.Fatalf("first frame must fade only the hidden-content edge: top=%v bottom=%v", top, bottom)
	}
	if got := pixel(40, 31); got.R != 0 {
		t.Fatalf("fade covered the header: %v", got)
	}
	if got := pixel(238, 130); got.R != 255 {
		t.Fatalf("content painted in native scrollbar lane: %v", got)
	}
	scroll.Y = 2
	tt.Frame()
	near := pixel(40, 33).R
	if near == 0 || near > 40 {
		t.Fatalf("near-start strength was not proportional: %d", near)
	}
	scroll.Y = 30
	tt.Frame()
	if full := pixel(40, 33).R; full < 220 || full <= near {
		t.Fatalf("scrolled top did not reach full fade: near=%d full=%d", near, full)
	}
	scroll.Y = scroll.MaxY
	tt.Frame()
	if got := pixel(40, 130); got.R != 0 {
		t.Fatalf("bottom fade remained at the end: %v", got)
	}
	short = true
	tt.Frame()
	if scroll.MaxY != 0 || pixel(40, 33).R != 0 {
		t.Fatalf("resize/shrink kept stale fade: maxY=%g top=%v", scroll.MaxY, pixel(40, 33))
	}
	short = false
	disabled = true
	tt.Frame()
	if got := pixel(40, 130); got.R != 0 {
		t.Fatalf("disabled effect painted: %v", got)
	}
}

func TestScrollEdgesDoNotInterceptInputOrSoftenHighContrast(t *testing.T) {
	var scroll ui.ScrollState
	clicks := 0
	white, black := ui.RGB(255, 255, 255), ui.RGB(0, 0, 0)
	tt := ui.NewTester(func(c *ui.Context) {
		viewport := ui.Scroll(c).FillWidth().Height(100).Background(white).TrackScroll(&scroll)
		mygoui.ScrollEdges(c, viewport, &scroll, mygoui.ScrollEdgesOptions{Surface: white})
		viewport.Children(func() {
			ui.ButtonBase(c).Label("Edge action").FillWidth().Height(32).Shrink(0).Background(black).
				OnClick(func() { clicks++ })
			ui.Box(c).FillWidth().Height(300).Shrink(0).Background(black)
		})
	}, 240, 160)
	if err := tt.Click("Edge action"); err != nil {
		t.Fatal(err)
	}
	if clicks != 1 {
		t.Fatal("paint decoration intercepted activation")
	}
	tt.SetPreferences(ui.Preferences{HighContrast: true, TextScale: 1, ReduceMotion: true})
	tt.Frame()
	if got := tt.Image().RGBAAt(40, 98); got.R != 0 {
		t.Fatalf("high contrast retained optional fade: %v", got)
	}
}
