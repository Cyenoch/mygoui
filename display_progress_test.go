package mygoui_test

import (
	"image"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// inkCount counts the pixels painted in c, within the palette's sRGB value.
func inkCount(img *image.RGBA, c ui.Color) int {
	n := 0
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl := pixelAt(img, x, y)
			if near(r, c.R) && near(g, c.G) && near(bl, c.B) {
				n++
			}
		}
	}
	return n
}

// TestProgressValueChangeEasesTheFill pins the 500 ms width transition: a
// settled determinate bar asks for no frames, a new value starts the fill
// from where it was instead of jumping, and it stops on the new value.
func TestProgressValueChangeEasesTheFill(t *testing.T) {
	const width, height = 400, 80
	value := 20.0
	var pal theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		pal = theme.Colors(c)
		mygoui.ProgressBar(c, mygoui.ProgressOptions{Label: "Upload", Value: value, Max: 100})
	}, width, height)
	tt.Frame()
	if tt.FramePending() {
		t.Fatal("a determinate bar at rest asks for frames")
	}
	low := inkCount(tt.Image(), pal.Primary)
	want := width * 6 / 5 // a fifth of the 6 DIP track
	if d := low - want; d < -want/8 || d > want/8 {
		t.Fatalf("fill at 20%% = %d px, want about %d", low, want)
	}
	value = 80
	tt.Frame()
	if !tt.FramePending() {
		t.Fatal("a value change asks for no frames")
	}
	if mid := inkCount(tt.Image(), pal.Primary); mid > low+low/4 {
		t.Fatalf("fill jumped from %d px to %d px instead of easing", low, mid)
	}
	tt.Advance(400 * time.Millisecond)
	if !tt.FramePending() {
		t.Fatal("the 500 ms fill transition settled before 400 ms")
	}
	tt.Advance(200 * time.Millisecond)
	if tt.FramePending() {
		t.Fatal("the fill transition never settles")
	}
	high := inkCount(tt.Image(), pal.Primary)
	want = width * 6 * 4 / 5
	if d := high - want; d < -want/8 || d > want/8 {
		t.Fatalf("fill at 80%% = %d px, want about %d", high, want)
	}
}

func TestProgressAndMeterUseDistinctTrackGeometry(t *testing.T) {
	for _, tc := range []struct {
		name   string
		meter  bool
		height float32
	}{
		{name: "progress", height: 6},
		{name: "meter", meter: true, height: 8},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var palette theme.Palette
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				palette = theme.Colors(c)
				ui.Column(c).Fill().Padding(20).Background(palette.Background).Children(func() {
					if tc.meter {
						mygoui.Meter(c, mygoui.MeterOptions{Value: 50}).Label("Indicator")
					} else {
						mygoui.ProgressBar(c, mygoui.ProgressOptions{Value: 50}).Label("Indicator")
					}
				})
			}, 160, 80)
			track := findRect(t, tt, "Indicator")
			if track.H != tc.height {
				t.Fatalf("track height = %v, want %v DIP", track.H, tc.height)
			}
			img := tt.Image()
			wantPixel(t, img, int(track.X)+10, int(track.Y+track.H/2), palette.Primary, "default indicator")
			wantPixel(t, img, int(track.X+track.W)-10, int(track.Y+track.H/2),
				palette.Input.Over(palette.Background), "unfilled track")
			corner := pixel(img, track.X, track.Y)
			if tc.meter {
				wantPixel(t, img, int(track.X), int(track.Y), palette.Primary, "square meter corner")
			} else if same(corner, palette.Primary) {
				t.Fatalf("pill progress corner is filled square: %v", corner)
			}
		})
	}
}

// TestReducedMotionJumpsTheFill proves a value change lands at once while the
// desktop asks for less motion, and that the indeterminate sweep stops.
func TestReducedMotionJumpsTheFill(t *testing.T) {
	value := 20.0
	var pal theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		pal = theme.Colors(c)
		mygoui.ProgressBar(c, mygoui.ProgressOptions{Value: value, Max: 100})
		mygoui.ProgressBar(c, mygoui.ProgressOptions{
			Indeterminate: true, Color: mygoui.IndicatorDanger,
		})
	}, 400, 120)
	tt.SetPreferences(ui.Preferences{ReduceMotion: true})
	tt.Frame()
	if tt.FramePending() {
		t.Fatal("reduced motion still animates")
	}
	low := inkCount(tt.Image(), pal.Primary)
	value = 80
	tt.Frame()
	if mid := inkCount(tt.Image(), pal.Primary); mid <= low+low/4 {
		t.Fatalf("fill stayed at %d px, want it at once near 80%% of the track", mid)
	}
	if tt.FramePending() {
		t.Fatal("the fill still animates without motion")
	}
}

// TestProgressIndeterminateKeepsPainting pins the looping band: the sweep
// asks for paint frames only, and keeps asking while it runs.
func TestProgressIndeterminateKeepsPainting(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.ProgressBar(c, mygoui.ProgressOptions{
			Label: "Indexing", Indeterminate: true, Size: mygoui.ProgressLarge,
		})
	}, 320, 80)
	tt.Frame()
	if !tt.FramePending() {
		t.Fatal("the indeterminate band does not animate")
	}
	tt.Advance(300 * time.Millisecond)
	if !tt.FramePending() {
		t.Fatal("the indeterminate band stopped on its own")
	}
}

