package mygoui_test

import (
	"image"
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

var (
	maskBackgroundA = ui.RGB(230, 40, 70)
	maskBackgroundB = ui.RGB(20, 200, 80)
	maskContent     = ui.RGB(30, 70, 230)
)

func maskPixel(t *testing.T, img *image.RGBA, x, y int, background ui.Color, alpha float32) {
	t.Helper()
	want := maskContent.Alpha(alpha).Over(background)
	got := img.RGBAAt(x, y)
	difference := func(a, b uint8) int {
		return max(int(a)-int(b), int(b)-int(a))
	}
	if difference(got.R, want.R) > 9 || difference(got.G, want.G) > 9 || difference(got.B, want.B) > 9 {
		t.Errorf("pixel (%d,%d) = %v, want %v at alpha %.2f", x, y, got, want, alpha)
	}
}

func TestScrollMaskFadeTracksOverflowAndBackground(t *testing.T) {
	var state ui.ScrollState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Box(c).Fill().Children(func() {
			ui.Box(c).Absolute().Left(0).Right(0).Top(0).Bottom(0).
				Draw(func(p *ui.Painter, r ui.Rect) {
					for y := range 2 {
						for x := range 4 {
							color := maskBackgroundA
							if (x+y)%2 != 0 {
								color = maskBackgroundB
							}
							p.Fill(ui.Rect{X: r.X + float32(x*30), Y: r.Y + float32(y*50), W: 30, H: 50}, color, 0)
						}
					}
				})
			viewport := ui.ScrollBoth(c).Absolute().Left(0).Top(0).Size(120, 100).TrackScroll(&state)
			viewport.Children(func() {
				ui.Box(c).Size(320, 300).Shrink(0).Background(maskContent)
			})
			mygoui.ScrollFade(c, viewport, &state, mygoui.ScrollFadeOptions{Mode: mygoui.MaskFade})
		})
	}, 120, 100)
	tt.Frame()
	if state.MaxX < 48 || state.MaxY < 48 {
		t.Fatalf("missing two-axis overflow: %+v", state)
	}
	img := tt.Image()
	maskPixel(t, img, 0, 40, maskBackgroundA, 1)
	maskPixel(t, img, 40, 0, maskBackgroundB, 1)
	maskPixel(t, img, 119, 40, maskBackgroundB, 0.5/24)
	maskPixel(t, img, 40, 99, maskBackgroundA, 0.5/24)

	state.X, state.Y = 6, 8
	tt.Frame()
	img = tt.Image()
	maskPixel(t, img, 0, 40, maskBackgroundA, 0.5/6)
	maskPixel(t, img, 5, 40, maskBackgroundA, 5.5/6)
	maskPixel(t, img, 10, 40, maskBackgroundA, 1)
	maskPixel(t, img, 40, 0, maskBackgroundB, 0.5/8)
	maskPixel(t, img, 40, 7, maskBackgroundB, 7.5/8)
	maskPixel(t, img, 40, 12, maskBackgroundB, 1)
	maskPixel(t, img, 3, 3, maskBackgroundA, (3.5/6)*(3.5/8))

	state.X, state.Y = state.MaxX, state.MaxY
	tt.Frame()
	img = tt.Image()
	maskPixel(t, img, 119, 40, maskBackgroundB, 1)
	maskPixel(t, img, 40, 99, maskBackgroundA, 1)
	maskPixel(t, img, 0, 40, maskBackgroundA, 0.5/24)
	maskPixel(t, img, 40, 0, maskBackgroundB, 0.5/24)
	if tt.FramePending() {
		t.Fatal("stationary mask requested another frame")
	}
}

func TestScrollMaskFadeUsesPhysicalEdgesInRTL(t *testing.T) {
	var state ui.ScrollState
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Box(c).Fill().Background(maskBackgroundA).Children(func() {
			viewport := ui.ScrollHorizontal(c).Absolute().Left(0).Top(0).Size(120, 100).TrackScroll(&state)
			viewport.Children(func() {
				ui.Box(c).Size(320, 100).Shrink(0).Background(maskContent)
			})
			mygoui.ScrollFade(c, viewport, &state, mygoui.ScrollFadeOptions{Mode: mygoui.MaskFade})
		})
	}, 120, 100)
	tt.SetDirection(ui.RightToLeft)
	tt.Frame()
	if state.MaxX < 48 {
		t.Fatalf("missing horizontal overflow: %+v", state)
	}
	img := tt.Image()
	maskPixel(t, img, 0, 40, maskBackgroundA, 0.5/24)
	maskPixel(t, img, 119, 40, maskBackgroundA, 1)
	state.X = state.MaxX
	tt.Frame()
	img = tt.Image()
	maskPixel(t, img, 0, 40, maskBackgroundA, 1)
	maskPixel(t, img, 119, 40, maskBackgroundA, 0.5/24)
}

