package form

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestCheckboxAccentRevealSettles checks the chosen face scales and fades in,
// and that reduced motion reaches it immediately.
func TestCheckboxAccentRevealSettles(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		name := "motion"
		if reduce {
			name = "reduced motion"
		}
		t.Run(name, func(t *testing.T) {
			var colors theme.Palette
			checked := false
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				colors = theme.Colors(c)
				ui.Column(c).FillWidth().FillHeight().Padding(24).
					Background(colors.Background).Children(func() {
					Checkbox(c, CheckboxOptions{Checked: &checked, Label: "Accept"})
				})
			}, 400, 140)
			if reduce {
				tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			}

			row, ok := tt.Find("Accept")
			if !ok {
				t.Fatal("checkbox missing")
			}
			// Inside the box, below the check mark's path.
			x := int(row.X) + 8
			y := int(row.Y+(row.H-float32(theme.Space16))/2) + 13
			wantFormHoverPixel(t, tt.Image(), x, y, colors.FieldBackground, "unchecked fill")

			if err := tt.Click("Accept"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
			if reduce {
				wantFormHoverPixel(t, tt.Image(), x, y, colors.AccentHover,
					"reduced motion checked fill")
				return
			}
			wantFormHoverPixel(t, tt.Image(), x, y, colors.FieldBackground,
				"entering accent layer")
			tt.Advance(200 * time.Millisecond)
			wantFormHoverPixel(t, tt.Image(), x, y, colors.AccentHover,
				"settled checked fill")
		})
	}
}

// TestRadioDotSelectionSettles checks the dot's color and scale move over
// 200 ms and snap under reduced motion.
func TestRadioDotSelectionSettles(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		name := "motion"
		if reduce {
			name = "reduced motion"
		}
		t.Run(name, func(t *testing.T) {
			var colors theme.Palette
			selected := "b"
			state := RadioGroupState[string]{Selected: &selected}
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				colors = theme.Colors(c)
				ui.Column(c).FillWidth().FillHeight().Padding(24).
					Background(colors.Background).Children(func() {
					RadioGroup(c, RadioGroupOptions{Label: "Group"}, func() {
						Radio(c, RadioOptions[string]{
							State: &state, Value: "a", Label: "Alpha",
						})
					})
				})
			}, 400, 200)
			if reduce {
				tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			}

			row, ok := tt.Find("Alpha")
			if !ok {
				t.Fatal("radio missing")
			}
			// Outside the 6 DIP chosen dot but inside the 14 DIP resting one.
			x, y := int(row.X)+13, int(row.Y+row.H/2)
			wantFormHoverPixel(t, tt.Image(), x, y, colors.FieldBackground,
				"unselected dot")

			if err := tt.Click("Alpha"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
			if reduce {
				wantFormHoverPixel(t, tt.Image(), x, y, colors.Accent,
					"reduced motion dot")
				return
			}
			wantFormHoverPixel(t, tt.Image(), x, y, colors.FieldBackground,
				"dot before it shrinks")
			tt.Advance(250 * time.Millisecond)
			wantFormHoverPixel(t, tt.Image(), x, y, colors.Accent,
				"settled dot")
		})
	}
}

// TestSwitchThumbSlidesAndSettles checks the 300 ms thumb travel reaches the
// far end and reduced motion arrives there at once.
func TestSwitchThumbSlidesAndSettles(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		name := "motion"
		if reduce {
			name = "reduced motion"
		}
		t.Run(name, func(t *testing.T) {
			var colors theme.Palette
			on := false
			tt := ui.NewTester(func(c *ui.Context) {
				theme.ApplyMode(c, false)
				colors = theme.Colors(c)
				ui.Column(c).FillWidth().FillHeight().Padding(24).
					Background(colors.Background).Children(func() {
					Switch(c, SwitchOptions{On: &on, Label: "Wifi"})
				})
			}, 400, 140)
			if reduce {
				tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
			}

			row, ok := tt.Find("Wifi")
			if !ok {
				t.Fatal("switch missing")
			}
			// Inside the thumb's resting place at the track's start.
			x, y := int(row.X)+5, int(row.Y+row.H/2)
			wantFormHoverPixel(t, tt.Image(), x, y, ui.RGB(255, 255, 255),
				"thumb at rest")

			if err := tt.Click("Wifi"); err != nil {
				t.Fatal(err)
			}
			tt.Frame()
			if reduce {
				wantFormHoverPixel(t, tt.Image(), x, y, colors.AccentHover,
					"reduced motion track")
				return
			}
			wantFormHoverPixel(t, tt.Image(), x, y, ui.RGB(255, 255, 255),
				"thumb before it slides")
			tt.Advance(400 * time.Millisecond)
			wantFormHoverPixel(t, tt.Image(), x, y, colors.AccentHover,
				"settled track under the thumb's old place")
		})
	}
}