func TestProgressIndeterminateStaysInsideTrack(t *testing.T) {
	var pal theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		pal = theme.Colors(c)
		ui.Column(c).Fill().Padding(40).Background(pal.Background).Children(func() {
			mygoui.ProgressBar(c, mygoui.ProgressOptions{Indeterminate: true}).Label("Track")
		})
	}, 300, 100)
	r := findRect(t, tt, "Track")
	for frame := 0; frame < 9; frame++ {
		tt.Advance(200 * time.Millisecond)
		img := tt.Image()
		for _, x := range []int{15, 30, int(r.X) - 1, int(r.X + r.W), 280} {
			wantPixel(t, img, x, int(r.Y+r.H/2), pal.Background, "outside progress track")
		}
	}
}

// TestProgressClampsAtItsBounds covers the empty, full and inverted ranges
// that the output text reports.
func TestProgressClampsAtItsBounds(t *testing.T) {
	cases := []struct {
		name string
		bar  mygoui.ProgressOptions
		want string
	}{
		{"zero", mygoui.ProgressOptions{Value: 0, ShowOutput: true}, "0%"},
		{"default midpoint", mygoui.ProgressOptions{Value: 60, ShowOutput: true}, "60%"},
		{"max", mygoui.ProgressOptions{Value: 100, ShowOutput: true}, "100%"},
		{"over max", mygoui.ProgressOptions{Value: 140, ShowOutput: true}, "100%"},
		{"under zero", mygoui.ProgressOptions{Value: -20, ShowOutput: true}, "0%"},
		{
			"inverted bounds",
			mygoui.ProgressOptions{Value: 5, Min: 10, Max: 10, ShowOutput: true}, "0%",
		},
		{"a third", mygoui.ProgressOptions{Value: 1, Max: 3, ShowOutput: true}, "33%"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				mygoui.ProgressBar(c, tc.bar)
			}, 320, 80)
			if !tt.HasText(tc.want) {
				t.Fatalf("output = %v, want %q", tt.Texts(), tc.want)
			}
		})
	}

	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.Meter(c, mygoui.MeterOptions{
			Label: "Seats", Value: 8, Max: 12, ShowOutput: true,
		})
	}, 320, 80)
	if !tt.HasText("67%") {
		t.Fatalf("meter output = %v, want 67%%", tt.Texts())
	}
}

// TestProgressCircleEasesAndCoversItsBounds pins the ring: no arc at zero, a
// full arc at the maximum, and a value change that paints from where it was.
func TestProgressCircleEasesAndCoversItsBounds(t *testing.T) {
	value := 0.0
	var pal theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		pal = theme.Colors(c)
		mygoui.ProgressCircle(c, mygoui.ProgressCircleOptions{Value: value, Max: 100})
	}, 80, 80)
	tt.Frame()
	track := pal.Default.Over(pal.Background)
	if n := arcPixels(tt.Image(), pal.Primary, track); n != 0 {
		t.Fatalf("an empty ring painted %d primary pixels", n)
	}
	value = 100
	tt.Frame()
	if n := arcPixels(tt.Image(), pal.Primary, track); n != 0 {
		t.Fatalf("the ring jumped to %d primary pixels", n)
	}
	tt.Advance(600 * time.Millisecond)
	full := arcPixels(tt.Image(), pal.Primary, track)
	if full == 0 {
		t.Fatal("a full ring paints no primary arc")
	}
	if tt.FramePending() {
		t.Fatal("the ring never settles")
	}
	// The arc covers the ring at its maximum: a quarter of the value paints
	// about a quarter of the arc, so the full ring measures far more.
	if quarter := arcPixels(ringFrame(t, 25), pal.Primary, track); full <= 2*quarter {
		t.Fatalf("full arc = %d px, quarter arc = %d px; want the full arc to cover the ring", full, quarter)
	}
	value = 0
	tt.Frame()
	if n := arcPixels(tt.Image(), pal.Primary, track); n != full {
		t.Fatalf("the ring emptied at once (%d of %d pixels)", n, full)
	}
	tt.Advance(600 * time.Millisecond)
	if n := arcPixels(tt.Image(), pal.Primary, track); n != 0 {
		t.Fatalf("the ring kept %d primary pixels at zero", n)
	}
}

// ringFrame paints one determinate ring at value and returns its frame.
func ringFrame(t *testing.T, value float64) *image.RGBA {
	t.Helper()
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.ProgressCircle(c, mygoui.ProgressCircleOptions{Value: value, Max: 100})
	}, 80, 80)
	return tt.Image()
}

// arcPixels counts the pixels painted on the fill's side of the fill/track
// boundary: the ring's arc, counted through the antialiasing of its dots,
// which leaves no pixel exactly at the fill color.
func arcPixels(img *image.RGBA, fill, track ui.Color) int {
	n := 0
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl := pixelAt(img, x, y)
			c := ui.Color{R: r, G: g, B: bl}
			if colorDistance(c, fill) < colorDistance(c, track) {
				n++
			}
		}
	}
	return n
}

// colorDistance is the squared sRGB distance between two colors.
func colorDistance(a, b ui.Color) int {
	dr, dg, db := int(a.R)-int(b.R), int(a.G)-int(b.G), int(a.B)-int(b.B)
	return dr*dr + dg*dg + db*db
}
