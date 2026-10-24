package form

import (
	"image"
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// formHoverPixel returns the colour at x, y in img.
func formHoverPixel(img *image.RGBA, x, y int) (uint8, uint8, uint8) {
	c := img.RGBAAt(x, y)
	return c.R, c.G, c.B
}

// formHoverNear reports whether got is within the 3-level antialias tolerance
// of want.
func formHoverNear(got, want uint8) bool {
	d := int(got) - int(want)
	return d >= -3 && d <= 3
}

// wantFormHoverPixel fails unless img paints c at x, y.
func wantFormHoverPixel(t *testing.T, img *image.RGBA, x, y int, c ui.Color, what string) {
	t.Helper()
	r, g, b := formHoverPixel(img, x, y)
	if !formHoverNear(r, c.R) || !formHoverNear(g, c.G) || !formHoverNear(b, c.B) {
		t.Fatalf("%s at (%d, %d) = %d,%d,%d, want %d,%d,%d", what, x, y, r, g, b, c.R, c.G, c.B)
	}
}

// TestCheckboxHoverAndPressedPaint checks a checked box follows the row's
// live hover and press faces and reverts when the pointer leaves.
func TestCheckboxHoverAndPressedPaint(t *testing.T) {
	var colors theme.Palette
	checked := true
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			Checkbox(c, CheckboxOptions{Checked: &checked, Label: "Accept"})
		})
	}, 400, 140)

	row, ok := tt.Find("Accept")
	if !ok {
		t.Fatal("checkbox missing")
	}
	// The 16 DIP box sits at the row's start; sample off its check mark.
	x, y := int(row.X)+3, int(row.Y+row.H/2)
	wantFormHoverPixel(t, tt.Image(), x, y, colors.Accent, "checkbox normal fill")

	tt.Move(row.X+row.W/2, row.Y+row.H/2)
	tt.Frame()
	wantFormHoverPixel(t, tt.Image(), x, y, colors.AccentHover, "checkbox hover fill")

	tt.Press(row.X+row.W/2, row.Y+row.H/2)
	tt.Frame()
	wantFormHoverPixel(t, tt.Image(), x, y, colors.AccentHover, "checkbox pressed fill")
	tt.Release(row.X+row.W/2, row.Y+row.H/2)
	// Releasing activates the box and clears the caller's value; restore it
	// so leaving still checks the checked face.
	checked = true
	tt.Frame()

	tt.Move(2, 2)
	tt.Frame()
	wantFormHoverPixel(t, tt.Image(), x, y, colors.Accent, "checkbox fill after leaving")
}

// TestRadioHoverAndPressedPaint checks an unselected radio takes the field
// hover face and a selected one takes the accent pressed face.
func TestRadioHoverAndPressedPaint(t *testing.T) {
	t.Run("unselected hover", func(t *testing.T) {
		var colors theme.Palette
		selected := "b"
		state := RadioGroupState[string]{Selected: &selected}
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			colors = theme.Colors(c)
			ui.Column(c).FillWidth().FillHeight().Padding(24).
				Background(colors.Background).Children(func() {
				RadioGroup(c, RadioGroupOptions{Label: "Group"}, func() {
					Radio(c, RadioOptions[string]{State: &state, Value: "a", Label: "Alpha"})
				})
			})
		}, 400, 200)

		row, ok := tt.Find("Alpha")
		if !ok {
			t.Fatal("radio missing")
		}
		x, y := int(row.X)+8, int(row.Y+row.H/2)
		wantFormHoverPixel(t, tt.Image(), x, y, colors.FieldBackground, "radio normal fill")

		tt.Move(row.X+row.W/2, row.Y+row.H/2)
		tt.Frame()
		wantFormHoverPixel(t, tt.Image(), x, y, colors.FieldHover, "radio hover fill")

		tt.Move(2, 2)
		tt.Frame()
		wantFormHoverPixel(t, tt.Image(), x, y, colors.FieldBackground, "radio fill after leaving")
	})

	t.Run("selected pressed", func(t *testing.T) {
		var colors theme.Palette
		selected := "a"
		state := RadioGroupState[string]{Selected: &selected}
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			colors = theme.Colors(c)
			ui.Column(c).FillWidth().FillHeight().Padding(24).
				Background(colors.Background).Children(func() {
				RadioGroup(c, RadioGroupOptions{Label: "Group"}, func() {
					Radio(c, RadioOptions[string]{State: &state, Value: "a", Label: "Alpha"})
				})
			})
		}, 400, 200)

		row, ok := tt.Find("Alpha")
		if !ok {
			t.Fatal("radio missing")
		}
		// The selected face is an accent ring around the white dot; sample
		// the ring, clear of the dot's pressed growth.
		x, y := int(row.X)+2, int(row.Y+row.H/2)
		wantFormHoverPixel(t, tt.Image(), x, y, colors.Accent, "radio selected fill")

		tt.Press(row.X+row.W/2, row.Y+row.H/2)
		tt.Frame()
		wantFormHoverPixel(t, tt.Image(), x, y, colors.AccentHover, "radio pressed fill")
		tt.Release(row.X+row.W/2, row.Y+row.H/2)
		tt.Frame()
		wantFormHoverPixel(t, tt.Image(), x, y, colors.Accent, "radio fill after release")
	})
}

