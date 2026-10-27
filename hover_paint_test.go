package mygoui_test

import (
	"testing"

	"github.com/Cyenoch/mygoui"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// hoverButtonSample returns a point left of the centered label, on the
// button's face.
func hoverButtonSample(r ui.Rect) (int, int) {
	return int(r.X) + 4, int(r.Y + r.H/2)
}

// TestButtonHoverAndPressedPaint checks the primary and ghost button faces
// follow live hover and press: MyGo rebuilds the view for a pointer move, so
// the build-time snapshot taken by style.Background must take the hover face,
// and a pressed button keeps it.
func TestButtonHoverAndPressedPaint(t *testing.T) {
	for _, dark := range []bool{false, true} {
		mode := "light"
		if dark {
			mode = "dark"
		}
		t.Run(mode, func(t *testing.T) {
			for _, tc := range []struct {
				name          string
				variant       mygoui.ButtonVariant
				normal, hover func(theme.Palette) ui.Color
			}{
				{
					name: "primary", variant: mygoui.ButtonPrimary,
					normal: func(p theme.Palette) ui.Color { return p.Primary },
					hover:  func(p theme.Palette) ui.Color { return p.Primary.Alpha(0.9).Over(p.Background) },
				},
				{
					name: "ghost", variant: mygoui.ButtonGhost,
					normal: func(p theme.Palette) ui.Color { return p.Background },
					hover:  func(p theme.Palette) ui.Color { return p.Accent.Over(p.Background) },
				},
			} {
				t.Run(tc.name, func(t *testing.T) {
					var colors theme.Palette
					tt := ui.NewTester(func(c *ui.Context) {
						theme.ApplyMode(c, dark)
						colors = theme.Colors(c)
						ui.Column(c).FillWidth().FillHeight().Padding(24).
							Background(colors.Background).Children(func() {
							mygoui.Button(c, mygoui.ButtonOptions{
								Label: "Go", Variant: tc.variant, OnClick: func() {},
							})
						})
					}, 260, 140)

					rect, ok := tt.Find("Go")
					if !ok {
						t.Fatal("button missing")
					}
					x, y := hoverButtonSample(rect)
					normal, hover := tc.normal(colors), tc.hover(colors)
					wantPixel(t, tt.Image(), x, y, normal, tc.name+" normal fill")

					tt.Move(rect.X+rect.W/2, rect.Y+rect.H/2)
					tt.Frame()
					wantPixel(t, tt.Image(), x, y, hover, tc.name+" hover fill")

					tt.Press(rect.X+rect.W/2, rect.Y+rect.H/2)
					tt.Frame()
					wantPixel(t, tt.Image(), x, y, hover, tc.name+" pressed fill")
					tt.Release(rect.X+rect.W/2, rect.Y+rect.H/2)

					tt.Move(2, 2)
					tt.Frame()
					wantPixel(t, tt.Image(), x, y, normal, tc.name+" fill after leaving")
				})
			}
		})
	}
}

// TestToggleButtonHoverAndPressedPaint checks the standalone toggle's default,
// selected and ghost faces follow hover and press.
func TestToggleButtonHoverAndPressedPaint(t *testing.T) {
	for _, dark := range []bool{false, true} {
		mode := "light"
		if dark {
			mode = "dark"
		}
		t.Run(mode, func(t *testing.T) {
			for _, tc := range []struct {
				name          string
				options       mygoui.ToggleButtonOptions
				normal, hover func(theme.Palette) ui.Color
			}{
				{
					name:   "default",
					normal: func(p theme.Palette) ui.Color { return p.Background },
					hover:  func(p theme.Palette) ui.Color { return p.Accent.Over(p.Background) },
				},
				{
					name:    "selected",
					options: mygoui.ToggleButtonOptions{Selected: new(bool)},
					normal: func(p theme.Palette) ui.Color {
						return p.Input.Alpha(0.64).Over(p.Background)
					},
					hover: func(p theme.Palette) ui.Color { return p.Input.Over(p.Background) },
				},
				{
					name:    "ghost",
					options: mygoui.ToggleButtonOptions{Variant: mygoui.ToggleGhost},
					normal:  func(p theme.Palette) ui.Color { return p.Background },
					hover:   func(p theme.Palette) ui.Color { return p.Accent.Over(p.Background) },
				},
			} {
				t.Run(tc.name, func(t *testing.T) {
					var colors theme.Palette
					selected := false
					opts := tc.options
					if opts.Selected != nil {
						selected = true
						opts.Selected = &selected
					}
					opts.Label = "Toggle"
					tt := ui.NewTester(func(c *ui.Context) {
						theme.ApplyMode(c, dark)
						colors = theme.Colors(c)
						ui.Column(c).FillWidth().FillHeight().Padding(24).
							Background(colors.Background).Children(func() {
							mygoui.ToggleButton(c, opts)
						})
					}, 260, 140)

					rect, ok := tt.Find("Toggle")
					if !ok {
						t.Fatal("toggle missing")
					}
					x, y := hoverButtonSample(rect)
					normal, hover := tc.normal(colors), tc.hover(colors)
					wantPixel(t, tt.Image(), x, y, normal, tc.name+" normal fill")

					tt.Move(rect.X+rect.W/2, rect.Y+rect.H/2)
					tt.Frame()
					wantPixel(t, tt.Image(), x, y, hover, tc.name+" hover fill")

					tt.Press(rect.X+rect.W/2, rect.Y+rect.H/2)
					tt.Frame()
					wantPixel(t, tt.Image(), x, y, hover, tc.name+" pressed fill")
					tt.Release(rect.X+rect.W/2, rect.Y+rect.H/2)
					if opts.Selected != nil {
						// Releasing activates the toggle and flips the caller's
						// state; restore it so leaving still checks the selected
						// face.
						selected = true
						tt.Frame()
					}

					tt.Move(2, 2)
					tt.Frame()
					wantPixel(t, tt.Image(), x, y, normal, tc.name+" fill after leaving")
				})
			}
		})
	}
}

// TestButtonGroupGhostHoverPaint checks a joined ghost group takes the neutral
// hover face through CornerFill's paint-time owner resolution.
func TestButtonGroupGhostHoverPaint(t *testing.T) {
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			mygoui.ButtonGroup(c, mygoui.ButtonGroupOptions{},
				mygoui.ButtonOptions{Label: "One", Variant: mygoui.ButtonGhost, OnClick: func() {}},
				mygoui.ButtonOptions{Label: "Two", Variant: mygoui.ButtonGhost, OnClick: func() {}},
			)
		})
	}, 320, 140)

	rect, ok := tt.Find("Two")
	if !ok {
		t.Fatal("group button missing")
	}
	// The last segment keeps its rounded end, so sample off the rounded
	// corner but inside the fill.
	x, y := int(rect.X)+8, int(rect.Y+rect.H/2)
	wantPixel(t, tt.Image(), x, y, colors.Background, "group normal fill")

	tt.Move(rect.X+rect.W/2, rect.Y+rect.H/2)
	tt.Frame()
	wantPixel(t, tt.Image(), x, y, colors.Accent.Over(colors.Background), "group hover fill")

	tt.Move(2, 2)
	tt.Frame()
	wantPixel(t, tt.Image(), x, y, colors.Background, "group fill after leaving")
}
