package style_test

import (
	"image/color"
	"testing"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestHeadingTrackingTightensTheLine pins tracking-tight on the heading roles:
// the same text set as an H2 measures narrower than the same 30/36 role
// without the token, while body copy keeps the font's own spacing.
func TestHeadingTrackingTightensTheLine(t *testing.T) {
	const text = "Heading"
	var tracked, plain ui.Handle
	var frames int
	var trackedW, plainW float32
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).AlignItems(ui.Start).Gap(theme.Space8).Children(func() {
			style.Text(ui.Text(c, text), c, theme.H2()).Bind(&tracked)
			loose := theme.H2()
			loose.Tracking = 0
			style.Text(ui.Text(c, text), c, loose).Bind(&plain)
		})
		if frames > 0 {
			trackedW, plainW = tracked.Bounds(c).W, plain.Bounds(c).W
		}
		frames++
	}, 400, 160)
	tt.Frame()
	tt.Frame()

	if trackedW <= 0 || plainW <= 0 {
		t.Fatalf("text bounds are empty (%v, %v)", trackedW, plainW)
	}
	if trackedW >= plainW {
		t.Fatalf("tracked heading = %v wide, untracked = %v; want the heading narrower", trackedW, plainW)
	}
}

// TestDisabledControlsKeepTheirRestingFace pins the shared helper against
// MyGo's disabled contract: Hovered and Pressed report false for a disabled
// element, as upstream's removed pointer events do, so a disabled row keeps
// its resting face under the pointer. Disabled comes before Background, as
// every component applies it, because Background snapshots the face for the
// frame. The enabled row proves the sampler sees hover when it should.
func TestDisabledControlsKeepTheirRestingFace(t *testing.T) {
	for _, mode := range []struct {
		name string
		dark bool
	}{
		{name: "light"},
		{name: "dark", dark: true},
	} {
		t.Run(mode.name, func(t *testing.T) {
			testDisabledControlsKeepTheirRestingFace(t, mode.dark)
		})
	}
}

func testDisabledControlsKeepTheirRestingFace(t *testing.T, dark bool) {
	t.Helper()
	var pal theme.Palette
	var hover ui.Color
	var idle, live ui.Handle
	var idleRect, liveRect ui.Rect
	frames := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, dark)
		pal = theme.Colors(c)
		// A distinct translucent hover face makes a broken state change
		// detectable even when the default accent is only four percent ink.
		hover = pal.Primary.Alpha(0.25)
		ui.Column(c).FillWidth().FillHeight().Padding(24).Gap(12).
			Background(pal.Background).Children(func() {
			off := ui.ButtonBase(c).FillWidth().MinHeight(40).Label("off")
			off = style.Disabled(off, true)
			off = style.Background(off, 8, pal.Accent, hover, hover)
			off.Bind(&idle)
			on := ui.ButtonBase(c).FillWidth().MinHeight(40).Label("on")
			on = style.Background(on, 8, pal.Accent, hover, hover)
			on.Bind(&live)
		})
		if frames > 0 {
			idleRect, liveRect = idle.Bounds(c), live.Bounds(c)
		}
		frames++
	}, 320, 180)
	tt.Frame()
	tt.Frame()
	if idleRect.W <= 0 || liveRect.W <= 0 {
		t.Fatalf("rows are empty (%+v, %+v)", idleRect, liveRect)
	}

	offX, offY := int(idleRect.X)+8, int(idleRect.Y+idleRect.H/2)
	onX, onY := int(liveRect.X)+8, int(liveRect.Y+liveRect.H/2)

	// With the pointer away, both rows rest on the accent face.
	tt.Move(2, 2)
	tt.Frame()
	img := tt.Image()
	restOff, restOn := img.RGBAAt(offX, offY), img.RGBAAt(onX, onY)
	if want := pal.Accent.Over(pal.Background); distance(restOn, want) > 1 {
		t.Fatalf("enabled row at rest = %v, want the composited accent face %v", restOn, want)
	}
	if pixelDistance(restOff, restOn) <= 1 {
		t.Fatalf("disabled row = %v is not dimmed like the resting row %v", restOff, restOn)
	}

	// The pointer leaves the disabled row exactly as it was.
	tt.Move(idleRect.X+idleRect.W/2, idleRect.Y+idleRect.H/2)
	tt.Frame()
	gotOff := tt.Image().RGBAAt(offX, offY)
	if pixelDistance(gotOff, restOff) > 1 {
		t.Fatalf("disabled row under the pointer = %v, want its resting face %v", gotOff, restOff)
	}

	// The enabled row still takes the hover face.
	tt.Move(liveRect.X+liveRect.W/2, liveRect.Y+liveRect.H/2)
	tt.Frame()
	got := tt.Image().RGBAAt(onX, onY)
	if want := hover.Over(pal.Background); distance(got, want) > 1 {
		t.Fatalf("enabled row under the pointer = %v, want the composited hover face %v", got, want)
	}
	if pixelDistance(got, restOn) <= 1 {
		t.Fatal("enabled hover did not change the rendered face")
	}
}

// distance is the largest channel difference between a rendered pixel and a
// colour, the tolerance paint tests use.
func distance(pixel color.RGBA, want ui.Color) int {
	d := 0
	for _, c := range [][2]uint8{{pixel.R, want.R}, {pixel.G, want.G}, {pixel.B, want.B}} {
		diff := int(c[0]) - int(c[1])
		if diff < 0 {
			diff = -diff
		}
		if diff > d {
			d = diff
		}
	}
	return d
}

// pixelDistance is distance between two rendered pixels.
func pixelDistance(a, b color.RGBA) int {
	return distance(a, ui.RGB(b.R, b.G, b.B))
}
