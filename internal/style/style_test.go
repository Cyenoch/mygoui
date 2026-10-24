package style_test

import (
	"image/color"
	"testing"

	"github.com/Cyenoch/mygoui/internal/style"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestHeadingTrackingTightensTheLine pins tracking-tight on the heading roles:
// the same text set as an H2 measures narrower than the same 20/28 role
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
	var pal theme.Palette
	var idle, live ui.Handle
	var idleRect, liveRect ui.Rect
	frames := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		pal = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).Gap(12).
			Background(pal.Background).Children(func() {
			off := ui.ButtonBase(c).FillWidth().MinHeight(40).Label("off")
			off = style.Disabled(off, true)
			off = style.Background(off, 8, pal.Accent, pal.AccentHover, pal.AccentHover)
			off.Bind(&idle)
			on := ui.ButtonBase(c).FillWidth().MinHeight(40).Label("on")
			on = style.Background(on, 8, pal.Accent, pal.AccentHover, pal.AccentHover)
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
	restOff, restOn := tt.Image().RGBAAt(offX, offY), tt.Image().RGBAAt(onX, onY)
	if distance(restOn, pal.Accent) > 3 {
		t.Fatalf("enabled row at rest = %v, want the accent face %v", restOn, pal.Accent)
	}
	if pixelDistance(restOff, restOn) <= 3 {
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
	if got := tt.Image().RGBAAt(onX, onY); distance(got, pal.AccentHover) > 3 {
		t.Fatalf("enabled row under the pointer = %v, want the hover face %v", got, pal.AccentHover)
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
