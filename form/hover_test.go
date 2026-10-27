package form

import (
	"image"
	"testing"

	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// formHoverPixel returns the colour at x, y in img.
func formHoverPixel(img *image.RGBA, x, y int) (uint8, uint8, uint8) {
	c := img.RGBAAt(x, y)
	return c.R, c.G, c.B
}

func TestEditablePopupHighlightsOnlyWithEditorFocus(t *testing.T) {
	items := []list.ListBoxItem{{Key: "a", Label: "Alpha"}, {Key: "b", Label: "Beta"}}
	for _, kind := range []string{"autocomplete", "combobox"} {
		t.Run(kind, func(t *testing.T) {
			auto := AutocompleteState{Open: true}
			combo := ComboBoxState{Open: true}
			var colors theme.Palette
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				colors = theme.Colors(c)
				ui.Column(c).Fill().Padding(16).Background(colors.Background).Children(func() {
					if kind == "autocomplete" {
						Autocomplete(c, AutocompleteOptions{Label: "Fruit", Items: items, State: &auto})
					} else {
						ComboBox(c, ComboBoxOptions{Label: "Fruit", Items: items, State: &combo})
					}
				})
			}, 360, 480)
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			alpha, ok := tt.Find("Alpha")
			if !ok {
				t.Fatal("popup has no first option")
			}
			x, y := int(alpha.X+4), int(alpha.Y+alpha.H/2)
			wantFormHoverPixel(t, tt.Image(), x, y, colors.Popover, "unfocused option")
			if err := tt.Click("Fruit"); err != nil {
				t.Fatal(err)
			}
			tt.Key(0, ui.KeyDown)
			beta, ok := tt.Find("Beta")
			if !ok {
				t.Fatal("popup has no second option")
			}
			wantFormHoverPixel(t, tt.Image(), int(beta.X+4), int(beta.Y+beta.H/2),
				colors.Accent.Over(colors.Popover), "keyboard highlight")
		})
	}
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
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})

	row, ok := tt.Find("Accept")
	if !ok {
		t.Fatal("checkbox missing")
	}
	// The 16 DIP box sits at the row's start; sample inside it, below the
	// check mark's path.
	x := int(row.X) + 8
	y := int(row.Y+(row.H-float32(theme.Space16))/2) + 13
	wantFormHoverPixel(t, tt.Image(), x, y, colors.Primary, "checkbox normal fill")

	tt.Move(row.X+row.W/2, row.Y+row.H/2)
	tt.Frame()
	wantFormHoverPixel(t, tt.Image(), x, y, colors.Primary, "checkbox hover fill")

	tt.Press(row.X+row.W/2, row.Y+row.H/2)
	tt.Frame()
	wantFormHoverPixel(t, tt.Image(), x, y, colors.Primary, "checkbox pressed fill")
	tt.Release(row.X+row.W/2, row.Y+row.H/2)
	// Releasing activates the box and clears the caller's value; restore it
	// so leaving still checks the checked face.
	checked = true
	tt.Frame()

	tt.Move(2, 2)
	tt.Frame()
	wantFormHoverPixel(t, tt.Image(), x, y, colors.Primary, "checkbox fill after leaving")
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
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})

		row, ok := tt.Find("Alpha")
		if !ok {
			t.Fatal("radio missing")
		}
		x, y := int(row.X)+8, int(row.Y+row.H/2)
		wantFormHoverPixel(t, tt.Image(), x, y, colors.FieldBackground, "radio normal fill")

		tt.Move(row.X+row.W/2, row.Y+row.H/2)
		tt.Frame()
		wantFormHoverPixel(t, tt.Image(), x, y, colors.Background, "radio hover fill")

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
		tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})

		row, ok := tt.Find("Alpha")
		if !ok {
			t.Fatal("radio missing")
		}
		// The selected face is an accent ring around the white dot; sample
		// the ring, clear of the dot's pressed growth.
		x, y := int(row.X)+2, int(row.Y+row.H/2)
		wantFormHoverPixel(t, tt.Image(), x, y, colors.Primary, "radio selected fill")

		tt.Press(row.X+row.W/2, row.Y+row.H/2)
		tt.Frame()
		wantFormHoverPixel(t, tt.Image(), x, y, colors.Primary, "radio pressed fill")
		tt.Release(row.X+row.W/2, row.Y+row.H/2)
		tt.Frame()
		wantFormHoverPixel(t, tt.Image(), x, y, colors.Primary, "radio fill after release")
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
			normal: func(p theme.Palette) ui.Color { return p.Primary },
			hover:  func(p theme.Palette) ui.Color { return p.Primary },
		},
		{
			name: "off", on: false, xOffset: 24,
			normal: func(p theme.Palette) ui.Color { return p.Input.Over(p.Background) },
			hover:  func(p theme.Palette) ui.Color { return p.Input.Over(p.Background) },
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
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})

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

// TestSliderThumbUsesWhiteFace checks the Coss circular thumb over the primary track.
func TestSliderThumbUsesWhiteFace(t *testing.T) {
	value := 50.0
	var colors theme.Palette
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).Fill().Padding(24).Background(colors.Background).Children(func() {
			Slider(c, SliderOptions{State: SliderState{Value: &value}, Label: "Volume", HideOutput: true})
		})
	}, 400, 200)
	header, _ := tt.Find("Volume")
	cx := header.X + (400-2*header.X)/2
	cy := header.Y + header.H + theme.Space4 + 10
	wantFormHoverPixel(t, tt.Image(), int(cx), int(cy), ui.RGB(255, 255, 255), "thumb face")
	wantFormHoverPixel(t, tt.Image(), int(cx)-25, int(cy), colors.Primary, "slider active track")
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
	tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})

	row, ok := tt.Find("Name")
	if !ok {
		t.Fatal("input group missing")
	}
	x, y := int(row.X+row.W-16), int(row.Y+row.H/2)
	wantFormHoverPixel(t, tt.Image(), x, y, colors.FieldBackground, "input group normal fill")

	tt.Move(row.X+row.W/2, row.Y+row.H/2)
	tt.Frame()
	wantFormHoverPixel(t, tt.Image(), x, y, colors.Background, "input group hover fill")

	tt.Move(2, 2)
	tt.Frame()
	wantFormHoverPixel(t, tt.Image(), x, y, colors.FieldBackground, "input group fill after leaving")
}
