package overlay

import (
	"fmt"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo"
	"github.com/egoist/mygo/ui"
)

func TestFlyoutOptions(t *testing.T) {
	parent := new(mygo.Window)
	tray := new(mygo.Tray)
	for _, tc := range []struct {
		name  string
		opts  FlyoutOptions
		valid bool
	}{
		{"parent", FlyoutOptions{Parent: parent, Width: 300, Height: 200}, true},
		{"tray", FlyoutOptions{Tray: tray, Width: 300, Height: 200}, true},
		{"missing owner", FlyoutOptions{Width: 300, Height: 200}, false},
		{"both owners", FlyoutOptions{Parent: parent, Tray: tray, Width: 300, Height: 200}, false},
		{"missing width", FlyoutOptions{Parent: parent, Height: 200}, false},
		{"negative height", FlyoutOptions{Parent: parent, Width: 300, Height: -1}, false},
		{"unknown placement", FlyoutOptions{Parent: parent, Width: 300, Height: 200, Placement: "invalid"}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := nativeFlyoutOptions(tc.opts)
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%v, error=%v", tc.valid, err)
			}
		})
	}
	o := FlyoutOptions{Parent: parent, Anchor: mygo.Rectangle{X: 10, Y: 20, Width: 30, Height: 40},
		Placement: mygo.PlacementTopEnd, Gap: 7, Width: 320, Height: 180, Focusable: true, Hidden: true}
	got, err := nativeFlyoutOptions(o)
	if err != nil {
		t.Fatal(err)
	}
	if got.Parent != o.Parent || got.Anchor != o.Anchor || got.Placement != o.Placement || got.Gap != o.Gap ||
		got.Width != o.Width || got.Height != o.Height || !got.Focusable || !got.Hidden || got.Content == nil {
		t.Fatalf("flyout options not preserved: %+v", got)
	}
}

func TestFlyoutAnchorZoom(t *testing.T) {
	var handle ui.Handle
	var anchor mygo.Rectangle
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Box(c).Absolute().Left(11).Top(21).Width(31).Height(41).Bind(&handle)
		anchor = FlyoutAnchor(c, handle)
	}, 300, 200)
	tt.SetZoom(1.5)
	tt.Frame()
	want := mygo.Rectangle{X: 16, Y: 31, Width: 47, Height: 62}
	if anchor != want {
		t.Fatalf("zoomed anchor = %+v, want %+v", anchor, want)
	}
}

func TestFlyoutSurfaceScrollingAndFocus(t *testing.T) {
	actions := 0
	var panel ui.Rect
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		p := FlyoutSurface(c, FlyoutSurfaceOptions{
			Title: "Details", Description: "Description",
			Body: func(c *ui.Context) {
				ui.Column(c).Gap(8).Children(func() {
					for i := range 12 {
						ui.Text(c, fmt.Sprintf("Row %d", i))
					}
				})
			},
			Footer: func(c *ui.Context) { ui.Button(c, "Last action").OnClick(func() { actions++ }) },
		})
		panel = p.Bounds()
	}, 280, 160)
	tt.Frame()
	if panel.W != 280 || panel.H != 160 {
		t.Fatalf("surface bounds = %+v", panel)
	}
	if !tt.HasText("Details") || !tt.HasText("Description") {
		t.Fatal("surface header missing")
	}
	tt.Scroll(140, 80, 0, 1000)
	tt.Frame()
	if err := tt.Click("Last action"); err != nil {
		t.Fatal(err)
	}
	if actions != 1 || !tt.Focused("Last action") {
		t.Fatalf("scrolled action: count=%d focused=%v", actions, tt.Focused("Last action"))
	}
}

func TestFlyoutSurfaceHeaderAndTheme(t *testing.T) {
	for _, dark := range []bool{false, true} {
		t.Run(fmt.Sprintf("dark=%v", dark), func(t *testing.T) {
			tt := ui.NewTester(func(c *ui.Context) {
				theme.Apply(c)
				FlyoutSurface(c, FlyoutSurfaceOptions{
					Title: "Replaced", Label: "Surface", Header: func(c *ui.Context) { ui.Text(c, "Custom header") },
				})
			}, 200, 120)
			tt.SetDark(dark)
			if !tt.HasText("Custom header") || tt.HasText("Replaced") {
				t.Fatal("custom header did not replace title")
			}
			got := tt.Image().RGBAAt(100, 100)
			if want := theme.ColorsForMode(dark).Popover; !sameColor(got, want, 3) {
				t.Fatalf("surface color = %v, want %v", got, want)
			}
		})
	}
}
