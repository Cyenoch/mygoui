package mygoui_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

var (
	maskBackgroundA = ui.RGB(230, 40, 70)
	maskContent     = ui.RGB(30, 70, 230)
)

// maskContentAt asserts the viewport content shows at full strength.
func maskContentAt(t *testing.T, img *image.RGBA, x, y int) {
	t.Helper()
	if got := img.RGBAAt(x, y); got != (color.RGBA{R: maskContent.R, G: maskContent.G, B: maskContent.B, A: 255}) {
		t.Errorf("pixel (%d,%d) = %v, want the unfaded content", x, y, got)
	}
}

// maskShadedAt asserts an edge is dimmed but the content stays readable.
func maskShadedAt(t *testing.T, img *image.RGBA, x, y int) {
	t.Helper()
	got := img.RGBAAt(x, y)
	if got.R == maskContent.R && got.G == maskContent.G && got.B == maskContent.B {
		t.Errorf("pixel (%d,%d) = %v, want a visible edge fade", x, y, got)
	}
	if got.A != 255 || got.B < maskContent.B/2 {
		t.Errorf("pixel (%d,%d) = %v, want the content still readable", x, y, got)
	}
}

// Without a GPU the tester's shader would run per pixel on the CPU: the mask
// falls back to the content under a shadow, kept only where a side has more
// to scroll.
func TestScrollMaskFadeWithoutGPUShadesOnlyScrollableEdges(t *testing.T) {
	var state ui.ScrollState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Box(c).Fill().Background(maskBackgroundA).Children(func() {
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
	maskContentAt(t, img, 0, 40)
	maskContentAt(t, img, 40, 0)
	maskShadedAt(t, img, 119, 40)
	maskShadedAt(t, img, 40, 99)

	state.X, state.Y = 6, 8
	tt.Frame()
	img = tt.Image()
	maskShadedAt(t, img, 0, 40)
	maskShadedAt(t, img, 40, 0)
	maskContentAt(t, img, 60, 50)

	state.X, state.Y = state.MaxX, state.MaxY
	tt.Frame()
	img = tt.Image()
	maskContentAt(t, img, 119, 40)
	maskContentAt(t, img, 40, 99)
	maskShadedAt(t, img, 0, 40)
	maskShadedAt(t, img, 40, 0)
	if tt.FramePending() {
		t.Fatal("stationary mask requested another frame")
	}
}

func TestScrollMaskFadeUsesPhysicalEdgesInRTL(t *testing.T) {
	var state ui.ScrollState
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
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
	maskShadedAt(t, img, 0, 40)
	maskContentAt(t, img, 119, 40)
	state.X = state.MaxX
	tt.Frame()
	img = tt.Image()
	maskContentAt(t, img, 0, 40)
	maskShadedAt(t, img, 119, 40)
}

func TestScrollMaskFadeWithoutOverflowAndWithSelectedEdges(t *testing.T) {
	for _, overflow := range []bool{false, true} {
		t.Run(map[bool]string{false: "no overflow", true: "selected edge"}[overflow], func(t *testing.T) {
			var state ui.ScrollState
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
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
			maskContentAt(t, img, 0, 40)
			maskContentAt(t, img, 40, 0)
			if overflow {
				maskShadedAt(t, img, 40, 99)
			} else {
				maskContentAt(t, img, 40, 99)
			}
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
