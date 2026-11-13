package effect

import (
	"image"
	"image/color"
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// The headless tester never reports a GPU, so the public calls reach only
// the fallbacks. These tests set the forced flag on the element's effect
// state: the frame then behaves as a GPU frame and MyGo interprets the
// shader on the CPU. They show what the shader computes, not what a GPU
// presents.

var (
	behind  = ui.RGB(230, 40, 70)
	content = ui.RGB(30, 70, 230)
)

func opaque(c ui.Color) color.RGBA { return color.RGBA{R: c.R, G: c.G, B: c.B, A: 255} }

// scene builds a content-colored element of w by h over the behind color
// and hands it to apply.
func scene(w, h float32, apply func(c *ui.Context, e ui.Element)) func(*ui.Context) {
	return func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Box(c).Fill().Background(behind).Children(func() {
			apply(c, ui.Box(c).Absolute().Left(0).Top(0).Size(w, h).Background(content))
		})
	}
}

func useTheme(t *testing.T, th theme.Theme) {
	t.Helper()
	if err := theme.Use(th); err != nil {
		t.Fatal(err)
	}
}

func forcedMask(o MaskOptions) func(*ui.Context, ui.Element) {
	return func(c *ui.Context, e ui.Element) {
		Mask(c, e, o)
		ui.Local(e, maskKey{}, newMask).forced = true
	}
}

func TestMaskLetsTheBackgroundThroughAtSelectedEdges(t *testing.T) {
	tt := ui.NewTester(scene(40, 100, forcedMask(MaskOptions{Edges: EdgeBottom, Size: 20})), 40, 100)
	tt.Frame()
	img := tt.Image()
	at := func(y int) color.RGBA { return img.RGBAAt(20, y) }
	if at(0) != opaque(content) || at(50) != opaque(content) || at(79) != opaque(content) {
		t.Errorf("unmasked area altered: %v %v %v", at(0), at(50), at(79))
	}
	redder := at(99).R > at(90).R && at(90).R > at(79).R
	lessBlue := at(99).B < at(90).B && at(90).B < at(79).B
	if !redder || !lessBlue {
		t.Errorf("bottom edge does not ramp to the background: %v %v %v", at(79), at(90), at(99))
	}
	if tt.FramePending() {
		t.Error("a still mask asked for a frame")
	}
}

// Opposite fades meet in the middle of a small element; they do not
// overlap and dim all of it.
func TestMaskKeepsTheMiddleOfASmallElementOpaque(t *testing.T) {
	tt := ui.NewTester(scene(100, 32, forcedMask(MaskOptions{})), 100, 32)
	middle, corner := tt.Image().RGBAAt(50, 16), tt.Image().RGBAAt(0, 0)
	if middle.B < 215 || middle.R > 50 {
		t.Errorf("middle = %v, want nearly all of the content %v", middle, opaque(content))
	}
	if corner.R < 200 {
		t.Errorf("corner = %v, want the background through it", corner)
	}
}

// Fallback keeps the shader off where it could run: the content stays whole
// and the edge is shaded, not let through.
func TestMaskFallbackShadesWhereTheShaderCouldRun(t *testing.T) {
	tt := ui.NewTester(scene(40, 100, forcedMask(MaskOptions{Edges: EdgeBottom, Size: 20, Fallback: true})), 40, 100)
	img := tt.Image()
	edge, inside := img.RGBAAt(20, 99), img.RGBAAt(20, 50)
	if inside != opaque(content) {
		t.Fatalf("content beyond the edge = %v, want it unchanged", inside)
	}
	if edge == opaque(content) || edge.R >= 100 {
		t.Fatalf("bottom edge = %v, want shaded content, not the background", edge)
	}
}

// Above the layer limit MyGo paints the element without its effect, so the
// shaded edge has to come from the element itself.
func TestMaskShadesItsEdgeAboveTheLayerLimit(t *testing.T) {
	for _, tc := range []struct {
		name   string
		width  float32
		shaded bool
	}{{"fits", 100, false}, {"oversized", 9000, true}} {
		t.Run(tc.name, func(t *testing.T) {
			tt := ui.NewTester(scene(tc.width, 40, forcedMask(MaskOptions{Edges: EdgeLeft})), 100, 40)
			tt.Frame()
			img := tt.Image()
			edge, inside := img.RGBAAt(0, 20), img.RGBAAt(60, 20)
			if inside != opaque(content) {
				t.Fatalf("content beyond the edge = %v, want it unchanged", inside)
			}
			if got := edge.R < 100 && edge != opaque(content); got != tc.shaded {
				t.Fatalf("left edge = %v, shaded content = %v, want %v", edge, got, tc.shaded)
			}
		})
	}
}

