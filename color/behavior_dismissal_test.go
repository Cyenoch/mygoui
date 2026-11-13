package color

import (
	"testing"

	"github.com/Cyenoch/mygoui/theme"
	"github.com/egoist/mygo/ui"
)

// TestColorPickerDismissalReturnsToTrigger checks the picker's panel:
// opening it moves focus in, and Escape closes it with focus back on the
// trigger, which reports whether the panel shows.
func TestColorPickerDismissalReturnsToTrigger(t *testing.T) {
	s := NewColorState(ui.Hex("#ff0000"))
	open := false
	tt := ui.NewTester(func(c *ui.Context) {
		theme.Apply(c)
		ColorPicker(c, ColorPickerOptions{Label: "Fill", State: s, Open: &open})
	}, 520, 760)
	expanded := func() bool {
		a, ok := tt.AccessByRole(ui.RoleColorWell, "Fill")
		if !ok {
			t.Fatal("the trigger is not a color well named Fill")
		}
		return a.Expanded
	}

	if err := tt.Click("Fill"); err != nil {
		t.Fatal(err)
	}
	if !open || !expanded() {
		t.Fatal("the trigger did not open the panel")
	}
	if !tt.Focused("Color area") {
		t.Fatal("opening did not move focus into the panel")
	}
	tt.Key(0, ui.KeyEscape)
	tt.Frame()
	if open || expanded() || !tt.Focused("Fill") {
		t.Fatalf("Escape: open %v, focus on trigger %v", open, tt.Focused("Fill"))
	}
	if s.Value != ui.Hex("#ff0000") {
		t.Fatalf("dismissing changed the color to %v", s.Value)
	}
}
