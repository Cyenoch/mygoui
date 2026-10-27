package mygoui_test

import (
	"fmt"
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

func TestScrollFadeDecoratesOnlyRequestedEdges(t *testing.T) {
	for _, tc := range []struct {
		name string
		edge mygoui.ScrollFadeEdges
	}{
		{"top", mygoui.FadeTop},
		{"right", mygoui.FadeRight},
		{"bottom", mygoui.FadeBottom},
		{"left", mygoui.FadeLeft},
		{"all", 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var state ui.ScrollState
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				viewport := ui.ScrollBoth(c).Size(120, 100).TrackScroll(&state)
				viewport.Children(func() {
					ui.Box(c).Size(360, 300).Shrink(0).Background(fadeContent)
				})
				mygoui.ScrollFade(c, viewport, &state, mygoui.ScrollFadeOptions{
					Mode: mygoui.SurfaceFade, Surface: fadeSurface, Edges: tc.edge,
				})
			}, 160, 140)
			if state.MaxX <= 48 || state.MaxY <= 48 {
				t.Fatalf("both axes must overflow: %+v", state)
			}
			state.X, state.Y = 48, 48
			tt.Frame()
			img := tt.Image()
			for _, sample := range []struct {
				edge mygoui.ScrollFadeEdges
				x, y int
			}{
				{mygoui.FadeTop, 60, 0},
				{mygoui.FadeRight, 119, 50},
				{mygoui.FadeBottom, 60, 99},
				{mygoui.FadeLeft, 0, 50},
			} {
				got := pixel(img, float32(sample.x), float32(sample.y))
				if tc.edge == 0 || tc.edge&sample.edge != 0 {
					if !closeTo(got, fadeSurface, 6) {
						t.Errorf("edge %v = %v, want fade toward %v", sample.edge, got, fadeSurface)
					}
				} else {
					wantPixel(t, img, sample.x, sample.y, fadeContent, "unrequested edge")
				}
			}
			wantPixel(t, img, 60, 50, fadeContent, "content beyond the 24 DIP fade")
			if tt.FramePending() {
				t.Fatal("a stationary fade schedules animation frames")
			}
		})
	}
}

func TestScrollFadeShadowUsesColorAndDepth(t *testing.T) {
	for _, tc := range []struct {
		name  string
		color ui.Color
	}{
		{"opaque", ui.RGB(15, 25, 35)},
		{"translucent", ui.RGBA(15, 25, 35, 0.4)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var state ui.ScrollState
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				viewport := ui.Scroll(c).Size(120, 100).TrackScroll(&state)
				viewport.Children(func() {
					ui.Box(c).FillWidth().Height(300).Shrink(0).Background(fadeContent)
				})
				mygoui.ScrollFade(c, viewport, &state, mygoui.ScrollFadeOptions{
					Edges: mygoui.FadeBottom, Size: 20, Color: tc.color,
				})
			}, 160, 140)
			img := tt.Image()
			wantPixel(t, img, 60, 79, fadeContent, "content outside shadow depth")
			edge := tc.color.Over(fadeContent)
			partial := pixel(img, 60, 90)
			if colorDistance(partial, edge) >= colorDistance(fadeContent, edge) || same(partial, edge) {
				t.Fatalf("shadow midpoint = %v, want a blend of %v and %v", partial, edge, fadeContent)
			}
			if got := pixel(img, 60, 99); !closeTo(got, edge, 6) {
				t.Fatalf("shadow edge = %v, want %v", got, edge)
			}
			state.Y = state.MaxY
			tt.Frame()
			wantPixel(t, tt.Image(), 60, 99, fadeContent, "shadow clears at content end")
		})
	}
}

