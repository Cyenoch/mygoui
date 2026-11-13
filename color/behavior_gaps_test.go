package color

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/locale"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestColorAreaPagesBrightnessAndMirrorsRightToLeft checks the area's page
// keys and that right to left the saturation axis, keys and pointer included,
// runs from the right.
func TestColorAreaPagesBrightnessAndMirrorsRightToLeft(t *testing.T) {
	s := NewColorState(ui.Hex("#ff0000"))
	s.SetSaturation(0.5)
	s.SetBrightness(0.5)
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ColorArea(c, ColorAreaOptions{Label: "Fill", State: s, Size: 200})
	}, 320, 320)
	tt.SetDirection(ui.RightToLeft)
	tt.Frame()
	area, _ := tt.Find("Fill")
	tt.Key(0, ui.KeyTab)
	for _, st := range []struct {
		key        ui.Key
		sat, value float64
	}{
		{ui.KeyPageUp, 0.5, 0.6},
		{ui.KeyPageDown, 0.5, 0.5},
		{ui.KeyLeft, 0.51, 0.5},
		{ui.KeyRight, 0.5, 0.5},
	} {
		tt.Key(0, st.key)
		if !near(s.Saturation(), st.sat, 0.001) || !near(s.Brightness(), st.value, 0.001) {
			t.Fatalf("after %v: saturation %.3f brightness %.3f, want %.2f %.2f",
				st.key, s.Saturation(), s.Brightness(), st.sat, st.value)
		}
	}
	tt.Press(area.X+area.W*0.25, area.Y+area.H/2)
	tt.Release(area.X+area.W*0.25, area.Y+area.H/2)
	if !near(s.Saturation(), 0.75, 0.02) {
		t.Fatalf("a press a quarter from the left set saturation %.3f, want 0.75 right to left", s.Saturation())
	}
}

// TestColorSliderUpAndDownStepAHorizontalSlider checks that Up raises and
// Down lowers a horizontal slider, as on every slider.
func TestColorSliderUpAndDownStepAHorizontalSlider(t *testing.T) {
	s := NewColorState(ui.Hex("#ff0000"))
	s.SetHue(180)
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ColorSlider(c, ColorSliderOptions{Label: "Hue", State: s, Channel: ChannelHue})
	}, 320, 120)
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyUp)
	if !near(s.Hue(), 181, 0.001) {
		t.Fatalf("Up set the hue to %v, want 181", s.Hue())
	}
	tt.Key(ui.Shift, ui.KeyDown)
	if !near(s.Hue(), 171, 0.001) {
		t.Fatalf("Shift+Down set the hue to %v, want 171", s.Hue())
	}
}

// TestColorSwatchSaysWhenItIsSelected checks that a selected swatch tells
// assistive technology so.
func TestColorSwatchSaysWhenItIsSelected(t *testing.T) {
	selected := true
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ColorSwatch(c, ColorSwatchOptions{Color: ui.Hex("#e0010f"), Label: "Crimson", Selected: selected})
	}, 120, 120)
	if a, _ := tt.AccessByRole(ui.RoleImage, "Crimson"); a.Description != locale.M().Selected {
		t.Fatalf("a selected swatch is described %q", a.Description)
	}
	selected = false
	tt.Frame()
	if a, _ := tt.AccessByRole(ui.RoleImage, "Crimson"); a.Description != "" {
		t.Fatalf("an unselected swatch is described %q", a.Description)
	}
}

// TestColorSwatchPickerChoosesWithoutAValue checks that a picker given no
// Value keeps the choice itself and reports it.
func TestColorSwatchPickerChoosesWithoutAValue(t *testing.T) {
	changes := 0
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ColorSwatchPicker(c, ColorSwatchPickerOptions{
			Label: "Palette", Colors: []ui.Color{ui.Hex("#ff0000"), ui.Hex("#00ff00")},
			Names: []string{"Red", "Green"}, OnChange: func() { changes++ },
		})
	}, 240, 120)
	if err := tt.Click("Green"); err != nil {
		t.Fatal(err)
	}
	if a, _ := tt.AccessByRole(ui.RoleRadio, "Green"); !a.Checked || changes != 1 {
		t.Fatalf("Green: checked %v, changes %d", a.Checked, changes)
	}
}

// TestSwatchHoverGrowsAtOnce checks F3: hover scales the swatch in the frame
// the pointer arrives, and the scale leaves as soon as the pointer does.
func TestSwatchHoverGrowsAtOnce(t *testing.T) {
	var colors theme.Palette
	red := ui.RGB(255, 0, 0)
	var chosen ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		colors = theme.Colors(c)
		ui.Column(c).Fill().Padding(24).Background(colors.Background).Children(func() {
			ColorSwatchPicker(c, ColorSwatchPickerOptions{
				Label: "Palette", Value: &chosen,
				Colors: []ui.Color{red, ui.RGB(0, 255, 0)}, Names: []string{"Red", "Green"},
			})
		})
	}, 400, 140)
	item, _ := tt.Find("Red")
	// One DIP inside the item's transparent border, covered only when the
	// swatch grows.
	x, y := int(item.X)+1, int(item.Y+item.H/2)
	wantColorHoverPixel(t, tt.Image(), x, y, colors.Background, "swatch at rest")
	tt.Move(item.X+item.W/2, item.Y+item.H/2)
	tt.Frame()
	wantColorHoverPixel(t, tt.Image(), x, y, red, "hovered swatch")
	tt.Move(2, 2)
	tt.Frame()
	wantColorHoverPixel(t, tt.Image(), x, y, colors.Background, "swatch after leaving")
}

// TestColorPickerMovesFocusIntoThePanelAndReportsOpening checks that opening
// the panel puts focus on its first control, and that opening and every way
// of closing are reported.
func TestColorPickerMovesFocusIntoThePanelAndReportsOpening(t *testing.T) {
	s := NewColorState(ui.Hex("#ff0000"))
	open := false
	var events []bool
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ColorPicker(c, ColorPickerOptions{Label: "Fill", State: s, Open: &open,
			OnOpenChange: func(v bool) { events = append(events, v) }})
	}, 520, 760)
	if err := tt.Click("Fill"); err != nil {
		t.Fatal(err)
	}
	if !tt.Focused("Color area") {
		t.Fatal("opening did not move focus into the panel")
	}
	tt.Key(0, ui.KeyEscape)
	tt.Advance(300 * time.Millisecond)
	if err := tt.Click("Fill"); err != nil {
		t.Fatal(err)
	}
	tt.Advance(300 * time.Millisecond)
	tt.ClickAt(510, 750)
	tt.Advance(300 * time.Millisecond)
	if open || len(events) != 4 || !events[0] || events[1] || !events[2] || events[3] {
		t.Fatalf("OnOpenChange reported %v, open %v", events, open)
	}
}
