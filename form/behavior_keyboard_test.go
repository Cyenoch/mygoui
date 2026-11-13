package form

import (
	"testing"

	"github.com/Cyenoch/mygoui/list"
	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// The radio group is one Tab stop that enters at the chosen radio; Home, End
// and the arrows choose, and the arrows wrap.
func TestRadioGroupTabEntersAtChoiceAndKeysChoose(t *testing.T) {
	choice := "pro"
	state := &RadioGroupState[string]{Selected: &choice}
	tt := formTextTester(t, func(c *ui.Context) {
		var before string
		Input(c, &before, InputOptions{Label: "Before"})
		RadioGroup(c, RadioGroupOptions{Label: "Plan"}, func() {
			Radio(c, RadioOptions[string]{State: state, Value: "free", Label: "Free"})
			Radio(c, RadioOptions[string]{State: state, Value: "pro", Label: "Pro"})
			Radio(c, RadioOptions[string]{State: state, Value: "team", Label: "Team"})
		})
	})
	formTextClick(t, tt, "Before")
	tt.Key(0, ui.KeyTab)
	if !tt.Focused("Pro") {
		t.Fatal("Tab did not enter the group at the chosen radio")
	}
	tt.Key(0, ui.KeyEnd)
	if choice != "team" || !tt.Focused("Team") {
		t.Fatalf("End chose %q", choice)
	}
	tt.Key(0, ui.KeyHome)
	if choice != "free" {
		t.Fatalf("Home chose %q", choice)
	}
	tt.Key(0, ui.KeyUp)
	if choice != "team" {
		t.Fatalf("Up from the first radio chose %q, want the last", choice)
	}
	tt.Key(0, ui.KeyTab)
	if tt.Focused("Team") {
		t.Fatal("Tab stayed inside the group")
	}
}

// A closed single select opens on Down at the first option and on Up at the
// last, and a letter chooses without opening it.
func TestSelectClosedKeysOpenAtEndsAndLettersChoose(t *testing.T) {
	state := &SelectState{}
	items := []list.ListBoxItem{{Key: "a", Label: "Alpha"}, {Key: "b", Label: "Beta"}, {Key: "c", Label: "Charlie"}}
	tt := ui.NewTester(func(c *ui.Context) {
		theme.ApplyMode(c, false)
		var before string
		Input(c, &before, InputOptions{Label: "Before"})
		Select(c, SelectOptions{Label: "Fruit", Placeholder: "Pick", Items: items, State: state})
	}, 360, 480)
	formTextClick(t, tt, "Before")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyB)
	if state.List.SelectedKey != "b" || state.Open {
		t.Fatalf("B chose %q, open %v", state.List.SelectedKey, state.Open)
	}
	tt.Key(0, ui.KeyUp)
	tt.Frame()
	if !state.Open || state.List.HighlightKey() != "c" {
		t.Fatalf("Up: open %v, highlight %q, want the last option", state.Open, state.List.HighlightKey())
	}
}

// Enter on a resting hotkey recorder starts recording, as a click does.
func TestHotkeyInputEnterStartsRecording(t *testing.T) {
	value := Hotkey{Mods: ui.Alt, Key: ui.KeyP}
	changes := 0
	tt := hotkeyTester(&value, &changes)
	formTextClick(t, tt, "Shortcut")
	tt.Key(0, ui.KeyEscape)
	tt.Key(0, ui.KeyEnter)
	tt.Key(ui.Ctrl, ui.KeyJ)
	if want := (Hotkey{Mods: ui.Ctrl, Key: ui.KeyJ}); value != want || changes != 1 {
		t.Fatalf("recorded %+v with %d changes, want %+v", value, changes, want)
	}
}

// A slider steps a tenth of its range with PageUp and PageDown (toolkit).
func TestSliderPageKeysStepATenth(t *testing.T) {
	level, before := 50.0, ""
	tt := formTextTester(t, func(c *ui.Context) {
		Input(c, &before, InputOptions{Label: "Before"})
		Slider(c, SliderOptions{State: SliderState{Value: &level}, Label: "Level", Step: 1})
	})
	formTextClick(t, tt, "Before")
	tt.Key(0, ui.KeyTab)
	tt.Key(0, ui.KeyPageUp)
	if level != 60 {
		t.Fatalf("PageUp moved the slider to %v", level)
	}
	tt.Key(0, ui.KeyPageDown)
	if level != 50 {
		t.Fatalf("PageDown moved the slider to %v", level)
	}
}
