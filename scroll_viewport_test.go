package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/egoist/mygo/ui"
)

func TestScrollEdgesReserveBothNativeBarLanesAndFadeToClip(t *testing.T) {
	for _, direction := range []ui.Direction{ui.LeftToRight, ui.RightToLeft} {
		t.Run(direction.String(), func(t *testing.T) {
			var scroll ui.ScrollState
			short, disabled := false, false
			white, black := ui.RGB(255, 255, 255), ui.RGB(0, 0, 0)
			tt := ui.NewTester(func(c *ui.Context) {
				theme := *c.Theme()
				theme.ScrollbarWidth = 8
				c.SetTheme(&theme)
				ui.Column(c).Fill().Background(white).Children(func() {
					ui.Box(c).Height(24).FillWidth().Background(black).Label("Viewport header")
					viewport := ui.ScrollBoth(c).Size(240, 120).TrackScroll(&scroll).Label("Both axes")
					mygoui.ScrollEdges(c, viewport, &scroll, mygoui.ScrollEdgesOptions{
						Surface: white, Depth: 24, Axes: ui.Horizontal | ui.Vertical, Disabled: disabled,
					})
					viewport.Children(func() {
						w, h := float32(600), float32(420)
						if short {
							w, h = 100, 60
						}
						ui.Box(c).Size(w, h).Shrink(0).Background(black)
					})
					ui.Box(c).Height(24).FillWidth().Background(black).Label("Viewport footer")
				})
			}, 260, 200)
			tt.SetDirection(direction)
			tt.Frame()
			if scroll.MaxX <= 0 || scroll.MaxY <= 0 {
				t.Fatalf("both-axis overflow missing: %+v", scroll)
			}
			viewport, _ := tt.Find("Both axes")
			left, right := int(viewport.X), int(viewport.X+viewport.W)
			if direction == ui.RightToLeft {
				left += 14
			} else {
				right -= 14
			}
			top, bottom := int(viewport.Y), int(viewport.Y+viewport.H)-14
			pixel := func(x, y int) uint8 { return tt.Image().RGBAAt(x, y).R }
			assertLanes := func() {
				t.Helper()
				for y := bottom; y < int(viewport.Y+viewport.H); y++ {
					for x := int(viewport.X); x < int(viewport.X+viewport.W); x++ {
						if pixel(x, y) != 255 {
							t.Fatalf("content reappeared in horizontal bar lane at %d,%d", x, y)
						}
					}
				}
				x := right + 2
				if direction == ui.RightToLeft {
					x = left - 2
				}
				for y := top; y < bottom; y++ {
					if pixel(x, y) != 255 {
						t.Fatalf("content entered vertical bar lane at %d,%d", x, y)
					}
				}
			}
			assertLanes()
			if pixel((left+right)/2, bottom-1) < 240 || pixel((left+right)/2, top) != 0 {
				t.Fatal("first-frame vertical fade does not end at the true content clip")
			}
			startX, endX := left, right-1
			if direction == ui.RightToLeft {
				startX, endX = endX, startX
			}
			if pixel(startX, top+50) != 0 || pixel(endX, top+50) < 240 {
				t.Fatal("horizontal fade does not follow the logical start/end")
			}
			scroll.X, scroll.Y = 2, 2
			tt.Frame()
			if near := pixel((left+right)/2, top); near == 0 || near > 40 {
				t.Fatalf("near-start vertical fade is not proportional: %d", near)
			}
			if near := pixel(startX, top+50); near == 0 || near > 40 {
				t.Fatalf("near-start horizontal fade is not proportional: %d", near)
			}
			scroll.X, scroll.Y = scroll.MaxX/2, scroll.MaxY/2
			tt.Frame()
			assertLanes()
			if pixel((left+right)/2, top) < 240 || pixel(startX, top+50) < 240 {
				t.Fatal("middle offsets did not fade both hidden starts")
			}
			scroll.X, scroll.Y = scroll.MaxX, scroll.MaxY
			tt.Frame()
			assertLanes()
			if pixel((left+right)/2, bottom-1) != 0 || pixel(endX, top+50) != 0 {
				t.Fatal("end fade hid fully revealed content")
			}
			short = true
			tt.Frame()
			if scroll.MaxX != 0 || scroll.MaxY != 0 || scroll.X != 0 || scroll.Y != 0 {
				t.Fatalf("content shrink retained stale extents or offsets: %+v", scroll)
			}
			assertLanes()
			short, disabled = false, true
			tt.Frame()
			if pixel((left+right)/2, bottom-1) != 0 {
				t.Fatal("disabled fade remained visible")
			}
			assertLanes()
			if pixel(40, int(viewport.Y)-1) != 0 || pixel(40, int(viewport.Y+viewport.H)) != 0 {
				t.Fatal("effect painted outside the content viewport")
			}
		})
	}
}

func TestScrollEdgesKeepNativeScrollbarDraggingAndTinyClip(t *testing.T) {
	var scroll ui.ScrollState
	height := float32(120)
	tt := ui.NewTester(func(c *ui.Context) {
		viewport := ui.ScrollBoth(c).Width(240).Height(height).TrackScroll(&scroll).Label("Drag viewport")
		mygoui.ScrollEdges(c, viewport, &scroll, mygoui.ScrollEdgesOptions{
			Surface: ui.RGB(255, 255, 255), Axes: ui.Horizontal | ui.Vertical,
		})
		viewport.Children(func() { ui.Box(c).Size(600, 420).Shrink(0).Background(ui.RGB(0, 0, 0)) })
	}, 260, 160)
	tt.Press(20, 114)
	tt.Move(180, 114)
	tt.Release(180, 114)
	if scroll.X <= 0 {
		t.Fatal("physical horizontal gutter blocked native scrollbar dragging")
	}
	tt.Press(234, 16)
	tt.Move(234, 90)
	tt.Release(234, 90)
	if scroll.Y <= 0 {
		t.Fatal("physical vertical gutter blocked native scrollbar dragging")
	}
	tt.Move(259, 159)
	height = 18
	tt.Frame()
	// The six-DIP content clip shares its height between the two fades;
	// neither gradient is allowed to reach the physical scrollbar gutter.
	for y := 6; y < 18; y++ {
		if got := tt.Image().RGBAAt(40, y); got.R != 255 || got.G != 255 || got.B != 255 {
			t.Fatalf("tiny viewport painted content/effect in the bar lane at y=%d: %v", y, got)
		}
	}
	tt.SetPreferences(ui.Preferences{HighContrast: true, TextScale: 1, ReduceMotion: true})
	tt.Frame()
	if got := tt.Image().RGBAAt(40, 1); got.R != 0 {
		t.Fatalf("high contrast retained a fade: %v", got)
	}
}