// TestSliderThumbPressSettles checks the cap's 250 ms shrink exposes the
// accent pill, without a one-frame snap.
func TestSliderThumbPressSettles(t *testing.T) {
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
	// The inner cap is 24×16 at rest; sample its centre and the right band
	// that only the resting 24 DIP width covers.
	centerX, centerY := int(cx), int(cy)
	edgeX, edgeY := int(cx+11), int(cy)
	wantFormHoverPixel(t, tt.Image(), centerX, centerY, colors.AccentForeground,
		"cap centre before the press")
	_, beforeG, _ := formHoverPixel(tt.Image(), edgeX, edgeY)

	tt.Press(cx, cy)
	tt.Frame()
	// The cap face stays neutral and does not snap to the accent face.
	wantFormHoverPixel(t, tt.Image(), centerX, centerY, colors.AccentForeground,
		"cap centre before the press transition")
	_, startG, _ := formHoverPixel(tt.Image(), edgeX, edgeY)
	if !formHoverNear(startG, beforeG) {
		t.Fatalf("the cap moved before its transition: edge green %d, was %d", startG, beforeG)
	}

	tt.Advance(300 * time.Millisecond)
	wantFormHoverPixel(t, tt.Image(), centerX, centerY, colors.AccentForeground,
		"cap centre after the press transition")
	_, pressedG, _ := formHoverPixel(tt.Image(), edgeX, edgeY)
	if pressedG >= beforeG {
		t.Fatalf("the cap did not shrink toward accent: edge green %d, was %d", pressedG, beforeG)
	}
	tt.Advance(100 * time.Millisecond)
	_, settledG, _ := formHoverPixel(tt.Image(), edgeX, edgeY)
	if !formHoverNear(settledG, pressedG) {
		t.Fatalf("the cap did not settle: green %d after 400 ms, %d after 300 ms", settledG, pressedG)
	}
	tt.Release(cx, cy)
}

// TestSelectKeyboardCommitReturnsFocus checks the popup takes the keyboard
// after the trigger opens it and returns focus to the trigger on a choice.
func TestSelectKeyboardCommitReturnsFocus(t *testing.T) {
	state := &SelectState{}
	items := []list.ListBoxItem{
		{Key: "a", Label: "Alpha"},
		{Key: "b", Label: "Beta"},
	}
	focused := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		Select(c, SelectOptions{
			Label: "Fruit", Placeholder: "Pick a fruit",
			Items: items, State: state,
		})
		focused = state.Trigger.Focused(c)
	}, 360, 480)

	if err := tt.Click("Pick a fruit"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if !state.Open {
		t.Fatal("the trigger did not open the popup")
	}
	tt.Key(0, ui.KeyEnter)
	tt.Frame()
	if state.List.SelectedKey != "a" || state.Open {
		t.Fatalf("keyboard commit: key %q, open %v", state.List.SelectedKey, state.Open)
	}
	tt.Frame()
	if !focused {
		t.Fatal("focus did not return to the trigger after the choice")
	}
}

// TestSelectDisabledDoesNotOpen checks that a disabled trigger ignores presses.
func TestSelectDisabledDoesNotOpen(t *testing.T) {
	state := &SelectState{}
	items := []list.ListBoxItem{{Key: "a", Label: "Alpha"}}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		Select(c, SelectOptions{
			Label: "Fruit", Placeholder: "Pick a fruit",
			Items: items, State: state, Disabled: true,
		})
	}, 360, 480)
	if err := tt.Click("Pick a fruit"); err != nil {
		t.Fatal(err)
	}
	tt.Frame()
	if state.Open {
		t.Fatal("a disabled select opened its popup")
	}
}

// TestComboBoxFilterFoldsCaseAndNewOptions checks the default filter ignores
// case, handles a non-ASCII label, and picks up replaced options.
func TestComboBoxFilterFoldsCaseAndNewOptions(t *testing.T) {
	state := &ComboBoxState{}
	items := []list.ListBoxItem{
		{Key: "a", Label: "Alpha"},
		{Key: "b", Label: "Beta"},
		{Key: "c", Label: "Bogotá"},
	}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		ComboBox(c, ComboBoxOptions{
			Label: "Language", Placeholder: "Pick a language",
			Items: items, State: state,
		})
	}, 360, 520)

	if err := tt.Click("Language"); err != nil {
		t.Fatal(err)
	}
	tt.Type("al")
	tt.Frame()
	if !tt.HasText("Alpha") {
		t.Fatalf("a lowercase query did not match its label: %q", tt.Texts())
	}
	if tt.HasText("Beta") {
		t.Fatalf("a non-matching option stayed in the list: %q", tt.Texts())
	}

	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("BOG")
	tt.Frame()
	if !tt.HasText("Bogotá") {
		t.Fatalf("a non-ASCII label did not match BOG: %q", tt.Texts())
	}

	tt.Key(ui.Cmd, ui.KeyA)
	tt.Type("alp")
	tt.Frame()
	if !tt.HasText("Alpha") {
		t.Fatalf("the old options did not match alp: %q", tt.Texts())
	}

	items = []list.ListBoxItem{
		{Key: "d", Label: "Alpaca"},
		{Key: "e", Label: "Beaver"},
		{Key: "f", Label: "Bison"},
	}
	tt.Frame()
	if !tt.HasText("Alpaca") {
		t.Fatalf("replaced options were not filtered: %q", tt.Texts())
	}
	if tt.HasText("Alpha") {
		t.Fatalf("the old filtered result survived the option change: %q", tt.Texts())
	}
}
