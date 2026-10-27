package form

import (
	"testing"
	"time"

	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// Checked choice controls use the primary face immediately, including when
// reduced motion is enabled.
func TestChoiceFacesUsePrimary(t *testing.T) {
	for _, reduce := range []bool{false, true} {
		checked := false
		var colors theme.Palette
		tt := ui.NewTester(func(c *ui.Context) {
			theme.ApplyMode(c, false)
			colors = theme.Colors(c)
			ui.Column(c).Fill().Padding(24).Background(colors.Background).Children(func() {
				Checkbox(c, CheckboxOptions{Checked: &checked, Label: "Accept"})
			})
		}, 400, 140)
		if reduce {
			tt.SetPreferences(ui.Preferences{ReduceMotion: true, TextScale: 1})
		}
		row, _ := tt.Find("Accept")
		if err := tt.Click("Accept"); err != nil {
			t.Fatal(err)
		}
		tt.Frame()
		x, y := int(row.X)+8, int(row.Y+(row.H-16)/2)+13
		wantFormHoverPixel(t, tt.Image(), x, y, colors.Primary, "checked face")
	}
}

// TestSwitchThumbSlidesAndSettles checks the thumb travel reaches the
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
				wantFormHoverPixel(t, tt.Image(), x, y, colors.Primary,
					"reduced motion track")
				return
			}
			wantFormHoverPixel(t, tt.Image(), x, y, ui.RGB(255, 255, 255),
				"thumb before it slides")
			tt.Advance(400 * time.Millisecond)
			wantFormHoverPixel(t, tt.Image(), x, y, colors.Primary,
				"settled track under the thumb's old place")
		})
	}
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
