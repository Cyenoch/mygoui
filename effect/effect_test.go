package effect_test

import (
	"image/color"
	"reflect"
	"testing"

	"github.com/Cyenoch/mygoui/effect"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// These tests use the public calls only. The tester has no GPU, so they
// cover the fallbacks and what an effect leaves alone.

var (
	behind  = ui.RGB(230, 40, 70)
	content = ui.RGB(30, 70, 230)
)

func opaque(c ui.Color) color.RGBA { return color.RGBA{R: c.R, G: c.G, B: c.B, A: 255} }

func scene(w, h float32, apply func(c *ui.Context, e ui.Element)) func(*ui.Context) {
	return func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ui.Box(c).Fill().Background(behind).Children(func() {
			apply(c, ui.Box(c).Absolute().Left(0).Top(0).Size(w, h).Background(content))
		})
	}
}

func TestMaskWithoutGPUShadesSelectedEdges(t *testing.T) {
	for _, tc := range []struct {
		name    string
		surface ui.Color
		lighter bool
	}{{"shadow", ui.Color{}, false}, {"surface", ui.RGB(255, 255, 255), true}} {
		t.Run(tc.name, func(t *testing.T) {
			tt := ui.NewTester(scene(40, 100, func(c *ui.Context, e ui.Element) {
				effect.Mask(c, e, effect.MaskOptions{Edges: effect.EdgeBottom, Surface: tc.surface})
			}), 40, 100)
			tt.Frame()
			img := tt.Image()
			if top, middle := img.RGBAAt(20, 0), img.RGBAAt(20, 50); top != opaque(content) || middle != opaque(content) {
				t.Fatalf("unselected area altered: top %v, middle %v", top, middle)
			}
			edge := img.RGBAAt(20, 99)
			if edge == opaque(content) || edge.A != 255 {
				t.Fatalf("bottom edge = %v, want it shaded", edge)
			}
			if got := edge.G > content.G; got != tc.lighter {
				t.Fatalf("bottom edge = %v, lighter than the content = %v, want %v", edge, got, tc.lighter)
			}
			if tt.FramePending() {
				t.Fatal("a still mask asked for a frame")
			}
		})
	}
}

func TestMaskRejectsATranslucentSurface(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("a translucent Surface did not fail")
		}
	}()
	tt := ui.NewTester(scene(40, 40, func(c *ui.Context, e ui.Element) {
		effect.Mask(c, e, effect.MaskOptions{Surface: ui.RGBA(255, 255, 255, 0.5)})
	}), 40, 40)
	tt.Frame()
}

func TestShimmerWithoutGPURests(t *testing.T) {
	tt := ui.NewTester(scene(200, 40, func(c *ui.Context, e ui.Element) {
		effect.Shimmer(c, e, effect.ShimmerOptions{})
	}), 200, 40)
	tt.Frame()
	if tt.FramePending() {
		t.Fatal("a shimmer without a GPU asked for a frame")
	}
	img := tt.Image()
	for x := range 200 {
		if got := img.RGBAAt(x, 20); got != opaque(content) {
			t.Fatalf("pixel %d = %v, want the content at rest", x, got)
		}
	}
}

func TestRevealWithoutGPUFadesInThenRests(t *testing.T) {
	duration := theme.Default().Motion.Reveal
	tt := ui.NewTester(scene(40, 100, func(c *ui.Context, e ui.Element) {
		effect.Reveal(c, e, effect.RevealOptions{})
	}), 40, 100)
	tt.Frame()
	tt.Advance(duration / 4)
	img := tt.Image()
	top, bottom := img.RGBAAt(20, 0), img.RGBAAt(20, 99)
	if top != bottom || top == opaque(content) || top == opaque(behind) {
		t.Fatalf("midway top = %v, bottom = %v: want one partial strength", top, bottom)
	}
	tt.Advance(duration)
	if got := tt.Image().RGBAAt(20, 99); got != opaque(content) {
		t.Fatalf("after the duration bottom = %v, want the content", got)
	}
	if tt.FramePending() {
		t.Fatal("a finished reveal asked for a frame")
	}
}

// An effect changes paint only: a control inside keeps its focus, its
// activation and what assistive technology reads.
func TestEffectsLeaveFocusActivationAndAccessibilityAlone(t *testing.T) {
	effects := map[string]func(*ui.Context, ui.Element){
		"none":    func(*ui.Context, ui.Element) {},
		"mask":    func(c *ui.Context, e ui.Element) { effect.Mask(c, e, effect.MaskOptions{}) },
		"shimmer": func(c *ui.Context, e ui.Element) { effect.Shimmer(c, e, effect.ShimmerOptions{}) },
		"reveal":  func(c *ui.Context, e ui.Element) { effect.Reveal(c, e, effect.RevealOptions{}) },
	}
	type observed struct {
		region, button ui.AccessInfo
	}
	seen := map[string]observed{}
	for name, apply := range effects {
		t.Run(name, func(t *testing.T) {
			runs := 0
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				region := ui.Column(c).Padding(12).Role(ui.RoleGroup).Label("Preview")
				region.Children(func() {
					ui.Button(c, "Run").OnClick(func() { runs++ })
				})
				apply(c, region)
			}, 200, 120)
			tt.Frame()
			if err := tt.Click("Run"); err != nil {
				t.Fatal(err)
			}
			tt.Key(0, ui.KeyEnter)
			if runs != 2 || !tt.Focused("Run") {
				t.Fatalf("runs = %d, focused = %v: want a click and a key activation on the focused button", runs, tt.Focused("Run"))
			}
			region, ok := tt.AccessByRole(ui.RoleGroup, "Preview")
			if !ok {
				t.Fatal("the region lost its role or name")
			}
			button, ok := tt.Access("Run")
			if !ok {
				t.Fatal("the button is missing from the accessibility tree")
			}
			// Shimmer marks the region busy, its one change.
			region.Busy = false
			seen[name] = observed{region, button}
		})
	}
	for name, got := range seen {
		if !reflect.DeepEqual(got, seen["none"]) {
			t.Errorf("%s changed the accessibility tree:\n got %+v\nwant %+v", name, got, seen["none"])
		}
	}
}