func TestMaskDepthFollowsTheFadeEdgeToken(t *testing.T) {
	t.Cleanup(func() { useTheme(t, theme.Default()) })
	custom := theme.Default()
	custom.Effects.FadeEdge = 2 * theme.Default().Effects.FadeEdge
	// 30 DIPs from the bottom: beyond the built-in depth, inside the doubled one.
	const y = 69
	render := func() color.RGBA {
		tt := ui.NewTester(scene(40, 100, forcedMask(MaskOptions{Edges: EdgeBottom})), 40, 100)
		tt.Frame()
		return tt.Image().RGBAAt(20, y)
	}
	if got := render(); got != opaque(content) {
		t.Fatalf("built-in depth reaches %d: %v", y, got)
	}
	useTheme(t, custom)
	if got := render(); got == opaque(content) {
		t.Fatalf("doubled FadeEdge did not deepen the mask at %d", y)
	}
}

// shimmerScene is a small window whose element shimmers while *on. The
// shader is interpreted per pixel on every frame it asks for, so the tests
// keep the scene small and move the clock before applying it.
func shimmerScene(on *bool) *ui.Tester {
	return ui.NewTester(scene(96, 8, func(c *ui.Context, e ui.Element) {
		if *on {
			Shimmer(c, e, ShimmerOptions{})
			ui.Local(e, shimmerKey{}, newShimmer).forced = true
		}
	}), 96, 8)
}

// brightest returns the largest red value in the row: the content has
// little red, the band mixes towards white.
func brightest(img *image.RGBA, y int) uint8 {
	var top uint8
	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
		top = max(top, img.RGBAAt(x, y).R)
	}
	return top
}

// toMidLoop moves a resting tester's clock to the middle of a shimmer
// loop, where the band crosses the middle of the window.
func toMidLoop(tt *ui.Tester) {
	period := theme.Effects().ShimmerPeriod
	into := time.Duration(tt.Now().UnixNano() % int64(period))
	tt.Advance((period/2 - into + period) % period)
}

func TestShimmerSweepsWhileAppliedAndRestsWhenAskedTo(t *testing.T) {
	for _, tc := range []struct {
		name  string
		prefs ui.Preferences
		moves bool
	}{
		{"default", ui.Preferences{}, true},
		{"reduced motion", ui.Preferences{ReduceMotion: true}, false},
		{"high contrast", ui.Preferences{HighContrast: true}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			on := false
			tt := shimmerScene(&on)
			tt.SetPreferences(tc.prefs)
			toMidLoop(tt)
			on = true
			tt.Frame()
			if got := brightest(tt.Image(), 4) > content.R; got != tc.moves {
				t.Errorf("highlight band shows = %v, want %v", got, tc.moves)
			}
			if got := tt.FramePending(); got != tc.moves {
				t.Errorf("asks for frames = %v, want %v", got, tc.moves)
			}
		})
	}
}

func TestShimmerStopsWithItsState(t *testing.T) {
	loading := true
	tt := shimmerScene(&loading)
	if !tt.FramePending() {
		t.Fatal("a running shimmer asked for no frame")
	}
	loading = false
	tt.Frame()
	if tt.FramePending() {
		t.Fatal("the shimmer kept asking for frames after loading ended")
	}
	if got := brightest(tt.Image(), 4); got != content.R {
		t.Fatalf("content at rest has a highlight: red up to %d", got)
	}
}

func TestShimmerStrengthFollowsThePeakToken(t *testing.T) {
	t.Cleanup(func() { useTheme(t, theme.Default()) })
	render := func() uint8 {
		on := false
		tt := shimmerScene(&on)
		toMidLoop(tt)
		on = true
		tt.Frame()
		return brightest(tt.Image(), 4)
	}
	builtin := render()
	custom := theme.Default()
	custom.Effects.ShimmerPeak = 1
	useTheme(t, custom)
	if got := render(); got <= builtin || got < 235 {
		t.Fatalf("band peak = %d with ShimmerPeak 1, built-in %d: want near white", got, builtin)
	}
}