func TestScrollMaskFadeWithoutOverflowAndWithSelectedEdges(t *testing.T) {
	for _, overflow := range []bool{false, true} {
		t.Run(map[bool]string{false: "no overflow", true: "selected edge"}[overflow], func(t *testing.T) {
			var state ui.ScrollState
			tt := ui.NewTester(func(c *ui.Context) {
				ui.Box(c).Fill().Background(maskBackgroundA).Children(func() {
					viewport := ui.ScrollBoth(c).Absolute().Left(0).Top(0).Size(120, 100).TrackScroll(&state)
					viewport.Children(func() {
						size := float32(100)
						if overflow {
							size = 300
						}
						ui.Box(c).Size(size, size).Shrink(0).Background(maskContent)
					})
					mygoui.ScrollFade(c, viewport, &state, mygoui.ScrollFadeOptions{
						Mode: mygoui.MaskFade, Edges: mygoui.FadeBottom,
					})
				})
			}, 120, 100)
			if overflow {
				state.X, state.Y = 30, 30
				tt.Frame()
			}
			img := tt.Image()
			maskPixel(t, img, 0, 40, maskBackgroundA, 1)
			maskPixel(t, img, 40, 0, maskBackgroundA, 1)
			maskPixel(t, img, 40, 99, maskBackgroundA, map[bool]float32{false: 1, true: 0.5 / 24}[overflow])
			if tt.FramePending() {
				t.Fatal("stationary mask requested another frame")
			}
		})
	}
}

func TestScrollMaskFadeKeepsNativeInput(t *testing.T) {
	var state ui.ScrollState
	calls := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.ScrollArea(c, &state, mygoui.ScrollAreaOptions{
			Height: 100, Label: "Actions",
			Fade: &mygoui.ScrollFadeOptions{Mode: mygoui.MaskFade},
		}, func() {
			ui.Column(c).Children(func() {
				mygoui.Button(c, mygoui.ButtonOptions{Label: "Run", OnClick: func() { calls++ }})
				ui.Box(c).Height(200).Shrink(0)
			})
		})
	}, 220, 140)
	if err := tt.Click("Run"); err != nil {
		t.Fatal(err)
	}
	if calls != 1 || !tt.Focused("Run") {
		t.Fatalf("masked button calls=%d focused=%v", calls, tt.Focused("Run"))
	}
	tt.Key(0, ui.KeyEnter)
	if calls != 2 {
		t.Fatalf("keyboard activations=%d, want 2", calls)
	}
	tt.Scroll(180, 50, 0, 40)
	if state.Y != 40 {
		t.Fatalf("masked viewport scroll=%v, want 40", state.Y)
	}
}

func TestScrollMaskFadePreservesIndeterminateProgressLabels(t *testing.T) {
	for _, circular := range []bool{false, true} {
		name := "bar"
		if circular {
			name = "circle"
		}
		t.Run(name, func(t *testing.T) {
			var state ui.ScrollState
			tt := ui.NewTester(func(c *ui.Context) {
				mygoui.ScrollArea(c, &state, mygoui.ScrollAreaOptions{
					Height: 80, Fade: &mygoui.ScrollFadeOptions{Mode: mygoui.MaskFade},
				}, func() {
					if circular {
						mygoui.ProgressCircle(c, mygoui.ProgressCircleOptions{Indeterminate: true})
					} else {
						mygoui.ProgressBar(c, mygoui.ProgressOptions{Indeterminate: true, ShowOutput: true})
					}
				})
			}, 180, 80)
			if _, ok := tt.Find(locale.M().Loading); !ok {
				t.Fatalf("loading label is missing: %v", tt.Texts())
			}
			if tt.HasText("0%") {
				t.Fatalf("indeterminate progress showed a numeric value: %v", tt.Texts())
			}
		})
	}
}
