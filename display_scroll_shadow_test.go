package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

var (
	fadeSurface = ui.RGB(250, 250, 250)
	fadeContent = ui.RGB(20, 120, 220)
)

// closeTo reports whether got is within tol of want on every channel.
func closeTo(got, want ui.Color, tol int) bool {
	d := func(a, b uint8) int {
		if a > b {
			return int(a - b)
		}
		return int(b - a)
	}
	return d(got.R, want.R) <= tol && d(got.G, want.G) <= tol && d(got.B, want.B) <= tol
}

// TestScrollShadowFadesTheScrolledEdges pins the edge fades against the real
// scroll position: the edge holding more content fades to the surface, the
// edge at the start or end of the content does not, and a partial scroll
// paints a partial fade.
func TestScrollShadowFadesTheScrolledEdges(t *testing.T) {
	var st ui.ScrollState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.ScrollShadow(c, &st, mygoui.ScrollShadowOptions{
			Orientation: mygoui.ScrollShadowVertical,
			Surface:     fadeSurface,
			Height:      100,
			Label:       "Activity",
		}, func() {
			ui.Column(c).Gap(0).Children(func() {
				for range 20 {
					ui.Box(c).FillWidth().Height(40).Background(fadeContent)
				}
			})
		})
	}, 220, 140)
	tt.Frame()
	if st.MaxY <= 0 {
		t.Fatalf("content does not overflow (MaxY %v)", st.MaxY)
	}
	if tt.FramePending() {
		t.Fatal("a scroll shadow at rest asks for frames")
	}
	img := tt.Image()
	if top := pixel(img, 20, 0); !same(top, fadeContent) {
		t.Fatalf("top edge at the start = %v, want the content %v", top, fadeContent)
	}
	if bottom := pixel(img, 20, 99); !closeTo(bottom, fadeSurface, 6) {
		t.Fatalf("bottom edge at the start = %v, want the surface %v", bottom, fadeSurface)
	}

	st.Y = st.MaxY
	tt.Frame()
	img = tt.Image()
	if top := pixel(img, 20, 0); !closeTo(top, fadeSurface, 6) {
		t.Fatalf("top edge at the end = %v, want the surface %v", top, fadeSurface)
	}
	if bottom := pixel(img, 20, 99); !same(bottom, fadeContent) {
		t.Fatalf("bottom edge at the end = %v, want the content %v", bottom, fadeContent)
	}

	st.Y = 10
	tt.Frame()
	img = tt.Image()
	top := pixel(img, 20, 0)
	if same(top, fadeContent) || closeTo(top, fadeSurface, 6) {
		t.Fatalf("top edge after 10 DIPs = %v, want a partial fade", top)
	}
}

// TestScrollShadowHorizontalFadesBothInlineEdges is the horizontal
// counterpart: the fade follows the inline scroll position.
func TestScrollShadowHorizontalFadesBothInlineEdges(t *testing.T) {
	var st ui.ScrollState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.ScrollShadow(c, &st, mygoui.ScrollShadowOptions{
			Orientation: mygoui.ScrollShadowHorizontal,
			Surface:     fadeSurface,
			Height:      60,
			Label:       "Cards",
		}, func() {
			ui.Row(c).Gap(0).Children(func() {
				for range 12 {
					ui.Box(c).Width(40).Height(60).Background(fadeContent)
				}
			})
		})
	}, 220, 140)
	tt.Frame()
	if st.MaxX <= 0 {
		t.Fatalf("content does not overflow (MaxX %v)", st.MaxX)
	}
	img := tt.Image()
	if left := pixel(img, 0, 20); !same(left, fadeContent) {
		t.Fatalf("left edge at the start = %v, want the content %v", left, fadeContent)
	}
	if right := pixel(img, 219, 20); !closeTo(right, fadeSurface, 6) {
		t.Fatalf("right edge at the start = %v, want the surface %v", right, fadeSurface)
	}

	st.X = st.MaxX
	tt.Frame()
	img = tt.Image()
	if left := pixel(img, 0, 20); !closeTo(left, fadeSurface, 6) {
		t.Fatalf("left edge at the end = %v, want the surface %v", left, fadeSurface)
	}
	if right := pixel(img, 219, 20); !same(right, fadeContent) {
		t.Fatalf("right edge at the end = %v, want the content %v", right, fadeContent)
	}
}

// TestScrollShadowWithoutOverflowPaintsNoFade keeps a short body clean: no
// scroll range, no edge fade.
func TestScrollShadowWithoutOverflowPaintsNoFade(t *testing.T) {
	var st ui.ScrollState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.ScrollShadow(c, &st, mygoui.ScrollShadowOptions{
			Surface: fadeSurface, Height: 100, Label: "Short",
		}, func() {
			ui.Box(c).FillWidth().Height(40).Background(fadeContent)
		})
	}, 220, 140)
	tt.Frame()
	if st.MaxY != 0 {
		t.Fatalf("a short body scrolls by %v", st.MaxY)
	}
	img := tt.Image()
	if top := pixel(img, 20, 0); !same(top, fadeContent) {
		t.Fatalf("top edge = %v, want the content %v", top, fadeContent)
	}
	if bottom := pixel(img, 20, 99); !closeTo(bottom, fadeSurface, 2) {
		t.Fatalf("bottom edge = %v, want the surface %v", bottom, fadeSurface)
	}
}