func TestScrollFadeRejectsInvalidBoundaryArguments(t *testing.T) {
	for _, tc := range []struct {
		name    string
		area    bool
		axis    mygoui.ScrollAxis
		state   *ui.ScrollState
		options mygoui.ScrollFadeOptions
		want    string
	}{
		{name: "fade nil state", want: "mygoui.ScrollFade: nil scroll state"},
		{name: "area nil state", area: true, want: "mygoui.ScrollArea: nil scroll state"},
		{
			name: "unknown axis", area: true, state: &ui.ScrollState{}, axis: mygoui.ScrollAxis(255),
			want: "mygoui.ScrollArea: unsupported scroll axis",
		},
		{
			name: "unknown mode", state: &ui.ScrollState{},
			options: mygoui.ScrollFadeOptions{Mode: mygoui.ScrollFadeMode(255)},
			want:    "mygoui.ScrollFade: unsupported mode",
		},
		{
			name: "transparent surface", state: &ui.ScrollState{},
			options: mygoui.ScrollFadeOptions{Mode: mygoui.SurfaceFade},
			want:    "mygoui.ScrollFade: SurfaceFade requires an opaque surface",
		},
		{
			name: "translucent surface", state: &ui.ScrollState{},
			options: mygoui.ScrollFadeOptions{Mode: mygoui.SurfaceFade, Surface: fadeSurface.Alpha(0.5)},
			want:    "mygoui.ScrollFade: SurfaceFade requires an opaque surface",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if got := recover(); fmt.Sprint(got) != tc.want {
					t.Errorf("panic = %v, want %q", got, tc.want)
				}
			}()
			ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				if tc.area {
					mygoui.ScrollArea(c, tc.state, mygoui.ScrollAreaOptions{Axis: tc.axis}, func() {})
					return
				}
				mygoui.ScrollFade(c, ui.Scroll(c), tc.state, tc.options)
			}, 160, 100)
		})
	}
}

func TestScrollAreaKeepsNativeAxesAndWheel(t *testing.T) {
	for _, tc := range []struct {
		name string
		axis mygoui.ScrollAxis
		x, y bool
	}{
		{name: "default vertical", y: true},
		{name: "horizontal", axis: mygoui.ScrollHorizontal, x: true},
		{name: "both", axis: mygoui.ScrollBoth, x: true, y: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var state ui.ScrollState
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				mygoui.ScrollArea(c, &state, mygoui.ScrollAreaOptions{
					Axis: tc.axis, Height: 100, Label: "Viewport",
				}, func() {
					ui.Box(c).Size(360, 300).Shrink(0)
				}).Width(120)
			}, 160, 140)
			viewport := findRect(t, tt, "Viewport")
			if viewport.W != 120 || viewport.H != 100 {
				t.Fatalf("viewport = %+v, want 120 x 100 DIP", viewport)
			}
			if (state.MaxX > 0) != tc.x || (state.MaxY > 0) != tc.y {
				t.Fatalf("scroll ranges = %+v, want horizontal=%v vertical=%v", state, tc.x, tc.y)
			}
			dx, dy := float32(0), float32(0)
			if tc.x {
				dx = 30
			}
			if tc.y {
				dy = 40
			}
			tt.Scroll(60, 50, dx, dy)
			if state.X != dx || state.Y != dy {
				t.Fatalf("native wheel position = (%v, %v), want (%v, %v)", state.X, state.Y, dx, dy)
			}
			if tc.x {
				tt.Key(0, ui.KeyRight)
				dx += 40
			}
			if tc.y {
				tt.Key(0, ui.KeyDown)
				dy += 40
			}
			if state.X != dx || state.Y != dy {
				t.Fatalf("native keyboard position = (%v, %v), want (%v, %v)", state.X, state.Y, dx, dy)
			}
		})
	}
}

func TestScrollAreaFadePreservesChildActivationAndFocus(t *testing.T) {
	var state ui.ScrollState
	calls := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		mygoui.ScrollArea(c, &state, mygoui.ScrollAreaOptions{
			Height: 100, Label: "Actions",
			Fade: &mygoui.ScrollFadeOptions{
				Mode: mygoui.SurfaceFade, Surface: theme.Colors(c).Background,
			},
		}, func() {
			ui.Column(c).Children(func() {
				ui.Box(c).Height(68).Shrink(0)
				mygoui.Button(c, mygoui.ButtonOptions{
					Label: "Run", OnClick: func() { calls++ },
				})
				ui.Box(c).Height(200).Shrink(0)
			})
		})
	}, 220, 140)
	if err := tt.Click("Run"); err != nil {
		t.Fatal(err)
	}
	if calls != 1 || !tt.Focused("Run") {
		t.Fatalf("faded button calls=%d focus=%v", calls, tt.Focused("Run"))
	}
	tt.Key(0, ui.KeyEnter)
	if calls != 2 {
		t.Fatalf("faded button keyboard activations=%d, want 2", calls)
	}
	tt.Scroll(180, 50, 0, 40)
	if state.Y != 40 {
		t.Fatalf("fade swallowed wheel input: Y=%v", state.Y)
	}
}