// A period of zero is the theme's way to keep placeholders still.
func TestShimmerRestsForAZeroPeriod(t *testing.T) {
	t.Cleanup(func() { useTheme(t, theme.Default()) })
	custom := theme.Default()
	custom.Effects.ShimmerPeriod = 0
	useTheme(t, custom)
	on := true
	tt := shimmerScene(&on)
	if tt.FramePending() {
		t.Fatal("a shimmer with no period asked for a frame")
	}
	if got := brightest(tt.Image(), 4); got != content.R {
		t.Fatalf("content has a highlight: red up to %d", got)
	}
}

func forcedReveal(o RevealOptions) func(*ui.Context, ui.Element) {
	return func(c *ui.Context, e ui.Element) {
		Reveal(c, e, o)
		ui.Local(e, revealKey{}, newReveal).forced = true
	}
}

func TestRevealWipesInFromItsEdgeOnceThenRests(t *testing.T) {
	duration := theme.Default().Motion.Reveal
	tt := ui.NewTester(scene(40, 100, forcedReveal(RevealOptions{})), 40, 100)
	tt.Frame()
	if got := tt.Image().RGBAAt(20, 50); got != opaque(behind) {
		t.Fatalf("content shows before the reveal starts: %v", got)
	}
	tt.Advance(duration / 4)
	img := tt.Image()
	if top, bottom := img.RGBAAt(20, 0), img.RGBAAt(20, 99); top != opaque(content) || bottom != opaque(behind) {
		t.Fatalf("midway top = %v, bottom = %v: want the top revealed first", top, bottom)
	}
	if !tt.FramePending() {
		t.Fatal("a reveal on its way asked for no frame")
	}
	tt.Advance(duration)
	img = tt.Image()
	if top, bottom := img.RGBAAt(20, 0), img.RGBAAt(20, 99); top != opaque(content) || bottom != opaque(content) {
		t.Fatalf("after the duration top = %v, bottom = %v: want all of it", top, bottom)
	}
	if tt.FramePending() {
		t.Fatal("a finished reveal asked for a frame")
	}
	tt.Frame()
	if got := tt.Image().RGBAAt(20, 99); got != opaque(content) || tt.FramePending() {
		t.Fatalf("a later build replayed the reveal: %v", got)
	}
}

// Reduced motion and Fallback keep the fade and drop the travelling front,
// also where the shader could run.
func TestRevealFadesUnderReducedMotionAndFallback(t *testing.T) {
	duration := theme.Default().Motion.Reveal
	for _, tc := range []struct {
		name    string
		prefs   ui.Preferences
		options RevealOptions
	}{
		{"reduced motion", ui.Preferences{ReduceMotion: true}, RevealOptions{}},
		{"fallback", ui.Preferences{}, RevealOptions{Fallback: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tt := ui.NewTester(scene(40, 100, forcedReveal(tc.options)), 40, 100)
			tt.SetPreferences(tc.prefs)
			tt.Frame()
			tt.Advance(duration / 4)
			img := tt.Image()
			top, bottom := img.RGBAAt(20, 0), img.RGBAAt(20, 99)
			if top != bottom || top == opaque(content) || top == opaque(behind) {
				t.Fatalf("midway top = %v, bottom = %v: want one partial strength", top, bottom)
			}
			tt.Advance(duration)
			if got := tt.Image().RGBAAt(20, 99); got != opaque(content) || tt.FramePending() {
				t.Fatalf("after the duration bottom = %v, frame pending = %v", got, tt.FramePending())
			}
		})
	}
}

// An element has one effect: the later one draws it.
func TestALaterEffectReplacesAnEarlierOne(t *testing.T) {
	tt := ui.NewTester(scene(40, 100, func(c *ui.Context, e ui.Element) {
		Shimmer(c, e, ShimmerOptions{})
		ui.Local(e, shimmerKey{}, newShimmer).forced = true
		forcedMask(MaskOptions{Edges: EdgeBottom})(c, e)
	}), 40, 100)
	tt.Frame()
	if tt.FramePending() {
		t.Error("the replaced shimmer still asks for frames")
	}
	if got := tt.Image().RGBAAt(20, 99); got.R < 200 {
		t.Errorf("bottom edge = %v, want the mask letting the background through", got)
	}
}