// TestSwitchHoverPaint checks both the on and off tracks follow hover, and
// revert when the pointer leaves.
func TestSwitchHoverPaint(t *testing.T) {
	for _, tc := range []struct {
		name          string
		on            bool
		xOffset       float32
		normal, hover func(theme.Palette) ui.Color
	}{
		{
			name: "on", on: true, xOffset: 5,
			normal: func(p theme.Palette) ui.Color { return p.Accent },
			hover:  func(p theme.Palette) ui.Color { return p.AccentHover },
		},
		{
			name: "off", on: false, xOffset: 35,
			normal: func(p theme.Palette) ui.Color { return p.Default },
			hover:  func(p theme.Palette) ui.Color { return p.DefaultHover },
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var colors theme.Palette
			on := tc.on
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				colors = theme.Colors(c)
				ui.Column(c).FillWidth().FillHeight().Padding(24).
					Background(colors.Background).Children(func() {
					Switch(c, SwitchOptions{On: &on, Label: "Wifi"})
				})
			}, 400, 140)

			row, ok := tt.Find("Wifi")
			if !ok {
				t.Fatal("switch missing")
			}
			x, y := int(row.X+tc.xOffset), int(row.Y+row.H/2)
			normal, hover := tc.normal(colors), tc.hover(colors)
			wantFormHoverPixel(t, tt.Image(), x, y, normal, tc.name+" track")

			tt.Move(row.X+row.W/2, row.Y+row.H/2)
			tt.Frame()
			wantFormHoverPixel(t, tt.Image(), x, y, hover, tc.name+" track under the pointer")

			tt.Move(2, 2)
			tt.Frame()
			wantFormHoverPixel(t, tt.Image(), x, y, normal, tc.name+" track after leaving")
		})
	}
}

// TestSliderPressedShrinksKnob checks the slider's pressed snapshot: HeroUI
// gives the thumb no hover face, so a press is the visible pointer state. The
// inner cap shrinks, exposing the accent pill just inside its resting edge.
func TestSliderPressedShrinksKnob(t *testing.T) {
	var colors theme.Palette
	value := 50.0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			Slider(c, SliderOptions{
				State: SliderState{Value: &value}, Label: "Volume",
				Min: 0, Max: 100, HideOutput: true,
			})
		})
	}, 400, 200)

	header, ok := tt.Find("Volume")
	if !ok {
		t.Fatal("slider label missing")
	}
	const pad = float32(14)
	trackY := header.Y + header.H + theme.Space4
	width := float32(400) - 2*header.X
	cx := header.X + pad + (width-2*pad)*0.5
	cy := trackY + 10
	// The pill's inner cap spans 24×16 DIP, so (11, -3) from the centre is
	// inside its resting edge and outside the pressed 21.6×14.4 cap.
	x, y := int(cx)+11, int(cy)-3

	tt.Move(cx, cy)
	tt.Frame()
	wantFormHoverPixel(t, tt.Image(), int(cx)+10, int(cy), colors.AccentForeground, "thumb cap at rest")

	tt.Press(cx, cy)
	tt.Frame()
	wantFormHoverPixel(t, tt.Image(), x, y, colors.Accent, "thumb pressed cap edge")
	tt.Release(cx, cy)
}

// TestInputGroupHoverPaint checks the native input group's field surface
// follows hover and reverts.
func TestInputGroupHoverPaint(t *testing.T) {
	var colors theme.Palette
	value := "hello"
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).FillWidth().FillHeight().Padding(24).
			Background(colors.Background).Children(func() {
			InputGroup(c, &value, InputGroupOptions{Label: "Name", FullWidth: true})
		})
	}, 460, 180)

	row, ok := tt.Find("Name")
	if !ok {
		t.Fatal("input group missing")
	}
	x, y := int(row.X+row.W-16), int(row.Y+row.H/2)
	wantFormHoverPixel(t, tt.Image(), x, y, colors.FieldBackground, "input group normal fill")

	tt.Move(row.X+row.W/2, row.Y+row.H/2)
	tt.Frame()
	wantFormHoverPixel(t, tt.Image(), x, y, colors.FieldHover, "input group hover fill")

	tt.Move(2, 2)
	tt.Frame()
	wantFormHoverPixel(t, tt.Image(), x, y, colors.FieldBackground, "input group fill after leaving")
}
